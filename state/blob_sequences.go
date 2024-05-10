package state

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrBlobSequenceIndex is returned when the blob sequence index is not correct
	ErrBlobSequenceIndex = errors.New("blob sequence index is not correct")
	// ErrBlobSequenceTime is returned when the blob sequence time is not correct
	ErrBlobSequenceTime = errors.New("blob sequence time is not correct")
)

// BlobSequence represents a blob sequence.
type BlobSequence struct {
	BlobSequenceIndex  uint64
	L2Coinbase         common.Address
	FinalAccInputHash  common.Hash
	FirstBlobSequenced uint64    // Is calculated from previous blob sequence
	LastBlobSequenced  uint64    // That comes from the event
	CreateAt           time.Time // time of the L1block
	ReceivedAt         time.Time // time when the blob sequence is received (typically Now())
	BlockNumber        uint64    // L1BlockNumber where appears this event
}

// AddBlobSequence adds a new blob sequence to the state.
// it override pgstorage.AddBlobSequence to add sanity checks
func (s *State) AddBlobSequence(ctx context.Context, blobSequence *BlobSequence, dbTx pgx.Tx) error {
	err := s.sanityCheckAddBlobSequence(ctx, blobSequence, dbTx)
	if err != nil {
		return err
	}
	return s.storage.AddBlobSequence(ctx, blobSequence, dbTx)
}

func (s *State) sanityCheckAddBlobSequence(ctx context.Context, blobSequence *BlobSequence, dbTx pgx.Tx) error {
	previousBlobSequence, err := s.GetLastBlobSequence(ctx, dbTx)
	if err != nil {
		return err
	}
	if previousBlobSequence == nil {
		// Is the  first one
		if blobSequence.BlobSequenceIndex != 1 {
			return fmt.Errorf("TThe firstBlobSequence  index must be 1, not %d. Err: %w", blobSequence.BlobSequenceIndex, ErrBlobSequenceIndex)
		}
		return nil
	}
	// The index must be the previous index + 1
	if previousBlobSequence.BlobSequenceIndex+1 != blobSequence.BlobSequenceIndex {
		return fmt.Errorf("last_index_on_db:%d try_to_insert:%d. Err: %w",
			previousBlobSequence.BlobSequenceIndex,
			blobSequence.BlobSequenceIndex,
			ErrBlobSequenceIndex)
	}
	// The new blob must be newer than the previous one
	if previousBlobSequence.CreateAt.After(blobSequence.CreateAt) {
		return fmt.Errorf("last_create_at_on_db:%d try_to_insert:%d. Err: %w",
			previousBlobSequence.CreateAt.Unix(),
			blobSequence.CreateAt.Unix(),
			ErrBlobSequenceTime)
	}
	return nil
}
