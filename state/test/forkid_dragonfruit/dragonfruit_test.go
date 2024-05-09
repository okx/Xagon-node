package dragonfruit_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	test "github.com/0xPolygonHermez/zkevm-node/state/test/forkid_common"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Counter"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testState *state.State
	forkID    = uint64(state.FORKID_DRAGONFRUIT)
	stateCfg  = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
)

func TestMain(m *testing.M) {
	testState = test.InitTestState(stateCfg)
	defer test.CloseTestState()
	result := m.Run()
	os.Exit(result)
}

func TestExecutorUnsignedTransactions(t *testing.T) {
	ctx := context.Background()
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	var chainIDSequencer = new(big.Int).SetInt64(1000)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var gasLimit = uint64(4000000)
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	scByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)

	// auth
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	// signed tx to deploy SC
	unsignedTxDeploy := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     common.Hex2Bytes(scByteCode),
	})
	signedTxDeploy, err := auth.Signer(auth.From, unsignedTxDeploy)
	require.NoError(t, err)

	incrementFnSignature := crypto.Keccak256Hash([]byte("increment()")).Bytes()[:4]
	retrieveFnSignature := crypto.Keccak256Hash([]byte("getCount()")).Bytes()[:4]

	// signed tx to call SC
	unsignedTxFirstIncrement := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		To:       &scAddress,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     incrementFnSignature,
	})
	signedTxFirstIncrement, err := auth.Signer(auth.From, unsignedTxFirstIncrement)
	require.NoError(t, err)

	unsignedTxFirstRetrieve := types.NewTx(&types.LegacyTx{
		Nonce:    2,
		To:       &scAddress,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     retrieveFnSignature,
	})
	signedTxFirstRetrieve, err := auth.Signer(auth.From, unsignedTxFirstRetrieve)
	require.NoError(t, err)

	dbTx, err := testState.BeginStateTransaction(context.Background())
	require.NoError(t, err)
	// Set genesis
	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: sequencerAddress.Hex(),
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
	}
	_, err = testState.SetGenesis(ctx, state.Block{}, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	batchCtx := state.ProcessingContext{
		BatchNumber: 1,
		Coinbase:    sequencerAddress,
		Timestamp:   time.Now(),
	}
	err = testState.OpenBatch(context.Background(), batchCtx, dbTx)
	require.NoError(t, err)
	signedTxs := []types.Transaction{
		*signedTxDeploy,
		*signedTxFirstIncrement,
		*signedTxFirstRetrieve,
	}
	threeEffectivePercentages := []uint8{state.MaxEffectivePercentage, state.MaxEffectivePercentage, state.MaxEffectivePercentage}
	batchL2Data, err := state.EncodeTransactions(signedTxs, threeEffectivePercentages, forkID)
	require.NoError(t, err)

	processBatchResponse, err := testState.ProcessSequencerBatch(context.Background(), 1, batchL2Data, metrics.SequencerCallerLabel, dbTx)
	require.NoError(t, err)
	// assert signed tx do deploy sc
	assert.Nil(t, processBatchResponse.BlockResponses[0].TransactionResponses[0].RomError)
	assert.Equal(t, scAddress, processBatchResponse.BlockResponses[0].TransactionResponses[0].CreateAddress)

	// assert signed tx to increment counter
	assert.Nil(t, processBatchResponse.BlockResponses[1].TransactionResponses[0].RomError)

	// assert signed tx to increment counter
	assert.Nil(t, processBatchResponse.BlockResponses[2].TransactionResponses[0].RomError)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", hex.EncodeToString(processBatchResponse.BlockResponses[2].TransactionResponses[0].ReturnValue))

	// Add txs to DB
	err = testState.StoreTransactions(context.Background(), 1, processBatchResponse.BlockResponses, nil, dbTx)
	require.NoError(t, err)
	// Close batch
	err = testState.CloseBatch(
		context.Background(),
		state.ProcessingReceipt{
			BatchNumber:   1,
			StateRoot:     processBatchResponse.NewStateRoot,
			LocalExitRoot: processBatchResponse.NewLocalExitRoot,
		}, dbTx,
	)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(context.Background()))

	unsignedTxSecondRetrieve := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       &scAddress,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     retrieveFnSignature,
	})
	l2BlockNumber := uint64(3)

	result, err := testState.ProcessUnsignedTransaction(context.Background(), unsignedTxSecondRetrieve, common.HexToAddress("0x1000000000000000000000000000000000000000"), &l2BlockNumber, true, nil)
	require.NoError(t, err)
	// assert unsigned tx
	assert.Nil(t, result.Err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", hex.EncodeToString(result.ReturnValue))
}

