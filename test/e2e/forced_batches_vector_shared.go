package e2e

import (
	"context"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func LaunchTestForcedBatchesVectorFilesGroup(t *testing.T, vectorFilesDir string) {

	//vectorFilesDir := "./../vectors/src/state-transition/forced-tx/group1"
	ctx := context.Background()
	genesisFileAsStr, err := config.LoadGenesisFileAsString("../../test/config/test.genesis.config.json")
	require.NoError(t, err)
	genesisConfig, err := config.LoadGenesisFromJSONString(genesisFileAsStr)
	require.NoError(t, err)
	err = filepath.Walk(vectorFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), "list.json") {

			t.Run(info.Name(), func(t *testing.T) {

				defer func() {
					require.NoError(t, operations.Teardown())
				}()

				// Load test vectors
				log.Info("=====================================================================")
				log.Info(path)
				log.Info("=====================================================================")
				testCase, err := vectors.LoadStateTransitionTestCaseV2(path)
				require.NoError(t, err)

				opsCfg := operations.GetDefaultOperationsConfig()
				opsCfg.State.MaxCumulativeGasUsed = 80000000000
				opsCfg.State.ForkIDIntervals = []state.ForkIDInterval{{
					FromBatchNumber: 0,
					ToBatchNumber:   math.MaxUint64,
					ForkId:          state.FORKID_ETROG,
					Version:         "",
				}}
				opsman, err := operations.NewManager(ctx, opsCfg)
				require.NoError(t, err)

				// Setting Genesis
				log.Info("###################")
				log.Info("# Setting Genesis #")
				log.Info("###################")
				genesisActions := vectors.GenerateGenesisActions(testCase.Genesis)
				require.NoError(t, opsman.SetGenesis(genesisConfig.Genesis.RollupBlockNumber, genesisActions))
				require.NoError(t, opsman.SetForkID(genesisConfig.Genesis.RollupBlockNumber, forkID6))
				actualOldStateRoot, err := opsman.State().GetLastStateRoot(ctx, nil)
				require.NoError(t, err)
				require.NoError(t, opsman.Setup())

				// Check initial root
				log.Info("################################")
				log.Info("# Verifying initial state root #")
				log.Info("################################")
				require.NoError(t, err)
				require.Equal(t, testCase.ExpectedOldStateRoot, actualOldStateRoot.Hex())
				decodedData, err := hex.DecodeHex(testCase.BatchL2Data)
				require.NoError(t, err)
				_, txBytes, _, err := state.DecodeTxs(decodedData, forkID6)
				require.NoError(t, err)
				forcedBatch, err := sendForcedBatchForVector(t, txBytes, opsman)
				require.NoError(t, err)
				actualNewStateRoot := forcedBatch.StateRoot
				isClosed, err := opsman.State().IsBatchClosed(ctx, forcedBatch.BatchNumber, nil)
				require.NoError(t, err)

				// wait until is closed
				for !isClosed {
					time.Sleep(1 * time.Second)
					isClosed, err = opsman.State().IsBatchClosed(ctx, forcedBatch.BatchNumber, nil)
					require.NoError(t, err)
				}

				log.Info("#######################")
				log.Info("# Verifying new leaves #")
				log.Info("#######################")
				merkleTree := opsman.State().GetTree()
				for _, expectedNewLeaf := range testCase.ExpectedNewLeafs {
					if expectedNewLeaf.IsSmartContract {
						log.Info("Smart Contract Address: ", expectedNewLeaf.Address)
					} else {
						log.Info("Account Address: ", expectedNewLeaf.Address)
					}
					log.Info("Verifying Balance...")
					actualBalance, err := merkleTree.GetBalance(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
					require.NoError(t, err)
					require.Equal(t, expectedNewLeaf.Balance.String(), actualBalance.String())

					log.Info("Verifying Nonce...")
					actualNonce, err := merkleTree.GetNonce(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
					require.NoError(t, err)
					require.Equal(t, expectedNewLeaf.Nonce, actualNonce.String())
					if expectedNewLeaf.IsSmartContract {
						log.Info("Verifying Storage...")
						for positionHex, expectedNewStorageHex := range expectedNewLeaf.Storage {
							position, ok := big.NewInt(0).SetString(positionHex[2:], 16)
							require.True(t, ok)
							expectedNewStorage, ok := big.NewInt(0).SetString(expectedNewStorageHex[2:], 16)
							require.True(t, ok)
							actualStorage, err := merkleTree.GetStorageAt(ctx, common.HexToAddress(expectedNewLeaf.Address), position, actualNewStateRoot.Bytes())
							require.NoError(t, err)
							require.Equal(t, expectedNewStorage, actualStorage)
						}

						log.Info("Verifying HashBytecode...")
						actualHashByteCode, err := merkleTree.GetCodeHash(ctx, common.HexToAddress(expectedNewLeaf.Address), actualNewStateRoot.Bytes())
						require.NoError(t, err)
						require.Equal(t, expectedNewLeaf.HashBytecode, common.BytesToHash(actualHashByteCode).String())
					}
				}
				return
			})

			return nil
		}
		return nil
	})
	require.NoError(t, err)
}
