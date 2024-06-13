package pool_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Revert"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	forkID6          = 6
	senderPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	senderAddress    = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
)

var (
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	poolDBCfg  = dbutils.NewPoolConfigFromEnv()
	genesis    = state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "90000000000000000000000000000000000000000000000000000000000",
			},
		},
	}
	cfg = pool.Config{
		MaxTxBytesSize:                    30132,
		MaxTxDataBytesSize:                30000,
		MinAllowedGasPriceInterval:        cfgTypes.NewDuration(5 * time.Minute),
		PollMinAllowedGasPriceInterval:    cfgTypes.NewDuration(15 * time.Second),
		DefaultMinGasPriceAllowed:         1000000000,
		IntervalToRefreshBlockedAddresses: cfgTypes.NewDuration(5 * time.Minute),
		IntervalToRefreshGasPrices:        cfgTypes.NewDuration(5 * time.Second),
		AccountQueue:                      15,
		GlobalQueue:                       20,
		TxFeeCap:                          1,
		EffectiveGasPrice: pool.EffectiveGasPriceCfg{
			Enabled:                     true,
			L1GasPriceFactor:            0.25,
			ByteGasCost:                 16,
			ZeroByteGasCost:             4,
			NetProfit:                   1,
			BreakEvenFactor:             1.1,
			FinalDeviationPct:           10,
			EthTransferGasPrice:         0,
			EthTransferL1GasPriceFactor: 0,
			L2GasPriceSuggesterFactor:   0.5,
		},
	}
	gasPrice   = big.NewInt(1000000000)
	l1GasPrice = big.NewInt(1000000000000)
	gasLimit   = uint64(21000)
	chainID    = big.NewInt(1337)
	bc         = state.BatchConstraintsCfg{
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
		MaxSHA256Hashes:      1596,
	}
	ip = "101.1.50.20"
)

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})

	code := m.Run()
	os.Exit(code)
}

type testData struct {
	pool *pool.Pool
	st   *state.State

	stateSqlDB *pgxpool.Pool
	poolSqlDB  *pgxpool.Pool
}

func Test_AddTxEGPAceptedBecauseGasPriceIsTheSuggested(t *testing.T) {
	ctx := context.Background()

	data := prepareToExecuteTx(t, chainID.Uint64())
	defer data.stateSqlDB.Close() //nolint:gosec,errcheck
	defer data.poolSqlDB.Close()  //nolint:gosec,errcheck

	b := make([]byte, cfg.MaxTxDataBytesSize-20)
	to := common.HexToAddress(senderAddress)
	gasPrice := big.NewInt(1000000000)
	gasLimitForThisTx := uint64(21000) + uint64(16)*uint64(len(b))
	tx := ethTypes.NewTransaction(0, to, big.NewInt(0), gasLimitForThisTx, gasPrice, b)

	// GetAuth configures and returns an auth object.
	auth, err := operations.GetAuth(senderPrivateKey, chainID.Uint64())
	require.NoError(t, err)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = data.pool.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)
}

