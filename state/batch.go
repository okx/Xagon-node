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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	cTrue             = 1
	cFalse            = 0
	noFlushID  uint64 = 0
	noProverID string = ""

	// MockL1InfoRootHex is used to send batches to the Executor
	// the number below represents this formula:
	//
	// 	mockL1InfoRoot := common.Hash{}
	// for i := 0; i < len(mockL1InfoRoot); i++ {
	// 	  mockL1InfoRoot[i] = byte(i)
	// }
	MockL1InfoRootHex = "0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
)

// Batch struct
type Batch struct {
	BatchNumber   uint64
	Coinbase      common.Address
	BatchL2Data   []byte
	StateRoot     common.Hash
	LocalExitRoot common.Hash
	AccInputHash  common.Hash
	// Timestamp (<=incaberry) -> batch time
	// 			 (>incaberry) -> minTimestamp used in batch creation, real timestamp is in virtual_batch.batch_timestamp
	Timestamp              time.Time
	Transactions           []types.Transaction
	GlobalExitRoot         common.Hash
	ForcedBatchNum         *uint64
	Resources              BatchResources
	HighReservedZKCounters ZKCounters
	// WIP: if WIP == true is a openBatch
	WIP bool
}

// ProcessingContext is the necessary data that a batch needs to provide to the runtime,
// without the historical state data (processing receipt from previous batch)
type ProcessingContext struct {
	BatchNumber    uint64
	Coinbase       common.Address
	Timestamp      time.Time
	GlobalExitRoot common.Hash
	ForcedBatchNum *uint64
	BatchL2Data    *[]byte
	ClosingReason  ClosingReason
}

// ClosingReason represents the reason why a batch is closed.
type ClosingReason string

const (
	// EmptyClosingReason is the closing reason used when a batch is not closed
	EmptyClosingReason ClosingReason = ""
	// MaxTxsClosingReason is the closing reason used when a batch reachs the max transactions per batch
	MaxTxsClosingReason ClosingReason = "Max transactions"
	// ResourceExhaustedClosingReason is the closing reason used when a batch has a resource (zkCounter or Bytes) exhausted
	ResourceExhaustedClosingReason ClosingReason = "Resource exhausted"
	// ResourceMarginExhaustedClosingReason is the closing reason used when a batch has a resource (zkCounter or Bytes) margin exhausted
	ResourceMarginExhaustedClosingReason ClosingReason = "Resource margin exhausted"
	// ForcedBatchClosingReason is the closing reason used when a batch is a forced batch
	ForcedBatchClosingReason ClosingReason = "Forced batch"
	// ForcedBatchDeadlineClosingReason is the closing reason used when forced batch deadline is reached
	ForcedBatchDeadlineClosingReason ClosingReason = "Forced batch deadline"
	// MaxDeltaTimestampClosingReason is the closing reason used when max delta batch timestamp is reached
	MaxDeltaTimestampClosingReason ClosingReason = "Max delta timestamp"
	// NoTxFitsClosingReason is the closing reason used when any of the txs in the pool (worker) fits in the remaining resources of the batch
	NoTxFitsClosingReason ClosingReason = "No transaction fits"
	// L2BlockReorgClonsingReason is the closing reason used when we have a L2 block reorg (unexpected error, like OOC, when processing L2 block)
	L2BlockReorgClonsingReason ClosingReason = "L2 block reorg"

	// Reason due Synchronizer
	// ------------------------------------------------------------------------------------------

	// SyncL1EventInitialBatchClosingReason is the closing reason used when a batch is closed by the synchronizer due to an initial batch (first batch mode forced)
	SyncL1EventInitialBatchClosingReason ClosingReason = "Sync L1: initial"
	// SyncL1EventSequencedBatchClosingReason is the closing reason used when a batch is closed by the synchronizer due to a sequenced batch event from L1
	SyncL1EventSequencedBatchClosingReason ClosingReason = "Sync L1: sequenced"
	// SyncL1EventSequencedForcedBatchClosingReason is the closing reason used when a batch is closed by the synchronizer due to a sequenced forced batch event from L1
	SyncL1EventSequencedForcedBatchClosingReason ClosingReason = "Sync L1: forced"
	// SyncL1EventUpdateEtrogSequenceClosingReason is the closing reason used when a batch is closed by the synchronizer due to an UpdateEtrogSequence event from L1 that inject txs
	SyncL1EventUpdateEtrogSequenceClosingReason ClosingReason = "Sync L1: injected"
	// SyncL2TrustedBatchClosingReason is the closing reason used when a batch is closed by the synchronizer due to a trusted batch from L2
	SyncL2TrustedBatchClosingReason ClosingReason = "Sync L2: trusted"
	// SyncGenesisBatchClosingReason is the closing reason used when genesis batch is created by synchronizer
	SyncGenesisBatchClosingReason ClosingReason = "Sync: genesis"
)

