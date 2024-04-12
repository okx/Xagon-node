package sequencer

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error
	DeleteFailedTransactionsOlderThan(ctx context.Context, date time.Time) error
	DeleteTransactionByHash(ctx context.Context, hash common.Hash) error
	MarkWIPTxsAsPending(ctx context.Context) error
	GetNonWIPPendingTxs(ctx context.Context, limit uint64) ([]pool.Transaction, error)
	UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool, failedReason *string) error
	GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, *state.ZKCounters, error)
	UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error
	GetGasPrices(ctx context.Context) (pool.GasPrices, error)
	GetDefaultMinGasPriceAllowed() uint64
	GetL1AndL2GasPrice() (uint64, uint64)
	GetEarliestProcessedTx(ctx context.Context) (common.Hash, error)
	CountPendingTransactions(ctx context.Context) (uint64, error)
	GetLevelByAddr(addr common.Address) int
	CountTransactionsByFromAndStatus(ctx context.Context, from common.Address, status ...pool.TxStatus) (uint64, error)
}

type workerInterface interface {
	GetBestFittingTx(resources state.BatchResources) (*TxTracker, error)
	UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.InfoReadWrite) []*TxTracker
	UpdateTxZKCounters(txHash common.Hash, from common.Address, usedZKCounters state.ZKCounters, reservedZKCounters state.ZKCounters)
	AddTxTracker(ctx context.Context, txTracker *TxTracker, nonce uint64) (replacedTx *TxTracker, dropReason error)
	MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int) []*TxTracker
	DeleteTx(txHash common.Hash, from common.Address)
	AddPendingTxToStore(txHash common.Hash, addr common.Address)
	DeletePendingTxToStore(txHash common.Hash, addr common.Address)
	NewTxTracker(tx pool.Transaction, usedZKcounters state.ZKCounters, reservedZKCouners state.ZKCounters, ip string) (*TxTracker, error)
	AddForcedTx(txHash common.Hash, addr common.Address)
	DeleteForcedTx(txHash common.Hash, addr common.Address)
}
