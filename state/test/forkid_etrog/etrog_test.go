package etrog_test

import (
	"context"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/ci/vectors"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	test "github.com/0xPolygonHermez/zkevm-node/state/test/forkid_common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	testsFolder = "../../../test/vectors/src/etrog/"
)

var (
	testState *state.State
	forkID    = uint64(state.FORKID_ETROG)
	stateCfg  = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
)

func TestMain(m *testing.M) {
	testState = test.InitTestState(stateCfg)
	defer test.CloseTestState()
	result := m.Run()
	os.Exit(result)
}

// TestStateTransition tests state using test vectors
func TestStateTransition(t *testing.T) {
	ctx := context.Background()

	// Get all tests vector in the etrog folder
	files, err := os.ReadDir(testsFolder)
	require.NoError(t, err)

	for _, file := range files {
		// Load test vectors
		testCases, err := vectors.LoadStateTransitionTestCasesEtrog(testsFolder + file.Name())
		require.NoError(t, err)

		// Run test cases
		for i, testCase := range testCases {
			block := state.Block{
				BlockNumber: uint64(i + 1),
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}

			genesisActions := vectors.GenerateGenesisActionsEtrog(testCase.Genesis)

			dbTx, err := testState.BeginStateTransaction(ctx)
			require.NoError(t, err)

			stateRoot, err := testState.SetGenesis(ctx, block, state.Genesis{Actions: genesisActions}, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.Equal(t, testCase.ExpectedOldStateRoot, stateRoot.String())
			err = dbTx.Rollback(ctx)
			require.NoError(t, err)

			// convert vector txs
			txs := make([]state.L2TxRaw, 0, len(testCase.Txs))
			for i := 0; i < len(testCase.Txs); i++ {
				vecTx := testCase.Txs[i]
				if vecTx.Type != 0x0b {
					tx, err := state.DecodeTx(vecTx.RawTx)
					require.NoError(t, err)
					l2Tx := state.L2TxRaw{
						Tx:                   *tx,
						EfficiencyPercentage: 255,
					}
					txs = append(txs, l2Tx)
				}
			}

			timestampLimit, ok := big.NewInt(0).SetString(testCase.TimestampLimit, 10)
			require.True(t, ok)

			if len(txs) > 0 {
				// Generate batchdata from the txs in the test and compared with the vector
				l2block := state.L2BlockRaw{
					ChangeL2BlockHeader: state.ChangeL2BlockHeader{
						DeltaTimestamp:  uint32(timestampLimit.Uint64()),
						IndexL1InfoTree: testCase.Txs[0].IndexL1InfoTree,
					},
					Transactions: txs,
				}

				batch := state.BatchRawV2{
					Blocks: []state.L2BlockRaw{l2block},
				}

				batchData, err := state.EncodeBatchV2(&batch)
				require.NoError(t, err)

				require.Equal(t, common.FromHex(testCase.BatchL2Data), batchData)
			}

			processRequest := state.ProcessRequest{
				BatchNumber:             uint64(i + 1),
				L1InfoRoot_V2:           common.HexToHash(testCase.L1InfoRoot),
				OldStateRoot:            stateRoot,
				OldAccInputHash:         common.HexToHash(testCase.OldAccInputHash),
				Transactions:            common.FromHex(testCase.BatchL2Data),
				TimestampLimit_V2:       timestampLimit.Uint64(),
				Coinbase:                common.HexToAddress(testCase.SequencerAddress),
				ForkID:                  testCase.ForkID,
				SkipVerifyL1InfoRoot_V2: testCase.L1InfoTree.SkipVerifyL1InfoRoot,
			}

			processResponse, _, _ := testState.ProcessBatchV2(ctx, processRequest, true)
			require.Nil(t, processResponse.ExecutorError)
			require.Equal(t, testCase.ExpectedNewStateRoot, processResponse.NewStateRoot.String())
		}
	}
}
