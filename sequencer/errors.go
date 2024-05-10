package sequencer

import "errors"

var (
	// ErrExpiredTransaction happens when the transaction is expired
	ErrExpiredTransaction = errors.New("transaction expired")
	// ErrEffectiveGasPriceReprocess happens when the effective gas price requires reexecution
	ErrEffectiveGasPriceReprocess = errors.New("effective gas price requires reprocessing the transaction")
	// ErrDuplicatedNonce is returned when adding a new tx to the worker and there is an existing tx
	// with the same nonce and higher gasPrice (in this case we keep the existing tx)
	ErrDuplicatedNonce = errors.New("duplicated nonce")
	// ErrReplacedTransaction is returned when an existing tx is replaced by a new tx with the same nonce and higher gasPrice
	ErrReplacedTransaction = errors.New("replaced transaction")
	// ErrGetBatchByNumber happens when we get an error trying to get a batch by number (GetBatchByNumber)
	ErrGetBatchByNumber = errors.New("get batch by number error")
	// ErrUpdateBatchAsChecked happens when we get an error trying to update a batch as checked (UpdateBatchAsChecked)
	ErrUpdateBatchAsChecked = errors.New("update batch as checked error")
	// ErrDecodeBatchL2Data happens when we get an error trying to decode BatchL2Data (DecodeTxs)
	ErrDecodeBatchL2Data = errors.New("decoding BatchL2Data error")
	// ErrProcessBatch happens when we get an error trying to process (executor) a batch
	ErrProcessBatch = errors.New("processing batch error")
	// ErrProcessBatchOOC happens when we get an OOC when processing (executor) a batch
	ErrProcessBatchOOC = errors.New("processing batch OOC")
	// ErrStateRootNoMatch happens when the SR returned for a full batch processing (sanity check) doesn't match
	// the SR calculated when filling a batch tx by tx
	ErrStateRootNoMatch = errors.New("state root no match")
	// ErrExecutorError happens when we got an executor error when processing a batch
	ErrExecutorError = errors.New("executor error")
	// ErrNoFittingTransaction happens when there is not a tx (from the txSortedList) that fits in the remaining batch resources
	ErrNoFittingTransaction = errors.New("no fit transaction")
	// ErrBatchResourceOverFlow happens when there is a tx that overlows remaining batch resources
	ErrBatchResourceOverFlow = errors.New("batch resource overflow")
	// ErrTransactionsListEmpty happens when txSortedList is empty
	ErrTransactionsListEmpty = errors.New("transactions list empty")
)
