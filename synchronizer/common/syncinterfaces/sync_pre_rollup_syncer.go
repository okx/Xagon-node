package syncinterfaces

import (
	"context"
)

// SyncPreRollupSyncer is the interface for synchronizing pre genesis rollup events
type SyncPreRollupSyncer interface {
	SynchronizePreGenesisRollupEvents(ctx context.Context) error
}
