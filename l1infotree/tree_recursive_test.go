package l1infotree_test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	l1InfoRootRecursiveHeight    = uint8(32)
	emptyL1InfoTreeRecursiveRoot = "0x27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"
	filenameTestData             = "../test/vectors/src/merkle-tree/l1-info-tree-recursive/smt-full-output.json"
)

type vectorTestData struct {
	GlobalExitRoot         common.Hash   `json:"globalExitRoot"`
	BlockHash              common.Hash   `json:"blockHash"`
	MinTimestamp           string        `json:"minTimestamp"`
	SmtProof               []common.Hash `json:"smtProof"`
	Index                  uint32        `json:"index"`
	PreviousIndex          uint32        `json:"previousIndex"`
	PreviousL1InfoTreeRoot common.Hash   `json:"previousL1InfoTreeRoot"`
	L1DataHash             common.Hash   `json:"l1DataHash"`
	L1InfoTreeRoot         common.Hash   `json:"l1InfoTreeRoot"`
	HistoricL1InfoRoot     common.Hash   `json:"historicL1InfoRoot"`
}

func readData(t *testing.T) []vectorTestData {
	data, err := os.ReadFile(filenameTestData)
	require.NoError(t, err)
	var mtTestVectors []vectorTestData
	err = json.Unmarshal(data, &mtTestVectors)
	require.NoError(t, err)
	return mtTestVectors
}

func TestEmptyL1InfoRootRecursive(t *testing.T) {
	mtr, err := l1infotree.NewL1InfoTreeRecursive(l1InfoRootRecursiveHeight)
	require.NoError(t, err)
	require.NotNil(t, mtr)
	root := mtr.GetRoot()
	require.Equal(t, common.Hash{}.String(), root.String())
}

func TestEmptyHistoricL1InfoRootRecursive(t *testing.T) {
	mtr, err := l1infotree.NewL1InfoTreeRecursive(l1InfoRootRecursiveHeight)
	require.NoError(t, err)
	require.NotNil(t, mtr)
	root := mtr.GetHistoricRoot()
	require.Equal(t, emptyL1InfoTreeRecursiveRoot, root.String())
}

func TestBuildTreeVectorData(t *testing.T) {
	data := readData(t)
	mtr, err := l1infotree.NewL1InfoTreeRecursive(l1InfoRootRecursiveHeight)
	require.NoError(t, err)
	for _, testVector := range data {
		minTimestamp, err := strconv.ParseUint(testVector.MinTimestamp, 10, 0)
		require.NoError(t, err)
		l1Data := l1infotree.HashLeafData(testVector.GlobalExitRoot, testVector.BlockHash, minTimestamp)
		l1DataHash := common.BytesToHash(l1Data[:])
		assert.Equal(t, testVector.L1DataHash.String(), l1DataHash.String(), "l1Data doesn't match leaf", testVector.Index)

		l1InfoTreeRoot, err := mtr.AddLeaf(testVector.Index-1, l1Data)
		require.NoError(t, err)
		assert.Equal(t, testVector.L1InfoTreeRoot.String(), l1InfoTreeRoot.String(), "l1InfoTreeRoot doesn't match leaf", testVector.Index)
		assert.Equal(t, testVector.L1InfoTreeRoot.String(), mtr.GetRoot().String(), "l1InfoTreeRoot doesn't match leaf", testVector.Index)
		assert.Equal(t, testVector.HistoricL1InfoRoot.String(), mtr.GetHistoricRoot().String(), "HistoricL1InfoTreeRoot doesn't match leaf", testVector.Index)
	}
}

func TestBuildTreeFromLeaves(t *testing.T) {
	data := readData(t)

	leaves := [][32]byte{}
	for _, testVector := range data {
		leaves = append(leaves, testVector.L1DataHash)
	}

	newMtr, err := l1infotree.NewL1InfoTreeRecursiveFromLeaves(l1InfoRootRecursiveHeight, leaves)
	require.NoError(t, err)
	assert.Equal(t, data[len(data)-1].L1InfoTreeRoot.String(), newMtr.GetRoot().String(), "L1InfoTreeRoot doesn't match leaf")
}

func TestProofsTreeVectorData(t *testing.T) {
	data := readData(t)
	mtr, err := l1infotree.NewL1InfoTreeRecursive(l1InfoRootRecursiveHeight)
	require.NoError(t, err)

	leaves := [][32]byte{}
	for _, testVector := range data {
		l1InfoTreeRoot, err := mtr.AddLeaf(testVector.Index-1, testVector.L1DataHash)
		require.NoError(t, err)

		leaves = append(leaves, l1InfoTreeRoot)

		mp, _, err := mtr.ComputeMerkleProof(testVector.Index, leaves)
		require.NoError(t, err)
		for i, v := range mp {
			c := common.Hash(v)
			if c.String() != testVector.SmtProof[i].String() {
				log.Info("MerkleProof: index ", testVector.Index, " mk:", i, " v:", c.String(), " expected:", testVector.SmtProof[i].String())
			}
		}
	}
}
