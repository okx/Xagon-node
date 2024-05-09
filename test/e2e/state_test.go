package e2e

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	// Load test vectors
	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/src/state-transition/no-data/general.json")
	require.NoError(t, err)
	genesisFileAsStr, err := config.LoadGenesisFileAsString("../config/test.genesis.config.json")
	require.NoError(t, err)
	genesisConfig, err := config.LoadGenesisFromJSONString(genesisFileAsStr)
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			opsCfg := &operations.Config{
				State: &state.Config{
					MaxCumulativeGasUsed: 800000,
				},
				SequenceSender: &operations.SequenceSenderConfig{
					SenderAddress:                            testCase.SequencerAddress,
					LastBatchVirtualizationTimeMaxWaitPeriod: "5s",
					WaitPeriodSendSequence:                   "5s",
					MaxTxSizeForL1:                           131072,
					PrivateKey:                               testCase.SequencerPrivateKey,
				},
			}
			opsman, err := operations.NewManager(ctx, opsCfg)
			require.NoError(t, err)

			genesisAccounts := make(map[string]big.Int)
			for _, gacc := range testCase.GenesisAccounts {
				genesisAccounts[gacc.Address] = gacc.Balance.Int
			}
			require.NoError(t, opsman.SetGenesisAccountsBalance(genesisConfig.Genesis.BlockNumber, genesisAccounts))

			// Check initial root
			require.NoError(t, opsman.CheckVirtualRoot(testCase.ExpectedOldRoot))

			require.NoError(t, opsman.Setup())

			// convert vector txs
			txs := make([]*types.Transaction, 0, len(testCase.Txs))
			for i := 0; i < len(testCase.Txs); i++ {
				vecTx := testCase.Txs[i]
				var tx types.Transaction
				err := rlp.DecodeBytes([]byte(vecTx.RawTx), &tx)
				require.NoError(t, err)
				txs = append(txs, &tx)
			}

			// send transactions
			_, err = operations.ApplyL2Txs(ctx, txs, nil, nil, operations.VerifiedConfirmationLevel)
			require.NoError(t, err)

			st := opsman.State()

			// Check leaves
			l2Block, err := st.GetLastL2Block(ctx, nil)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(ctx, addr, l2Block.Root())
				require.NoError(t, err)
				require.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("addr: %s expected: %s found: %s", addr.Hex(), leaf.Balance.Text(encoding.Base10), actualBalance.Text(encoding.Base10)))

				actualNonce, err := st.GetNonce(ctx, addr, l2Block.Root())
				require.NoError(t, err)
				require.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, encoding.Base10), fmt.Sprintf("addr: %s expected: %s found: %d", addr.Hex(), leaf.Nonce, actualNonce))
			}

			// Check virtual root against the expected state
			require.NoError(t, opsman.CheckVirtualRoot(testCase.ExpectedNewRoot))

			// Check that last virtual and consolidated l2 block are the same
			lastConsolidatedL2BlockNumber, err := st.GetLastConsolidatedL2BlockNumber(ctx, nil)
			require.NoError(t, err)
			lastVirtualL2BlockNumber, err := st.GetLastL2BlockNumber(ctx, nil)
			require.NoError(t, err)
			t.Logf("lastConsolidatedL2BlockNumber: %d lastVirtualL2BlockNumber: %d", lastConsolidatedL2BlockNumber, lastVirtualL2BlockNumber)
			require.Equal(t, lastConsolidatedL2BlockNumber, lastVirtualL2BlockNumber)

			require.NoError(t, operations.Teardown())
		})
	}
}
