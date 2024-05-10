package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Counter"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/CounterAndBlock"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/FailureTest"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Read"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		_, scTx, sc, err := Counter.DeployCounter(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		count, err := sc.GetCount(&bind.CallOpts{Pending: false})
		require.NoError(t, err)

		assert.Equal(t, 0, count.Cmp(big.NewInt(0)))

		scCallTx, err := sc.Increment(auth)
		require.NoError(t, err)

		logTx(scCallTx)
		err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		count, err = sc.GetCount(&bind.CallOpts{Pending: false})
		require.NoError(t, err)
		assert.Equal(t, 0, count.Cmp(big.NewInt(1)))
	}
}

func TestEmitLog2(t *testing.T) {
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

	type testCase struct {
		name                 string
		logsFromSubscription chan types.Log
		subscribe            func(*testing.T, *ethclient.Client, *testCase, common.Address) ethereum.Subscription
		getLogs              func(*testing.T, *ethclient.Client, *testCase, common.Address, *types.Receipt, ethereum.Subscription) []types.Log
		validate             func(*testing.T, context.Context, []types.Log, *EmitLog2.EmitLog2)
	}

	testCases := []testCase{
		{
			name: "validate logs by block number",
			getLogs: func(t *testing.T, client *ethclient.Client, tc *testCase, scAddr common.Address, scCallTxReceipt *types.Receipt, sub ethereum.Subscription) []types.Log {
				filterBlock := scCallTxReceipt.BlockNumber
				logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: filterBlock, ToBlock: filterBlock,
					Addresses: []common.Address{scAddr},
				})
				require.NoError(t, err)
				return logs
			},
			validate: func(t *testing.T, ctx context.Context, logs []types.Log, sc *EmitLog2.EmitLog2) {
				assert.Equal(t, 4, len(logs))

				log0 := getLogByIndex(0, logs)
				assert.Equal(t, 0, len(log0.Topics))

				_, err = sc.ParseLog(getLogByIndex(1, logs))
				require.NoError(t, err)

				logA, err := sc.ParseLogA(getLogByIndex(2, logs))
				require.NoError(t, err)
				expectedA := big.NewInt(1)
				assert.Equal(t, 0, logA.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logA.A.String())

				logABCD, err := sc.ParseLogABCD(getLogByIndex(3, logs))
				require.NoError(t, err)
				expectedA = big.NewInt(1)
				expectedB := big.NewInt(2)
				expectedC := big.NewInt(3)
				expectedD := big.NewInt(4)
				assert.Equal(t, 0, logABCD.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logABCD.A.String())
				assert.Equal(t, 0, logABCD.B.Cmp(expectedB), "B expected to be: %v found: %v", expectedA.String(), logABCD.B.String())
				assert.Equal(t, 0, logABCD.C.Cmp(expectedC), "C expected to be: %v found: %v", expectedA.String(), logABCD.C.String())
				assert.Equal(t, 0, logABCD.D.Cmp(expectedD), "D expected to be: %v found: %v", expectedA.String(), logABCD.D.String())
			},
		},
		{
			name: "validate logs by block number and topics",
			getLogs: func(t *testing.T, client *ethclient.Client, tc *testCase, scAddr common.Address, scCallTxReceipt *types.Receipt, sub ethereum.Subscription) []types.Log {
				filterBlock := scCallTxReceipt.BlockNumber
				logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: filterBlock, ToBlock: filterBlock,
					Addresses: []common.Address{scAddr},
					Topics: [][]common.Hash{
						{
							common.HexToHash("0xe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003"),
						},
					},
				})
				require.NoError(t, err)
				return logs
			},
			validate: func(t *testing.T, ctx context.Context, logs []types.Log, sc *EmitLog2.EmitLog2) {
				assert.Equal(t, 1, len(logs))

				logABCD, err := sc.ParseLogABCD(getLogByIndex(3, logs))
				require.NoError(t, err)
				expectedA := big.NewInt(1)
				expectedB := big.NewInt(2)
				expectedC := big.NewInt(3)
				expectedD := big.NewInt(4)
				assert.Equal(t, 0, logABCD.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logABCD.A.String())
				assert.Equal(t, 0, logABCD.B.Cmp(expectedB), "B expected to be: %v found: %v", expectedA.String(), logABCD.B.String())
				assert.Equal(t, 0, logABCD.C.Cmp(expectedC), "C expected to be: %v found: %v", expectedA.String(), logABCD.C.String())
				assert.Equal(t, 0, logABCD.D.Cmp(expectedD), "D expected to be: %v found: %v", expectedA.String(), logABCD.D.String())
			},
		},
		{
			name: "validate logs by block hash",
			getLogs: func(t *testing.T, client *ethclient.Client, tc *testCase, scAddr common.Address, scCallTxReceipt *types.Receipt, sub ethereum.Subscription) []types.Log {
				filterBlock := scCallTxReceipt.BlockHash
				logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
					BlockHash: &filterBlock,
					Addresses: []common.Address{scAddr},
				})
				require.NoError(t, err)
				return logs
			},
			validate: func(t *testing.T, ctx context.Context, logs []types.Log, sc *EmitLog2.EmitLog2) {
				assert.Equal(t, 4, len(logs))

				log0 := getLogByIndex(0, logs)
				assert.Equal(t, 0, len(log0.Topics))

				_, err = sc.ParseLog(getLogByIndex(1, logs))
				require.NoError(t, err)

				logA, err := sc.ParseLogA(getLogByIndex(2, logs))
				require.NoError(t, err)
				expectedA := big.NewInt(1)
				assert.Equal(t, 0, logA.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logA.A.String())

				logABCD, err := sc.ParseLogABCD(getLogByIndex(3, logs))
				require.NoError(t, err)
				expectedA = big.NewInt(1)
				expectedB := big.NewInt(2)
				expectedC := big.NewInt(3)
				expectedD := big.NewInt(4)
				assert.Equal(t, 0, logABCD.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logABCD.A.String())
				assert.Equal(t, 0, logABCD.B.Cmp(expectedB), "B expected to be: %v found: %v", expectedA.String(), logABCD.B.String())
				assert.Equal(t, 0, logABCD.C.Cmp(expectedC), "C expected to be: %v found: %v", expectedA.String(), logABCD.C.String())
				assert.Equal(t, 0, logABCD.D.Cmp(expectedD), "D expected to be: %v found: %v", expectedA.String(), logABCD.D.String())
			},
		},
		{
			name: "validate logs by block hash and topics",
			getLogs: func(t *testing.T, client *ethclient.Client, tc *testCase, scAddr common.Address, scCallTxReceipt *types.Receipt, sub ethereum.Subscription) []types.Log {
				filterBlock := scCallTxReceipt.BlockHash
				logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
					BlockHash: &filterBlock,
					Addresses: []common.Address{scAddr},
					Topics: [][]common.Hash{
						{
							common.HexToHash("0xe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003"),
						},
					},
				})
				require.NoError(t, err)
				return logs
			},
			validate: func(t *testing.T, ctx context.Context, logs []types.Log, sc *EmitLog2.EmitLog2) {
				assert.Equal(t, 1, len(logs))

				logABCD, err := sc.ParseLogABCD(getLogByIndex(3, logs))
				require.NoError(t, err)
				expectedA := big.NewInt(1)
				expectedB := big.NewInt(2)
				expectedC := big.NewInt(3)
				expectedD := big.NewInt(4)
				assert.Equal(t, 0, logABCD.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logABCD.A.String())
				assert.Equal(t, 0, logABCD.B.Cmp(expectedB), "B expected to be: %v found: %v", expectedA.String(), logABCD.B.String())
				assert.Equal(t, 0, logABCD.C.Cmp(expectedC), "C expected to be: %v found: %v", expectedA.String(), logABCD.C.String())
				assert.Equal(t, 0, logABCD.D.Cmp(expectedD), "D expected to be: %v found: %v", expectedA.String(), logABCD.D.String())
			},
		},
		{
			name: "validate logs by subscription",
			subscribe: func(t *testing.T, c *ethclient.Client, tc *testCase, scAddr common.Address) ethereum.Subscription {
				query := ethereum.FilterQuery{Addresses: []common.Address{scAddr}}
				sub, err := c.SubscribeFilterLogs(context.Background(), query, tc.logsFromSubscription)
				require.NoError(t, err)
				return sub
			},
			getLogs: func(t *testing.T, c *ethclient.Client, tc *testCase, a common.Address, r *types.Receipt, sub ethereum.Subscription) []types.Log {
				logs := []types.Log{}
				for {
					select {
					case err := <-sub.Err():
						require.NoError(t, err)
					case vLog, closed := <-tc.logsFromSubscription:
						logs = append(logs, vLog)
						if len(logs) == 4 && closed {
							return logs
						}
					}
				}
			},
			validate: func(t *testing.T, ctx context.Context, logs []types.Log, sc *EmitLog2.EmitLog2) {
				assert.Equal(t, 4, len(logs))

				log0 := getLogByIndex(0, logs)
				assert.Equal(t, 0, len(log0.Topics))

				logWithoutParameters, err := sc.ParseLog(getLogByIndex(1, logs))
				require.NoError(t, err)
				assert.Equal(t, 1, len(logWithoutParameters.Raw.Topics))

				logA, err := sc.ParseLogA(getLogByIndex(2, logs))
				require.NoError(t, err)
				expectedA := big.NewInt(1)
				assert.Equal(t, 0, logA.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logA.A.String())

				logABCD, err := sc.ParseLogABCD(getLogByIndex(3, logs))
				require.NoError(t, err)
				expectedA = big.NewInt(1)
				expectedB := big.NewInt(2)
				expectedC := big.NewInt(3)
				expectedD := big.NewInt(4)
				assert.Equal(t, 0, logABCD.A.Cmp(expectedA), "A expected to be: %v found: %v", expectedA.String(), logABCD.A.String())
				assert.Equal(t, 0, logABCD.B.Cmp(expectedB), "B expected to be: %v found: %v", expectedA.String(), logABCD.B.String())
				assert.Equal(t, 0, logABCD.C.Cmp(expectedC), "C expected to be: %v found: %v", expectedA.String(), logABCD.C.String())
				assert.Equal(t, 0, logABCD.D.Cmp(expectedD), "D expected to be: %v found: %v", expectedA.String(), logABCD.D.String())
			},
		},
	}

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		wsClient := operations.MustGetClient(network.WebSocketURL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		// deploy sc
		scAddr, scTx, sc, err := EmitLog2.DeployEmitLog2(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.logsFromSubscription = make(chan types.Log)

				var sub ethereum.Subscription
				if tc.subscribe != nil {
					sub = tc.subscribe(t, wsClient, &tc, scAddr)
				}

				// emit logs
				scCallTx, err := sc.EmitLogs(auth)
				require.NoError(t, err)

				logTx(scCallTx)
				err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
				require.NoError(t, err)

				scCallTxReceipt, err := client.TransactionReceipt(ctx, scCallTx.Hash())
				require.NoError(t, err)

				logs := tc.getLogs(t, client, &tc, scAddr, scCallTxReceipt, sub)

				tc.validate(t, ctx, logs, sc)
			})
		}
	}
}

