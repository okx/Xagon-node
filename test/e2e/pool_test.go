package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/HasOpCode"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/pendingtx/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/pendingtx/sequencer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestRepeatedNonce(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	receiverAddr := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	amount := big.NewInt(1000)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(fromPriKey, network.ChainID)

		nonce, err := client.NonceAt(ctx, auth.From, nil)
		require.NoError(t, err)

		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  auth.From,
			To:    &receiverAddr,
			Value: amount,
		})
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)

		correctNonceTx := types.NewTransaction(nonce, receiverAddr, amount.Add(amount, amount), gasLimit+gasLimit, gasPrice.Add(gasPrice, gasPrice), nil)
		correctNonceSignedTx, err := auth.Signer(auth.From, correctNonceTx)
		require.NoError(t, err)

		repeatedNonceTx := types.NewTransaction(nonce, receiverAddr, amount, gasLimit, gasPrice, nil)
		repeatedNonceSignedTx, err := auth.Signer(auth.From, repeatedNonceTx)
		require.NoError(t, err)

		log.Debug("sending correct nonce tx")
		err = client.SendTransaction(ctx, correctNonceSignedTx)
		require.NoError(t, err)

		log.Debug("sending repeated nonce tx")
		err = client.SendTransaction(ctx, repeatedNonceSignedTx)
		require.Equal(t, "replacement transaction underpriced", err.Error())

		log.Debug("waiting correct nonce tx to be mined")
		err = operations.WaitTxToBeMined(ctx, client, correctNonceSignedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}
}

func TestRepeatedTx(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	receiverAddr := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	amount := big.NewInt(1000)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(fromPriKey, network.ChainID)

		nonce, err := client.NonceAt(ctx, auth.From, nil)
		require.NoError(t, err)

		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  auth.From,
			To:    &receiverAddr,
			Value: amount,
		})
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)

		tx := types.NewTransaction(nonce, receiverAddr, amount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)

		log.Debug("sending tx")
		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		log.Debug("re sending tx")
		err = client.SendTransaction(ctx, signedTx)
		require.Equal(t, "already known", err.Error())

		log.Debug("waiting correct nonce tx to be mined")
		err = operations.WaitTxToBeMined(ctx, client, signedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}
}

func TestPendingNonce(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	receiverAddr := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	amount := big.NewInt(1000)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(fromPriKey, network.ChainID)

		nonce, err := client.NonceAt(ctx, auth.From, nil)
		require.NoError(t, err)
		log.Debug("nonce: ", nonce)

		pendingNonce, err := client.PendingNonceAt(ctx, auth.From)
		require.Equal(t, nonce, pendingNonce)
		require.NoError(t, err)
		log.Debug("pending Nonce: ", pendingNonce)

		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  auth.From,
			To:    &receiverAddr,
			Value: amount,
		})
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			txNonce := pendingNonce + uint64(i)
			log.Debugf("creating transaction with nonce %v: ", txNonce)
			tx := types.NewTransaction(txNonce, receiverAddr, amount, gasLimit, gasPrice, nil)
			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			log.Debug("sending tx")
			err = client.SendTransaction(ctx, signedTx)
			require.NoError(t, err)

			newPendingNonce, err := client.PendingNonceAt(ctx, auth.From)
			require.NoError(t, err)
			log.Debug("newPendingNonce: ", newPendingNonce)
			require.Equal(t, txNonce+1, newPendingNonce)
		}
	}
}

func TestHasOpCode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	client := operations.MustGetClient(operations.DefaultL2NetworkURL)
	auth := operations.MustGetAuth(fromPriKey, operations.DefaultL2ChainID)

	time.Sleep(2 * time.Second)

	log.Debug("deploying HasOpCode SC")
	_, scTx, sc, err := HasOpCode.DeployHasOpCode(auth, client)
	require.NoError(t, err)

	logTx(scTx)
	err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	log.Debug("calling opGasPrice")
	scCallOpGasPrice, err := sc.OpGasPrice(auth)
	require.NoError(t, err)

	logTx(scCallOpGasPrice)
	err = operations.WaitTxToBeMined(ctx, client, scCallOpGasPrice, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	log.Debug("get tx EGPLog")
	egpLog, err := opsMan.State().GetTransactionEGPLogByHash(ctx, scCallOpGasPrice.Hash(), nil)
	require.NoError(t, err)
	require.Equal(t, egpLog.GasPriceOC, true)

	log.Debug("calling opBalance")
	scCallBalance, err := sc.OpBalance(auth)
	require.NoError(t, err)

	logTx(scCallBalance)
	err = operations.WaitTxToBeMined(ctx, client, scCallBalance, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	log.Debug("get tx EGPLog")
	egpLog, err = opsMan.State().GetTransactionEGPLogByHash(ctx, scCallBalance.Hash(), nil)
	require.NoError(t, err)
	require.Equal(t, egpLog.BalanceOC, true)
}

func TestSpeedTx(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	//defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	//opsCfg := operations.GetDefaultOperationsConfig()
	//opsMan, err := operations.NewManager(ctx, opsCfg)
	//require.NoError(t, err)
	//err = opsMan.Setup()
	//require.NoError(t, err)

	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(dbutils.NewPoolConfigFromEnv())
	require.NoError(t, err)

	cfg, err := config.Default()
	require.NoError(t, err)
	cfg.Log.Level = "error"
	log.Init(cfg.Log)

	poolInstance := pool.NewPool(cfg.Pool, poolStorage)

	levelCount, txPerAddr, addrPerLevel := 50, 100, 5
	poolInstance.PrepareTx(ctx, levelCount, addrPerLevel, txPerAddr)

	time.Sleep(3 * time.Second)
	count, err := poolInstance.CountPendingTransactions(ctx)
	require.NoError(t, err)
	require.Equal(t, levelCount*txPerAddr*addrPerLevel, int(count))

	finishedCh := make(chan int, 1)
	cfg.Sequencer.QueryPendingTxsLimit = 10000
	seq, err := sequencer.New(cfg.Sequencer, cfg.State.Batch, poolInstance)
	require.NoError(t, err)
	go seq.Start(ctx, levelCount, addrPerLevel*txPerAddr, finishedCh)

	//time.Sleep(1 * time.Second)
	//go poolInstance.Speed(ctx)

	finishedLevel := levelCount - 1
	for {
		select {
		case level := <-finishedCh:
			require.Equal(t, finishedLevel, level)
			finishedLevel--
			if finishedLevel == -1 {
				return
			}
		}
	}
}