func Test_EGPValidateEffectiveGasPrice(t *testing.T) {
	tests := []struct {
		name                string
		egpEnabled          bool
		gasPriceTx          *big.Int
		preExecutionGasUsed uint64
		gasPrices           pool.GasPrices
		expectedError       error
	}{
		{
			name:                "Reject transaction if below break-even and below current estimated L2 gas price",
			egpEnabled:          true,
			gasPriceTx:          big.NewInt(1000000000),
			preExecutionGasUsed: uint64(21000) * 2000,
			gasPrices: pool.GasPrices{
				L1GasPrice: uint64(1000000000000),
				L2GasPrice: uint64(1000000000 + 1),
			},
			expectedError: pool.ErrEffectiveGasPriceGasPriceTooLow,
		},
		{
			name:                "Accept transaction if below break-even and below current estimated L2 gas price if EGP is disabled",
			egpEnabled:          false,
			gasPriceTx:          big.NewInt(1000000000),
			preExecutionGasUsed: uint64(21000) * 2000,
			gasPrices: pool.GasPrices{
				L1GasPrice: uint64(1000000000000),
				L2GasPrice: uint64(1000000000 + 1),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.EffectiveGasPrice.Enabled = tt.egpEnabled
			data := prepareToExecuteTx(t, chainID.Uint64())
			dataLen := cfg.MaxTxDataBytesSize - 20
			signedTx := createSignedTx(t, dataLen, tt.gasPriceTx, uint64(21000)+uint64(16)*uint64(dataLen))

			err := data.pool.ValidateBreakEvenGasPrice(context.Background(), *signedTx, tt.preExecutionGasUsed, tt.gasPrices)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func Test_AddTx(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)

	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	const chainID = 2576980377
	p := setupPool(t, cfg, bc, s, st, chainID, ctx, eventLog)

	tx := new(ethTypes.Transaction)
	expectedTxEncoded := "0xf86880843b9aca008252089400000000000000000000000000000000000000008080850133333355a03ee24709870c8dbc67884c9c8acb864c1aceaaa7332b9a3db0d7a5d7c68eb8e4a0302980b070f5e3ffca3dc27b07daf69d66ab27d4df648e0b3ed059cf23aa168d"
	b, err := hex.DecodeHex(expectedTxEncoded)
	require.NoError(t, err)
	tx.UnmarshalBinary(b) //nolint:gosec,errcheck

	err = p.AddTx(ctx, *tx, ip)
	require.NoError(t, err)

	rows, err := poolSqlDB.Query(ctx, "SELECT hash, encoded, decoded, status, used_steps FROM pool.transaction")
	require.NoError(t, err)
	defer rows.Close() // nolint:staticcheck

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, status string
		var usedSteps int
		err := rows.Scan(&hash, &encoded, &decoded, &status, &usedSteps)
		require.NoError(t, err)
		b, _ := tx.MarshalJSON()

		assert.Equal(t, "0x3c499a6308dbf4e67bd4e949b0b609e3a0a5a7fd6a497acb23e37ae7f0a923cc", hash, "invalid hash")
		assert.Equal(t, expectedTxEncoded, encoded, "invalid encoded")
		assert.JSONEq(t, string(b), decoded, "invalid decoded")
		assert.Equal(t, string(pool.TxStatusPending), status, "invalid tx status")
		assert.Greater(t, usedSteps, 0, "invalid used steps")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}

func Test_AddTx_OversizedData(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	const chainID = 2576980377
	p := pool.NewPool(cfg, bc, s, st, chainID, eventLog)

	b := make([]byte, cfg.MaxTxBytesSize+1)
	to := common.HexToAddress(operations.DefaultSequencerAddress)
	tx := ethTypes.NewTransaction(0, to, big.NewInt(0), gasLimit, big.NewInt(0), b)

	// GetAuth configures and returns an auth object.
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, chainID)
	require.NoError(t, err)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.EqualError(t, err, pool.ErrOversizedData.Error())
}

func Test_AddPreEIP155Tx(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
			{
				Address: "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "200000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	const chainID = 2576980377
	p := setupPool(t, cfg, bc, s, st, chainID, ctx, eventLog)

	batchL2Data := "0xe580843b9aca00830186a0941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77c6b39bdc5f8e458aba689f2a1ff8c543a94e4817bda40f3fe34080c4ab26c1e3c2fc2cda93bc32f0a79940501fd505dcf48d94abfde932ebf1417f502cb0d9de81bff"
	b, err := hex.DecodeHex(batchL2Data)
	require.NoError(t, err)
	txs, _, _, err := state.DecodeTxs(b, forkID6)
	require.NoError(t, err)

	tx := txs[0]

	err = p.AddTx(ctx, tx, ip)
	require.NoError(t, err)

	rows, err := poolSqlDB.Query(ctx, "SELECT hash, encoded, decoded, status FROM pool.transaction")
	require.NoError(t, err)
	defer rows.Close() // nolint:staticcheck

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, status string
		err := rows.Scan(&hash, &encoded, &decoded, &status)
		require.NoError(t, err)

		b, err := tx.MarshalBinary()
		require.NoError(t, err)

		bJSON, err := tx.MarshalJSON()
		require.NoError(t, err)

		assert.Equal(t, tx.Hash().String(), hash, "invalid hash")
		assert.Equal(t, hex.EncodeToHex(b), encoded, "invalid encoded")
		assert.JSONEq(t, string(bJSON), decoded, "invalid decoded")
		assert.Equal(t, string(pool.TxStatusPending), status, "invalid tx status")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}

func Test_GetPendingTxs(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	const txsCount = 10
	const limit = 5

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := ethTypes.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
	}

	txs, err := p.GetPendingTxs(ctx, limit)
	require.NoError(t, err)

	assert.Equal(t, limit, len(txs))

	for i := 0; i < txsCount; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_GetPendingTxsZeroPassed(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	const txsCount = 10
	const limit = 0

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := ethTypes.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
	}

	txs, err := p.GetPendingTxs(ctx, limit)
	require.NoError(t, err)

	assert.Equal(t, txsCount, len(txs))

	for i := 0; i < txsCount; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_GetTopPendingTxByProfitabilityAndZkCounters(t *testing.T) {
	ctx := context.Background()
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close()

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := ethTypes.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, big.NewInt(gasPrice.Int64()+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
	}

	txs, err := p.GetTxs(ctx, pool.TxStatusPending, 1, 10)
	require.NoError(t, err)
	// bcs it's sorted by nonce, tx with the lowest nonce is expected here
	assert.Equal(t, txs[0].Transaction.Nonce(), uint64(0))
}

func Test_UpdateTxsStatus(t *testing.T) {
	ctx := context.Background()

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := ethTypes.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx1, ip)
	require.NoError(t, err)

	tx2 := ethTypes.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx2, ip)
	require.NoError(t, err)

	expectedFailedReason := "failed"
	newStatus := pool.TxStatusInvalid
	err = p.UpdateTxsStatus(ctx, []pool.TxStatusUpdateInfo{
		{
			Hash:         signedTx1.Hash(),
			NewStatus:    newStatus,
			IsWIP:        false,
			FailedReason: &expectedFailedReason,
		},
		{
			Hash:         signedTx2.Hash(),
			NewStatus:    newStatus,
			IsWIP:        false,
			FailedReason: &expectedFailedReason,
		},
	})
	if err != nil {
		t.Error(err)
	}

	var count int
	rows, err := poolSqlDB.Query(ctx, "SELECT status, failed_reason FROM pool.transaction WHERE hash = ANY($1)", []string{signedTx1.Hash().String(), signedTx2.Hash().String()})
	if err != nil {
		t.Error(err)
	}
	defer rows.Close() // nolint:staticcheck
	var state, failedReason string
	for rows.Next() {
		count++
		if err := rows.Scan(&state, &failedReason); err != nil {
			t.Error(err)
		}
	}
	assert.Equal(t, 2, count)
}

func Test_UpdateTxStatus(t *testing.T) {
	ctx := context.Background()

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	tx := ethTypes.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx, ip); err != nil {
		t.Error(err)
	}
	expectedFailedReason := "failed"
	err = p.UpdateTxStatus(ctx, signedTx.Hash(), pool.TxStatusInvalid, false, &expectedFailedReason)
	if err != nil {
		t.Error(err)
	}

	rows, err := poolSqlDB.Query(ctx, "SELECT status, failed_reason FROM pool.transaction WHERE hash = $1", signedTx.Hash().Hex())
	require.NoError(t, err)

	defer rows.Close() // nolint:staticcheck
	var state, failedReason string
	rows.Next()
	if err := rows.Scan(&state, &failedReason); err != nil {
		t.Error(err)
	}

	assert.Equal(t, pool.TxStatusInvalid, pool.TxStatus(state))
	assert.Equal(t, expectedFailedReason, failedReason)
}

