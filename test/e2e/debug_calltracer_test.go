package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	dockersArePreLaunchedForCallTracerTests = false
)

func TestDebugTraceTransactionCallTracer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	const l2NetworkURL = "http://localhost:8124"
	const l2ExplorerRPCComponentName = "l2-explorer-json-rpc"

	var err error
	if !dockersArePreLaunchedForCallTracerTests {
		err = operations.Teardown()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, operations.Teardown())
			require.NoError(t, operations.StopComponent(l2ExplorerRPCComponentName))
		}()
	}

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	if !dockersArePreLaunchedForCallTracerTests {
		opsMan, err := operations.NewManager(ctx, opsCfg)
		require.NoError(t, err)
		err = opsMan.Setup()
		require.NoError(t, err)
		err = operations.StartComponent(l2ExplorerRPCComponentName, func() (bool, error) { return operations.NodeUpCondition(l2NetworkURL) })
		require.NoError(t, err)
	} else {
		log.Info("Using pre-launched dockers: no reset Database")
	}

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	results := map[string]json.RawMessage{}

	type testCase struct {
		name           string
		prepare        func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error)
		createSignedTx func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error)
	}
	testCases := []testCase{
		// successful transactions
		{name: "eth transfer", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment", createSignedTx: createScDeploySignedTx},
		{name: "sc call", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		{name: "create", prepare: prepareCreate, createSignedTx: createCreateSignedTx},
		{name: "create2", prepare: prepareCreate, createSignedTx: createCreate2SignedTx},
		{name: "call", prepare: prepareCalls, createSignedTx: createCallSignedTx},
		{name: "delegate call", prepare: prepareCalls, createSignedTx: createDelegateCallSignedTx},
		{name: "multi call", prepare: prepareCalls, createSignedTx: createMultiCallSignedTx},
		{name: "pre ecrecover 0", prepare: prepareCalls, createSignedTx: createPreEcrecover0SignedTx},
		{name: "chain call", prepare: prepareChainCalls, createSignedTx: createChainCallSignedTx},
		{name: "delegate transfers", prepare: prepareChainCalls, createSignedTx: createDelegateTransfersSignedTx},
		{name: "memory", prepare: prepareMemory, createSignedTx: createMemorySignedTx},
		{name: "bridge", prepare: prepareBridge, createSignedTx: createBridgeSignedTx},
		{name: "deploy create 0", createSignedTx: createDeployCreate0SignedTx},
		{name: "log0 all zeros", prepare: prepareLog0, createSignedTx: createLog0AllZeros},
		{name: "log0 empty", prepare: prepareLog0, createSignedTx: createLog0Empty},
		{name: "log0 short", prepare: prepareLog0, createSignedTx: createLog0Short},

		// failed transactions
		{name: "sc deployment reverted", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc deployment out of gas", createSignedTx: createScDeployOutOfGasSignedTx},
		// PENDING {name: "sc creation storage out of gas", createSignedTx: createScCreationCodeStorageOutOfGasSignedTx},
		{name: "sc call reverted", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
		{name: "invalid static call less parameters", prepare: prepareCalls, createSignedTx: createInvalidStaticCallLessParametersSignedTx},
		{name: "invalid static call more parameters", prepare: prepareCalls, createSignedTx: createInvalidStaticCallMoreParametersSignedTx},
		{name: "invalid static call with inner call", prepare: prepareCalls, createSignedTx: createInvalidStaticCallWithInnerCallSignedTx},
		{name: "chain call reverted", prepare: prepareChainCalls, createSignedTx: createChainCallRevertedSignedTx},
		{name: "chain delegate call reverted", prepare: prepareChainCalls, createSignedTx: createChainDelegateCallRevertedSignedTx},
		{name: "depth reverted", prepare: prepareDepth, createSignedTx: createDepthSignedTx},
	}
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	for _, network := range networks {
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(network.ChainID))
		require.NoError(t, err)

		ethereumClient := operations.MustGetClient(network.URL)
		sourceAuth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		nonce, err := ethereumClient.NonceAt(ctx, sourceAuth.From, nil)
		require.NoError(t, err)

		balance, err := ethereumClient.BalanceAt(ctx, sourceAuth.From, nil)
		require.NoError(t, err)

		gasPrice, err := ethereumClient.SuggestGasPrice(ctx)
		require.NoError(t, err)

		value := big.NewInt(0).Quo(balance, big.NewInt(2))

		gas, err := ethereumClient.EstimateGas(ctx, ethereum.CallMsg{
			From:     sourceAuth.From,
			To:       &auth.From,
			GasPrice: gasPrice,
			Value:    value,
		})
		require.NoError(t, err)

		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			To:       &auth.From,
			Nonce:    nonce,
			GasPrice: gasPrice,
			Value:    value,
			Gas:      gas,
		})

		signedTx, err := sourceAuth.Signer(sourceAuth.From, tx)
		require.NoError(t, err)

		err = ethereumClient.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}

	for tcIdx, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")

			for _, network := range networks {
				log.Debug("------------------------ ", network.Name, " ------------------------")
				ethereumClient := operations.MustGetClient(network.URL)
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(network.ChainID))
				require.NoError(t, err)

				var customData map[string]interface{}
				if tc.prepare != nil {
					customData, err = tc.prepare(t, ctx, auth, ethereumClient)
					require.NoError(t, err)
				}

				var receipt *ethTypes.Receipt
				var signedTx *ethTypes.Transaction
				forceTxIndexDifferentFromZero := tcIdx%2 == 0
				for {
					log.Debugf("forceTxIndexDifferentFromZero: %v", forceTxIndexDifferentFromZero)
					var err error
					if forceTxIndexDifferentFromZero {
						// send eth transfers txs to make the trace tx to not be the index 0 in the block
						sendEthTransfersWithoutWaiting(t, ctx, ethereumClient, auth, common.HexToAddress(operations.DefaultSequencerAddress), big.NewInt(1), 3)
					}
					signedTx, err = tc.createSignedTx(t, ctx, auth, ethereumClient, customData)
					require.NoError(t, err)

					err = ethereumClient.SendTransaction(ctx, signedTx)
					require.NoError(t, err)

					log.Debugf("tx sent: %v", signedTx.Hash().String())

					err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
					if err != nil && !strings.HasPrefix(err.Error(), "transaction has failed, reason:") {
						require.NoError(t, err)
					}

					if forceTxIndexDifferentFromZero {
						receipt, err = ethereumClient.TransactionReceipt(ctx, signedTx.Hash())
						require.NoError(t, err)
						if receipt.TransactionIndex != 0 {
							log.Debugf("tx receipt has tx index %v, accepted", receipt.TransactionIndex)
							break
						} else {
							log.Debugf("tx receipt has tx index 0, retrying")
						}
					} else {
						break
					}
				}
				debugOptions := map[string]interface{}{
					"tracer": "callTracer",
					"tracerConfig": map[string]interface{}{
						"onlyTopCall": false,
						"withLog":     true,
					},
				}

				response, err := client.JSONRPCCall(network.URL, "debug_traceTransaction", signedTx.Hash().String(), debugOptions)
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)

				results[network.Name] = response.Result
				log.Debug(string(response.Result))

				saveTraceResultToFile(t, fmt.Sprintf("callTracer_%v_%v", tcIdx, tc.name), network.Name, signedTx, response.Result, true)
			}

			referenceValueMap := map[string]interface{}{}
			err = json.Unmarshal(results[l1NetworkName], &referenceValueMap)
			require.NoError(t, err)

			for networkName, result := range results {
				if networkName == l1NetworkName {
					continue
				}

				resultMap := map[string]interface{}{}
				err = json.Unmarshal(result, &resultMap)
				require.NoError(t, err)

				compareCallFrame(t, referenceValueMap, resultMap, networkName)
			}
		})
	}
}

