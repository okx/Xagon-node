package eigenda

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
)

// BlobRetriever is used to retrieve data availability message from EigenDA blob request ID
type BlobRetriever interface {
	dataavailability.BatchDataProvider
	GetBatchL2DataFromRequestId(ctx context.Context, id []byte) ([][]byte, error)
	GetDataAvailabilityMessageFromId(ctx context.Context, requestID []byte) ([]byte, error)
}
