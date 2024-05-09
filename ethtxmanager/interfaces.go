package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type ethermanInterface interface {
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) (bool, error)
	SendTx(ctx context.Context, tx *types.Transaction) error
	PendingNonce(ctx context.Context, account common.Address) (uint64, error)
	CurrentNonce(ctx context.Context, account common.Address) (uint64, error)
	SuggestedGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error)
	CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error)
	SignTx(ctx context.Context, sender common.Address, tx *types.Transaction) (*types.Transaction, error)
	GetRevertMessage(ctx context.Context, tx *types.Transaction) (string, error)
}

type storageInterface interface {
	Add(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error
	Get(ctx context.Context, owner, id string, dbTx pgx.Tx) (monitoredTx, error)
	GetByStatus(ctx context.Context, owner *string, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]monitoredTx, error)
	GetBySenderAndStatus(ctx context.Context, sender common.Address, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]monitoredTx, error)
	GetByBlock(ctx context.Context, fromBlock, toBlock *uint64, dbTx pgx.Tx) ([]monitoredTx, error)
	Update(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error
}

type stateInterface interface {
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
}
