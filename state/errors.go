package state

import (
	"errors"
	"fmt"

	"github.com/umbracle/ethgo/abi"
)

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
	// ErrUnexpectedBatch indicates that the batch is unexpected
	ErrUnexpectedBatch = errors.New("unexpected batch")
	// ErrStateNotSynchronized indicates the state database may be empty
	ErrStateNotSynchronized = errors.New("state not synchronized")
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")
	// ErrNilDBTransaction indicates the db transaction has not been properly initialized
	ErrNilDBTransaction = errors.New("database transaction not properly initialized")
	// ErrAlreadyInitializedDBTransaction indicates the db transaction was already initialized
	ErrAlreadyInitializedDBTransaction = errors.New("database transaction already initialized")
	// ErrNotEnoughIntrinsicGas indicates the gas is not enough to cover the intrinsic gas cost
	ErrNotEnoughIntrinsicGas = fmt.Errorf("not enough gas supplied for intrinsic gas costs")
	// ErrParsingExecutorTrace indicates an error occurred while parsing the executor trace
	ErrParsingExecutorTrace = fmt.Errorf("error while parsing executor trace")
	// ErrInvalidBatchNumber indicates the provided batch number is not the latest in db
	ErrInvalidBatchNumber = errors.New("provided batch number is not latest")
	// ErrLastBatchShouldBeClosed indicates that last batch needs to be closed before adding a new one
	ErrLastBatchShouldBeClosed = errors.New("last batch needs to be closed before adding a new one")
	// ErrBatchAlreadyClosed indicates that batch is already closed
	ErrBatchAlreadyClosed = errors.New("batch is already closed")
	// ErrClosingBatchWithoutTxs indicates that the batch attempted to close does not have txs.
	ErrClosingBatchWithoutTxs = errors.New("can not close a batch without transactions")
	// ErrTimestampGE indicates that timestamp needs to be greater or equal
	ErrTimestampGE = errors.New("timestamp needs to be greater or equal")
	// ErrDBTxNil indicates that the method requires a dbTx that is not nil
	ErrDBTxNil = errors.New("the method requires a dbTx that is not nil")
	// ErrExistingTxGreaterThanProcessedTx indicates that we have more txs stored
	// in db than the txs we want to process.
	ErrExistingTxGreaterThanProcessedTx = errors.New("there are more transactions in the database than in the processed transaction set")
	// ErrOutOfOrderProcessedTx indicates the the processed transactions of an
	// ongoing batch are not in the same order as the transactions stored in the
	// database for the same batch.
	ErrOutOfOrderProcessedTx = errors.New("the processed transactions are not in the same order as the stored transactions")
	// ErrInsufficientFundsForTransfer is returned if the transaction sender doesn't
	// have enough funds for transfer(topmost call only).
	ErrInsufficientFundsForTransfer = errors.New("insufficient funds for transfer")
	// ErrExecutorNil indicates that the method requires an executor that is not nil
	ErrExecutorNil = errors.New("the method requires an executor that is not nil")
	// ErrStateTreeNil indicates that the method requires a state tree that is not nil
	ErrStateTreeNil = errors.New("the method requires a state tree that is not nil")
	// ErrUnsupportedDuration is returned if the provided unit for a time
	// interval is not supported by our conversion mechanism.
	ErrUnsupportedDuration = errors.New("unsupported time duration")
	// ErrInvalidData is the error when the raw txs is unexpected
	ErrInvalidData = errors.New("invalid data")
	// ErrInvalidBlockRange returned when the selected block range is invalid, generally
	// because the toBlock is bigger than the fromBlock
	ErrInvalidBlockRange = errors.New("invalid block range")
	// ErrMaxLogsCountLimitExceeded returned when the number of logs is bigger than the
	// configured limit
	ErrMaxLogsCountLimitExceeded = errors.New("query returned more than %v results")
	// ErrMaxLogsBlockRangeLimitExceeded returned when the range between block number range
	// to filter logs is bigger than the configured limit
	ErrMaxLogsBlockRangeLimitExceeded = errors.New("logs are limited to a %v block range")
	// ErrMaxNativeBlockHashBlockRangeLimitExceeded returned when the range between block number range
	// to filter native block hashes is bigger than the configured limit
	ErrMaxNativeBlockHashBlockRangeLimitExceeded = errors.New("native block hashes are limited to a %v block range")
)

// ConstructErrorFromRevert extracts the reverted reason from the provided returnValue
// and creates an instance of error that wraps the original error + the reverted reason
func ConstructErrorFromRevert(err error, returnValue []byte) error {
	revertErrMsg, unpackErr := abi.UnpackRevertError(returnValue)
	if unpackErr != nil {
		return err
	}

	return fmt.Errorf("%w: %s", err, revertErrMsg)
}

// BatchRemainingResourcesUnderflowError happens when the execution of a batch runs out of counters
type BatchRemainingResourcesUnderflowError struct {
	Message      string
	Code         int
	Err          error
	ResourceName string
}

// Error returns the error message
func (b BatchRemainingResourcesUnderflowError) Error() string {
	return constructErrorMsg(b.ResourceName)
}

func constructErrorMsg(resourceName string) string {
	return fmt.Sprintf("underflow of remaining resources for current batch. Resource %s", resourceName)
}
