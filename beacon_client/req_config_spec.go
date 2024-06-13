package beaconclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/hex"
)

// ConfigSpec returns the current beacon chain configuration
// Curl example:
// curl -X 'GET' \
// 'http://localhost/eth/v1/config/spec' \
// -H 'accept: application/json'
const configSpecPath = "/eth/v1/config/spec"

// ConfigSpecNodeResponse represents the response of the config spec endpoint
type ConfigSpecNodeResponse struct {
	SecondsPerSlot      uint64
	SecondsPerEth1Block uint64
}

type configSpecNodeResponseInternal struct {
	Data struct {
		SecondsPerSlot      string `json:"SECONDS_PER_SLOT"`
		SecondsPerEth1Block string `json:"SECONDS_PER_ETH1_BLOCK"`
	}
}

func convertConfigSpecResponseInternal(data configSpecNodeResponseInternal) (ConfigSpecNodeResponse, error) {
	tmpSecondsPerSlot, err := strconv.ParseUint(data.Data.SecondsPerSlot, 0, hex.BitSize64)
	if err != nil {
		return ConfigSpecNodeResponse{}, fmt.Errorf("error parsing SecondsPerSlot: %v", err)
	}
	tmpSecondsPerEth1Block, err := strconv.ParseUint(data.Data.SecondsPerEth1Block, 0, hex.BitSize64)
	if err != nil {
		return ConfigSpecNodeResponse{}, fmt.Errorf("error parsing SecondsPerSlot: %v", err)
	}
	res := ConfigSpecNodeResponse{
		SecondsPerSlot:      tmpSecondsPerSlot,
		SecondsPerEth1Block: tmpSecondsPerEth1Block,
	}
	return res, nil
}

// ConfigSpec returns the current beacon chain configuration
func (c *BeaconAPIClient) ConfigSpec(ctx context.Context) (*ConfigSpecNodeResponse, error) {
	response, err := JSONRPCBeaconCall(ctx, c.urlBase, configSpecPath)
	if err != nil {
		return nil, err
	}

	internalStruct, err := unserializeGenericResponse[configSpecNodeResponseInternal](response)
	if err != nil {
		return nil, err
	}

	responseData, err := convertConfigSpecResponseInternal(internalStruct)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
