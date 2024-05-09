package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	const httpUrl = "https://zkevm-rpc.com"
	const wsUrl = "wss://ws.zkevm-rpc.com"

	const numberOfConnections = 10
	const intervalToCheckBlockNumber = 2 * time.Second

	const enableLogSubscription = true

	wg := sync.WaitGroup{}
	wg.Add(numberOfConnections)
	for connID := 0; connID < numberOfConnections; connID++ {
		go func(connID int) {
			ctx := context.Background()

			logf(connID, "connecting to: %v\n", httpUrl)
			httpClient, err := ethclient.Dial(httpUrl)
			chkErr(connID, err)
			logf(connID, "connected to: %v\n", httpUrl)

			latestBlockNumber, err := httpClient.BlockNumber(ctx)
			chkErr(connID, err)

			logf(connID, "connecting to: %v\n", wsUrl)
			wsClient, err := ethclient.Dial(wsUrl)
			chkErr(connID, err)
			logf(connID, "connected to: %v\n", wsUrl)

			signals := make(chan os.Signal, 100)
			signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

			lastWSBlockNumber := uint64(0)
			numberOfLogsReceived := uint64(0)

			// concurrently check block synchronization and logs received
			go func(connID int, httpClient *ethclient.Client) {
				for {
					if lastWSBlockNumber != 0 {
						httpBlockNumber, err := httpClient.BlockNumber(ctx)
						if err != nil {
							logf(connID, "%v failed to check block sync, retrying...\n", time.Now().Format(time.RFC3339Nano))
							time.Sleep(intervalToCheckBlockNumber)
							continue
						}

						wsBlockNumber := atomic.LoadUint64(&lastWSBlockNumber)

						diff := httpBlockNumber - wsBlockNumber
						logf(connID, "%v wsBlockNumber: %v httpBlockNumber: %v diff: %v\n", time.Now().Format(time.RFC3339Nano), wsBlockNumber, httpBlockNumber, diff)
					}
					if numberOfLogsReceived > 0 {
						logf(connID, "%v logs received: %v\n", time.Now().Format(time.RFC3339Nano), numberOfLogsReceived)
					}

					time.Sleep(intervalToCheckBlockNumber)
				}
			}(connID, httpClient)

			newHeaders := make(chan *types.Header)
			subHeaders, err := wsClient.SubscribeNewHead(ctx, newHeaders)
			chkErr(connID, err)
			logf(connID, "subscribed to newHeads\n")

			newLogs := make(chan types.Log)
			var subLogs ethereum.Subscription = &rpc.ClientSubscription{}
			if enableLogSubscription {
				subLogs, err = wsClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(0).SetUint64(latestBlockNumber),
					ToBlock:   big.NewInt(0).SetUint64(latestBlockNumber + 10000),
				}, newLogs)
				chkErr(connID, err)
				logf(connID, "subscribed to filterLogs\n")
			}

			// concurrently infinite sending messages
			go func(connID int, ctx context.Context, wsClient *ethclient.Client) {
				for {
					//bn, err := wsClient.BlockNumber(ctx)
					_, err := wsClient.BlockNumber(ctx)
					if err != nil {
						errorf(connID, "ERROR: %v\n", err.Error())
					}
					// logf(connID, "block number retrieved via message: %v\n", bn)
					time.Sleep(time.Second)
				}
			}(connID, ctx, wsClient)

		out:
			for {
				select {
				case err := <-subHeaders.Err():
					if err != nil {
						errorf(connID, "%v\n", err.Error())
						wg.Done()
						break out
					}
				case err := <-subLogs.Err():
					if err != nil {
						errorf(connID, "%v\n", err.Error())
						wg.Done()
						break out
					}
				case header := <-newHeaders:
					atomic.StoreUint64(&lastWSBlockNumber, header.Number.Uint64())
					// logf(connID, "%v L2 Block Received: %v\n", time.Now().Format(time.RFC3339Nano), header.Number.Uint64())
				case <-newLogs:
					atomic.AddUint64(&numberOfLogsReceived, 1)
					// logf(connID, "%v Log Received: %v - %v\n", time.Now().Format(time.RFC3339Nano), log.TxHash.String(), log.Index)
				case <-signals:
					subHeaders.Unsubscribe()
					if enableLogSubscription {
						subLogs.Unsubscribe()
					}
					logf(connID, "unsubscribed\n")
					close(newHeaders)
					close(newLogs)
					wg.Done()
					break out
				}
			}
		}(connID)
	}
	wg.Wait()
}

func chkErr(connID int, err error) {
	if err != nil {
		errorf(connID, err.Error())
		os.Exit(0)
	}
}

func logf(connID int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[connID: %v] %v", connID, msg)
}

func errorf(connID int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg = fmt.Sprintf("*****ERROR: %v", msg)
	logf(connID, msg)
}
