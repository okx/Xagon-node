package ethtxmanager

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
)

// CustodialAssetsConfig is the config of the custodial assets
type CustodialAssetsConfig struct {
	// Enable is the flag to enable the custodial assets
	Enable bool `mapstructure:"Enable"`

	// URL is the url to sign the custodial assets
	URL string `mapstructure:"URL"`

	// Symbol is the symbol of the network, 2 prd, 2882 devnet
	Symbol int `mapstructure:"Symbol"`

	// SequencerAddr is the address of the sequencer
	SequencerAddr common.Address `mapstructure:"SequencerAddr"`

	// AggregatorAddr is the address of the aggregator
	AggregatorAddr common.Address `mapstructure:"AggregatorAddr"`

	// WaitResultTimeout is the timeout to wait for the result of the custodial assets
	WaitResultTimeout types.Duration `mapstructure:"WaitResultTimeout"`

	// OperateTypeSeq is the operate type of the custodial assets for the sequencer
	OperateTypeSeq int `mapstructure:"OperateTypeSeq"`

	// OperateTypeAgg is the operate type of the custodial assets for the aggregator
	OperateTypeAgg int `mapstructure:"OperateTypeAgg"`

	// ProjectSymbol is the project symbol of the custodial assets
	ProjectSymbol int `mapstructure:"ProjectSymbol"`

	// OperateSymbol is the operate symbol of the custodial assets
	OperateSymbol int `mapstructure:"OperateSymbol"`

	// SysFrom is the sys from of the custodial assets
	SysFrom int `mapstructure:"SysFrom"`

	// UserID is the user id of the custodial assets
	UserID int `mapstructure:"UserID"`

	// OperateAmount is the operate amount of the custodial assets
	OperateAmount int `mapstructure:"OperateAmount"`

	// RequestSignURI is the request sign uri of the custodial assets
	RequestSignURI string `mapstructure:"RequestSignURI"`

	// QuerySignURI is the query sign uri of the custodial assets
	QuerySignURI string `mapstructure:"QuerySignURI"`

	// AccessKey is the access key of the custodial assets
	AccessKey string `mapstructure:"AccessKey"`

	// SecretKey is the secret key of the custodial assets
	SecretKey string `mapstructure:"SecretKey"`
}
