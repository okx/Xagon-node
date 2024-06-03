package l1_check_block

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// L1BlockPoint is an enum that represents the point of the L1 block
type L1BlockPoint int

const (
	// FinalizedBlockNumber is the finalized block number
	FinalizedBlockNumber L1BlockPoint = 3
	// SafeBlockNumber is the safe block number
	SafeBlockNumber L1BlockPoint = 2
	// PendingBlockNumber is the pending block number
	PendingBlockNumber L1BlockPoint = 1
	// LastBlockNumber is the last block number
	LastBlockNumber L1BlockPoint = 0
)

// ToString converts a L1BlockPoint to a string
func (v L1BlockPoint) ToString() string {
	switch v {
	case FinalizedBlockNumber:
		return "finalized"
	case SafeBlockNumber:
		return "safe"
	case PendingBlockNumber:
		return "pending"
	case LastBlockNumber:
		return "latest"
	}
	return "Unknown"
}

// StringToL1BlockPoint converts a string to a L1BlockPoint
func StringToL1BlockPoint(s string) L1BlockPoint {
	switch s {
	case "finalized":
		return FinalizedBlockNumber
	case "safe":
		return SafeBlockNumber
	case "pending":
		return PendingBlockNumber
	case "latest":
		return LastBlockNumber
	default:
		return FinalizedBlockNumber
	}
}

// ToGethRequest converts a L1BlockPoint to a big.Int used for request to GETH
func (v L1BlockPoint) ToGethRequest() *big.Int {
	switch v {
	case FinalizedBlockNumber:
		return big.NewInt(int64(rpc.FinalizedBlockNumber))
	case PendingBlockNumber:
		return big.NewInt(int64(rpc.PendingBlockNumber))
	case SafeBlockNumber:
		return big.NewInt(int64(rpc.SafeBlockNumber))
	case LastBlockNumber:
		return nil
	}
	return big.NewInt(int64(v))
}

// SafeL1BlockNumberFetch is a struct that implements a safe L1 block number fetch
type SafeL1BlockNumberFetch struct {
	// SafeBlockPoint is the block number that is reference to l1 Block
	SafeBlockPoint L1BlockPoint
	// Offset is a vaule add to the L1 block
	Offset int
}

// NewSafeL1BlockNumberFetch creates a new SafeL1BlockNumberFetch
func NewSafeL1BlockNumberFetch(safeBlockPoint L1BlockPoint, offset int) *SafeL1BlockNumberFetch {
	return &SafeL1BlockNumberFetch{
		SafeBlockPoint: safeBlockPoint,
		Offset:         offset,
	}
}

// Description returns a string representation of SafeL1BlockNumberFetch
func (p *SafeL1BlockNumberFetch) Description() string {
	return fmt.Sprintf("%s/%d", p.SafeBlockPoint.ToString(), p.Offset)
}

// GetSafeBlockNumber gets the safe block number from L1
func (p *SafeL1BlockNumberFetch) GetSafeBlockNumber(ctx context.Context, requester L1Requester) (uint64, error) {
	l1SafePointBlock, err := requester.HeaderByNumber(ctx, p.SafeBlockPoint.ToGethRequest())
	if err != nil {
		log.Errorf("%s: Error getting L1 block %d. err: %s", logPrefix, p.String(), err.Error())
		return uint64(0), err
	}
	result := l1SafePointBlock.Number.Uint64()
	if p.Offset < 0 {
		if result < uint64(-p.Offset) {
			result = 0
		} else {
			result += uint64(p.Offset)
		}
	} else {
		result = l1SafePointBlock.Number.Uint64() + uint64(p.Offset)
	}
	if p.SafeBlockPoint == LastBlockNumber {
		result = min(result, l1SafePointBlock.Number.Uint64())
	}

	return result, nil
}

// String returns a string representation of SafeL1BlockNumberFetch
func (p *SafeL1BlockNumberFetch) String() string {
	return fmt.Sprintf("SafeBlockPoint: %s, Offset: %d", p.SafeBlockPoint.ToString(), p.Offset)
}
