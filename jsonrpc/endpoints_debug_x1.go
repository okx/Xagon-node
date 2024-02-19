package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

func (d *DebugEndpoints) buildInnerTransaction(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (interface{}, types.Error) {
	//traceCfg := defaultTraceConfig
	//tracer := "callTracer"
	//traceCfg.Tracer = &tracer
	//
	//// check tracer
	//if traceCfg.Tracer != nil && *traceCfg.Tracer != "" && !isBuiltInTracer(*traceCfg.Tracer) && !isJSCustomTracer(*traceCfg.Tracer) {
	//	return RPCErrorResponse(types.DefaultErrorCode, "invalid tracer", nil, true)
	//}
	//
	//stateTraceConfig := state.TraceConfig{
	//	DisableStack:     traceCfg.DisableStack,
	//	DisableStorage:   traceCfg.DisableStorage,
	//	EnableMemory:     traceCfg.EnableMemory,
	//	EnableReturnData: traceCfg.EnableReturnData,
	//	Tracer:           traceCfg.Tracer,
	//	TracerConfig:     traceCfg.TracerConfig,
	//}
	//result, err := d.state.DebugTransaction(ctx, hash, stateTraceConfig, dbTx)
	//if errors.Is(err, state.ErrNotFound) {
	//	return RPCErrorResponse(types.DefaultErrorCode, "transaction not found", nil, true)
	//} else if err != nil {
	//	const errorMessage = "failed to get trace"
	//	log.Errorf("%v: %v", errorMessage, err)
	//	return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	//}
	//
	//// if a tracer was specified, then return the trace result
	//if stateTraceConfig.Tracer != nil && *stateTraceConfig.Tracer != "" && len(result.ExecutorTraceResult) > 0 {
	//	return result.ExecutorTraceResult, nil
	//}
	//
	//receipt, err := d.state.GetTransactionReceipt(ctx, hash, dbTx)
	//if err != nil {
	//	const errorMessage = "failed to tx receipt"
	//	log.Errorf("%v: %v", errorMessage, err)
	//	return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	//}
	//
	//failed := receipt.Status == ethTypes.ReceiptStatusFailed
	//var returnValue interface{}
	//if stateTraceConfig.EnableReturnData {
	//	returnValue = common.Bytes2Hex(result.ReturnValue)
	//}
	//
	//structLogs := d.buildStructLogs(result.StructLogs, *traceCfg)
	//
	//resp := traceTransactionResponse{
	//	Gas:         result.GasUsed,
	//	Failed:      failed,
	//	ReturnValue: returnValue,
	//	StructLogs:  structLogs,
	//}

	//return resp, nil

	//TODO
	return nil, nil
}
