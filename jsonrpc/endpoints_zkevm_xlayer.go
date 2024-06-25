package jsonrpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
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

// GetBatchDataByNumbers returns the batch data for batches by numbers, XLayer
func (z *ZKEVMEndpoints) GetBatchDataByNumbers(filter types.BatchFilter) (interface{}, types.Error) {
	ctx := context.Background()

	var batchNumbers []uint64
	for _, bn := range filter.Numbers {
		n, rpcErr := bn.GetNumericBatchNumber(ctx, z.state, z.etherman, nil)
		if rpcErr != nil {
			return nil, rpcErr
		}
		batchNumbers = append(batchNumbers, n)
	}

	batchesData, err := z.state.GetBatchL2DataByNumbers(ctx, batchNumbers, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode,
			fmt.Sprintf("couldn't load batch data from state by numbers %v", filter.Numbers), err, true)
	}

	var ret []*types.BatchData
	for _, n := range batchNumbers {
		data := &types.BatchData{Number: types.ArgUint64(n)}
		if b, ok := batchesData[n]; ok {
			data.BatchL2Data = b
			data.Empty = false
		} else {
			data.Empty = true
		}
		ret = append(ret, data)
	}

	return types.BatchDataResult{Data: ret}, nil
}
