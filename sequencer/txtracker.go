package sequencer

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxTracker is a struct that contains all the tx data needed to be managed by the worker
type TxTracker struct {
	Hash               common.Hash
	HashStr            string
	From               common.Address
	FromStr            string
	Nonce              uint64
	Gas                uint64 // To check if it fits into a batch
	GasPrice           *big.Int
	Cost               *big.Int // Cost = Amount + Benefit
	Bytes              uint64
	UsedZKCounters     state.ZKCounters
	ReservedZKCounters state.ZKCounters
	RawTx              []byte
	ReceivedAt         time.Time // To check if it has been in the txSortedList for too long
	IP                 string    // IP of the tx sender
	FailedReason       *string   // FailedReason is the reason why the tx failed, if it failed
	EffectiveGasPrice  *big.Int
	EGPPercentage      byte
	IsLastExecution    bool
	EGPLog             state.EffectiveGasPriceLog
	L1GasPrice         uint64
	L2GasPrice         uint64
}

// newTxTracker creates and inti a TxTracker
func newTxTracker(tx types.Transaction, usedZKCounters state.ZKCounters, reservedZKCounters state.ZKCounters, ip string) (*TxTracker, error) {
	addr, err := state.GetSender(tx)
	if err != nil {
		return nil, err
	}

	rawTx, err := state.EncodeTransactionWithoutEffectivePercentage(tx)
	if err != nil {
		return nil, err
	}

	txTracker := &TxTracker{
		Hash:               tx.Hash(),
		HashStr:            tx.Hash().String(),
		From:               addr,
		FromStr:            addr.String(),
		Nonce:              tx.Nonce(),
		Gas:                tx.Gas(),
		GasPrice:           tx.GasPrice(),
		Cost:               tx.Cost(),
		Bytes:              uint64(len(rawTx)) + state.EfficiencyPercentageByteLength,
		UsedZKCounters:     usedZKCounters,
		ReservedZKCounters: reservedZKCounters,
		RawTx:              rawTx,
		ReceivedAt:         time.Now(),
		IP:                 ip,
		EffectiveGasPrice:  new(big.Int).SetUint64(0),
		EGPLog: state.EffectiveGasPriceLog{
			ValueFinal:     new(big.Int).SetUint64(0),
			ValueFirst:     new(big.Int).SetUint64(0),
			ValueSecond:    new(big.Int).SetUint64(0),
			FinalDeviation: new(big.Int).SetUint64(0),
			MaxDeviation:   new(big.Int).SetUint64(0),
			GasPrice:       new(big.Int).SetUint64(0),
		},
	}

	return txTracker, nil
}

// updateZKCounters updates the used and reserved ZKCounters of the tx
func (tx *TxTracker) updateZKCounters(usedZKCounters state.ZKCounters, reservedZKCounters state.ZKCounters) {
	tx.UsedZKCounters = usedZKCounters
	tx.ReservedZKCounters = reservedZKCounters
}
