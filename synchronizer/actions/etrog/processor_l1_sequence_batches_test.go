package etrog

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	syncMocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	hashExamplesValues = []string{"0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92",
		"0x9c8fa7ce2e197f9f1b3c30de9f93de3c1cb290e6c118a18446f47a9e1364c3ab",
		"0x896cfc0684057d0560e950dee352189528167f4663609678d19c7a506a03fe4e",
		"0xde6d2dac4b6e0cb39ed1924db533558a23e5c56ab60fadac8c7d21e7eceb121a",
		"0x9883711e78d02992ac1bd6f19de3bf7bb3f926742d4601632da23525e33f8555"}

	addrExampleValues = []string{"0x8dAF17A20c9DBA35f005b6324F493785D239719d",
		"0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e",
		"0x5FbDB2315678afecb367f032d93F642f64180aa3",
		"0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"}
)

type mocksEtrogProcessorL1 struct {
	Etherman             *mock_syncinterfaces.EthermanFullInterface
	State                *mock_syncinterfaces.StateFullInterface
	Synchronizer         *mock_syncinterfaces.SynchronizerFullInterface
	DbTx                 *syncMocks.DbTxMock
	TimeProvider         *syncCommon.MockTimerProvider
	CriticalErrorHandler *mock_syncinterfaces.CriticalErrorHandler
}

func createMocks(t *testing.T) *mocksEtrogProcessorL1 {
	mocks := &mocksEtrogProcessorL1{
		Etherman:             mock_syncinterfaces.NewEthermanFullInterface(t),
		State:                mock_syncinterfaces.NewStateFullInterface(t),
		Synchronizer:         mock_syncinterfaces.NewSynchronizerFullInterface(t),
		DbTx:                 syncMocks.NewDbTxMock(t),
		TimeProvider:         &syncCommon.MockTimerProvider{},
		CriticalErrorHandler: mock_syncinterfaces.NewCriticalErrorHandler(t),
	}
	return mocks
}

func createSUT(mocks *mocksEtrogProcessorL1) *ProcessorL1SequenceBatchesEtrog {
	return NewProcessorL1SequenceBatches(mocks.State, mocks.Synchronizer,
		mocks.TimeProvider, mocks.CriticalErrorHandler)
}

func TestL1SequenceBatchesNoData(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	err := sut.Process(ctx, etherman.Order{}, nil, mocks.DbTx)
	require.ErrorIs(t, err, actions.ErrInvalidParams)
}

func TestL1SequenceBatchesWrongOrder(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	l1Block := etherman.Block{
		SequencedBatches: [][]etherman.SequencedBatch{},
	}
	err := sut.Process(ctx, etherman.Order{Pos: 1}, &l1Block, mocks.DbTx)
	require.Error(t, err)
}

func TestL1SequenceBatchesPermissionlessNewBatchSequenced(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	batch := newStateBatch(3)
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	expectationsPreExecution(t, mocks, ctx, batch, state.ErrNotFound)
	executionResponse := newProcessBatchResponseV2(batch)
	expectationsProcessAndStoreClosedBatchV2(t, mocks, ctx, executionResponse, nil)
	expectationsAddSequencedBatch(t, mocks, ctx, executionResponse)
	mocks.Synchronizer.EXPECT().PendingFlushID(mock.Anything, mock.Anything)
	err := sut.Process(ctx, etherman.Order{Pos: 1}, newL1Block(mocks, batch, l1InfoRoot), mocks.DbTx)
	require.NoError(t, err)
}

func TestL1SequenceBatchesTrustedBatchSequencedThatAlreadyExistsHappyPath(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	batch := newStateBatch(3)
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	l1Block := newL1Block(mocks, batch, l1InfoRoot)
	expectationsPreExecution(t, mocks, ctx, batch, nil)
	executionResponse := newProcessBatchResponseV2(batch)
	expectationsForExecution(t, mocks, ctx, l1Block.SequencedBatches[1][0], l1Block.ReceivedAt, executionResponse)
	mocks.State.EXPECT().AddAccumulatedInputHash(ctx, executionResponse.NewBatchNum, common.BytesToHash(executionResponse.NewAccInputHash), mocks.DbTx).Return(nil)
	expectationsAddSequencedBatch(t, mocks, ctx, executionResponse)
	err := sut.Process(ctx, etherman.Order{Pos: 1}, l1Block, mocks.DbTx)
	require.NoError(t, err)
}

