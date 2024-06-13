package eip4844

import (
	"context"
	"fmt"

	beaconclient "github.com/0xPolygonHermez/zkevm-node/beacon_client"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// EthermanEIP4844 represents the EIP-4844 implementation
type EthermanEIP4844 struct {
	beaconClient   *beaconclient.BeaconAPIClient
	initialized    bool
	genesisTime    uint64
	secondsPerSlot uint64
}

// NewEthermanEIP4844 creates a new EthermanEIP4844
func NewEthermanEIP4844(beaconClient *beaconclient.BeaconAPIClient) *EthermanEIP4844 {
	return &EthermanEIP4844{
		beaconClient: beaconClient,
		initialized:  false,
	}
}

// IsInitialized returns if the EthermanEIP4844 is initialized
func (e *EthermanEIP4844) IsInitialized() bool {
	return e.initialized && e.genesisTime != 0 && e.secondsPerSlot != 0
}

// Initialize initializes the EthermanEIP4844
func (e *EthermanEIP4844) Initialize(ctx context.Context) error {
	// You can initialize multiples times and will fetch again the data

	configSpec, err := e.beaconClient.ConfigSpec(ctx)
	if err != nil {
		return fmt.Errorf("error fetching config spec: %v", err)
	}

	e.secondsPerSlot = configSpec.SecondsPerSlot

	genesis, err := e.beaconClient.BeaconGenesis(ctx)
	if err != nil {
		return fmt.Errorf("error fetching beacon genesis: %v", err)
	}
	e.genesisTime = genesis.GenesisTime
	if e.secondsPerSlot == 0 || e.genesisTime == 0 {
		return fmt.Errorf("genesisTime:%d or secondsPerSlot: %d is 0", e.genesisTime, e.secondsPerSlot)
	}
	e.initialized = true

	return nil
}

// GetBlobSidecar returns the blob sidecar for a given blockTime and kzgCommitment
func (e *EthermanEIP4844) GetBlobSidecar(ctx context.Context, blockTime uint64, kzgCommitment string) ([]byte, error) {
	slot, err := e.CalculateSlot(blockTime)
	if err != nil {
		errComposed := fmt.Errorf("error calculating Slot  blob sidecars: %v", err)
		log.Error(errComposed.Error())
		return nil, errComposed
	}
	sidecars, err := e.beaconClient.BeaconBlobSidecars(ctx, slot)
	if err != nil {
		errComposed := fmt.Errorf("error fetching beacon blob sidecars: %v", err)
		log.Error(errComposed.Error())
		return nil, errComposed
	}
	for _, sidecar := range sidecars.Sidecars {
		if sidecar.KzgCommitment == kzgCommitment {
			return sidecar.Blob, nil
		}
	}
	err = fmt.Errorf("sidecar not found")
	log.Error(err.Error())
	return nil, err
}

// CalculateSlot calculates the slot for a given blockTime
func (e *EthermanEIP4844) CalculateSlot(blockTime uint64) (uint64, error) {
	if !e.IsInitialized() {
		return 0, fmt.Errorf("EIP-4844 not initialized,please call Initialize(..) function first")
	}
	return (blockTime - e.genesisTime) / e.secondsPerSlot, nil
}