func Test_SetAndGetGasPrice(t *testing.T) {
	initOrResetDB(t)

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(t, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	p := pool.NewPool(cfg, bc, s, nil, chainID.Uint64(), eventLog)

	nBig, err := rand.Int(rand.Reader, big.NewInt(0).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	expectedGasPrice := pool.GasPrices{nBig.Uint64(), nBig.Uint64()}
	ctx := context.Background()
	err = p.SetGasPrices(ctx, expectedGasPrice.L2GasPrice, expectedGasPrice.L1GasPrice)
	require.NoError(t, err)

	gasPrice, err := p.GetGasPrices(ctx)
	require.NoError(t, err)

	assert.Equal(t, expectedGasPrice, gasPrice)
}

func TestDeleteGasPricesHistoryOlderThan(t *testing.T) {
	initOrResetDB(t)

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(t, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	p := pool.NewPool(cfg, bc, s, nil, chainID.Uint64(), eventLog)

	ctx := context.Background()

	// set first gas price
	expectedL2GasPrice1 := uint64(1)
	expectedL1GasPrice1 := expectedL2GasPrice1 * 2
	err = p.SetGasPrices(ctx, expectedL2GasPrice1, expectedL1GasPrice1)
	require.NoError(t, err)
	gasPrices, err := p.GetGasPrices(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedL2GasPrice1, gasPrices.L2GasPrice)
	assert.Equal(t, expectedL1GasPrice1, gasPrices.L1GasPrice)

	// set second gas price
	expectedL2GasPrice2 := uint64(2)
	expectedL1GasPrice2 := uint64(2) * 2
	err = p.SetGasPrices(ctx, expectedL2GasPrice2, expectedL1GasPrice2)
	require.NoError(t, err)
	gasPrices, err = p.GetGasPrices(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedL2GasPrice2, gasPrices.L2GasPrice)
	assert.Equal(t, expectedL1GasPrice2, gasPrices.L1GasPrice)

	// min gas price should be the first one
	date := time.Now().UTC().Add(-time.Second * 2)
	min, err := p.MinL2GasPriceSince(ctx, date)
	require.NoError(t, err)
	require.Equal(t, expectedL2GasPrice1, min)

	// deleting the gas price history should keep at least the last one gas price (the second one)
	err = p.DeleteGasPricesHistoryOlderThan(ctx, time.Now().UTC().Add(time.Second))
	require.NoError(t, err)

	min, err = p.MinL2GasPriceSince(ctx, date)
	require.NoError(t, err)
	require.Equal(t, expectedL2GasPrice2, min)
}

func TestGetPendingTxSince(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	txsAddedHashes := []common.Hash{}
	txsAddedTime := []time.Time{}

	timeBeforeTxs := time.Now()
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := ethTypes.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		txsAddedTime = append(txsAddedTime, time.Now())
		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
		txsAddedHashes = append(txsAddedHashes, signedTx.Hash())
		time.Sleep(1 * time.Second)
	}

	txHashes, err := p.GetPendingTxHashesSince(ctx, timeBeforeTxs)
	require.NoError(t, err)
	assert.Equal(t, txsCount, len(txHashes))
	for i, txHash := range txHashes {
		assert.Equal(t, txHash.Hex(), txsAddedHashes[i].Hex())
	}

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[5])
	require.NoError(t, err)
	assert.Equal(t, 5, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[5].Hex())
	assert.Equal(t, txHashes[1].Hex(), txsAddedHashes[6].Hex())
	assert.Equal(t, txHashes[2].Hex(), txsAddedHashes[7].Hex())
	assert.Equal(t, txHashes[3].Hex(), txsAddedHashes[8].Hex())
	assert.Equal(t, txHashes[4].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[8])
	require.NoError(t, err)
	assert.Equal(t, 2, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[8].Hex())
	assert.Equal(t, txHashes[1].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[9])
	require.NoError(t, err)
	assert.Equal(t, 1, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[9].Add(1*time.Second))
	require.NoError(t, err)
	assert.Equal(t, 0, len(txHashes))
}

func Test_DeleteTransactionsByHashes(t *testing.T) {
	ctx := context.Background()
	initOrResetDB(t)
	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := ethTypes.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx1, ip)
	require.NoError(t, err)

	tx2 := ethTypes.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx2, ip)
	require.NoError(t, err)

	err = p.DeleteTransactionsByHashes(ctx, []common.Hash{signedTx1.Hash(), signedTx2.Hash()})
	require.NoError(t, err)

	var count int
	err = poolSqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.transaction").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func Test_TryAddIncompatibleTxs(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	initialBalance, _ := big.NewInt(0).SetString(encoding.MaxUint256StrNumber, encoding.Base10)
	initialBalance = initialBalance.Add(initialBalance, initialBalance)
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: operations.DefaultSequencerAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   initialBalance.String(),
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	type testCase struct {
		name                 string
		createIncompatibleTx func() ethTypes.Transaction
		expectedError        error
	}

	auth := operations.MustGetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(operations.DefaultSequencerPrivateKey, "0x"))
	require.NoError(t, err)

	chainIdOver64Bits := big.NewInt(0).SetUint64(math.MaxUint64)
	chainIdOver64Bits = chainIdOver64Bits.Add(chainIdOver64Bits, big.NewInt(2))
	authChainIdOver64Bits, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdOver64Bits)
	require.NoError(t, err)

	bigIntOver256Bits, _ := big.NewInt(0).SetString(encoding.MaxUint256StrNumber, encoding.Base10)
	bigIntOver256Bits = bigIntOver256Bits.Add(bigIntOver256Bits, big.NewInt(2))

	testCases := []testCase{
		{
			name: "Gas price over 256 bits",
			createIncompatibleTx: func() ethTypes.Transaction {
				tx := ethTypes.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), gasLimit, bigIntOver256Bits, nil)
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: pool.ErrInsufficientFunds,
		},
		{
			name: "Value over 256 bits",
			createIncompatibleTx: func() ethTypes.Transaction {
				tx := ethTypes.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					bigIntOver256Bits, gasLimit, gasPrice, nil)
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: pool.ErrInsufficientFunds,
		},
		{
			name: "data over 30k bytes",
			createIncompatibleTx: func() ethTypes.Transaction {
				data := [30001]byte{}
				tx := ethTypes.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), 141004, gasPrice, data[:])
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: fmt.Errorf("data size bigger than allowed, current size is %v bytes and max allowed is %v bytes", 30001, 30000),
		},
		{
			name: "chain id over 64 bits",
			createIncompatibleTx: func() ethTypes.Transaction {
				tx := ethTypes.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), gasLimit, gasPrice, nil)
				signedTx, err := authChainIdOver64Bits.Signer(authChainIdOver64Bits.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: fmt.Errorf("chain id higher than allowed, max allowed is %v", uint64(math.MaxUint64)),
		},
	}
	c := cfg
	c.TxFeeCap = 0
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			incompatibleTx := testCase.createIncompatibleTx()
			p := setupPool(t, c, bc, s, st, incompatibleTx.ChainId().Uint64(), ctx, eventLog)
			err = p.AddTx(ctx, incompatibleTx, ip)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func newState(sqlDB *pgxpool.Pool, eventLog *event.EventLog) *state.State {
	ctx := context.Background()
	stCfg := state.Config{MaxCumulativeGasUsed: 800000, ChainID: chainID.Uint64(), ForkIDIntervals: []state.ForkIDInterval{{
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
		ForkId:          5,
		Version:         "",
	}}}

	stateDb := pgstatestorage.NewPostgresStorage(stCfg, sqlDB)
	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}
	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	executorClient, _, _ := executor.NewExecutorClient(ctx, executorServerConfig)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	stateTree := merkletree.NewStateTree(stateDBClient)

	st := state.NewState(stCfg, stateDb, executorClient, stateTree, eventLog, nil, nil)
	return st
}

