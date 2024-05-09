package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	flagSequencesName = "sequences"
	flagWaitName      = "wait"
	flagVerboseName   = "verbose"
)

var (
	flagSequences = cli.Uint64Flag{
		Name:     flagSequencesName,
		Aliases:  []string{"s"},
		Usage:    "send batches for the provided number of sequences.",
		Required: false,
	}
	flagWait = cli.BoolFlag{
		Name:     flagWaitName,
		Aliases:  []string{"w"},
		Usage:    "wait batch transaction to be confirmed",
		Required: false,
	}
	flagVerbose = cli.BoolFlag{
		Name:     flagVerboseName,
		Aliases:  []string{"v"},
		Usage:    "output verbose logs",
		Required: false,
	}
)

func main() {
	batchsender := cli.NewApp()
	batchsender.Name = "batchsender"
	batchsender.Usage = "send batch transactions to L1"
	batchsender.Description = `This tool allows to send a specified number of batch transactions to L1. 
Optionally it can wait for the batches to be validated.`
	batchsender.DefaultCommand = "send"
	batchsender.Flags = []cli.Flag{&flagSequences, &flagWait, &flagVerbose}
	batchsender.Commands = []*cli.Command{
		{
			Before:  setLogLevel,
			Name:    "send",
			Aliases: []string{},
			Usage:   "Sends the specified number of batch transactions to L1",
			Description: `This command sends the specified number of batches to L1.
If --sequences flag is used, the number of batches is repeated for the number of sequences provided.
If --wait flag is used, it waits for the corresponding validation transaction.`,
			ArgsUsage: "number of batches to be sent (default: 1)",
			Action:    sendBatches,
		},
	}

	err := batchsender.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func setLogLevel(ctx *cli.Context) error {
	logLevel := "info"
	if ctx.Bool(flagVerboseName) {
		logLevel = "debug"
	}

	log.Init(log.Config{
		Level:   logLevel,
		Outputs: []string{"stderr"},
	})
	return nil
}

func sendBatches(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	// retrieve default configuration
	var cfg config.Config
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(config.DefaultValues)))
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return err
	}

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}

	ethMan, err := etherman.NewClient(cfg.Etherman, cfg.NetworkConfig.L1Config)
	if err != nil {
		return err
	}

	err = ethMan.AddOrReplaceAuth(*auth)
	if err != nil {
		return err
	}
	log.Info("Using address: ", auth.From)

	wait := cliCtx.Bool(flagWaitName)

	nBatches := 1 // send 1 batch by default
	if cliCtx.NArg() > 0 {
		nBatchesArgStr := cliCtx.Args().Get(0)
		nBatchesArg, err := strconv.Atoi(nBatchesArgStr)
		if err == nil {
			nBatches = nBatchesArg
		}
	}

	nSequences := int(cliCtx.Uint64(flagSequencesName))

	var sentTxs []*ethTypes.Transaction
	sentTxsMap := make(map[common.Hash]struct{})

	var duration time.Duration
	if nBatches > 1 {
		duration = 500
	}

	// here the behavior is different:
	// - if the `--sequences` flag is used we send ns sequences filled with nb batches each
	// - if the flag is not used we send one sequence for each batch
	var ns, nb int
	if nSequences == 0 {
		ns = nBatches
		nb = 1
	} else {
		ns = nSequences
		nb = nBatches
	}

	nonce, err := ethMan.CurrentNonce(ctx, auth.From)
	if err != nil {
		err := fmt.Errorf("failed to get current nonce: %w", err)
		log.Error(err.Error())
		return err
	}

	for i := 0; i < ns; i++ {
		currentBlock, err := ethMan.EthClient.BlockByNumber(ctx, nil)
		if err != nil {
			return err
		}
		log.Debug("currentBlock.Time(): ", currentBlock.Time())

		seqs := make([]ethmanTypes.Sequence, 0, nBatches)
		for i := 0; i < nb; i++ {
			// empty rollup
			seqs = append(seqs, ethmanTypes.Sequence{
				BatchNumber:          uint64(i),
				GlobalExitRoot:       common.HexToHash("0x"),
				BatchL2Data:          []byte{},
				LastL2BLockTimestamp: int64(currentBlock.Time() - 1), // fit in latest-sequence < > current-block rage
			})
		}

		// send to L1
		firstSequence := seqs[0]
		lastSequence := seqs[len(seqs)-1]
		to, data, err := ethMan.BuildSequenceBatchesTxData(auth.From, seqs, uint64(lastSequence.LastL2BLockTimestamp), firstSequence.BatchNumber, auth.From)
		if err != nil {
			return err
		}
		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			To:   to,
			Data: data,
		})
		signedTx, err := ethMan.SignTx(ctx, auth.From, tx)
		if err != nil {
			return err
		}
		err = ethMan.SendTx(ctx, signedTx)
		if err != nil {
			return err
		}
		gas, err := ethMan.EstimateGas(ctx, auth.From, to, nil, data)
		if err != nil {
			err := fmt.Errorf("failed to estimate gas: %w", err)
			log.Error(err.Error())
			return err
		}
		// get gas price
		gasPrice, err := ethMan.SuggestedGasPrice(ctx)
		if err != nil {
			err := fmt.Errorf("failed to get suggested gas price: %w", err)
			log.Error(err.Error())
			return err
		}
		tx = ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    nonce,
			Gas:      gas + uint64(i),
			GasPrice: gasPrice,
			To:       to,
			Data:     data,
		})
		signedTx, err = ethMan.SignTx(ctx, auth.From, tx)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = ethMan.SendTx(ctx, signedTx)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		log.Info("TxHash: ", signedTx.Hash())
		sentTxs = append(sentTxs, signedTx)
		sentTxsMap[signedTx.Hash()] = struct{}{}

		time.Sleep(duration * time.Millisecond)
	}

	sentBatches := len(sentTxs)

	if wait { // wait proofs
		log.Info("Waiting for transactions to be confirmed...")
		time.Sleep(time.Second)

		virtualBatches := make(map[uint64]common.Hash)
		verifiedBatches := make(map[uint64]struct{})
		loggedBatches := make(map[uint64]struct{})

		miningTimeout := 180 * time.Second                           //nolint:gomnd
		waitTimeout := time.Duration(180*len(sentTxs)) * time.Second //nolint:gomnd
		done := make(chan struct{})

		for _, tx := range sentTxs {
			err := operations.WaitTxToBeMined(ctx, ethMan.EthClient, tx, miningTimeout)
			if err != nil {
				return err
			}
		}

		for {
			select {
			case <-time.After(waitTimeout):
				return errors.New("deadline exceeded")
			case <-done:
				success(sentBatches)
				return nil
			default:
			txLoop:
				for _, tx := range sentTxs {
					// get rollup tx block number
					receipt, err := ethMan.EthClient.TransactionReceipt(ctx, tx.Hash())
					if err != nil {
						return err
					}

					fromBlock := receipt.BlockNumber
					toBlock := new(big.Int).Add(fromBlock, new(big.Int).SetUint64(cfg.Synchronizer.SyncChunkSize))
					query := ethereum.FilterQuery{
						FromBlock: fromBlock,
						ToBlock:   toBlock,
						Addresses: ethMan.SCAddresses,
					}
					logs, err := ethMan.EthClient.FilterLogs(ctx, query)
					if err != nil {
						return err
					}
					for _, vLog := range logs {
						switch vLog.Topics[0] {
						case etherman.SequencedBatchesSigHash():
							if vLog.TxHash == tx.Hash() { // ignore other txs happening on L1
								sb, err := ethMan.EtrogZkEVM.ParseSequenceBatches(vLog)
								if err != nil {
									return err
								}

								virtualBatches[sb.NumBatch] = vLog.TxHash

								if _, logged := loggedBatches[sb.NumBatch]; !logged {
									log.Infof("Batch [%d] virtualized in TxHash [%v]", sb.NumBatch, vLog.TxHash)
									loggedBatches[sb.NumBatch] = struct{}{}
								}
							}
						case etherman.TrustedVerifyBatchesSigHash():
							vb, err := ethMan.EtrogZkEVM.ParseVerifyBatches(vLog)
							if err != nil {
								return err
							}

							if _, verified := verifiedBatches[vb.NumBatch]; !verified {
								log.Infof("Batch [%d] verified in TxHash [%v]", vb.NumBatch, vLog.TxHash)
								verifiedBatches[vb.NumBatch] = struct{}{}
							}

							// batch is verified, remove it from the txs set
							delete(sentTxsMap, virtualBatches[vb.NumBatch])
							if len(sentTxsMap) == 0 {
								close(done)
								break txLoop
							}
						}
					}

					// wait for verifications
					time.Sleep(time.Second) //nolint:gomnd
				}
			}
		}
	}

	success(sentBatches)
	return nil
}

func success(nBatches int) {
	if nBatches > 1 {
		log.Infof("Successfully sent %d batches", nBatches)
	} else {
		log.Info("Successfully sent 1 batch")
	}
}
