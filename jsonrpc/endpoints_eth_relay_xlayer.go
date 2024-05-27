package jsonrpc

import (
	"encoding/json"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

func (e *EthEndpoints) relayCall(method string, arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	dstURI := getRelayDestURI(e.cfg.ApiRelay.DestURI)

	res, err := client.JSONRPCCall(dstURI, method, arg, toString(blockArg))
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to relay tx to "+dstURI, err, true)
	}

	if res.Error != nil {
		if res.Error.Data == nil {
			return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
		}

		return RPCErrorResponseWithData(res.Error.Code, res.Error.Message, *res.Error.Data, nil, false)
	}
	var result types.ArgBytes
	err = json.Unmarshal(res.Result, &result)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to unmarshal result", err, true)
	}

	return result, nil
}

func toString(blockArg *types.BlockNumberOrHash) string {
	if blockArg == nil {
		return ""
	}
	if blockArg.IsHash() {
		return blockArg.Hash().Hash().String()
	}
	return blockArg.Number().StringOrHex()
}
