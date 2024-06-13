package sequencesender

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
)

// Config represents the configuration of a sequence sender
type Config struct {
	// WaitPeriodSendSequence is the time the sequencer waits until
	// trying to send a sequence to L1
	WaitPeriodSendSequence types.Duration `mapstructure:"WaitPeriodSendSequence"`
	// LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent
	LastBatchVirtualizationTimeMaxWaitPeriod types.Duration `mapstructure:"LastBatchVirtualizationTimeMaxWaitPeriod"`
	// L1BlockTimestampMargin is the time difference (margin) that must exists between last L1 block and last L2 block in the sequence before
	// to send the sequence to L1. If the difference is lower than this value then sequencesender will wait until the difference is equal or greater
	L1BlockTimestampMargin types.Duration `mapstructure:"L1BlockTimestampMargin"`
	// MaxTxSizeForL1 is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions than 128KB are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	MaxTxSizeForL1 uint64 `mapstructure:"MaxTxSizeForL1"`
	// SenderAddress defines which private key the eth tx manager needs to use
	// to sign the L1 txs
	SenderAddress common.Address
	// L2Coinbase defines which address is going to receive the fees
	L2Coinbase common.Address `mapstructure:"L2Coinbase"`
	// PrivateKey defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKey types.KeystoreFileConfig `mapstructure:"PrivateKey"`
	// Batch number where there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64
	// GasOffset is the amount of gas to be added to the gas estimation in order
	// to provide an amount that is higher than the estimated one. This is used
	// to avoid the TX getting reverted in case something has changed in the network
	// state after the estimation which can cause the TX to require more gas to be
	// executed.
	//
	// ex:
	// gas estimation: 1000
	// gas offset: 100
	// final gas: 1100
	GasOffset uint64 `mapstructure:"GasOffset"`
	// SequenceL1BlockConfirmations is number of blocks to consider a sequence sent to L1 as final
	SequenceL1BlockConfirmations uint64 `mapstructure:"SequenceL1BlockConfirmations"`
}
