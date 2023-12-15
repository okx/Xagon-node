package synchronizer

import (
	context "context"
	"math/big"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	cProverIDExecution = "PROVER_ID-EXE001"
)

type mocks struct {
	Etherman                   *ethermanMock
	State                      *stateMock
	Pool                       *poolMock
	EthTxManager               *ethTxManagerMock
	DbTx                       *dbTxMock
	ZKEVMClient                *zkEVMClientMock
	DataCommitteeClientFactory *dataCommitteeClientFactoryMock
	//EventLog     *eventLogMock
}

// Feature #2220 and  #2239: Optimize Trusted state synchronization
//
//	this Check partially point 2: Use previous batch stored in memory to avoid getting from database
func TestGivenPermissionlessNodeWhenSyncronizeAgainSameBatchThenUseTheOneInMemoryInstaeadOfGettingFromDb(t *testing.T) {
	genesis, cfg, m := setupGenericTest(t)
	ethermanForL1 := []EthermanInterface{m.Etherman}
	m.Etherman.
		On("GetCurrentDataCommittee").
		Return(&etherman.DataCommittee{}, nil)
	syncInterface, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, nil, *genesis, *cfg, nil, false)
	require.NoError(t, err)
	sync, ok := syncInterface.(*ClientSynchronizer)
	require.EqualValues(t, true, ok, "Can't convert to underlaying struct the interface of syncronizer")
	lastBatchNumber := uint64(10)
	batch10With2Tx := createBatch(t, lastBatchNumber, 2)
	batch10With3Tx := createBatch(t, lastBatchNumber, 3)
	previousBatch09 := createBatch(t, lastBatchNumber-1, 1)

	expectedCallsForsyncTrustedState(t, m, sync, nil, batch10With2Tx, previousBatch09, true, false)
	// Is the first time that appears this batch, so it need to OpenBatch
	expectedCallsForOpenBatch(t, m, sync, lastBatchNumber)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)
	expectedCallsForsyncTrustedState(t, m, sync, batch10With2Tx, batch10With3Tx, previousBatch09, true, true)
	expectedCallsForOpenBatch(t, m, sync, lastBatchNumber)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)
	require.Equal(t, sync.trustedState.lastTrustedBatches[0], rpcBatchTostateBatch(batch10With3Tx))
}

// Feature #2220 and  #2239: Optimize Trusted state synchronization
//
//	this Check partially point 2: Store last batch in memory (CurrentTrustedBatch)
func TestGivenPermissionlessNodeWhenSyncronizeFirstTimeABatchThenStoreItInALocalVar(t *testing.T) {
	genesis, cfg, m := setupGenericTest(t)
	ethermanForL1 := []EthermanInterface{m.Etherman}
	m.Etherman.
		On("GetCurrentDataCommittee").
		Return(&etherman.DataCommittee{}, nil)
	syncInterface, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, nil, *genesis, *cfg, nil, false)
	require.NoError(t, err)
	sync, ok := syncInterface.(*ClientSynchronizer)
	require.EqualValues(t, true, ok, "Can't convert to underlaying struct the interface of syncronizer")
	lastBatchNumber := uint64(10)
	batch10With1Tx := createBatch(t, lastBatchNumber, 1)
	batch10With2Tx := createBatch(t, lastBatchNumber, 2)
	previousBatch09 := createBatch(t, lastBatchNumber-1, 1)

	expectedCallsForsyncTrustedState(t, m, sync, batch10With1Tx, batch10With2Tx, previousBatch09, true, true)
	expectedCallsForOpenBatch(t, m, sync, lastBatchNumber)
	err = sync.syncTrustedState(lastBatchNumber)
	require.NoError(t, err)
	require.Equal(t, sync.trustedState.lastTrustedBatches[0], rpcBatchTostateBatch(batch10With2Tx))
}