func TestExecutorEstimateGas(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetUint64(stateCfg.ChainID)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var sequencerBalance = 4000000
	scRevertByteCode, err := testutils.ReadBytecode("Revert2/Revert2.bin")
	require.NoError(t, err)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	nonce := uint64(0)

	// Deploy revert.sol
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scRevertByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Call SC method
	nonce++
	tx1 := types.NewTransaction(nonce, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("4abbb40a"))
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx0, *signedTx1}, constants.TwoEffectivePercentages, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 0,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.NotEqual(t, "", processBatchResponse.Responses[0].Error)

	convertedResponse, err := testState.TestConvertToProcessBatchResponse(processBatchResponse)
	require.NoError(t, err)
	log.Debugf("%v", len(convertedResponse.BlockResponses))

	// Store processed txs into the batch
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	processingContext := state.ProcessingContext{
		BatchNumber:    processBatchRequest.OldBatchNum + 1,
		Coinbase:       common.Address{},
		Timestamp:      time.Now(),
		GlobalExitRoot: common.BytesToHash(processBatchRequest.GlobalExitRoot),
	}

	err = testState.OpenBatch(ctx, processingContext, dbTx)
	require.NoError(t, err)

	err = testState.StoreTransactions(ctx, processBatchRequest.OldBatchNum+1, convertedResponse.BlockResponses, nil, dbTx)
	require.NoError(t, err)

	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   processBatchRequest.OldBatchNum + 1,
		StateRoot:     convertedResponse.NewStateRoot,
		LocalExitRoot: convertedResponse.NewLocalExitRoot,
	}

	err = testState.CloseBatch(ctx, processingReceipt, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// l2BlockNumber := uint64(2)
	nonce++
	tx2 := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scRevertByteCode),
	})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	blockNumber, err := testState.GetLastL2BlockNumber(ctx, nil)
	require.NoError(t, err)

	estimatedGas, _, err := testState.EstimateGas(signedTx2, sequencerAddress, &blockNumber, nil)
	require.NoError(t, err)
	log.Debugf("Estimated gas = %v", estimatedGas)

	nonce++
	tx3 := types.NewTransaction(nonce, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("4abbb40a"))
	signedTx3, err := auth.Signer(auth.From, tx3)
	require.NoError(t, err)
	_, _, err = testState.EstimateGas(signedTx3, sequencerAddress, &blockNumber, nil)
	require.Error(t, err)
}

// TODO: Uncomment once the executor properly returns gas refund
/*
func TestExecutorGasRefund(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(1000)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var sequencerBalance = 4000000
	scStorageByteCode, err := testutils.ReadBytecode("Storage/Storage.bin")
	require.NoError(t, err)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "100000000000000000000000",
			},
			{
				Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "100000000000000000000000",
			},
			{
				Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "100000000000000000000000",
			},
		},
	}

	test.InitOrResetDB(stateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesisAccountsBalance(ctx, block, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// Deploy contract
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scStorageByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Call SC method to set value to 123456
	tx1 := types.NewTransaction(1, scAddress, new(big.Int), 80000, new(big.Int).SetUint64(0), common.Hex2Bytes("6057361d000000000000000000000000000000000000000000000000000000000001e240"))
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx0, *signedTx1})
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		BatchNum:         1,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     stateRoot,
		globalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldLocalExitRoot: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.Equal(t, pb.Error_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)
	assert.Equal(t, pb.Error_ERROR_NO_ERROR, processBatchResponse.Responses[1].Error)

	// Preparation to be able to estimate gas
	convertedResponse, err := state.TestConvertToProcessBatchResponse([]types.Transaction{*signedTx0, *signedTx1}, processBatchResponse)
	require.NoError(t, err)
	log.Debugf("%v", len(convertedResponse.Responses))

	// Store processed txs into the batch
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	processingContext := state.ProcessingContext{
		BatchNumber:    processBatchRequest.BatchNum,
		Coinbase:       common.Address{},
		Timestamp:      time.Now(),
		globalExitRoot: common.BytesToHash(processBatchRequest.globalExitRoot),
	}

	err = testState.OpenBatch(ctx, processingContext, dbTx)
	require.NoError(t, err)

	err = testState.StoreTransactions(ctx, processBatchRequest.BatchNum, convertedResponse.Responses, dbTx)
	require.NoError(t, err)

	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   processBatchRequest.BatchNum,
		StateRoot:     convertedResponse.NewStateRoot,
		LocalExitRoot: convertedResponse.NewLocalExitRoot,
	}

	err = testState.CloseBatch(ctx, processingReceipt, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// Retrieve Value
	tx2 := types.NewTransaction(2, scAddress, new(big.Int), 80000, new(big.Int).SetUint64(0), common.Hex2Bytes("2e64cec1"))
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	estimatedGas, _, err := testState.EstimateGas(signedTx2, sequencerAddress, nil, nil)
	require.NoError(t, err)
	log.Debugf("Estimated gas = %v", estimatedGas)

	tx2 = types.NewTransaction(2, scAddress, new(big.Int), estimatedGas, new(big.Int).SetUint64(0), common.Hex2Bytes("2e64cec1"))
	signedTx2, err = auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	batchL2Data, err = state.EncodeTransactions([]types.Transaction{*signedTx2})
	require.NoError(t, err)

	processBatchRequest = &executor.ProcessBatchRequest{
		BatchNum:         2,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     processBatchResponse.NewStateRoot,
		globalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldLocalExitRoot: common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err = executorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.Equal(t, pb.Error_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)
	assert.LessOrEqual(t, processBatchResponse.Responses[0].GasUsed, estimatedGas)
	assert.NotEqual(t, uint64(0), processBatchResponse.Responses[0].GasRefunded)
	assert.Equal(t, new(big.Int).SetInt64(123456), new(big.Int).SetBytes(processBatchResponse.Responses[0].ReturnValue))
}
*/