func TestL1SequenceBatchesPermissionlessBatchSequencedThatAlreadyExistsHappyPath(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	batch := newStateBatch(3)
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	l1Block := newL1Block(mocks, batch, l1InfoRoot)
	expectationsPreExecution(t, mocks, ctx, batch, nil)
	executionResponse := newProcessBatchResponseV2(batch)
	expectationsForExecution(t, mocks, ctx, l1Block.SequencedBatches[1][0], l1Block.ReceivedAt, executionResponse)
	mocks.State.EXPECT().AddAccumulatedInputHash(ctx, executionResponse.NewBatchNum, common.BytesToHash(executionResponse.NewAccInputHash), mocks.DbTx).Return(nil)
	expectationsAddSequencedBatch(t, mocks, ctx, executionResponse)
	err := sut.Process(ctx, etherman.Order{Pos: 1}, l1Block, mocks.DbTx)
	require.NoError(t, err)
}

// CASE: A permissionless process a L1 sequenced batch that already is in state (presumably synced from Trusted)
// - Execute it
// - Check if match state batch
// - Don't match -> Reorg Pool and reset trusted state
// - Reprocess again as a new batch
func TestL1SequenceBatchesPermissionlessBatchSequencedThatAlreadyExistsMismatch(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	batch := newStateBatch(3)
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	l1Block := newL1Block(mocks, batch, l1InfoRoot)
	expectationsPreExecution(t, mocks, ctx, batch, nil)
	executionResponse := newProcessBatchResponseV2(batch)
	executionResponse.NewStateRoot = common.HexToHash(hashExamplesValues[2]).Bytes()
	expectationsForExecution(t, mocks, ctx, l1Block.SequencedBatches[1][0], l1Block.ReceivedAt, executionResponse)
	mocks.State.EXPECT().AddAccumulatedInputHash(ctx, executionResponse.NewBatchNum, common.BytesToHash(executionResponse.NewAccInputHash), mocks.DbTx).Return(nil)
	mocks.Synchronizer.EXPECT().IsTrustedSequencer().Return(false)
	mocks.State.EXPECT().AddTrustedReorg(ctx, mock.Anything, mocks.DbTx).Return(nil)
	mocks.State.EXPECT().ResetTrustedState(ctx, batch.BatchNumber-1, mocks.DbTx).Return(nil)
	mocks.Synchronizer.EXPECT().CleanTrustedState()

	// Reexecute it as a new batch
	expectationsProcessAndStoreClosedBatchV2(t, mocks, ctx, executionResponse, nil)
	expectationsAddSequencedBatch(t, mocks, ctx, executionResponse)
	mocks.Synchronizer.EXPECT().PendingFlushID(mock.Anything, mock.Anything)
	err := sut.Process(ctx, etherman.Order{Pos: 1}, l1Block, mocks.DbTx)
	require.NoError(t, err)
}

type CriticalErrorHandlerPanic struct {
}

func (c CriticalErrorHandlerPanic) CriticalError(ctx context.Context, err error) {
	panic("CriticalError")
}

