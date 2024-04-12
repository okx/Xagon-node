package pool

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type storage interface {
	AddTx(ctx context.Context, tx pool.Transaction) error
	CountTransactionsByStatus(ctx context.Context, status ...pool.TxStatus) (uint64, error)
	CountTransactionsByFromAndStatus(ctx context.Context, from common.Address, status ...pool.TxStatus) (uint64, error)
	DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error
	GetGasPrices(ctx context.Context) (uint64, uint64, error)
	GetNonce(ctx context.Context, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
	GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]pool.Transaction, error)
	GetTxsByStatus(ctx context.Context, state pool.TxStatus, limit uint64) ([]pool.Transaction, error)
	GetNonWIPPendingTxs(ctx context.Context, limit uint64) ([]pool.Transaction, error)
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	DeleteFailedTransactionsOlderThan(ctx context.Context, date time.Time) error
	UpdateTxStatus(ctx context.Context, updateInfo pool.TxStatusUpdateInfo) error
	UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error
	GetTxs(ctx context.Context, filterStatus pool.TxStatus, minGasPrice, limit uint64) ([]*pool.Transaction, error)
	GetTxFromAddressFromByHash(ctx context.Context, hash common.Hash) (common.Address, uint64, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*pool.Transaction, error)
	GetTransactionByL2Hash(ctx context.Context, hash common.Hash) (*pool.Transaction, error)
	GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, *state.ZKCounters, error)
	DeleteTransactionByHash(ctx context.Context, hash common.Hash) error
	MarkWIPTxsAsPending(ctx context.Context) error
	GetEarliestProcessedTx(ctx context.Context) (common.Hash, error)
}