func TestExecutorGasEstimationMultisig(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetInt64(1000)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var erc20SCAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var multisigSCAddress = common.HexToAddress("0x85e844b762a271022b692cf99ce5c59ba0650ac8")
	var multisigParameter = "00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000617b3a3528F9cDd6630fd3301B9c8911F7Bf063D000000000000000000000000B2D0a21D2b14679331f67F3FAB36366ef2270312000000000000000000000000B2bF7Ef15AFfcd23d99A9FB41a310992a70Ed7720000000000000000000000005b6C62FF5dC5De57e9B1a36B64BE3ef4Ac9b08fb"
	var sequencerBalance = 4000000
	scERC20ByteCode, err := testutils.ReadBytecode("../compiled/ERC20Token/ERC20Token.bin")
	require.NoError(t, err)
	scMultiSigByteCode, err := testutils.ReadBytecode("../compiled/MultiSigWallet/MultiSigWallet.bin")
	require.NoError(t, err)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// Deploy contract
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scERC20ByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Deploy contract
	tx1 := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scMultiSigByteCode + multisigParameter),
	})

	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	// Transfer Ownership
	tx2 := types.NewTransaction(2, erc20SCAddress, new(big.Int), 80000, new(big.Int).SetUint64(0), common.Hex2Bytes("f2fde38b00000000000000000000000085e844b762a271022b692cf99ce5c59ba0650ac8"))
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	// Transfer balance to multisig smart contract
	tx3 := types.NewTx(&types.LegacyTx{
		Nonce:    3,
		To:       &multisigSCAddress,
		Value:    new(big.Int).SetUint64(1000000000),
		Gas:      uint64(30000),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     nil,
	})
	signedTx3, err := auth.Signer(auth.From, tx3)
	require.NoError(t, err)

	// Submit Transaction
	tx4 := types.NewTransaction(4, multisigSCAddress, new(big.Int), 150000, new(big.Int).SetUint64(0), common.Hex2Bytes("c64274740000000000000000000000001275fbb540c8efc58b812ba83b0d0b8b9917ae98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000014352ca32838ab928d9e55bd7d1a39cb7fbd453ab1000000000000000000000000"))
	signedTx4, err := auth.Signer(auth.From, tx4)
	require.NoError(t, err)

	// Confirm transaction
	tx5 := types.NewTransaction(5, multisigSCAddress, new(big.Int), 150000, new(big.Int).SetUint64(0), common.Hex2Bytes("c01a8c840000000000000000000000000000000000000000000000000000000000000000"))
	signedTx5, err := auth.Signer(auth.From, tx5)
	require.NoError(t, err)

	transactions := []types.Transaction{*signedTx0, *signedTx1, *signedTx2, *signedTx3, *signedTx4, *signedTx5}
	effectivePercentages := make([]uint8, 0, len(transactions))
	for range transactions {
		effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
	}
	batchL2Data, err := state.EncodeTransactions(transactions, effectivePercentages, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[1].Error)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[2].Error)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[3].Error)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[4].Error)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[5].Error)

	// Check SC code
	// Check Smart Contracts Code
	stateTree := testState.GetTree()
	code, err := stateTree.GetCode(ctx, erc20SCAddress, processBatchResponse.NewStateRoot)
	require.NoError(t, err)
	require.NotEmpty(t, code)
	code, err = stateTree.GetCode(ctx, multisigSCAddress, processBatchResponse.NewStateRoot)
	require.NoError(t, err)
	require.NotEmpty(t, code)

	// Check Smart Contract Balance
	balance, err := stateTree.GetBalance(ctx, multisigSCAddress, processBatchResponse.NewStateRoot)
	require.NoError(t, err)
	require.Equal(t, uint64(1000000000), balance.Uint64())

	// Preparation to be able to estimate gas
	convertedResponse, err := testState.TestConvertToProcessBatchResponse(processBatchResponse)
	require.NoError(t, err)
	log.Debugf("%v", len(convertedResponse.BlockResponses))

	// Store processed txs into the batch
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	processingContext := state.ProcessingContext{
		BatchNumber:    processBatchRequest.OldBatchNum + 1,
		Coinbase:       common.Address{},
		Timestamp:      time.Now(),
		GlobalExitRoot: common.BytesToHash(processBatchRequest.GlobalExitRoot),
	}

	err = testState.OpenBatch(ctx, processingContext, dbTx)
	require.NoError(t, err)

	err = testState.StoreTransactions(ctx, processBatchRequest.OldBatchNum+1, convertedResponse.BlockResponses, nil, dbTx)
	require.NoError(t, err)

	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   processBatchRequest.OldBatchNum + 1,
		StateRoot:     convertedResponse.NewStateRoot,
		LocalExitRoot: convertedResponse.NewLocalExitRoot,
	}

	err = testState.CloseBatch(ctx, processingReceipt, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// Revoke Confirmation
	tx6 := types.NewTransaction(6, multisigSCAddress, new(big.Int), 50000, new(big.Int).SetUint64(0), common.Hex2Bytes("20ea8d860000000000000000000000000000000000000000000000000000000000000000"))
	signedTx6, err := auth.Signer(auth.From, tx6)
	require.NoError(t, err)

	blockNumber, err := testState.GetLastL2BlockNumber(ctx, nil)
	require.NoError(t, err)

	estimatedGas, _, err := testState.EstimateGas(signedTx6, sequencerAddress, &blockNumber, nil)
	require.NoError(t, err)
	log.Debugf("Estimated gas = %v", estimatedGas)

	tx6 = types.NewTransaction(6, multisigSCAddress, new(big.Int), estimatedGas, new(big.Int).SetUint64(0), common.Hex2Bytes("20ea8d860000000000000000000000000000000000000000000000000000000000000000"))
	signedTx6, err = auth.Signer(auth.From, tx6)
	require.NoError(t, err)

	batchL2Data, err = state.EncodeTransactions([]types.Transaction{*signedTx6}, constants.EffectivePercentage, forkID)
	require.NoError(t, err)

	processBatchRequest = &executor.ProcessBatchRequest{
		OldBatchNum:      1,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     processBatchResponse.NewStateRoot,
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err = test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)
	log.Debugf("Used gas = %v", processBatchResponse.Responses[0].GasUsed)
}

