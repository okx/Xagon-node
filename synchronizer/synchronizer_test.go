package synchronizer

import (
	context "context"
	"math"
	"math/big"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync"
	syncMocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	cProverIDExecution             = "PROVER_ID-EXE001"
	ETROG_MODE_FLAG                = true
	RETRIEVE_BATCH_FROM_DB_FLAG    = true
	RETRIEVE_BATCH_FROM_CACHE_FLAG = false
	PROCESS_BATCH_SELECTOR_ENABLED = false
)

type mocks struct {
	Etherman                      *mock_syncinterfaces.EthermanFullInterface
	State                         *mock_syncinterfaces.StateFullInterface
	Pool                          *mock_syncinterfaces.PoolInterface
	EthTxManager                  *mock_syncinterfaces.EthTxManager
	DbTx                          *syncMocks.DbTxMock
	ZKEVMClient                   *mock_syncinterfaces.ZKEVMClientInterface
	zkEVMClientEthereumCompatible *mock_syncinterfaces.ZKEVMClientEthereumCompatibleInterface
	//EventLog     *eventLogMock
}

// Feature #2220 and  #2239: Optimize Trusted state synchronization
//
//	this Check partially point 2: Use previous batch stored in memory to avoid getting from database
func TestGivenPermissionlessNodeWhenSyncronizeAgainSameBatchThenUseTheOneInMemoryInstaeadOfGettingFromDb(t *testing.T) {
	genesis, cfg, m := setupGenericTest(t)
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	syncInterface, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, *genesis, *cfg, false)
	require.NoError(t, err)
	sync, ok := syncInterface.(*ClientSynchronizer)
	require.EqualValues(t, true, ok, "Can't convert to underlaying struct the interface of syncronizer")
	lastBatchNumber := uint64(10)
	batch10With2Tx := createBatch(t, lastBatchNumber, 2, ETROG_MODE_FLAG)
	batch10With3Tx := createBatch(t, lastBatchNumber, 3, ETROG_MODE_FLAG)
	previousBatch09 := createBatch(t, lastBatchNumber-1, 1, ETROG_MODE_FLAG)

	forkIdInterval := state.ForkIDInterval{
		FromBatchNumber: 0,
		ToBatchNumber:   ^uint64(0),
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(7)).Return(&forkIdInterval)
	m.State.EXPECT().GetForkIDByBatchNumber(lastBatchNumber + 1).Return(uint64(7))

	expectedCallsForsyncTrustedState(t, m, sync, nil, batch10With2Tx, previousBatch09, RETRIEVE_BATCH_FROM_DB_FLAG, ETROG_MODE_FLAG)
	// Is the first time that appears this batch, so it need to OpenBatch
	expectedCallsForOpenBatch(t, m, sync, lastBatchNumber)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)
	// Check that all mock expectations are satisfied before continue with next call
	m.checkExpectedCalls(t)

	// This call is going to be a incremental process of the batch using the cache data
	expectedCallsForsyncTrustedState(t, m, sync, batch10With2Tx, batch10With3Tx, previousBatch09, RETRIEVE_BATCH_FROM_CACHE_FLAG, ETROG_MODE_FLAG)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)

	cachedBatch := sync.syncTrustedStateExecutor.GetCachedBatch(uint64(batch10With3Tx.Number))
	require.True(t, cachedBatch != nil)
	require.Equal(t, rpcBatchTostateBatch(batch10With3Tx), cachedBatch)
}

// Feature #2220 and  #2239: Optimize Trusted state synchronization
//
//	this Check partially point 2: Store last batch in memory (CurrentTrustedBatch)
func TestGivenPermissionlessNodeWhenSyncronizeFirstTimeABatchThenStoreItInALocalVar(t *testing.T) {
	genesis, cfg, m := setupGenericTest(t)
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	syncInterface, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, *genesis, *cfg, false)
	require.NoError(t, err)
	sync, ok := syncInterface.(*ClientSynchronizer)
	require.EqualValues(t, true, ok, "Can't convert to underlaying struct the interface of syncronizer")
	lastBatchNumber := uint64(10)
	batch10With1Tx := createBatch(t, lastBatchNumber, 1, ETROG_MODE_FLAG)
	batch10With2Tx := createBatch(t, lastBatchNumber, 2, ETROG_MODE_FLAG)
	previousBatch09 := createBatch(t, lastBatchNumber-1, 1, ETROG_MODE_FLAG)

	forkIdInterval := state.ForkIDInterval{
		FromBatchNumber: 0,
		ToBatchNumber:   ^uint64(0),
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(7)).Return(&forkIdInterval)
	m.State.EXPECT().GetForkIDByBatchNumber(lastBatchNumber + 1).Return(uint64(7))

	// This is a incremental process, permissionless have batch10With1Tx and we add a new block
	// but the cache doesnt have this information so it need to get from db
	expectedCallsForsyncTrustedState(t, m, sync, batch10With1Tx, batch10With2Tx, previousBatch09, RETRIEVE_BATCH_FROM_DB_FLAG, ETROG_MODE_FLAG)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)

	cachedBatch := sync.syncTrustedStateExecutor.GetCachedBatch(uint64(batch10With2Tx.Number))
	require.True(t, cachedBatch != nil)
	require.Equal(t, rpcBatchTostateBatch(batch10With2Tx), cachedBatch)
}

