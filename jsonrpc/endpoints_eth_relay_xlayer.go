package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

func (e *EthEndpoints) relayCall(method string, arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	dstURI := getRelayDestURI(e.cfg.ApiRelay.DestURI)
	log.Infof("Relaying tx to %s %s %v %v", dstURI, method, arg, blockArg)

	res, err := client.JSONRPCCall(dstURI, method, arg, blockArg)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to relay tx to "+dstURI, err, true)
	}

	return res, nil
}