func TestExecuteWithoutUpdatingMT(t *testing.T) {
	ctx := context.Background()
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	var chainIDSequencer = new(big.Int).SetUint64(stateCfg.ChainID)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var gasLimit = uint64(4000000)
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	scByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)

	// auth
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	// signed tx to deploy SC
	unsignedTxDeploy := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     common.Hex2Bytes(scByteCode),
	})
	signedTxDeploy, err := auth.Signer(auth.From, unsignedTxDeploy)
	require.NoError(t, err)

	signedTxs := []types.Transaction{
		*signedTxDeploy,
	}

	batchL2Data, err := state.EncodeTransactions(signedTxs, constants.EffectivePercentage, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 0,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	// assert signed tx do deploy sc
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)
	assert.Equal(t, scAddress, common.HexToAddress(processBatchResponse.Responses[0].CreateAddress))

	log.Debug(processBatchResponse)

	incrementFnSignature := crypto.Keccak256Hash([]byte("increment()")).Bytes()[:4]
	retrieveFnSignature := crypto.Keccak256Hash([]byte("getCount()")).Bytes()[:4]

	// signed tx to call SC
	unsignedTxFirstIncrement := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		To:       &scAddress,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     incrementFnSignature,
	})

	signedTxFirstIncrement, err := auth.Signer(auth.From, unsignedTxFirstIncrement)
	require.NoError(t, err)

	unsignedTxFirstRetrieve := types.NewTx(&types.LegacyTx{
		Nonce:    2,
		To:       &scAddress,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		Data:     retrieveFnSignature,
	})

	signedTxFirstRetrieve, err := auth.Signer(auth.From, unsignedTxFirstRetrieve)
	require.NoError(t, err)

	signedTxs2 := []types.Transaction{
		*signedTxFirstIncrement,
		*signedTxFirstRetrieve,
	}

	batchL2Data2, err := state.EncodeTransactions(signedTxs2, constants.TwoEffectivePercentages, forkID)
	require.NoError(t, err)

	// Create Batch 2
	processBatchRequest = &executor.ProcessBatchRequest{
		OldBatchNum:      1,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data2,
		OldStateRoot:     processBatchResponse.NewStateRoot,
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 0,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	processBatchResponse, err = test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	log.Debug(processBatchResponse)

	// assert signed tx to increment counter
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[0].Error)

	// assert signed tx to increment counter
	assert.Equal(t, executor.RomError_ROM_ERROR_NO_ERROR, processBatchResponse.Responses[1].Error)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", hex.EncodeToString(processBatchResponse.Responses[1].ReturnValue))
}