func compareCallFrame(t *testing.T, referenceValueMap, resultMap map[string]interface{}, networkName string) {
	require.Equal(t, referenceValueMap["from"], resultMap["from"], fmt.Sprintf("invalid `from` for network %s", networkName))
	// TODO: after we fix the full trace and the gas values for create commands, we can enable this check again.
	// require.Equal(t, referenceValueMap["gas"], resultMap["gas"], fmt.Sprintf("invalid `gas` for network %s", networkName))
	// require.Equal(t, referenceValueMap["gasUsed"], resultMap["gasUsed"], fmt.Sprintf("invalid `gasUsed` for network %s", networkName))
	require.Equal(t, referenceValueMap["input"], resultMap["input"], fmt.Sprintf("invalid `input` for network %s", networkName))
	require.Equal(t, referenceValueMap["output"], resultMap["output"], fmt.Sprintf("invalid `output` for network %s", networkName))
	require.Equal(t, referenceValueMap["value"], resultMap["value"], fmt.Sprintf("invalid `value` for network %s", networkName))
	require.Equal(t, referenceValueMap["type"], resultMap["type"], fmt.Sprintf("invalid `type` for network %s", networkName))
	require.Equal(t, referenceValueMap["error"], resultMap["error"], fmt.Sprintf("invalid `error` for network %s", networkName))
	if _, found := referenceValueMap["revertReason"]; found {
		require.Equal(t, referenceValueMap["revertReason"], resultMap["revertReason"], fmt.Sprintf("invalid `revertReason` for network %s", networkName))
	}

	referenceLogs, found := referenceValueMap["logs"].([]interface{})
	if found {
		resultLogs := resultMap["logs"].([]interface{})
		require.Equal(t, len(referenceLogs), len(resultLogs), "logs size doesn't match")
		for logIndex := range referenceLogs {
			referenceLog := referenceLogs[logIndex].(map[string]interface{})
			resultLog := resultLogs[logIndex].(map[string]interface{})

			require.Equal(t, referenceLog["data"], resultLog["data"], fmt.Sprintf("log index %v data doesn't match", logIndex))
			referenceTopics, found := referenceLog["topics"].([]interface{})
			if found {
				resultTopics := resultLog["topics"].([]interface{})
				require.Equal(t, len(referenceTopics), len(resultTopics), "log index %v topics size doesn't match", logIndex)
				for topicIndex := range referenceTopics {
					require.Equal(t, referenceTopics[topicIndex], resultTopics[topicIndex], fmt.Sprintf("log index %v topic index %v doesn't match", logIndex, topicIndex))
				}
			}
		}
	}

	referenceCalls, found := referenceValueMap["calls"].([]interface{})
	if found {
		resultCalls := resultMap["calls"].([]interface{})
		require.Equal(t, len(referenceCalls), len(resultCalls), "logs size doesn't match")
		for callIndex := range referenceCalls {
			referenceCall := referenceCalls[callIndex].(map[string]interface{})
			resultCall := resultCalls[callIndex].(map[string]interface{})

			compareCallFrame(t, referenceCall, resultCall, networkName)
		}
	}
}