// ProcessingReceipt indicates the outcome (StateRoot, AccInputHash) of processing a batch
type ProcessingReceipt struct {
	BatchNumber    uint64
	StateRoot      common.Hash
	LocalExitRoot  common.Hash
	GlobalExitRoot common.Hash
	AccInputHash   common.Hash
	// Txs           []types.Transaction
	BatchL2Data            []byte
	ClosingReason          ClosingReason
	BatchResources         BatchResources
	HighReservedZKCounters ZKCounters
}

// VerifiedBatch represents a VerifiedBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	TxHash      common.Hash
	StateRoot   common.Hash
	IsTrusted   bool
}

// VirtualBatch represents a VirtualBatch
type VirtualBatch struct {
	BatchNumber   uint64
	TxHash        common.Hash
	Coinbase      common.Address
	SequencerAddr common.Address
	BlockNumber   uint64
	L1InfoRoot    *common.Hash
	// TimestampBatchEtrog etrog: Batch timestamp comes from L1 block timestamp
	//  for previous batches is NULL because the batch timestamp is in batch table
	TimestampBatchEtrog *time.Time
}

// Sequence represents the sequence interval
type Sequence struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
}

// OpenBatch adds a new batch into the state, with the necessary data to start processing transactions within it.
// It's meant to be used by sequencers, since they don't necessarily know what transactions are going to be added
// in this batch yet. In other words it's the creation of a WIP batch.
// Note that this will add a batch with batch number N + 1, where N it's the greatest batch number on the state.
func (s *State) OpenBatch(ctx context.Context, processingContext ProcessingContext, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}
	// Check if the batch that is being opened has batch num + 1 compared to the latest batch
	lastBatchNum, err := s.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum+1 != processingContext.BatchNumber {
		return fmt.Errorf("%w number %d, should be %d", ErrUnexpectedBatch, processingContext.BatchNumber, lastBatchNum+1)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if !isLastBatchClosed {
		return ErrLastBatchShouldBeClosed
	}
	// Check that timestamp is equal or greater compared to previous batch
	prevTimestamp, err := s.GetLastBatchTime(ctx, dbTx)
	if err != nil {
		return err
	}
	if prevTimestamp.Unix() > processingContext.Timestamp.Unix() {
		return fmt.Errorf(" oldBatch(%d) tstamp=%d > openingBatch(%d)=%d err: %w", lastBatchNum, prevTimestamp.Unix(), processingContext.BatchNumber, processingContext.Timestamp.Unix(), ErrTimestampGE)
	}
	return s.OpenBatchInStorage(ctx, processingContext, dbTx)
}

// OpenWIPBatch adds a new WIP batch into the state
func (s *State) OpenWIPBatch(ctx context.Context, batch Batch, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	//TODO: Use s.GetLastBatch to retrieve number and time and avoid to do 2 queries
	// Check if the batch that is being opened has batch num + 1 compared to the latest batch
	lastBatchNum, err := s.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum+1 != batch.BatchNumber {
		return fmt.Errorf("%w number %d, should be %d", ErrUnexpectedBatch, batch.BatchNumber, lastBatchNum+1)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if !isLastBatchClosed {
		return ErrLastBatchShouldBeClosed
	}
	// Check that timestamp is equal or greater compared to previous batch
	prevTimestamp, err := s.GetLastBatchTime(ctx, dbTx)
	if err != nil {
		return err
	}
	if prevTimestamp.Unix() > batch.Timestamp.Unix() {
		return ErrTimestampGE
	}
	return s.OpenWIPBatchInStorage(ctx, batch, dbTx)
}

// GetWIPBatch returns the wip batch in the state
func (s *State) GetWIPBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	return s.GetWIPBatchInStorage(ctx, batchNumber, dbTx)
}

