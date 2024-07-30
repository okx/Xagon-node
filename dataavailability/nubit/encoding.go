package nubit

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeSequence is the helper function to encode sequence data and their metadata into 1D byte array.
// The encoding scheme is ensured to be lossless.
//
// When encoding the blob data, the first 8-bytes stores the size of the batches (n) in the sequence. The
// next n slots of sized 40 bytes stores the metadata of the batches data.
// The first 8-bytes of the batches metadata stores the batches data length, and the next 32-bytes stores
// the batches hash.
//
// The remaining n slots contains the batches data, each slot length is specified in the retrieved batch
// metadata.
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

	return sequence
}

// DecodeSequence is the helper function to decode the 1D byte array into sequence data and the batches
// metadata. The decoding sceheme is ensured to be lossless and follows the encoding scheme specified in
// the EncodeSequence function.
func DecodeSequence(blobData []byte) ([][]byte, []common.Hash) {
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
