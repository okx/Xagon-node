package sequencer

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	validIP = "10.23.100.1"
)

var (
	// Init ZKEVM resourceCostMax values
	rcMax = state.BatchConstraintsCfg{
		MaxCumulativeGasUsed: 10,
		MaxArithmetics:       10,
		MaxBinaries:          10,
		MaxKeccakHashes:      10,
		MaxMemAligns:         10,
		MaxPoseidonHashes:    10,
		MaxPoseidonPaddings:  10,
		MaxSteps:             10,
		MaxSHA256Hashes:      10,
		MaxBatchBytesSize:    10,
	}
)

type workerAddTxTestCase struct {
	name                 string
	from                 common.Address
	txHash               common.Hash
	nonce                uint64
	cost                 *big.Int
	reservedZKCounters   state.ZKCounters
	usedBytes            uint64
	gasPrice             *big.Int
	expectedTxSortedList []common.Hash
	ip                   string
	expectedErr          error
}

type workerAddrQueueInfo struct {
	from    common.Address
	nonce   *big.Int
	balance *big.Int
}

func processWorkerAddTxTestCases(ctx context.Context, t *testing.T, worker *Worker, testCases []workerAddTxTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tx := TxTracker{}

			tx.Hash = testCase.txHash
			tx.HashStr = testCase.txHash.String()
			tx.From = testCase.from
			tx.FromStr = testCase.from.String()
			tx.Nonce = testCase.nonce
			tx.Cost = testCase.cost
			tx.Bytes = testCase.usedBytes
			tx.GasPrice = testCase.gasPrice
			tx.updateZKCounters(testCase.reservedZKCounters, testCase.reservedZKCounters)
			if testCase.ip == "" {
				// A random valid IP Address
				tx.IP = validIP
			} else {
				tx.IP = testCase.ip
			}
			t.Logf("%s=%d", testCase.name, tx.GasPrice)

			_, err := worker.AddTxTracker(ctx, &tx)
			if err != nil && testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				return
			}

			el := worker.txSortedList
			if el.len() != len(testCase.expectedTxSortedList) {
				t.Fatalf("Error txSortedList.len(%d) != expectedTxSortedList.len(%d)", el.len(), len(testCase.expectedTxSortedList))
			}
			for i := 0; i < el.len(); i++ {
				if el.getByIndex(i).HashStr != testCase.expectedTxSortedList[i].String() {
					t.Fatalf("Error txSortedList(%d). Expected=%s, Actual=%s", i, testCase.expectedTxSortedList[i].String(), el.getByIndex(i).HashStr)
				}
			}
		})
	}
}

func TestWorkerAddTx(t *testing.T) {
	var nilErr error

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock, rcMax)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/gp:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1, gasPrice: new(big.Int).SetInt64(10),
			cost:      new(big.Int).SetInt64(5),
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/gp:4", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(4),
			cost:      new(big.Int).SetInt64(5),
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{1}, {2},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/gp:20", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(20),
			cost:      new(big.Int).SetInt64(5),
			usedBytes: 5,
			expectedTxSortedList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Adding from:0x03, tx:0x03/gp:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1, gasPrice: new(big.Int).SetInt64(25),
			cost:      new(big.Int).SetInt64(5),
			usedBytes: 2,
			expectedTxSortedList: []common.Hash{
				{3}, {2}, {1},
			},
		},
		{
			name: "Invalid IP address", from: common.Address{5}, txHash: common.Hash{5}, nonce: 1,
			usedBytes:   1,
			ip:          "invalid IP",
			expectedErr: pool.ErrInvalidIP,
		},
		{
			name: "Out Of Counters Err",
			from: common.Address{5}, txHash: common.Hash{5}, nonce: 1,
			cost: new(big.Int).SetInt64(5),
			// Here, we intentionally set the reserved counters such that they violate the constraints
			reservedZKCounters: state.ZKCounters{
				GasUsed:          worker.batchConstraints.MaxCumulativeGasUsed + 1,
				KeccakHashes:     worker.batchConstraints.MaxKeccakHashes + 1,
				PoseidonHashes:   worker.batchConstraints.MaxPoseidonHashes + 1,
				PoseidonPaddings: worker.batchConstraints.MaxPoseidonPaddings + 1,
				MemAligns:        worker.batchConstraints.MaxMemAligns + 1,
				Arithmetics:      worker.batchConstraints.MaxArithmetics + 1,
				Binaries:         worker.batchConstraints.MaxBinaries + 1,
				Steps:            worker.batchConstraints.MaxSteps + 1,
				Sha256Hashes_V2:  worker.batchConstraints.MaxSHA256Hashes + 1,
			},
			usedBytes:   1,
			expectedErr: pool.ErrOutOfCounters,
		},
		{
			name: "Adding from:0x04, tx:0x04/gp:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1, gasPrice: new(big.Int).SetInt64(100),
			cost:      new(big.Int).SetInt64(5),
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{4}, {3}, {2}, {1},
			},
		},
	}

	processWorkerAddTxTestCases(ctx, t, worker, addTxsTC)
}