// ProcessSequencerBatch is used by the sequencers to process transactions into an open batch
func (s *State) ProcessSequencerBatch(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessSequencerBatch start")

	processBatchResponse, err := s.processBatch(ctx, batchNumber, batchL2Data, caller, dbTx)
	if err != nil {
		return nil, err
	}

	result, err := s.convertToProcessBatchResponse(processBatchResponse)
	if err != nil {
		return nil, err
	}
	log.Debugf("ProcessSequencerBatch end")
	log.Debugf("*******************************************")
	return result, nil
}

// ProcessBatch processes a batch
func (s *State) ProcessBatch(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessBatch start")

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	forkID := s.GetForkIDByBatchNumber(request.BatchNumber)

	// Create Batch
	var processBatchRequest = &executor.ProcessBatchRequest{
		OldBatchNum:      request.BatchNumber - 1,
		Coinbase:         request.Coinbase.String(),
		BatchL2Data:      request.Transactions,
		OldStateRoot:     request.OldStateRoot.Bytes(),
		GlobalExitRoot:   request.GlobalExitRoot_V1.Bytes(),
		OldAccInputHash:  request.OldAccInputHash.Bytes(),
		EthTimestamp:     uint64(request.Timestamp_V1.Unix()),
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}
	res, err := s.sendBatchRequestToExecutor(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, err
	}

	var result *ProcessBatchResponse
	result, err = s.convertToProcessBatchResponse(res)
	if err != nil {
		return nil, err
	}

	log.Debugf("ProcessBatch end")
	log.Debugf("*******************************************")

	return result, nil
}