func TestExecutorUnsignedTransactionsWithCorrectL2BlockStateRoot(t *testing.T) {
	ctx := context.Background()
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	// auth
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(operations.DefaultSequencerPrivateKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(stateCfg.ChainID))
	require.NoError(t, err)

	auth.Nonce = big.NewInt(0)
	auth.Value = nil
	auth.GasPrice = big.NewInt(0)
	auth.GasLimit = uint64(4000000)
	auth.NoSend = true

	// Set genesis
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: operations.DefaultSequencerAddress,
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
	}
	_, err = testState.SetGenesis(ctx, state.Block{}, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)

	scAddr, scTx, sc, err := Counter.DeployCounter(auth, &ethclient.Client{})
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(context.Background()))

	// deploy SC
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	batchCtx := state.ProcessingContext{
		BatchNumber: 1,
		Coinbase:    common.HexToAddress(operations.DefaultSequencerAddress),
		Timestamp:   time.Now(),
	}
	err = testState.OpenBatch(context.Background(), batchCtx, dbTx)
	require.NoError(t, err)
	signedTxs := []types.Transaction{
		*scTx,
	}
	effectivePercentages := make([]uint8, 0, len(signedTxs))
	for range signedTxs {
		effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
	}

	batchL2Data, err := state.EncodeTransactions(signedTxs, effectivePercentages, forkID)
	require.NoError(t, err)

	processBatchResponse, err := testState.ProcessSequencerBatch(context.Background(), 1, batchL2Data, metrics.SequencerCallerLabel, dbTx)
	require.NoError(t, err)
	// assert signed tx do deploy sc
	assert.Nil(t, processBatchResponse.BlockResponses[0].TransactionResponses[0].RomError)
	assert.NotEqual(t, state.ZeroAddress, processBatchResponse.BlockResponses[0].TransactionResponses[0].CreateAddress.Hex())
	assert.Equal(t, scAddr.Hex(), processBatchResponse.BlockResponses[0].TransactionResponses[0].CreateAddress.Hex())

	// assert signed tx to increment counter
	assert.Nil(t, processBatchResponse.BlockResponses[0].TransactionResponses[0].RomError)

	// Add txs to DB
	err = testState.StoreTransactions(context.Background(), 1, processBatchResponse.BlockResponses, nil, dbTx)
	require.NoError(t, err)
	// Close batch
	err = testState.CloseBatch(
		context.Background(),
		state.ProcessingReceipt{
			BatchNumber:   1,
			StateRoot:     processBatchResponse.NewStateRoot,
			LocalExitRoot: processBatchResponse.NewLocalExitRoot,
		}, dbTx,
	)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(context.Background()))

	// increment
	for n := int64(1); n <= 3; n++ {
		batchNumber := uint64(n + 1)
		dbTx, err := testState.BeginStateTransaction(ctx)
		require.NoError(t, err)

		auth.Nonce = big.NewInt(n)
		tx, err := sc.Increment(auth)
		require.NoError(t, err)

		batchCtx := state.ProcessingContext{
			BatchNumber: batchNumber,
			Coinbase:    common.HexToAddress(operations.DefaultSequencerAddress),
			Timestamp:   time.Now(),
		}
		err = testState.OpenBatch(context.Background(), batchCtx, dbTx)
		require.NoError(t, err)
		signedTxs := []types.Transaction{
			*tx,
		}
		effectivePercentages := make([]uint8, 0, len(signedTxs))
		for range signedTxs {
			effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
		}

		batchL2Data, err := state.EncodeTransactions(signedTxs, effectivePercentages, forkID)
		require.NoError(t, err)

		processBatchResponse, err := testState.ProcessSequencerBatch(context.Background(), batchNumber, batchL2Data, metrics.SequencerCallerLabel, dbTx)
		require.NoError(t, err)
		// assert signed tx to increment counter
		assert.Nil(t, processBatchResponse.BlockResponses[0].TransactionResponses[0].RomError)

		// Add txs to DB
		err = testState.StoreTransactions(context.Background(), batchNumber, processBatchResponse.BlockResponses, nil, dbTx)
		require.NoError(t, err)
		// Close batch
		err = testState.CloseBatch(
			context.Background(),
			state.ProcessingReceipt{
				BatchNumber:   batchNumber,
				StateRoot:     processBatchResponse.NewStateRoot,
				LocalExitRoot: processBatchResponse.NewLocalExitRoot,
			}, dbTx,
		)
		require.NoError(t, err)
		require.NoError(t, dbTx.Commit(context.Background()))
	}

	getCountFnSignature := crypto.Keccak256Hash([]byte("getCount()")).Bytes()[:4]
	getCountUnsignedTx := types.NewTx(&types.LegacyTx{
		To:   &processBatchResponse.BlockResponses[0].TransactionResponses[0].CreateAddress,
		Gas:  uint64(100000),
		Data: getCountFnSignature,
	})

	l2BlockNumber := uint64(1)
	result, err := testState.ProcessUnsignedTransaction(context.Background(), getCountUnsignedTx, auth.From, &l2BlockNumber, true, nil)
	require.NoError(t, err)
	// assert unsigned tx
	assert.Nil(t, result.Err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(result.ReturnValue))

	l2BlockNumber = uint64(2)
	result, err = testState.ProcessUnsignedTransaction(context.Background(), getCountUnsignedTx, auth.From, &l2BlockNumber, true, nil)
	require.NoError(t, err)
	// assert unsigned tx
	assert.Nil(t, result.Err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", hex.EncodeToString(result.ReturnValue))

	l2BlockNumber = uint64(3)
	result, err = testState.ProcessUnsignedTransaction(context.Background(), getCountUnsignedTx, auth.From, &l2BlockNumber, true, nil)
	require.NoError(t, err)
	// assert unsigned tx
	assert.Nil(t, result.Err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000002", hex.EncodeToString(result.ReturnValue))

	l2BlockNumber = uint64(4)
	result, err = testState.ProcessUnsignedTransaction(context.Background(), getCountUnsignedTx, auth.From, &l2BlockNumber, true, nil)
	require.NoError(t, err)
	// assert unsigned tx
	assert.Nil(t, result.Err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000003", hex.EncodeToString(result.ReturnValue))
}

func TestBigDataTx(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 4000000

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       &sequencerAddress,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     make([]byte, 120000), // large data
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	// Encode transaction
	batchL2Data, err := state.EncodeTransaction(*signedTx, state.MaxEffectivePercentage, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	response, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	require.Equal(t, executor.ExecutorError_EXECUTOR_ERROR_INVALID_BATCH_L2_DATA, response.Error)
}

func TestExecutorTxHashAndRLP(t *testing.T) {
	ctx := context.Background()
	// Test Case
	type TxHashTestCase struct {
		Nonce    string `json:"nonce"`
		GasPrice string `json:"gasPrice"`
		GasLimit string `json:"gasLimit"`
		To       string `json:"to"`
		Value    string `json:"value"`
		Data     string `json:"data"`
		ChainID  string `json:"chainId"`
		V        string `json:"v"`
		R        string `json:"r"`
		S        string `json:"s"`
		From     string `json:"from"`
		Hash     string `json:"hash"`
		Link     string `json:"link"`
	}

	var testCases, testCases2 []TxHashTestCase

	jsonFile, err := os.Open(filepath.Clean("../../test/vectors/src/tx-hash-ethereum/uniswap_formated.json"))
	require.NoError(t, err)
	defer func() { _ = jsonFile.Close() }()

	bytes, err := io.ReadAll(jsonFile)
	require.NoError(t, err)

	err = json.Unmarshal(bytes, &testCases)
	require.NoError(t, err)

	jsonFile2, err := os.Open(filepath.Clean("../../test/vectors/src/tx-hash-ethereum/rlp.json"))
	require.NoError(t, err)
	defer func() { _ = jsonFile2.Close() }()

	bytes2, err := io.ReadAll(jsonFile2)
	require.NoError(t, err)

	err = json.Unmarshal(bytes2, &testCases2)
	require.NoError(t, err)
	testCases = append(testCases, testCases2...)

	for x, testCase := range testCases {
		var stateRoot = state.ZeroHash
		var receiverAddress = common.HexToAddress(testCase.To)
		receiver := &receiverAddress
		if testCase.To == "0x" {
			receiver = nil
		}

		v, ok := new(big.Int).SetString(testCase.V, 0)
		require.Equal(t, true, ok)

		r, ok := new(big.Int).SetString(testCase.R, 0)
		require.Equal(t, true, ok)

		s, ok := new(big.Int).SetString(testCase.S, 0)
		require.Equal(t, true, ok)

		var value *big.Int

		if testCase.Value != "0x" {
			value, ok = new(big.Int).SetString(testCase.Value, 0)
			require.Equal(t, true, ok)
		}

		gasPrice, ok := new(big.Int).SetString(testCase.GasPrice, 0)
		require.Equal(t, true, ok)

		gasLimit, ok := new(big.Int).SetString(testCase.GasLimit, 0)
		require.Equal(t, true, ok)

		nonce, ok := new(big.Int).SetString(testCase.Nonce, 0)
		require.Equal(t, true, ok)

		// Create transaction
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce.Uint64(),
			To:       receiver,
			Value:    value,
			Gas:      gasLimit.Uint64(),
			GasPrice: gasPrice,
			Data:     common.Hex2Bytes(strings.TrimPrefix(testCase.Data, "0x")),
			V:        v,
			R:        r,
			S:        s,
		})
		t.Log("chainID: ", tx.ChainId())
		t.Log("txHash: ", tx.Hash())

		require.Equal(t, testCase.Hash, tx.Hash().String())

		batchL2Data, err := state.EncodeTransactions([]types.Transaction{*tx}, constants.EffectivePercentage, forkID)
		require.NoError(t, err)

		// Create Batch
		processBatchRequest := &executor.ProcessBatchRequest{
			OldBatchNum:      uint64(x),
			Coinbase:         receiverAddress.String(),
			BatchL2Data:      batchL2Data,
			OldStateRoot:     stateRoot.Bytes(),
			GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
			OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
			EthTimestamp:     uint64(0),
			UpdateMerkleTree: 1,
			ChainId:          stateCfg.ChainID,
			ForkId:           forkID,
			ContextId:        uuid.NewString(),
		}

		// Process batch
		processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
		require.NoError(t, err)

		// TX Hash
		log.Debugf("TX Hash=%v", tx.Hash().String())
		log.Debugf("Response TX Hash=%v", common.BytesToHash(processBatchResponse.Responses[0].TxHash).String())

		// RPL Encoding
		b, err := tx.MarshalBinary()
		require.NoError(t, err)
		log.Debugf("TX RLP=%v", hex.EncodeToHex(b))
		log.Debugf("Response TX RLP=%v", "0x"+common.Bytes2Hex(processBatchResponse.Responses[0].RlpTx))

		require.Equal(t, tx.Hash(), common.BytesToHash(processBatchResponse.Responses[0].TxHash))
		require.Equal(t, hex.EncodeToHex(b), "0x"+common.Bytes2Hex(processBatchResponse.Responses[0].RlpTx))
	}
}

func TestExecuteTransaction(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 4000000
	scCounterByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)

	// Deploy counter.sol
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scCounterByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	// Encode transaction
	v, r, s := signedTx.RawSignatureValues()
	sign := 1 - (v.Uint64() & 1)

	txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
		signedTx.Nonce(),
		signedTx.GasPrice(),
		signedTx.Gas(),
		signedTx.To(),
		signedTx.Value(),
		signedTx.Data(),
		signedTx.ChainId(), uint(0), uint(0),
	})
	require.NoError(t, err)

	newV := new(big.Int).Add(big.NewInt(test.Ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	batchL2Data, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 1,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	log.Debugf("%v", processBatchRequest)

	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	log.Debug(processBatchResponse)
	// TODO: assert processBatchResponse to make sure that the response makes sense
}

func TestExecutorInvalidNonce(t *testing.T) {
	ctx := context.Background()
	chainID := new(big.Int).SetInt64(1000)
	senderPvtKey := "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	receiverAddress := common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB")

	// authorization
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	senderAddress := auth.From

	type testCase struct {
		name         string
		currentNonce uint64
		txNonce      uint64
	}

	testCases := []testCase{
		{
			name:         "tx nonce is greater than expected",
			currentNonce: 1,
			txNonce:      2,
		},
		{
			name:         "tx nonce is less than expected",
			currentNonce: 5,
			txNonce:      4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			test.InitOrResetDB(test.StateDBCfg)

			// Set Genesis
			block := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}
			test.Genesis.Actions = []*state.GenesisAction{
				{
					Address: senderAddress.String(),
					Type:    int(merkletree.LeafTypeBalance),
					Value:   "10000000",
				},
				{
					Address: senderAddress.String(),
					Type:    int(merkletree.LeafTypeNonce),
					Value:   strconv.FormatUint(testCase.currentNonce, encoding.Base10),
				},
			}
			dbTx, err := testState.BeginStateTransaction(ctx)
			require.NoError(t, err)
			stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

			stateTree := testState.GetTree()

			// Read Sender Balance
			currentNonce, err := stateTree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
			require.NoError(t, err)
			assert.Equal(t, testCase.currentNonce, currentNonce.Uint64())

			// Create transaction
			tx := types.NewTransaction(testCase.txNonce, receiverAddress, new(big.Int).SetUint64(2), uint64(30000), big.NewInt(1), nil)
			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			// encode txs
			batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx}, constants.EffectivePercentage, forkID)
			require.NoError(t, err)

			// Create Batch
			processBatchRequest := &executor.ProcessBatchRequest{
				OldBatchNum:      1,
				Coinbase:         receiverAddress.String(),
				BatchL2Data:      batchL2Data,
				OldStateRoot:     stateRoot.Bytes(),
				GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
				OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
				EthTimestamp:     uint64(0),
				UpdateMerkleTree: 1,
				ChainId:          chainID.Uint64(),
				ForkId:           forkID,
				ContextId:        uuid.NewString(),
			}

			// Process batch
			processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
			require.NoError(t, err)

			transactionResponses := processBatchResponse.GetResponses()
			assert.Equal(t, true, executor.IsIntrinsicError(transactionResponses[0].Error), "invalid tx Error, it is expected to be INVALID TX")
		})
	}
}

