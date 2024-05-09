package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// add here the url of the nodes you want to check
// against the trusted node
var networkURLsToCheck = []string{
	// "https://add.your.node.url.here",
	// "https://add.your.node.url.here",
	// "https://add.your.node.url.here",
}

// set the from and to block numbers you want to verify
const fromBlockNumber uint64 = 10
const toBlockNumber uint64 = 20

// pick the correct trusted Node URL depending on the network you are testing

// mainnet
const trustedNodeURL = "https://zkevm-rpc.com"

// cardona
// const trustedNodeURL = "https://rpc.cardona.zkevm-rpc.com/"

func main() {
	fmt.Printf("connecting to network: %v ...", trustedNodeURL)
	trustedNodeClient, err := ethclient.Dial(trustedNodeURL)
	chkErr(err)
	fmt.Print("connected")
	fmt.Println()

	networkClients := map[string]*ethclient.Client{}
	for _, networkURL := range networkURLsToCheck {
		fmt.Printf("connecting to network: %v ...", networkURL)
		client, err := ethclient.Dial(networkURL)
		chkErr(err)
		networkClients[networkURL] = client
		fmt.Print("connected")
		fmt.Println()
	}

	for blockNumberU64 := fromBlockNumber; blockNumberU64 <= toBlockNumber; blockNumberU64++ {
		ctx := context.Background()
		blockNumber := big.NewInt(0).SetUint64(blockNumberU64)
		fmt.Println()
		fmt.Println("block to verify: ", blockNumberU64)

		// load blocks from trusted node
		trustedNodeBlockHeader, err := trustedNodeClient.HeaderByNumber(ctx, blockNumber)
		chkErr(err)
		const logPattern = "block: %v hash: %v parentHash: %v network: %v\n"
		trustedNodeBlockHash := trustedNodeBlockHeader.Hash().String()
		trustedNodeParentBlockHash := trustedNodeBlockHeader.ParentHash.String()

		// load blocks from networks to verify
		blocks := sync.Map{}
		wg := sync.WaitGroup{}
		wg.Add(len(networkURLsToCheck))
		for _, networkURL := range networkURLsToCheck {
			go func(networkURL string) {
				defer wg.Done()
				c := networkClients[networkURL]

				blockHeader, err := c.HeaderByNumber(ctx, blockNumber)
				if errors.Is(err, ethereum.NotFound) {
					return
				} else {
					chkErr(err)
				}

				blocks.Store(networkURL, blockHeader)
			}(networkURL)
		}
		wg.Wait()

		failed := false
		blocks.Range(func(networkURLValue, blockValue any) bool {
			networkURL, block := networkURLValue.(string), blockValue.(*types.Header)

			// when block is not found
			if block == nil {
				fmt.Printf(logPattern, blockNumberU64, "NOT FOUND", "NOT FOUND", networkURL)
				return true
			}

			blockHash := block.Hash().String()
			parentBlockHash := block.ParentHash.String()

			if trustedNodeBlockHash != blockHash || trustedNodeParentBlockHash != parentBlockHash {
				failed = true
				fmt.Printf(logPattern, blockNumberU64, trustedNodeBlockHash, trustedNodeParentBlockHash, trustedNodeURL)
				fmt.Printf(logPattern, blockNumberU64, blockHash, parentBlockHash, networkURL)
				fmt.Printf("ERROR block information mismatch for network: %v\n", networkURL)
			} else {
				fmt.Printf("%v: OK\n", networkURL)
			}

			return true
		})
		if failed {
			panic("block information mismatch")
		}

		// avoid getting blocked by request rate limit
		time.Sleep(time.Second)
	}
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}