// ExecuteBatch is used by the synchronizer to reprocess batches to compare generated state root vs stored one
// It is also used by the sequencer in order to calculate used zkCounter of a WIPBatch
func (s *State) ExecuteBatch(ctx context.Context, batch Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error) {
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
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:     batch.BatchNumber - 1,
		Coinbase:        batch.Coinbase.String(),
		BatchL2Data:     batch.BatchL2Data,
		OldStateRoot:    previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:  batch.GlobalExitRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),
		EthTimestamp:    uint64(batch.Timestamp.Unix()),
		// Changed for new sequencer strategy
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		ContextId:        uuid.NewString(),
	}

	// Send Batch to the Executor
	log.Debugf("ExecuteBatch[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("ExecuteBatch[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
	log.Debugf("ExecuteBatch[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("ExecuteBatch[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("ExecuteBatch[processBatchRequest.GlobalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("ExecuteBatch[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("ExecuteBatch[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("ExecuteBatch[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("ExecuteBatch[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("ExecuteBatch[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("ExecuteBatch[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("ExecuteBatch[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
	}

	return processBatchResponse, err
}

func (s *State) processBatch(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error) {
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

	isBatchClosed, err := s.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if isBatchClosed {
		return nil, ErrBatchAlreadyClosed
	}

	// Check provided batch number is the latest in db
	if lastBatch.BatchNumber != batchNumber {
		return nil, ErrInvalidBatchNumber
	}
	forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber - 1,
		Coinbase:         lastBatch.Coinbase.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
		UpdateMerkleTree: cTrue,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	return s.sendBatchRequestToExecutor(ctx, processBatchRequest, caller)
}

func (s *State) sendBatchRequestToExecutor(ctx context.Context, processBatchRequest *executor.ProcessBatchRequest, caller metrics.CallerLabel) (*executor.ProcessBatchResponse, error) {
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}
	// Send Batch to the Executor
	if caller != metrics.DiscardCallerLabel {
		log.Debugf("processBatch[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		log.Debugf("processBatch[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("processBatch[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("processBatch[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("processBatch[processBatchRequest.GlobalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("processBatch[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("processBatch[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("processBatch[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("processBatch[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("processBatch[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("processBatch[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
		log.Debugf("processBatch[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)
	}
	now := time.Now()
	res, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Errorf("Error s.executorClient.ProcessBatch: %v", err)
		log.Errorf("Error s.executorClient.ProcessBatch: %s", err.Error())
		log.Errorf("Error s.executorClient.ProcessBatch response: %v", res)
	} else if res.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(res.Error)
		s.eventLog.LogExecutorError(ctx, res.Error, processBatchRequest)
	}
	elapsed := time.Since(now)
	if caller != metrics.DiscardCallerLabel {
		metrics.ExecutorProcessingTime(string(caller), elapsed)
	}
	log.Infof("Batch: %d took %v to be processed by the executor ", processBatchRequest.OldBatchNum+1, elapsed)

	return res, err
}

func (s *State) isBatchClosable(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	// Check if the batch that is being closed is the last batch
	lastBatchNum, err := s.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum != receipt.BatchNumber {
		return fmt.Errorf("%w number %d, should be %d", ErrUnexpectedBatch, receipt.BatchNumber, lastBatchNum)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if isLastBatchClosed {
		return ErrBatchAlreadyClosed
	}

	return nil
}

// CloseBatch is used to close a batch
func (s *State) CloseBatch(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	err := s.isBatchClosable(ctx, receipt, dbTx)
	if err != nil {
		return err
	}

	return s.CloseBatchInStorage(ctx, receipt, dbTx)
}

// CloseWIPBatch is used by sequencer to close the wip batch
func (s *State) CloseWIPBatch(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	return s.CloseWIPBatchInStorage(ctx, receipt, dbTx)
}

// ProcessAndStoreClosedBatch is used by the Synchronizer to add a closed batch into the data base. Values returned are the new stateRoot,
// the flushID (incremental value returned by executor),
// the ProverID (executor running ID) the result of closing the batch.
func (s *State) ProcessAndStoreClosedBatch(ctx context.Context, processingCtx ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error) {
	BatchL2Data := processingCtx.BatchL2Data
	if BatchL2Data == nil {
		log.Warnf("Batch %v: ProcessAndStoreClosedBatch: processingCtx.BatchL2Data is nil, assuming is empty", processingCtx.BatchNumber)
		var BatchL2DataEmpty []byte
		BatchL2Data = &BatchL2DataEmpty
	}
	// Decode transactions
	forkID := s.GetForkIDByBatchNumber(processingCtx.BatchNumber)
	decodedTransactions, _, _, err := DecodeTxs(*BatchL2Data, forkID)
	if err != nil && !errors.Is(err, ErrInvalidData) {
		log.Debugf("error decoding transactions: %v", err)
		return common.Hash{}, noFlushID, noProverID, err
	}

	// Open the batch and process the txs
	if dbTx == nil {
		return common.Hash{}, noFlushID, noProverID, ErrDBTxNil
	}
	// Avoid writing twice to the DB the BatchL2Data that is going to be written also in the call closeBatch
	processingCtx.BatchL2Data = nil
	if err := s.OpenBatch(ctx, processingCtx, dbTx); err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}
	processed, err := s.processBatch(ctx, processingCtx.BatchNumber, *BatchL2Data, caller, dbTx)
	if err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}

	// Sanity check
	if len(decodedTransactions) != len(processed.Responses) {
		log.Errorf("number of decoded (%d) and processed (%d) transactions do not match", len(decodedTransactions), len(processed.Responses))
	}

	// Filter unprocessed txs and decode txs to store metadata
	// note that if the batch is not well encoded it will result in an empty batch (with no txs)
	for i := 0; i < len(processed.Responses); i++ {
		if !IsStateRootChanged(processed.Responses[i].Error) {
			if executor.IsROMOutOfCountersError(processed.Responses[i].Error) {
				processed.Responses = []*executor.ProcessTransactionResponse{}
				break
			}

			// Remove unprocessed tx
			if i == len(processed.Responses)-1 {
				processed.Responses = processed.Responses[:i]
			} else {
				processed.Responses = append(processed.Responses[:i], processed.Responses[i+1:]...)
			}
			i--
		}
	}

	processedBatch, err := s.convertToProcessBatchResponse(processed)
	if err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}

	if len(processedBatch.BlockResponses) > 0 {
		// Store processed txs into the batch
		err = s.StoreTransactions(ctx, processingCtx.BatchNumber, processedBatch.BlockResponses, nil, dbTx)
		if err != nil {
			return common.Hash{}, noFlushID, noProverID, err
		}
	}

	// Close batch
	return common.BytesToHash(processed.NewStateRoot), processed.FlushId, processed.ProverId, s.CloseBatchInStorage(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
		AccInputHash:  processedBatch.NewAccInputHash,
		BatchL2Data:   *BatchL2Data,
	}, dbTx)
}

// GetLastBatch gets latest batch (closed or not) on the data base
func (s *State) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error) {
	batches, err := s.GetLastNBatches(ctx, 1, dbTx)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, ErrNotFound
	}
	return batches[0], nil
}

// GetBatchTimestamp returns the batch timestamp.
//
//	   for >= etrog is stored on virtual_batch.batch_timestamp
//		  previous batches is stored on batch.timestamp
func (s *State) GetBatchTimestamp(ctx context.Context, batchNumber uint64, forcedForkId *uint64, dbTx pgx.Tx) (*time.Time, error) {
	var forkid uint64
	if forcedForkId != nil {
		forkid = *forcedForkId
	} else {
		forkid = s.GetForkIDByBatchNumber(batchNumber)
	}
	batchTimestamp, virtualTimestamp, err := s.GetRawBatchTimestamps(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if forkid >= FORKID_ETROG {
		return virtualTimestamp, nil
	}
	return batchTimestamp, nil
}

// GetL1InfoTreeDataFromBatchL2Data returns a map with the L1InfoTreeData used in the L2 blocks included in the batchL2Data, the last L1InfoRoot used and the highest globalExitRoot used in the batch
func (s *State) GetL1InfoTreeDataFromBatchL2Data(ctx context.Context, batchL2Data []byte, dbTx pgx.Tx) (map[uint32]L1DataV2, common.Hash, common.Hash, error) {
	batchRaw, err := DecodeBatchV2(batchL2Data)
	if err != nil {
		return nil, ZeroHash, ZeroHash, err
	}
	if len(batchRaw.Blocks) == 0 {
		return map[uint32]L1DataV2{}, ZeroHash, ZeroHash, nil
	}

	l1InfoTreeData := map[uint32]L1DataV2{}
	maxIndex := findMax(batchRaw.Blocks)
	l1InfoTreeExitRoot, err := s.GetL1InfoRootLeafByIndex(ctx, maxIndex, dbTx)
	if err != nil {
		return nil, ZeroHash, ZeroHash, err
	}
	maxGER := l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot
	if maxIndex == 0 {
		maxGER = ZeroHash
	}

	l1InfoRoot := l1InfoTreeExitRoot.L1InfoTreeRoot
	for _, l2blockRaw := range batchRaw.Blocks {
		// Index 0 is a special case, it means that the block is not changing GlobalExitRoot.
		// it must not be included in l1InfoTreeData. If all index are 0 L1InfoRoot == ZeroHash
		if l2blockRaw.IndexL1InfoTree > 0 {
			_, found := l1InfoTreeData[l2blockRaw.IndexL1InfoTree]
			if !found {
				l1InfoTreeExitRootStorageEntry, err := s.GetL1InfoRootLeafByIndex(ctx, l2blockRaw.IndexL1InfoTree, dbTx)
				if err != nil {
					return nil, l1InfoRoot, maxGER, err
				}

				l1Data := L1DataV2{
					GlobalExitRoot: l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.GlobalExitRoot.GlobalExitRoot,
					BlockHashL1:    l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.PreviousBlockHash,
					MinTimestamp:   uint64(l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.GlobalExitRoot.Timestamp.Unix()),
				}

				l1InfoTreeData[l2blockRaw.IndexL1InfoTree] = l1Data
			}
		}
	}

	return l1InfoTreeData, l1InfoRoot, maxGER, nil
}

func findMax(blocks []L2BlockRaw) uint32 {
	maxIndex := blocks[0].IndexL1InfoTree
	for _, b := range blocks {
		if b.IndexL1InfoTree > maxIndex {
			maxIndex = b.IndexL1InfoTree
		}
	}
	return maxIndex
}

var mockL1InfoRoot = common.HexToHash(MockL1InfoRootHex)

// GetMockL1InfoRoot returns an instance of common.Hash set
// with the value provided by the const MockL1InfoRootHex
func GetMockL1InfoRoot() common.Hash {
	return mockL1InfoRoot
}
