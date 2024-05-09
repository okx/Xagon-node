package beaconclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
)

// /eth/v1/beacon/genesis
const beaconGenesisPath = "/eth/v1/beacon/genesis"

// BeaconGenesisResponse represents the response of the beacon genesis endpoint
type BeaconGenesisResponse struct {
	GenesisTime           uint64
	GenesisValidatorsRoot common.Address
	GenesisForkVersion    string
}

type beaconGenesisResponseInternal struct {
	Data struct {
		GenesisTime           string `json:"genesis_time"`
		GenesisValidatorsRoot string `json:"genesis_validators_root"`
		GenesisForkVersion    string `json:"genesis_fork_version"`
	} `json:"data"`
}

func convertBeaconGenesisResponseInternal(data beaconGenesisResponseInternal) (BeaconGenesisResponse, error) {
	genesisTime, err := strconv.ParseUint(data.Data.GenesisTime, 0, hex.BitSize64)
	if err != nil {
		return BeaconGenesisResponse{}, fmt.Errorf("error parsing genesisTime: %v", err)
	}
	res := BeaconGenesisResponse{
		GenesisTime:           genesisTime,
		GenesisValidatorsRoot: common.HexToAddress(data.Data.GenesisValidatorsRoot),
		GenesisForkVersion:    data.Data.GenesisForkVersion,
	}
	return res, nil
}

// BeaconGenesis request the current beacon chain genesis
func (c *BeaconAPIClient) BeaconGenesis(ctx context.Context) (*BeaconGenesisResponse, error) {
	response, err := JSONRPCBeaconCall(ctx, c.urlBase, beaconGenesisPath)
	if err != nil {
		return nil, err
	}

	internalStruct, err := unserializeGenericResponse[beaconGenesisResponseInternal](response)
	if err != nil {
		return nil, err
	}

	responseData, err := convertBeaconGenesisResponseInternal(internalStruct)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
