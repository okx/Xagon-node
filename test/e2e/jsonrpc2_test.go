package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Counter"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Revert"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Revert2"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Storage"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Misc(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	setup()
	defer teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		ethereumClient, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		// ChainId()
		chainId, err := ethereumClient.ChainID(ctx)
		require.NoError(t, err)
		require.Equal(t, network.ChainID, chainId.Uint64())

		// Syncing()
		progress, err := ethereumClient.SyncProgress(ctx)
		require.NoError(t, err)
		if progress != nil {
			log.Info("Its syncing")
			blockNumber, err := ethereumClient.BlockNumber(ctx)
			require.NoError(t, err)
			// if it's actually syncing
			require.Equal(t, uint64(0x0), progress.StartingBlock)
			require.Equal(t, blockNumber, progress.CurrentBlock)
		}

		// GetStorageAt()

		// first deploy sample smart contract
		sc_payload := int64(42)
		sc_retrieve := common.HexToHash("0x2a")
		auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, network.ChainID)
		require.NoError(t, err)
		contractAddress, tx, storageSC, err := Storage.DeployStorage(auth, ethereumClient)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
		tx, err = storageSC.Store(auth, big.NewInt(sc_payload))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		storage, err := ethereumClient.StorageAt(ctx, contractAddress, common.HexToHash("0x0"), nil)
		require.NoError(t, err)
		// in bytes but has to be hash 0x0...42
		require.Equal(t, sc_retrieve, common.BytesToHash(storage))

		// eth_getCode

		scBytecode, err := ethereumClient.CodeAt(ctx, contractAddress, nil)
		require.NoError(t, err)
		require.Contains(t, Storage.StorageMetaData.Bin, common.Bytes2Hex(scBytecode))

		emptyBytecode, err := ethereumClient.CodeAt(ctx, common.HexToAddress("0xdeadbeef"), nil)
		require.NoError(t, err)
		require.Empty(t, emptyBytecode)

		// check for request having more params than required:

		response, err := client.JSONRPCCall(network.URL, "eth_chainId", common.HexToHash("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"), "latest")
		require.NoError(t, err)
		require.NotNil(t, response.Error)
		require.Nil(t, response.Result)
		require.Equal(t, invalidParamsErrorCode, response.Error.Code)
		require.Equal(t, "too many arguments, want at most 0", response.Error.Message)
	}
}

func Test_WebSocketsRequest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	setup()
	defer teardown()

	acc := common.HexToAddress(operations.DefaultSequencerAddress)

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		client, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		expectedBalance, err := client.BalanceAt(ctx, acc, nil)
		require.NoError(t, err)

		wsConn, _, err := websocket.DefaultDialer.Dial(network.WebSocketURL, nil)
		require.NoError(t, err)

		receivedMessages := make(chan []byte)
		go func() {
			for {
				_, message, err := wsConn.ReadMessage()
				require.NoError(t, err)
				receivedMessages <- message
				wsConn.Close()
				break
			}
		}()

		params := []string{acc.String(), "latest"}
		jParam, err := json.Marshal(params)
		require.NoError(t, err)

		req := types.Request{JSONRPC: "2.0", ID: float64(1), Method: "eth_getBalance", Params: jParam}
		jReq, _ := json.Marshal(req)

		err = wsConn.WriteMessage(websocket.TextMessage, jReq)
		require.NoError(t, err)

		receivedMessage := <-receivedMessages

		resp := types.Response{}
		err = json.Unmarshal(receivedMessage, &resp)
		require.NoError(t, err)

		assert.Equal(t, req.JSONRPC, resp.JSONRPC)
		assert.Equal(t, req.ID, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		result := ""
		err = json.Unmarshal(resp.Result, &result)
		require.NoError(t, err)

		str := strings.TrimPrefix(result, "0x")
		balance := hex.DecodeBig(str)
		require.NoError(t, err)

		assert.Equal(t, expectedBalance.String(), balance.String())
	}
}

