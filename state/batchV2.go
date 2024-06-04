package state

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrExecutingBatchOOC process batch fails because OOC (Out of counters)
	ErrExecutingBatchOOC = errors.New("Batch execution fails because: out of counters")
)

// ProcessingContextV2 is the necessary data that a batch needs to provide to the runtime,
// without the historical state data (processing receipt from previous batch)
type ProcessingContextV2 struct {
	BatchNumber          uint64
	Coinbase             common.Address
	Timestamp            *time.Time // Batch timeStamp and also TimestampLimit
	L1InfoRoot           common.Hash
	L1InfoTreeData       map[uint32]L1DataV2
	ForcedBatchNum       *uint64
	BatchL2Data          *[]byte
	ForcedBlockHashL1    *common.Hash
	SkipVerifyL1InfoRoot uint32
	GlobalExitRoot       common.Hash // GlobalExitRoot is not use for execute but use to OpenBatch (data on  DB)
	ExecutionMode        uint64
	ClosingReason        ClosingReason
}

// ProcessBatchV2 processes a batch for forkID >= ETROG
func (s *State) ProcessBatchV2(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponse, string, error) {
	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	l1InfoTreeData := make(map[uint32]*executor.L1DataV2)

	for k, v := range request.L1InfoTreeData_V2 {
		l1InfoTreeData[k] = &executor.L1DataV2{
			GlobalExitRoot: v.GlobalExitRoot.Bytes(),
			BlockHashL1:    v.BlockHashL1.Bytes(),
			MinTimestamp:   v.MinTimestamp,
		}
	}

	// Create Batch
	var processBatchRequest = &executor.ProcessBatchRequestV2{
		OldBatchNum:       request.BatchNumber - 1,
		Coinbase:          request.Coinbase.String(),
		ForcedBlockhashL1: request.ForcedBlockHashL1.Bytes(),
		BatchL2Data:       request.Transactions,
		OldStateRoot:      request.OldStateRoot.Bytes(),
		L1InfoRoot:        request.L1InfoRoot_V2.Bytes(),
		L1InfoTreeData:    l1InfoTreeData,
		OldAccInputHash:   request.OldAccInputHash.Bytes(),
		TimestampLimit:    request.TimestampLimit_V2,
		UpdateMerkleTree:  updateMT,
		ChainId:           s.cfg.ChainID,
		ForkId:            request.ForkID,
		ContextId:         uuid.NewString(),
	}

	if request.SkipFirstChangeL2Block_V2 {
		processBatchRequest.SkipFirstChangeL2Block = cTrue
	}

	if request.SkipWriteBlockInfoRoot_V2 {
		processBatchRequest.SkipWriteBlockInfoRoot = cTrue
	}

	if request.SkipVerifyL1InfoRoot_V2 {
		processBatchRequest.SkipVerifyL1InfoRoot = cTrue
	}

	res, err := s.sendBatchRequestToExecutorV2(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, "", err
	}

	var result *ProcessBatchResponse
	result, err = s.convertToProcessBatchResponseV2(res)
	if err != nil {
		return nil, "", err
	}

	return result, processBatchRequest.ContextId, nil
}

