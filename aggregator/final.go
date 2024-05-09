package aggregator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// tryBuildFinalProof checks if the provided proof is eligible to be used to
// build the final proof.  If no proof is provided it looks for a previously
// generated proof.  If the proof is eligible, then the final proof generation
// is triggered.
func (a *Aggregator) tryBuildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) (bool, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
	)
	log.Debug("tryBuildFinalProof start")

	var err error
	if !a.canVerifyProof() {
		log.Debug("Time to verify proof not reached or proof verification in progress")
		return false, nil
	}
	log.Debug("Send final proof time reached")

	for !a.isSynced(ctx, nil) {
		log.Info("Waiting for synchronizer to sync...")
		time.Sleep(a.cfg.RetryTime.Duration)
		continue
	}

	var lastVerifiedBatchNum uint64
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return false, fmt.Errorf("failed to get last verified batch, %w", err)
	}
	if lastVerifiedBatch != nil {
		lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
	}

	if proof == nil {
		// we don't have a proof generating at the moment, check if we
		// have a proof ready to verify

		proof, err = a.getAndLockProofReadyForFinal(ctx, prover, lastVerifiedBatchNum)
		if errors.Is(err, state.ErrNotFound) {
			// nothing to verify, swallow the error
			log.Debug("No proof ready to verify")
			return false, nil
		}
		if err != nil {
			return false, err
		}

		defer func() {
			if err != nil {
				// Set the generating state to false for the proof ("unlock" it)
				proof.GeneratingSince = nil
				err2 := a.State.UpdateBatchProof(a.ctx, proof, nil)
				if err2 != nil {
					log.Errorf("Failed to unlock proof: %v", err2)
				}
			}
		}()
	} else {
		// we do have a proof generating at the moment, check if it is
		// eligible to be verified
		eligible, err := a.validateEligibleFinalProof(ctx, proof, lastVerifiedBatchNum)
		if err != nil {
			return false, fmt.Errorf("failed to validate eligible final proof, %w", err)
		}
		if !eligible {
			return false, nil
		}
	}

	log = log.WithFields(
		"proofId", *proof.ProofID,
		"batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal),
	)

	// at this point we have an eligible proof, build the final one using it
	finalProof, err := a.buildFinalProof(ctx, prover, proof)
	if err != nil {
		err = fmt.Errorf("failed to build final proof, %w", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	msg := finalProofMsg{
		proverName:     proverName,
		proverID:       proverID,
		recursiveProof: proof,
		finalProof:     finalProof,
	}

	select {
	case <-a.ctx.Done():
		return false, a.ctx.Err()
	case a.finalProof <- msg:
	}

	log.Debug("tryBuildFinalProof end")
	return true, nil
}

// buildFinalProof builds and return the final proof for an aggregated/batch proof.
func (a *Aggregator) buildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) (*prover.FinalProof, error) {
	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
		"recursiveProofId", *proof.ProofID,
		"batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal),
	)
	log.Info("Generating final proof")

	finalProofID, err := prover.FinalProof(proof.Proof, a.cfg.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get final proof id: %w", err)
	}
	proof.ProofID = finalProofID

	log.Infof("Final proof ID for batches [%d-%d]: %s", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)
	log = log.WithFields("finalProofId", finalProofID)

	finalProof, err := prover.WaitFinalProof(ctx, *proof.ProofID)
	if err != nil {
		return nil, fmt.Errorf("failed to get final proof from prover: %w", err)
	}

	log.Info("Final proof generated")

	// mock prover sanity check
	if string(finalProof.Public.NewStateRoot) == mockedStateRoot && string(finalProof.Public.NewLocalExitRoot) == mockedLocalExitRoot {
		// This local exit root and state root come from the mock
		// prover, use the one captured by the executor instead
		finalBatch, err := a.State.GetBatchByNumber(ctx, proof.BatchNumberFinal, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve batch with number [%d]", proof.BatchNumberFinal)
		}
		log.Warnf("NewLocalExitRoot and NewStateRoot look like a mock values, using values from executor instead: LER: %v, SR: %v",
			finalBatch.LocalExitRoot.TerminalString(), finalBatch.StateRoot.TerminalString())
		finalProof.Public.NewStateRoot = finalBatch.StateRoot.Bytes()
		finalProof.Public.NewLocalExitRoot = finalBatch.LocalExitRoot.Bytes()
	}

	return finalProof, nil
}

