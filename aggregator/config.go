package aggregator

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
)

// TokenAmountWithDecimals is a wrapper type that parses token amount with decimals to big int
type TokenAmountWithDecimals struct {
	*big.Int `validate:"required"`
}

// UnmarshalText unmarshal token amount from float string to big int
func (t *TokenAmountWithDecimals) UnmarshalText(data []byte) error {
	amount, ok := new(big.Float).SetString(string(data))
	if !ok {
		return fmt.Errorf("failed to unmarshal string to float")
	}
	coin := new(big.Float).SetInt(big.NewInt(encoding.TenToThePowerOf18))
	bigval := new(big.Float).Mul(amount, coin)
	result := new(big.Int)
	bigval.Int(result)
	t.Int = result

	return nil
}

// Config represents the configuration of the aggregator
type Config struct {
	// Host for the grpc server
	Host string `mapstructure:"Host"`
	// Port for the grpc server
	Port int `mapstructure:"Port"`

	// RetryTime is the time the aggregator main loop sleeps if there are no proofs to aggregate
	// or batches to generate proofs. It is also used in the isSynced loop
	RetryTime types.Duration `mapstructure:"RetryTime"`

	// VerifyProofInterval is the interval of time to verify/send an proof in L1
	VerifyProofInterval types.Duration `mapstructure:"VerifyProofInterval"`

	// ProofStatePollingInterval is the interval time to polling the prover about the generation state of a proof
	ProofStatePollingInterval types.Duration `mapstructure:"ProofStatePollingInterval"`

	// TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
	// possible values: base/acceptall
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`

	// TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
	// this parameter is used for the base tx profitability checker
	TxProfitabilityMinReward TokenAmountWithDecimals `mapstructure:"TxProfitabilityMinReward"`

	// IntervalAfterWhichBatchConsolidateAnyway this is interval for the main sequencer, that will check if there is no transactions
	IntervalAfterWhichBatchConsolidateAnyway types.Duration `mapstructure:"IntervalAfterWhichBatchConsolidateAnyway"`

	// ChainID is the L2 ChainID provided by the Network Config
	ChainID uint64

	// ForkID is the L2 ForkID provided by the Network Config
	ForkId uint64

	// SenderAddress defines which private key the eth tx manager needs to use
	// to sign the L1 txs
	SenderAddress string `mapstructure:"SenderAddress"`

	// CleanupLockedProofsInterval is the interval of time to clean up locked proofs.
	CleanupLockedProofsInterval types.Duration `mapstructure:"CleanupLockedProofsInterval"`

	// GeneratingProofCleanupThreshold represents the time interval after
	// which a proof in generating state is considered to be stuck and
	// allowed to be cleared.
	GeneratingProofCleanupThreshold string `mapstructure:"GeneratingProofCleanupThreshold"`

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

	// UpgradeEtrogBatchNumber is the number of the first batch after upgrading to etrog
	UpgradeEtrogBatchNumber uint64 `mapstructure:"UpgradeEtrogBatchNumber"`

	// BatchProofL1BlockConfirmations is number of L1 blocks to consider we can generate the proof for a virtual batch
	BatchProofL1BlockConfirmations uint64 `mapstructure:"BatchProofL1BlockConfirmations"`
}
