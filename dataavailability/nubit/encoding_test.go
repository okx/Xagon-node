package nubit

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeSequenceToAndFromStringBlob(t *testing.T) {
	mock_string_data := "hihihihihihihihihihihihihihihihihihi"
	data := []byte(mock_string_data)
	hash := crypto.Keccak256Hash(data)

	// Generate mock sequence data
	mockSeqData := [][]byte{}
	for i := 0; i < 10; i++ {
		mockSeqData = append(mockSeqData, data)
	}
	blob := EncodeSequence(mockSeqData)

	// Decode blob
	decodedBatchesData, decodedBatchesHash := DecodeSequence(blob)

	// Assert decoded sequence length is correct
	n_data := len(decodedBatchesData)
	n_hash := len(decodedBatchesHash)
	assert.Equal(t, 10, n_data)
	assert.Equal(t, 10, n_hash)

	// Assert decoded sequence data is correct
	for _, batchData := range decodedBatchesData {
		data_decoded := string(batchData)
		assert.Equal(t, mock_string_data, data_decoded)
	}

	// Assert decoded batches' hash is correct
	for _, batchHash := range decodedBatchesHash {
		assert.Equal(t, hash, batchHash)
	}
}

func TestEncodeDecodeSequenceToAndFromRandomBlob(t *testing.T) {
	// Define Different DataSizes
	dataSize := []int{100000, 200000, 1000, 80, 30000}

	// Generate mock sequence data
	mockSeqData := [][]byte{}
	mockSeqHash := []common.Hash{}
	for i := 0; i < 10; i++ {
		// Disperse Blob with different DataSizes
		rand.Seed(time.Now().UnixNano())                         //nolint:gosec,staticcheck
		data := make([]byte, dataSize[rand.Intn(len(dataSize))]) //nolint:gosec,staticcheck
		_, err := rand.Read(data)                                //nolint:gosec,staticcheck
		assert.NoError(t, err)
		mockSeqData = append(mockSeqData, data)

		// Get batch hash
		hash := crypto.Keccak256Hash(data)
		mockSeqHash = append(mockSeqHash, hash)
	}
	blob := EncodeSequence(mockSeqData)

	// Decode blob
	decodedBatchesData, decodedBatchesHash := DecodeSequence(blob)

	// Assert decoded sequence length is correct
	n_data := len(decodedBatchesData)
	n_hash := len(decodedBatchesHash)
	assert.Equal(t, 10, n_data)
	assert.Equal(t, 10, n_hash)

	// Assert decoded sequence data is correct
	for i := 0; i < n_data; i++ {
		assert.Equal(t, mockSeqData[i], decodedBatchesData[i])
	}

	// Assert decoded batches' hash is correct
	for i := 0; i < n_hash; i++ {
		assert.Equal(t, mockSeqHash[i], decodedBatchesHash[i])
	}
}
