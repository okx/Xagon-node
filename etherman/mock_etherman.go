package etherman

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BuildMockSequenceBatchesTxData builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildMockSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence, l2Coinbase common.Address, committeeSignaturesAndAddrs []byte, txHashs [][32]byte) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceMockBatches(opts, sequences, l2Coinbase, committeeSignaturesAndAddrs, txHashs)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceMockBatches(opts bind.TransactOpts, sequences []ethmanTypes.Sequence, l2Coinbase common.Address, committeeSignaturesAndAddrs []byte, txHashs [][32]byte) (*types.Transaction, error) {
	var batches []polygonzkevm.PolygonZkEVMBatchData

	var tx *types.Transaction
	var err error
	if len(committeeSignaturesAndAddrs) > 0 {
		for index, seq := range sequences {
			batch := polygonzkevm.PolygonZkEVMBatchData{
				TransactionsHash:   txHashs[index],
				GlobalExitRoot:     seq.GlobalExitRoot,
				Timestamp:          uint64(seq.Timestamp),
				MinForcedTimestamp: uint64(seq.ForcedBatchTimestamp),
			}

			batches = append(batches, batch)
		}

		log.Infof("Sequence batches with validium.")
		tx, err = etherMan.ZkEVM.SequenceBatches(&opts, batches, l2Coinbase, committeeSignaturesAndAddrs)
	} else {
		for _, seq := range sequences {
			batch := polygonzkevm.PolygonZkEVMBatchData{
				Transactions:       seq.BatchL2Data,
				GlobalExitRoot:     seq.GlobalExitRoot,
				Timestamp:          uint64(seq.Timestamp),
				MinForcedTimestamp: uint64(seq.ForcedBatchTimestamp),
			}

			batches = append(batches, batch)
		}

		log.Infof("Sequence batches with rollup.")
		tx, err = etherMan.ZkEVM.SequenceBatches(&opts, batches, l2Coinbase, nil)
	}

	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		err = fmt.Errorf(
			"error sequencing batches: %w, committeeSignaturesAndAddrs %s",
			err, common.Bytes2Hex(committeeSignaturesAndAddrs),
		)
	}

	return tx, err
}