func TestLogTxIndex(t *testing.T) {
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

	assertTxHashAndIndex := func(t *testing.T, log types.Log, tx *types.Transaction, receipt *types.Receipt) {
		assert.Equal(t, tx.Hash().String(), log.TxHash.String())
		assert.Equal(t, receipt.TxHash.String(), log.TxHash.String())
		assert.Equal(t, receipt.TransactionIndex, log.TxIndex)
	}

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		wsClient := operations.MustGetClient(network.WebSocketURL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		// deploy sc
		scAddr, scTx, sc, err := EmitLog2.DeployEmitLog2(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		if network.Name == "Local L2" {
			// stops sequencer
			err = operations.StopComponent("seq")
			require.NoError(t, err)
		}

		logsFromSubscription := make(chan types.Log)
		query := ethereum.FilterQuery{Addresses: []common.Address{scAddr}}
		sub, err := wsClient.SubscribeFilterLogs(context.Background(), query, logsFromSubscription)
		require.NoError(t, err)

		// send transfer
		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)
		nonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(t, err)
		tx := types.NewTx(&types.LegacyTx{
			To:       state.Ptr(common.HexToAddress("0x1")),
			Gas:      30000,
			GasPrice: gasPrice,
			Value:    big.NewInt(1000),
			Nonce:    nonce,
		})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		// send log tx
		auth.Nonce = big.NewInt(0).SetUint64(nonce + 1)
		scCallTx, err := sc.EmitLogs(auth)
		require.NoError(t, err)
		logTx(scCallTx)

		time.Sleep(time.Second)

		if network.Name == "Local L2" {
			// starts sequencer and wait log tx to get mined
			err = operations.StartComponent("seq", func() (done bool, err error) {
				err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
				return true, err
			})
			require.NoError(t, err)
		} else {
			err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
			require.NoError(t, err)
		}

		scCallTxReceipt, err := client.TransactionReceipt(ctx, scCallTx.Hash())
		require.NoError(t, err)

		if network.Name == "Local L2" {
			assert.Equal(t, uint(1), scCallTxReceipt.TransactionIndex)
		}

		// validate logs from filterLogs
		filterBlock := scCallTxReceipt.BlockNumber
		logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: filterBlock, ToBlock: filterBlock,
			Addresses: []common.Address{scAddr},
		})
		require.NoError(t, err)

		assert.Equal(t, 4, len(logs))
		for i := range logs {
			l := getLogByIndex(i, logs)
			assertTxHashAndIndex(t, l, scCallTx, scCallTxReceipt)
		}

		// validate logs from receipt
		logs = make([]types.Log, len(scCallTxReceipt.Logs))
		for i, log := range scCallTxReceipt.Logs {
			logs[i] = *log
		}

		assert.Equal(t, 4, len(logs))
		for i := range logs {
			l := getLogByIndex(i, logs)
			assertTxHashAndIndex(t, l, scCallTx, scCallTxReceipt)
		}

		// validate logs by subscription
		logs = []types.Log{}
	out:
		for {
			select {
			case err := <-sub.Err():
				require.NoError(t, err)
			case vLog, closed := <-logsFromSubscription:
				logs = append(logs, vLog)
				if len(logs) == 4 && closed {
					break out
				}
			}
		}

		assert.Equal(t, 4, len(logs))
		for i := range logs {
			l := getLogByIndex(i, logs)
			assertTxHashAndIndex(t, l, scCallTx, scCallTxReceipt)
		}
	}
}