func initOrResetDB(t *testing.T) {
	err := dbutils.InitOrResetState(stateDBCfg)
	require.NoError(t, err)

	err = dbutils.InitOrResetPool(poolDBCfg)
	require.NoError(t, err)
}

func Test_AddTxWithIntrinsicGasTooLow(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)
	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert transaction
	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(0),
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(0),
		To:       nil,
		Value:    big.NewInt(10),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(0),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21000),
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)

	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21223),
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21224),
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)

	txs, err := p.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(txs))

	for i := 0; i < 2; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_AddTx_GasPriceErr(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	auth.NoSend = true
	auth.GasLimit = 53000
	auth.GasPrice = big.NewInt(0)
	auth.Nonce = big.NewInt(0)

	require.NoError(t, err)
	testCases := []struct {
		name          string
		nonce         uint64
		to            *common.Address
		gasLimit      uint64
		gasPrice      *big.Int
		data          []byte
		expectedError error
	}{
		{
			name:          "GasPriceTooLowErr",
			nonce:         0,
			to:            nil,
			gasLimit:      gasLimit, // Is a contract 53000
			gasPrice:      big.NewInt(0).SetUint64(gasPrice.Uint64() - uint64(1)),
			data:          []byte{},
			expectedError: pool.ErrGasPrice,
		},
	}

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initOrResetDB(t)

			stateSqlDB, err := db.NewSQLDB(stateDBCfg)
			if err != nil {
				panic(err)
			}
			defer stateSqlDB.Close() //nolint:gosec,errcheck

			poolSqlDB, err := db.NewSQLDB(poolDBCfg)
			require.NoError(t, err)
			defer poolSqlDB.Close() //nolint:gosec,errcheck

			st := newState(stateSqlDB, eventLog)

			genesisBlock := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}
			genesis := state.Genesis{
				Actions: []*state.GenesisAction{
					{
						Address: senderAddress,
						Type:    int(merkletree.LeafTypeBalance),
						Value:   "1000000000000000000000",
					},
				},
			}
			ctx := context.Background()
			dbTx, err := st.BeginStateTransaction(ctx)
			require.NoError(t, err)
			_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

			s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
			require.NoError(t, err)

			p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)
			tx := ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce:    tc.nonce,
				To:       tc.to,
				Value:    big.NewInt(0),
				Gas:      tc.gasLimit,
				GasPrice: tc.gasPrice,
				Data:     tc.data,
			})
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
			require.NoError(t, err)

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainID.Uint64())))
			require.NoError(t, err)

			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			err = p.AddTx(ctx, *signedTx, ip)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_AddRevertedTx(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert transaction
	revertScData, err := hex.DecodeHex(Revert.RevertBin)
	require.NoError(t, err)
	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    uint64(0),
		Gas:      uint64(1000000),
		GasPrice: gasPrice,
		Data:     revertScData,
	})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)

	txs, err := p.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(txs))

	for i := 0; i < 1; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_BlockedAddress(t *testing.T) {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	auth := operations.MustGetAuth(operations.DefaultSequencerPrivateKey, chainID.Uint64())

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: auth.From.String(),
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)

	require.NoError(t, err)

	cfg := pool.Config{
		MaxTxBytesSize:                    30132,
		MaxTxDataBytesSize:                30000,
		MinAllowedGasPriceInterval:        cfgTypes.NewDuration(5 * time.Minute),
		PollMinAllowedGasPriceInterval:    cfgTypes.NewDuration(15 * time.Second),
		DefaultMinGasPriceAllowed:         1000000000,
		IntervalToRefreshBlockedAddresses: cfgTypes.NewDuration(5 * time.Second),
		IntervalToRefreshGasPrices:        cfgTypes.NewDuration(5 * time.Second),
		AccountQueue:                      64,
		GlobalQueue:                       1024,
	}

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)
	p.StartRefreshingBlockedAddressesPeriodically()

	gasPrices, err := p.GetGasPrices(ctx)
	require.NoError(t, err)

	// Add tx while address is not blocked
	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    0,
		GasPrice: big.NewInt(0).SetInt64(int64(gasPrices.L2GasPrice)),
		Gas:      24000,
		To:       &auth.From,
		Value:    big.NewInt(1000),
	})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)

	// block address
	_, err = poolSqlDB.Exec(ctx, "INSERT INTO pool.blocked(addr) VALUES($1)", auth.From.String())
	require.NoError(t, err)

	// wait it to refresh
	time.Sleep(cfg.IntervalToRefreshBlockedAddresses.Duration)

	// get blocked when try to add new tx
	tx = ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(0).SetInt64(int64(gasPrices.L2GasPrice)),
		Gas:      24000,
		To:       &auth.From,
		Value:    big.NewInt(1000),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.Equal(t, pool.ErrBlockedSender, err)

	// remove block
	_, err = poolSqlDB.Exec(ctx, "DELETE FROM pool.blocked WHERE addr = $1", auth.From.String())
	require.NoError(t, err)

	// wait it to refresh
	time.Sleep(cfg.IntervalToRefreshBlockedAddresses.Duration)

	// allowed to add tx again
	err = p.AddTx(ctx, *signedTx, ip)
	require.NoError(t, err)
}

