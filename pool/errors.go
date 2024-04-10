package pool

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// ErrInvalidChainID is returned when the transaction has a different chain id
	// than the chain id of the network
	ErrInvalidChainID = errors.New("invalid chain id")

	// ErrTxTypeNotSupported is returned if a transaction is not supported in the
	// current network configuration.
	ErrTxTypeNotSupported = types.ErrTxTypeNotSupported

	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")

	// ErrNegativeValue is a sanity error to ensure no one is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")

	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// ErrBlockedSender is returned if the transaction is sent by a blocked account.
	ErrBlockedSender = errors.New("blocked sender")

	// ErrBlockedReceiver is returned if the transaction is received by a blocked account.
	ErrBlockedReceiver = errors.New("blocked receiver")

	// ErrNoWhitelistedSender is returned if the transaction is sent by a no whitelisted account.
	ErrNoWhitelistedSender = errors.New("You are not allowed to send transactions on the X Layer as we are under the phase 1, X layer will be open to the public soon")

	// ErrGasLimit is returned if a transaction's requested gas limit exceeds the
	// maximum allowance of the current block.
	ErrGasLimit = errors.New("exceeds block gas limit")

	// ErrTxPoolAccountOverflow is returned if the account sending the transaction
	// has already reached the limit of transactions in the pool set by the config
	// AccountQueue and can't accept another remote transaction.
	ErrTxPoolAccountOverflow = errors.New("account has reached the tx limit in the txpool")

	// ErrTxPoolOverflow is returned if the transaction pool is full and can't accept
	// another remote transaction.
	ErrTxPoolOverflow = errors.New("txpool is full")

	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = errors.New("nonce too low")

	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// current + the configured AccountQueue.
	ErrNonceTooHigh = errors.New("nonce too high")

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient funds for gas * price + value")

	// ErrIntrinsicGas is returned if the transaction is specified to use less gas
	// than required to start the invocation.
	ErrIntrinsicGas = errors.New("intrinsic gas too low")

	// ErrGasUintOverflow is returned when calculating gas usage.
	ErrGasUintOverflow = errors.New("gas uint64 overflow")

	// ErrGasPrice is returned if the transaction has specified lower gas price than the minimum allowed.
	ErrGasPrice = errors.New("gas price too low")

	// ErrReceivedZeroL1GasPrice is returned if the L1 gas price is 0.
	ErrReceivedZeroL1GasPrice = errors.New("received L1 gas price 0")

	// ErrInvalidIP is returned if the IP address is invalid.
	ErrInvalidIP = errors.New("invalid IP address")

	// ErrOutOfCounters is returned if the pool is out of counters.
	ErrOutOfCounters = errors.New("out of counters")

	// ErrZeroL1GasPrice is returned if the L1 gas price is 0.
	ErrZeroL1GasPrice = errors.New("L1 gas price 0")
)
