package eigenda

import (
	"encoding/binary"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeSequence is the helper function to encode sequence data into 1D byte array. The
// encoding scheme is ensured to be lossless.
//
// The first n+1 8-bytes of the blob contains the metadata of the batches data.
// The first 8-bytes stores the size of the sequence, and the next 8-bytes will store the
// byte array length of every batch data.
func EncodeSequence(batchesData [][]byte) []byte {
	sequence := []byte{}
	metadata := []byte{}
	n := uint64(len(batchesData))
	bn := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(bn, n)
	metadata = append(metadata, bn...)

	for _, seq := range batchesData {
		// Add batch data to byte array
		sequence = append(sequence, seq...)

		// Add batch metadata to byte array
		// Batch metadata contains the byte array length and the Keccak256 hash of the
		// batch data
		n := uint64(len(seq))
		bn := make([]byte, 8) //nolint:gomnd
		binary.BigEndian.PutUint64(bn, n)
		hash := crypto.Keccak256Hash(seq)
		metadata = append(metadata, bn...)
		metadata = append(metadata, hash.Bytes()...)
	}
	sequence = append(metadata, sequence...)

	// Blob serialization
	sequence = codec.ConvertByPaddingEmptyByte(sequence)

	return sequence
}

// DecodeSequence is the helper function to decode 1D byte array into sequence data. The
// encoding scheme is ensured to be lossless.
//
// When decoding the blob data, the first n+1 8-bytes of the blob contains the metadata of
// the batches data.
// The first 8-bytes stores the size of the sequence, and the next 8-bytes will store the
// byte array length of every batch data.
func DecodeSequence(blobData []byte) ([][]byte, []common.Hash) {
	// Blob deserialization
	blobData = codec.RemoveEmptyByteFromPaddedBytes(blobData)

	bn := blobData[:8]
	n := binary.BigEndian.Uint64(bn)
	// Each batch metadata contains the batch data byte array length (8 byte) and the
	// batch data hash (32 byte)
	metadata := blobData[8 : 40*n+8]
	sequence := blobData[40*n+8:]

	batchesData := [][]byte{}
	batchesHash := []common.Hash{}
	idx := uint64(0)
	for i := uint64(0); i < n; i++ {
		// Get batch metadata
		bn := metadata[40*i : 40*i+8]
		n := binary.BigEndian.Uint64(bn)

		hash := common.BytesToHash(metadata[40*i+8 : 40*(i+1)])
		batchesHash = append(batchesHash, hash)

		// Get batch data
		batchesData = append(batchesData, sequence[idx:idx+n])
		idx += n
	}

	return batchesData, batchesHash
}
