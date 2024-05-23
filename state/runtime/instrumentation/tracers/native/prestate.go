// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package native

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run github.com/fjl/gencodec -type account -field-override accountMarshaling -out gen_account_json.go

func init() {
	tracers.DefaultDirectory.Register("prestateTracer", NewPrestateTracer, false)
}

type state = map[common.Address]*account

type account struct {
	Balance *big.Int                    `json:"balance,omitempty"`
	Code    []byte                      `json:"code,omitempty"`
	Nonce   uint64                      `json:"nonce,omitempty"`
	Storage map[common.Hash]common.Hash `json:"storage,omitempty"`
}

func (a *account) exists() bool {
	return a.Nonce > 0 || len(a.Code) > 0 || len(a.Storage) > 0 || (a.Balance != nil && a.Balance.Sign() != 0)
}

type accountMarshaling struct {
	Balance *hexutil.Big
	Code    hexutil.Bytes
}

type prestateTracer struct {
	noopTracer
	env       *fakevm.FakeEVM
	pre       state
	post      state
	create    bool
	to        common.Address
	gasLimit  uint64 // Amount of gas bought for the whole tx
	config    prestateTracerConfig
	interrupt uint32 // Atomic flag to signal execution interruption
	reason    error  // Textual reason for the interruption
	created   map[common.Address]bool
	deleted   map[common.Address]bool
}

type prestateTracerConfig struct {
	DiffMode bool `json:"diffMode"` // If true, this tracer will return state modifications
}

func NewPrestateTracer(ctx *tracers.Context, cfg json.RawMessage) (tracers.Tracer, error) {
	var config prestateTracerConfig
	if cfg != nil {
		if err := json.Unmarshal(cfg, &config); err != nil {
			return nil, err
		}
	}
	return &prestateTracer{
		pre:     state{},
		post:    state{},
		config:  config,
		created: make(map[common.Address]bool),
		deleted: make(map[common.Address]bool),
	}, nil
}

// CaptureStart implements the EVMLogger interface to initialize the tracing operation.
func (t *prestateTracer) CaptureStart(env *fakevm.FakeEVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	t.env = env
	t.create = create
	t.to = to

	t.lookupAccount(from)
	t.lookupAccount(to)
	t.lookupAccount(env.Context.Coinbase)

	// The recipient balance includes the value transferred.
	toBal := new(big.Int).Sub(t.pre[to].Balance, value)
	t.pre[to].Balance = toBal

	// The sender balance is after reducing: value and gasLimit.
	// We need to re-add them to get the pre-tx balance.
	fromBal := new(big.Int).Set(t.pre[from].Balance)
	gasPrice := env.TxContext.GasPrice
	consumedGas := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(t.gasLimit))
	fromBal.Add(fromBal, new(big.Int).Add(value, consumedGas))
	t.pre[from].Balance = fromBal
	t.pre[from].Nonce--

	if create && t.config.DiffMode {
		t.created[to] = true
	}
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (t *prestateTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	if t.config.DiffMode {
		return
	}

	if t.create {
		// Keep existing account prior to contract creation at that address
		if s := t.pre[t.to]; s != nil && !s.exists() {
			// Exclude newly created contract.
			delete(t.pre, t.to)
		}
	}
}

// CaptureState implements the EVMLogger interface to trace a single step of VM execution.
func (t *prestateTracer) CaptureState(pc uint64, op fakevm.OpCode, gas, cost uint64, scope *fakevm.ScopeContext, rData []byte, depth int, err error) {
	stack := scope.Stack
	stackData := stack.Data()
	stackLen := len(stackData)
	caller := scope.Contract.Address()
	switch {
	case stackLen >= 1 && (op == fakevm.SLOAD || op == fakevm.SSTORE):
		slot := common.Hash(stackData[stackLen-1].Bytes32())
		t.lookupStorage(caller, slot)
	case stackLen >= 1 && (op == fakevm.EXTCODECOPY || op == fakevm.EXTCODEHASH || op == fakevm.EXTCODESIZE || op == fakevm.BALANCE || op == fakevm.SELFDESTRUCT):
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.lookupAccount(addr)
		if op == fakevm.SELFDESTRUCT {
			t.deleted[caller] = true
		}
	case stackLen >= 5 && (op == fakevm.DELEGATECALL || op == fakevm.CALL || op == fakevm.STATICCALL || op == fakevm.CALLCODE):
		addr := common.Address(stackData[stackLen-2].Bytes20())
		t.lookupAccount(addr)
	case op == fakevm.CREATE:
		nonce := t.env.StateDB.GetNonce(caller)
		addr := crypto.CreateAddress(caller, nonce)
		t.lookupAccount(addr)
		t.created[addr] = true
	case stackLen >= 4 && op == fakevm.CREATE2:
		offset := stackData[stackLen-2]
		size := stackData[stackLen-3]
		init := scope.Memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		inithash := crypto.Keccak256(init)
		salt := stackData[stackLen-4]
		addr := crypto.CreateAddress2(caller, salt.Bytes32(), inithash)
		t.lookupAccount(addr)
		t.created[addr] = true
	}
}

