package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	agglayerTypes "github.com/0xPolygon/agglayer/rpc/types"
	"github.com/0xPolygon/agglayer/tx"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/peer"
)

const minParaCount = 2

func (a *Aggregator) settleDirect(
	ctx context.Context,
	proof *state.Proof,
	inputs ethmanTypes.FinalProofInputs,
) (success bool) {
	// add batch verification to be monitored
	sender := common.HexToAddress(a.cfg.SenderAddress)

	to, data, err := a.Ethman.BuildTrustedVerifyBatchesTxData(
		proof.BatchNumber-1,
		proof.BatchNumberFinal,
		&inputs,
		sender,
	)
	if err != nil {
		log.Errorf("Error estimating batch verification to add to eth tx manager: %v", err)
		a.handleFailureToAddVerifyBatchToBeMonitored(ctx, proof)

		return false
	}

	monitoredTxID := buildMonitoredTxID(proof.BatchNumber, proof.BatchNumberFinal)
	err = a.EthTxManager.Add(
		ctx,
		ethTxManagerOwner,
		monitoredTxID,
		sender,
		to,
		nil,
		data,
		a.cfg.GasOffset,
		nil,
	)
	if err != nil {
		mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, sender, to)
		mTxLogger.Errorf("Error to add batch verification tx to eth tx manager: %v", err)
		a.handleFailureToAddVerifyBatchToBeMonitored(ctx, proof)

		return false
	}

	// process monitored batch verifications before starting a next cycle
	a.EthTxManager.ProcessPendingMonitoredTxs(
		ctx,
		ethTxManagerOwner,
		func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
			a.handleMonitoredTxResult(result)
		},
		nil,
	)

	return true
}

func (a *Aggregator) settleWithAggLayer(
	ctx context.Context,
	proof *state.Proof,
	inputs ethmanTypes.FinalProofInputs,
) (success bool) {
	proofStrNo0x := strings.TrimPrefix(inputs.FinalProof.Proof, "0x")
	proofBytes := common.Hex2Bytes(proofStrNo0x)
	tx := tx.Tx{
		LastVerifiedBatch: agglayerTypes.ArgUint64(proof.BatchNumber - 1),
		NewVerifiedBatch:  agglayerTypes.ArgUint64(proof.BatchNumberFinal),
		ZKP: tx.ZKP{
			NewStateRoot:     common.BytesToHash(inputs.NewStateRoot),
			NewLocalExitRoot: common.BytesToHash(inputs.NewLocalExitRoot),
			Proof:            agglayerTypes.ArgBytes(proofBytes),
		},
		RollupID: a.Ethman.GetRollupId(),
	}
	signedTx, err := tx.Sign(a.sequencerPrivateKey)

	if err != nil {
		log.Errorf("failed to sign tx: %v", err)
		a.handleFailureToSendToAggLayer(ctx, proof)

		return false
	}

	log.Debug("final proof signedTx: ", signedTx.Tx.ZKP.Proof.Hex())
	txHash, err := a.AggLayerClient.SendTx(*signedTx)
	if err != nil {
		log.Errorf("failed to send tx to the interop: %v", err)
		a.handleFailureToSendToAggLayer(ctx, proof)

		return false
	}

	log.Infof("tx %s sent to agglayer, waiting to be mined", txHash.Hex())
	log.Debugf("Timeout set to %f seconds", a.cfg.AggLayerTxTimeout.Duration.Seconds())
	waitCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(a.cfg.AggLayerTxTimeout.Duration))
	defer cancelFunc()
	if err := a.AggLayerClient.WaitTxToBeMined(txHash, waitCtx); err != nil {
		log.Errorf("interop didn't mine the tx: %v", err)
		a.handleFailureToSendToAggLayer(ctx, proof)

		return false
	}

	// TODO: wait for synchronizer to catch up
	return true
}

func (a *Aggregator) handleFailureToSendToAggLayer(ctx context.Context, proof *state.Proof) {
	log := log.WithFields("proofId", proof.ProofID, "batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal))
	proof.GeneratingSince = nil

	err := a.State.UpdateGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Errorf("Failed updating proof state (false): %v", err)
	}

	a.endProofVerification()
}

func (a *Aggregator) channelParallel(stream prover.AggregatorService_ChannelServer) error {
	metrics.ConnectedProver()
	defer metrics.DisconnectedProver()

	ctx := stream.Context()
	var proverAddr net.Addr
	p, ok := peer.FromContext(ctx)
	if ok {
		proverAddr = p.Addr
	}
	prover, err := prover.New(stream, proverAddr, a.cfg.ProofStatePollingInterval)
	if err != nil {
		return err
	}

	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)
	log.Info("Establishing stream connection with prover")

	// Check if prover supports the required Fork ID
	if !prover.SupportsForkID(forkId9) {
		err := errors.New("prover does not support required fork ID")
		log.Warn(FirstToUpper(err.Error()))
		return err
	}

	// We start multi batch proof routines, one aggregate proof routine and one final proof routine in parallel.
	paraCount := a.cfg.ParaCount
	if paraCount < minParaCount {
		paraCount = minParaCount
	}
	for i := uint64(0); i < paraCount; i++ {
		go a.generateBatchProofRoutine(ctx, prover)
	}
	go a.aggregateProofsRoutine(ctx, prover)
	go a.buildFinalProofRoutine(ctx, prover)

	select {
	case <-a.ctx.Done():
		// server disconnected
		return a.ctx.Err()
	case <-ctx.Done():
		// client disconnected
		return ctx.Err()
	}
}

func (a *Aggregator) aggregateProofsRoutine(ctx context.Context, prover *prover.Prover) {
	for {
		select {
		case <-a.ctx.Done():
			// server disconnected
			return
		case <-ctx.Done():
			// client disconnected
			return

		default:
			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}

			proofGenerated, err := a.tryAggregateProofs(ctx, prover)
			if err != nil {
				log.Errorf("Error trying to aggregate proofs: %v", err)
			}

			if !proofGenerated {
				// if no proof was generated (aggregated or batch) wait some time before retry
				time.Sleep(a.cfg.RetryTime.Duration)
			} // if proof was generated we retry immediately as probably we have more proofs to process
		}
	}
}

func (a *Aggregator) generateBatchProofRoutine(ctx context.Context, prover *prover.Prover) {
	for {
		select {
		case <-a.ctx.Done():
			// server disconnected
			return
		case <-ctx.Done():
			// client disconnected
			return

		default:
			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}

			proofGenerated, err := a.tryGenerateBatchProof(ctx, prover)
			if err != nil {
				log.Errorf("Error trying to generate proof: %v", err)
			}
			if !proofGenerated {
				// if no proof was generated (aggregated or batch) wait some time before retry
				time.Sleep(a.cfg.RetryTime.Duration)
			} // if proof was generated we retry immediately as probably we have more proofs to process
		}
	}
}

func (a *Aggregator) buildFinalProofRoutine(ctx context.Context, prover *prover.Prover) {
	for {
		select {
		case <-a.ctx.Done():
			// server disconnected
			return
		case <-ctx.Done():
			// client disconnected
			return

		default:
			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}

			proofGenerated, err := a.tryBuildFinalProof(ctx, prover, nil)
			if err != nil {
				log.Errorf("Error checking proofs to verify: %v", err)
			}

			if !proofGenerated {
				// if no proof was generated (aggregated or batch) wait some time before retry
				time.Sleep(a.cfg.RetryTime.Duration)
			} // if proof was generated we retry immediately as probably we have more proofs to process
		}
	}
}