func Test_WebSocketsSubscription(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		wsConn, _, err := websocket.DefaultDialer.Dial(network.WebSocketURL, nil)
		require.NoError(t, err)

		receivedMessages := make(chan []byte)
		go func() {
			for {
				_, message, err := wsConn.ReadMessage()
				require.NoError(t, err)
				receivedMessages <- message
				break
			}
		}()

		params := []string{"newHeads"}
		jParam, err := json.Marshal(params)
		require.NoError(t, err)

		req := types.Request{JSONRPC: "2.0", ID: float64(1), Method: "eth_subscribe", Params: jParam}
		jReq, _ := json.Marshal(req)

		err = wsConn.WriteMessage(websocket.TextMessage, jReq)
		require.NoError(t, err)

		subscriptionMessage := <-receivedMessages

		resp := types.Response{}
		err = json.Unmarshal(subscriptionMessage, &resp)
		require.NoError(t, err)

		assert.Equal(t, req.JSONRPC, resp.JSONRPC)
		assert.Equal(t, req.ID, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		subscription := ""
		err = json.Unmarshal(resp.Result, &subscription)
		require.NoError(t, err)

		assert.NotEmpty(t, subscription)

		const numberOfBlocks = 3

		go func() {
			for i := 0; i <= numberOfBlocks; i++ {
				_, message, err := wsConn.ReadMessage()
				require.NoError(t, err)
				receivedMessages <- message
			}
			wsConn.Close()
		}()

		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)
		for i := 0; i <= numberOfBlocks; i++ {
			tx, err := createTX(client, auth, toAddress, big.NewInt(1000000000))
			require.NoError(t, err)
			err = operations.WaitTxToBeMined(context.Background(), client, tx, operations.DefaultTimeoutTxToBeMined)
			require.NoError(t, err)
		}

		for i := 0; i <= numberOfBlocks; i++ {
			receivedMessage := <-receivedMessages
			resp := types.SubscriptionResponse{}

			err = json.Unmarshal(receivedMessage, &resp)
			require.NoError(t, err)

			assert.Equal(t, req.JSONRPC, resp.JSONRPC)
			assert.Equal(t, "eth_subscription", resp.Method)
			assert.Equal(t, subscription, resp.Params.Subscription)

			block := map[string]interface{}{}
			err = json.Unmarshal(resp.Params.Result, &block)
			require.NoError(t, err)
			assert.NotEmpty(t, block["hash"].(string))
		}
	}
}

func Test_RevertOnConstructorTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	ctx := context.Background()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		auth.GasLimit = 1000000

		_, scTx, _, err := Revert.DeployRevert(auth, client)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		errMsg := err.Error()
		log.Debugf("Error: %v", errMsg)
		prefix := "transaction has failed, reason: execution reverted: Today is not juernes"
		hasPrefix := strings.HasPrefix(errMsg, prefix)
		require.True(t, hasPrefix)

		receipt, err := client.TransactionReceipt(ctx, scTx.Hash())
		require.NoError(t, err)

		assert.Equal(t, receipt.Status, ethTypes.ReceiptStatusFailed)

		msg := ethereum.CallMsg{
			From: auth.From,
			To:   scTx.To(),
			Gas:  scTx.Gas(),

			Value: scTx.Value(),
			Data:  scTx.Data(),
		}
		result, err := client.CallContract(ctx, msg, receipt.BlockNumber)
		require.NotNil(t, err)
		require.Nil(t, result)
		rpcErr := err.(rpc.Error)
		assert.Equal(t, 3, rpcErr.ErrorCode())
		assert.Equal(t, "execution reverted: Today is not juernes", rpcErr.Error())

		dataErr := err.(rpc.DataError)
		data := dataErr.ErrorData().(string)
		decodedData := hex.DecodeBig(data)
		unpackedData, err := abi.UnpackRevert(decodedData.Bytes())
		require.NoError(t, err)
		assert.Equal(t, "Today is not juernes", unpackedData)
	}
}

