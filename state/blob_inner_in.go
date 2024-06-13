package state

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/jackc/pgx/v4"
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

func (b BlobType) String() string {
	switch b {
	case TypeCallData:
		return "call_data"
	case TypeBlobTransaction:
		return "blob"
	case TypeForcedBlob:
		return "forced"
	default:
		return "Unknown"
	}
}

// BlobBlobTypeParams is the data for a SequenceBlob stored as a Blob
type BlobBlobTypeParams struct {
	BlobIndex  uint64
	Z          []byte
	Y          []byte
	Commitment kzg4844.Commitment
	Proof      kzg4844.Proof
}

// BlobInner struct
type BlobInner struct {
	BlobSequenceIndex    uint64              // Index of the blobSequence in DB (is a internal number)
	BlobInnerNum         uint64              // Incremental value, starts from 1
	Type                 BlobType            // Type of the blob
	MaxSequenceTimestamp time.Time           // it comes from SequenceBlobs call to contract
	ZkGasLimit           uint64              // it comes from SequenceBlobs call to contract
	L1InfoLeafIndex      uint32              // it comes from SequenceBlobs call to contract
	L1InfoTreeRoot       common.Hash         // obtained from the L1InfoTree
	BlobDataHash         common.Hash         // Hash of the data
	BlobBlobTypeParams   *BlobBlobTypeParams // Field only valid if BlobType == BlobTransaction
	//HowManyBatches 	uint64 			// Number of batches in the blob
	//FirstBatchNumber   uint64              // First batch number of the blob
	//LastBatchNumber    uint64              // Last batch number of the blob
	// We don't need blockNumber because is in BlobSequence
	//BlockNumber             uint64
	//PreviousL1InfoTreeIndex uint32      // ?? we need that?
	//PreviousL1InfoTreeRoot  common.Hash // ?? we need that?
}

func (b *BlobInner) String() string {
	res := fmt.Sprintf("BlobInner{BlobSequenceIndex:%d, BlobInnerNum:%d, Type:%s, MaxSequenceTimestamp:%s, ZkGasLimit:%d, L1InfoLeafIndex:%d, L1InfoTreeRoot:%s, BlobDataHash:%s",
		b.BlobSequenceIndex, b.BlobInnerNum, b.Type.String(), b.MaxSequenceTimestamp.String(), b.ZkGasLimit, b.L1InfoLeafIndex, b.L1InfoTreeRoot.String(), b.BlobDataHash.String())
	if b.BlobBlobTypeParams != nil {
		res += ", BlobBlobTypeParams: " + b.BlobBlobTypeParams.String()
	}
	res += "}"
	return res
}

func (b *BlobBlobTypeParams) String() string {
	return "BlobBlobTypeParams{" +
		"BlobIndex: " + fmt.Sprintf("%d", b.BlobIndex) +
		", Z: " + common.Bytes2Hex(b.Z) +
		", Y: " + common.Bytes2Hex(b.Y) +
		", Commitment: " + common.Bytes2Hex(b.Commitment[:]) +
		", Proof: " + common.Bytes2Hex(b.Proof[:]) +
		"}"
}

// IsEqual compares two BlobInner
func (b *BlobInner) IsEqual(other *BlobInner) bool {
	if b == nil && other == nil {
		return true
	}
	if b == nil || other == nil {
		return false
	}
	return b.String() == other.String()
}

// AddBlobInner adds a blob inner to the database, currently is just a call to storage
func (s *State) AddBlobInner(ctx context.Context, blobInner *BlobInner, dbTx pgx.Tx) error {
	return s.storage.AddBlobInner(ctx, blobInner, dbTx)
}
