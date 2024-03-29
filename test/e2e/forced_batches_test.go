package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/pol"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	// dockersArePreLaunched is a flag that indicates if dockers are pre-launched, used for local development
	// avoiding launch time and reset database time at end (so you can check the database after the test)
	dockersArePreLaunched = false
	gerFinalityBlocks     = uint64(9223372036854775807) // The biggeset uint64
)

type l1Stuff struct {
	ethClient       *ethclient.Client
	authSequencer   *bind.TransactOpts
	authForcedBatch *bind.TransactOpts
	zkEvmAddr       common.Address
	zkEvm           *polygonzkevm.Polygonzkevm
}

type l2Stuff struct {
	opsman        *operations.Manager
	authSequencer *bind.TransactOpts
	client        *ethclient.Client
	amount        *big.Int
	gasLimit      uint64
	gasPrice      *big.Int
	nonce         uint64
}

//TODO: Fix test ETROG
/*func TestForcedBatches(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	log.Infof("Running TestForcedBatches ==========================")
	if !dockersArePreLaunched {
		defer func() {
			require.NoError(t, operations.Teardown())
		}()
	}

	var err error
	nTxs := 10
	ctx := context.Background()
	l2 := setupEnvironment(ctx, t)
	l1 := setupEnvironmentL1(ctx, t)
	l2BlockNumbersTxsBeforeForcedBatch := generateTxsBeforeSendingForcedBatch(ctx, t, nTxs, l2)
	time.Sleep(2 * time.Second)
	l2.amount = big.NewInt(0).Add(l2.amount, big.NewInt(10))
	encodedTxs := generateSignedAndEncodedTxForForcedBatch(ctx, t, l2)
	forcedBatch, err := sendForcedBatch(ctx, t, encodedTxs, l2.opsman, l1)
	require.NoError(t, err)
	checkThatPreviousTxsWereProcessedWithinPreviousClosedBatch(ctx, t, l2.opsman.State(), l2BlockNumbersTxsBeforeForcedBatch, forcedBatch.BatchNumber)
}*/

func generateTxsBeforeSendingForcedBatch(ctx context.Context, t *testing.T, nTxs int, l2 *l2Stuff) []*big.Int {
	txs := make([]*types.Transaction, 0, nTxs)
	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(l2.nonce, toAddress, l2.amount, l2.gasLimit, l2.gasPrice, nil)
		l2.nonce = l2.nonce + 1
		txs = append(txs, tx)
	}

	var l2BlockNumbers []*big.Int
	l2BlockNumbers, err := operations.ApplyL2Txs(ctx, txs, l2.authSequencer, l2.client, operations.VerifiedConfirmationLevel)
	require.NoError(t, err)
	return l2BlockNumbers
}

func checkThatPreviousTxsWereProcessedWithinPreviousClosedBatch(ctx context.Context, t *testing.T, state *state.State, l2BlockNumbers []*big.Int, forcedBatchNumber uint64) {
	// Checking if all txs sent before the forced batch were processed within previous closed batch
	for _, l2blockNum := range l2BlockNumbers {
		batch, err := state.GetBatchByL2BlockNumber(ctx, l2blockNum.Uint64(), nil)
		require.NoError(t, err)
		require.Less(t, batch.BatchNumber, forcedBatchNumber)
	}
}

func generateSignedAndEncodedTxForForcedBatch(ctx context.Context, t *testing.T, l2 *l2Stuff) []byte {
	unsignedTx := types.NewTransaction(l2.nonce, toAddress, l2.amount, l2.gasLimit, l2.gasPrice, nil)
	signedTx, err := l2.authSequencer.Signer(l2.authSequencer.From, unsignedTx)
	require.NoError(t, err)
	log.Info("Forced Batch: 1 tx -> ", signedTx.Hash())
	encodedTxs, err := state.EncodeTransactions([]types.Transaction{*signedTx}, constants.EffectivePercentage, forkID6)
	require.NoError(t, err)
	return encodedTxs
}

