package jsonrpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// GetBatchDataByNumbers returns the batch data for batches by numbers
func (z *ZKEVMEndpoints) GetBatchDataByNumbers(filter types.BatchFilter) (interface{}, types.Error) {
	return z.txMan.NewDbTxScope(z.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var batchNumbers []uint64
		for _, bn := range filter.Numbers {
			n, rpcErr := bn.GetNumericBatchNumber(ctx, z.state, z.etherman, dbTx)
			if rpcErr != nil {
				return nil, rpcErr
			}
			batchNumbers = append(batchNumbers, n)
		}

		batchesData, err := z.state.GetBatchL2DataByNumbers(ctx, batchNumbers, dbTx)
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
	})
}
