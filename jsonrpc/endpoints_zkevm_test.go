package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			ExpectedResult: state.Ptr(uint64(10)),
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("GetLastConsolidatedL2BlockNumber", context.Background(), nil).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get consolidated block number",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get last consolidated block number from state"),
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("GetLastConsolidatedL2BlockNumber", context.Background(), nil).
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
				m.State.
					On("IsL2BlockConsolidated", context.Background(), uint64(1), nil).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the consolidation status",
			ExpectedResult: false,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to check if the block is consolidated"),
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("IsL2BlockConsolidated", context.Background(), uint64(1), nil).
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
				m.State.
					On("IsL2BlockVirtualized", context.Background(), uint64(1), nil).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the virtualization status",
			ExpectedResult: false,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to check if the block is virtualized"),
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("IsL2BlockVirtualized", context.Background(), uint64(1), nil).
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
				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, nil).
					Return(batchNumber, nil).
					Once()
			},
		},
		{
			Name:           "failed to get batch number",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get batch number from block number"),
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, nil).
					Return(uint64(0), errors.New("failed to get batch number of l2 batchNum")).
					Once()
			},
		},
		{
			Name:           "batch number not found",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, nil).
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
				m.State.
					On("GetLastBatchNumber", context.Background(), nil).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("GetLastBatchNumber", context.Background(), nil).
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
				m.State.
					On("GetLastVirtualBatchNum", context.Background(), nil).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get virtual batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last virtual batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("GetLastVirtualBatchNum", context.Background(), nil).
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
				m.State.
					On("GetLastVerifiedBatch", context.Background(), nil).
					Return(&state.VerifiedBatch{BatchNumber: uint64(10)}, nil).
					Once()
			},
		},
		{
			Name:           "failed to get verified batch number",
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get the last verified batch number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocksWrapper) {
				m.State.
					On("GetLastVerifiedBatch", context.Background(), nil).
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
				m.State.
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
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
				SendSequencesTxHash: state.Ptr(common.HexToHash("0x10")),
				VerifyBatchTxHash:   state.Ptr(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				txs := []*ethTypes.Transaction{
					signTx(ethTypes.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
					signTx(ethTypes.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
				}

				batchTxs := make([]ethTypes.Transaction, 0, len(txs))
				effectivePercentages := make([]uint8, 0, len(txs))
				tc.ExpectedResult.Transactions = []types.TransactionOrHash{}
				receipts := []*ethTypes.Receipt{}
				blocks := []state.L2Block{}
				for i, tx := range txs {
					block := state.NewL2BlockWithHeader(state.NewL2Header(&ethTypes.Header{Number: big.NewInt(int64(i))})).WithBody([]*ethTypes.Transaction{tx}, []*state.L2Header{})
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
								BlockHash:   state.Ptr(receipt.BlockHash),
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
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(batch, nil).
					Once()

				m.State.
					On("GetBatchTimestamp", mock.Anything, mock.Anything, (*uint64)(nil), nil).
					Return(&batch.Timestamp, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, nil).
					Return(&ger, nil).
					Once()

				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), nil).
						Return(receipts[i], nil).
						Once()
					m.State.
						On("GetL2TxHashByTxHash", context.Background(), tx.Hash(), nil).
						Return(state.Ptr(tx.Hash()), nil).
						Once()
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(batchTxs, effectivePercentages, nil).
					Once()

				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
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
				SendSequencesTxHash: state.Ptr(common.HexToHash("0x10")),
				VerifyBatchTxHash:   state.Ptr(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				txs := []*ethTypes.Transaction{
					signTx(ethTypes.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
					signTx(ethTypes.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.ChainID()),
				}

				batchTxs := make([]ethTypes.Transaction, 0, len(txs))
				effectivePercentages := make([]uint8, 0, len(txs))
				tc.ExpectedResult.Transactions = []types.TransactionOrHash{}

				receipts := []*ethTypes.Receipt{}
				blocks := []state.L2Block{}
				for i, tx := range txs {
					block := state.NewL2BlockWithHeader(state.NewL2Header(&ethTypes.Header{Number: big.NewInt(int64(i))})).WithBody([]*ethTypes.Transaction{tx}, []*state.L2Header{})
					blocks = append(blocks, *block)
					receipt := ethTypes.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = block.Number()
					receipt.BlockHash = block.Hash()
					receipts = append(receipts, receipt)

					tc.ExpectedResult.Transactions = append(tc.ExpectedResult.Transactions,
						types.TransactionOrHash{
							Hash: state.Ptr(tx.Hash()),
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
					On("GetBatchByNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(batch, nil).
					Once()

				m.State.
					On("GetBatchTimestamp", mock.Anything, mock.Anything, (*uint64)(nil), nil).
					Return(&batch.Timestamp, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, nil).
					Return(&ger, nil).
					Once()
				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), nil).
						Return(receipts[i], nil).
						Once()
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
					Return(batchTxs, effectivePercentages, nil).
					Once()

				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), hex.DecodeBig(tc.Number).Uint64(), nil).
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
				SendSequencesTxHash: state.Ptr(common.HexToHash("0x10")),
				VerifyBatchTxHash:   state.Ptr(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetLastClosedBatchNumber", context.Background(), nil).
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
				blocks := []state.L2Block{}
				for i, tx := range txs {
					block := state.NewL2BlockWithHeader(state.NewL2Header(&ethTypes.Header{Number: big.NewInt(int64(i))})).WithBody([]*ethTypes.Transaction{tx}, []*state.L2Header{})
					blocks = append(blocks, *block)
					receipt := ethTypes.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = block.Number()
					receipt.BlockHash = block.Hash()
					receipts = append(receipts, receipt)
					from, _ := state.GetSender(*tx)
					V, R, S := tx.RawSignatureValues()
					l2Hash := common.HexToHash("0x987654321")

					rpcReceipt, err := types.NewReceipt(*tx, receipt, &l2Hash)
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
								BlockHash:   state.Ptr(receipt.BlockHash),
								TxIndex:     ptrArgUint64FromUint(receipt.TransactionIndex),
								ChainID:     types.ArgBig(*tx.ChainId()),
								Type:        types.ArgUint64(tx.Type()),
								V:           types.ArgBig(*V),
								R:           types.ArgBig(*R),
								S:           types.ArgBig(*S),
								Receipt:     &rpcReceipt,
								L2Hash:      &l2Hash,
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
					On("GetBatchByNumber", context.Background(), uint64(tc.ExpectedResult.Number), nil).
					Return(batch, nil).
					Once()

				m.State.
					On("GetBatchTimestamp", mock.Anything, mock.Anything, (*uint64)(nil), nil).
					Return(&batch.Timestamp, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), uint64(tc.ExpectedResult.Number), nil).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), uint64(tc.ExpectedResult.Number), nil).
					Return(verifiedBatch, nil).
					Once()

				ger := state.GlobalExitRoot{
					MainnetExitRoot: common.HexToHash("0x4"),
					RollupExitRoot:  common.HexToHash("0x4"),
					GlobalExitRoot:  common.HexToHash("0x4"),
				}
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), batch.GlobalExitRoot, nil).
					Return(&ger, nil).
					Once()

				for i, tx := range txs {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), nil).
						Return(receipts[i], nil).
						Once()

					m.State.
						On("GetL2TxHashByTxHash", context.Background(), tx.Hash(), nil).
						Return(state.Ptr(tx.Hash()), nil).
						Once()
				}

				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), uint64(tc.ExpectedResult.Number), nil).
					Return(batchTxs, effectivePercentages, nil).
					Once()
				m.State.
					On("GetL2BlocksByBatchNumber", context.Background(), uint64(tc.ExpectedResult.Number), nil).
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
				m.State.
					On("GetLastClosedBatchNumber", context.Background(), nil).
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
				m.State.
					On("GetLastClosedBatchNumber", context.Background(), nil).
					Return(uint64(1), nil).
					Once()

				m.State.
					On("GetBatchByNumber", context.Background(), uint64(1), nil).
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

	st := trie.NewStackTrie(nil)
	testCases := []testCase{
		{
			Name:           "Block not found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:           "Failed get block from state",
			Hash:           common.HexToHash("0x234"),
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get block by hash from state"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, nil).
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
				st,
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				uncles := make([]*state.L2Header, 0, len(tc.ExpectedResult.Uncles()))
				for _, uncle := range tc.ExpectedResult.Uncles() {
					uncles = append(uncles, state.NewL2Header(uncle))
				}
				st := trie.NewStackTrie(nil)
				block := state.NewL2Block(state.NewL2Header(tc.ExpectedResult.Header()), tc.ExpectedResult.Transactions(), uncles, []*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))}, st)

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, nil).
					Return(block, nil).
					Once()

				for _, tx := range tc.ExpectedResult.Transactions() {
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), nil).
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
				assert.Equal(t, state.Ptr(tc.ExpectedResult.Hash()), result.Hash)
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
		ExpectedResult *types.Block
		ExpectedError  *types.RPCError
		SetupMocks     func(*mocksWrapper, *testCase)
	}

	transactions := []*ethTypes.Transaction{
		ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    1,
			GasPrice: big.NewInt(2),
			Gas:      3,
			To:       state.Ptr(common.HexToAddress("0x4")),
			Value:    big.NewInt(5),
			Data:     types.ArgBytes{6},
		}),
		ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    2,
			GasPrice: big.NewInt(3),
			Gas:      4,
			To:       state.Ptr(common.HexToAddress("0x5")),
			Value:    big.NewInt(6),
			Data:     types.ArgBytes{7},
		}),
	}

	auth := operations.MustGetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	var signedTransactions []*ethTypes.Transaction
	for _, tx := range transactions {
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		signedTransactions = append(signedTransactions, signedTx)
	}

	uncles := []*state.L2Header{
		state.NewL2Header(&ethTypes.Header{Number: big.NewInt(222)}),
		state.NewL2Header(&ethTypes.Header{Number: big.NewInt(333)}),
	}

	receipts := []*ethTypes.Receipt{}
	for _, tx := range signedTransactions {
		receipts = append(receipts, &ethTypes.Receipt{
			TxHash: tx.Hash(),
		})
	}

	header := &ethTypes.Header{
		ParentHash:  common.HexToHash("0x1"),
		UncleHash:   common.HexToHash("0x2"),
		Coinbase:    common.HexToAddress("0x3"),
		Root:        common.HexToHash("0x4"),
		TxHash:      common.HexToHash("0x5"),
		ReceiptHash: common.HexToHash("0x6"),
		Difficulty:  big.NewInt(8),
		Number:      big.NewInt(9),
		GasLimit:    10,
		GasUsed:     11,
		Time:        12,
		Extra:       types.ArgBytes{13},
		MixDigest:   common.HexToHash("0x14"),
		Nonce:       ethTypes.EncodeNonce(15),
		Bloom:       ethTypes.CreateBloom(receipts),
	}

	l2Header := state.NewL2Header(header)
	l2Header.GlobalExitRoot = common.HexToHash("0x16")
	l2Header.BlockInfoRoot = common.HexToHash("0x17")
	st := trie.NewStackTrie(nil)
	l2Block := state.NewL2Block(l2Header, signedTransactions, uncles, receipts, st)

	for _, receipt := range receipts {
		receipt.BlockHash = l2Block.Hash()
		receipt.BlockNumber = l2Block.Number()
	}

	rpcTransactions := []types.TransactionOrHash{}
	for _, tx := range signedTransactions {
		sender, _ := state.GetSender(*tx)
		rpcTransactions = append(rpcTransactions,
			types.TransactionOrHash{
				Tx: &types.Transaction{
					Nonce:    types.ArgUint64(tx.Nonce()),
					GasPrice: types.ArgBig(*tx.GasPrice()),
					Gas:      types.ArgUint64(tx.Gas()),
					To:       tx.To(),
					Value:    types.ArgBig(*tx.Value()),
					Input:    tx.Data(),

					Hash:        tx.Hash(),
					From:        sender,
					BlockHash:   state.Ptr(l2Block.Hash()),
					BlockNumber: state.Ptr(types.ArgUint64(l2Block.Number().Uint64())),
				},
			})
	}

	rpcUncles := []common.Hash{}
	for _, uncle := range uncles {
		rpcUncles = append(rpcUncles, uncle.Hash())
	}

	var miner *common.Address
	if l2Block.Coinbase().String() != state.ZeroAddress.String() {
		cb := l2Block.Coinbase()
		miner = &cb
	}

	n := big.NewInt(0).SetUint64(l2Block.Nonce())
	rpcBlockNonce := types.ArgBytes(common.LeftPadBytes(n.Bytes(), 8)) //nolint:gomnd

	difficulty := types.ArgUint64(0)
	var totalDifficulty *types.ArgUint64
	if l2Block.Difficulty() != nil {
		difficulty = types.ArgUint64(l2Block.Difficulty().Uint64())
		totalDifficulty = &difficulty
	}

	rpcBlock := &types.Block{
		ParentHash:      l2Block.ParentHash(),
		Sha3Uncles:      l2Block.UncleHash(),
		Miner:           miner,
		StateRoot:       l2Block.Root(),
		TxRoot:          l2Block.TxHash(),
		ReceiptsRoot:    l2Block.ReceiptHash(),
		LogsBloom:       ethTypes.CreateBloom(receipts),
		Difficulty:      difficulty,
		TotalDifficulty: totalDifficulty,
		Size:            types.ArgUint64(l2Block.Size()),
		Number:          types.ArgUint64(l2Block.NumberU64()),
		GasLimit:        types.ArgUint64(l2Block.GasLimit()),
		GasUsed:         types.ArgUint64(l2Block.GasUsed()),
		Timestamp:       types.ArgUint64(l2Block.Time()),
		ExtraData:       l2Block.Extra(),
		MixHash:         l2Block.MixDigest(),
		Nonce:           &rpcBlockNonce,
		Hash:            state.Ptr(l2Block.Hash()),
		GlobalExitRoot:  state.Ptr(l2Block.GlobalExitRoot()),
		BlockInfoRoot:   state.Ptr(l2Block.BlockInfoRoot()),
		Uncles:          rpcUncles,
		Transactions:    rpcTransactions,
	}

	testCases := []testCase{
		{
			Name:           "Block not found",
			Number:         "0x7B",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetL2BlockByNumber", context.Background(), hex.DecodeUint64(tc.Number), nil).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "get specific block successfully",
			Number:         "0x159",
			ExpectedResult: rpcBlock,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetL2BlockByNumber", context.Background(), hex.DecodeUint64(tc.Number), nil).
					Return(l2Block, nil).
					Once()

				for _, receipt := range receipts {
					m.State.
						On("GetTransactionReceipt", context.Background(), receipt.TxHash, nil).
						Return(receipt, nil).
						Once()
				}
			},
		},
		{
			Name:           "get latest block successfully",
			Number:         "latest",
			ExpectedResult: rpcBlock,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				blockNumber := uint64(1)

				m.State.
					On("GetLastL2BlockNumber", context.Background(), nil).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), blockNumber, nil).
					Return(l2Block, nil).
					Once()

				for _, receipt := range receipts {
					m.State.
						On("GetTransactionReceipt", context.Background(), receipt.TxHash, nil).
						Return(receipt, nil).
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
				m.State.
					On("GetLastL2BlockNumber", context.Background(), nil).
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
				m.State.
					On("GetLastL2BlockNumber", context.Background(), nil).
					Return(uint64(1), nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), uint64(1), nil).
					Return(nil, errors.New("failed to load block by number")).
					Once()
			},
		},
		{
			Name:          "get pending block successfully",
			Number:        "pending",
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				lastBlockHeader := &ethTypes.Header{Number: big.NewInt(0).SetUint64(uint64(rpcBlock.Number))}
				lastBlockHeader.Number.Sub(lastBlockHeader.Number, big.NewInt(1))
				st := trie.NewStackTrie(nil)
				lastBlock := state.NewL2Block(state.NewL2Header(lastBlockHeader), nil, nil, nil, st)

				tc.ExpectedResult = &types.Block{}
				tc.ExpectedResult.ParentHash = lastBlock.Hash()
				tc.ExpectedResult.Number = types.ArgUint64(lastBlock.Number().Uint64() + 1)
				tc.ExpectedResult.TxRoot = ethTypes.EmptyRootHash
				tc.ExpectedResult.Sha3Uncles = ethTypes.EmptyUncleHash
				tc.ExpectedResult.Size = 501
				tc.ExpectedResult.ExtraData = []byte{}
				tc.ExpectedResult.GlobalExitRoot = state.Ptr(common.Hash{})
				tc.ExpectedResult.BlockInfoRoot = state.Ptr(common.Hash{})
				tc.ExpectedResult.Hash = nil
				tc.ExpectedResult.Miner = nil
				tc.ExpectedResult.Nonce = nil
				tc.ExpectedResult.TotalDifficulty = nil

				m.State.
					On("GetLastL2Block", context.Background(), nil).
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
				m.State.
					On("GetLastL2Block", context.Background(), nil).
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

				assert.Equal(t, tc.ExpectedResult.ParentHash.String(), result.ParentHash.String())
				assert.Equal(t, tc.ExpectedResult.Sha3Uncles.String(), result.Sha3Uncles.String())
				assert.Equal(t, tc.ExpectedResult.StateRoot.String(), result.StateRoot.String())
				assert.Equal(t, tc.ExpectedResult.TxRoot.String(), result.TxRoot.String())
				assert.Equal(t, tc.ExpectedResult.ReceiptsRoot.String(), result.ReceiptsRoot.String())
				assert.Equal(t, tc.ExpectedResult.LogsBloom, result.LogsBloom)
				assert.Equal(t, tc.ExpectedResult.Difficulty, result.Difficulty)
				assert.Equal(t, tc.ExpectedResult.Size, result.Size)
				assert.Equal(t, tc.ExpectedResult.Number, result.Number)
				assert.Equal(t, tc.ExpectedResult.GasLimit, result.GasLimit)
				assert.Equal(t, tc.ExpectedResult.GasUsed, result.GasUsed)
				assert.Equal(t, tc.ExpectedResult.Timestamp, result.Timestamp)
				assert.Equal(t, tc.ExpectedResult.ExtraData, result.ExtraData)
				assert.Equal(t, tc.ExpectedResult.MixHash, result.MixHash)
				assert.Equal(t, tc.ExpectedResult.GlobalExitRoot, result.GlobalExitRoot)
				assert.Equal(t, tc.ExpectedResult.BlockInfoRoot, result.BlockInfoRoot)

				if tc.ExpectedResult.Hash != nil {
					assert.Equal(t, tc.ExpectedResult.Hash.String(), result.Hash.String())
				} else {
					assert.Nil(t, result.Hash)
				}
				if tc.ExpectedResult.Miner != nil {
					assert.Equal(t, tc.ExpectedResult.Miner.String(), result.Miner.String())
				} else {
					assert.Nil(t, result.Miner)
				}
				if tc.ExpectedResult.Nonce != nil {
					assert.Equal(t, tc.ExpectedResult.Nonce, result.Nonce)
				} else {
					assert.Nil(t, result.Nonce)
				}
				if tc.ExpectedResult.TotalDifficulty != nil {
					assert.Equal(t, tc.ExpectedResult.TotalDifficulty, result.TotalDifficulty)
				} else {
					assert.Nil(t, result.TotalDifficulty)
				}

				assert.Equal(t, len(tc.ExpectedResult.Transactions), len(result.Transactions))
				assert.Equal(t, len(tc.ExpectedResult.Uncles), len(result.Uncles))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				rpcErr := res.Error.RPCError()
				assert.Equal(t, tc.ExpectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.ExpectedError.Error(), rpcErr.Error())
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
			ExpectedResult: state.Ptr([]string{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				fromBlock, _ := tc.Filter.FromBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)
				toBlock, _ := tc.Filter.ToBlock.GetNumericBlockNumber(context.Background(), nil, nil, nil)

				m.State.
					On("GetNativeBlockHashesInRange", context.Background(), fromBlock, toBlock, nil).
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
			ExpectedResult: state.Ptr([]string{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
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
					On("GetNativeBlockHashesInRange", context.Background(), fromBlock, toBlock, nil).
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

func TestGetTransactionByL2Hash(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name            string
		Hash            common.Hash
		ExpectedPending bool
		ExpectedResult  *types.Transaction
		ExpectedError   *types.RPCError
		SetupMocks      func(m *mocksWrapper, tc testCase)
	}

	chainID := big.NewInt(1)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx := ethTypes.NewTransaction(1, common.HexToAddress("0x111"), big.NewInt(2), 3, big.NewInt(4), []byte{5, 6, 7, 8})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	blockHash := common.HexToHash("0x1")
	blockNumber := blockNumOne

	receipt := &ethTypes.Receipt{
		TxHash:           signedTx.Hash(),
		BlockHash:        blockHash,
		BlockNumber:      blockNumber,
		TransactionIndex: 0,
	}

	txV, txR, txS := signedTx.RawSignatureValues()

	l2Hash := common.HexToHash("0x987654321")

	rpcTransaction := types.Transaction{
		Nonce:    types.ArgUint64(signedTx.Nonce()),
		GasPrice: types.ArgBig(*signedTx.GasPrice()),
		Gas:      types.ArgUint64(signedTx.Gas()),
		To:       signedTx.To(),
		Value:    types.ArgBig(*signedTx.Value()),
		Input:    signedTx.Data(),

		Hash:        signedTx.Hash(),
		From:        auth.From,
		BlockHash:   state.Ptr(blockHash),
		BlockNumber: state.Ptr(types.ArgUint64(blockNumber.Uint64())),
		V:           types.ArgBig(*txV),
		R:           types.ArgBig(*txR),
		S:           types.ArgBig(*txS),
		TxIndex:     state.Ptr(types.ArgUint64(0)),
		ChainID:     types.ArgBig(*chainID),
		Type:        0,
		L2Hash:      state.Ptr(l2Hash),
	}

	testCases := []testCase{
		{
			Name:            "Get TX Successfully from state",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  &rpcTransaction,
			ExpectedError:   nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(receipt, nil).
					Once()

				m.State.
					On("GetL2TxHashByTxHash", context.Background(), signedTx.Hash(), nil).
					Return(&l2Hash, nil).
					Once()
			},
		},
		{
			Name:            "Get TX Successfully from pool",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: true,
			ExpectedResult:  &rpcTransaction,
			ExpectedError:   nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				tc.ExpectedResult.BlockHash = nil
				tc.ExpectedResult.BlockNumber = nil
				tc.ExpectedResult.TxIndex = nil
				tc.ExpectedResult.L2Hash = nil

				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash).
					Return(&pool.Transaction{Transaction: *signedTx, Status: pool.TxStatusPending}, nil).
					Once()
			},
		},
		{
			Name:            "TX Not Found",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash).
					Return(nil, pool.ErrNotFound).
					Once()
			},
		},
		{
			Name:            "TX failed to load from the state",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   types.NewRPCError(types.DefaultErrorCode, "failed to load transaction by l2 hash from state"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, errors.New("failed to load transaction by l2 hash from state")).
					Once()
			},
		},
		{
			Name:            "TX failed to load from the pool",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   types.NewRPCError(types.DefaultErrorCode, "failed to load transaction by l2 hash from pool"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash).
					Return(nil, errors.New("failed to load transaction by l2 hash from pool")).
					Once()
			},
		},
		{
			Name:            "TX receipt Not Found",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   types.NewRPCError(types.DefaultErrorCode, "transaction receipt not found"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:            "TX receipt failed to load",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   types.NewRPCError(types.DefaultErrorCode, "failed to load transaction receipt from state"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(nil, errors.New("failed to load transaction receipt from state")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("zkevm_getTransactionByL2Hash", tc.Hash.String())
			require.NoError(t, err)

			if testCase.ExpectedResult != nil {
				require.NotNil(t, res.Result)
				require.Nil(t, res.Error)

				var result types.Transaction
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				require.Equal(t, tc.ExpectedResult.Nonce, result.Nonce)
				require.Equal(t, tc.ExpectedResult.GasPrice, result.GasPrice)
				require.Equal(t, tc.ExpectedResult.Gas, result.Gas)
				require.Equal(t, tc.ExpectedResult.To, result.To)
				require.Equal(t, tc.ExpectedResult.Value, result.Value)
				require.Equal(t, tc.ExpectedResult.Input, result.Input)

				require.Equal(t, tc.ExpectedResult.Hash, result.Hash)
				require.Equal(t, tc.ExpectedResult.From, result.From)
				require.Equal(t, tc.ExpectedResult.BlockHash, result.BlockHash)
				require.Equal(t, tc.ExpectedResult.BlockNumber, result.BlockNumber)
				require.Equal(t, tc.ExpectedResult.V, result.V)
				require.Equal(t, tc.ExpectedResult.R, result.R)
				require.Equal(t, tc.ExpectedResult.S, result.S)
				require.Equal(t, tc.ExpectedResult.TxIndex, result.TxIndex)
				require.Equal(t, tc.ExpectedResult.ChainID, result.ChainID)
				require.Equal(t, tc.ExpectedResult.Type, result.Type)
				require.Equal(t, tc.ExpectedResult.L2Hash, result.L2Hash)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				rpcErr := res.Error.RPCError()
				assert.Equal(t, tc.ExpectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.ExpectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func TestGetTransactionReceiptByL2Hash(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Hash           common.Hash
		ExpectedResult *types.Receipt
		ExpectedError  *types.RPCError
		SetupMocks     func(m *mocksWrapper, tc testCase)
	}

	chainID := big.NewInt(1)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx := ethTypes.NewTransaction(1, common.HexToAddress("0x111"), big.NewInt(2), 3, big.NewInt(4), []byte{5, 6, 7, 8})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	l2Hash := common.HexToHash("0x987654321")

	log := &ethTypes.Log{Topics: []common.Hash{common.HexToHash("0x1")}, Data: []byte{}}
	logs := []*ethTypes.Log{log}

	stateRoot := common.HexToHash("0x112233")

	receipt := &ethTypes.Receipt{
		Type:              signedTx.Type(),
		PostState:         stateRoot.Bytes(),
		CumulativeGasUsed: 1,
		BlockNumber:       big.NewInt(2),
		GasUsed:           3,
		TxHash:            signedTx.Hash(),
		TransactionIndex:  4,
		ContractAddress:   common.HexToAddress("0x223344"),
		Logs:              logs,
		Status:            ethTypes.ReceiptStatusSuccessful,
		EffectiveGasPrice: big.NewInt(5),
		BlobGasUsed:       6,
		BlobGasPrice:      big.NewInt(7),
		BlockHash:         common.HexToHash("0x1"),
	}

	receipt.Bloom = ethTypes.CreateBloom(ethTypes.Receipts{receipt})

	rpcReceipt := types.Receipt{
		Root:              &stateRoot,
		CumulativeGasUsed: types.ArgUint64(receipt.CumulativeGasUsed),
		LogsBloom:         receipt.Bloom,
		Logs:              receipt.Logs,
		Status:            types.ArgUint64(receipt.Status),
		TxHash:            receipt.TxHash,
		TxL2Hash:          &l2Hash,
		TxIndex:           types.ArgUint64(receipt.TransactionIndex),
		BlockHash:         receipt.BlockHash,
		BlockNumber:       types.ArgUint64(receipt.BlockNumber.Uint64()),
		GasUsed:           types.ArgUint64(receipt.GasUsed),
		FromAddr:          auth.From,
		ToAddr:            signedTx.To(),
		ContractAddress:   state.Ptr(receipt.ContractAddress),
		Type:              types.ArgUint64(receipt.Type),
		EffectiveGasPrice: state.Ptr(types.ArgBig(*receipt.EffectiveGasPrice)),
	}

	testCases := []testCase{
		{
			Name:           "Get TX receipt Successfully",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: &rpcReceipt,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(receipt, nil).
					Once()

				m.State.
					On("GetL2TxHashByTxHash", context.Background(), signedTx.Hash(), nil).
					Return(&l2Hash, nil).
					Once()
			},
		},
		{
			Name:           "Get TX receipt but tx not found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get TX receipt but failed to get tx",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get tx from state"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(nil, errors.New("failed to get tx")).
					Once()
			},
		},
		{
			Name:           "TX receipt Not Found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "TX receipt failed to load",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to get tx receipt from state"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(nil, errors.New("failed to get tx receipt from state")).
					Once()
			},
		},
		{
			Name:           "Get TX but failed to build response Successfully",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "failed to build the receipt response"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.
					On("GetTransactionByL2Hash", context.Background(), tc.Hash, nil).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, nil).
					Return(ethTypes.NewReceipt([]byte{}, false, 0), nil).
					Once()

				m.State.
					On("GetL2TxHashByTxHash", context.Background(), tx.Hash(), nil).
					Return(&l2Hash, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("zkevm_getTransactionReceiptByL2Hash", tc.Hash.String())
			require.NoError(t, err)

			if testCase.ExpectedResult != nil {
				require.NotNil(t, res.Result)
				require.Nil(t, res.Error)

				var result types.Receipt
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				assert.Equal(t, rpcReceipt.Root.String(), result.Root.String())
				assert.Equal(t, rpcReceipt.CumulativeGasUsed, result.CumulativeGasUsed)
				assert.Equal(t, rpcReceipt.LogsBloom, result.LogsBloom)
				assert.Equal(t, len(rpcReceipt.Logs), len(result.Logs))
				for i := 0; i < len(rpcReceipt.Logs); i++ {
					assert.Equal(t, rpcReceipt.Logs[i].Address, result.Logs[i].Address)
					assert.Equal(t, rpcReceipt.Logs[i].Topics, result.Logs[i].Topics)
					assert.Equal(t, rpcReceipt.Logs[i].Data, result.Logs[i].Data)
					assert.Equal(t, rpcReceipt.Logs[i].BlockNumber, result.Logs[i].BlockNumber)
					assert.Equal(t, rpcReceipt.Logs[i].TxHash, result.Logs[i].TxHash)
					assert.Equal(t, rpcReceipt.Logs[i].TxIndex, result.Logs[i].TxIndex)
					assert.Equal(t, rpcReceipt.Logs[i].BlockHash, result.Logs[i].BlockHash)
					assert.Equal(t, rpcReceipt.Logs[i].Index, result.Logs[i].Index)
					assert.Equal(t, rpcReceipt.Logs[i].Removed, result.Logs[i].Removed)
				}
				assert.Equal(t, rpcReceipt.Status, result.Status)
				assert.Equal(t, rpcReceipt.TxHash, result.TxHash)
				assert.Equal(t, rpcReceipt.TxL2Hash, result.TxL2Hash)
				assert.Equal(t, rpcReceipt.TxIndex, result.TxIndex)
				assert.Equal(t, rpcReceipt.BlockHash, result.BlockHash)
				assert.Equal(t, rpcReceipt.BlockNumber, result.BlockNumber)
				assert.Equal(t, rpcReceipt.GasUsed, result.GasUsed)
				assert.Equal(t, rpcReceipt.FromAddr, result.FromAddr)
				assert.Equal(t, rpcReceipt.ToAddr, result.ToAddr)
				assert.Equal(t, rpcReceipt.ContractAddress, result.ContractAddress)
				assert.Equal(t, rpcReceipt.Type, result.Type)
				assert.Equal(t, rpcReceipt.EffectiveGasPrice, result.EffectiveGasPrice)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				rpcErr := res.Error.RPCError()
				assert.Equal(t, tc.ExpectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.ExpectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func ptrArgUint64FromUint(n uint) *types.ArgUint64 {
	tmp := types.ArgUint64(n)
	return &tmp
}

func ptrArgUint64FromUint64(n uint64) *types.ArgUint64 {
	tmp := types.ArgUint64(n)
	return &tmp
}

func signTx(tx *ethTypes.Transaction, chainID uint64) *ethTypes.Transaction {
	privateKey, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(chainID))
	signedTx, _ := auth.Signer(auth.From, tx)
	return signedTx
}

func TestGetExitRootsByGER(t *testing.T) {
	type testCase struct {
		Name           string
		GER            common.Hash
		ExpectedResult *types.ExitRoots
		ExpectedError  types.Error
		SetupMocks     func(*mockedServer, *mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "GER not found",
			GER:            common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), tc.GER, nil).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:           "get exit roots fails to load exit roots from state",
			GER:            common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), tc.GER, nil).
					Return(nil, fmt.Errorf("failed to load exit roots from state"))
			},
		},
		{
			Name: "get exit roots successfully",
			GER:  common.HexToHash("0x345"),
			ExpectedResult: &types.ExitRoots{
				BlockNumber:     100,
				Timestamp:       types.ArgUint64(time.Now().Unix()),
				MainnetExitRoot: common.HexToHash("0x1"),
				RollupExitRoot:  common.HexToHash("0x2"),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocksWrapper, tc *testCase) {
				er := &state.GlobalExitRoot{
					BlockNumber:     uint64(tc.ExpectedResult.BlockNumber),
					Timestamp:       time.Unix(int64(tc.ExpectedResult.Timestamp), 0),
					MainnetExitRoot: tc.ExpectedResult.MainnetExitRoot,
					RollupExitRoot:  tc.ExpectedResult.RollupExitRoot,
				}

				m.State.
					On("GetExitRootByGlobalExitRoot", context.Background(), tc.GER, nil).
					Return(er, nil)
			},
		},
	}
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	zkEVMClient := client.NewClient(s.ServerURL)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(s, m, &tc)

			exitRoots, err := zkEVMClient.ExitRootsByGER(context.Background(), tc.GER)
			require.NoError(t, err)

			if exitRoots != nil || tc.ExpectedResult != nil {
				assert.Equal(t, tc.ExpectedResult.BlockNumber.Hex(), exitRoots.BlockNumber.Hex())
				assert.Equal(t, tc.ExpectedResult.Timestamp.Hex(), exitRoots.Timestamp.Hex())
				assert.Equal(t, tc.ExpectedResult.MainnetExitRoot.String(), exitRoots.MainnetExitRoot.String())
				assert.Equal(t, tc.ExpectedResult.RollupExitRoot.String(), exitRoots.RollupExitRoot.String())
			}

			if err != nil || tc.ExpectedError != nil {
				rpcErr := err.(types.RPCError)
				assert.Equal(t, tc.ExpectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.ExpectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func TestGetLatestGlobalExitRoot(t *testing.T) {
	type testCase struct {
		Name           string
		ExpectedResult *common.Hash
		ExpectedError  types.Error
		SetupMocks     func(*mocksWrapper, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "failed to load GER from state",
			ExpectedResult: nil,
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "couldn't load the last global exit root"),
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetLatestBatchGlobalExitRoot", context.Background(), nil).
					Return(nil, fmt.Errorf("failed to load GER from state")).
					Once()
			},
		},
		{
			Name:           "Get latest GER successfully",
			ExpectedResult: state.Ptr(common.HexToHash("0x1")),
			ExpectedError:  nil,
			SetupMocks: func(m *mocksWrapper, tc *testCase) {
				m.State.
					On("GetLatestBatchGlobalExitRoot", context.Background(), nil).
					Return(common.HexToHash("0x1"), nil).
					Once()
			},
		},
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	zkEVMClient := client.NewClient(s.ServerURL)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			ger, err := zkEVMClient.GetLatestGlobalExitRoot(context.Background())

			if tc.ExpectedResult != nil {
				assert.Equal(t, tc.ExpectedResult.String(), ger.String())
			}

			if err != nil || tc.ExpectedError != nil {
				rpcErr := err.(types.RPCError)
				assert.Equal(t, tc.ExpectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.ExpectedError.Error(), rpcErr.Error())
			}
		})
	}
}