// issue #2220
// TODO: this is running against old sequential L1 sync, need to update to parallel L1 sync.
// but it used a feature that is not implemented in new one that is asking beyond the last block on L1
func TestForcedBatchEtrog(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:    mock_syncinterfaces.NewEthermanFullInterface(t),
		State:       mock_syncinterfaces.NewStateFullInterface(t),
		Pool:        mock_syncinterfaces.NewPoolInterface(t),
		DbTx:        syncMocks.NewDbTxMock(t),
		ZKEVMClient: mock_syncinterfaces.NewZKEVMClientInterface(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		FromBatchNumber: 0,
		ToBatchNumber:   ^uint64(0),
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(7)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(7), nil).
				Maybe()

			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock0, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Times(2)

			n := big.NewInt(rpc.LatestBlockNumber.Int64())
			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader1, nil).
				Once()

			t := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
			//t := time.Now().Round(time.Second)
			sequencedBatch := etherman.SequencedBatch{
				BatchNumber:   uint64(2),
				Coinbase:      common.HexToAddress("0x222"),
				SequencerAddr: common.HexToAddress("0x00"),
				TxHash:        common.HexToHash("0x333"),
				PolygonRollupBaseEtrogBatchData: &etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
					Transactions:         []byte{},
					ForcedGlobalExitRoot: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					ForcedTimestamp:      uint64(t.Unix()),
					ForcedBlockHashL1:    common.HexToHash("0x444"),
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock1.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				RawTxsData:        sequencedBatch.PolygonRollupBaseEtrogBatchData.Transactions,
				ForcedAt:          time.Unix(int64(sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0),
			}}

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  t,
				BlockHash:   ethBlock0.Hash(),
			}
			ethermanBlock1 := etherman.Block{
				BlockNumber:      1,
				ReceivedAt:       t,
				BlockHash:        ethBlock1.Hash(),
				SequencedBatches: [][]etherman.SequencedBatch{{sequencedBatch}},
				ForcedBatches:    forceb,
			}
			blocks := []etherman.Block{ethermanBlock0, ethermanBlock1}
			order := map[common.Hash][]etherman.Order{
				ethBlock1.Hash(): {
					{
						Name: etherman.ForcedBatchesOrder,
						Pos:  0,
					},
					{
						Name: etherman.SequenceBatchesOrder,
						Pos:  0,
					},
				},
			}

			fromBlock := ethBlock0.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock1.NumberU64() {
				toBlock = ethBlock1.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil)

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock := &state.Block{
				BlockNumber: ethermanBlock1.BlockNumber,
				BlockHash:   ethermanBlock1.BlockHash,
				ParentHash:  ethermanBlock1.ParentHash,
				ReceivedAt:  ethermanBlock1.ReceivedAt,
				Checked:     true,
			}

			executionResponse := executor.ProcessBatchResponseV2{
				NewStateRoot:     common.Hash{}.Bytes(),
				NewAccInputHash:  common.Hash{}.Bytes(),
				NewLocalExitRoot: common.Hash{}.Bytes(),
			}
			m.State.EXPECT().ExecuteBatchV2(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(&executionResponse, nil).
				Times(1)

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock1.NumberU64(), nil).
				Once()

			m.State.
				On("AddBlock", ctx, stateBlock, m.DbTx).
				Return(nil).
				Once()

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock1.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				RawTxsData:        sequencedBatch.PolygonRollupBaseEtrogBatchData.Transactions,
				ForcedAt:          time.Unix(int64(sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0),
			}}

			m.State.
				On("AddForcedBatch", ctx, &fb[0], m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetNextForcedBatches", ctx, 1, m.DbTx).
				Return(fb, nil).
				Once()

			trustedBatch := &state.Batch{
				BatchL2Data:    sequencedBatch.PolygonRollupBaseEtrogBatchData.Transactions,
				GlobalExitRoot: sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				Timestamp:      time.Unix(int64(sequencedBatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0),
				Coinbase:       sequencedBatch.Coinbase,
			}

			m.State.
				On("GetBatchByNumber", ctx, sequencedBatch.BatchNumber, m.DbTx).
				Return(trustedBatch, nil).
				Once()

			var forcedGER common.Hash = sequencedBatch.ForcedGlobalExitRoot
			virtualBatch := &state.VirtualBatch{
				BatchNumber:         sequencedBatch.BatchNumber,
				TxHash:              sequencedBatch.TxHash,
				Coinbase:            sequencedBatch.Coinbase,
				BlockNumber:         ethermanBlock1.BlockNumber,
				TimestampBatchEtrog: &t,
				L1InfoRoot:          &forcedGER,
			}

			m.State.
				On("AddVirtualBatch", ctx, virtualBatch, m.DbTx).
				Return(nil).
				Once()

			seq := state.Sequence{
				FromBatchNumber: 2,
				ToBatchNumber:   2,
			}
			m.State.
				On("AddSequence", ctx, seq, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("AddAccumulatedInputHash", ctx, sequencedBatch.BatchNumber, common.Hash{}, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Run(func(args mock.Arguments) { sync.Stop() }).
				Return(nil).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

// TODO: this is running against old sequential L1 sync, need to update to parallel L1 sync.
// but it used a feature that is not implemented in new one that is asking beyond the last block on L1
func TestSequenceForcedBatchIncaberry(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
	}

	m := mocks{
		Etherman:    mock_syncinterfaces.NewEthermanFullInterface(t),
		State:       mock_syncinterfaces.NewStateFullInterface(t),
		Pool:        mock_syncinterfaces.NewPoolInterface(t),
		DbTx:        syncMocks.NewDbTxMock(t),
		ZKEVMClient: mock_syncinterfaces.NewZKEVMClientInterface(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(true, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}
			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(1), nil).
				Maybe()

			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock0, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil).
				Once()

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())
			m.Etherman.
				On("HeaderByNumber", ctx, n).
				Return(ethHeader1, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			sequencedForceBatch := etherman.SequencedForceBatch{
				BatchNumber: uint64(2),
				Coinbase:    common.HexToAddress("0x222"),
				TxHash:      common.HexToHash("0x333"),
				PolygonRollupBaseEtrogBatchData: etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
					Transactions:         []byte{},
					ForcedGlobalExitRoot: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					ForcedTimestamp:      1000, //ForcedBatch
					ForcedBlockHashL1:    common.HexToHash("0x444"),
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock1.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0),
			}}

			ethermanBlock0 := etherman.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			ethermanBlock1 := etherman.Block{
				BlockNumber:           ethBlock1.NumberU64(),
				BlockHash:             ethBlock1.Hash(),
				ParentHash:            ethBlock1.ParentHash(),
				SequencedForceBatches: [][]etherman.SequencedForceBatch{{sequencedForceBatch}},
				ForcedBatches:         forceb,
			}
			blocks := []etherman.Block{ethermanBlock0, ethermanBlock1}
			order := map[common.Hash][]etherman.Order{
				ethBlock1.Hash(): {
					{
						Name: etherman.ForcedBatchesOrder,
						Pos:  0,
					},
					{
						Name: etherman.SequenceForceBatchesOrder,
						Pos:  0,
					},
				},
			}

			fromBlock := ethBlock0.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock1.NumberU64() {
				toBlock = ethBlock1.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", ctx, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock1.NumberU64(), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock := &state.Block{
				BlockNumber: ethermanBlock1.BlockNumber,
				BlockHash:   ethermanBlock1.BlockHash,
				ParentHash:  ethermanBlock1.ParentHash,
				ReceivedAt:  ethermanBlock1.ReceivedAt,
				Checked:     true,
			}

			m.State.
				On("AddBlock", ctx, stateBlock, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetForkIDByBlockNumber", stateBlock.BlockNumber).
				Return(uint64(9), nil).
				Once()

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock1.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0),
			}}

			m.State.
				On("AddForcedBatch", ctx, &fb[0], m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetLastVirtualBatchNum", ctx, m.DbTx).
				Return(uint64(1), nil).
				Once()

			m.State.
				On("ResetTrustedState", ctx, uint64(1), m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetNextForcedBatches", ctx, 1, m.DbTx).
				Return(fb, nil).
				Once()

			f := uint64(1)
			processingContext := state.ProcessingContext{
				BatchNumber:    sequencedForceBatch.BatchNumber,
				Coinbase:       sequencedForceBatch.Coinbase,
				Timestamp:      ethBlock1.ReceivedAt,
				GlobalExitRoot: sequencedForceBatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				ForcedBatchNum: &f,
				BatchL2Data:    &sequencedForceBatch.PolygonRollupBaseEtrogBatchData.Transactions,
			}

			m.State.
				On("ProcessAndStoreClosedBatch", ctx, processingContext, sequencedForceBatch.PolygonRollupBaseEtrogBatchData.Transactions, m.DbTx, metrics.SynchronizerCallerLabel).
				Return(common.Hash{}, uint64(1), cProverIDExecution, nil).
				Once()
			virtualBatch := &state.VirtualBatch{
				BatchNumber:   sequencedForceBatch.BatchNumber,
				TxHash:        sequencedForceBatch.TxHash,
				Coinbase:      sequencedForceBatch.Coinbase,
				SequencerAddr: sequencedForceBatch.Coinbase,
				BlockNumber:   ethermanBlock1.BlockNumber,
			}

			m.State.
				On("AddVirtualBatch", ctx, virtualBatch, m.DbTx).
				Return(nil).
				Once()

			seq := state.Sequence{
				FromBatchNumber: 2,
				ToBatchNumber:   2,
			}
			m.State.
				On("AddSequence", ctx, seq, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Return(nil).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

func setupGenericTest(t *testing.T) (*state.Genesis, *Config, *mocks) {
	genesis := state.Genesis{
		BlockNumber: uint64(123456),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1ParallelSynchronization: L1ParallelSynchronizationConfig{
			MaxClients:                             2,
			MaxPendingNoProcessedBlocks:            2,
			RequestLastBlockPeriod:                 cfgTypes.Duration{Duration: 1 * time.Second},
			RequestLastBlockTimeout:                cfgTypes.Duration{Duration: 1 * time.Second},
			RequestLastBlockMaxRetries:             1,
			StatisticsPeriod:                       cfgTypes.Duration{Duration: 1 * time.Second},
			TimeOutMainLoop:                        cfgTypes.Duration{Duration: 1 * time.Second},
			RollupInfoRetriesSpacing:               cfgTypes.Duration{Duration: 1 * time.Second},
			FallbackToSequentialModeOnSynchronized: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:                      mock_syncinterfaces.NewEthermanFullInterface(t),
		State:                         mock_syncinterfaces.NewStateFullInterface(t),
		Pool:                          mock_syncinterfaces.NewPoolInterface(t),
		DbTx:                          syncMocks.NewDbTxMock(t),
		ZKEVMClient:                   mock_syncinterfaces.NewZKEVMClientInterface(t),
		zkEVMClientEthereumCompatible: mock_syncinterfaces.NewZKEVMClientEthereumCompatibleInterface(t),
		EthTxManager:                  mock_syncinterfaces.NewEthTxManager(t),
		//EventLog:    newEventLogMock(t),
	}
	return &genesis, &cfg, &m
}

func (m mocks) checkExpectedCalls(t *testing.T) {
	m.Etherman.AssertExpectations(t)
	m.State.AssertExpectations(t)
	m.Pool.AssertExpectations(t)
	m.DbTx.AssertExpectations(t)
	m.ZKEVMClient.AssertExpectations(t)
	m.EthTxManager.AssertExpectations(t)
	//m.EventLog.AssertExpectations(t)
}

func transactionToTxData(t types.Transaction) *ethTypes.Transaction {
	inner := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(t.Nonce),
		GasPrice: (*big.Int)(&t.GasPrice),
		Gas:      uint64(t.Gas),
		To:       t.To,
		Value:    (*big.Int)(&t.Value),
		V:        (*big.Int)(&t.V),
		R:        (*big.Int)(&t.R),
		S:        (*big.Int)(&t.S),
	})
	return inner
}

func createTransaction(txIndex uint64) types.Transaction {
	r, _ := new(big.Int).SetString("0x07445CC110033D6A44AD1736ECDF76D26CAB8AB20B9DABB1022EA9BF0707A14E", 0)
	s, _ := new(big.Int).SetString("0x675F2042E60C2D09A4B9C4862693596A75DDE508FCD8E6C7A95283FAD94372EC", 0)
	to := common.HexToAddress("530C75b2E17ac4d1DF146845cF905AEfB31c3607")
	block_hash := common.Hash([common.HashLength]byte{102, 231, 81, 89, 126, 43, 201, 5, 72, 85, 63, 88, 132, 194, 77, 155, 206, 246, 224, 205, 132, 229, 190, 32, 116, 150, 59, 88, 201, 248, 128, 99})
	block_number := types.ArgUint64(1)
	tx_index := types.ArgUint64(txIndex)
	transaction := types.Transaction{
		Nonce:       types.ArgUint64(8),
		GasPrice:    types.ArgBig(*big.NewInt(1000000000)),
		Gas:         types.ArgUint64(21000),
		To:          &to,
		Value:       types.ArgBig(*big.NewInt(2000000000000000000)),
		V:           types.ArgBig(*big.NewInt(2037)),
		R:           types.ArgBig(*r),
		S:           types.ArgBig(*s),
		Hash:        common.Hash([common.HashLength]byte{30, 184, 220, 207, 103, 194, 81, 217, 185, 173, 187, 253, 136, 201, 218, 21, 192, 0, 116, 182, 60, 68, 209, 250, 178, 183, 117, 113, 44, 41, 249, 43}),
		From:        common.HexToAddress("f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		BlockHash:   &block_hash,
		BlockNumber: &block_number,
		TxIndex:     &tx_index,
		ChainID:     types.ArgBig(*big.NewInt(1001)),
		Type:        types.ArgUint64(0),
	}
	return transaction
}

func createBatchL2DataIncaberry(howManyTx int) ([]byte, []types.TransactionOrHash, error) {
	transactions := []types.TransactionOrHash{}
	transactions_state := []ethTypes.Transaction{}
	for i := 0; i < howManyTx; i++ {
		t := createTransaction(uint64(i + 1))
		transaction := types.TransactionOrHash{Tx: &t}
		transactions = append(transactions, transaction)
		transactions_state = append(transactions_state, *transactionToTxData(t))
	}
	encoded, err := state.EncodeTransactions(transactions_state, nil, 4)
	return encoded, transactions, err
}

func createBatchL2DataEtrog(howManyBlocks int, howManyTx int) ([]byte, []types.TransactionOrHash, error) {
	batchV2 := state.BatchRawV2{Blocks: []state.L2BlockRaw{}}
	transactions := []types.TransactionOrHash{}
	for nBlock := 0; nBlock < howManyBlocks; nBlock++ {
		block := state.L2BlockRaw{
			ChangeL2BlockHeader: state.ChangeL2BlockHeader{
				DeltaTimestamp:  123,
				IndexL1InfoTree: 456,
			},
			Transactions: []state.L2TxRaw{},
		}
		for i := 0; i < howManyTx; i++ {
			tx := createTransaction(uint64(i + 1))
			transactions = append(transactions, types.TransactionOrHash{Tx: &tx})
			l2Tx := state.L2TxRaw{
				Tx: *transactionToTxData(tx),
			}

			block.Transactions = append(block.Transactions, l2Tx)
		}
		batchV2.Blocks = append(batchV2.Blocks, block)
	}
	encoded, err := state.EncodeBatchV2(&batchV2)
	return encoded, transactions, err
}

func createBatch(t *testing.T, batchNumber uint64, howManyTx int, etrogMode bool) *types.Batch {
	var err error
	var batchL2Data []byte
	var transactions []types.TransactionOrHash
	if etrogMode {
		batchL2Data, transactions, err = createBatchL2DataEtrog(howManyTx, 1)
		require.NoError(t, err)
	} else {
		batchL2Data, transactions, err = createBatchL2DataIncaberry(howManyTx)
		require.NoError(t, err)
	}
	batch := &types.Batch{
		Number:       types.ArgUint64(batchNumber),
		Coinbase:     common.Address([common.AddressLength]byte{243, 159, 214, 229, 26, 173, 136, 246, 244, 206, 106, 184, 130, 114, 121, 207, 255, 185, 34, 102}),
		Timestamp:    types.ArgUint64(1687854474), // Creation timestamp
		Transactions: transactions,
		BatchL2Data:  batchL2Data,
		StateRoot:    common.HexToHash("0x444"),
	}
	return batch
}

func rpcBatchTostateBatch(rpcBatch *types.Batch) *state.Batch {
	if rpcBatch == nil {
		return nil
	}
	return &state.Batch{
		BatchNumber:    uint64(rpcBatch.Number),
		Coinbase:       rpcBatch.Coinbase,
		StateRoot:      rpcBatch.StateRoot,
		BatchL2Data:    rpcBatch.BatchL2Data,
		GlobalExitRoot: rpcBatch.GlobalExitRoot,
		LocalExitRoot:  rpcBatch.MainnetExitRoot,
		Timestamp:      time.Unix(int64(rpcBatch.Timestamp), 0),
		WIP:            true,
	}
}

func expectedCallsForOpenBatch(t *testing.T, m *mocks, sync *ClientSynchronizer, batchNumber uint64) {
	m.State.
		On("OpenBatch", sync.ctx, mock.Anything, m.DbTx).
		Return(nil).
		Once()
}

func expectedCallsForsyncTrustedState(t *testing.T, m *mocks, sync *ClientSynchronizer,
	batchInPermissionLess *types.Batch, batchInTrustedNode *types.Batch, previousBatchInPermissionless *types.Batch,
	needToRetrieveBatchFromDatabase bool, etrogMode bool) {
	m.State.EXPECT().GetForkIDByBatchNumber(mock.Anything).Return(uint64(7)).Times(1)
	batchNumber := uint64(batchInTrustedNode.Number)
	m.ZKEVMClient.
		On("BatchNumber", mock.Anything).
		Return(batchNumber, nil).
		Once()

	m.ZKEVMClient.
		On("BatchByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).
		Run(func(args mock.Arguments) {
			param := args.Get(1).(*big.Int)
			expected := big.NewInt(int64(batchNumber))
			assert.Equal(t, *expected, *param)
		}).
		Return(batchInTrustedNode, nil).
		Once()

	m.State.
		On("BeginStateTransaction", mock.Anything).
		Return(m.DbTx, nil).
		Once()

	stateBatchInTrustedNode := rpcBatchTostateBatch(batchInTrustedNode)
	stateBatchInPermissionLess := rpcBatchTostateBatch(batchInPermissionLess)
	statePreviousBatchInPermissionless := rpcBatchTostateBatch(previousBatchInPermissionless)

	if needToRetrieveBatchFromDatabase {
		if statePreviousBatchInPermissionless != nil {
			m.State.
				On("GetBatchByNumber", mock.Anything, uint64(batchInTrustedNode.Number-1), mock.Anything).
				Return(statePreviousBatchInPermissionless, nil).
				Once()
		} else {
			m.State.
				On("GetBatchByNumber", mock.Anything, uint64(batchInTrustedNode.Number-1), mock.Anything).
				Return(nil, state.ErrNotFound).
				Once()
		}
		if stateBatchInPermissionLess != nil {
			m.State.
				On("GetBatchByNumber", mock.Anything, uint64(batchInTrustedNode.Number), mock.Anything).
				Return(stateBatchInPermissionLess, nil).
				Once()
		} else {
			m.State.
				On("GetBatchByNumber", mock.Anything, uint64(batchInTrustedNode.Number), mock.Anything).
				Return(nil, state.ErrNotFound).
				Once()
		}
	}
	tx1 := state.ProcessTransactionResponse{}
	block1 := state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{&tx1},
	}
	processedBatch := state.ProcessBatchResponse{
		FlushID:        1,
		ProverID:       cProverIDExecution,
		BlockResponses: []*state.ProcessBlockResponse{&block1},
		//NewStateRoot:   common.HexToHash("0x444"),
		NewStateRoot: batchInTrustedNode.StateRoot,
	}
	if etrogMode {
		m.State.EXPECT().GetL1InfoTreeDataFromBatchL2Data(mock.Anything, mock.Anything, mock.Anything).Return(map[uint32]state.L1DataV2{}, common.Hash{}, common.Hash{}, nil).Times(1)
		m.State.EXPECT().ProcessBatchV2(mock.Anything, mock.Anything, mock.Anything).
			Return(&processedBatch, "", nil).Times(1)
		m.State.EXPECT().StoreL2Block(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(common.Hash{}, nil).Times(1)
		m.State.EXPECT().UpdateWIPBatch(mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Times(1)
		m.State.EXPECT().GetBatchByNumber(mock.Anything, mock.Anything, mock.Anything).
			Return(stateBatchInTrustedNode, nil).Maybe()
	} else {
		m.State.
			On("ProcessBatch", mock.Anything, mock.Anything, true).
			Return(&processedBatch, nil).
			Once()
		m.State.
			On("StoreTransaction", sync.ctx, stateBatchInTrustedNode.BatchNumber, mock.Anything, stateBatchInTrustedNode.Coinbase, uint64(batchInTrustedNode.Timestamp), common.Hash{}, common.Hash{}, mock.Anything, m.DbTx).
			Return(&state.L2Header{}, nil).
			Once()
	}

	m.State.
		On("GetStoredFlushID", mock.Anything).
		Return(uint64(1), cProverIDExecution, nil).
		Once()

	m.DbTx.
		On("Commit", mock.Anything).
		Return(nil).
		Once()
}

func TestReorg(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         3,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		ZKEVMClient:  mock_syncinterfaces.NewZKEVMClientInterface(t),
		EthTxManager: mock_syncinterfaces.NewEthTxManager(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		ForkId:          9,
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(9)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1bis := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash(), Time: 10, GasUsed: 20, Root: common.HexToHash("0x234")}
			ethBlock1bis := ethTypes.NewBlockWithHeader(ethHeader1bis)
			ethHeader2bis := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1bis.Hash()}
			ethBlock2bis := ethTypes.NewBlockWithHeader(ethHeader2bis)
			ethHeader3bis := &ethTypes.Header{Number: big.NewInt(3), ParentHash: ethBlock2bis.Hash()}
			ethBlock3bis := ethTypes.NewBlockWithHeader(ethHeader3bis)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			ethHeader2 := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1.Hash()}
			ethBlock2 := ethTypes.NewBlockWithHeader(ethHeader2)
			ethHeader3 := &ethTypes.Header{Number: big.NewInt(3), ParentHash: ethBlock2.Hash()}
			ethBlock3 := ethTypes.NewBlockWithHeader(ethHeader3)

			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(9), nil).
				Maybe()
			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock1, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())
			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3bis, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			ti := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)

			ethermanBlock1bis := etherman.Block{
				BlockNumber: 1,
				ReceivedAt:  ti,
				BlockHash:   ethBlock1bis.Hash(),
				ParentHash:  ethBlock1bis.ParentHash(),
			}
			ethermanBlock2bis := etherman.Block{
				BlockNumber: 2,
				ReceivedAt:  ti,
				BlockHash:   ethBlock2bis.Hash(),
				ParentHash:  ethBlock2bis.ParentHash(),
			}
			blocks := []etherman.Block{ethermanBlock1bis, ethermanBlock2bis}
			order := map[common.Hash][]etherman.Order{}

			fromBlock := ethBlock1.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock3.NumberU64() {
				toBlock = ethBlock3.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			var depth uint64 = 1
			stateBlock0 := &state.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, m.DbTx).
				Return(stateBlock0, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			m.State.
				On("Reset", ctx, ethBlock0.NumberU64(), m.DbTx).
				Return(nil).
				Once()

			m.EthTxManager.
				On("Reorg", ctx, ethBlock0.NumberU64()+1, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3bis, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  ti,
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			ethermanBlock3bis := etherman.Block{
				BlockNumber: 3,
				ReceivedAt:  ti,
				BlockHash:   ethBlock3bis.Hash(),
				ParentHash:  ethBlock3bis.ParentHash(),
			}
			fromBlock = 0
			blocks2 := []etherman.Block{ethermanBlock0, ethermanBlock1bis, ethermanBlock2bis, ethermanBlock3bis}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks2, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock2bis.NumberU64(), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock1bis := &state.Block{
				BlockNumber: ethermanBlock1bis.BlockNumber,
				BlockHash:   ethermanBlock1bis.BlockHash,
				ParentHash:  ethermanBlock1bis.ParentHash,
				ReceivedAt:  ethermanBlock1bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock1bis, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock2bis := &state.Block{
				BlockNumber: ethermanBlock2bis.BlockNumber,
				BlockHash:   ethermanBlock2bis.BlockHash,
				ParentHash:  ethermanBlock2bis.ParentHash,
				ReceivedAt:  ethermanBlock2bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock2bis, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock3bis := &state.Block{
				BlockNumber: ethermanBlock3bis.BlockNumber,
				BlockHash:   ethermanBlock3bis.BlockHash,
				ParentHash:  ethermanBlock3bis.ParentHash,
				ReceivedAt:  ethermanBlock3bis.ReceivedAt,
				Checked:     false,
			}
			m.State.
				On("AddBlock", ctx, stateBlock3bis, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

func TestLatestSyncedBlockEmpty(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         3,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		ZKEVMClient:  mock_syncinterfaces.NewZKEVMClientInterface(t),
		EthTxManager: mock_syncinterfaces.NewEthTxManager(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		ForkId:          9,
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(9)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			ethHeader2 := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1.Hash()}
			ethBlock2 := ethTypes.NewBlockWithHeader(ethHeader2)
			ethHeader3 := &ethTypes.Header{Number: big.NewInt(3), ParentHash: ethBlock2.Hash()}
			ethBlock3 := ethTypes.NewBlockWithHeader(ethHeader3)

			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(9), nil).
				Maybe()
			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock1, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())
			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			blocks := []etherman.Block{}
			order := map[common.Hash][]etherman.Order{}

			fromBlock := ethBlock1.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock3.NumberU64() {
				toBlock = ethBlock3.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			ti := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
			var depth uint64 = 1
			stateBlock0 := &state.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, nil).
				Return(stateBlock0, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			m.State.
				On("Reset", ctx, ethBlock0.NumberU64(), m.DbTx).
				Return(nil).
				Once()

			m.EthTxManager.
				On("Reorg", ctx, ethBlock0.NumberU64()+1, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  ti,
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			blocks = []etherman.Block{ethermanBlock0}
			fromBlock = 0
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock3.NumberU64(), nil).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

func TestRegularReorg(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         3,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		ZKEVMClient:  mock_syncinterfaces.NewZKEVMClientInterface(t),
		EthTxManager: mock_syncinterfaces.NewEthTxManager(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		ForkId:          9,
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(9)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1bis := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash(), Time: 10, GasUsed: 20, Root: common.HexToHash("0x234")}
			ethBlock1bis := ethTypes.NewBlockWithHeader(ethHeader1bis)
			ethHeader2bis := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1bis.Hash()}
			ethBlock2bis := ethTypes.NewBlockWithHeader(ethHeader2bis)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			ethHeader2 := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1.Hash()}
			ethBlock2 := ethTypes.NewBlockWithHeader(ethHeader2)

			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(9), nil).
				Maybe()
			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock1, nil).
				Once()

			// After a ResetState get lastblock that must be block 0
			m.State.
				On("GetLastBlock", ctx, nil).
				Return(lastBlock0, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1bis, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			ti := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
			var depth uint64 = 1
			stateBlock0 := &state.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, m.DbTx).
				Return(stateBlock0, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			m.State.
				On("Reset", ctx, ethBlock0.NumberU64(), m.DbTx).
				Return(nil).
				Once()

			m.EthTxManager.
				On("Reorg", ctx, ethBlock0.NumberU64()+1, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader2bis, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  ti,
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			ethermanBlock1bis := etherman.Block{
				BlockNumber: 1,
				ReceivedAt:  ti,
				BlockHash:   ethBlock1bis.Hash(),
				ParentHash:  ethBlock1bis.ParentHash(),
			}
			ethermanBlock2bis := etherman.Block{
				BlockNumber: 2,
				ReceivedAt:  ti,
				BlockHash:   ethBlock2bis.Hash(),
				ParentHash:  ethBlock2bis.ParentHash(),
			}
			blocks := []etherman.Block{ethermanBlock0, ethermanBlock1bis, ethermanBlock2bis}
			order := map[common.Hash][]etherman.Order{}

			fromBlock := ethBlock0.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock2.NumberU64() {
				toBlock = ethBlock2.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock2bis.NumberU64(), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock1bis := &state.Block{
				BlockNumber: ethermanBlock1bis.BlockNumber,
				BlockHash:   ethermanBlock1bis.BlockHash,
				ParentHash:  ethermanBlock1bis.ParentHash,
				ReceivedAt:  ethermanBlock1bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock1bis, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock2bis := &state.Block{
				BlockNumber: ethermanBlock2bis.BlockNumber,
				BlockHash:   ethermanBlock2bis.BlockHash,
				ParentHash:  ethermanBlock2bis.ParentHash,
				ReceivedAt:  ethermanBlock2bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock2bis, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Return(nil).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

func TestLatestSyncedBlockEmptyWithExtraReorg(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         3,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		ZKEVMClient:  mock_syncinterfaces.NewZKEVMClientInterface(t),
		EthTxManager: mock_syncinterfaces.NewEthTxManager(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		ForkId:          9,
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(9)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			ethHeader1bis := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash(), Time: 0, GasUsed: 10}
			ethBlock1bis := ethTypes.NewBlockWithHeader(ethHeader1bis)
			ethHeader2 := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1.Hash()}
			ethBlock2 := ethTypes.NewBlockWithHeader(ethHeader2)
			ethHeader3 := &ethTypes.Header{Number: big.NewInt(3), ParentHash: ethBlock2.Hash()}
			ethBlock3 := ethTypes.NewBlockWithHeader(ethHeader3)

			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}
			lastBlock2 := &state.Block{BlockHash: ethBlock2.Hash(), BlockNumber: ethBlock2.Number().Uint64(), ParentHash: ethBlock2.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(9), nil).
				Maybe()
			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock2, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock2.BlockNumber).
				Return(ethBlock2, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())
			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock2.BlockNumber).
				Return(ethBlock2, nil).
				Once()

			blocks := []etherman.Block{}
			order := map[common.Hash][]etherman.Order{}

			fromBlock := ethBlock2.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock3.NumberU64() {
				toBlock = ethBlock3.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			ti := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
			var depth uint64 = 1
			stateBlock1 := &state.Block{
				BlockNumber: ethBlock1.NumberU64(),
				BlockHash:   ethBlock1.Hash(),
				ParentHash:  ethBlock1.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, nil).
				Return(stateBlock1, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1bis, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock0 := &state.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, m.DbTx).
				Return(stateBlock0, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			m.State.
				On("Reset", ctx, ethBlock0.NumberU64(), m.DbTx).
				Return(nil).
				Once()

			m.EthTxManager.
				On("Reorg", ctx, ethBlock0.NumberU64()+1, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader3, nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  ti,
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			ethermanBlock1bis := etherman.Block{
				BlockNumber: 1,
				ReceivedAt:  ti,
				BlockHash:   ethBlock1.Hash(),
				ParentHash:  ethBlock1.ParentHash(),
			}
			blocks = []etherman.Block{ethermanBlock0, ethermanBlock1bis}
			fromBlock = 0
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock3.NumberU64(), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock1bis := &state.Block{
				BlockNumber: ethermanBlock1bis.BlockNumber,
				BlockHash:   ethermanBlock1bis.BlockHash,
				ParentHash:  ethermanBlock1bis.ParentHash,
				ReceivedAt:  ethermanBlock1bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock1bis, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}

func TestCallFromEmptyBlockAndReorg(t *testing.T) {
	genesis := state.Genesis{
		BlockNumber: uint64(0),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         3,
		L1SynchronizationMode: SequentialMode,
		SyncBlockProtection:   "latest",
		L1BlockCheck: L1BlockCheckConfig{
			Enabled: false,
		},
		L2Synchronization: l2_sync.Config{
			Enabled: true,
		},
	}

	m := mocks{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		ZKEVMClient:  mock_syncinterfaces.NewZKEVMClientInterface(t),
		EthTxManager: mock_syncinterfaces.NewEthTxManager(t),
	}
	ethermanForL1 := []syncinterfaces.EthermanFullInterface{m.Etherman}
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, m.zkEVMClientEthereumCompatible, nil, genesis, cfg, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	forkIdInterval := state.ForkIDInterval{
		ForkId:          9,
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
	}
	m.State.EXPECT().GetForkIDInMemory(uint64(9)).Return(&forkIdInterval)

	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader0 := &ethTypes.Header{Number: big.NewInt(0), ParentHash: parentHash}
			ethBlock0 := ethTypes.NewBlockWithHeader(ethHeader0)
			ethHeader1bis := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash(), Time: 10, GasUsed: 20, Root: common.HexToHash("0x234")}
			ethBlock1bis := ethTypes.NewBlockWithHeader(ethHeader1bis)
			ethHeader2bis := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1bis.Hash()}
			ethBlock2bis := ethTypes.NewBlockWithHeader(ethHeader2bis)
			ethHeader1 := &ethTypes.Header{Number: big.NewInt(1), ParentHash: ethBlock0.Hash()}
			ethBlock1 := ethTypes.NewBlockWithHeader(ethHeader1)
			ethHeader2 := &ethTypes.Header{Number: big.NewInt(2), ParentHash: ethBlock1.Hash()}
			ethBlock2 := ethTypes.NewBlockWithHeader(ethHeader2)

			lastBlock0 := &state.Block{BlockHash: ethBlock0.Hash(), BlockNumber: ethBlock0.Number().Uint64(), ParentHash: ethBlock0.ParentHash()}
			lastBlock1 := &state.Block{BlockHash: ethBlock1.Hash(), BlockNumber: ethBlock1.Number().Uint64(), ParentHash: ethBlock1.ParentHash()}

			m.State.
				On("GetForkIDByBatchNumber", mock.Anything).
				Return(uint64(9), nil).
				Maybe()
			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock1, nil).
				Once()

			m.State.
				On("GetLastBatchNumber", ctx, m.DbTx).
				Return(uint64(10), nil).
				Once()

			m.State.
				On("SetInitSyncBatch", ctx, uint64(10), m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("GetLatestBatchNumber").
				Return(uint64(10), nil)

			var nilDbTx pgx.Tx
			m.State.
				On("GetLastBatchNumber", ctx, nilDbTx).
				Return(uint64(10), nil)

			m.Etherman.
				On("GetLatestVerifiedBatchNum").
				Return(uint64(10), nil)

			m.State.
				On("SetLastBatchInfoSeenOnEthereum", ctx, uint64(10), uint64(10), nilDbTx).
				Return(nil)

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			n := big.NewInt(rpc.LatestBlockNumber.Int64())

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock1.BlockNumber).
				Return(ethBlock1, nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader2bis, nil).
				Once()

			ti := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)

			ethermanBlock0 := etherman.Block{
				BlockNumber: 0,
				ReceivedAt:  ti,
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
			}
			ethermanBlock2bis := etherman.Block{
				BlockNumber: 2,
				ReceivedAt:  ti,
				BlockHash:   ethBlock2bis.Hash(),
				ParentHash:  ethBlock2bis.ParentHash(),
			}
			blocks := []etherman.Block{ethermanBlock2bis}
			order := map[common.Hash][]etherman.Order{}

			fromBlock := ethBlock1.NumberU64()
			toBlock := fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock2.NumberU64() {
				toBlock = ethBlock2.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			var depth uint64 = 1
			stateBlock0 := &state.Block{
				BlockNumber: ethBlock0.NumberU64(),
				BlockHash:   ethBlock0.Hash(),
				ParentHash:  ethBlock0.ParentHash(),
				ReceivedAt:  ti,
			}
			m.State.
				On("GetPreviousBlock", ctx, depth, m.DbTx).
				Return(stateBlock0, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			m.State.
				On("Reset", ctx, ethBlock0.NumberU64(), m.DbTx).
				Return(nil).
				Once()

			m.EthTxManager.
				On("Reorg", ctx, ethBlock0.NumberU64()+1, m.DbTx).
				Return(nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Return(nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock0.BlockNumber).
				Return(ethBlock0, nil).
				Once()

			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader2bis, nil).
				Once()

			blocks = []etherman.Block{ethermanBlock0, ethermanBlock2bis}
			fromBlock = ethBlock0.NumberU64()
			toBlock = fromBlock + cfg.SyncChunkSize
			if toBlock > ethBlock2.NumberU64() {
				toBlock = ethBlock2.NumberU64()
			}
			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.Etherman.
				On("GetFinalizedBlockNumber", ctx).
				Return(ethBlock2bis.NumberU64(), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock2bis := &state.Block{
				BlockNumber: ethermanBlock2bis.BlockNumber,
				BlockHash:   ethermanBlock2bis.BlockHash,
				ParentHash:  ethermanBlock2bis.ParentHash,
				ReceivedAt:  ethermanBlock2bis.ReceivedAt,
				Checked:     true,
			}
			m.State.
				On("AddBlock", ctx, stateBlock2bis, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetStoredFlushID", ctx).
				Return(uint64(1), cProverIDExecution, nil).
				Once()

			m.DbTx.
				On("Commit", ctx).
				Run(func(args mock.Arguments) {
					sync.Stop()
					ctx.Done()
				}).
				Return(nil).
				Once()
		}).
		Return(m.DbTx, nil).
		Once()

	err = sync.Sync()
	require.NoError(t, err)
}
