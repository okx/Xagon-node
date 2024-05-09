//nolint:deadcode,unused,varcheck
package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/stretchr/testify/require"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	invalidParamsErrorCode = -32602
	toAddressHex           = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	forkID6                = 6
)

var (
	toAddress = common.HexToAddress(toAddressHex)
	opsMan    *operations.Manager
)

var networks = []struct {
	Name         string
	URL          string
	WebSocketURL string
	ChainID      uint64
	PrivateKey   string
}{
	{
		Name:         "Local L1",
		URL:          operations.DefaultL1NetworkURL,
		WebSocketURL: operations.DefaultL1NetworkWebSocketURL,
		ChainID:      operations.DefaultL1ChainID,
		PrivateKey:   operations.DefaultSequencerPrivateKey,
	},
	{
		Name:         "Local L2",
		URL:          operations.DefaultL2NetworkURL,
		WebSocketURL: operations.DefaultL2NetworkWebSocketURL,
		ChainID:      operations.DefaultL2ChainID,
		PrivateKey:   operations.DefaultSequencerPrivateKey,
	},
}

func setup() {
	var err error
	ctx := context.Background()
	err = operations.Teardown()
	if err != nil {
		panic(err)
	}

	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err = operations.NewManager(ctx, opsCfg)
	if err != nil {
		panic(err)
	}
	err = opsMan.Setup()
	if err != nil {
		panic(err)
	}
}

func teardown() {
	err := operations.Teardown()
	if err != nil {
		panic(err)
	}
}

func createTX(client *ethclient.Client, auth *bind.TransactOpts, to common.Address, amount *big.Int) (*ethTypes.Transaction, error) {
	nonce, err := client.NonceAt(context.Background(), auth.From, nil)
	if err != nil {
		return nil, err
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &to, Value: amount})
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	log.Infof("\nTX details:\n\tNonce:    %d\n\tGasLimit: %d\n\tGasPrice: %d", nonce, gasLimit, gasPrice)
	if gasLimit != uint64(21000) { //nolint:gomnd
		return nil, fmt.Errorf("gasLimit %d != 21000", gasLimit)
	}
	tx := ethTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}
	log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func logTx(tx *ethTypes.Transaction) {
	sender, _ := state.GetSender(*tx)
	log.Debugf("********************")
	log.Debugf("Hash: %v", tx.Hash())
	log.Debugf("From: %v", sender)
	log.Debugf("Nonce: %v", tx.Nonce())
	log.Debugf("ChainId: %v", tx.ChainId())
	log.Debugf("To: %v", tx.To())
	log.Debugf("Gas: %v", tx.Gas())
	log.Debugf("GasPrice: %v", tx.GasPrice())
	log.Debugf("Cost: %v", tx.Cost())

	// b, _ := tx.MarshalBinary()
	//log.Debugf("RLP: ", hex.EncodeToHex(b))
	log.Debugf("********************")
}

func sendForcedBatchForVector(t *testing.T, txs []byte, opsman *operations.Manager) (*state.Batch, error) {
	ctx := context.Background()
	st := opsman.State()
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)

	// Create smc client
	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
	zkEvm, err := etrogpolygonzkevm.NewEtrogpolygonzkevm(zkEvmAddr, ethClient)
	require.NoError(t, err)

	rollupManagerAddr := common.HexToAddress(operations.DefaultL1RollupManagerSmartContract)
	rollupManager, err := etrogpolygonrollupmanager.NewEtrogpolygonrollupmanager(rollupManagerAddr, ethClient)
	require.NoError(t, err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)

	log.Info("Using address: ", auth.From)
	num, err := zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	log.Info("Number of forceBatches in the smc: ", num)

	// Get tip
	tip, err := rollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	tx, err := zkEvm.SetForceBatchAddress(auth, common.Address{})
	require.NoError(t, err)
	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	// Send forceBatch
	tx, err = zkEvm.ForceBatch(auth, txs, tip)
	require.NoError(t, err)

	log.Info("Forced Batch Submit to L1 TxHash: ", tx.Hash())
	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	query := ethereum.FilterQuery{
		FromBlock: currentBlock.Number(),
		Addresses: []common.Address{zkEvmAddr},
	}
	logs, err := ethClient.FilterLogs(ctx, query)
	require.NoError(t, err)

	var forcedBatch *state.Batch
	for _, vLog := range logs {
		if vLog.Topics[0] != constants.ForcedBatchSignatureHash {
			logs, err = ethClient.FilterLogs(ctx, query)
			require.NoError(t, err)
			continue
		}
		fb, err := zkEvm.ParseForceBatch(vLog)
		if err != nil {
			log.Errorf("failed to parse force batch log event, err: ", err)
		}
		log.Debugf("log decoded: %+v", fb)
		ger := fb.LastGlobalExitRoot
		log.Info("GlobalExitRoot: ", ger)
		log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
		fullBlock, err := ethClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
			return nil, err
		}
		log.Info("MinForcedTimestamp: ", fullBlock.Time())
		forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
		for err == state.ErrStateNotSynchronized {
			time.Sleep(1 * time.Second)
			forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
		}
		require.NoError(t, err)
		require.NotNil(t, forcedBatch)

		log.Info("Waiting Forced Batch to be virtualized ...")
		err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)

		log.Info("Waiting Forced Batch to be consolidated ...")
		err = operations.WaitBatchToBeConsolidated(forcedBatch.BatchNumber, 4*time.Minute, st)
		require.NoError(t, err)
	}

	return forcedBatch, nil
}