func getLogByIndex(index int, logs []types.Log) types.Log {
	for _, log := range logs {
		if int(log.Index) == index {
			return log
		}
	}
	return types.Log{}
}

func TestFailureTest(t *testing.T) {
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		log.Debug("deploying SC")
		_, scTx, sc, err := FailureTest.DeployFailureTest(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("storing value")
		scCallTx, err := sc.Store(auth, big.NewInt(1))
		require.NoError(t, err)

		logTx(scCallTx)
		err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("storing value with revert")
		_, err = sc.StoreAndFail(auth, big.NewInt(2))
		assert.Equal(t, "execution reverted: this method always fails", err.Error())
	}
}

func TestRead(t *testing.T) {
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		const ownerName = "this is the owner name"
		callOpts := &bind.CallOpts{Pending: false}

		log.Debug("deploying SC")
		_, scTx, sc, err := Read.DeployRead(auth, client, ownerName)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("read string public variable directly")
		ownerNameValue, err := sc.OwnerName(callOpts)
		require.NoError(t, err)
		require.Equal(t, ownerName, ownerNameValue)

		log.Debug("read address public variable directly")
		ownerValue, err := sc.Owner(callOpts)
		require.NoError(t, err)
		require.Equal(t, auth.From, ownerValue)

		tA := Read.Readtoken{
			Name:     "Token A",
			Quantity: big.NewInt(50),
			Address:  common.HexToAddress("0x1"),
		}

		tB := Read.Readtoken{
			Name:     "Token B",
			Quantity: big.NewInt(30),
			Address:  common.HexToAddress("0x2"),
		}

		log.Debug("public add token")
		tx, err := sc.PublicAddToken(auth, tA)
		require.NoError(t, err)
		logTx(tx)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("external add token")
		tx, err = sc.ExternalAddToken(auth, tB)
		require.NoError(t, err)
		logTx(tx)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("read mapping public variable directly")
		tk, err := sc.Tokens(callOpts, tA.Address)
		require.NoError(t, err)
		require.Equal(t, tA.Name, tk.Name)
		require.Equal(t, tA.Quantity, tk.Quantity)
		require.Equal(t, tA.Address, tk.Address)

		tk, err = sc.Tokens(callOpts, tB.Address)
		require.NoError(t, err)
		require.Equal(t, tB.Name, tk.Name)
		require.Equal(t, tB.Quantity, tk.Quantity)
		require.Equal(t, tB.Address, tk.Address)

		log.Debug("public struct read")
		tk, err = sc.PublicGetToken(callOpts, tA.Address)
		require.NoError(t, err)
		require.Equal(t, tA.Name, tk.Name)
		require.Equal(t, tA.Quantity, tk.Quantity)
		require.Equal(t, tA.Address, tk.Address)

		log.Debug("external struct read")
		tk, err = sc.ExternalGetToken(callOpts, tB.Address)
		require.NoError(t, err)
		require.Equal(t, tB.Name, tk.Name)
		require.Equal(t, tB.Quantity, tk.Quantity)
		require.Equal(t, tB.Address, tk.Address)

		log.Debug("public uint256 read")
		value, err := sc.PublicRead(callOpts)
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(1).Cmp(value))

		log.Debug("external uint256 read")
		value, err = sc.ExternalRead(callOpts)
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(1).Cmp(value))

		log.Debug("public uint256 read with parameter")
		value, err = sc.PublicReadWParams(callOpts, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(2).Cmp(value))

		log.Debug("external uint256 read with parameter")
		value, err = sc.ExternalReadWParams(callOpts, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(2).Cmp(value))
	}
}

