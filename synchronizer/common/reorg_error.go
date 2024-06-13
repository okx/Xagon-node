package common

import "fmt"

// ReorgError is an error that is raised when a reorg is detected
type ReorgError struct {
	// BlockNumber is the block number that caused the reorg
	BlockNumber uint64
	Err         error
}

// NewReorgError creates a new ReorgError
func NewReorgError(blockNumber uint64, err error) *ReorgError {
	return &ReorgError{
		BlockNumber: blockNumber,
		Err:         err,
	}
}

func (e *ReorgError) Error() string {
	return fmt.Sprintf("%s blockNumber: %d", e.Err.Error(), e.BlockNumber)
}

// IsReorgError checks if an error is a ReorgError
func IsReorgError(err error) bool {
	_, ok := err.(*ReorgError)
	return ok
}

// GetReorgErrorBlockNumber returns the block number that caused the reorg
func GetReorgErrorBlockNumber(err error) uint64 {
	if reorgErr, ok := err.(*ReorgError); ok {
		return reorgErr.BlockNumber
	}
	return 0
}

// GetReorgError returns the error that caused the reorg
func GetReorgError(err error) error {
	if reorgErr, ok := err.(*ReorgError); ok {
		return reorgErr.Err
	}
	return nil
}
