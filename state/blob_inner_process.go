package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// ProcessBlobInner processes a blobInner and returns the splitted batches
func (s *State) ProcessBlobInner(ctx context.Context, request ProcessBlobInnerProcessRequest, data []byte) (*ProcessBlobInnerResponse, error) {
	requestExecutor := convertBlobInnerProcessRequestToExecutor(request, data)
	processResponse, err := s.executorClient.ProcessBlobInnerV3(ctx, requestExecutor)
	if err != nil {
		log.Errorf("Error processing blobInner: %v", err)
		return nil, err
	}
	return newProcessBlobInnerProcessResponse(processResponse), nil
}
