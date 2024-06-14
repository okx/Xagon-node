package dataavailability

// DABackendType is the data availability protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
	// Celestia is the DA protocol backend for Celestia
	Celestia DABackendType = "Celestia"
)