// issue #2220
// TODO: this is running against old sequential L1 sync, need to update to parallel L1 sync.
// but it used a feature that is not implemented in new one that is asking beyond the last block on L1
func TestForcedBatch(t *testing.T) {
	genesis := state.Genesis{
		GenesisBlockNum: uint64(123456),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
	}

	m := mocks{
		Etherman:    newEthermanMock(t),
		State:       newStateMock(t),
		Pool:        newPoolMock(t),
		DbTx:        newDbTxMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
	}
	ethermanForL1 := []EthermanInterface{m.Etherman}
	m.Etherman.
		On("GetCurrentDataCommittee").
		Return(&etherman.DataCommittee{}, nil)
	sync, err := NewSynchronizer(false, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, nil, genesis, cfg, nil, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader := &ethTypes.Header{Number: big.NewInt(1), ParentHash: parentHash}
			ethBlock := ethTypes.NewBlockWithHeader(ethHeader)
			lastBlock := &state.Block{BlockHash: ethBlock.Hash(), BlockNumber: ethBlock.Number().Uint64()}

			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock, nil).
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
				Return(uint64(10), nil).
				Once()

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

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock.BlockNumber).
				Return(ethBlock, nil).
				Once()

			var n *big.Int
			m.Etherman.
				On("HeaderByNumber", mock.Anything, n).
				Return(ethHeader, nil).
				Once()

			t := time.Now()
			txs := []byte{}
			sequencedBatch := etherman.SequencedBatch{
				BatchNumber:   uint64(2),
				Coinbase:      common.HexToAddress("0x222"),
				SequencerAddr: common.HexToAddress("0x00"),
				TxHash:        common.HexToHash("0x333"),
				PolygonZkEVMBatchData: polygonzkevm.PolygonZkEVMBatchData{
					TransactionsHash:   crypto.Keccak256Hash(txs),
					GlobalExitRoot:     [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					Timestamp:          uint64(t.Unix()),
					MinForcedTimestamp: 1000, //ForcedBatch
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.GlobalExitRoot,
				RawTxsData:        txs,
				ForcedAt:          time.Unix(int64(sequencedBatch.MinForcedTimestamp), 0),
			}}

			ethermanBlock := etherman.Block{
				BlockHash:        ethBlock.Hash(),
				SequencedBatches: [][]etherman.SequencedBatch{{sequencedBatch}},
				ForcedBatches:    forceb,
			}
			blocks := []etherman.Block{ethermanBlock}
			order := map[common.Hash][]etherman.Order{
				ethBlock.Hash(): {
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

			fromBlock := ethBlock.NumberU64() + 1
			toBlock := fromBlock + cfg.SyncChunkSize

			m.Etherman.
				On("GetRollupInfoByBlockRange", mock.Anything, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock := &state.Block{
				BlockNumber: ethermanBlock.BlockNumber,
				BlockHash:   ethermanBlock.BlockHash,
				ParentHash:  ethermanBlock.ParentHash,
				ReceivedAt:  ethermanBlock.ReceivedAt,
			}

			m.State.
				On("AddBlock", ctx, stateBlock, m.DbTx).
				Return(nil).
				Once()
			m.State.
				On("GetBatchL2DataByNumber", ctx, uint64(2), nil).
				Return(txs, nil).
				Once()

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.GlobalExitRoot,
				RawTxsData:        []byte{},
				ForcedAt:          time.Unix(int64(sequencedBatch.MinForcedTimestamp), 0),
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
				BatchL2Data:    txs,
				GlobalExitRoot: sequencedBatch.GlobalExitRoot,
				Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
				Coinbase:       sequencedBatch.Coinbase,
			}

			m.State.
				On("GetBatchByNumber", ctx, sequencedBatch.BatchNumber, m.DbTx).
				Return(trustedBatch, nil).
				Once()

			var forced uint64 = 1
			sbatch := state.Batch{
				BatchNumber:    sequencedBatch.BatchNumber,
				Coinbase:       common.HexToAddress("0x222"),
				BatchL2Data:    []byte{},
				Timestamp:      time.Unix(int64(t.Unix()), 0),
				Transactions:   nil,
				GlobalExitRoot: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
				ForcedBatchNum: &forced,
			}
			m.State.On("ExecuteBatch", ctx, sbatch, false, m.DbTx).
				Return(&executor.ProcessBatchResponse{NewStateRoot: trustedBatch.StateRoot.Bytes()}, nil).
				Once()

			virtualBatch := &state.VirtualBatch{
				BatchNumber: sequencedBatch.BatchNumber,
				TxHash:      sequencedBatch.TxHash,
				Coinbase:    sequencedBatch.Coinbase,
				BlockNumber: ethermanBlock.BlockNumber,
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
func TestSequenceForcedBatch(t *testing.T) {
	genesis := state.Genesis{
		GenesisBlockNum: uint64(123456),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
	}

	m := mocks{
		Etherman:    newEthermanMock(t),
		State:       newStateMock(t),
		Pool:        newPoolMock(t),
		DbTx:        newDbTxMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
	}
	ethermanForL1 := []EthermanInterface{m.Etherman}
	m.Etherman.
		On("GetCurrentDataCommittee").
		Return(nil, nil)
	sync, err := NewSynchronizer(true, m.Etherman, ethermanForL1, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, nil, genesis, cfg, nil, false)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader := &ethTypes.Header{Number: big.NewInt(1), ParentHash: parentHash}
			ethBlock := ethTypes.NewBlockWithHeader(ethHeader)
			lastBlock := &state.Block{BlockHash: ethBlock.Hash(), BlockNumber: ethBlock.Number().Uint64()}

			m.State.
				On("GetLastBlock", ctx, m.DbTx).
				Return(lastBlock, nil).
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
				Return(uint64(10), nil).
				Once()

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

			m.Etherman.
				On("EthBlockByNumber", ctx, lastBlock.BlockNumber).
				Return(ethBlock, nil).
				Once()

			var n *big.Int
			m.Etherman.
				On("HeaderByNumber", ctx, n).
				Return(ethHeader, nil).
				Once()

			sequencedForceBatch := etherman.SequencedForceBatch{
				BatchNumber: uint64(2),
				Coinbase:    common.HexToAddress("0x222"),
				TxHash:      common.HexToHash("0x333"),
				PolygonZkEVMForcedBatchData: polygonzkevm.PolygonZkEVMForcedBatchData{
					Transactions:       []byte{},
					GlobalExitRoot:     [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					MinForcedTimestamp: 1000, //ForcedBatch
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.GlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.MinForcedTimestamp), 0),
			}}

			ethermanBlock := etherman.Block{
				BlockHash:             ethBlock.Hash(),
				SequencedForceBatches: [][]etherman.SequencedForceBatch{{sequencedForceBatch}},
				ForcedBatches:         forceb,
			}
			blocks := []etherman.Block{ethermanBlock}
			order := map[common.Hash][]etherman.Order{
				ethBlock.Hash(): {
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

			fromBlock := ethBlock.NumberU64() + 1
			toBlock := fromBlock + cfg.SyncChunkSize

			m.Etherman.
				On("GetRollupInfoByBlockRange", ctx, fromBlock, &toBlock).
				Return(blocks, order, nil).
				Once()

			m.State.
				On("BeginStateTransaction", ctx).
				Return(m.DbTx, nil).
				Once()

			stateBlock := &state.Block{
				BlockNumber: ethermanBlock.BlockNumber,
				BlockHash:   ethermanBlock.BlockHash,
				ParentHash:  ethermanBlock.ParentHash,
				ReceivedAt:  ethermanBlock.ReceivedAt,
			}

			m.State.
				On("AddBlock", ctx, stateBlock, m.DbTx).
				Return(nil).
				Once()

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.GlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.MinForcedTimestamp), 0),
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
				Timestamp:      ethBlock.ReceivedAt,
				GlobalExitRoot: sequencedForceBatch.GlobalExitRoot,
				ForcedBatchNum: &f,
				BatchL2Data:    &sequencedForceBatch.Transactions,
			}

			m.State.
				On("ProcessAndStoreClosedBatch", ctx, processingContext, sequencedForceBatch.Transactions, m.DbTx, metrics.SynchronizerCallerLabel).
				Return(common.Hash{}, uint64(1), cProverIDExecution, nil).
				Once()

			virtualBatch := &state.VirtualBatch{
				BatchNumber:   sequencedForceBatch.BatchNumber,
				TxHash:        sequencedForceBatch.TxHash,
				Coinbase:      sequencedForceBatch.Coinbase,
				SequencerAddr: sequencedForceBatch.Coinbase,
				BlockNumber:   ethermanBlock.BlockNumber,
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
				Run(func(args mock.Arguments) { sync.Stop() }).
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
		GenesisBlockNum: uint64(123456),
	}
	cfg := Config{
		SyncInterval:          cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:         10,
		L1SynchronizationMode: SequentialMode,
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
	}

	m := mocks{
		Etherman:    newEthermanMock(t),
		State:       newStateMock(t),
		Pool:        newPoolMock(t),
		DbTx:        newDbTxMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
		//EventLog:    newEventLogMock(t),
	}
	return &genesis, &cfg, &m
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

func createBatch(t *testing.T, batchNumber uint64, howManyTx int) *types.Batch {
	transactions := []types.TransactionOrHash{}
	transactions_state := []ethTypes.Transaction{}
	for i := 0; i < howManyTx; i++ {
		t := createTransaction(uint64(i + 1))
		transaction := types.TransactionOrHash{Tx: &t}
		transactions = append(transactions, transaction)
		transactions_state = append(transactions_state, *transactionToTxData(t))
	}
	batchL2Data, err := state.EncodeTransactions(transactions_state, nil, 4)
	require.NoError(t, err)

	batch := &types.Batch{
		Number:       types.ArgUint64(batchNumber),
		Coinbase:     common.Address([common.AddressLength]byte{243, 159, 214, 229, 26, 173, 136, 246, 244, 206, 106, 184, 130, 114, 121, 207, 255, 185, 34, 102}),
		Timestamp:    types.ArgUint64(1687854474), // Creation timestamp
		Transactions: transactions,
		BatchL2Data:  batchL2Data,
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
	needToRetrieveBatchFromDatabase bool, needUpdateL2Data bool) {
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
	if needUpdateL2Data {
		m.State.
			On("ResetTrustedState", sync.ctx, batchNumber-1, mock.Anything).
			Return(nil).
			Once()

		m.State.
			On("UpdateBatchL2Data", mock.Anything, batchNumber, stateBatchInTrustedNode.BatchL2Data, mock.Anything).
			Return(nil).
			Once()
	}
	tx1 := state.ProcessTransactionResponse{}
	processedBatch := state.ProcessBatchResponse{
		FlushID:   1,
		ProverID:  cProverIDExecution,
		Responses: []*state.ProcessTransactionResponse{&tx1},
	}
	m.State.
		On("ProcessBatch", mock.Anything, mock.Anything, true).
		Return(&processedBatch, nil).
		Once()

	m.State.
		On("StoreTransaction", sync.ctx, uint64(stateBatchInTrustedNode.BatchNumber), mock.Anything, stateBatchInTrustedNode.Coinbase, uint64(batchInTrustedNode.Timestamp), mock.Anything, m.DbTx).
		Return(&ethTypes.Header{}, nil).
		Once()

	m.State.
		On("GetStoredFlushID", mock.Anything).
		Return(uint64(1), cProverIDExecution, nil).
		Once()

	m.DbTx.
		On("Commit", mock.Anything).
		Return(nil).
		Once()
}

func TestSsZeroByteArray(t *testing.T) {
	var hash1 [32]byte
	assert.Equal(t, isZeroByteArray(hash1), true)

	hash2 := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	assert.Equal(t, isZeroByteArray(hash2), false)
}
