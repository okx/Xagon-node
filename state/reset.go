package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

// Reset resets the state to the given L1 block number
func (s *State) Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	// Reset from DB to L1 block number, this will remove in cascade:
	//  - VirtualBatches
	//  - VerifiedBatches
	//  - Entries in exit_root table
	err := s.ResetToL1BlockNumber(ctx, blockNumber, dbTx)
	if err != nil {
		log.Error("error resetting L1BlockNumber. Error: ", err)
		return err
	}
	s.ResetL1InfoTree()
	return nil
}

// ResetL1InfoTree resets the L1InfoTree
func (s *State) ResetL1InfoTree() {
	// Discard L1InfoTree cache
	// We can't rebuild cache, because we are inside a transaction, so we dont known
	// is going to be a commit or a rollback. So is going to be rebuild on the next
	// request that needs it.
	s.l1InfoTree = nil
}