// ExecuteBatchV2 is used by the synchronizer to reprocess batches to compare generated state root vs stored one
func (s *State) ExecuteBatchV2(ctx context.Context, batch Batch, L1InfoTreeRoot common.Hash, l1InfoTreeData map[uint32]L1DataV2, timestampLimit time.Time, updateMerkleTree bool, skipVerifyL1InfoRoot uint32, forcedBlockHashL1 *common.Hash, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}

	// Get previous batch to get state root and local exit root
	previousBatch, err := s.GetBatchByNumber(ctx, batch.BatchNumber-1, dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:     batch.BatchNumber - 1,
		Coinbase:        batch.Coinbase.String(),
		BatchL2Data:     batch.BatchL2Data,
		OldStateRoot:    previousBatch.StateRoot.Bytes(),
		L1InfoRoot:      L1InfoTreeRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),
		TimestampLimit:  uint64(timestampLimit.Unix()),
		// Changed for new sequencer strategy
		UpdateMerkleTree:     updateMT,
		ChainId:              s.cfg.ChainID,
		ForkId:               forkId,
		ContextId:            uuid.NewString(),
		SkipVerifyL1InfoRoot: skipVerifyL1InfoRoot,
	}

	if forcedBlockHashL1 != nil {
		processBatchRequest.ForcedBlockhashL1 = forcedBlockHashL1.Bytes()
	} else {
		l1InfoTree := make(map[uint32]*executor.L1DataV2)
		for i, v := range l1InfoTreeData {
			l1InfoTree[i] = &executor.L1DataV2{
				GlobalExitRoot: v.GlobalExitRoot.Bytes(),
				BlockHashL1:    v.BlockHashL1.Bytes(),
				MinTimestamp:   v.MinTimestamp,
			}
		}
		processBatchRequest.L1InfoTreeData = l1InfoTree
	}

	// Send Batch to the Executor
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("ExecuteBatchV2[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
	log.Debugf("ExecuteBatchV2[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("ExecuteBatchV2[processBatchRequest.L1InfoRoot]: %v", hex.EncodeToHex(processBatchRequest.L1InfoRoot))
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("ExecuteBatchV2[processBatchRequest.TimestampLimit]: %v", processBatchRequest.TimestampLimit)
	log.Debugf("ExecuteBatchV2[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("ExecuteBatchV2[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.SkipVerifyL1InfoRoot]: %v", processBatchRequest.SkipVerifyL1InfoRoot)
	log.Debugf("ExecuteBatchV2[processBatchRequest.L1InfoTreeData]: %+v", l1InfoTreeData)

	processBatchResponse, err := s.executorClient.ProcessBatchV2(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorErrorV2(ctx, processBatchResponse.Error, processBatchRequest)
	}

	return processBatchResponse, err
}

func (s *State) processBatchV2(ctx context.Context, processingCtx *ProcessingContextV2, caller metrics.CallerLabel, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}

	lastBatches, err := s.GetLastNBatches(ctx, 2, dbTx) // nolint:gomnd
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	isBatchClosed, err := s.IsBatchClosed(ctx, processingCtx.BatchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if isBatchClosed {
		return nil, ErrBatchAlreadyClosed
	}

	// Check provided batch number is the latest in db
	if lastBatch.BatchNumber != processingCtx.BatchNumber {
		return nil, ErrInvalidBatchNumber
	}
	forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)

	var timestampLimitUnix uint64
	if processingCtx.Timestamp != nil {
		timestampLimitUnix = uint64(processingCtx.Timestamp.Unix())
	} else {
		timestampLimitUnix = uint64(time.Now().Unix())
	}
	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:          lastBatch.BatchNumber - 1,
		Coinbase:             lastBatch.Coinbase.String(),
		BatchL2Data:          *processingCtx.BatchL2Data,
		OldStateRoot:         previousBatch.StateRoot.Bytes(),
		OldAccInputHash:      previousBatch.AccInputHash.Bytes(),
		TimestampLimit:       timestampLimitUnix,
		UpdateMerkleTree:     cTrue,
		ChainId:              s.cfg.ChainID,
		ForkId:               forkID,
		ContextId:            uuid.NewString(),
		SkipVerifyL1InfoRoot: processingCtx.SkipVerifyL1InfoRoot,
		L1InfoRoot:           processingCtx.L1InfoRoot.Bytes(),
	}

	if processingCtx.ForcedBlockHashL1 != nil {
		log.Debug("Setting ForcedBlockhashL1: ", processingCtx.ForcedBlockHashL1)
		processBatchRequest.ForcedBlockhashL1 = processingCtx.ForcedBlockHashL1.Bytes()
	} else {
		l1InfoTreeData := make(map[uint32]*executor.L1DataV2)

		for k, v := range processingCtx.L1InfoTreeData {
			l1InfoTreeData[k] = &executor.L1DataV2{
				GlobalExitRoot: v.GlobalExitRoot.Bytes(),
				BlockHashL1:    v.BlockHashL1.Bytes(),
				MinTimestamp:   v.MinTimestamp,
			}
		}
		processBatchRequest.L1InfoTreeData = l1InfoTreeData
	}

	if processingCtx.L1InfoRoot != (common.Hash{}) {
		processBatchRequest.L1InfoRoot = processingCtx.L1InfoRoot.Bytes()
	} else {
		currentl1InfoRoot, err := s.GetCurrentL1InfoRoot(ctx, dbTx)
		if err != nil {
			log.Errorf("error getting current L1InfoRoot: %v", err)
			return nil, err
		}
		processBatchRequest.L1InfoRoot = currentl1InfoRoot.Bytes()
	}

	return s.sendBatchRequestToExecutorV2(ctx, processBatchRequest, caller)
}

