package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	forkID6 = 6
)

func TestConsolidatedBlockNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult *uint64
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "Get consolidated block number successfully",
			ExpectedResult: ptrUint64(10),
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastConsolidatedL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get consolidated block number",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get last consolidated block number from state"),
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastConsolidatedL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last consolidated block number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_consolidatedBlockNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result types.ArgUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, *tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestIsBlockConsolidated(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult bool
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "Query status of block number successfully",
			ExpectedResult: true,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("IsL2BlockConsolidated", context.Background(), uint64(1), m.DbTx).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the consolidation status",
			ExpectedResult: false,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to check if the block is consolidated"),
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("IsL2BlockConsolidated", context.Background(), uint64(1), m.DbTx).
					Return(false, errors.New("failed to check if the block is consolidated")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isBlockConsolidated", "0x1")
			require.NoError(t, err)

			if res.Result != nil {
				var result bool
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestIsBlockVirtualized(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult bool
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "Query status of block number successfully",
			ExpectedResult: true,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("IsL2BlockVirtualized", context.Background(), uint64(1), m.DbTx).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the virtualization status",
			ExpectedResult: false,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to check if the block is virtualized"),
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("IsL2BlockVirtualized", context.Background(), uint64(1), m.DbTx).
					Return(false, errors.New("failed to check if the block is virtualized")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isBlockVirtualized", "0x1")
			require.NoError(t, err)

			if res.Result != nil {
				var result bool
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestBatchNumberByBlockNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()
	blockNumber := uint64(1)
	batchNumber := uint64(1)

	type testCase struct {
		Name           string
		ExpectedResult *uint64
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "get batch number by block number successfully",
			ExpectedResult: &batchNumber,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, m.DbTx).
					Return(batchNumber, nil).
					Once()
			},
		},
		{
			Name:           "failed to get batch number",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get batch number from block number"),
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, m.DbTx).
					Return(uint64(0), errors.New("failed to get batch number of l2 batchNum")).
					Once()
			},
		},
		{
			Name:           "batch number not found",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, m.DbTx).
					Return(uint64(0), state.ErrNotFound).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_batchNumberByBlockNumber", hex.EncodeUint64(blockNumber))
			require.NoError(t, err)

			if tc.ExpectedResult != nil {
				var result types.ArgUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, *tc.ExpectedResult, uint64(result))
			} else {
				if res.Result == nil {
					assert.Nil(t, res.Result)
				} else {
					var result *uint64
					err = json.Unmarshal(res.Result, &result)
					require.NoError(t, err)
					assert.Nil(t, result)
				}
			}

			if tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			} else {
				assert.Nil(t, res.Error)
			}
		})
	}
}

func TestBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "get batch number successfully",
			ExpectedError:  nil,
			ExpectedResult: 10,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastBatchNumber", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastBatchNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_batchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result types.ArgUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestVirtualBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "get virtual batch number successfully",
			ExpectedError:  nil,
			ExpectedResult: 10,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastVirtualBatchNum", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get virtual batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last virtual batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastVirtualBatchNum", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_virtualBatchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result types.ArgUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestVerifiedBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper)
	}

	testCases := []testCase{
		{
			Name:           "get verified batch number successfully",
			ExpectedError:  nil,
			ExpectedResult: 10,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastVerifiedBatch", context.Background(), m.DbTx).
					Return(&state.VerifiedBatch{BatchNumber: uint64(10)}, nil).
					Once()
			},
		},
		{
			Name:           "failed to get verified batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last verified batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastVerifiedBatch", context.Background(), m.DbTx).
					Return(nil, errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_verifiedBatchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result types.ArgUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetBatchByNumber(t *testing.T) {
	type testCase struct {
		Name           string
		Number         string
		WithTxDetail   bool
		ExpectedResult *types.Batch
		ExpectedError  types.Error
		SetupMocks     func(*mockedServer, *mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "Batch not found",
			Number:         "0x123",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:         "get specific batch successfully with tx detail",
			Number:       "0x345",
			WithTxDetail: true,
			ExpectedResult: &types.Batch{
				Number:              1,
				Coinbase:            common.HexToAddress("0x1"),
				StateRoot:           common.HexToHash("0x2"),
				AccInputHash:        common.HexToHash("0x3"),
				GlobalExitRoot:      common.HexToHash("0x4"),
				Timestamp:           1,
				SendSequencesTxHash: ptrHash(common.HexToHash("0x10")),
				VerifyBatchTxHash:   ptrHash(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				txs := []*ethTypes.Transaction{
					signTx(ethTypes.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
					signTx(ethTypes.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
				}

				batchTxs := make([]ethTypes.Transaction, 0, len(txs))
				effectivePercentages := make([]uint8, 0, len(txs))
				tc.ExpectedResult.Transactions = []types.TransactionOrHash{}
				receipts := []*ethTypes.Receipt{}
				blocks := []ethTypes.Block{}
				for i, tx := range txs {
					block := ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(i))}).WithBody([]*ethTypes.Transaction{tx}, []*ethTypes.Header{})
					blocks = append(blocks, *block)
					receipt := ethTypes.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = block.Number()
					receipt.BlockHash = block.Hash()
					receipts = append(receipts, receipt)
					from, _ := state.GetSender(*tx)
					V, R, S := tx.RawSignatureValues()

					tc.ExpectedResult.Transactions = append(tc.ExpectedResult.Transactions,
						types.TransactionOrHash{
							Tx: &types.Transaction{
								Nonce:       types.ArgUint64(tx.Nonce()),
								GasPrice:    types.ArgBig(*tx.GasPrice()),
								Gas:         types.ArgUint64(tx.Gas()),
								To:          tx.To(),
								Value:       types.ArgBig(*tx.Value()),
								Input:       tx.Data(),
								Hash:        tx.Hash(),
								From:        from,
								BlockNumber: ptrArgUint64FromUint64(block.NumberU64()),
								BlockHash:   ptrHash(receipt.BlockHash),
								TxIndex:     ptrArgUint64FromUint(receipt.TransactionIndex),
								ChainID:     types.ArgBig(*tx.ChainId()),
								Type:        types.ArgUint64(tx.Type()),
								V:           types.ArgBig(*V),
								R:           types.ArgBig(*R),
								S:           types.ArgBig(*S),
							},
						},
					)

					batchTxs = append(batchTxs, *tx)
					effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
				}
				batchL2Data, err := state.EncodeTransactions(batchTxs, effectivePercentages, forkID6)
				require.NoError(t, err)
				tc.ExpectedResult.BatchL2Data = batchL2Data
				batch := &state.Batch{
					BatchNumber:    1,
					Coinbase:       common.HexToAddress("0x1"),
					StateRoot:      common.HexToHash("0x2"),
					AccInputHash:   common.HexToHash("0x3"),
					GlobalExitRoot: common.HexToHash("0x4"),
					Timestamp:      time.Unix(1, 0),
					BatchL2Data:    batchL2Data,
				}

				m.State.
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(batch, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, m.DbTx).
					Return(&ger, nil).
					Once()

				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(receipts[i], nil).
						Once()
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(batchTxs, effectivePercentages, nil).
					Once()

				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(blocks, nil).
					Once()
			},
		},
		{
			Name:         "get specific batch successfully without tx detail",
			Number:       "0x345",
			WithTxDetail: false,
			ExpectedResult: &types.Batch{
				Number:              1,
				Coinbase:            common.HexToAddress("0x1"),
				StateRoot:           common.HexToHash("0x2"),
				AccInputHash:        common.HexToHash("0x3"),
				GlobalExitRoot:      common.HexToHash("0x4"),
				Timestamp:           1,
				SendSequencesTxHash: ptrHash(common.HexToHash("0x10")),
				VerifyBatchTxHash:   ptrHash(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				txs := []*ethTypes.Transaction{
					signTx(ethTypes.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
					signTx(ethTypes.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
				}

				batchTxs := make([]ethTypes.Transaction, 0, len(txs))
				effectivePercentages := make([]uint8, 0, len(txs))
				tc.ExpectedResult.Transactions = []types.TransactionOrHash{}

				receipts := []*ethTypes.Receipt{}
				blocks := []ethTypes.Block{}
				for i, tx := range txs {
					block := ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(i))}).WithBody([]*ethTypes.Transaction{tx}, []*ethTypes.Header{})
					blocks = append(blocks, *block)
					receipt := ethTypes.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = block.Number()
					receipt.BlockHash = block.Hash()
					receipts = append(receipts, receipt)

					tc.ExpectedResult.Transactions = append(tc.ExpectedResult.Transactions,
						types.TransactionOrHash{
							Hash: state.HashPtr(tx.Hash()),
						},
					)

					batchTxs = append(batchTxs, *tx)
					effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
				}
				batchL2Data, err := state.EncodeTransactions(batchTxs, effectivePercentages, forkID6)
				require.NoError(t, err)

				batch := &state.Batch{
					BatchNumber:    1,
					Coinbase:       common.HexToAddress("0x1"),
					StateRoot:      common.HexToHash("0x2"),
					AccInputHash:   common.HexToHash("0x3"),
					GlobalExitRoot: common.HexToHash("0x4"),
					Timestamp:      time.Unix(1, 0),
					BatchL2Data:    batchL2Data,
				}

				m.State.
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(batch, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, m.DbTx).
					Return(&ger, nil).
					Once()
				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(receipts[i], nil).
						Once()
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(batchTxs, effectivePercentages, nil).
					Once()

				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), m.DbTx).
					Return(blocks, nil).
					Once()

				tc.ExpectedResult.BatchL2Data = batchL2Data
			},
		},
		{
			Name:         "get latest batch successfully",
			Number:       "latest",
			WithTxDetail: true,
			ExpectedResult: &types.Batch{
				Number:              1,
				ForcedBatchNumber:   ptrArgUint64FromUint64(1),
				Coinbase:            common.HexToAddress("0x1"),
				StateRoot:           common.HexToHash("0x2"),
				AccInputHash:        common.HexToHash("0x3"),
				GlobalExitRoot:      common.HexToHash("0x4"),
				Timestamp:           1,
				SendSequencesTxHash: ptrHash(common.HexToHash("0x10")),
				VerifyBatchTxHash:   ptrHash(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastClosedBatchNumber", context.Background(), m.DbTx).
					Return(uint64(tc.ExpectedResult.Number), nil).
					Once()

				txs := []*ethTypes.Transaction{
					signTx(ethTypes.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
					signTx(ethTypes.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
				}

				batchTxs := make([]ethTypes.Transaction, 0, len(txs))
				effectivePercentages := make([]uint8, 0, len(txs))
				tc.ExpectedResult.Transactions = []types.TransactionOrHash{}

				receipts := []*ethTypes.Receipt{}
				blocks := []ethTypes.Block{}
				for i, tx := range txs {
					block := ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(i))}).WithBody([]*ethTypes.Transaction{tx}, []*ethTypes.Header{})
					blocks = append(blocks, *block)
					receipt := ethTypes.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = block.Number()
					receipt.BlockHash = block.Hash()
					receipts = append(receipts, receipt)
					from, _ := state.GetSender(*tx)
					V, R, S := tx.RawSignatureValues()

					rpcReceipt, err := types.NewReceipt(*tx, receipt)
					require.NoError(t, err)

					tc.ExpectedResult.Transactions = append(tc.ExpectedResult.Transactions,
						types.TransactionOrHash{
							Tx: &types.Transaction{
								Nonce:       types.ArgUint64(tx.Nonce()),
								GasPrice:    types.ArgBig(*tx.GasPrice()),
								Gas:         types.ArgUint64(tx.Gas()),
								To:          tx.To(),
								Value:       types.ArgBig(*tx.Value()),
								Input:       tx.Data(),
								Hash:        tx.Hash(),
								From:        from,
								BlockNumber: ptrArgUint64FromUint64(block.NumberU64()),
								BlockHash:   ptrHash(receipt.BlockHash),
								TxIndex:     ptrArgUint64FromUint(receipt.TransactionIndex),
								ChainID:     types.ArgBig(*tx.ChainId()),
								Type:        types.ArgUint64(tx.Type()),
								V:           types.ArgBig(*V),
								R:           types.ArgBig(*R),
								S:           types.ArgBig(*S),
								Receipt:     &rpcReceipt,
							},
						},
					)

					batchTxs = append(batchTxs, *tx)
					effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
				}
				batchL2Data, err := state.EncodeTransactions(batchTxs, effectivePercentages, forkID6)
				require.NoError(t, err)
				var fb uint64 = 1
				batch := &state.Batch{
					BatchNumber:    1,
					ForcedBatchNum: &fb,
					Coinbase:       common.HexToAddress("0x1"),
					StateRoot:      common.HexToHash("0x2"),
					AccInputHash:   common.HexToHash("0x3"),
					GlobalExitRoot: common.HexToHash("0x4"),
					Timestamp:      time.Unix(1, 0),
					BatchL2Data:    batchL2Data,
				}

				m.State.
					On("GetBatchByNumber", context.Background(), uint64(tc.ExpectedResult.Number), m.DbTx).
					Return(batch, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), uint64(tc.ExpectedResult.Number), m.DbTx).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), uint64(tc.ExpectedResult.Number), m.DbTx).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, m.DbTx).
					Return(&ger, nil).
					Once()

				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(receipts[i], nil).
						Once()
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), uint64(tc.ExpectedResult.Number), m.DbTx).
					Return(batchTxs, effectivePercentages, nil).
					Once()
				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), uint64(tc.ExpectedResult.Number), m.DbTx).
					Return(blocks, nil).
					Once()
				tc.ExpectedResult.BatchL2Data = batchL2Data
			},
		},
		{
			Name:           "get latest batch fails to compute batch number",
			Number:         "latest",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last batch number from state"),
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastClosedBatchNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
		{
			Name:           "get latest batch fails to load batch by number",
			Number:         "latest",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "couldn't load batch from state by number 1"),
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastClosedBatchNumber", context.Background(), m.DbTx).
					Return(uint64(1), nil).
					Once()

				m.State.
					On("GetBatchByNumber", context.Background(), uint64(1), m.DbTx).
					Return(nil, errors.New("failed to load batch by number")).
					Once()
			},
		},
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(s, m, &tc)

			res, err := s.JSONRPCCall("zkevm_getBatchByNumber", tc.Number, tc.WithTxDetail)
			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result interface{}
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				if result != nil || testCase.ExpectedResult != nil {
					var batch map[string]interface{}
					err = json.Unmarshal(res.Result, &batch)
					require.NoError(t, err)
					assert.Equal(t, tc.ExpectedResult.Number.Hex(), batch["number"].(string))
					if tc.ExpectedResult.ForcedBatchNumber != nil {
						assert.Equal(t, tc.ExpectedResult.ForcedBatchNumber.Hex(), batch["forcedBatchNumber"].(string))
					}
					assert.Equal(t, tc.ExpectedResult.Coinbase.String(), batch["coinbase"].(string))
					assert.Equal(t, tc.ExpectedResult.StateRoot.String(), batch["stateRoot"].(string))
					assert.Equal(t, tc.ExpectedResult.GlobalExitRoot.String(), batch["globalExitRoot"].(string))
					assert.Equal(t, tc.ExpectedResult.LocalExitRoot.String(), batch["localExitRoot"].(string))
					assert.Equal(t, tc.ExpectedResult.AccInputHash.String(), batch["accInputHash"].(string))
					assert.Equal(t, tc.ExpectedResult.Timestamp.Hex(), batch["timestamp"].(string))
					assert.Equal(t, tc.ExpectedResult.SendSequencesTxHash.String(), batch["sendSequencesTxHash"].(string))
					assert.Equal(t, tc.ExpectedResult.VerifyBatchTxHash.String(), batch["verifyBatchTxHash"].(string))
					batchTxs := batch["transactions"].([]interface{})
					for i, txOrHash := range tc.ExpectedResult.Transactions {
						switch batchTxOrHash := batchTxs[i].(type) {
						case string:
							assert.Equal(t, txOrHash.Hash.String(), batchTxOrHash)
						case map[string]interface{}:
							tx := txOrHash.Tx
							assert.Equal(t, tx.Nonce.Hex(), batchTxOrHash["nonce"].(string))
							assert.Equal(t, tx.GasPrice.Hex(), batchTxOrHash["gasPrice"].(string))
							assert.Equal(t, tx.Gas.Hex(), batchTxOrHash["gas"].(string))
							assert.Equal(t, tx.To.String(), batchTxOrHash["to"].(string))
							assert.Equal(t, tx.Value.Hex(), batchTxOrHash["value"].(string))
							assert.Equal(t, tx.Input.Hex(), batchTxOrHash["input"].(string))
							assert.Equal(t, tx.V.Hex(), batchTxOrHash["v"].(string))
							assert.Equal(t, tx.R.Hex(), batchTxOrHash["r"].(string))
							assert.Equal(t, tx.S.Hex(), batchTxOrHash["s"].(string))
							assert.Equal(t, tx.Hash.String(), batchTxOrHash["hash"].(string))
							assert.Equal(t, strings.ToLower(tx.From.String()), strings.ToLower(batchTxOrHash["from"].(string)))
							assert.Equal(t, tx.BlockHash.String(), batchTxOrHash["blockHash"].(string))
							assert.Equal(t, tx.BlockNumber.Hex(), batchTxOrHash["blockNumber"].(string))
							assert.Equal(t, tx.TxIndex.Hex(), batchTxOrHash["transactionIndex"].(string))
							assert.Equal(t, tx.ChainID.Hex(), batchTxOrHash["chainId"].(string))
							assert.Equal(t, tx.Type.Hex(), batchTxOrHash["type"].(string))
						}
					}
					expectedBatchL2DataHex := "0x" + common.Bytes2Hex(testCase.ExpectedResult.BatchL2Data)
					assert.Equal(t, expectedBatchL2DataHex, batch["batchL2Data"].(string))
				}
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetL2FullBlockByHash(t *testing.T) {
	type testCase struct {
		Name           string
		Hash           common.Hash
		ExpectedResult *ethTypes.Block
		ExpectedError  interface{}
		SetupMocks     func(*mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "Block not found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:           "Failed get block from state",
			Hash:           common.HexToHash("0x234"),
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get block by hash from state"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to get block from state")).
					Once()
			},
		},
		{
			Name: "get block successfully",
			Hash: common.HexToHash("0x345"),
			ExpectedResult: ethTypes.NewBlock(
				&ethTypes.Header{Number: big.NewInt(1), UncleHash: ethTypes.EmptyUncleHash, Root: ethTypes.EmptyRootHash},
				[]*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				block := ethTypes.NewBlock(ethTypes.CopyHeader(tc.ExpectedResult.Header()), tc.ExpectedResult.Transactions(), tc.ExpectedResult.Uncles(), []*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))}, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(block, nil).
					Once()

				for _, tx := range tc.ExpectedResult.Transactions() {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil).
						Once()
				}
			},
		},
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			res, err := s.JSONRPCCall("zkevm_getFullBlockByHash", tc.Hash.String())
			require.NoError(t, err)

			if tc.ExpectedResult != nil {
				require.NotNil(t, res.Result)
				require.Nil(t, res.Error)

				var result types.Block
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				assert.Equal(t, tc.ExpectedResult.Number().Uint64(), uint64(result.Number))
				assert.Equal(t, len(tc.ExpectedResult.Transactions()), len(result.Transactions))
				assert.Equal(t, tc.ExpectedResult.Hash(), result.Hash)
			}

			if tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*types.RPCError); ok {
					assert.Equal(t, expectedErr.ErrorCode(), res.Error.Code)
					assert.Equal(t, expectedErr.Error(), res.Error.Message)
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetL2FullBlockByNumber(t *testing.T) {
	type testCase struct {
		Name           string
		Number         string
		ExpectedResult *ethTypes.Block
		ExpectedError  interface{}
		SetupMocks     func(*mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "Block not found",
			Number:         "0x7B",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), hex.DecodeUint64(tc.Number), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:   "get specific block successfully",
			Number: "0x159",
			ExpectedResult: ethTypes.NewBlock(
				&ethTypes.Header{Number: big.NewInt(1), UncleHash: ethTypes.EmptyUncleHash, Root: ethTypes.EmptyRootHash},
				[]*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				block := ethTypes.NewBlock(ethTypes.CopyHeader(tc.ExpectedResult.Header()), tc.ExpectedResult.Transactions(),
					tc.ExpectedResult.Uncles(), []*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))}, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), hex.DecodeUint64(tc.Number), m.DbTx).
					Return(block, nil).
					Once()

				for _, tx := range tc.ExpectedResult.Transactions() {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil).
						Once()
				}
			},
		},
		{
			Name:   "get latest block successfully",
			Number: "latest",
			ExpectedResult: ethTypes.NewBlock(
				&ethTypes.Header{Number: big.NewInt(2), UncleHash: ethTypes.EmptyUncleHash, Root: ethTypes.EmptyRootHash},
				[]*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				blockNumber := uint64(1)

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), blockNumber, m.DbTx).
					Return(tc.ExpectedResult, nil).
					Once()

				for _, tx := range tc.ExpectedResult.Transactions() {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil).
						Once()
				}
			},
		},
		{
			Name:           "get latest block fails to compute block number",
			Number:         "latest",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last block number from state"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "get latest block fails to load block by number",
			Number:         "latest",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "couldn't load block from state by number 1"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(1), nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), uint64(1), m.DbTx).
					Return(nil, errors.New("failed to load block by number")).
					Once()
			},
		},
		{
			Name:           "get pending block successfully",
			Number:         "pending",
			ExpectedResult: ethTypes.NewBlock(&ethTypes.Header{Number: big.NewInt(2)}, nil, nil, nil, &trie.StackTrie{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				lastBlockHeader := ethTypes.CopyHeader(tc.ExpectedResult.Header())
				lastBlockHeader.Number.Sub(lastBlockHeader.Number, big.NewInt(1))
				lastBlock := ethTypes.NewBlock(lastBlockHeader, nil, nil, nil, &trie.StackTrie{})

				expectedResultHeader := ethTypes.CopyHeader(tc.ExpectedResult.Header())
				expectedResultHeader.ParentHash = lastBlock.Hash()
				tc.ExpectedResult = ethTypes.NewBlock(expectedResultHeader, nil, nil, nil, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2Block", context.Background(), m.DbTx).
					Return(lastBlock, nil).
					Once()
			},
		},
		{
			Name:           "get pending block fails",
			Number:         "pending",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "couldn't load last block from state to compute the pending block"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2Block", context.Background(), m.DbTx).
					Return(nil, errors.New("failed to load last block")).
					Once()
			},
		},
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			res, err := s.JSONRPCCall("zkevm_getFullBlockByNumber", tc.Number)
			require.NoError(t, err)

			if tc.ExpectedResult != nil {
				require.NotNil(t, res.Result)
				require.Nil(t, res.Error)

				var result types.Block
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				assert.Equal(t, tc.ExpectedResult.Number().Uint64(), uint64(result.Number))
				assert.Equal(t, len(tc.ExpectedResult.Transactions()), len(result.Transactions))
				assert.Equal(t, tc.ExpectedResult.Hash(), result.Hash)
			}

			if tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*types.RPCError); ok {
					assert.Equal(t, expectedErr.ErrorCode(), res.Error.Code)
					assert.Equal(t, expectedErr.Error(), res.Error.Message)
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetNativeBlockHashesInRange(t *testing.T) {
	type testCase struct {
		Name           string
		Filter         NativeBlockHashBlockRangeFilter
		ExpectedResult *[]string
		ExpectedError  interface{}
		SetupMocks     func(*mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name: "Block not found",
			Filter: NativeBlockHashBlockRangeFilter{
				FromBlock: types.BlockNumber(0),
				ToBlock:   types.BlockNumber(10),
			},
			ExpectedResult: ptr([]string{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				fromBlock, _ := tc.Filter.FromBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)
				toBlock, _ := tc.Filter.ToBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)

				m.State.
					On("GetNativeBlockHashesInRange", context.Background(), fromBlock, toBlock, m.DbTx).
					Return([]common.Hash{}, nil).
					Once()
			},
		},
		{
			Name: "native block hash range returned successfully",
			Filter: NativeBlockHashBlockRangeFilter{
				FromBlock: types.BlockNumber(0),
				ToBlock:   types.BlockNumber(10),
			},
			ExpectedResult: ptr([]string{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				fromBlock, _ := tc.Filter.FromBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)
				toBlock, _ := tc.Filter.ToBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)
				hashes := []common.Hash{}
				expectedResult := []string{}
				for i := fromBlock; i < toBlock; i++ {
					sHash := hex.EncodeUint64(i)
					hash := common.HexToHash(sHash)
					hashes = append(hashes, hash)
					expectedResult = append(expectedResult, hash.String())
				}
				tc.ExpectedResult = &expectedResult

				m.State.
					On("GetNativeBlockHashesInRange", context.Background(), fromBlock, toBlock, m.DbTx).
					Return(hashes, nil).
					Once()
			},
		},
		{
			Name: "native block hash range fails due to invalid range",
			Filter: NativeBlockHashBlockRangeFilter{
				FromBlock: types.BlockNumber(10),
				ToBlock:   types.BlockNumber(0),
			},
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.InvalidParamsErrorCode, "invalid block range"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()
			},
		},
		{
			Name: "native block hash range fails due to range limit",
			Filter: NativeBlockHashBlockRangeFilter{
				FromBlock: types.BlockNumber(0),
				ToBlock:   types.BlockNumber(60001),
			},
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.InvalidParamsErrorCode, "native block hashes are limited to a 60000 block range"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()
			},
		},
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			res, err := s.JSONRPCCall("zkevm_getNativeBlockHashesInRange", tc.Filter)
			require.NoError(t, err)

			if tc.ExpectedResult != nil {
				require.NotNil(t, res.Result)
				require.Nil(t, res.Error)

				var result []string
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				assert.Equal(t, len(*tc.ExpectedResult), len(result))
				assert.ElementsMatch(t, *tc.ExpectedResult, result)
			}

			if tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*types.RPCError); ok {
					assert.Equal(t, expectedErr.ErrorCode(), res.Error.Code)
					assert.Equal(t, expectedErr.Error(), res.Error.Message)
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func ptrUint64(n uint64) *uint64 {
	return &n
}

func ptrArgUint64FromUint(n uint) *types.ArgUint64 {
	tmp := types.ArgUint64(n)
	return &tmp
}

func ptrArgUint64FromUint64(n uint64) *types.ArgUint64 {
	tmp := types.ArgUint64(n)
	return &tmp
}

func ptrHash(h common.Hash) *common.Hash {
	return &h
}

func ptr[T any](v T) *T {
	return &v
}

func signTx(tx *ethTypes.Transaction, chainID uint64) *ethTypes.Transaction {
	privateKey, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(chainID))
	signedTx, _ := auth.Signer(auth.From, tx)
	return signedTx
}
