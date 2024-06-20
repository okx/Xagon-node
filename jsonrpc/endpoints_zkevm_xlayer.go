package jsonrpc

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/jackc/pgx/v4"
)

// GetBatchSealTime returns the seal time
func (z *ZKEVMEndpoints) GetBatchSealTime(batchNumber types.BatchNumber) (interface{}, types.Error) {
	return z.txMan.NewDbTxScope(z.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var err error
		batchNumber, rpcErr := batchNumber.GetNumericBatchNumber(ctx, z.state, z.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sealTime, err := z.state.GetLastL2BlockTimeByBatchNumber(ctx, batchNumber, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't get last l2 block time from state batch number %v, error: %v", batchNumber, err), nil, false)
		}

		return hex.EncodeUint64(sealTime), nil
	})
}
