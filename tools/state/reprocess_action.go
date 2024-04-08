package main

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
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
	forkID := r.st.GetForkIDByBatchNumber(batch2.BatchNumber)
	//l1data, l1hash, _, err := r.st.GetL1InfoTreeDataFromBatchL2Data(context.Background(), batch2.BatchL2Data, dbTx)
	//if err != nil {
	//	log.Errorf("error getting L1InfoTreeData from batch. Error: %v", err)
	//	return batch2, nil, err
	//}
	l2blocks, err := r.st.GetL2BlocksByBatchNumber(context.Background(), batch2.BatchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting L2BlocksByBatchNumber. Error: %v", err)
		return batch2, nil, err
	}
	//for _, block := range l2blocks {
	//	log.Infof("L2Block %d: %s", block.BlockInfoRoot(), block.DeltaTimestamp)
	//}

	l2data, err := state.DecodeBatchV2(batch2.BatchL2Data)
	for ii, block := range l2data.Blocks {

		l2block := l2blocks[ii]
		batchL2Data := []byte{}

		// Add changeL2Block to batchL2Data
		changeL2BlockBytes := r.st.BuildChangeL2Block(block.DeltaTimestamp, block.IndexL1InfoTree)
		batchL2Data = append(batchL2Data, changeL2BlockBytes...)

		// Add transactions data to batchL2Data
		var txs []types.Transaction
		var epg []uint8
		for _, tx := range block.Transactions {
			txs = append(txs, tx.Tx)
			epg = append(epg, tx.EfficiencyPercentage)
		}
		transactions, err := state.EncodeTransactions(txs, epg, forkID)
		if err != nil {
			log.Errorf("error encoding transactions. Error: %v", err)
			return batch2, nil, err
		}
		batchL2Data = append(batchL2Data, transactions...)

		index, err := r.st.GetL1InfoRootLeafByIndex(context.Background(), block.IndexL1InfoTree, dbTx)
		if err != nil {
			log.Errorf("error getting L1InfoRootLeafByIndex. Error: %v", err)
			return batch2, nil, err
		}
		//r.output.numOfTransactionsInBatch(len(transactions))

		//request := state.ProcessRequest{
		//	BatchNumber:       batch2.BatchNumber,
		//	OldStateRoot:      oldStateRoot,
		//	OldAccInputHash:   oldAccInputHash,
		//	Coinbase:          batch2.Coinbase,
		//	TimestampLimit_V2: l2block.Time(),
		//	L1InfoRoot_V2:     state.GetMockL1InfoRoot(),
		//	L1InfoTreeData_V2: map[uint32]state.L1DataV2{},
		//	ForkID:            forkID,
		//
		//	GlobalExitRoot_V1:       batch2.GlobalExitRoot,
		//	Transactions:            batchL2Data,
		//	SkipVerifyL1InfoRoot_V2: true,
		//}
		request := state.ProcessRequest{
			BatchNumber:               batch2.BatchNumber,
			OldStateRoot:              oldStateRoot,
			Coinbase:                  batch2.Coinbase,
			L1InfoRoot_V2:             state.GetMockL1InfoRoot(),
			TimestampLimit_V2:         uint64(time.Now().Unix()),
			Transactions:              batchL2Data,
			SkipFirstChangeL2Block_V2: false,
			SkipWriteBlockInfoRoot_V2: false,
			Caller:                    stateMetrics.DiscardCallerLabel,
			ForkID:                    forkID,
			SkipVerifyL1InfoRoot_V2:   true,
			L1InfoTreeData_V2:         map[uint32]state.L1DataV2{},
		}
		//request.L1InfoTreeData_V2[block.IndexL1InfoTree] = state.L1DataV2{
		//	GlobalExitRoot: batch2.GlobalExitRoot,
		//	BlockHashL1:    common.Hash{},
		//	MinTimestamp:   0,
		//}

		request.L1InfoTreeData_V2 = map[uint32]state.L1DataV2{
			block.IndexL1InfoTree: {
				GlobalExitRoot: index.GlobalExitRoot.GlobalExitRoot,
				BlockHashL1:    index.PreviousBlockHash,
				MinTimestamp:   uint64(index.GlobalExitRoot.Timestamp.Unix()),
			},
		}
		log.Infof("deltatimestamp: %d, minTimestamp:%d, current: %d", block.DeltaTimestamp, index.GlobalExitRoot.Timestamp.Unix(), time.Now().Unix())
		//request.L1InfoRoot_V2 = l1hash

		response, err := r.st.ProcessBatchV2(context.Background(), request, false)
		if err != nil {
			log.Errorf("error processing batch %d. Error: %v", i, err)
			return batch2, nil, err
		}
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
		log.Infof("Verified batch %d: ntx:%d StateRoot:%s", i, len(transactions), response.NewStateRoot)
	}

	//
	//request := state.ProcessRequest{
	//	BatchNumber:       batch2.BatchNumber,
	//	OldStateRoot:      oldStateRoot,
	//	OldAccInputHash:   oldAccInputHash,
	//	Coinbase:          batch2.Coinbase,
	//	Timestamp_V1:      batch2.Timestamp,
	//	TimestampLimit_V2: uint64(batch2.Timestamp.Unix()),
	//	L1InfoRoot_V2:     state.GetMockL1InfoRoot(),
	//	L1InfoTreeData_V2: map[uint32]state.L1DataV2{},
	//
	//	GlobalExitRoot_V1: batch2.GlobalExitRoot,
	//	Transactions:      batch2.BatchL2Data,
	//}
	//
	//
	//log.Debugf("Processing batch %d: ntx:%d StateRoot:%s", batch2.BatchNumber, len(batch2.BatchL2Data), batch2.StateRoot)
	//
	//fmt.Printf("fork id: %d\n", forkID)
	//request.ForkID = forkID
	//
	////syncedTxs, _, _, err := state.DecodeTxs(v2.Blocks, forkID)
	////if err != nil {
	////	log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, batch2.BatchL2Data)
	////	return batch2, nil, err
	////} else {
	//
	////}
	//var response *state.ProcessBatchResponse
	//
	//log.Infof("id:%d len_trs:%d oldStateRoot:%s", batch2.BatchNumber, len(syncedTxs), request.OldStateRoot)
	//response, err = r.st.ProcessBatchV2(r.ctx, request, r.updateHasbDB)

	//if response.NewStateRoot != batch2.StateRoot {
	//	if rollbackErr := dbTx.Rollback(r.ctx); rollbackErr != nil {
	//		return batch2, response, fmt.Errorf(
	//			"failed to rollback dbTx: %s. Rollback err: %w",
	//			rollbackErr.Error(), err,
	//		)
	//	}
	//	log.Errorf("error processing batch %d. Error: state root differs: calculated: %s  != expected: %s", i, response.NewStateRoot, batch2.StateRoot)
	//	return batch2, response, fmt.Errorf("missmatch state root calculated: %s  != expected: %s", response.NewStateRoot, batch2.StateRoot)
	//}

	if commitErr := dbTx.Commit(r.ctx); commitErr != nil {
		return batch2, nil, fmt.Errorf(
			"failed to commit dbTx: %s. Commit err: %w",
			commitErr.Error(), err,
		)
	}

	//log.Infof("Verified batch %d: ntx:%d StateRoot:%s", i, len(syncedTxs), batch2.StateRoot)

	return batch2, nil, nil
}
