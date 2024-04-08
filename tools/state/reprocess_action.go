package main

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type reprocessAction struct {
	firstBatchNumber uint64
	lastBatchNumber  uint64
	l2ChainId        uint64
	// If true, when execute a batch write the MT in hashDB
	updateHasbDB             bool
	stopOnError              bool
	preferExecutionStateRoot bool

	st          *state.State
	ctx         context.Context
	output      reprocessingOutputer
	flushIdCtrl flushIDController
}

func (r *reprocessAction) start() error {
	lastBatch := r.lastBatchNumber
	firstBatchNumber := r.firstBatchNumber

	batch, err := getBatchByNumber(r.ctx, r.st, firstBatchNumber-1)
	if err != nil {
		log.Errorf("no batch %d. Error: %v", 0, err)
		return err
	}
	oldStateRoot := batch.StateRoot
	oldAccInputHash := batch.AccInputHash

	for i := firstBatchNumber; i < lastBatch; i++ {
		r.output.startProcessingBatch(i)
		batchOnDB, response, err := r.stepWithFlushId(i, oldStateRoot, oldAccInputHash)
		if response != nil {
			r.output.finishProcessingBatch(response.NewStateRoot, err)
		} else {
			r.output.finishProcessingBatch(common.Hash{}, err)
		}
		if batchOnDB != nil {
			oldStateRoot = batchOnDB.StateRoot
			oldAccInputHash = batchOnDB.AccInputHash
		}

		if r.preferExecutionStateRoot && response != nil {
			// If there is a response use that instead of the batch on DB
			log.Infof("Using as oldStateRoot the execution state root: %s", response.NewStateRoot)
			oldStateRoot = response.NewStateRoot
			oldAccInputHash = response.NewAccInputHash
		}
		if r.stopOnError && err != nil {
			log.Fatalf("error processing batch %d. Error: %v", i, err)
		}
	}
	return nil
}

func (r *reprocessAction) stepWithFlushId(i uint64, oldStateRoot common.Hash, oldAccInputHash common.Hash) (*state.Batch, *state.ProcessBatchResponse, error) {
	if r.updateHasbDB {
		err := r.flushIdCtrl.BlockUntilLastFlushIDIsWritten()
		if err != nil {
			return nil, nil, err
		}
	}
	batchOnDB, response, err := r.step(i, oldStateRoot, oldAccInputHash)
	if r.updateHasbDB && err == nil && response != nil {
		r.flushIdCtrl.SetPendingFlushIDAndCheckProverID(response.FlushID, response.ProverID, "reprocessAction")
	}
	return batchOnDB, response, err
}

// returns:
// - state.Batch -> batch on DB
// - *ProcessBatchResponse -> response of reprocessing batch with EXECUTOR
func (r *reprocessAction) step(i uint64, oldStateRoot common.Hash, oldAccInputHash common.Hash) (*state.Batch, *state.ProcessBatchResponse, error) {
	dbTx, err := r.st.BeginStateTransaction(r.ctx)
	if err != nil {
		log.Errorf("error creating db transaction to get latest block. Error: %v", err)
		return nil, nil, err
	}

	batch2, err := r.st.GetBatchByNumber(r.ctx, i, dbTx)
	if err != nil {
		log.Errorf("no batch %d. Error: %v", i, err)
		return batch2, nil, err
	}

	request := state.ProcessRequest{
		BatchNumber:     batch2.BatchNumber,
		OldStateRoot:    oldStateRoot,
		OldAccInputHash: oldAccInputHash,
		Coinbase:        batch2.Coinbase,
		Timestamp_V1:    batch2.Timestamp,

		GlobalExitRoot_V1: batch2.GlobalExitRoot,
		Transactions:      batch2.BatchL2Data,
	}
	log.Debugf("Processing batch %d: ntx:%d StateRoot:%s", batch2.BatchNumber, len(batch2.BatchL2Data), batch2.StateRoot)
	forkID := r.st.GetForkIDByBatchNumber(batch2.BatchNumber)
	fmt.Printf("")
	batchV2, err := state.DecodeBatchV2(batch2.BatchL2Data)
	var syncedTxs []types.Transaction
	for _, block := range batchV2.Blocks {
		for _, transaction := range block.Transactions {
			syncedTxs = append(syncedTxs, transaction.Tx)
		}
	}
	//syncedTxs, _, _, err := state.DecodeTxs(v2.Blocks, forkID)
	//if err != nil {
	//	log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, batch2.BatchL2Data)
	//	return batch2, nil, err
	//} else {
	r.output.numOfTransactionsInBatch(len(syncedTxs))
	//}
	var response *state.ProcessBatchResponse

	log.Infof("id:%d len_trs:%d oldStateRoot:%s", batch2.BatchNumber, len(syncedTxs), request.OldStateRoot)
	response, err = r.st.ProcessBatchV2(r.ctx, request, r.updateHasbDB)
	for _, blockResponse := range response.BlockResponses {
		for tx_i, txresponse := range blockResponse.TransactionResponses {
			if txresponse.RomError != nil {
				r.output.addTransactionError(tx_i, txresponse.RomError)
				log.Errorf("error processing batch %d. tx:%d Error: %v stateroot:%s", i, tx_i, txresponse.RomError, response.NewStateRoot)
				//return txresponse.RomError
			}
		}
	}

	if err != nil {
		r.output.isWrittenOnHashDB(false, response.FlushID)
		if rollbackErr := dbTx.Rollback(r.ctx); rollbackErr != nil {
			return batch2, response, fmt.Errorf(
				"failed to rollback dbTx: %s. Rollback err: %w",
				rollbackErr.Error(), err,
			)
		}
		log.Errorf("error processing batch %d. Error: %v", i, err)
		return batch2, response, err
	} else {
		r.output.isWrittenOnHashDB(r.updateHasbDB, response.FlushID)
	}
	if response.NewStateRoot != batch2.StateRoot {
		if rollbackErr := dbTx.Rollback(r.ctx); rollbackErr != nil {
			return batch2, response, fmt.Errorf(
				"failed to rollback dbTx: %s. Rollback err: %w",
				rollbackErr.Error(), err,
			)
		}
		log.Errorf("error processing batch %d. Error: state root differs: calculated: %s  != expected: %s", i, response.NewStateRoot, batch2.StateRoot)
		return batch2, response, fmt.Errorf("missmatch state root calculated: %s  != expected: %s", response.NewStateRoot, batch2.StateRoot)
	}

	if commitErr := dbTx.Commit(r.ctx); commitErr != nil {
		return batch2, response, fmt.Errorf(
			"failed to commit dbTx: %s. Commit err: %w",
			commitErr.Error(), err,
		)
	}

	log.Infof("Verified batch %d: ntx:%d StateRoot:%s", i, len(syncedTxs), batch2.StateRoot)

	return batch2, response, nil
}