func TestWorkerGetBestTx(t *testing.T) {
	var nilErr error

	rc := state.BatchResources{
		ZKCounters: state.ZKCounters{GasUsed: 10, KeccakHashes: 10, PoseidonHashes: 10, PoseidonPaddings: 10, MemAligns: 10, Arithmetics: 10, Binaries: 10, Steps: 10, Sha256Hashes_V2: 10},
		Bytes:      10,
	}

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock, rcMax)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/gp:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1, gasPrice: new(big.Int).SetInt64(10),
			cost:               new(big.Int).SetInt64(5),
			reservedZKCounters: state.ZKCounters{GasUsed: 1, KeccakHashes: 1, PoseidonHashes: 1, PoseidonPaddings: 1, MemAligns: 1, Arithmetics: 1, Binaries: 1, Steps: 1, Sha256Hashes_V2: 1},
			usedBytes:          1,
			expectedTxSortedList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/gp:12", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(12),
			cost:               new(big.Int).SetInt64(5),
			reservedZKCounters: state.ZKCounters{GasUsed: 5, KeccakHashes: 5, PoseidonHashes: 5, PoseidonPaddings: 5, MemAligns: 5, Arithmetics: 5, Binaries: 5, Steps: 5, Sha256Hashes_V2: 5},
			usedBytes:          5,
			expectedTxSortedList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Readding from:0x03, tx:0x03/gp:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1, gasPrice: new(big.Int).SetInt64(25),
			cost:               new(big.Int).SetInt64(5),
			reservedZKCounters: state.ZKCounters{GasUsed: 2, KeccakHashes: 2, PoseidonHashes: 2, PoseidonPaddings: 2, MemAligns: 2, Arithmetics: 2, Binaries: 2, Steps: 2, Sha256Hashes_V2: 2},
			usedBytes:          2,
			expectedTxSortedList: []common.Hash{
				{3}, {2}, {1},
			},
		},
		{
			name: "Adding from:0x04, tx:0x04/gp:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1, gasPrice: new(big.Int).SetInt64(100),
			cost:               new(big.Int).SetInt64(5),
			reservedZKCounters: state.ZKCounters{GasUsed: 4, KeccakHashes: 4, PoseidonHashes: 4, PoseidonPaddings: 4, MemAligns: 4, Arithmetics: 4, Binaries: 4, Steps: 4, Sha256Hashes_V2: 4},
			usedBytes:          4,
			expectedTxSortedList: []common.Hash{
				{4}, {3}, {2}, {1},
			},
		},
	}

	processWorkerAddTxTestCases(ctx, t, worker, addTxsTC)

	expectedGetBestTx := []common.Hash{{4}, {3}, {1}}
	ct := 0

	for {
		tx, _ := worker.GetBestFittingTx(rc, state.ZKCounters{})
		if tx != nil {
			if ct >= len(expectedGetBestTx) {
				t.Fatalf("Error getting more best tx than expected. Expected=%d, Actual=%d", len(expectedGetBestTx), ct+1)
			}
			if tx.HashStr != expectedGetBestTx[ct].String() {
				t.Fatalf("Error GetBestFittingTx(%d). Expected=%s, Actual=%s", ct, expectedGetBestTx[ct].String(), tx.HashStr)
			}
			overflow, _ := rc.Sub(state.BatchResources{ZKCounters: tx.ReservedZKCounters, Bytes: tx.Bytes})
			assert.Equal(t, false, overflow)

			touch := make(map[common.Address]*state.InfoReadWrite)
			var newNonce uint64 = tx.Nonce + 1
			touch[tx.From] = &state.InfoReadWrite{Address: tx.From, Nonce: &newNonce, Balance: new(big.Int).SetInt64(10)}
			worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, touch)
			ct++
		} else {
			if ct < len(expectedGetBestTx) {
				t.Fatalf("Error expecting more best tx. Expected=%d, Actual=%d", len(expectedGetBestTx), ct)
			}
			break
		}
	}
}

func initWorker(stateMock *StateMock, rcMax state.BatchConstraintsCfg) *Worker {
	worker := NewWorker(stateMock, rcMax, newTimeoutCond(&sync.Mutex{}))
	return worker
}
