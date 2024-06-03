package sequencer

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/ethereum/go-ethereum/common"
)

// addrQueue is a struct that stores the ready and notReady txs for a specific from address
type addrQueue struct {
	from              common.Address
	fromStr           string
	currentNonce      uint64
	currentBalance    *big.Int
	readyTx           *TxTracker
	notReadyTxs       map[uint64]*TxTracker
	forcedTxs         map[common.Hash]struct{}
	pendingTxsToStore map[common.Hash]struct{}
}

// newAddrQueue creates and init a addrQueue
func newAddrQueue(addr common.Address, nonce uint64, balance *big.Int) *addrQueue {
	return &addrQueue{
		from:              addr,
		fromStr:           addr.String(),
		currentNonce:      nonce,
		currentBalance:    balance,
		readyTx:           nil,
		notReadyTxs:       make(map[uint64]*TxTracker),
		forcedTxs:         make(map[common.Hash]struct{}),
		pendingTxsToStore: make(map[common.Hash]struct{}),
	}
}

// addTx adds a tx to the addrQueue and updates the ready a notReady Txs. Also if the new tx matches
// an existing tx with the same nonce but the new tx has better or equal gasPrice, we will return in the replacedTx
// the existing tx with lower gasPrice (the replacedTx will be later set as failed in the pool).
// If the existing tx has better gasPrice then we will drop the new tx (dropReason = ErrDuplicatedNonce)
func (a *addrQueue) addTx(tx *TxTracker) (newReadyTx, prevReadyTx, replacedTx *TxTracker, dropReason error) {
	var repTx *TxTracker

	if a.currentNonce == tx.Nonce { // Is a possible readyTx
		// We set the tx as readyTx if we do not have one assigned or if the gasPrice is better or equal than the current readyTx
		if a.readyTx == nil || ((a.readyTx != nil) && (tx.GasPrice.Cmp(a.readyTx.GasPrice) >= 0)) {
			oldReadyTx := a.readyTx
			if (oldReadyTx != nil) && (oldReadyTx.HashStr != tx.HashStr) {
				// if it is a different tx then we need to return the replaced tx to set as failed in the pool
				repTx = oldReadyTx
			}
			if a.currentBalance.Cmp(tx.Cost) >= 0 {
				a.readyTx = tx
				return tx, oldReadyTx, repTx, nil
			} else { // If there is not enough balance we set the new tx as notReadyTxs
				a.readyTx = nil
				a.notReadyTxs[tx.Nonce] = tx
				return nil, oldReadyTx, repTx, nil
			}
		} else { // We have an already readytx with the same nonce and better gas price, we discard the new tx
			return nil, nil, nil, ErrDuplicatedNonce
		}
	} else if a.currentNonce > tx.Nonce {
		return nil, nil, nil, runtime.ErrIntrinsicInvalidNonce
	}

	nrTx, found := a.notReadyTxs[tx.Nonce]
	if !found || ((found) && (tx.GasPrice.Cmp(nrTx.GasPrice) >= 0)) {
		a.notReadyTxs[tx.Nonce] = tx
		if (found) && (nrTx.HashStr != tx.HashStr) {
			// if it is a different tx then we need to return the replaced tx to set as failed in the pool
			repTx = nrTx
		}
		return nil, nil, repTx, nil
	} else {
		// We have an already notReadytx with the same nonce and better gas price, we discard the new tx
		return nil, nil, nil, ErrDuplicatedNonce
	}
}

// addForcedTx adds a forced tx to the list of forced txs
func (a *addrQueue) addForcedTx(txHash common.Hash) {
	a.forcedTxs[txHash] = struct{}{}
}

// addPendingTxToStore adds a tx to the list of pending txs to store in the DB (trusted state)
func (a *addrQueue) addPendingTxToStore(txHash common.Hash) {
	a.pendingTxsToStore[txHash] = struct{}{}
}

// ExpireTransactions removes the txs that have been in the queue for more than maxTime
func (a *addrQueue) ExpireTransactions(maxTime time.Duration) ([]*TxTracker, *TxTracker) {
	var (
		txs         []*TxTracker
		prevReadyTx *TxTracker
	)

	for _, txTracker := range a.notReadyTxs {
		if txTracker.ReceivedAt.Add(maxTime).Before(time.Now()) {
			txs = append(txs, txTracker)
			delete(a.notReadyTxs, txTracker.Nonce)
			log.Debugf("deleting notReadyTx %s from addrQueue %s", txTracker.HashStr, a.fromStr)
		}
	}

	if a.readyTx != nil && a.readyTx.ReceivedAt.Add(maxTime).Before(time.Now()) {
		prevReadyTx = a.readyTx
		txs = append(txs, a.readyTx)
		a.readyTx = nil
		log.Debugf("deleting readyTx %s from addrQueue %s", prevReadyTx.HashStr, a.fromStr)
	}

	return txs, prevReadyTx
}

// IsEmpty returns true if the addrQueue is empty
func (a *addrQueue) IsEmpty() bool {
	return a.readyTx == nil && len(a.notReadyTxs) == 0 && len(a.forcedTxs) == 0 && len(a.pendingTxsToStore) == 0
}

