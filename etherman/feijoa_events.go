package etherman

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// BlobType is the type of the blob type
type BlobType uint8

const (
	// TypeCallData The data is stored on call data directly
	TypeCallData BlobType = 0
	// TypeBlobTransaction The data is stored on a blob
	TypeBlobTransaction BlobType = 1
	// TypeForcedBlob The data is a forced Blob
	TypeForcedBlob BlobType = 2
)

// SequenceBlob is for each Blob inside a SequenceBlobs
type SequenceBlob struct {
	Type   BlobType
	Params BlobCommonParams
	Data   []byte
	// Field only valid if BlobType == BlobTransaction
	BlobBlobTypeParams *BlobBlobTypeParams
}

func (s *SequenceBlob) String() string {
	return fmt.Sprintf("Type: %d, Params: %v, Data: %v, BlobBlobTypeParams: %v", s.Type, s.Params, s.Data, s.BlobBlobTypeParams)
}

// BlobCommonParams is the data for a SequenceBlob
type BlobCommonParams struct {
	MaxSequenceTimestamp uint64
	ZkGasLimit           uint64
	L1InfoLeafIndex      uint32
}

// BlobBlobTypeParams is the data for a SequenceBlob stored on a Blob
// case: if (currentBlob.blobType ==> BLOBTX_BLOB_TYPE)
// sames as calldata plus BlobIndex, ...
type BlobBlobTypeParams struct {
	BlobIndex  *big.Int
	Z          []byte
	Y          []byte
	Commitment kzg4844.Commitment
	Proof      kzg4844.Proof
}

// SequenceBlobs is the data in the event SequenceBlobs
type SequenceBlobs struct {
	Blobs             []SequenceBlob
	L2Coinbase        common.Address // from Calldata
	FinalAccInputHash common.Hash
	EventData         *SequenceBlobsEventData
}

// SequenceBlobsEventData is the data in the event SequenceBlobs
type SequenceBlobsEventData struct {
	// LastBlobSequenced is the count of blob sequenced after process this event
	//  if the first event have 1 blob -> lastBlobSequenced=1
	LastBlobSequenced uint64
}

func (s *SequenceBlobs) thereIsAnyBlobType() bool {
	for blobIndex := range s.Blobs {
		if s.Blobs[blobIndex].Type == TypeBlobTransaction {
			return true
		}
	}
	return false
}
