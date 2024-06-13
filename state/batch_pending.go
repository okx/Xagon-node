package state

import "time"

// PendingBatch represents a batch pending to be executed
type PendingBatch struct {
	BatchNumber  uint64
	BlobInnerNum uint64
	CreatedAt    time.Time
	Processed    bool
}