/*
func Test_AddTx_GasOverBatchLimit(t *testing.T) {
	testCases := []struct {
		name          string
		nonce         uint64
		to            *common.Address
		value         *big.Int
		gasLimit      uint64
		gasPrice      *big.Int
		data          []byte
		expectedError error
	}{
		{
			name:          "Gas over batch limit",
			nonce:         0,
			to:            nil,
			value:         big.NewInt(0),
			gasLimit:      uint64(30000001),
			gasPrice:      big.NewInt(1000000000000),
			data:          []byte{},
			expectedError: pool.ErrGasLimit,
		},
	}

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initOrResetDB(t)

			stateSqlDB, err := db.NewSQLDB(stateDBCfg)
			if err != nil {
				panic(err)
			}
			defer stateSqlDB.Close() //nolint:gosec,errcheck

			poolSqlDB, err := db.NewSQLDB(poolDBCfg)
			require.NoError(t, err)
			defer poolSqlDB.Close() //nolint:gosec,errcheck

			st := newState(stateSqlDB, eventLog)

			genesisBlock := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}
			genesis := state.Genesis{
				GenesisActions: []*state.GenesisAction{
					{
						Address: senderAddress,
						Type:    int(merkletree.LeafTypeBalance),
						Value:   "1000000000000000000000",
					},
				},
				FirstBatchData: genesis.FirstBatchData,
			}
			ctx := context.Background()
			dbTx, err := st.BeginStateTransaction(ctx)
			require.NoError(t, err)
			genesis.FirstBatchData.Timestamp = uint64(time.Now().Unix())
			_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

			s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
			require.NoError(t, err)

			p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)
			tx := ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce:    tc.nonce,
				To:       tc.to,
				Value:    tc.value,
				Gas:      tc.gasLimit,
				GasPrice: tc.gasPrice,
				Data:     tc.data,
			})
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
			require.NoError(t, err)

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)

			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			err = p.AddTx(ctx, *signedTx, ip)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
*/

