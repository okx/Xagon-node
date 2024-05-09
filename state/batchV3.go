package state

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/google/uuid"
)

// ProcessBatchV3 processes a batch for forkID >= FEIJOA
func (s *State) ProcessBatchV3(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponse, error) {
	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	l1InfoTreeData := make(map[uint32]*executor.L1DataV3)

	for k, v := range request.L1InfoTreeData_V3 {
		l1InfoTreeData[k] = &executor.L1DataV3{
			GlobalExitRoot:        v.GlobalExitRoot.Bytes(),
			BlockHashL1:           v.BlockHashL1.Bytes(),
			MinTimestamp:          v.MinTimestamp,
			SmtProofPreviousIndex: v.SmtProofPreviousIndex,
			InitialHistoricRoot:   v.InitialHistoricRoot.Bytes(),
		}
	}

	// Create Batch
	var processBatchRequest = &executor.ProcessBatchRequestV3{
		OldStateRoot:            request.OldStateRoot.Bytes(),
		OldAccInputHash:         request.OldAccInputHash.Bytes(),
		PreviousL1InfoTreeRoot:  request.PreviousL1InfoTreeRoot_V3.Bytes(),
		PreviousL1InfoTreeIndex: request.PreviousL1InfoTreeIndex_V3,
		ChainId:                 s.cfg.ChainID,
		ForkId:                  request.ForkID,
		BatchL2Data:             request.Transactions,
		Coinbase:                request.Coinbase.String(),
		UpdateMerkleTree:        updateMT,
		L1InfoTreeData:          l1InfoTreeData,
		ContextId:               uuid.NewString(),
	}

	if request.SkipFirstChangeL2Block_V2 {
		processBatchRequest.SkipFirstChangeL2Block = cTrue
	}

	if request.SkipWriteBlockInfoRoot_V2 {
		processBatchRequest.SkipWriteBlockInfoRoot = cTrue
	}

	res, err := s.sendBatchRequestToExecutorV3(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, err
	}

	var result *ProcessBatchResponse
	result, err = s.convertToProcessBatchResponseV3(res)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *State) sendBatchRequestToExecutorV3(ctx context.Context, batchRequest *executor.ProcessBatchRequestV3, caller metrics.CallerLabel) (*executor.ProcessBatchResponseV3, error) {
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}

	l1DataStr := ""
	for i, l1Data := range batchRequest.L1InfoTreeData {
		l1DataStr += fmt.Sprintf("[%d]{GlobalExitRoot: %v, BlockHashL1: %v, MinTimestamp: %v},", i, hex.EncodeToHex(l1Data.GlobalExitRoot), hex.EncodeToHex(l1Data.BlockHashL1), l1Data.MinTimestamp)
	}
	if l1DataStr != "" {
		l1DataStr = l1DataStr[:len(l1DataStr)-1]
	}

	// Log the batch request
	batchRequestLog := "OldStateRoot: %v, OldAccInputHash: %v, PreviousL1InfoTreeRoot: %v, PreviousL1InfoTreeIndex: %v, ChainId: %v, ForkId: %v, BatchL2Data: %v, Coinbase: %v, UpdateMerkleTree: %v, L1InfoTreeData: %+v, ContextId: %v, SkipFirstChangeL2Block: %v, SkipWriteBlockInfoRoot: %v"
	batchRequestLog = fmt.Sprintf(batchRequestLog, hex.EncodeToHex(batchRequest.OldStateRoot), hex.EncodeToHex(batchRequest.OldAccInputHash), hex.EncodeToHex(batchRequest.PreviousL1InfoTreeRoot), batchRequest.PreviousL1InfoTreeIndex, batchRequest.ChainId, batchRequest.ForkId, len(batchRequest.BatchL2Data), batchRequest.Coinbase, batchRequest.UpdateMerkleTree, l1DataStr, batchRequest.ContextId, batchRequest.SkipFirstChangeL2Block, batchRequest.SkipWriteBlockInfoRoot)

	log.Debugf("executor batch request, %s", batchRequestLog)

	now := time.Now()
	batchResponse, err := s.executorClient.ProcessBatchV3(ctx, batchRequest)
	elapsed := time.Since(now)

	// workarroundDuplicatedBlock(res)
	if caller != metrics.DiscardCallerLabel {
		metrics.ExecutorProcessingTime(string(caller), elapsed)
	}

	if err != nil {
		log.Errorf("error executor ProcessBatchV3: %v", err)
		log.Errorf("error executor ProcessBatchV3: %s", err.Error())
		log.Errorf("error executor ProcessBatchV3 response: %v", batchResponse)
	} else {
		batchResponseToString := processBatchResponseV3ToString(batchResponse, elapsed)
		if batchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(batchResponse.Error)
			log.Warnf("executor batch response, executor error: %v", err)
			log.Warn(batchResponseToString)
			s.eventLog.LogExecutorError(ctx, batchResponse.Error, batchRequest)
		} else if batchResponse.ErrorRom != executor.RomError_ROM_ERROR_NO_ERROR && executor.IsROMOutOfCountersError(batchResponse.ErrorRom) {
			err = executor.RomErr(batchResponse.ErrorRom)
			log.Warnf("executor batch response, ROM OOC, error: %v", err)
			log.Warn(batchResponseToString)
		} else if batchResponse.ErrorRom != executor.RomError_ROM_ERROR_NO_ERROR {
			err = executor.RomErr(batchResponse.ErrorRom)
			log.Warnf("executor batch response, ROM error: %v", err)
			log.Warn(batchResponseToString)
		} else {
			log.Debug(batchResponseToString)
		}
	}

	return batchResponse, err
}

func processBatchResponseV3ToString(batchResponse *executor.ProcessBatchResponseV3, executionTime time.Duration) string {
	batchResponseLog := "executor batch response, Time: %v, NewStateRoot: %v, NewAccInputHash: %v, NewLocalExitRoot: %v, GasUsed: %v, FlushId: %v, StoredFlushId: %v, ProverId:%v, ForkId:%v, Error: %v\n"
	batchResponseLog = fmt.Sprintf(batchResponseLog, executionTime, hex.EncodeToHex(batchResponse.NewStateRoot), hex.EncodeToHex(batchResponse.NewAccInputHash), hex.EncodeToHex(batchResponse.NewLocalExitRoot),
		batchResponse.GasUsed, batchResponse.FlushId, batchResponse.StoredFlushId, batchResponse.ProverId, batchResponse.ForkId, batchResponse.Error)

	for blockIndex, block := range batchResponse.BlockResponses {
		prefix := "  " + fmt.Sprintf("block[%v]: ", blockIndex)
		batchResponseLog += blockResponseV2ToString(block, prefix)
	}

	return batchResponseLog
}
