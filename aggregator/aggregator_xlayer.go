package aggregator

import (
	"context"
	"fmt"
	"strings"
	"time"

	agglayerTypes "github.com/0xPolygon/agglayer/rpc/types"
	"github.com/0xPolygon/agglayer/tx"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

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
