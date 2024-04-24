package syncinterfaces

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// EthermanFullInterface contains the methods required to interact with ethereum.
type EthermanFullInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*ethTypes.Block, error)
	GetTrustedSequencerURL() (string, error)
	VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error)
	GetLatestVerifiedBatchNum() (uint64, error)

	EthermanGetLatestBatchNumber
	GetFinalizedBlockNumber(ctx context.Context) (uint64, error)
}

type EthermanGetLatestBatchNumber interface {
	GetLatestBatchNumber() (uint64, error)
}
