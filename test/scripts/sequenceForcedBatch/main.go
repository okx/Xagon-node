package main

import (
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	flagL1URLName        = "url"
	flagSmcAddrName      = "smc"
	flagGerName          = "ger"
	flagTimestampName    = "timestamp"
	flagTransactionsName = "transactions"
	miningTimeout        = 180
)

var (
	flagL1URL = cli.StringFlag{
		Name:     flagL1URLName,
		Aliases:  []string{"u"},
		Usage:    "L1 node url",
		Required: true,
	}
	flagSmcAddr = cli.StringFlag{
		Name:     flagSmcAddrName,
		Aliases:  []string{"a"},
		Usage:    "Smart contract address",
		Required: true,
	}
	flagGer = cli.StringFlag{
		Name:     flagGerName,
		Aliases:  []string{"g"},
		Usage:    "Global exit root",
		Required: true,
	}
	flagTimestamp = cli.StringFlag{
		Name:     flagTimestampName,
		Aliases:  []string{"t"},
		Usage:    "MinForcedTimestamp",
		Required: true,
	}
	flagTransactions = cli.StringFlag{
		Name:     flagTransactionsName,
		Aliases:  []string{"tx"},
		Usage:    "Transactions",
		Required: true,
	}
)

func main() {
	fbatchsender := cli.NewApp()
	fbatchsender.Name = "sequenceForcedBatchsender"
	fbatchsender.Usage = "send sequencer forced batch transactions to L1"
	fbatchsender.DefaultCommand = "send"
	flags := []cli.Flag{&flagL1URL, &flagSmcAddr, &flagGer, &flagTimestamp, &flagTransactions}
	fbatchsender.Commands = []*cli.Command{
		{
			Before:  setLogLevel,
			Name:    "send",
			Aliases: []string{},
			Flags:   flags,
			Action:  sendForcedBatches,
		},
	}

	err := fbatchsender.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setLogLevel(ctx *cli.Context) error {
	logLevel := "debug"
	log.Init(log.Config{
		Level:   logLevel,
		Outputs: []string{"stderr"},
	})
	return nil
}

func sendForcedBatches(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	url := cliCtx.String(flagL1URLName)
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", url, err)
		return err
	}
	// Create smc client
	zkevmAddr := common.HexToAddress(cliCtx.String(flagSmcAddrName))
	zkevm, err := etrogpolygonzkevm.NewEtrogpolygonzkevm(zkevmAddr, ethClient)
	if err != nil {
		return err
	}

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}

	log.Info("Using address: ", auth.From)

	num, err := zkevm.LastForceBatch(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("error getting lastForBatch number. Error : ", err)
		return err
	}
	log.Info("Number of forceBatches in the smc: ", num)

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		log.Error("error getting blockByNumber. Error: ", err)
		return err
	}
	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	transactions, err := hex.DecodeHex(cliCtx.String(flagTransactionsName))
	if err != nil {
		log.Error("error decoding txs. Error: ", err)
		return err
	}
	fbData := []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{{
		Transactions:       transactions,
		ForcedGlobalExitRoot:     common.HexToHash(cliCtx.String(flagGerName)),
		ForcedTimestamp: cliCtx.Uint64(flagTimestampName),
	}}
	log.Warnf("%v, %+v", cliCtx.String(flagTransactionsName), fbData)
	// Send forceBatch
	tx, err := zkevm.SequenceForceBatches(auth, fbData)
	if err != nil {
		log.Error("error sending forceBatch. Error: ", err)
		return err
	}

	log.Info("TxHash: ", tx.Hash())

	time.Sleep(1 * time.Second)

	return operations.WaitTxToBeMined(ctx, ethClient, tx, miningTimeout*time.Second)
}