func Test_RevertOnSCCallTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	ctx := context.Background()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		auth.GasLimit = 1000000

		_, scTx, sc, err := Revert2.DeployRevert2(auth, client)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		tx, err := sc.GenerateError(auth)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		errMsg := err.Error()
		prefix := "transaction has failed, reason: execution reverted: Today is not juernes"
		hasPrefix := strings.HasPrefix(errMsg, prefix)
		require.True(t, hasPrefix)

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		assert.Equal(t, receipt.Status, ethTypes.ReceiptStatusFailed)

		msg := ethereum.CallMsg{
			From: auth.From,
			To:   tx.To(),
			Gas:  tx.Gas(),

			Value: tx.Value(),
			Data:  tx.Data(),
		}
		result, err := client.CallContract(ctx, msg, receipt.BlockNumber)
		require.NotNil(t, err)
		require.Nil(t, result)
		rpcErr := err.(rpc.Error)
		assert.Equal(t, 3, rpcErr.ErrorCode())
		assert.Equal(t, "execution reverted: Today is not juernes", rpcErr.Error())

		dataErr := err.(rpc.DataError)
		data := dataErr.ErrorData().(string)
		decodedData := hex.DecodeBig(data)
		unpackedData, err := abi.UnpackRevert(decodedData.Bytes())
		require.NoError(t, err)
		assert.Equal(t, "Today is not juernes", unpackedData)
	}
}

func Test_RevertOnSCCallGasEstimation(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	ctx := context.Background()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		auth.GasLimit = 1000000

		_, scTx, sc, err := Revert2.DeployRevert2(auth, client)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		tx, err := sc.GenerateError(auth)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		errMsg := err.Error()
		prefix := "transaction has failed, reason: execution reverted: Today is not juernes"
		hasPrefix := strings.HasPrefix(errMsg, prefix)
		require.True(t, hasPrefix)

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		assert.Equal(t, receipt.Status, ethTypes.ReceiptStatusFailed)

		msg := ethereum.CallMsg{
			From: auth.From,
			To:   tx.To(),
			Gas:  tx.Gas(),

			Value: tx.Value(),
			Data:  tx.Data(),
		}
		result, err := client.EstimateGas(ctx, msg)
		require.NotNil(t, err)
		require.Equal(t, uint64(0), result)
		rpcErr := err.(rpc.Error)
		assert.Equal(t, 3, rpcErr.ErrorCode())
		assert.Equal(t, "execution reverted: Today is not juernes", rpcErr.Error())

		dataErr := err.(rpc.DataError)
		data := dataErr.ErrorData().(string)
		decodedData := hex.DecodeBig(data)
		unpackedData, err := abi.UnpackRevert(decodedData.Bytes())
		require.NoError(t, err)
		assert.Equal(t, "Today is not juernes", unpackedData)
	}
}

func TestCallMissingParameters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	type testCase struct {
		name          string
		params        []interface{}
		expectedError types.ErrorObject
	}

	testCases := []testCase{
		{
			name:          "params is empty",
			params:        []interface{}{},
			expectedError: types.ErrorObject{Code: types.InvalidParamsErrorCode, Message: "missing value for required argument 0"},
		},
		{
			name:   "params has only first parameter",
			params: []interface{}{map[string]interface{}{"value": "0x1", "from": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "to": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92267"}},
		},
	}

	for _, network := range networks {
		t.Logf("Network %s", network.Name)
		for tc, testCase := range testCases {
			t.Logf("testCase %d", tc)
			t.Run(network.Name+testCase.name, func(t *testing.T) {
				response, err := client.JSONRPCCall(network.URL, "eth_call", testCase.params...)
				require.NoError(t, err)
				if (testCase.expectedError != types.ErrorObject{}) {
					require.NotNil(t, response.Error)
					require.Nil(t, response.Result)
					require.Equal(t, testCase.expectedError.Code, response.Error.Code)
					require.Equal(t, testCase.expectedError.Message, response.Error.Message)
				} else {
					require.Nil(t, response.Error)
					require.NotNil(t, response.Result)
				}
			})
		}
	}
}

func TestWebSocketsConcurrentWrites(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	const msgQty = 1000

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		wsConn, _, err := websocket.DefaultDialer.Dial(network.WebSocketURL, nil)
		require.NoError(t, err)
		defer func() {
			err := wsConn.Close()
			require.NoError(t, err)
		}()

		wg := sync.WaitGroup{}
		wg.Add(msgQty)
		receivedMessages := make(chan []byte, msgQty)
		go func(t *testing.T, wsConn *websocket.Conn) {
			c := 0
			for i := 0; i < msgQty; i++ {
				_, message, err := wsConn.ReadMessage()
				require.NoError(t, err)
				receivedMessages <- message
				wg.Done()
				c++
			}
		}(t, wsConn)

		params := []string{}
		jParam, err := json.Marshal(params)
		require.NoError(t, err)

		req := types.Request{JSONRPC: "2.0", ID: float64(1), Method: "eth_blockNumber", Params: jParam}
		jReq, _ := json.Marshal(req)

		for i := 0; i < msgQty; i++ {
			err = wsConn.WriteMessage(websocket.TextMessage, jReq)
			require.NoError(t, err)
		}

		if waitTimeout(&wg, time.Second*5) {
			t.Error("timeout reached")
			return
		}

		assert.Equal(t, msgQty, len(receivedMessages))
	}
}

