package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

func (e *EthEndpoints) relayDebugTransaction(method string, hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	dstURI := getRelayDestURI(e.cfg.ApiRelay.DestURI)

	res, err := client.JSONRPCCall(dstURI, method, hash, cfg)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to relay tx to "+dstURI, err, true)
	}

	return res, nil
}
