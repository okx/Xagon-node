package l2_sync_etrog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrNotImplemented is returned when a method is not implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrFailExecuteBatch is returned when the batch is not executed correctly
	ErrFailExecuteBatch = errors.New("fail execute batch")
	// ErrCriticalClosedBatchDontContainExpectedData is returnted when try to close a batch that is already close but data doesnt match
	ErrCriticalClosedBatchDontContainExpectedData = errors.New("when closing the batch, the batch is already close, but  the data on state doesnt match the expected")
	// ErrCantReprocessBatchMissingPreviousStateBatch can't reprocess a divergent batch because is missing previous state batch
	ErrCantReprocessBatchMissingPreviousStateBatch = errors.New("cant reprocess batch because is missing previous state batch")
)

// StateInterface contains the methods required to interact with the state.
type StateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	UpdateWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, string, error)
	StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) (common.Hash, error)
	GetL1InfoTreeDataFromBatchL2Data(ctx context.Context, batchL2Data []byte, dbTx pgx.Tx) (map[uint32]state.L1DataV2, common.Hash, common.Hash, error)
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
}

// SyncTrustedBatchExecutorForEtrog is the implementation of the SyncTrustedStateBatchExecutorSteps that
// have the functions to sync a fullBatch, incrementalBatch and reprocessBatch
type SyncTrustedBatchExecutorForEtrog struct {
	state StateInterface
	sync  syncinterfaces.SynchronizerFlushIDManager
}

// NewSyncTrustedBatchExecutorForEtrog creates a new SyncTrustedBatchExecutorForEtrog
func NewSyncTrustedBatchExecutorForEtrog(stateBatchExecutor StateInterface,
	sync syncinterfaces.SynchronizerFlushIDManager) *SyncTrustedBatchExecutorForEtrog {
	return &SyncTrustedBatchExecutorForEtrog{
		state: stateBatchExecutor,
		sync:  sync,
	}
}

// NothingProcess process a batch that is already on database and no new L2batchData, so it is not going to be processed again.
// Maybe it needs to be close
func (b *SyncTrustedBatchExecutorForEtrog) NothingProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	isEqual, strResult := l2_shared.AreEqualStateBatchAndTrustedBatch(data.StateBatch, data.TrustedBatch, l2_shared.CMP_BATCH_IGNORE_TSTAMP+l2_shared.CMP_BATCH_IGNORE_WIP)
	if !isEqual {
		log.Warnf("%s Nothing new to process but the TrustedBatch differ: %s. Forcing a reprocess", data.DebugPrefix, strResult)
		if data.StateBatch.WIP {
			if data.PreviousStateBatch != nil {
				data.OldAccInputHash = data.PreviousStateBatch.AccInputHash
				data.OldStateRoot = data.PreviousStateBatch.StateRoot
				return b.ReProcess(ctx, data, dbTx)
			} else {
				log.Warnf("%s PreviousStateBatch is nil. Can't reprocess", data.DebugPrefix)
				return nil, ErrCantReprocessBatchMissingPreviousStateBatch
			}
		} else {
			log.Warnf("%s StateBatch is not WIP. Can't reprocess", data.DebugPrefix)
			return nil, ErrCriticalClosedBatchDontContainExpectedData
		}
	}
	res := l2_shared.NewProcessResponse()
	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err := b.CloseBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Error("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
		data.StateBatch.WIP = false
		res.UpdateCurrentBatch(data.StateBatch)
	}

	return &res, nil
}

