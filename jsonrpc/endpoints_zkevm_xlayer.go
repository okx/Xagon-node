package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"fmt"
	"github.com/jackc/pgx/v4"
	"context"
)

// GetBatchSealTime returns the seal time
func (z *ZKEVMEndpoints) GetBatchSealTime(batchNumber types.BatchNumber) (interface{}, types.Error) {
	return z.txMan.NewDbTxScope(z.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var err error
		batchNumber, rpcErr := batchNumber.GetNumericBatchNumber(ctx, z.state, z.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sealTime, err := z.state.GetLastL2BlockCreateTimeBatchNumber(ctx, batchNumber, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't get last l2 block create time from state by batch number %v", batchNumber), err, true)
		}
		return sealTime.Unix(), nil
	})
}