func TestCounterAndBlock(t *testing.T) {
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		_, scTx, sc, err := CounterAndBlock.DeployCounterAndBlock(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scReceipt, err := client.TransactionReceipt(ctx, scTx.Hash())
		require.NoError(t, err)

		scBlock, err := client.BlockByNumber(ctx, scReceipt.BlockNumber)
		require.NoError(t, err)

		count, ts, err := sc.GetCount(&bind.CallOpts{Pending: false, BlockNumber: scBlock.Number()})
		require.NoError(t, err)

		assert.Equal(t, 0, count.Cmp(big.NewInt(0)))
		assert.Equal(t, ts.Uint64(), scBlock.Time())

		const numberOfIncrements = 5
		type result struct {
			tx            *types.Transaction
			receipt       *types.Receipt
			block         *types.Block
			expectedCount *big.Int
		}

		results := make([]result, 0, numberOfIncrements)
		for i := 0; i < numberOfIncrements; i++ {
			tx, err := sc.Increment(auth)
			require.NoError(t, err)

			logTx(tx)
			err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
			require.NoError(t, err)

			receipt, err := client.TransactionReceipt(ctx, tx.Hash())
			require.NoError(t, err)

			block, err := client.BlockByNumber(ctx, receipt.BlockNumber)
			require.NoError(t, err)

			results = append(results, result{
				tx:            tx,
				expectedCount: big.NewInt(int64(i) + 1),
				receipt:       receipt,
				block:         block,
			})
		}

		const numberOfChecks = 2

		// checks against first increment
		for _, r := range results {
			for i := 0; i < numberOfChecks; i++ {
				count, ts, err = sc.GetCount(&bind.CallOpts{Pending: false, BlockNumber: r.block.Number()})
				require.NoError(t, err)
				assert.Equal(t, r.expectedCount.Uint64(), count.Uint64())
				assert.Equal(t, r.block.Time(), ts.Uint64())

				time.Sleep(time.Second)
			}
		}

		latestIncrement := results[len(results)-1]
		// checks against second increment with latest block
		for i := 0; i < numberOfChecks; i++ {
			latestBlock, err := client.BlockByNumber(ctx, nil)
			require.NoError(t, err)

			count, ts, err = sc.GetCount(&bind.CallOpts{Pending: false})
			require.NoError(t, err)
			assert.Equal(t, latestIncrement.expectedCount.Uint64(), count.Uint64())
			assert.Equal(t, latestBlock.Time(), ts.Uint64())

			time.Sleep(time.Second)
		}
	}
}