func TestExecutorRevert(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetInt64(1000)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var sequencerBalance = 4000000
	scRevertByteCode, err := testutils.ReadBytecode("Revert2/Revert2.bin")
	require.NoError(t, err)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: sequencerAddress.String(),
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "10000000",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)

	// Deploy revert.sol
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scRevertByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Call SC method
	tx1 := types.NewTransaction(1, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("4abbb40a"))
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx0, *signedTx1}, constants.TwoEffectivePercentages, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      1,
		Coinbase:         sequencerAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(time.Now().Unix()),
		UpdateMerkleTree: 0,
		ChainId:          stateCfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}
	fmt.Println("batchL2Data: ", batchL2Data)
	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)
	assert.Equal(t, runtime.ErrExecutionReverted, executor.RomErr(processBatchResponse.Responses[1].Error))

	// Unsigned
	receipt := &types.Receipt{
		Type:              signedTx0.Type(),
		PostState:         processBatchResponse.Responses[0].StateRoot,
		CumulativeGasUsed: processBatchResponse.Responses[0].GasUsed,
		BlockNumber:       big.NewInt(0),
		GasUsed:           processBatchResponse.Responses[0].GasUsed,
		TxHash:            signedTx0.Hash(),
		TransactionIndex:  0,
		Status:            types.ReceiptStatusSuccessful,
	}

	receipt1 := &types.Receipt{
		Type:              signedTx1.Type(),
		PostState:         processBatchResponse.Responses[1].StateRoot,
		CumulativeGasUsed: processBatchResponse.Responses[0].GasUsed + processBatchResponse.Responses[1].GasUsed,
		BlockNumber:       big.NewInt(0),
		GasUsed:           signedTx1.Gas(),
		TxHash:            signedTx1.Hash(),
		TransactionIndex:  1,
		Status:            types.ReceiptStatusSuccessful,
	}

	header := state.NewL2Header(&types.Header{
		Number:     big.NewInt(2),
		ParentHash: state.ZeroHash,
		Coinbase:   state.ZeroAddress,
		Root:       common.BytesToHash(processBatchResponse.NewStateRoot),
		GasUsed:    receipt1.GasUsed,
		GasLimit:   receipt1.GasUsed,
		Time:       uint64(time.Now().Unix()),
	})

	receipts := []*types.Receipt{receipt, receipt1}
	imStateRoots := []common.Hash{common.BytesToHash(processBatchResponse.Responses[0].StateRoot), common.BytesToHash(processBatchResponse.Responses[1].StateRoot)}

	transactions := []*types.Transaction{signedTx0, signedTx1}

	// Create block to be able to calculate its hash
	st := trie.NewStackTrie(nil)
	l2Block := state.NewL2Block(header, transactions, []*state.L2Header{}, receipts, st)
	l2Block.ReceivedAt = time.Now()

	receipt.BlockHash = l2Block.Hash()
	receipt1.BlockHash = l2Block.Hash()

	numTxs := len(transactions)
	storeTxsEGPData := make([]state.StoreTxEGPData, numTxs)
	txsL2Hash := make([]common.Hash, numTxs)
	for i := range transactions {
		storeTxsEGPData[i] = state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: state.MaxEffectivePercentage}
		txsL2Hash[i] = common.HexToHash(fmt.Sprintf("0x%d", i))
	}

	err = testState.AddL2Block(ctx, 0, l2Block, receipts, txsL2Hash, storeTxsEGPData, imStateRoots, dbTx)
	require.NoError(t, err)
	l2Block, err = testState.GetL2BlockByHash(ctx, l2Block.Hash(), dbTx)
	require.NoError(t, err)

	require.NoError(t, dbTx.Commit(ctx))

	lastL2BlockNumber := l2Block.NumberU64()

	unsignedTx := types.NewTransaction(2, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("4abbb40a"))

	result, err := testState.ProcessUnsignedTransaction(ctx, unsignedTx, auth.From, &lastL2BlockNumber, false, nil)
	require.NoError(t, err)
	require.NotNil(t, result.Err)
	assert.Equal(t, fmt.Errorf("execution reverted: Today is not juernes").Error(), result.Err.Error())
}