// CASE: A TRUSTED SYNCHRONIZER process a L1 sequenced batch that already is in state but it doesnt match with the trusted State
// - Execute it
// - Check if match state batch
// - Don't match -> HALT
func TestL1SequenceBatchesTrustedBatchSequencedThatAlreadyExistsMismatch(t *testing.T) {
	mocks := createMocks(t)
	CriticalErrorHandlerPanic := CriticalErrorHandlerPanic{}
	sut := NewProcessorL1SequenceBatches(mocks.State, mocks.Synchronizer,
		mocks.TimeProvider, CriticalErrorHandlerPanic)
	ctx := context.Background()
	batch := newStateBatch(3)
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	l1Block := newL1Block(mocks, batch, l1InfoRoot)
	expectationsPreExecution(t, mocks, ctx, batch, nil)
	executionResponse := newProcessBatchResponseV2(batch)
	executionResponse.NewStateRoot = common.HexToHash(hashExamplesValues[2]).Bytes()
	expectationsForExecution(t, mocks, ctx, l1Block.SequencedBatches[1][0], l1Block.ReceivedAt, executionResponse)
	mocks.State.EXPECT().AddAccumulatedInputHash(ctx, executionResponse.NewBatchNum, common.BytesToHash(executionResponse.NewAccInputHash), mocks.DbTx).Return(nil)
	mocks.Synchronizer.EXPECT().IsTrustedSequencer().Return(true)

	// CriticalError call in a real implementation is a blocking call, in the test is going to panic
	assertPanic(t, func() { sut.Process(ctx, etherman.Order{Pos: 1}, l1Block, mocks.DbTx) }) //nolint
}

func TestL1SequenceForcedBatchesNum1TrustedBatch(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	batch := newStateBatch(3)
	forcedTime := mocks.TimeProvider.Now()
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	forcedGlobalExitRoot := common.HexToHash(hashExamplesValues[1])
	forcedBlockHash := common.HexToHash(hashExamplesValues[2])
	sequencedForcedBatch := newForcedSequenceBatch(batch, l1InfoRoot, forcedTime, forcedGlobalExitRoot, forcedBlockHash)

	l1Block := newComposedL1Block(mocks, sequencedForcedBatch, l1InfoRoot)

	mocks.State.EXPECT().GetNextForcedBatches(ctx, int(1), mocks.DbTx).Return([]state.ForcedBatch{
		{
			BlockNumber:       32,
			ForcedBatchNumber: 4,
			Sequencer:         common.HexToAddress(addrExampleValues[0]),
			GlobalExitRoot:    forcedGlobalExitRoot,
			RawTxsData:        []byte{},
			ForcedAt:          forcedTime,
		},
	}, nil)
	expectationsPreExecution(t, mocks, ctx, batch, state.ErrNotFound)

	executionResponse := newProcessBatchResponseV2(batch)
	executionResponse.NewStateRoot = common.HexToHash(hashExamplesValues[2]).Bytes()

	expectationsProcessAndStoreClosedBatchV2(t, mocks, ctx, executionResponse, nil)
	expectationsAddSequencedBatch(t, mocks, ctx, executionResponse)
	mocks.Synchronizer.EXPECT().PendingFlushID(mock.Anything, mock.Anything)

	err := sut.Process(ctx, etherman.Order{Pos: 1}, l1Block, mocks.DbTx)
	require.NoError(t, err)
}

// --------------------- Helper functions ----------------------------------------------------------------------------------------------------

func expectationsPreExecution(t *testing.T, mocks *mocksEtrogProcessorL1, ctx context.Context, trustedBatch *state.Batch, responseError error) {
	mocks.State.EXPECT().GetL1InfoTreeDataFromBatchL2Data(ctx, mock.Anything, mocks.DbTx).Return(map[uint32]state.L1DataV2{}, state.ZeroHash, state.ZeroHash, nil).Maybe()
	mocks.State.EXPECT().GetBatchByNumber(ctx, trustedBatch.BatchNumber, mocks.DbTx).Return(trustedBatch, responseError)
}

func expectationsAddSequencedBatch(t *testing.T, mocks *mocksEtrogProcessorL1, ctx context.Context, response *executor.ProcessBatchResponseV2) {
	mocks.State.EXPECT().AddVirtualBatch(ctx, mock.Anything, mocks.DbTx).Return(nil)
	mocks.State.EXPECT().AddSequence(ctx, state.Sequence{FromBatchNumber: 3, ToBatchNumber: 3}, mocks.DbTx).Return(nil)
}