func setupEnvironment(ctx context.Context, t *testing.T) *l2Stuff {
	if !dockersArePreLaunched {
		err := operations.Teardown()
		require.NoError(t, err)
	}
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000

	var opsman *operations.Manager
	var err error

	if !dockersArePreLaunched {
		log.Info("Launching dockers and resetting Database")
		opsman, err = operations.NewManager(ctx, opsCfg)
		require.NoError(t, err)
		log.Info("Setting Genesis")
		setInitialState(t, opsman)
	} else {
		log.Info("Using pre-launched dockers: no reset Database")
		opsman, err = operations.NewManagerNoInitDB(ctx, opsCfg)
		require.NoError(t, err)
	}

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)
	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(t, err)
	// Send txs
	amount := big.NewInt(10000)
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(t, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	log.Infof("Receiver Addr: %v", toAddress.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)
	return &l2Stuff{opsman, auth, client, amount, gasLimit, gasPrice, senderNonce}
}

func setupEnvironmentL1(ctx context.Context, t *testing.T) *l1Stuff {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)
	authSequencer, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)
	authForcedBatch, err := operations.GetAuth(operations.DefaultForcedBatchesPrivateKey, operations.DefaultL1ChainID)
	require.NoError(t, err)
	polSmc, err := pol.NewPol(common.HexToAddress(operations.DefaultL1PolSmartContract), ethClient)
	require.NoError(t, err)
	polAmount, _ := big.NewInt(0).SetString("9999999999999999999999", 0)
	log.Debugf("Charging pol from sequencer -> forcedBatchesAddress")
	txValue, err := polSmc.Transfer(authSequencer, common.HexToAddress(operations.DefaultForcedBatchesAddress), polAmount)
	require.NoError(t, err)
	log.Debugf("Waiting for tx %s to be mined (transfer of pol from sequencer -> forcedBatches)", txValue.Hash().String())
	err = operations.WaitTxToBeMined(ctx, ethClient, txValue, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)
	balance, err := polSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultSequencerAddress))
	require.NoError(t, err)
	log.Debugf("Account (sequencer) %s pol balance %s", operations.DefaultSequencerAddress, balance.String())

	balance, err = polSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultForcedBatchesAddress))
	require.NoError(t, err)
	log.Debugf("Account (force_batches) %s pol balance %s", operations.DefaultForcedBatchesAddress, balance.String())
	log.Debugf("Approve to zkEVM SMC to spend %s pol", polAmount.String())
	_, err = polSmc.Approve(authForcedBatch, common.HexToAddress(operations.DefaultL1ZkEVMSmartContract), polAmount)
	require.NoError(t, err)

	zkEvmAddr := common.HexToAddress(operations.DefaultL1ZkEVMSmartContract)
	zkEvm, err := polygonzkevm.NewPolygonzkevm(zkEvmAddr, ethClient)
	require.NoError(t, err)
	return &l1Stuff{ethClient: ethClient, authSequencer: authSequencer, authForcedBatch: authForcedBatch, zkEvmAddr: zkEvmAddr, zkEvm: zkEvm}
}

func setInitialState(t *testing.T, opsman *operations.Manager) {
	genesisFileAsStr, err := config.LoadGenesisFileAsString("../../test/config/test.genesis.config.json")
	require.NoError(t, err)
	genesisConfig, err := config.LoadGenesisFromJSONString(genesisFileAsStr)
	require.NoError(t, err)
	require.NoError(t, opsman.SetForkID(genesisConfig.Genesis.RollupBlockNumber, forkID6))
	err = opsman.Setup()
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
}