func Test_AddTx_AccountQueueLimit(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	nonce := uint64(0)
	for nonce < cfg.AccountQueue {
		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    nonce,
			Value:    big.NewInt(0),
			Gas:      uint64(1000000),
			GasPrice: gasPrice,
		})

		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)

		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
		nonce++
	}

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    nonce,
		Value:    big.NewInt(0),
		Gas:      uint64(1000000),
		GasPrice: gasPrice,
	})

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err, pool.ErrNonceTooHigh)
}

func Test_AddTx_GlobalQueueLimit(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB, eventLog)

	// generate accounts
	accounts := map[common.Address]*ecdsa.PrivateKey{}
	genesisActions := []*state.GenesisAction{
		{
			Address: senderAddress,
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "1000000000000000000000",
		},
	}
	for i := 0; i < int(cfg.GlobalQueue); i++ {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		publicKey := privateKey.Public().(*ecdsa.PublicKey)
		fromAddress := crypto.PubkeyToAddress(*publicKey)
		accounts[fromAddress] = privateKey
		genesisActions = append(genesisActions, &state.GenesisAction{
			Address: fromAddress.String(),
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "1000000000000000000000",
		})
	}

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: genesisActions,
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	for _, privateKey := range accounts {
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		require.NoError(t, err)
		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    0,
			Value:    big.NewInt(0),
			Gas:      uint64(1000000),
			GasPrice: gasPrice,
		})

		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)

		err = p.AddTx(ctx, *signedTx, ip)
		require.NoError(t, err)
	}

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    0,
		Value:    big.NewInt(0),
		Gas:      uint64(1000000),
		GasPrice: big.NewInt(1),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err, pool.ErrTxPoolOverflow)
}