// CreateEmptyBatch create a new empty batch (no batchL2Data and WIP)
func (b *SyncTrustedBatchExecutorForEtrog) CreateEmptyBatch(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	log.Debugf("%s The Batch is a empty (batchl2data=0 bytes), so just creating a DB entry", data.DebugPrefix)
	err := b.openBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
	if err != nil {
		log.Errorf("%s error openning batch. Error: %v", data.DebugPrefix, err)
		return nil, err
	}
	if data.BatchMustBeClosed {
		log.Infof("%s Closing empty batch (no execution)", data.DebugPrefix)
		err = b.CloseBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Error("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	} else {
		log.Debugf("%s updateWIPBatch", data.DebugPrefix)
		err = b.updateWIPBatch(ctx, data, data.TrustedBatch.StateRoot, dbTx)
		if err != nil {
			log.Errorf("%s error updateWIPBatch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}

	res := l2_shared.NewProcessResponse()
	stateBatch := syncCommon.RpcBatchToStateBatch(data.TrustedBatch)
	res.UpdateCurrentBatch(stateBatch)
	return &res, nil
}

// FullProcess process a batch that is not on database, so is the first time we process it
func (b *SyncTrustedBatchExecutorForEtrog) FullProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	log.Debugf("%s FullProcess", data.DebugPrefix)
	if len(data.TrustedBatch.BatchL2Data) == 0 {
		data.DebugPrefix += " (emptyBatch) "
		return b.CreateEmptyBatch(ctx, data, dbTx)
	}
	err := b.openBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
	if err != nil {
		log.Errorf("%s error openning batch. Error: %v", data.DebugPrefix, err)
		return nil, err
	}

	leaves, l1InfoRoot, _, err := b.state.GetL1InfoTreeDataFromBatchL2Data(ctx, data.TrustedBatch.BatchL2Data, dbTx)
	if err != nil {
		log.Errorf("%s error getting GetL1InfoTreeDataFromBatchL2Data: %v. Error:%w", data.DebugPrefix, l1InfoRoot, err)
		return nil, err
	}
	debugStr := data.DebugPrefix
	processBatchResp, err := b.processAndStoreTxs(ctx, b.getProcessRequest(data, leaves, l1InfoRoot), dbTx, debugStr)
	if err != nil {
		log.Error("%s error procesingAndStoringTxs. Error: ", debugStr, err)
		return nil, err
	}

	err = batchResultSanityCheck(data, processBatchResp, debugStr)
	if err != nil {
		log.Errorf("%s error batchResultSanityCheck. Error: %s", data.DebugPrefix, err.Error())
		return nil, err
	}

	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err = b.CloseBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Error("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	} else {
		log.Debugf("%s updateWIPBatch", data.DebugPrefix)
		err = b.updateWIPBatch(ctx, data, processBatchResp.NewStateRoot, dbTx)
		if err != nil {
			log.Errorf("%s error updateWIPBatch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}

	resultBatch, err := b.state.GetBatchByNumber(ctx, uint64(data.TrustedBatch.Number), dbTx)
	if err != nil {
		log.Error("%s error getting batch. Error: ", data.DebugPrefix, err)
		return nil, err
	}
	res := l2_shared.NewProcessResponse()
	res.UpdateCurrentBatchWithExecutionResult(resultBatch, processBatchResp)
	return &res, nil
}

// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
func (b *SyncTrustedBatchExecutorForEtrog) IncrementalProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	var err error
	if data == nil || data.TrustedBatch == nil || data.StateBatch == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if err := checkThatL2DataIsIncremental(data); err != nil {
		log.Errorf("%s error checkThatL2DataIsIncremental. Error: %v", data.DebugPrefix, err)
		return nil, err
	}

	PartialBatchL2Data, err := b.composePartialBatch(data.StateBatch, data.TrustedBatch)
	if err != nil {
		log.Errorf("%s error composePartialBatch batch Error:%w", data.DebugPrefix, err)
		return nil, err
	}

	leaves, l1InfoRoot, _, err := b.state.GetL1InfoTreeDataFromBatchL2Data(ctx, PartialBatchL2Data, dbTx)
	if err != nil {
		log.Errorf("%s error getting GetL1InfoTreeDataFromBatchL2Data: %v. Error:%w", data.DebugPrefix, l1InfoRoot, err)
		// TODO: Need to refine, depending of the response of GetL1InfoTreeDataFromBatchL2Data
		// if some leaf is missing, we need to resync from L1 to get the missing events and then process again
		return nil, syncinterfaces.ErrMissingSyncFromL1
	}
	debugStr := fmt.Sprintf("%s: Batch %d:", data.Mode, uint64(data.TrustedBatch.Number))
	processReq := b.getProcessRequest(data, leaves, l1InfoRoot)
	processReq.Transactions = PartialBatchL2Data
	processBatchResp, err := b.processAndStoreTxs(ctx, processReq, dbTx, debugStr)
	if err != nil {
		log.Errorf("%s error procesingAndStoringTxs. Error: ", data.DebugPrefix, err)
		return nil, err
	}

	err = batchResultSanityCheck(data, processBatchResp, debugStr)
	if err != nil {
		log.Errorf("%s error batchResultSanityCheck. Error: %s", data.DebugPrefix, err.Error())
		return nil, err
	}
	log.Debugf("%s updateWIPBatch ", data.DebugPrefix)
	err = b.updateWIPBatch(ctx, data, processBatchResp.NewStateRoot, dbTx)
	if err != nil {
		log.Errorf("%s error updateWIPBatch. Error: ", data.DebugPrefix, err)
		return nil, err
	}

	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err = b.CloseBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Errorf("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}

	updatedBatch := *data.StateBatch
	updatedBatch.LocalExitRoot = data.TrustedBatch.LocalExitRoot
	updatedBatch.AccInputHash = data.TrustedBatch.AccInputHash
	updatedBatch.GlobalExitRoot = data.TrustedBatch.GlobalExitRoot
	updatedBatch.BatchL2Data = data.TrustedBatch.BatchL2Data
	updatedBatch.WIP = !data.BatchMustBeClosed

	res := l2_shared.NewProcessResponse()
	res.UpdateCurrentBatchWithExecutionResult(&updatedBatch, processBatchResp)
	return &res, nil
}

func (b *SyncTrustedBatchExecutorForEtrog) updateWIPBatch(ctx context.Context, data *l2_shared.ProcessData, NewStateRoot common.Hash, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:    data.BatchNumber,
		StateRoot:      NewStateRoot,
		LocalExitRoot:  data.TrustedBatch.LocalExitRoot,
		BatchL2Data:    data.TrustedBatch.BatchL2Data,
		AccInputHash:   data.TrustedBatch.AccInputHash,
		GlobalExitRoot: data.TrustedBatch.GlobalExitRoot,
	}

	err := b.state.UpdateWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		log.Errorf("%s error UpdateWIPBatch. Error: ", data.DebugPrefix, err)
		return err
	}
	return err
}

// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
func (b *SyncTrustedBatchExecutorForEtrog) ReProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	log.Warnf("%s needs to be reprocessed! deleting batches from this batch, because it was partially processed but the intermediary stateRoot is lost", data.DebugPrefix)
	// Check that there are no VirtualBatches neither VerifiedBatches that are newer than this batch
	lastVirtualBatchNum, err := b.state.GetLastVirtualBatchNum(ctx, dbTx)
	if err != nil {
		log.Errorf("%s error getting lastVirtualBatchNum. Error: %v", data.DebugPrefix, err)
		return nil, err
	}
	if lastVirtualBatchNum >= uint64(data.TrustedBatch.Number) {
		log.Errorf("%s there are newer or equal virtualBatches than this batch. Can't reprocess because then will delete a virtualBatch", data.DebugPrefix)
		return nil, syncinterfaces.ErrMissingSyncFromL1
	}
	err = b.state.ResetTrustedState(ctx, uint64(data.TrustedBatch.Number)-1, dbTx)
	if err != nil {
		log.Warnf("%s error deleting batches from this batch: %v", data.DebugPrefix, err)
		return nil, err
	}
	// From this point is like a new trusted batch
	return b.FullProcess(ctx, data, dbTx)
}

func batchResultSanityCheck(data *l2_shared.ProcessData, processBatchResp *state.ProcessBatchResponse, debugStr string) error {
	if processBatchResp == nil {
		return nil
	}
	if processBatchResp.NewStateRoot == state.ZeroHash {
		return fmt.Errorf("%s processBatchResp.NewStateRoot is ZeroHash. Err: %w", debugStr, l2_shared.ErrFatalBatchDesynchronized)
	}
	if processBatchResp.NewStateRoot != data.TrustedBatch.StateRoot {
		return fmt.Errorf("%s processBatchResp.NewStateRoot(%s) != data.TrustedBatch.StateRoot(%s). Err: %w", debugStr,
			processBatchResp.NewStateRoot.String(), data.TrustedBatch.StateRoot.String(), l2_shared.ErrFatalBatchDesynchronized)
	}
	if processBatchResp.NewLocalExitRoot != data.TrustedBatch.LocalExitRoot {
		return fmt.Errorf("%s processBatchResp.NewLocalExitRoot(%s) != data.StateBatch.LocalExitRoot(%s). Err: %w", debugStr,
			processBatchResp.NewLocalExitRoot.String(), data.TrustedBatch.LocalExitRoot.String(), l2_shared.ErrFatalBatchDesynchronized)
	}
	// We can't check AccInputHash because we dont have timeLimit neither L1InfoRoot used to create the batch
	// is going to be update from L1
	// if processBatchResp.NewAccInputHash != data.TrustedBatch.AccInputHash {
	// 	return fmt.Errorf("%s processBatchResp.	if processBatchResp.NewAccInputHash(%s) != data.TrustedBatch.AccInputHash(%s). Err:%w", debugStr,
	// 		processBatchResp.NewAccInputHash.String(), data.TrustedBatch.AccInputHash.String(), ErrNotExpectedBathResult)
	// }
	return nil
}

// CloseBatch close a batch
func (b *SyncTrustedBatchExecutorForEtrog) CloseBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx, debugStr string) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   uint64(trustedBatch.Number),
		StateRoot:     trustedBatch.StateRoot,
		LocalExitRoot: trustedBatch.LocalExitRoot,
		BatchL2Data:   trustedBatch.BatchL2Data,
		AccInputHash:  trustedBatch.AccInputHash,
		ClosingReason: state.SyncL2TrustedBatchClosingReason,
	}
	log.Debugf("%s closing batch %v", debugStr, trustedBatch.Number)
	// This update SET state_root = $1, local_exit_root = $2, acc_input_hash = $3, raw_txs_data = $4, batch_resources = $5, closing_reason = $6, wip = FALSE
	if err := b.state.CloseBatch(ctx, receipt, dbTx); err != nil {
		// This is a workaround to avoid closing a batch that was already closed
		if err.Error() != state.ErrBatchAlreadyClosed.Error() {
			log.Errorf("%s error closing batch %d", debugStr, trustedBatch.Number)
			return err
		} else {
			log.Warnf("%s CASE 02: the batch [%d] looks like were not close but in STATE was closed", debugStr, trustedBatch.Number)
			// Check that the fields have the right values
			dbBatch, err := b.state.GetBatchByNumber(ctx, uint64(trustedBatch.Number), dbTx)
			if err != nil {
				log.Errorf("%s error getting local batch %d", debugStr, trustedBatch.Number)
				return err
			}
			equals, str := l2_shared.AreEqualStateBatchAndTrustedBatch(dbBatch, trustedBatch, l2_shared.CMP_BATCH_IGNORE_TSTAMP)
			if !equals {
				// This is a situation impossible to reach!, if it happens we halt sync and we need to develop a recovery process
				err := fmt.Errorf("%s the batch data on state doesnt match the expected (%s) error:%w", debugStr, str, ErrCriticalClosedBatchDontContainExpectedData)
				log.Warnf(err.Error())
				return err
			}
		}
	}
	return nil
}

