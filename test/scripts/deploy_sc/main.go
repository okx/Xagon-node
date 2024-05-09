package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Counter"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Storage"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	txTimeout = 60 * time.Second
)

func main() {
	var networks = []struct {
		Name       string
		URL        string
		ChainID    uint64
		PrivateKey string
	}{
		{Name: "Local L1", URL: operations.DefaultL1NetworkURL, ChainID: operations.DefaultL1ChainID, PrivateKey: operations.DefaultSequencerPrivateKey},
		{Name: "Local L2", URL: operations.DefaultL2NetworkURL, ChainID: operations.DefaultL2ChainID, PrivateKey: operations.DefaultSequencerPrivateKey},
	}

	for _, network := range networks {
		ctx := context.Background()

		log.Infof("connecting to %v: %v", network.Name, network.URL)
		client, err := ethclient.Dial(network.URL)
		chkErr(err)
		log.Infof("connected")

		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)
		chkErr(err)

		const receiverAddr = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"

		balance, err := client.BalanceAt(ctx, auth.From, nil)
		chkErr(err)
		log.Debugf("ETH Balance for %v: %v", auth.From, balance)

		// Counter
		log.Debugf("Sending TX to deploy Counter SC")
		_, tx, counterSC, err := Counter.DeployCounter(auth, client)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		log.Debugf("Calling Increment method from Counter SC")
		tx, err = counterSC.Increment(auth)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		fmt.Println()

		// EmitLog
		log.Debugf("Sending TX to deploy EmitLog SC")
		_, tx, emitLogSC, err := EmitLog.DeployEmitLog(auth, client)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		log.Debugf("Calling EmitLogs method from EmitLog SC")
		tx, err = emitLogSC.EmitLogs(auth)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		fmt.Println()

		// ERC20
		mintAmount, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
		log.Debugf("Sending TX to deploy ERC20 SC")
		_, tx, erc20SC, err := ERC20.DeployERC20(auth, client, "Test Coin", "TCO")
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		log.Debugf("Sending TX to do a ERC20 mint")
		tx, err = erc20SC.Mint(auth, mintAmount)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		transferAmount, _ := big.NewInt(0).SetString("900000000000000000000", encoding.Base10)
		log.Debugf("Sending TX to do a ERC20 transfer")
		tx, err = erc20SC.Transfer(auth, common.HexToAddress(receiverAddr), transferAmount)
		chkErr(err)
		auth.Nonce = big.NewInt(0).SetUint64(tx.Nonce() + 1)
		log.Debugf("Sending invalid TX to do a ERC20 transfer")
		invalidTx, err := erc20SC.Transfer(auth, common.HexToAddress(receiverAddr), transferAmount)
		chkErr(err)
		log.Debugf("Invalid ERC20 tx hash: %v", invalidTx.Hash())
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		operations.WaitTxToBeMined(ctx, client, invalidTx, txTimeout) //nolint:errcheck
		chkErr(err)
		auth.Nonce = nil
		fmt.Println()

		// Storage
		const numberToStore = 22
		log.Debugf("Sending TX to deploy Storage SC")
		_, tx, storageSC, err := Storage.DeployStorage(auth, client)
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		log.Debugf("Calling Store method from Storage SC")
		tx, err = storageSC.Store(auth, big.NewInt(numberToStore))
		chkErr(err)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		fmt.Println()

		// Valid ETH Transfer
		balance, err = client.BalanceAt(ctx, auth.From, nil)
		log.Debugf("ETH Balance for %v: %v", auth.From, balance)
		chkErr(err)
		const halfDivision = 2
		transferAmount = balance.Quo(balance, big.NewInt(halfDivision))
		log.Debugf("Transfer Amount: %v", transferAmount)

		log.Debugf("Sending TX to transfer ETH")
		to := common.HexToAddress(receiverAddr)
		tx = ethTransfer(ctx, client, auth, to, transferAmount, nil)
		err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
		chkErr(err)
		fmt.Println()
	}
}

func ethTransfer(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, to common.Address, amount *big.Int, nonce *uint64) *types.Transaction {
	if nonce == nil {
		log.Infof("reading nonce for account: %v", auth.From.Hex())
		var err error
		n, err := client.NonceAt(ctx, auth.From, nil)
		log.Infof("nonce: %v", n)
		chkErr(err)
		nonce = &n
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	chkErr(err)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{To: &to})
	chkErr(err)

	tx := types.NewTransaction(*nonce, to, amount, gasLimit, gasPrice, nil)

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("sending transfer tx")
	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	log.Infof("tx sent: %v", signedTx.Hash().Hex())

	rlp, err := signedTx.MarshalBinary()
	chkErr(err)

	log.Infof("tx rlp: %v", hex.EncodeToHex(rlp))

	return signedTx
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