func TestDebugTraceBlockCallTracer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	const l2NetworkURL = "http://localhost:8124"
	const l2ExplorerRPCComponentName = "l2-explorer-json-rpc"

	var err error
	if !dockersArePreLaunchedForCallTracerTests {
		err = operations.Teardown()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, operations.Teardown())
			require.NoError(t, operations.StopComponent(l2ExplorerRPCComponentName))
		}()
	}

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	if !dockersArePreLaunchedForCallTracerTests {
		opsMan, err := operations.NewManager(ctx, opsCfg)
		require.NoError(t, err)
		err = opsMan.Setup()
		require.NoError(t, err)
		err = operations.StartComponent(l2ExplorerRPCComponentName, func() (bool, error) { return operations.NodeUpCondition(l2NetworkURL) })
		require.NoError(t, err)
	} else {
		log.Info("Using pre-launched dockers: no reset Database")
	}

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}
	require.Equal(t, len(networks), 2, "only support 2 networks!")
	//var results map[string]map[string]interface{}
	results := map[string]map[string]interface{}{}

	type testCase struct {
		name              string
		blockNumberOrHash string
		prepare           func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error)
		createSignedTx    func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error)
	}
	testCases := []testCase{
		// successful transactions
		// by block number
		{name: "eth transfer by number", blockNumberOrHash: "number", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by number", blockNumberOrHash: "number", createSignedTx: createScDeploySignedTx},
		{name: "sc call by number", blockNumberOrHash: "number", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by number", blockNumberOrHash: "number", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		// by block hash
		{name: "eth transfer by hash", blockNumberOrHash: "hash", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by hash", blockNumberOrHash: "hash", createSignedTx: createScDeploySignedTx},
		{name: "sc call by hash", blockNumberOrHash: "hash", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by hash", blockNumberOrHash: "hash", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},

		// failed transactions
		// by block number
		{name: "sc deployment reverted by number", blockNumberOrHash: "number", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc call reverted by number", blockNumberOrHash: "number", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted by number", blockNumberOrHash: "number", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
		// by block hash
		{name: "sc deployment reverted by hash", blockNumberOrHash: "hash", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc call reverted by hash", blockNumberOrHash: "hash", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted by hash", blockNumberOrHash: "hash", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")

			for _, network := range networks {
				debugID := fmt.Sprintf("[%s/%s]", tc.name, network.Name)
				log.Debug("------------------------ ", network.Name, " ------------------------")
				ethereumClient := operations.MustGetClient(network.URL)
				auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

				var customData map[string]interface{}
				if tc.prepare != nil {
					customData, err = tc.prepare(t, ctx, auth, ethereumClient)
					require.NoError(t, err)
				}

				signedTx, err := tc.createSignedTx(t, ctx, auth, ethereumClient, customData)
				require.NoError(t, err)

				err = ethereumClient.SendTransaction(ctx, signedTx)
				require.NoError(t, err)

				log.Debugf("%s tx sent: %v", debugID, signedTx.Hash().String())

				err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
				if err != nil && !strings.HasPrefix(err.Error(), "transaction has failed, reason:") {
					require.NoError(t, err)
				}

				receipt, err := ethereumClient.TransactionReceipt(ctx, signedTx.Hash())
				require.NoError(t, err)

				debugOptions := map[string]interface{}{
					"tracer": "callTracer",
					"tracerConfig": map[string]interface{}{
						"onlyTopCall": false,
						"withLog":     true,
					},
				}

				var response types.Response
				if tc.blockNumberOrHash == "number" {
					log.Infof("%s debug_traceBlockByNumber %v", debugID, receipt.BlockNumber)
					response, err = client.JSONRPCCall(network.URL, "debug_traceBlockByNumber", hex.EncodeBig(receipt.BlockNumber), debugOptions)
				} else {
					log.Infof("%s debug_traceBlockByHash %v", debugID, receipt.BlockHash.String())
					response, err = client.JSONRPCCall(network.URL, "debug_traceBlockByHash", receipt.BlockHash.String(), debugOptions)
				}
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)
				// log.Debugf("%s response:%s", debugID, string(response.Result))

				txHash := signedTx.Hash().String()
				resultForTx := findTxInResponse(t, response.Result, txHash, debugID)
				results[network.Name] = resultForTx
			}

			referenceTransactions := results[l1NetworkName]
			resultTransactions := results[l2NetworkName]
			compareCallFrame(t, referenceTransactions, resultTransactions, l2NetworkName)
		})
	}
}

func findTxInResponse(t *testing.T, response json.RawMessage, txHash string, debugPrefix string) map[string]interface{} {
	valueMap := []interface{}{}
	err := json.Unmarshal(response, &valueMap)
	require.NoError(t, err)
	log.Infof("%s Reponse Length: %d", debugPrefix, len(valueMap))
	for transactionIndex := range valueMap {
		if valueMap[transactionIndex].(map[string]interface{})["txHash"] == txHash {
			return valueMap[transactionIndex].(map[string]interface{})
		}
	}
	log.Infof("%s Transaction not found: %s, returning first index", debugPrefix, txHash)
	return valueMap[0].(map[string]interface{})
}
