package incaberry

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "missmatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

type zkEVMClientInterface interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}

func isZeroByteArray(bytesArray [32]byte) bool {
	var zero = [32]byte{}
	return bytes.Equal(bytesArray[:], zero[:])
}

func (s *ProcessorL1SequenceBatches) getDataFromTrustedSequencer(ctx context.Context, batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	b, err := s.zkEVMClient.BatchByNumber(ctx, big.NewInt(int64(batchNum)))
	if err != nil {
		return nil, fmt.Errorf("failed to get batch num %d from trusted sequencer: %w", batchNum, err)
	}
	actualTransactionsHash := crypto.Keccak256Hash(b.BatchL2Data)
	if expectedTransactionsHash != actualTransactionsHash {
		return nil, fmt.Errorf(
			unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash,
		)
	}
	return b.BatchL2Data, nil
}
