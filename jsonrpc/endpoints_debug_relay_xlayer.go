package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

func (d *DebugEndpoints) relayDebugTransaction(method string, hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	dstURI := getRelayDestURI(d.cfg.ApiRelay.DestURI)

	res, err := client.JSONRPCCall(dstURI, method, hash, cfg)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to relay tx to "+dstURI, err, true)
	}
	if res.Error != nil {
		if res.Error.Data == nil {
			return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
		}

		return RPCErrorResponseWithData(res.Error.Code, res.Error.Message, *res.Error.Data, nil, false)
	}

	return res.Result, nil
}