func (s *State) sendBatchRequestToExecutorV2(ctx context.Context, batchRequest *executor.ProcessBatchRequestV2, caller metrics.CallerLabel) (*executor.ProcessBatchResponseV2, error) {
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}

	batchRequestLog := "OldBatchNum: %v, From: %v, OldStateRoot: %v, L1InfoRoot: %v, OldAccInputHash: %v, TimestampLimit: %v, Coinbase: %v, UpdateMerkleTree: %v, SkipFirstChangeL2Block: %v, SkipWriteBlockInfoRoot: %v, ChainId: %v, ForkId: %v, ContextId: %v, SkipVerifyL1InfoRoot: %v, ForcedBlockhashL1: %v, L1InfoTreeData: %+v, BatchL2Data: %v"

	l1DataStr := ""
	for i, l1Data := range batchRequest.L1InfoTreeData {
		l1DataStr += fmt.Sprintf("[%d]{GlobalExitRoot: %v, BlockHashL1: %v, MinTimestamp: %v},", i, hex.EncodeToHex(l1Data.GlobalExitRoot), hex.EncodeToHex(l1Data.BlockHashL1), l1Data.MinTimestamp)
	}
	if l1DataStr != "" {
		l1DataStr = l1DataStr[:len(l1DataStr)-1]
	}

	batchRequestLog = fmt.Sprintf(batchRequestLog, batchRequest.OldBatchNum, batchRequest.From, hex.EncodeToHex(batchRequest.OldStateRoot), hex.EncodeToHex(batchRequest.L1InfoRoot),
		hex.EncodeToHex(batchRequest.OldAccInputHash), batchRequest.TimestampLimit, batchRequest.Coinbase, batchRequest.UpdateMerkleTree, batchRequest.SkipFirstChangeL2Block,
		batchRequest.SkipWriteBlockInfoRoot, batchRequest.ChainId, batchRequest.ForkId, batchRequest.ContextId, batchRequest.SkipVerifyL1InfoRoot, hex.EncodeToHex(batchRequest.ForcedBlockhashL1),
		l1DataStr, hex.EncodeToHex(batchRequest.BatchL2Data))

	newBatchNum := batchRequest.OldBatchNum + 1
	log.Debugf("executor batch %d request, %s", newBatchNum, batchRequestLog)

	now := time.Now()
	batchResponse, err := s.executorClient.ProcessBatchV2(ctx, batchRequest)
	elapsed := time.Since(now)

	//workarroundDuplicatedBlock(res)
	if caller != metrics.DiscardCallerLabel {
		metrics.ExecutorProcessingTime(string(caller), elapsed)
	}

	if err != nil {
		log.Errorf("error executor ProcessBatchV2: %v", err)
		log.Errorf("error executor ProcessBatchV2: %s", err.Error())
		log.Errorf("error executor ProcessBatchV2 response: %v", batchResponse)
	} else {
		batchResponseToString := processBatchResponseToString(newBatchNum, batchResponse, elapsed)
		if batchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(batchResponse.Error)
			log.Warnf("executor batch %d response, executor error: %v", newBatchNum, err)
			log.Warn(batchResponseToString)
			s.eventLog.LogExecutorErrorV2(ctx, batchResponse.Error, batchRequest)
		} else if batchResponse.ErrorRom != executor.RomError_ROM_ERROR_NO_ERROR && executor.IsROMOutOfCountersError(batchResponse.ErrorRom) {
			err = executor.RomErr(batchResponse.ErrorRom)
			log.Warnf("executor batch %d response, ROM OOC, error: %v", newBatchNum, err)
			log.Warn(batchResponseToString)
		} else if batchResponse.ErrorRom != executor.RomError_ROM_ERROR_NO_ERROR {
			err = executor.RomErr(batchResponse.ErrorRom)
			log.Warnf("executor batch %d response, ROM error: %v", newBatchNum, err)
			log.Warn(batchResponseToString)
		} else {
			log.Debug(batchResponseToString)
		}
	}

	return batchResponse, err
}