func TestWebSocketsReadLimit(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	wsConn, _, err := websocket.DefaultDialer.Dial(operations.DefaultL2NetworkWebSocketURL, nil)
	require.NoError(t, err)
	defer func() {
		err := wsConn.Close()
		require.NoError(t, err)
	}()

	jReq := make([]byte, 104857601)
	err = wsConn.WriteMessage(websocket.TextMessage, jReq)
	require.NoError(t, err)

	_, _, err = wsConn.ReadMessage()
	require.NotNil(t, err)
	require.Equal(t, websocket.CloseMessageTooBig, err.(*websocket.CloseError).Code)
}

func TestEstimateTxWithDataBiggerThanMaxAllowed(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	ctx := context.Background()

	ethereumClient, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(t, err)

	sender := common.HexToAddress(operations.DefaultSequencerAddress)
	receiver := common.HexToAddress(operations.DefaultSequencerAddress)

	balance, err := ethereumClient.BalanceAt(ctx, sender, nil)
	require.NoError(t, err)

	_, err = ethereumClient.EstimateGas(ctx, ethereum.CallMsg{
		From:     sender,
		To:       &receiver,
		Value:    new(big.Int),
		Gas:      balance.Uint64(),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     make([]byte, 120000), // large data
	})
	rpcErr := err.(rpc.Error)
	assert.Equal(t, -32000, rpcErr.ErrorCode())
	assert.Equal(t, "batch_l2_data is invalid", rpcErr.Error())
}