func expectationsProcessAndStoreClosedBatchV2(t *testing.T, mocks *mocksEtrogProcessorL1, ctx context.Context, response *executor.ProcessBatchResponseV2, responseError error) {
	newStateRoot := common.BytesToHash(response.NewStateRoot)
	mocks.State.EXPECT().ProcessAndStoreClosedBatchV2(ctx, mock.Anything, mocks.DbTx, mock.Anything).Return(newStateRoot, response.FlushId, response.ProverId, responseError)
}

func expectationsForExecution(t *testing.T, mocks *mocksEtrogProcessorL1, ctx context.Context, sequencedBatch etherman.SequencedBatch, timestampLimit time.Time, response *executor.ProcessBatchResponseV2) {
	mocks.State.EXPECT().ExecuteBatchV2(ctx,
		mock.Anything, *sequencedBatch.L1InfoRoot, mock.Anything, timestampLimit, false,
		uint32(1), (*common.Hash)(nil), mocks.DbTx).Return(response, nil)
}

func newProcessBatchResponseV2(batch *state.Batch) *executor.ProcessBatchResponseV2 {
	return &executor.ProcessBatchResponseV2{
		NewBatchNum:     batch.BatchNumber,
		NewAccInputHash: batch.AccInputHash[:],
		NewStateRoot:    batch.StateRoot[:],
		FlushId:         uint64(1234),
		ProverId:        "prover-id",
	}
}

func newStateBatch(number uint64) *state.Batch {
	return &state.Batch{
		BatchNumber: number,
		StateRoot:   common.HexToHash(hashExamplesValues[3]),
		Coinbase:    common.HexToAddress(addrExampleValues[0]),
	}
}

func newForcedSequenceBatch(batch *state.Batch, l1InfoRoot common.Hash, forcedTimestamp time.Time, forcedGlobalExitRoot, forcedBlockHashL1 common.Hash) *etherman.SequencedBatch {
	return &etherman.SequencedBatch{
		BatchNumber:   batch.BatchNumber,
		L1InfoRoot:    &l1InfoRoot,
		TxHash:        state.HashByteArray(batch.BatchL2Data),
		Coinbase:      batch.Coinbase,
		SequencerAddr: common.HexToAddress(addrExampleValues[0]),
		PolygonRollupBaseEtrogBatchData: &etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions:         []byte{},
			ForcedTimestamp:      uint64(forcedTimestamp.Unix()),
			ForcedGlobalExitRoot: forcedGlobalExitRoot,
			ForcedBlockHashL1:    forcedBlockHashL1,
		},
	}
}

func newL1Block(mocks *mocksEtrogProcessorL1, batch *state.Batch, l1InfoRoot common.Hash) *etherman.Block {
	sbatch := etherman.SequencedBatch{
		BatchNumber:   batch.BatchNumber,
		L1InfoRoot:    &l1InfoRoot,
		TxHash:        state.HashByteArray(batch.BatchL2Data),
		Coinbase:      batch.Coinbase,
		SequencerAddr: common.HexToAddress(addrExampleValues[0]),
		PolygonRollupBaseEtrogBatchData: &etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions: []byte{},
		},
	}

	return newComposedL1Block(mocks, &sbatch, l1InfoRoot)
}

func newComposedL1Block(mocks *mocksEtrogProcessorL1, forcedBatch *etherman.SequencedBatch, l1InfoRoot common.Hash) *etherman.Block {
	l1Block := etherman.Block{
		BlockNumber:      123,
		ReceivedAt:       mocks.TimeProvider.Now(),
		SequencedBatches: [][]etherman.SequencedBatch{},
	}
	l1Block.SequencedBatches = append(l1Block.SequencedBatches, []etherman.SequencedBatch{})
	l1Block.SequencedBatches = append(l1Block.SequencedBatches, []etherman.SequencedBatch{
		*forcedBatch,
	})
	return &l1Block
}

// https://stackoverflow.com/questions/31595791/how-to-test-panics
func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
