package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/urfave/cli/v2"
)

const (
	forkID4 = 4
)

func main() {
	app := cli.NewApp()
	app.Name = "RlpTool"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "decodeFull",
			Aliases: []string{},
			Usage:   "decode full callData rlp",
			Action:  decodeFull,
		},
		{
			Name:    "decode",
			Aliases: []string{},
			Usage:   "decode rlp",
			Action:  decode,
		},
		{
			Name:    "encode",
			Aliases: []string{},
			Usage:   "encode tx with rlp",
			Action:  encode,
		},
	}
	app.Flags = []cli.Flag{
		&cli.Uint64Flag{
			Name:    "forkID",
			Aliases: []string{"forkid"},
			Usage:   "forkID number",
			Value:   forkID4,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("\nError: %v\n", err)
		os.Exit(1)
	}
}

func decodeFull(ctx *cli.Context) error {
	callData := ctx.Args().First()
	bytesCallData, err := hex.DecodeHex(callData)
	if err != nil {
		log.Error("error decoding callData: ", err)
		return err
	}
	txs, rawTxs, err := decodeFullCallDataToTxs(bytesCallData, ctx.Uint64("forkID"))
	if err != nil {
		return err
	}
	printTxs(txs, rawTxs)
	return nil
}

func decode(ctx *cli.Context) error {
	rawTxs := ctx.Args().First()
	bytesRawTxs, err := hex.DecodeHex(rawTxs)
	if err != nil {
		log.Error("error decoding rawTxs: ", err)
		return err
	}
	txs, _, _, err := state.DecodeTxs(bytesRawTxs, ctx.Uint64("forkID"))
	if err != nil {
		log.Error("error decoding tx callData: ", err)
		return err
	}
	printTxs(txs, bytesRawTxs)
	return nil
}

func encode(ctx *cli.Context) error {
	fmt.Println("Nonce : ")
	var nonceS string
	if _, err := fmt.Scanln(&nonceS); err != nil {
		return err
	}
	nonce, err := strconv.ParseUint(nonceS, encoding.Base10, 64) //nolint:gomnd
	if err != nil {
		log.Error("error decoding nonce: ", err)
		return err
	}
	fmt.Println("Nonce: ", nonce)

	fmt.Println("GasPrice : ")
	var gasPriceS string
	if _, err := fmt.Scanln(&gasPriceS); err != nil {
		return err
	}
	gasPrice, _ := new(big.Int).SetString(gasPriceS, encoding.Base10)
	fmt.Println("GasPrice: ", gasPrice)

	fmt.Println("Gas : ")
	var gasS string
	if _, err := fmt.Scanln(&gasS); err != nil {
		return err
	}
	gas, err := strconv.ParseUint(gasS, encoding.Base10, 64) //nolint:gomnd
	if err != nil {
		log.Error("error decoding gas: ", err)
		return err
	}
	fmt.Println("Gas: ", gas)

	fmt.Println("To : ")
	var toS string
	if _, err := fmt.Scanln(&toS); err != nil {
		return err
	}
	to := common.HexToAddress(toS)
	fmt.Println("To: ", to)

	fmt.Println("Value : ")
	var valueS string
	if _, err := fmt.Scanln(&valueS); err != nil {
		return err
	}
	value, _ := new(big.Int).SetString(valueS, encoding.Base10)
	fmt.Println("Value: ", value)

	fmt.Println("Data : ")
	var dataS string
	if _, err := fmt.Scanln(&dataS); err != nil {
		if err.Error() != "unexpected newline" {
			return err
		}
	}
	var data []byte
	if dataS != "" {
		data, err = hex.DecodeHex(dataS)
		if err != nil {
			log.Error("error decoding data: ", err)
			return err
		}
	}
	fmt.Println("Data: ", data)

	fmt.Println("V: ")
	var vS string
	if _, err := fmt.Scanln(&vS); err != nil {
		return err
	}
	v, _ := new(big.Int).SetString(vS, encoding.Base10)
	fmt.Println("V: ", v)

	fmt.Println("R: ")
	var rS string
	if _, err := fmt.Scanln(&rS); err != nil {
		return err
	}
	r, _ := new(big.Int).SetString(rS, encoding.Base10)
	fmt.Println("R: ", r)

	fmt.Println("S: ")
	var sS string
	if _, err := fmt.Scanln(&sS); err != nil {
		return err
	}
	s, _ := new(big.Int).SetString(sS, encoding.Base10)
	fmt.Println("S: ", s)

	var txLegacy = types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	}
	tx := types.NewTx(&txLegacy)

	rawBytes, err := state.EncodeTransactions([]types.Transaction{*tx}, constants.EffectivePercentage, ctx.Uint64("forkID"))
	if err != nil {
		log.Error("error encoding txs: ", err)
		return err
	}
	fmt.Println("encoded tx with signature using RLP in []byte: ", rawBytes)
	fmt.Println("rawtx with signature using RLP in hex: ", hex.EncodeToString(rawBytes))

	return nil
}

func printTxs(txs []types.Transaction, rawTxs []byte) {
	fmt.Println("RawTxs: ", hex.EncodeToHex(rawTxs))
	for _, tx := range txs {
		fmt.Println("#######################################################################")
		fmt.Println("#######################################################################")
		fmt.Printf("Decoded tx: %+v\n", tx)
		fmt.Println("ChainID: ", tx.ChainId())
		fmt.Println("Cost: ", tx.Cost())
		fmt.Println("Data: ", hex.EncodeToString(tx.Data()))
		fmt.Println("Gas: ", tx.Gas())
		fmt.Println("GasPrice: ", tx.GasPrice())
		fmt.Println("Hash: ", tx.Hash())
		fmt.Println("Nonce: ", tx.Nonce())
		v, r, s := tx.RawSignatureValues()
		fmt.Println("V: ", v)
		fmt.Println("R: ", r)
		fmt.Println("S: ", s)
		fmt.Println("To: ", tx.To())
		fmt.Println("Type: ", tx.Type())
		fmt.Println("Value: ", tx.Value())
	}
}

func decodeFullCallDataToTxs(txsData []byte, forkID uint64) ([]types.Transaction, []byte, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is a chunk of concatenated rawTx. This rawTx is the encoded tx information in rlp + the signature information (v, r, s).
	//So, txs data will look like: txRLP+r+s+v+txRLP2+r2+s2+v2

	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
	if err != nil {
		log.Fatal("error reading smart contract abi: ", err)
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txsData[:4])
	if err != nil {
		log.Fatal("error getting abi method: ", err)
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txsData[4:])
	if err != nil {
		log.Fatal("error reading call data: ", err)
	}

	txsData = data[0].([]byte)

	txs, _, _, err := state.DecodeTxs(txsData, forkID)
	return txs, txsData, err
}
