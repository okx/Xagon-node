package elderberry_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/elderberry"
	mock_elderberry "github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/elderberry/mocks"
)

func TestProcessorL1InfoTreeUpdate_Process(t *testing.T) {
	mockState := mock_elderberry.NewStateL1SequenceBatchesElderberry(t)
	mockPreviousProcessor := mock_elderberry.NewPreviousProcessor(t)

	_ = elderberry.NewProcessorL1SequenceBatchesElderberry(mockPreviousProcessor, mockState)
}