func TestEstimateGas(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	log.Info("TestEstimateGas-------------------------------------------------------")
	setup()
	defer opsMan.ShowDockerLogs()
	defer teardown()

	ctx := context.Background()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		ethereumClient, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		// deploy a smart contract
		_, tx, sc, err := Counter.DeployCounter(auth, ethereumClient)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		gasPrice, err := ethereumClient.SuggestGasPrice(ctx)
		require.NoError(t, err)

		// prepare a tx information to be estimated
		auth.NoSend = true // force the tx to not be sent while using the sc method
		auth.GasLimit = 1  // force gas limit to avoid estimation while building the tx
		txToMsg, err := sc.Increment(auth)
		require.NoError(t, err)

		// addresses the test needs to have balance
		addressesToAddBalance := map[common.Address]*big.Int{
			// add funds to address 0x111...111 which is the default address
			// when estimating TXs without specifying the sender
			common.HexToAddress("0x1111111111111111111111111111111111111111"): big.NewInt(3000000000000000),

			// add funds to address 0x000...001
			common.HexToAddress("0x1"): big.NewInt(1000),
		}

		for addr, value := range addressesToAddBalance {
			nonce, err := ethereumClient.NonceAt(ctx, auth.From, nil)
			require.NoError(t, err)
			value := value
			require.NoError(t, err)
			tx = ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce:    nonce,
				To:       state.Ptr(addr),
				Value:    value,
				Gas:      24000,
				GasPrice: gasPrice,
			})
			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)
			err = ethereumClient.SendTransaction(ctx, signedTx)
			require.NoError(t, err)
			err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
			require.NoError(t, err)

			balance, err := ethereumClient.BalanceAt(ctx, addr, nil)
			require.NoError(t, err)
			log.Debugf("%v balance: %v", addr.String(), balance.String())
		}

		type testCase struct {
			name          string
			address       *common.Address
			value         *int64
			setGasPrice   bool
			expectedError rpc.Error
		}

		testCases := []testCase{
			{
				name:          "with gasPrice set and address with enough balance",
				address:       state.Ptr(auth.From),
				setGasPrice:   true,
				expectedError: nil,
			},
			{
				name:          "with gasPrice set and address without enough balance",
				address:       state.Ptr(common.HexToAddress("0x1")),
				setGasPrice:   true,
				expectedError: types.NewRPCError(-32000, "gas required exceeds allowance"),
			},
			{
				name:          "with gasPrice set and address with balance zero",
				address:       state.Ptr(common.HexToAddress("0x2")),
				setGasPrice:   true,
				expectedError: types.NewRPCError(-32000, "gas required exceeds allowance"),
			},
			{
				name:          "with gasPrice set and without from address",
				address:       nil,
				setGasPrice:   true,
				expectedError: nil,
			},
			{
				name:          "with gasPrice and value set and address with enough balance",
				address:       state.Ptr(auth.From),
				value:         state.Ptr(int64(1)),
				setGasPrice:   true,
				expectedError: types.NewRPCError(-32000, "execution reverted"),
			},
			{
				name:          "with gasPrice and value set and address without enough balance",
				address:       state.Ptr(common.HexToAddress("0x1")),
				value:         state.Ptr(int64(-1)),
				setGasPrice:   true,
				expectedError: types.NewRPCError(-32000, "insufficient funds for transfer"),
			},
			{
				name:          "with gasPrice and value set and address with balance zero",
				address:       state.Ptr(common.HexToAddress("0x2")),
				value:         state.Ptr(int64(-1)),
				setGasPrice:   true,
				expectedError: types.NewRPCError(-32000, "insufficient funds for transfer"),
			},
			// TODO = Review the test below in future versions of geth.
			//
			// Geth is returning -32000, "insufficient funds for transfer"
			// zkEVM is returning 3, "execution reverted"
			//
			// Since the tx has value, the method increment is not payable
			// and the default account has balance, the tx should revert
			//
			// {
			// 	name:          "with gasPrice and value set and without from address",
			// 	address:       nil,
			// 	value:         state.Ptr(int64(-1)),
			// 	setGasPrice:   true,
			// 	expectedError: types.NewRPCError(-32000, "insufficient funds for transfer"),
			// },
			{
				name:          "without gasPrice set and address with enough balance",
				address:       state.Ptr(auth.From),
				setGasPrice:   false,
				expectedError: nil,
			},
			{
				name:          "without gasPrice set and address without enough balance",
				address:       state.Ptr(common.HexToAddress("0x1")),
				setGasPrice:   false,
				expectedError: nil,
			},
			{
				name:          "without gasPrice set and address with balance zero",
				address:       state.Ptr(common.HexToAddress("0x2")),
				setGasPrice:   false,
				expectedError: nil,
			},
			{
				name:          "without gasPrice set and without from address",
				address:       nil,
				setGasPrice:   false,
				expectedError: nil,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				log.Info("TestEstimateGas------------------------------------------------------- subcase: ", testCase.name)
				msg := ethereum.CallMsg{
					To:   txToMsg.To(),
					Data: txToMsg.Data(),
				}
				if testCase.address != nil {
					msg.From = *testCase.address
				}
				balance, err := ethereumClient.BalanceAt(ctx, msg.From, nil)
				require.NoError(t, err)

				if testCase.value != nil {
					v := *testCase.value
					if v == -1 { //set the value as acc balance + 1 to force overflow
						msg.Value = common.Big0.Add(balance, common.Big1)
					} else {
						msg.Value = big.NewInt(0).SetInt64(v)
					}
				}

				if testCase.setGasPrice {
					msg.GasPrice = gasPrice
				}

				gas, err := ethereumClient.EstimateGas(ctx, msg)
				t.Log("testCase: ", testCase.name)
				t.Log("err: ", err)
				t.Log("gas: ", gas)
				if testCase.expectedError != nil {
					rpcErr := err.(rpc.Error)
					errMsg := fmt.Sprintf("[%v] expected: %v %v found: %v %v", network.Name, testCase.expectedError.ErrorCode(), testCase.expectedError.Error(), rpcErr.ErrorCode(), rpcErr.Error())
					assert.Equal(t, testCase.expectedError.ErrorCode(), rpcErr.ErrorCode(), errMsg)
					assert.True(t, strings.HasPrefix(rpcErr.Error(), testCase.expectedError.Error()), errMsg)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
