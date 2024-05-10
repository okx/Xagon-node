package merkletree

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetCode(t *testing.T) {
	ctx := context.Background()
	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")

	cfg := Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	c, _, _ := NewMTDBServiceClient(ctx, cfg)
	sTree := NewStateTree(c)

	type testCase struct {
		name           string
		addr           common.Address
		root           []byte
		expectedResult []byte
		expectedError  error
		setup          func(*testing.T, *testCase, *StateTree)
	}

	testCases := []testCase{
		{
			name:           "get existent code successfully",
			addr:           common.HexToAddress("0x1"),
			root:           common.HexToHash("0x0").Bytes(),
			expectedResult: hex.DecodeBig(EmitLog2.EmitLog2Bin).Bytes(),
			expectedError:  nil,
			setup: func(t *testing.T, tc *testCase, sTree *StateTree) {
				txID := uuid.NewString()

				err := sTree.StartBlock(ctx, common.Hash(tc.root), txID)
				require.NoError(t, err)

				newRoot, _, err := sTree.SetCode(ctx, tc.addr, tc.expectedResult, tc.root, txID)
				require.NoError(t, err)
				tc.root = newRoot

				err = sTree.FinishBlock(ctx, common.Hash(tc.root), txID)
				require.NoError(t, err)

				err = sTree.Flush(ctx, common.Hash(newRoot), txID)
				require.NoError(t, err)
			},
		},
		{
			name:           "get non-existent code successfully",
			addr:           common.HexToAddress("0x2"),
			root:           common.HexToHash("0x0").Bytes(),
			expectedResult: []byte{},
			expectedError:  nil,
			setup: func(t *testing.T, tc *testCase, sTree *StateTree) {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			tc.setup(t, &tc, sTree)

			result, err := sTree.GetCode(ctx, tc.addr, tc.root)
			require.NoError(t, err)

			if tc.expectedResult != nil || result != nil {
				require.Equal(t, len(tc.expectedResult), len(result))
				require.ElementsMatch(t, tc.expectedResult, result)
			}

			if tc.expectedError != nil || err != nil {
				require.Equal(t, tc.expectedError, err)
			}
		})
	}
}
