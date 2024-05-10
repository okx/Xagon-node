package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	GlobalExitRoot, StateRoot, LocalExitRoot common.Hash
	AccInputHash                             common.Hash
	LastL2BLockTimestamp                     int64
	BatchL2Data                              []byte
	IsSequenceTooBig                         bool
	BatchNumber                              uint64
	ForcedBatchTimestamp                     int64
	PrevBlockHash                            common.Hash
}

// IsEmpty checks is sequence struct is empty
func (s Sequence) IsEmpty() bool {
	return reflect.DeepEqual(s, Sequence{})
}
