package dataavailability

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "mismatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

// DataAvailability implements an abstract data availability integration
type DataAvailability struct {
	isTrustedSequencer bool

	state       stateInterface
	zkEVMClient ZKEVMClientTrustedBatchesGetter
	backend     DABackender

	ctx context.Context
}

// New creates a DataAvailability instance
func New(
	isTrustedSequencer bool,
	backend DABackender,
	state stateInterface,
	zkEVMClient ZKEVMClientTrustedBatchesGetter,
) (*DataAvailability, error) {
	da := &DataAvailability{
		isTrustedSequencer: isTrustedSequencer,
		backend:            backend,
		state:              state,
		zkEVMClient:        zkEVMClient,
		ctx:                context.Background(),
	}
	err := da.backend.Init()
	return da, err
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (d *DataAvailability) PostSequence(ctx context.Context, sequences []types.Sequence) ([]byte, error) {
	batchesData := [][]byte{}
	for _, batch := range sequences {
		// Do not send to the DA backend data that will be stored to L1
		if batch.ForcedBatchTimestamp == 0 {
			batchesData = append(batchesData, batch.BatchL2Data)
		}
	}
	return d.backend.PostSequence(ctx, batchesData)
}

// GetBatchL2Data tries to return the data from a batch, in the following priorities
// 1. From local DB
// 2. From Sequencer
// 3. From DA backend
func (d *DataAvailability) GetBatchL2Data(batchNums []uint64, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	if len(batchNums) != len(batchHashes) {
		return nil, fmt.Errorf("invalid L2 batch data retrieval arguments, %d != %d", len(batchNums), len(batchHashes))
	}

	data, err := d.localData(batchNums, batchHashes)
	if err == nil {
		return data, nil
	}

	if !d.isTrustedSequencer {
		log.Infof("Try to get data from trusted sequencer for batches %v", batchNums)
		data, err = d.trustedSequencerData(batchNums, batchHashes)
		if err != nil {
			log.Warnf("trusted sequencer failed to return data for batches %v: %s", batchNums, err.Error())
		} else {
			return data, nil
		}
	}
	log.Infof("Try to get data from DA backend for batches %v", batchNums)
	return d.backend.GetSequence(d.ctx, batchHashes, dataAvailabilityMessage)
}

// localData retrieves batches from local database and returns an error unless all are found
func (d *DataAvailability) localData(numbers []uint64, hashes []common.Hash) ([][]byte, error) {
	data, err := d.state.GetBatchL2DataByNumbers(d.ctx, numbers, nil)
	if err != nil {
		return nil, err
	}
	var batches [][]byte
	for i := 0; i < len(numbers); i++ {
		batchNumber := numbers[i]
		expectedHash := hashes[i]
		batchData, ok := data[batchNumber]
		if !ok {
			return nil, fmt.Errorf("missing batch %v", batchNumber)
		}
		actualHash := crypto.Keccak256Hash(batchData)
		if actualHash != expectedHash {
			err = fmt.Errorf(unexpectedHashTemplate, batchNumber, expectedHash, actualHash)
			log.Warnf("wrong local data for hash: %s", err.Error())
			return nil, err
		} else {
			batches = append(batches, batchData)
		}
	}
	return batches, nil
}

// trustedSequencerData retrieved batch data from the trusted sequencer and returns an error unless all are found
func (d *DataAvailability) trustedSequencerData(batchNums []uint64, expectedHashes []common.Hash) ([][]byte, error) {
	if len(batchNums) != len(expectedHashes) {
		return nil, fmt.Errorf("invalid arguments, len of batch numbers does not equal length of expected hashes: %d != %d",
			len(batchNums), len(expectedHashes))
	}
	var nums []*big.Int
	for _, n := range batchNums {
		nums = append(nums, new(big.Int).SetUint64(n))
	}
	batchData, err := d.zkEVMClient.BatchesByNumbers(d.ctx, nums)
	if err != nil {
		return nil, err
	}
	if len(batchData) != len(batchNums) {
		return nil, fmt.Errorf("missing batch data, expected %d, got %d", len(batchNums), len(batchData))
	}
	var result [][]byte
	for i := 0; i < len(batchNums); i++ {
		number := batchNums[i]
		batch := batchData[i]
		expectedTransactionsHash := expectedHashes[i]
		actualTransactionsHash := crypto.Keccak256Hash(batch.BatchL2Data)
		if expectedTransactionsHash != actualTransactionsHash {
			return nil, fmt.Errorf(unexpectedHashTemplate, number, expectedTransactionsHash, actualTransactionsHash)
		}
		result = append(result, batch.BatchL2Data)
	}
	return result, nil
}
