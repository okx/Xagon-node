package main

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/pol"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	flagL1URLName             = "url"
	flagZkevmAddrName         = "zkevm"
	flagRollupManagerAddrName = "rollupmanager"
	miningTimeout             = 180
)

var (
	flagL1URL = cli.StringFlag{
		Name:     flagL1URLName,
		Aliases:  []string{"u"},
		Usage:    "L1 node url",
		Required: true,
	}
	flagZkevmAddr = cli.StringFlag{
		Name:     flagZkevmAddrName,
		Aliases:  []string{"zk"},
		Usage:    "Zkevm smart contract address",
		Required: true,
	}
	flagRollupManagerAddr = cli.StringFlag{
		Name:     flagRollupManagerAddrName,
		Aliases:  []string{"r"},
		Usage:    "RollupmManager smart contract address",
		Required: true,
	}
)

func main() {
	fbatchsender := cli.NewApp()
	fbatchsender.Name = "forcedBatchsender"
	fbatchsender.Usage = "send forced batch transactions to L1"
	fbatchsender.DefaultCommand = "send"
	flags := []cli.Flag{&flagL1URL, &flagZkevmAddr, &flagRollupManagerAddr}
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

func transferERC20Pol(ctx context.Context, ethClient *ethclient.Client, authSequencer, authForcedBatch *bind.TransactOpts, zkevmAddr common.Address) error {
	log.Infof("Transfering POL from sequencer to forcedBatchesAddress")
	polSmc, err := pol.NewPol(common.HexToAddress(operations.DefaultL1PolSmartContract), ethClient)
	if err != nil {
		return err
	}
	polAmount, _ := big.NewInt(0).SetString("9999999999999999999999", 0)
	log.Debugf("Charging pol from sequencer -> forcedBatchesAddress")
	txValue, err := polSmc.Transfer(authSequencer, common.HexToAddress(operations.DefaultForcedBatchesAddress), polAmount)
	if err != nil {
		return err
	}
	log.Debugf("Waiting for tx %s to be mined (transfer of pol from sequencer -> forcedBatches)", txValue.Hash().String())
	err = operations.WaitTxToBeMined(ctx, ethClient, txValue, operations.DefaultTimeoutTxToBeMined)
	if err != nil {
		return err
	}
	balance, err := polSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultSequencerAddress))
	if err != nil {
		return err
	}
	log.Debugf("Account (sequencer) %s pol balance %s", operations.DefaultSequencerAddress, balance.String())
	balance, err = polSmc.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress(operations.DefaultForcedBatchesAddress))
	if err != nil {
		return err
	}
	log.Debugf("Account (force_batches) %s pol balance %s", operations.DefaultForcedBatchesAddress, balance.String())
	log.Debugf("Approve to zkEVM SMC (%s) to spend %s pol", zkevmAddr, polAmount.String())
	_, err = polSmc.Approve(authForcedBatch, zkevmAddr, polAmount)
	if err != nil {
		return err
	}

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
	zkevmAddr := common.HexToAddress(cliCtx.String(flagZkevmAddrName))
	zkevm, err := etrogpolygonzkevm.NewEtrogpolygonzkevm(zkevmAddr, ethClient)
	if err != nil {
		return err
	}

	rollupManagerAddr := common.HexToAddress(cliCtx.String(flagRollupManagerAddrName))
	rollupManager, err := etrogpolygonrollupmanager.NewEtrogpolygonrollupmanager(rollupManagerAddr, ethClient)
	if err != nil {
		return err
	}
	authSeq, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}
	auth, err := operations.GetAuth(operations.DefaultForcedBatchesPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}
	err = transferERC20Pol(ctx, ethClient, authSeq, auth, zkevmAddr)
	if err != nil {
		log.Error("error transfering pol. Error: ", err)
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

	// Get tip
	tip, err := rollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("error getting tip. Error: ", err)
		return err
	}
	log.Info("Tip: ", tip)
	tx, err := zkevm.SetForceBatchAddress(authSeq, common.Address{})
	if err != nil {
		log.Error("error sending SetForceBatchAddress. Error: ", err)
		return err
	}
	err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
	if err != nil {

		log.Error("error waiting tx to be mined. Error: ", err)
		return err
	}

	txs := "ee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e980801186622d03b6b8da7cf111d1ccba5bb185c56deae6a322cebc6dda0556f3cb9700910c26408b64b51c5da36ba2f38ef55ba1cee719d5a6c012259687999074321bff"
	// Send forceBatch
	data := common.Hex2Bytes(txs)
	log.Info("Data: ", data)
	tx, err = zkevm.ForceBatch(auth, data, tip)
	if err != nil {
		log.Error("error sending forceBatch. Error: ", err)
		return err
	}

	log.Info("TxHash: ", tx.Hash())

	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, ethClient, tx, miningTimeout*time.Second)
	if err != nil {
		return err
	}

	query := ethereum.FilterQuery{
		FromBlock: currentBlock.Number(),
		Addresses: []common.Address{zkevmAddr},
	}
	logs, err := ethClient.FilterLogs(ctx, query)
	if err != nil {
		return err
	}
	for _, vLog := range logs {
		fb, err := zkevm.ParseForceBatch(vLog)
		if err == nil {
			log.Debugf("log decoded: %+v", fb)
			ger := fb.LastGlobalExitRoot
			log.Info("GlobalExitRoot: ", ger)
			log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
			fullBlock, err := ethClient.BlockByHash(ctx, vLog.BlockHash)
			if err != nil {
				log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
				return err
			}
			log.Info("MinForcedTimestamp: ", fullBlock.Time())
		}
	}

	return nil
}
