package jsonrpc

import (
	"encoding/json"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

type ban func(request *types.Request) bool
type banList []ban

var banListInst banList

func init() {
	banListInst.register(debugTraceTransactionNotPrestateTracer)
}

func pass(request *types.Request) bool {
	return !banListInst.ban(request)
}

func (c *banList) register(handler ban) {
	*c = append(*c, handler)
}

func (c *banList) ban(request *types.Request) bool {
	for _, handler := range *c {
		if handler(request) {
			return true
		}
	}

	return false
}

// debutTraceTransactionNotPrestateTracer checks if the request is a debug_traceTransaction and the tracer is not preStateTracer
func debugTraceTransactionNotPrestateTracer(req *types.Request) bool {
	if req == nil || req.Method != "debug_traceTransaction" {
		return false
	}

	// check params passed by request match function params
	var testStruct []interface{}
	if err := json.Unmarshal(req.Params, &testStruct); err != nil {
		return false
	}
	inputs := make([]interface{}, len(testStruct))
	if err := json.Unmarshal(req.Params, &inputs); err != nil {
		return false
	}

	// for debug_tranceTransaction, at least 2 param is required
	if len(inputs) < 2 { // nolint:gomnd
		return false
	}
	var cfg traceConfig
	err := json.Unmarshal(inputs[1].([]byte), &cfg) // the second param is the trace config
	if err != nil {
		return false
	}
	if cfg.Tracer != nil && *cfg.Tracer != "prestateTracer" {
		return true
	}

	return false
}
