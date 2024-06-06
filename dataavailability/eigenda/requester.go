package eigenda

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type DataAvailabilityRetriever struct {
	backend BlobRetriever
}

// Get batch data from EigenDA request ID
func (d *DataAvailabilityRetriever) GetBatchL2DataFromRequestId(ctx context.Context, id []byte) ([][]byte, error) {
	msg, err := d.backend.GetDataAvailabilityMessageFromId(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return d.backend.GetBatchL2Data([]uint64{}, []common.Hash{}, msg)
	}
}

// Get data availability message from EigenDA request ID
func (d *DataAvailabilityRetriever) GetDataAvailabilityMessageFromRequestId(ctx context.Context, id []byte) ([]byte, error) {
	return d.backend.GetDataAvailabilityMessageFromId(ctx, id)
}
