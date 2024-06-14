package eigenda

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// Set longer timeout flag for test case
func TestDisperseBlobWithStringDataUsingProvider(t *testing.T) {
	cfg := dataavailability.Config{
		Hostname:          "disperser-holesky.eigenda.xyz",
		Port:              "443",
		Timeout:           types.NewDuration(30 * time.Second),
		UseSecureGrpcFlag: true,
	}
	provider := NewDataAvailabilityProvider(cfg)

	// Generate mock string batch data
	stringData := "hihihihihihihihihihihihihihihihihihi"
	data := []byte(stringData)

	// Generate mock string sequence
	mockBatches := [][]byte{}
	for i := 0; i < 10; i++ {
		mockBatches = append(mockBatches, data)
	}

	// send mock sequence with provider
	msg, err := provider.PostSequence(context.Background(), mockBatches)
	fmt.Println("DA msg: ", msg)
	assert.NoError(t, err)
	blobData, err := TryDecodeFromDataAvailabilityMessage(msg)
	assert.NoError(t, err)
	assert.NotZero(t, blobData.BlobVerificationProof.BatchId)
	assert.NotZero(t, blobData.BlobVerificationProof.BlobIndex)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.BlobHeadersRoot.Bytes())
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.BlobHeadersRoot.Bytes())
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.SignedStakeForQuorums)
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.SignedStakeForQuorums)
	assert.NotZero(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.SignatoryRecordHash.Bytes())
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.SignatoryRecordHash.Bytes())
	assert.NotZero(t, blobData.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber)
	assert.NotNil(t, blobData.BlobVerificationProof.InclusionProof)
	assert.NotEmpty(t, blobData.BlobVerificationProof.InclusionProof)
	assert.NotNil(t, blobData.BlobVerificationProof.QuorumIndices)
	assert.NotEmpty(t, blobData.BlobVerificationProof.QuorumIndices)
	fmt.Println("Decoding DA msg successful")

	// Retrieve sequence with provider
	batchesData, err := provider.GetSequence(context.Background(), []common.Hash{}, msg)

	// Validate retrieved data
	assert.NoError(t, err)
	assert.Equal(t, 10, len(batchesData))
	for _, batchData := range batchesData {
		assert.Equal(t, stringData, string(batchData))
	}
}

// Set longer timeout flag for test case
func TestDisperseBlobWithRandomDataUsingProvider(t *testing.T) {
	cfg := dataavailability.Config{
		Hostname:          "disperser-holesky.eigenda.xyz",
		Port:              "443",
		Timeout:           types.NewDuration(30 * time.Second),
		UseSecureGrpcFlag: true,
	}
	provider := NewDataAvailabilityProvider(cfg)

	// Define Different DataSizes
	dataSize := []int{100000, 200000, 1000, 80, 30000}

	// Disperse Blob with different DataSizes
	rand.Seed(time.Now().UnixNano())                         //nolint:gosec,staticcheck
	data := make([]byte, dataSize[rand.Intn(len(dataSize))]) //nolint:gosec,staticcheck
	_, err := rand.Read(data)                                //nolint:gosec,staticcheck
	assert.NoError(t, err)

	// Generate mock string sequence
	mockBatches := [][]byte{}
	for i := 0; i < 10; i++ {
		mockBatches = append(mockBatches, data)
	}

	// send mock sequence with provider
	msg, err := provider.PostSequence(context.Background(), mockBatches)
	fmt.Println("DA msg: ", msg)
	assert.NoError(t, err)
	blobData, err := TryDecodeFromDataAvailabilityMessage(msg)
	assert.NoError(t, err)
	assert.NotZero(t, blobData.BlobVerificationProof.BatchId)
	assert.NotZero(t, blobData.BlobVerificationProof.BlobIndex)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.BlobHeadersRoot.Bytes())
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.BlobHeadersRoot.Bytes())
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.SignedStakeForQuorums)
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.SignedStakeForQuorums)
	assert.NotZero(t, blobData.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber)
	assert.NotNil(t, blobData.BlobVerificationProof.BatchMetadata.SignatoryRecordHash.Bytes())
	assert.NotEmpty(t, blobData.BlobVerificationProof.BatchMetadata.SignatoryRecordHash.Bytes())
	assert.NotZero(t, blobData.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber)
	assert.NotNil(t, blobData.BlobVerificationProof.InclusionProof)
	assert.NotEmpty(t, blobData.BlobVerificationProof.InclusionProof)
	assert.NotNil(t, blobData.BlobVerificationProof.QuorumIndices)
	assert.NotEmpty(t, blobData.BlobVerificationProof.QuorumIndices)
	fmt.Println("Decoding DA msg successful")

	// Retrieve sequence with provider
	batchesData, err := provider.GetSequence(context.Background(), []common.Hash{}, msg)

	// Validate retrieved data
	assert.NoError(t, err)
	assert.Equal(t, 10, len(batchesData))
	for idx, batchData := range batchesData {
		assert.Equal(t, mockBatches[idx], batchData)
	}
}
