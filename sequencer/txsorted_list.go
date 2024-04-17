package sequencer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// txSortedList represents a list of tx sorted by gasPrice
type txSortedList struct {
	list   map[string]*TxTracker
	sorted []*TxTracker
	mutex  sync.Mutex
}

// newTxSortedList creates and init an txSortedList
func newTxSortedList() *txSortedList {
	return &txSortedList{
		list:   make(map[string]*TxTracker),
		sorted: []*TxTracker{},
	}
}

// add adds a tx to the txSortedList
func (e *txSortedList) add(tx *TxTracker) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, found := e.list[tx.HashStr]; !found {
		e.list[tx.HashStr] = tx
		e.addSort(tx)
		return true
	}
	return false
}

// delete deletes the tx from the txSortedList
func (e *txSortedList) delete(tx *TxTracker) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if tx, found := e.list[tx.HashStr]; found {
		sLen := len(e.sorted)
		i := sort.Search(sLen, func(i int) bool {
			return e.isGreaterOrEqualThan(tx, e.list[e.sorted[i].HashStr])
		})

		// i is the index of the first tx that has equal (or lower) gasPrice than the tx. From here we need to go down in the list
		// looking for the sorted[i].HashStr equal to tx.HashStr to get the index of tx in the sorted slice.
		// We need to go down until we find the tx or we have a tx with different (lower) gasPrice or we reach the end of the list
		for {
			if i == sLen {
				log.Warnf("error deleting tx %s from txSortedList, we reach the end of the list", tx.HashStr)
				return false
			}

			if (e.sorted[i].GasPrice.Cmp(tx.GasPrice)) != 0 {
				// we have a tx with different (lower) GasPrice than the tx we are looking for, therefore we haven't found the tx
				log.Warnf("error deleting tx %s from txSortedList, not found in the list of txs with same gasPrice: %s", tx.HashStr)
				return false
			}

			if e.sorted[i].HashStr == tx.HashStr {
				break
			}

			i = i + 1
		}

		delete(e.list, tx.HashStr)

		copy(e.sorted[i:], e.sorted[i+1:])
		e.sorted[sLen-1] = nil
		e.sorted = e.sorted[:sLen-1]

		return true
	}
	return false
}

// getByIndex retrieves the tx at the i position in the sorted txSortedList
func (e *txSortedList) getByIndex(i int) *TxTracker {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	tx := e.sorted[i]

	return tx
}

// len returns the length of the txSortedList
func (e *txSortedList) len() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	l := len(e.sorted)

	return l
}

// print prints the contents of the txSortedList
func (e *txSortedList) Print() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	fmt.Println("Len: ", len(e.sorted))
	for _, txi := range e.sorted {
		fmt.Printf("Hash=%s, gasPrice=%d\n", txi.HashStr, txi.GasPrice)
	}
}

// addSort adds the tx to the txSortedList in a sorted way
func (e *txSortedList) addSort(tx *TxTracker) {
	i := sort.Search(len(e.sorted), func(i int) bool {
		return e.isGreaterThan(tx, e.list[e.sorted[i].HashStr])
	})

	e.sorted = append(e.sorted, nil)
	copy(e.sorted[i+1:], e.sorted[i:])
	e.sorted[i] = tx
	log.Debugf("added tx %s with  gasPrice %d to txSortedList at index %d from total %d", tx.HashStr, tx.GasPrice, i, len(e.sorted))
}

// isGreaterThan returns true if the tx1 has greater gasPrice than tx2
func (e *txSortedList) isGreaterThan(tx1 *TxTracker, tx2 *TxTracker) bool {
	cmp := tx1.GasPrice.Cmp(tx2.GasPrice)
	if cmp == 1 {
		return true
	} else {
		return false
	}
}

// isGreaterOrEqualThan returns true if the tx1 has greater or equal gasPrice than tx2
func (e *txSortedList) isGreaterOrEqualThan(tx1 *TxTracker, tx2 *TxTracker) bool {
	cmp := tx1.GasPrice.Cmp(tx2.GasPrice)
	if cmp >= 0 {
		return true
	} else {
		return false
	}
}

// GetSorted returns the sorted list of tx
func (e *txSortedList) GetSorted() []*TxTracker {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sorted
}