func Test_AddTx_NonceTooHigh(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB, eventLog)

	// generate accounts
	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    cfg.AccountQueue,
		Value:    big.NewInt(0),
		Gas:      uint64(1000000),
		GasPrice: big.NewInt(1),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, ip)
	require.Error(t, err, pool.ErrNonceTooHigh)
}

func Test_AddTx_IPValidation(t *testing.T) {
	var tests = []struct {
		name     string
		ip       string
		expected error
	}{
		{"Valid IPv4", "127.0.0.1", nil},
		{"Valid IPv6", "2001:db8:0:1:1:1:1:1", nil},
		{"Invalid IP", "300.0.0.1", pool.ErrInvalidIP},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			initOrResetDB(t)

			stateSqlDB, err := db.NewSQLDB(stateDBCfg)
			require.NoError(t, err)
			defer stateSqlDB.Close() //nolint:gosec,errcheck

			poolSqlDB, err := db.NewSQLDB(poolDBCfg)
			require.NoError(t, err)

			defer poolSqlDB.Close() //nolint:gosec,errcheck

			eventStorage, err := nileventstorage.NewNilEventStorage()
			if err != nil {
				log.Fatal(err)
			}
			eventLog := event.NewEventLog(event.Config{}, eventStorage)

			st := newState(stateSqlDB, eventLog)

			genesisBlock := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}
			ctx := context.Background()
			dbTx, err := st.BeginStateTransaction(ctx)
			require.NoError(t, err)
			_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

			s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
			require.NoError(t, err)

			const chainID = 2576980377
			p := setupPool(t, cfg, bc, s, st, chainID, ctx, eventLog)

			tx := new(ethTypes.Transaction)
			expectedTxEncoded := "0xf86880843b9aca008252089400000000000000000000000000000000000000008080850133333355a03ee24709870c8dbc67884c9c8acb864c1aceaaa7332b9a3db0d7a5d7c68eb8e4a0302980b070f5e3ffca3dc27b07daf69d66ab27d4df648e0b3ed059cf23aa168d"
			b, err := hex.DecodeHex(expectedTxEncoded)
			require.NoError(t, err)
			tx.UnmarshalBinary(b) //nolint:gosec,errcheck

			err = p.AddTx(context.Background(), *tx, tc.ip)

			if tc.expected != nil {
				assert.ErrorIs(t, err, tc.expected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AddTx_TxFeeCap(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainID.Uint64(), ctx, eventLog)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	type testCase struct {
		name          string
		nonce         uint64
		gas           uint64
		gasPrice      string
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "add tx with fee under cap",
			nonce:         0,
			gas:           uint64(100000),
			gasPrice:      "9999999999999",
			expectedError: nil,
		},
		{
			name:          "add tx with fee exactly as cap",
			nonce:         0,
			gas:           uint64(100000),
			gasPrice:      "10000000000000",
			expectedError: nil,
		},
		{
			name:          "add tx with fee over the cap",
			nonce:         0,
			gas:           uint64(100000),
			gasPrice:      "10000000000001",
			expectedError: fmt.Errorf("tx fee (1.0000000000001 ether) exceeds the configured cap (1.00 ether)"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gasPrice, ok := big.NewInt(0).SetString(tc.gasPrice, encoding.Base10)
			require.True(t, ok)
			tx := ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce:    tc.nonce,
				Gas:      tc.gas,
				GasPrice: gasPrice,
			})

			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			err = p.AddTx(ctx, *signedTx, ip)
			if tc.expectedError != nil {
				require.Equal(t, err.Error(), tc.expectedError.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func setupPool(t *testing.T, cfg pool.Config, constraintsCfg state.BatchConstraintsCfg, s *pgpoolstorage.PostgresPoolStorage, st *state.State, chainID uint64, ctx context.Context, eventLog *event.EventLog) *pool.Pool {
	err := s.SetGasPrices(ctx, gasPrice.Uint64(), l1GasPrice.Uint64())
	require.NoError(t, err)
	p := pool.NewPool(cfg, constraintsCfg, s, st, chainID, eventLog)
	p.StartPollingMinSuggestedGasPrice(ctx)
	return p
}

func prepareToExecuteTx(t *testing.T, chainIDToCreate uint64) testData {
	initOrResetDB(t)

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	//defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	require.NoError(t, err)

	//defer poolSqlDB.Close() //nolint:gosec,errcheck

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		stateSqlDB.Close() //nolint:gosec,errcheck
		poolSqlDB.Close()  //nolint:gosec,errcheck
		log.Fatal(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	st := newState(stateSqlDB, eventLog)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	require.NoError(t, err)

	p := setupPool(t, cfg, bc, s, st, chainIDToCreate, ctx, eventLog)
	return testData{
		pool:       p,
		st:         st,
		stateSqlDB: stateSqlDB,
		poolSqlDB:  poolSqlDB,
	}
}

func createSignedTx(t *testing.T, dataLen int, gasPrice *big.Int, gasLimit uint64) *ethTypes.Transaction {
	b := make([]byte, cfg.MaxTxDataBytesSize-20)
	to := common.HexToAddress(senderAddress)
	//gasPrice := big.NewInt(1000000000)
	//gasLimitForThisTx := uint64(21000) + uint64(16)*uint64(len(b))
	tx := ethTypes.NewTransaction(0, to, big.NewInt(0), gasLimit, gasPrice, b)
	auth, err := operations.GetAuth(senderPrivateKey, chainID.Uint64())
	require.NoError(t, err)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	return signedTx
}
