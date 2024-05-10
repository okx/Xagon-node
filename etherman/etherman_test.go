package etherman

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevmbridge"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	forkID6 = 6
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (ethman *Client, ethBackend *simulated.Backend, auth *bind.TransactOpts, polAddr common.Address, br *etrogpolygonzkevmbridge.Etrogpolygonzkevmbridge) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, ethBackend, polAddr, br, err = NewSimulatedEtherman(Config{ForkIDChunkSize: 10}, auth)
	if err != nil {
		log.Fatal(err)
	}
	err = ethman.AddOrReplaceAuth(*auth)
	if err != nil {
		log.Fatal(err)
	}
	return ethman, ethBackend, auth, polAddr, br
}

func TestGEREvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount := big.NewInt(1000000000000000)
	auth.Value = amount
	_, err = br.BridgeAsset(auth, 1, auth.From, amount, common.Address{}, true, []byte{})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, uint64(8), blocks[0].L1InfoTree[0].BlockNumber)
	assert.NotEqual(t, common.Hash{}, blocks[0].L1InfoTree[0].MainnetExitRoot)
	assert.Equal(t, common.Hash{}, blocks[0].L1InfoTree[0].RollupExitRoot)
}

func TestForcedBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount, err := etherman.EtrogRollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.EtrogZkEVM.ForceBatch(auth, data, amount)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, uint64(8), blocks[0].BlockNumber)
	assert.Equal(t, uint64(8), blocks[0].ForcedBatches[0].BlockNumber)
	assert.NotEqual(t, common.Hash{}, blocks[0].ForcedBatches[0].GlobalExitRoot)
	assert.NotEqual(t, time.Time{}, blocks[0].ForcedBatches[0].ForcedAt)
	assert.Equal(t, uint64(1), blocks[0].ForcedBatches[0].ForcedBatchNumber)
	assert.Equal(t, rawTxs, hex.EncodeToString(blocks[0].ForcedBatches[0].RawTxsData))
	assert.Equal(t, auth.From, blocks[0].ForcedBatches[0].Sequencer)
}

func TestSequencedBatchesEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Make a bridge tx
	auth.Value = big.NewInt(1000000000000000)
	_, err = br.BridgeAsset(auth, 1, auth.From, auth.Value, common.Address{}, true, []byte{})
	require.NoError(t, err)
	ethBackend.Commit()
	auth.Value = big.NewInt(0)

	amount, err := etherman.EtrogRollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.EtrogZkEVM.ForceBatch(auth, data, amount)
	require.NoError(t, err)
	require.NoError(t, err)
	ethBackend.Commit()

	// Now read the event
	currentBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	currentBlockNumber := currentBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &currentBlockNumber)
	require.NoError(t, err)
	t.Log("Blocks: ", blocks)
	var sequences []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData
	sequences = append(sequences, etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
		Transactions: common.Hex2Bytes(rawTxs),
	}, etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
		Transactions: common.Hex2Bytes(rawTxs),
	})
	_, err = etherman.EtrogZkEVM.SequenceBatches(auth, sequences, uint64(time.Now().Unix()), uint64(1), auth.From)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, 3, len(blocks))
	assert.Equal(t, 1, len(blocks[2].SequencedBatches))
	assert.Equal(t, common.Hex2Bytes(rawTxs), blocks[2].SequencedBatches[0][1].PolygonRollupBaseEtrogBatchData.Transactions)
	assert.Equal(t, uint64(0), blocks[2].SequencedBatches[0][0].ForcedTimestamp)
	assert.Equal(t, [32]byte{}, blocks[2].SequencedBatches[0][0].ForcedGlobalExitRoot)
	assert.Equal(t, auth.From, blocks[2].SequencedBatches[0][0].Coinbase)
	assert.Equal(t, auth.From, blocks[2].SequencedBatches[0][0].SequencerAddr)
	assert.NotEqual(t, common.Hash{}, blocks[2].SequencedBatches[0][0].ForcedBlockHashL1)
	assert.Equal(t, 0, order[blocks[2].BlockHash][0].Pos)
}

func TestVerifyBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()

	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	tx := etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
		Transactions: common.Hex2Bytes(rawTxs),
	}
	//TODO: Fix params
	_, err = etherman.EtrogZkEVM.SequenceBatches(auth, []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{tx}, uint64(time.Now().Unix()), uint64(1), auth.From)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	_, err = etherman.EtrogRollupManager.VerifyBatchesTrustedAggregator(auth, 1, uint64(0), uint64(0), uint64(1), [32]byte{}, [32]byte{}, auth.From, [24][32]byte{})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v, \nOrder: %+v", blocks, order)
	assert.Equal(t, uint64(9), blocks[1].BlockNumber)
	assert.Equal(t, uint64(1), blocks[1].VerifiedBatches[0].BatchNumber)
	assert.NotEqual(t, common.Address{}, blocks[1].VerifiedBatches[0].Aggregator)
	assert.NotEqual(t, common.Hash{}, blocks[1].VerifiedBatches[0].TxHash)
	assert.Equal(t, L1InfoTreeOrder, order[blocks[1].BlockHash][1].Name)
	assert.Equal(t, VerifyBatchOrder, order[blocks[1].BlockHash][0].Name)
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
	assert.Equal(t, 0, order[blocks[1].BlockHash][1].Pos)
}

func TestSequenceForceBatchesEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount, err := etherman.EtrogRollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	data, err := hex.DecodeString(rawTxs)
	require.NoError(t, err)
	_, err = etherman.EtrogZkEVM.ForceBatch(auth, data, amount)
	require.NoError(t, err)
	ethBackend.Commit()
	ethBackend.Commit()

	err = ethBackend.AdjustTime((24*7 + 1) * time.Hour)
	require.NoError(t, err)
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)

	forcedGer := blocks[0].ForcedBatches[0].GlobalExitRoot
	forcedTimestamp := uint64(blocks[0].ForcedBatches[0].ForcedAt.Unix())
	prevBlock, err := etherman.EthClient.BlockByNumber(ctx, big.NewInt(0).SetUint64(blocks[0].BlockNumber-1))
	require.NoError(t, err)
	forcedBlockHashL1 := prevBlock.Hash()
	forceBatchData := etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
		Transactions:         blocks[0].ForcedBatches[0].RawTxsData,
		ForcedGlobalExitRoot: forcedGer,
		ForcedTimestamp:      forcedTimestamp,
		ForcedBlockHashL1:    forcedBlockHashL1,
	}
	_, err = etherman.EtrogZkEVM.SequenceForceBatches(auth, []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{forceBatchData})
	require.NoError(t, err)
	ethBackend.Commit()

	// Now read the event
	finalBlock, err = etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber = finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, uint64(12), blocks[1].BlockNumber)
	assert.Equal(t, uint64(2), blocks[1].SequencedForceBatches[0][0].BatchNumber)
	assert.Equal(t, forcedGer, common.BytesToHash(blocks[1].SequencedForceBatches[0][0].ForcedGlobalExitRoot[:]))
	assert.Equal(t, forcedTimestamp, blocks[1].SequencedForceBatches[0][0].ForcedTimestamp)
	assert.Equal(t, forcedBlockHashL1, common.BytesToHash(blocks[1].SequencedForceBatches[0][0].ForcedBlockHashL1[:]))
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
}

func TestSendSequences(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Make a bridge tx
	auth.Value = big.NewInt(1000000000000000)
	_, err = br.BridgeAsset(auth, 1, auth.From, auth.Value, common.Address{}, true, []byte{})
	require.NoError(t, err)
	ethBackend.Commit()
	auth.Value = big.NewInt(0)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*tx1}, constants.EffectivePercentage, forkID6)
	require.NoError(t, err)
	sequence := ethmanTypes.Sequence{
		BatchNumber:          0,
		BatchL2Data:          batchL2Data,
		LastL2BLockTimestamp: time.Now().Unix(),
	}
	lastL2BlockTStamp := tx1.Time().Unix()
	// TODO: fix params
	tx, err := etherman.sequenceBatches(*auth, []ethmanTypes.Sequence{sequence}, uint64(lastL2BlockTStamp), uint64(1), auth.From)
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, 2, len(blocks))
	assert.Equal(t, 1, len(blocks[1].SequencedBatches))
	assert.Equal(t, [32]byte{}, blocks[1].SequencedBatches[0][0].ForcedGlobalExitRoot)
	assert.Equal(t, [32]byte{}, blocks[1].SequencedBatches[0][0].ForcedBlockHashL1)
	assert.Equal(t, auth.From, blocks[1].SequencedBatches[0][0].Coinbase)
	assert.Equal(t, auth.From, blocks[1].SequencedBatches[0][0].SequencerAddr)
	assert.Equal(t, uint64(0), blocks[1].SequencedBatches[0][0].ForcedTimestamp)
	assert.Equal(t, 0, order[blocks[1].BlockHash][0].Pos)
}

func TestGasPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	etherscanM := new(etherscanMock)
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}
	ctx := context.Background()

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(1448795322), nil)
	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(1448795321), nil)
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(1448795322), gp)

	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, ethGasStationM}

	gp = etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(1448795321), gp)
}

func TestErrorEthGasStationPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, ethGasStationM}
	ctx := context.Background()

	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(0), fmt.Errorf("error getting gasPrice from ethGasStation"))
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(1392695906), gp)

	etherscanM := new(etherscanMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(1448795322), nil)
	gp = etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(1448795322), gp)
}

func TestErrorEtherScanPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	etherscanM := new(etherscanMock)
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}
	ctx := context.Background()

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(0), fmt.Errorf("error getting gasPrice from etherscan"))
	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(1448795321), nil)
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(1448795321), gp)
}

func TestGetForks(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	ctx := context.Background()
	forks, err := etherman.GetForks(ctx, 0, 132)
	require.NoError(t, err)
	assert.Equal(t, 1, len(forks))
	assert.Equal(t, uint64(6), forks[0].ForkId)
	assert.Equal(t, uint64(1), forks[0].FromBatchNumber)
	assert.Equal(t, uint64(math.MaxUint64), forks[0].ToBatchNumber)
	assert.Equal(t, "", forks[0].Version)
	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, 0, &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, 1, len(blocks))
	assert.Equal(t, 1, len(blocks[0].ForkIDs))
	assert.Equal(t, 0, order[blocks[0].BlockHash][0].Pos)
	assert.Equal(t, ForkIDsOrder, order[blocks[0].BlockHash][0].Name)
	assert.Equal(t, uint64(0), blocks[0].ForkIDs[0].BatchNumber)
	assert.Equal(t, uint64(6), blocks[0].ForkIDs[0].ForkID)
	assert.Equal(t, "", blocks[0].ForkIDs[0].Version)
}

func TestProof(t *testing.T) {
	proof := "0x20227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0520227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a05"
	p, err := convertProof(proof)
	require.NoError(t, err)
	str := "20227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a05" //nolint:gosec
	proofReference, err := encoding.DecodeBytes(&str)
	require.NoError(t, err)
	var expected [32]byte
	copy(expected[:], proofReference)
	for i := 0; i < 24; i++ {
		assert.Equal(t, expected, p[i])
	}
	t.Log("Proof: ", p)
}