func TestExecutorTransfer(t *testing.T) {
	ctx := context.Background()
	var senderAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var senderPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var receiverAddress = common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB")
	var chainID = new(big.Int).SetUint64(stateCfg.ChainID)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "10000000",
		},
	}
	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	// Create transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       &receiverAddress,
		Value:    new(big.Int).SetUint64(2),
		Gas:      uint64(30000),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     nil,
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx}, constants.EffectivePercentage, forkID)
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      0,
		Coinbase:         receiverAddress.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldAccInputHash:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:     uint64(0),
		UpdateMerkleTree: 1,
		ChainId:          chainID.Uint64(),
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	stateTree := testState.GetTree()

	// Read Sender Balance before execution
	balance, err := stateTree.GetBalance(ctx, senderAddress, processBatchRequest.OldStateRoot)
	require.NoError(t, err)
	require.Equal(t, uint64(10000000), balance.Uint64())

	// Read Receiver Balance before execution
	balance, err = stateTree.GetBalance(ctx, receiverAddress, processBatchRequest.OldStateRoot)
	require.NoError(t, err)
	require.Equal(t, uint64(0), balance.Uint64())

	// Process batch
	processBatchResponse, err := test.ExecutorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	// Read Sender Balance
	balance, err = stateTree.GetBalance(ctx, senderAddress, processBatchResponse.Responses[0].StateRoot)
	require.NoError(t, err)
	require.Equal(t, uint64(9978998), balance.Uint64())

	// Read Receiver Balance
	balance, err = stateTree.GetBalance(ctx, receiverAddress, processBatchResponse.Responses[0].StateRoot)
	require.NoError(t, err)
	require.Equal(t, uint64(21002), balance.Uint64())

	// Read Modified Addresses directly from response
	readWriteAddresses := processBatchResponse.ReadWriteAddresses
	log.Debug(receiverAddress.String())
	data := readWriteAddresses[strings.ToLower(receiverAddress.String())]
	require.Equal(t, "21002", data.Balance)

	// Read Modified Addresses from converted response
	converted, err := testState.TestConvertToProcessBatchResponse(processBatchResponse)
	require.NoError(t, err)
	convertedData := converted.ReadWriteAddresses[receiverAddress]
	require.Equal(t, uint64(21002), convertedData.Balance.Uint64())
	require.Equal(t, receiverAddress, convertedData.Address)
	require.Equal(t, (*uint64)(nil), convertedData.Nonce)
}
