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
	// Reset L1InfoTree siblings and leaves
	allLeaves, err := s.GetAllL1InfoRootEntries(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves to reset l1InfoTree. Error: ", err)
		return err
	}
	var leaves [][32]byte
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	mt, err := s.l1InfoTree.ResetL1InfoTree(leaves)
	if err != nil {
		log.Error("error resetting l1InfoTree. Error: ", err)
		return err
	}
	s.l1InfoTree = mt
	return nil
}

// AddL1InfoTreeLeaf adds a new leaf to the L1InfoTree and returns the entry and error
func (s *State) AddL1InfoTreeLeaf(ctx context.Context, l1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeExitRootStorageEntry, error) {
	var stateTx *StateTx
	if dbTx != nil {
		var ok bool
		stateTx, ok = dbTx.(*StateTx)
		if !ok {
			return nil, fmt.Errorf("error casting dbTx to stateTx")
		}
	}
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
	if stateTx != nil {
		stateTx.SetL1InfoTreeModified()
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