func sendForcedBatch(ctx context.Context, t *testing.T, txs []byte, opsman *operations.Manager, l1 *l1Stuff) (*state.Batch, error) {
	st := opsman.State()

	initialGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
	require.NoError(t, err)

	log.Info("Using address: ", l1.authForcedBatch.From)

	num, err := l1.zkEvm.LastForceBatch(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	log.Info("Number of forceBatches in the smc: ", num)

	rollupManagerAddr := common.HexToAddress(operations.DefaultL1RollupManagerSmartContract)
	rollupManager, err := polygonrollupmanager.NewPolygonrollupmanager(rollupManagerAddr, l1.ethClient)
	require.NoError(t, err)

	// Get tip
	tip, err := rollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	log.Infof("Foced Batch Fee:%s", tip.String())
	managerAddress, err := l1.zkEvm.GlobalExitRootManager(&bind.CallOpts{Pending: false})
	require.NoError(t, err)

	manager, err := polygonzkevmglobalexitroot.NewPolygonzkevmglobalexitroot(managerAddress, l1.ethClient)
	require.NoError(t, err)

	rootInContract, err := manager.GetLastGlobalExitRoot(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	rootInContractHash := common.BytesToHash(rootInContract[:])

	log.Infof("Activating forced batches...")
	tx, err := l1.zkEvm.SetForceBatchAddress(l1.authSequencer, common.Address{})
	require.NoError(t, err)
	log.Infof("Forced batch is disallowed. Activated. Waiting for tx %s to be mined", tx.Hash())
	err = operations.WaitTxToBeMined(ctx, l1.ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	currentBlock, err := l1.ethClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	log.Debugf("L1: currentBlock: number:%s Time():%s ", currentBlock.Number().String(), currentBlock.Time())

	// Send forceBatch
	tx, err = l1.zkEvm.ForceBatch(l1.authForcedBatch, txs, tip)
	require.NoError(t, err)

	log.Info("TxHash: ", tx.Hash())
	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, l1.ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	fb, vLog, err := findForcedBatchInL1Logs(ctx, t, currentBlock.Number(), l1)
	if err != nil {
		log.Errorf("failed to parse force batch log event, err: ", err)
	}
	ger := fb.LastGlobalExitRoot

	log.Debugf("log decoded: %+v", fb)
	log.Info("GlobalExitRoot: ", ger)
	log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
	log.Info("ForcedBatchNum: ", fb.ForceBatchNum)
	fullBlock, err := l1.ethClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
		return nil, err
	}
	log.Info("MinForcedTimestamp: ", fullBlock.Time())
	forcedBatch, err := st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
	for err == state.ErrStateNotSynchronized {
		log.Infof("state not synced, waiting...")
		time.Sleep(1 * time.Second)
		forcedBatch, err = st.GetBatchByForcedBatchNum(ctx, fb.ForceBatchNum, nil)
	}

	require.NoError(t, err)
	require.NotNil(t, forcedBatch)

	log.Info("Waiting for batch to be virtualized...")
	err = operations.WaitBatchToBeVirtualized(forcedBatch.BatchNumber, 4*time.Minute, st)
	require.NoError(t, err)

	log.Info("Waiting for batch to be consolidated...")
	err = operations.WaitBatchToBeConsolidated(forcedBatch.BatchNumber, 4*time.Minute, st)
	require.NoError(t, err)

	if rootInContractHash != initialGer.GlobalExitRoot {
		log.Info("Checking if global exit root is updated...")
		finalGer, _, err := st.GetLatestGer(ctx, gerFinalityBlocks)
		require.NoError(t, err)
		require.Equal(t, rootInContractHash, finalGer.GlobalExitRoot, "global exit root is not updated")
	}

	return forcedBatch, nil
}

func findForcedBatchInL1Logs(ctx context.Context, t *testing.T, fromBlock *big.Int, l1 *l1Stuff) (*polygonzkevm.PolygonzkevmForceBatch, *types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: []common.Address{l1.zkEvmAddr},
	}

	found := false
	for found != true {
		log.Debugf("Looking for forced batch in logs from block %s", fromBlock.String())
		logs, err := l1.ethClient.FilterLogs(ctx, query)
		require.NoError(t, err)
		for _, vLog := range logs {
			if vLog.Topics[0] == constants.ForcedBatchSignatureHash {
				fb, err := l1.zkEvm.ParseForceBatch(vLog)
				return fb, &vLog, err
			}
		}
		log.Info("Forced batch not found in logs. Waiting 1 second...")
		time.Sleep(1 * time.Second)
	}
	return nil, nil, nil

}