func (a *Aggregator) getAndLockProofReadyForFinal(ctx context.Context, prover proverInterface, lastVerifiedBatchNum uint64) (*state.Proof, error) {
	a.StateDBMutex.Lock()
	defer a.StateDBMutex.Unlock()

	// Get proof ready to be verified
	proofToVerify, err := a.State.GetProofReadyForFinal(ctx, lastVerifiedBatchNum, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().Round(time.Microsecond)
	proofToVerify.GeneratingSince = &now

	err = a.State.UpdateBatchProof(ctx, proofToVerify, nil)
	if err != nil {
		return nil, err
	}

	return proofToVerify, nil
}

func (a *Aggregator) validateEligibleFinalProof(ctx context.Context, proof *state.Proof, lastVerifiedBatchNum uint64) (bool, error) {
	batchNumberToVerify := lastVerifiedBatchNum + 1

	if proof.BatchNumber != batchNumberToVerify {
		if proof.BatchNumber < batchNumberToVerify && proof.BatchNumberFinal >= batchNumberToVerify {
			// We have a proof that contains some batches below the last batch verified, anyway can be eligible as final proof
			log.Warnf("Proof %d-%d contains some batches lower than last batch verified %d. Check anyway if it is eligible", proof.BatchNumber, proof.BatchNumberFinal, lastVerifiedBatchNum)
		} else if proof.BatchNumberFinal < batchNumberToVerify {
			// We have a proof that contains batches below that the last batch verified, we need to delete this proof
			log.Warnf("Proof %d-%d lower than next batch to verify %d. Deleting it", proof.BatchNumber, proof.BatchNumberFinal, batchNumberToVerify)
			err := a.State.DeleteBatchProofs(ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
			if err != nil {
				return false, fmt.Errorf("failed to delete discarded proof, err: %w", err)
			}
			return false, nil
		} else {
			log.Debugf("Proof batch number %d is not the following to last verfied batch number %d", proof.BatchNumber, lastVerifiedBatchNum)
			return false, nil
		}
	}

	bComplete, err := a.State.CheckProofContainsCompleteSequences(ctx, proof, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check if proof contains complete sequences, %w", err)
	}
	if !bComplete {
		log.Infof("Recursive proof %d-%d not eligible to be verified: not containing complete sequences", proof.BatchNumber, proof.BatchNumberFinal)
		return false, nil
	}
	return true, nil
}

// This function waits to receive a final proof from a prover. Once it receives
// the proof, it performs these steps in order:
// - send the final proof to L1
// - wait for the synchronizer to catch up
// - clean up the cache of recursive proofs
func (a *Aggregator) sendFinalProof() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case msg := <-a.finalProof:
			ctx := a.ctx
			proof := msg.recursiveProof

			log.WithFields("proofId", proof.ProofID, "batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal))
			log.Info("Verifying final proof with ethereum smart contract")

			a.startProofVerification()

			finalBatch, err := a.State.GetBatchByNumber(ctx, proof.BatchNumberFinal, nil)
			if err != nil {
				log.Errorf("Failed to retrieve batch with number [%d]: %v", proof.BatchNumberFinal, err)
				a.endProofVerification()
				continue
			}

			inputs := ethmanTypes.FinalProofInputs{
				FinalProof:       msg.finalProof,
				NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
				NewStateRoot:     finalBatch.StateRoot.Bytes(),
			}

			log.Infof("Final proof inputs: NewLocalExitRoot [%#x], NewStateRoot [%#x]", inputs.NewLocalExitRoot, inputs.NewStateRoot)

			// add batch verification to be monitored
			sender := common.HexToAddress(a.cfg.SenderAddress)
			to, data, err := a.Ethman.BuildTrustedVerifyBatchesTxData(proof.BatchNumber-1, proof.BatchNumberFinal, &inputs, sender)
			if err != nil {
				log.Errorf("Error estimating batch verification to add to eth tx manager: %v", err)
				a.handleErrorSendFinalProof(ctx, proof)
				continue
			}
			monitoredTxID := buildMonitoredTxID(proof.BatchNumber, proof.BatchNumberFinal)
			err = a.EthTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, sender, to, nil, data, a.cfg.GasOffset, nil)
			if err != nil {
				mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, sender, to)
				mTxLogger.Errorf("Error to add batch verification tx to eth tx manager: %v", err)
				a.handleErrorSendFinalProof(ctx, proof)
				continue
			}

			// process monitored batch verifications before starting a next cycle
			a.EthTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
				a.handleMonitoredTxResult(result)
			}, nil)

			a.resetVerifyProofTime()
			a.endProofVerification()
		}
	}
}

func (a *Aggregator) handleErrorSendFinalProof(ctx context.Context, proof *state.Proof) {
	log := log.WithFields("proofId", proof.ProofID, "batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal))
	proof.GeneratingSince = nil
	err := a.State.UpdateBatchProof(ctx, proof, nil)
	if err != nil {
		log.Errorf("Failed updating proof state (false): %v", err)
	}
	a.endProofVerification()
}
