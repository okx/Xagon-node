package setup

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	sleepDuration                         = 5 * time.Second
	minAllowedGasPriceIntervalMinutes     = 5
	pollMinAllowedGasPriceIntervalSeconds = 15
	defaultGasPrice                       = 1000000000
)

var (
	bc = state.BatchConstraintsCfg{
		MaxTxsPerBatch:       300,
		MaxBatchBytesSize:    120000,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      2145,
		MaxPoseidonHashes:    252357,
		MaxPoseidonPaddings:  135191,
		MaxMemAligns:         236585,
		MaxArithmetics:       236585,
		MaxBinaries:          473170,
		MaxSteps:             7570538,
	}
)

// Environment sets up the environment for the benchmark
func Environment(ctx context.Context, b *testing.B) (*operations.Manager, *ethclient.Client, *pool.Pool, *bind.TransactOpts) {
	if testing.Short() {
		b.Skip()
	}

	err := operations.Teardown()
	require.NoError(b, err)

	params.OpsCfg.State.MaxCumulativeGasUsed = params.MaxCumulativeGasUsed
	opsman, err := operations.NewManager(ctx, params.OpsCfg)
	require.NoError(b, err)

	err = Components(opsman)
	require.NoError(b, err)
	time.Sleep(sleepDuration)

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(b, err)

	// Load params client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(b, err)

	st := opsman.State()
	s, err := pgpoolstorage.NewPostgresPoolStorage(params.PoolDbConfig)
	require.NoError(b, err)
	config := pool.Config{
		DB:                             params.PoolDbConfig,
		MinAllowedGasPriceInterval:     types.NewDuration(minAllowedGasPriceIntervalMinutes * time.Minute),
		PollMinAllowedGasPriceInterval: types.NewDuration(pollMinAllowedGasPriceIntervalSeconds * time.Second),
	}

	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(b, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	pl := pool.NewPool(config, bc, s, st, common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), params.ChainID, eventLog)

	// Print Info before send
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(b, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(b, err)

	// Print Initial Stats
	fmt.Printf("Receiver Addr: %v\n", params.To.String())
	fmt.Printf("Sender Addr: %v\n", auth.From.String())
	fmt.Printf("Sender Balance: %v\n", senderBalance.String())
	fmt.Printf("Sender Nonce: %v\n", senderNonce)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)

	if gasPrice == nil || gasPrice.Int64() == 0 {
		gasPrice = big.NewInt(defaultGasPrice)
	}

	// PrivateKey is the private key of the sender
	// Auth is the auth of the sender
	auth, err = bind.NewKeyedTransactorWithChainID(params.PrivateKey, new(big.Int).SetUint64(params.ChainID))
	if err != nil {
		panic(err)
	}
	auth.GasPrice = gasPrice

	return opsman, client, pl, auth
}

// Components runs the network container, starts synchronizer and JSON-RPC components, and approves matic
func Components(opsman *operations.Manager) error {
	// Run network container
	err := opsman.StartNetwork()
	if err != nil {
		return err
	}

	// Approve matic
	err = operations.ApproveMatic()
	if err != nil {
		return err
	}

	err = operations.StartComponent("sync")
	if err != nil {
		return err
	}

	err = operations.StartComponent("json-rpc")
	if err != nil {
		return err
	}
	time.Sleep(sleepDuration)

	return nil
}

// BootstrapSequencer starts the sequencer and waits for it to be ready
func BootstrapSequencer(b *testing.B, opsman *operations.Manager) {
	fmt.Println("Starting sequencer ....")
	err := operations.StartComponent("seq")
	require.NoError(b, err)
	fmt.Println("Sequencer Started!")
}