func processBatchResponseToString(batchNum uint64, batchResponse *executor.ProcessBatchResponseV2, executionTime time.Duration) string {
	batchResponseLog := "executor batch %d response, Time: %v, NewStateRoot: %v, NewAccInputHash: %v, NewLocalExitRoot: %v, NewBatchNumber: %v, GasUsed: %v, FlushId: %v, StoredFlushId: %v, ProverId:%v, ForkId:%v, Error: %v\n"
	batchResponseLog = fmt.Sprintf(batchResponseLog, batchNum, executionTime, hex.EncodeToHex(batchResponse.NewStateRoot), hex.EncodeToHex(batchResponse.NewAccInputHash), hex.EncodeToHex(batchResponse.NewLocalExitRoot),
		batchResponse.NewBatchNum, batchResponse.GasUsed, batchResponse.FlushId, batchResponse.StoredFlushId, batchResponse.ProverId, batchResponse.ForkId, batchResponse.Error)

	for blockIndex, block := range batchResponse.BlockResponses {
		prefix := "  " + fmt.Sprintf("block[%v]: ", blockIndex)
		batchResponseLog += blockResponseToString(block, prefix)
	}

	return batchResponseLog
}
func blockResponseToString(blockResponse *executor.ProcessBlockResponseV2, prefix string) string {
	blockResponseLog := prefix + "ParentHash: %v, Coinbase: %v, GasLimit: %v, BlockNumber: %v, Timestamp: %v, GlobalExitRoot: %v, BlockHashL1: %v, GasUsed: %v, BlockInfoRoot: %v, BlockHash: %v\n"
	blockResponseLog = fmt.Sprintf(blockResponseLog, common.BytesToHash(blockResponse.ParentHash), blockResponse.Coinbase, blockResponse.GasLimit, blockResponse.BlockNumber, blockResponse.Timestamp,
		common.BytesToHash(blockResponse.Ger), common.BytesToHash(blockResponse.BlockHashL1), blockResponse.GasUsed, common.BytesToHash(blockResponse.BlockInfoRoot), common.BytesToHash(blockResponse.BlockHash))

	for txIndex, tx := range blockResponse.Responses {
		prefix := "    " + fmt.Sprintf("tx[%v]: ", txIndex)
		blockResponseLog += transactionResponseToString(tx, prefix)
	}

	return blockResponseLog
}

func transactionResponseToString(txResponse *executor.ProcessTransactionResponseV2, prefix string) string {
	txResponseLog := prefix + "TxHash: %v, TxHashL2: %v, Type: %v, StateRoot:%v, GasUsed: %v, GasLeft: %v, GasRefund: %v, Error: %v\n"
	txResponseLog = fmt.Sprintf(txResponseLog, common.BytesToHash(txResponse.TxHash), common.BytesToHash(txResponse.TxHashL2), txResponse.Type,
		common.BytesToHash(txResponse.StateRoot), txResponse.GasUsed, txResponse.GasLeft, txResponse.GasRefunded, txResponse.Error)

	return txResponseLog
}

// ProcessAndStoreClosedBatchV2 is used by the Synchronizer to add a closed batch into the data base. Values returned are the new stateRoot,
// the flushID (incremental value returned by executor),
// the ProverID (executor running ID) the result of closing the batch.
func (s *State) ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx ProcessingContextV2, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error) {
	debugPrefix := fmt.Sprint("Batch ", processingCtx.BatchNumber, ": ProcessAndStoreClosedBatchV2: ")

	BatchL2Data := processingCtx.BatchL2Data
	if BatchL2Data == nil {
		log.Warnf("%s processingCtx.BatchL2Data is nil, assuming is empty", debugPrefix, processingCtx.BatchNumber)
		var BatchL2DataEmpty []byte
		BatchL2Data = &BatchL2DataEmpty
	}

	if dbTx == nil {
		return common.Hash{}, noFlushID, noProverID, ErrDBTxNil
	}
	// Avoid writing twice to the DB the BatchL2Data that is going to be written also in the call closeBatch
	// TODO: check if is need this

	convertedProcessingContextV1, err := convertProcessingContext(&processingCtx)
	if err != nil {
		log.Errorf("%s error convertProcessingContext: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}
	convertedProcessingContextV1.BatchL2Data = nil
	if err := s.OpenBatch(ctx, *convertedProcessingContextV1, dbTx); err != nil {
		log.Errorf("%s error OpenBatch: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}
	processed, err := s.processBatchV2(ctx, &processingCtx, caller, dbTx)
	if err != nil && processed.ErrorRom == executor.RomError_ROM_ERROR_NO_ERROR {
		log.Errorf("%s error processBatchV2: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}

	processedBatch, err := s.convertToProcessBatchResponseV2(processed)
	if err != nil {
		log.Errorf("%s error convertToProcessBatchResponseV2: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}
	if processedBatch.IsRomOOCError {
		log.Errorf("%s error isRomOOCError: %v", debugPrefix, err)
	}

	if len(processedBatch.BlockResponses) > 0 && !processedBatch.IsRomOOCError && processedBatch.RomError_V2 == nil {
		for _, blockResponse := range processedBatch.BlockResponses {
			_, err = s.StoreL2Block(ctx, processingCtx.BatchNumber, blockResponse, nil, dbTx)
			if err != nil {
				log.Errorf("%s error StoreL2Block: %v", debugPrefix, err)
				return common.Hash{}, noFlushID, noProverID, err
			}
		}
	}
	return common.BytesToHash(processed.NewStateRoot), processed.FlushId, processed.ProverId, s.CloseBatchInStorage(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
		AccInputHash:  processedBatch.NewAccInputHash,
		BatchL2Data:   *BatchL2Data,
		ClosingReason: processingCtx.ClosingReason,
	}, dbTx)
}