func (b *SyncTrustedBatchExecutorForEtrog) openBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx, debugStr string) error {
	log.Debugf("%s Opening batch %d", debugStr, trustedBatch.Number)
	var batchL2Data []byte = trustedBatch.BatchL2Data
	processCtx := state.ProcessingContext{
		BatchNumber: uint64(trustedBatch.Number),
		Coinbase:    common.HexToAddress(trustedBatch.Coinbase.String()),
		// Instead of using trustedBatch.Timestamp use now, because the prevBatch could have a newer timestamp because
		// use the tstamp of the L1Block where is the virtualization event
		Timestamp:      time.Now(),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	if trustedBatch.ForcedBatchNumber != nil {
		fb := uint64(*trustedBatch.ForcedBatchNumber)
		processCtx.ForcedBatchNum = &fb
	}
	err := b.state.OpenBatch(ctx, processCtx, dbTx)
	if err != nil {
		log.Error("%s error opening batch: ", debugStr, trustedBatch.Number)
		return err
	}
	return nil
}

func (b *SyncTrustedBatchExecutorForEtrog) processAndStoreTxs(ctx context.Context, request state.ProcessRequest, dbTx pgx.Tx, debugPrefix string) (*state.ProcessBatchResponse, error) {
	if request.OldStateRoot == state.ZeroHash {
		log.Warnf("%s Processing batch with oldStateRoot == zero....", debugPrefix)
	}
	processBatchResp, _, err := b.state.ProcessBatchV2(ctx, request, true)
	if err != nil {
		log.Errorf("%s error processing sequencer batch for batch: %v error:%v ", debugPrefix, request.BatchNumber, err)
		return nil, err
	}
	b.sync.PendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("%s Storing %d blocks for batch %v", debugPrefix, len(processBatchResp.BlockResponses), request.BatchNumber)
	if processBatchResp.IsExecutorLevelError {
		log.Warnf("%s executorLevelError detected. Avoid store txs...", debugPrefix)
		return nil, fmt.Errorf("%s executorLevelError detected err: %w", debugPrefix, ErrFailExecuteBatch)
	} else if processBatchResp.IsRomOOCError {
		log.Warnf("%s romOOCError detected. Avoid store txs...", debugPrefix)
		return nil, fmt.Errorf("%s romOOCError detected.err: %w", debugPrefix, ErrFailExecuteBatch)
	}
	for _, block := range processBatchResp.BlockResponses {
		log.Debugf("%s Storing trusted tx %d", debugPrefix, block.BlockNumber)
		if _, err = b.state.StoreL2Block(ctx, request.BatchNumber, block, nil, dbTx); err != nil {
			newErr := fmt.Errorf("%s failed to store l2block: %v  err:%w", debugPrefix, block.BlockNumber, err)
			log.Error(newErr.Error())
			return nil, newErr
		}
	}
	log.Infof("%s Batch %v: batchl2data len:%d processed and stored: %s oldStateRoot: %s -> newStateRoot:%s", debugPrefix, request.BatchNumber, len(request.Transactions), getResponseInfo(processBatchResp),
		request.OldStateRoot.String(), processBatchResp.NewStateRoot.String())
	return processBatchResp, nil
}

func getResponseInfo(response *state.ProcessBatchResponse) string {
	if len(response.BlockResponses) == 0 {
		return "no blocks, no txs"
	}
	minBlock := response.BlockResponses[0].BlockNumber
	maxBlock := response.BlockResponses[len(response.BlockResponses)-1].BlockNumber
	totalTx := 0
	for _, block := range response.BlockResponses {
		totalTx += len(block.TransactionResponses)
	}
	return fmt.Sprintf(" l2block[%v-%v] txs[%v]", minBlock, maxBlock, totalTx)
}

func (b *SyncTrustedBatchExecutorForEtrog) getProcessRequest(data *l2_shared.ProcessData, l1InfoTreeLeafs map[uint32]state.L1DataV2, l1InfoTreeRoot common.Hash) state.ProcessRequest {
	request := state.ProcessRequest{
		BatchNumber:             uint64(data.TrustedBatch.Number),
		OldStateRoot:            data.OldStateRoot,
		OldAccInputHash:         data.OldAccInputHash,
		Coinbase:                common.HexToAddress(data.TrustedBatch.Coinbase.String()),
		L1InfoRoot_V2:           l1InfoTreeRoot,
		L1InfoTreeData_V2:       l1InfoTreeLeafs,
		TimestampLimit_V2:       uint64(data.TrustedBatch.Timestamp),
		Transactions:            data.TrustedBatch.BatchL2Data,
		ForkID:                  b.state.GetForkIDByBatchNumber(uint64(data.TrustedBatch.Number)),
		SkipVerifyL1InfoRoot_V2: true,
	}
	return request
}

func checkThatL2DataIsIncremental(data *l2_shared.ProcessData) error {
	newDataFlag, err := l2_shared.ThereAreNewBatchL2Data(data.StateBatch.BatchL2Data, data.TrustedBatch.BatchL2Data)
	if err != nil {
		return err
	}
	if !newDataFlag {
		return l2_shared.ErrBatchDataIsNotIncremental
	}
	return nil
}

func (b *SyncTrustedBatchExecutorForEtrog) composePartialBatch(previousBatch *state.Batch, newBatch *types.Batch) ([]byte, error) {
	debugStr := " composePartialBatch: "
	rawPreviousBatch, err := state.DecodeBatchV2(previousBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}
	debugStr += fmt.Sprintf("previousBatch.blocks: %v (%v) ", len(rawPreviousBatch.Blocks), len(previousBatch.BatchL2Data))
	if len(previousBatch.BatchL2Data) >= len(newBatch.BatchL2Data) {
		return nil, fmt.Errorf("previousBatch.BatchL2Data (%d)>=newBatch.BatchL2Data (%d)", len(previousBatch.BatchL2Data), len(newBatch.BatchL2Data))
	}
	newData := newBatch.BatchL2Data[len(previousBatch.BatchL2Data):]
	rawPartialBatch, err := state.DecodeBatchV2(newData)
	if err != nil {
		return nil, err
	}
	debugStr += fmt.Sprintf(" deltaBatch.blocks: %v (%v) ", len(rawPartialBatch.Blocks), len(newData))

	newBatchEncoded, err := state.EncodeBatchV2(rawPartialBatch)
	if err != nil {
		return nil, err
	}
	log.Debug(debugStr)
	return newBatchEncoded, nil
}