func (t *prestateTracer) CaptureTxStart(gasLimit uint64) {
	t.gasLimit = gasLimit
}

func (t *prestateTracer) CaptureTxEnd(restGas uint64) {
	if !t.config.DiffMode {
		return
	}

	tStart := time.Now()
	var tAccBalance time.Duration
	var tAccNonce time.Duration
	var tAccCode time.Duration
	var tStorage time.Duration
	var accCount int
	for addr, state := range t.pre {
		// The deleted account's state is pruned from `post` but kept in `pre`
		if _, ok := t.deleted[addr]; ok {
			continue
		}
		accCount++
		modified := false
		postAccount := &account{Storage: make(map[common.Hash]common.Hash)}
		tAccBalanceStart := time.Now()
		newBalance := t.env.StateDB.GetBalance(addr)
		tAccNonceStart := time.Now()
		newNonce := t.env.StateDB.GetNonce(addr)
		tAccCodeStart := time.Now()
		newCode := t.env.StateDB.GetCode(addr)
		tAccEnd := time.Now()
		tAccBalance += tAccNonceStart.Sub(tAccBalanceStart)
		tAccNonce += tAccCodeStart.Sub(tAccNonceStart)
		tAccCode += tAccEnd.Sub(tAccCodeStart)

		if newBalance.Cmp(t.pre[addr].Balance) != 0 {
			modified = true
			postAccount.Balance = newBalance
		}
		if newNonce != t.pre[addr].Nonce {
			modified = true
			postAccount.Nonce = newNonce
		}
		if !bytes.Equal(newCode, t.pre[addr].Code) {
			modified = true
			postAccount.Code = newCode
		}

		tStorageStart := time.Now()
		for key, val := range state.Storage {
			// don't include the empty slot
			if val == (common.Hash{}) {
				delete(t.pre[addr].Storage, key)
			}

			newVal := t.env.StateDB.GetState(addr, key)
			if val == newVal {
				// Omit unchanged slots
				delete(t.pre[addr].Storage, key)
			} else {
				modified = true
				if newVal != (common.Hash{}) {
					postAccount.Storage[key] = newVal
				}
			}
		}
		tStorageEnd := time.Now()
		tStorage += tStorageEnd.Sub(tStorageStart)

		if modified {
			t.post[addr] = postAccount
		} else {
			// if state is not modified, then no need to include into the pre state
			delete(t.pre, addr)
		}
	}
	tEnd := time.Now()
	log.Infof("CaptureTxEnd took %v, tAccBalance %v, tAccNonce %v, tAccCode %v, tStorage %v counter %v",
		tEnd.Sub(tStart), tAccBalance, tAccNonce, tAccCode, tStorage, accCount)
	// the new created contracts' prestate were empty, so delete them
	for a := range t.created {
		// the created contract maybe exists in statedb before the creating tx
		if s := t.pre[a]; s != nil && !s.exists() {
			delete(t.pre, a)
		}
	}
}

// GetResult returns the json-encoded nested list of call traces, and any
// error arising from the encoding or forceful termination (via `Stop`).
func (t *prestateTracer) GetResult() (json.RawMessage, error) {
	var res []byte
	var err error
	if t.config.DiffMode {
		res, err = json.Marshal(struct {
			Post state `json:"post"`
			Pre  state `json:"pre"`
		}{t.post, t.pre})
	} else {
		res, err = json.Marshal(t.pre)
	}
	if err != nil {
		return nil, err
	}
	return json.RawMessage(res), t.reason
}

// Stop terminates execution of the tracer at the first opportune moment.
func (t *prestateTracer) Stop(err error) {
	t.reason = err
	atomic.StoreUint32(&t.interrupt, 1)
}

// lookupAccount fetches details of an account and adds it to the prestate
// if it doesn't exist there.
func (t *prestateTracer) lookupAccount(addr common.Address) {
	if _, ok := t.pre[addr]; ok {
		return
	}

	t.pre[addr] = &account{
		Balance: t.env.StateDB.GetBalance(addr),
		Nonce:   t.env.StateDB.GetNonce(addr),
		Code:    t.env.StateDB.GetCode(addr),
		Storage: make(map[common.Hash]common.Hash),
	}
}

// lookupStorage fetches the requested storage slot and adds
// it to the prestate of the given contract. It assumes `lookupAccount`
// has been performed on the contract before.
func (t *prestateTracer) lookupStorage(addr common.Address, key common.Hash) {
	if _, ok := t.pre[addr].Storage[key]; ok {
		return
	}
	t.pre[addr].Storage[key] = t.env.StateDB.GetState(addr, key)
}
