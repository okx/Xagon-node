package nubit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBlobData(t *testing.T) {
	data := BlobData{
		BlobID:    []byte{10},
		Signature: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
	}
	msg, err := TryEncodeToDataAvailabilityMessage(data)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg)
}

func TestEncodeDecodeBlobData(t *testing.T) {
	data := BlobData{
		BlobID:    []byte{10},
		Signature: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
	}
	msg, err := TryEncodeToDataAvailabilityMessage(data)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg)

	// Check blob ID
	decoded_data, err := TryDecodeFromDataAvailabilityMessage(msg)
	assert.NoError(t, err)
	assert.Equal(t, data.BlobID, decoded_data.BlobID)
	assert.Equal(t, data.Signature, decoded_data.Signature)
}
