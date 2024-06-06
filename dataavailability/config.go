package dataavailability

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// DABackendType is the data availability protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
	// DataAvailabilityEigenDA is the EigenDA protocol backend
	DataAvailabilityEigenDA DABackendType = "EigenDA"
)

// Config is the EigenDA network config
type Config struct {
	Hostname                   string         `mapstructure:"Hostname"`
	Port                       string         `mapstructure:"Port"`
	Timeout                    types.Duration `mapstructure:"Timeout"`
	UseSecureGrpcFlag          bool           `mapstructure:"UseSecureGrpcFlag"`
	RetrieveBlobStatusPeriod   types.Duration `mapstructure:"RetrieveBlobStatusPeriod"`
	BlobStatusConfirmedTimeout types.Duration `mapstructure:"BlobStatusConfirmedTimeout"`
}
