package jsonrpc

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// GetBatchSealTime returns the seal time
func (z *ZKEVMEndpoints) GetBatchSealTime(batchNumber types.BatchNumber) (interface{}, types.Error) {
	var err error
	ctx := context.Background()
	batchNumberResult, rpcErr := batchNumber.GetNumericBatchNumber(ctx, z.state, z.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	sealTime, err := z.state.GetLastL2BlockTimeByBatchNumber(ctx, batchNumberResult, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't get batch number %v's seal time, error: %v", batchNumber, err), nil, false)
	}

	return hex.EncodeUint64(sealTime), nil
}
