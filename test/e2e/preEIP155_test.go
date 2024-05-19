package e2e

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestPreEIP155Tx(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// Edit config
	const path = "../../test/config/test.node.config.toml"
	require.NoError(t,
		exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 1000000000/DefaultMinGasPriceAllowed = 0/g", path).Run(),
	)
	require.NoError(t,
		exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = true/EnableL2SuggestedGasPricePolling = false/g", path).Run(),
	)
	// Undo edit config
	defer func() {
		require.NoError(t,
			exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 0/DefaultMinGasPriceAllowed = 1000000000/g", path).Run(),
		)
		require.NoError(t,
			exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = false/EnableL2SuggestedGasPricePolling = true/g", path).Run(),
		)
	}()

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		priKey := network.PrivateKey
		if network.Name == "Local L2" {
			priKey = "0x815405dddb0e2a99b12af775fd2929e526704e1d1aea6a0b4e74dc33e2f7fcd2"
		}
		auth := operations.MustGetAuth(priKey, network.ChainID)

		nonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)

		to := common.HexToAddress("0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98")
		data := hex.DecodeBig("0x64fbb77c").Bytes()

		gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From: auth.From,
			To:   &to,
			Data: data,
		})
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			GasPrice: gasPrice,
			Gas:      gas,
			Data:     data,
		})

		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(priKey, "0x"))
		require.NoError(t, err)

		signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
		require.NoError(t, err)

		//log.Debug("privateKey:", privateKey, ", from:", auth.From)
		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, signedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		require.NoError(t, err)

		// wait for l2 block to be virtualized
		if network.ChainID == operations.DefaultL2ChainID {
			log.Infof("waiting for the block number %v to be virtualized", receipt.BlockNumber.String())
			err = operations.WaitL2BlockToBeVirtualized(receipt.BlockNumber, 4*time.Minute) //nolint:gomnd
			require.NoError(t, err)
		}
	}
}

func TestFakeEIP155With_V_As35(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// Edit config
	const path = "../../test/config/test.node.config.toml"
	require.NoError(t,
		exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 1000000000/DefaultMinGasPriceAllowed = 0/g", path).Run(),
	)
	require.NoError(t,
		exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = true/EnableL2SuggestedGasPricePolling = false/g", path).Run(),
	)
	// Undo edit config
	defer func() {
		require.NoError(t,
			exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 0/DefaultMinGasPriceAllowed = 1000000000/g", path).Run(),
		)
		require.NoError(t,
			exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = false/EnableL2SuggestedGasPricePolling = true/g", path).Run(),
		)
	}()

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)

		privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		require.NoError(t, err)
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		require.NoError(t, err)

		toAddress := common.HexToAddress("0x1234")
		tx := &types.LegacyTx{
			Nonce: nonce,
			To:    &toAddress,
			Value: big.NewInt(0),
			Gas:   uint64(21000),

			GasPrice: big.NewInt(10000000000000),
			Data:     nil,
		}

		// set the chainID to 0 to fake a pre EIP155 tx
		signer := types.NewEIP155Signer(big.NewInt(0))

		// sign tx
		h := signer.Hash(types.NewTx(tx))
		sig, err := crypto.Sign(h[:], privateKey)
		require.NoError(t, err)
		r, s, _, err := signer.SignatureValues(types.NewTx(tx), sig)
		require.NoError(t, err)

		// set the value V of the signature to 35
		tx.V = big.NewInt(35)
		tx.R = r
		tx.S = s

		signedTx := types.NewTx(tx)
		err = client.SendTransaction(context.Background(), signedTx)
		require.Equal(t, "invalid sender", err.Error())
	}
}
