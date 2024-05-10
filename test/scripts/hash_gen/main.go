package main

import (
	"encoding/json"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

const (
	networkURL              = "https://zkevm-rpc.com"
	startBlockNumber uint64 = 10
	endBlockNumber   uint64 = 20
)

func main() {
	for blockNumber := startBlockNumber; blockNumber <= endBlockNumber; blockNumber++ {
		printfLn("getting block: %v", blockNumber)
		blockResponse, err := client.JSONRPCCall(networkURL, "eth_getBlockByNumber", hex.EncodeUint64(blockNumber), true)
		chkErr(err)
		chkRespErr(blockResponse.Error)

		rawBlock := map[string]interface{}{}
		err = json.Unmarshal(blockResponse.Result, &rawBlock)
		chkErr(err)

		// create header
		rawBlockHash := rawBlock["hash"].(string)
		number := hex.DecodeBig(rawBlock["number"].(string))
		parentHash := common.HexToHash(rawBlock["parentHash"].(string))
		coinbase := common.HexToAddress(rawBlock["miner"].(string))
		root := common.HexToHash(rawBlock["stateRoot"].(string))
		gasUsed := hex.DecodeUint64(rawBlock["gasUsed"].(string))
		gasLimit := hex.DecodeUint64(rawBlock["gasLimit"].(string))
		timeStamp := hex.DecodeUint64(rawBlock["timestamp"].(string))

		header := &ethTypes.Header{
			Number: number, ParentHash: parentHash, Coinbase: coinbase,
			Root: root, GasUsed: gasUsed, GasLimit: gasLimit, Time: timeStamp,
		}

		// create txs and receipts
		rawTransactions := rawBlock["transactions"].([]interface{})
		txs := make([]*ethTypes.Transaction, 0, len(rawTransactions))
		receipts := make([]*ethTypes.Receipt, 0, len(rawTransactions))
		for i, rawTransaction := range rawTransactions {
			if i == 1 {
				continue
			}
			rawTransactionMap := rawTransaction.(map[string]interface{})

			nonce := hex.DecodeUint64(rawTransactionMap["nonce"].(string))
			gasPrice := hex.DecodeBig(rawTransactionMap["gasPrice"].(string))
			gas := hex.DecodeUint64(rawTransactionMap["gas"].(string))
			var to *common.Address
			if rawTransactionMap["to"] != nil {
				aux := common.HexToAddress(rawTransactionMap["to"].(string))
				to = &aux
			}
			value := hex.DecodeBig(rawTransactionMap["value"].(string))
			data, _ := hex.DecodeHex(rawTransactionMap["input"].(string))
			v := hex.DecodeBig(rawTransactionMap["v"].(string))
			r := hex.DecodeBig(rawTransactionMap["r"].(string))
			s := hex.DecodeBig(rawTransactionMap["s"].(string))

			tx := ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce: nonce, GasPrice: gasPrice, Gas: gas, To: to,
				Value: value, Data: data, V: v, R: r, S: s,
			})
			txs = append(txs, tx)

			hash := rawTransactionMap["hash"].(string)
			printfLn("getting receipt for tx: %v", hash)
			receiptResponse, err := client.JSONRPCCall(networkURL, "eth_getTransactionReceipt", hash)
			chkErr(err)
			chkRespErr(receiptResponse.Error)

			rawReceipt := map[string]interface{}{}
			err = json.Unmarshal(receiptResponse.Result, &rawReceipt)
			chkErr(err)

			receiptType := uint8(hex.DecodeUint64(rawReceipt["type"].(string)))
			postState := common.HexToHash(rawReceipt["root"].(string)).Bytes()
			status := hex.DecodeUint64(rawReceipt["status"].(string))
			cumulativeGasUsed := hex.DecodeUint64(rawReceipt["cumulativeGasUsed"].(string))
			txHash := common.HexToHash(rawReceipt["transactionHash"].(string))
			var contractAddress common.Address
			if rawReceipt["contractAddress"] != nil {
				contractAddress = common.HexToAddress(rawReceipt["contractAddress"].(string))
			}
			gasUsed := hex.DecodeUint64(rawReceipt["gasUsed"].(string))
			blockHash := common.HexToHash(rawReceipt["blockHash"].(string))
			blockNumber := hex.DecodeBig(rawReceipt["blockNumber"].(string))
			transactionIndex := uint(hex.DecodeUint64(rawReceipt["transactionIndex"].(string)))

			receipt := &ethTypes.Receipt{
				Type: receiptType, PostState: postState, Status: status, CumulativeGasUsed: cumulativeGasUsed,
				TxHash: txHash, ContractAddress: contractAddress, GasUsed: gasUsed,
				BlockHash: blockHash, BlockNumber: blockNumber, TransactionIndex: transactionIndex,
			}

			rawLogs := rawReceipt["logs"].([]interface{})
			logs := make([]*ethTypes.Log, 0, len(rawLogs))
			printfLn("logs: %v", len(rawLogs))
			for _, rawLog := range rawLogs {
				rawLogMap := rawLog.(map[string]interface{})

				address := common.HexToAddress(rawLogMap["address"].(string))
				data, _ := hex.DecodeHex(rawLogMap["data"].(string))
				blockNumber := hex.DecodeUint64(rawLogMap["blockNumber"].(string))
				txHash := common.HexToHash(rawLogMap["transactionHash"].(string))
				txIndex := uint(hex.DecodeUint64(rawLogMap["transactionIndex"].(string)))
				blockHash := common.HexToHash(rawLogMap["blockHash"].(string))
				index := uint(hex.DecodeUint64(rawLogMap["logIndex"].(string)))
				removed := rawLogMap["removed"].(bool)

				log := &ethTypes.Log{
					Address:     address,
					Data:        data,
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     txIndex,
					BlockHash:   blockHash,
					Index:       index,
					Removed:     removed,
				}
				logs = append(logs, log)

				rawTopics := rawLogMap["topics"].([]interface{})
				topics := make([]common.Hash, 0, len(rawTopics))
				for _, rawTopic := range rawTopics {
					topic := common.HexToHash(rawTopic.(string))
					topics = append(topics, topic)
				}
				log.Topics = topics
			}
			receipt.Logs = logs

			// RPC is not setting the receipt bloom when computing the block hash
			// receipt.Bloom = ethTypes.CreateBloom([]*ethTypes.Receipt{receipt})

			receipts = append(receipts, receipt)
		}

		uncles := []*ethTypes.Header{}

		builtBlock := ethTypes.NewBlock(header, txs, uncles, receipts, &trie.StackTrie{})

		match := rawBlockHash == builtBlock.Hash().String()

		log.Infof("     RPC block hash: %v", rawBlockHash)
		log.Infof("Computed block hash: %v", builtBlock.Hash().String())
		if !match {
			log.Errorf(" block hashes DO NOT match")
		} else {
			log.Infof(" block hashes MATCH")
		}
	}
}

func chkRespErr(err *types.ErrorObject) {
	if err != nil {
		errMsg := fmt.Sprintf("%v %v", err.Code, err.Message)
		errorfLn(errMsg)
		panic(err)
	}
}

func chkErr(err error) {
	if err != nil {
		errorfLn(err.Error())
		panic(err)
	}
}

func errorfLn(format string, args ...interface{}) {
	printfLn("ERROR: "+format, args...)
}

func printfLn(format string, args ...interface{}) {
	fmt.Printf(format+" \n", args...)
}