// deleteTx deletes the tx from the addrQueue
func (a *addrQueue) deleteTx(txHash common.Hash) (deletedTx *TxTracker, isReady bool) {
	txHashStr := txHash.String()

	if (a.readyTx != nil) && (a.readyTx.HashStr == txHashStr) {
		log.Infof("deleting readyTx %s from addrQueue %s", txHashStr, a.fromStr)
		prevReadyTx := a.readyTx
		a.readyTx = nil
		return prevReadyTx, true
	} else {
		var deletedTx *TxTracker
		for _, txTracker := range a.notReadyTxs {
			if txTracker.HashStr == txHashStr {
				deletedTx = txTracker
				log.Infof("deleting notReadyTx %s from addrQueue %s", txHashStr, a.fromStr)
				delete(a.notReadyTxs, txTracker.Nonce)
				break
			}
		}
		return deletedTx, false
	}
}

// deleteForcedTx deletes the tx from the addrQueue
func (a *addrQueue) deleteForcedTx(txHash common.Hash) {
	if _, found := a.forcedTxs[txHash]; found {
		delete(a.forcedTxs, txHash)
	} else {
		log.Warnf("tx %s not found in forcedTxs list", txHash.String())
	}
}

// deletePendingTxToStore delete a tx from the list of pending txs to store in the DB (trusted state)
func (a *addrQueue) deletePendingTxToStore(txHash common.Hash) {
	if _, found := a.pendingTxsToStore[txHash]; found {
		delete(a.pendingTxsToStore, txHash)
	} else {
		log.Warnf("tx %s not found in pendingTxsToStore list", txHash.String())
	}
}

func (a *addrQueue) getTransactions() []*TxTracker {
	// TODO: Add test for this function

	txsList := []*TxTracker{}

	if a.readyTx != nil {
		txsList = append(txsList, a.readyTx)
	}

	for _, tx := range a.notReadyTxs {
		txsList = append(txsList, tx)
	}

	return txsList
}

// updateCurrentNonceBalance updates the nonce and balance of the addrQueue and updates the ready and notReady txs
func (a *addrQueue) updateCurrentNonceBalance(nonce *uint64, balance *big.Int) (newReadyTx, prevReadyTx *TxTracker, toDelete []*TxTracker) {
	var oldReadyTx *TxTracker = nil
	txsToDelete := make([]*TxTracker, 0)

	if balance != nil {
		log.Debugf("updating balance for addrQueue %s from %s to %s", a.fromStr, a.currentBalance.String(), balance.String())
		a.currentBalance = balance
	}

	if nonce != nil {
		if a.currentNonce != *nonce {
			a.currentNonce = *nonce
			for _, txTracker := range a.notReadyTxs {
				if txTracker.Nonce < a.currentNonce {
					reason := runtime.ErrIntrinsicInvalidNonce.Error()
					txTracker.FailedReason = &reason
					txsToDelete = append(txsToDelete, txTracker)
				}
			}
			for _, txTracker := range txsToDelete {
				log.Infof("deleting notReadyTx with nonce %d from addrQueue %s, reason: %s", txTracker.Nonce, a.fromStr, *txTracker.FailedReason)
				delete(a.notReadyTxs, txTracker.Nonce)
			}
		}
	}

	if a.readyTx != nil {
		// If readyTX.nonce is not the currentNonce or currentBalance is less that the readyTx.Cost
		// set readyTx=nil. Later we will move the tx to notReadyTxs
		if (a.readyTx.Nonce != a.currentNonce) || (a.currentBalance.Cmp(a.readyTx.Cost) < 0) {
			oldReadyTx = a.readyTx
			a.readyTx = nil
		}
	}

	// We check if we have a new readyTx from the notReadyTxs (at this point, to optimize the code,
	// we are not including the oldReadyTx in notReadyTxs, as it can match again if the nonce has not changed)
	if a.readyTx == nil {
		nrTx, found := a.notReadyTxs[a.currentNonce]
		if found {
			if a.currentBalance.Cmp(nrTx.Cost) >= 0 {
				a.readyTx = nrTx
				log.Infof("set notReadyTx %s as readyTx for addrQueue %s", nrTx.HashStr, a.fromStr)
				delete(a.notReadyTxs, a.currentNonce)
			}
		}
	}

	// We add the oldReadyTx to notReadyTxs (if it has a valid nonce) at this point to avoid check it again in the previous if statement
	if oldReadyTx != nil && oldReadyTx.Nonce > a.currentNonce {
		log.Infof("set readyTx %s as notReadyTx from addrQueue %s", oldReadyTx.HashStr, a.fromStr)
		a.notReadyTxs[oldReadyTx.Nonce] = oldReadyTx
	} else if oldReadyTx != nil { // if oldReadyTx doesn't have a valid nonce then we add it to the txsToDelete
		reason := runtime.ErrIntrinsicInvalidNonce.Error()
		oldReadyTx.FailedReason = &reason
		txsToDelete = append(txsToDelete, oldReadyTx)
	}

	return a.readyTx, oldReadyTx, txsToDelete
}

// UpdateTxZKCounters updates the ZKCounters for the given tx (txHash)
func (a *addrQueue) UpdateTxZKCounters(txHash common.Hash, usedZKCounters state.ZKCounters, reservedZKCounters state.ZKCounters) {
	txHashStr := txHash.String()

	if (a.readyTx != nil) && (a.readyTx.HashStr == txHashStr) {
		log.Debugf("updating readyTx %s with new ZKCounters from addrQueue %s", txHashStr, a.fromStr)
		a.readyTx.updateZKCounters(usedZKCounters, reservedZKCounters)
	} else {
		for _, txTracker := range a.notReadyTxs {
			if txTracker.HashStr == txHashStr {
				log.Debugf("updating notReadyTx %s with new ZKCounters from addrQueue %s", txHashStr, a.fromStr)
				txTracker.updateZKCounters(usedZKCounters, reservedZKCounters)
				break
			}
		}
	}
}
