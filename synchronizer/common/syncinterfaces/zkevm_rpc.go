package syncinterfaces

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
)

// ZkEVMClientInterface contains the methods required to interact with zkEVM-RPC
type ZKEVMClientTrustedBatchesGetter interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}

// ZkEVMClientInterface contains the methods required to interact with zkEVM-RPC for obtain GlobalExitRoot information
type ZKEVMClientGlobalExitRootGetter interface {
	ExitRootsByGER(ctx context.Context, globalExitRoot common.Hash) (*types.ExitRoots, error)
}

type ZKEVMClientGetL2BlockByNumber interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type ZKEVMClientInterface interface {
	ZKEVMClientTrustedBatchesGetter
	ZKEVMClientGlobalExitRootGetter
	ZKEVMClientGetL2BlockByNumber
}
