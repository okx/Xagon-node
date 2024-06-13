package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// L1InfoTreeRecursiveExitRootStorageEntry leaf of the L1InfoTreeRecursive
type L1InfoTreeRecursiveExitRootStorageEntry L1InfoTreeExitRootStorageEntry

func (s *State) buildL1InfoTreeRecursiveCacheIfNeed(ctx context.Context, dbTx pgx.Tx) error {
	if s.l1InfoTreeRecursive != nil {
		return nil
	}
	log.Debugf("Building L1InfoTree cache")
	allLeaves, err := s.GetAllL1InfoTreeRecursiveRootEntries(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves. Error: ", err)
		return fmt.Errorf("error getting all leaves. Error: %w", err)
	}

	var leaves [][32]byte
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	mt, err := l1infotree.NewL1InfoTreeRecursiveFromLeaves(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTreeRecursive = mt
	return nil
}

// AddL1InfoTreeRecursiveLeaf adds a new leaf to the L1InfoTree and returns the entry and error
func (s *State) AddL1InfoTreeRecursiveLeaf(ctx context.Context, l1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeExitRootStorageEntry, error) {
	var newIndex uint32
	gerIndex, err := s.GetLatestIndex(ctx, dbTx)
	if err != nil && !errors.Is(err, ErrNotFound) {
		log.Error("error getting latest L1InfoTreeRecursive index. Error: ", err)
		return nil, err
	} else if err == nil {
		newIndex = gerIndex + 1
	}
	err = s.buildL1InfoTreeCacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTreeRecursive cache. Error: ", err)
		return nil, err
	}
	log.Debug("latestIndex: ", gerIndex)
	l1InfoTreeRoot, err := s.l1InfoTreeRecursive.AddLeaf(newIndex, l1InfoTreeLeaf.Hash())
	if err != nil {
		log.Error("error add new leaf to the L1InfoTreeRecursive. Error: ", err)
		return nil, err
	}
	entry := L1InfoTreeExitRootStorageEntry{
		L1InfoTreeLeaf:  *l1InfoTreeLeaf,
		L1InfoTreeRoot:  l1InfoTreeRoot,
		L1InfoTreeIndex: newIndex,
	}
	err = s.AddL1InfoRootToExitRoot(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	return &entry, nil
}

// GetCurrentL1InfoTreeRecursiveRoot Return current L1InfoRoot
func (s *State) GetCurrentL1InfoTreeRecursiveRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	err := s.buildL1InfoTreeRecursiveCacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTree cache. Error: ", err)
		return ZeroHash, err
	}
	return s.l1InfoTreeRecursive.GetRoot(), nil
}
