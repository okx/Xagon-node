package l1infotree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	emptyHistoricL1InfoTreeRoot = "0x27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"
)

// L1InfoTreeRecursive is a recursive implementation of the L1InfoTree of Feijoa
type L1InfoTreeRecursive struct {
	historicL1InfoTree *L1InfoTree
	currentLeaf        common.Hash
}

// L1InfoTreeRecursiveSnapshot provides the information generated when a new
// leaf is added to the tree
type L1InfoTreeRecursiveSnapshot struct {
	HistoricL1InfoTreeRoot common.Hash
	L1Data                 common.Hash
	L1InfoTreeRoot         common.Hash
}

// NewL1InfoTreeRecursive creates a new empty L1InfoTreeRecursive
func NewL1InfoTreeRecursive(height uint8) (*L1InfoTreeRecursive, error) {
	historic, err := NewL1InfoTree(height, nil)
	if err != nil {
		return nil, err
	}

	mtr := &L1InfoTreeRecursive{
		historicL1InfoTree: historic,
		currentLeaf:        common.Hash{},
	}
	return mtr, nil
}

// NewL1InfoTreeRecursiveFromLeaves creates a new L1InfoTreeRecursive from leaves as they are
func NewL1InfoTreeRecursiveFromLeaves(height uint8, leaves [][32]byte) (*L1InfoTreeRecursive, error) {
	mtr, err := NewL1InfoTreeRecursive(height)
	if err != nil {
		return nil, err
	}

	for i, leaf := range leaves {
		_, err := mtr.AddLeaf(uint32(i), leaf)
		if err != nil {
			return nil, err
		}
		mtr.currentLeaf = leaf
	}
	return mtr, nil
}

// AddLeaf hashes the current historicL1InfoRoot + currentLeaf data into the new historicLeaf value,
// then adds it to the historicL1InfoTree and finally stores the new leaf as the currentLeaf
func (mt *L1InfoTreeRecursive) AddLeaf(index uint32, leaf [32]byte) (common.Hash, error) {
	// adds the current l1InfoTreeRoot into the historic tree to generate
	// the next historicL2InfoTreeRoot
	l1InfoTreeRoot := mt.GetRoot()
	_, err := mt.historicL1InfoTree.AddLeaf(index, l1InfoTreeRoot)
	if err != nil {
		return common.Hash{}, err
	}

	mt.currentLeaf = leaf

	return mt.GetRoot(), nil
}

// GetRoot returns the root of the L1InfoTreeRecursive
func (mt *L1InfoTreeRecursive) GetRoot() common.Hash {
	// if the historicL1InfoTree is empty and the the current leaf is also empty
	// returns the root as all zeros 0x0000...0000
	if mt.historicL1InfoTree.GetRoot().String() == emptyHistoricL1InfoTreeRoot &&
		mt.currentLeaf.Cmp(common.Hash{}) == 0 {
		return common.Hash{}
	}

	l1InfoTreeRoot := crypto.Keccak256Hash(mt.historicL1InfoTree.GetRoot().Bytes(), mt.currentLeaf[:])
	return l1InfoTreeRoot
}

// GetHistoricRoot returns the root of the HistoricL1InfoTree
func (mt *L1InfoTreeRecursive) GetHistoricRoot() common.Hash {
	return mt.historicL1InfoTree.GetRoot()
}

// ComputeMerkleProof computes the Merkle proof from the leaves
func (mt *L1InfoTreeRecursive) ComputeMerkleProof(gerIndex uint32, leaves [][32]byte) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, leaves)
}
