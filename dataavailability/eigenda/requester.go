package eigenda

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// DataAvailabilityRetriever is the EigenDA retriever client to retrieve and
// decode EigenDA blobs into L2 batches data
type DataAvailabilityRetriever struct {
	backend BlobRetriever
}

// GetBatchL2DataFromRequestId gets the batch data from the EigenDA request ID
func (d *DataAvailabilityRetriever) GetBatchL2DataFromRequestId(ctx context.Context, id []byte) ([][]byte, error) {
	msg, err := d.backend.GetDataAvailabilityMessageFromId(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return d.backend.GetBatchL2Data([]uint64{}, []common.Hash{}, msg)
	}
}

// GetDataAvailabilityMessageFromRequestId gets the data availability message
// from the EigenDA request ID
func (d *DataAvailabilityRetriever) GetDataAvailabilityMessageFromRequestId(ctx context.Context, id []byte) ([]byte, error) {
	return d.backend.GetDataAvailabilityMessageFromId(ctx, id)
}
