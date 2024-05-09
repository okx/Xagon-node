package syncinterfaces

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// ZKEVMClientEthereumCompatibleInterface contains the methods required to interact with zkEVM-RPC as a ethereum-API compatible
//
//	Reason behind: the zkEVMClient have some extensions to ethereum-API that are not compatible with all nodes. So if you need to maximize
//	 the compatibility the idea is to use a regular ethereum-API compatible client
type ZKEVMClientEthereumCompatibleInterface interface {
	ZKEVMClientEthereumCompatibleL2BlockGetter
}

// ZKEVMClientEthereumCompatibleL2BlockGetter contains the methods required to interact with zkEVM-RPC as a ethereum-API compatible for obtain Block information
type ZKEVMClientEthereumCompatibleL2BlockGetter interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}
