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

// L1InfoTreeLeaf leaf of the L1InfoTree
type L1InfoTreeLeaf struct {
	GlobalExitRoot
	PreviousBlockHash common.Hash
}

// L1InfoTreeExitRootStorageEntry entry of the Database
type L1InfoTreeExitRootStorageEntry struct {
	L1InfoTreeLeaf
	L1InfoTreeRoot  common.Hash
	L1InfoTreeIndex uint32
}

// Hash returns the hash of the leaf
func (l *L1InfoTreeLeaf) Hash() common.Hash {
	timestamp := uint64(l.Timestamp.Unix())
	return l1infotree.HashLeafData(l.GlobalExitRoot.GlobalExitRoot, l.PreviousBlockHash, timestamp)
}

func (s *State) buildL1InfoTreeCacheIfNeed(ctx context.Context, dbTx pgx.Tx) error {
	if s.l1InfoTree != nil {
		return nil
	}
	log.Debugf("Building L1InfoTree cache")
	allLeaves, err := s.storage.GetAllL1InfoRootEntries(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves. Error: ", err)
		return fmt.Errorf("error getting all leaves. Error: %w", err)
	}
	var leaves [][32]byte
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	mt, err := l1infotree.NewL1InfoTree(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTree = mt
	return nil
}

// AddL1InfoTreeLeaf adds a new leaf to the L1InfoTree and returns the entry and error
func (s *State) AddL1InfoTreeLeaf(ctx context.Context, l1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeExitRootStorageEntry, error) {
	var newIndex uint32
	gerIndex, err := s.GetLatestIndex(ctx, dbTx)
	if err != nil && !errors.Is(err, ErrNotFound) {
		log.Error("error getting latest l1InfoTree index. Error: ", err)
		return nil, err
	} else if err == nil {
		newIndex = gerIndex + 1
	}
	err = s.buildL1InfoTreeCacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTree cache. Error: ", err)
		return nil, err
	}
	log.Debug("latestIndex: ", gerIndex)
	root, err := s.l1InfoTree.AddLeaf(newIndex, l1InfoTreeLeaf.Hash())
	if err != nil {
		log.Error("error add new leaf to the L1InfoTree. Error: ", err)
		return nil, err
	}
	entry := L1InfoTreeExitRootStorageEntry{
		L1InfoTreeLeaf:  *l1InfoTreeLeaf,
		L1InfoTreeRoot:  root,
		L1InfoTreeIndex: newIndex,
	}
	err = s.AddL1InfoRootToExitRoot(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	return &entry, nil
}

// GetCurrentL1InfoRoot Return current L1InfoRoot
func (s *State) GetCurrentL1InfoRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	err := s.buildL1InfoTreeCacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTree cache. Error: ", err)
		return ZeroHash, err
	}
	root, _, _ := s.l1InfoTree.GetCurrentRootCountAndSiblings()
	return root, nil
}
