package beaconclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/hex"
)

const beaconBlobSidecarsPath = "/eth/v1/beacon/blob_sidecars/"

// BeaconBlobSidecarsResponse represents the response of the beacon blob sidecars endpoint
type BeaconBlobSidecarsResponse struct {
	Sidecars map[uint64]BeaconBlobSidecarResponse
}

// BeaconBlobSidecarResponse represents the response of the config spec endpoint
type BeaconBlobSidecarResponse struct {
	Index         uint64
	KzgCommitment string
	Blob          []byte
}

type beaconBlobSidecarsResponseInternal struct {
	Data []struct {
		Index             string `json:"index"`
		Blob              string `json:"blob"`
		KzgCommitment     string `json:"kzg_commitment"`
		KzgProof          string `json:"kzg_proof"`
		SignedBlockHeader struct {
			Message struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"signed_block_header"`
		KzgCommitmentInclusionProof []string `json:"kzg_commitment_inclusion_proof"`
	} `json:"data"`
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func convertBeaconBlobSidecarsResponseInternal(data beaconBlobSidecarsResponseInternal) (*BeaconBlobSidecarsResponse, error) {
	response := BeaconBlobSidecarsResponse{
		Sidecars: make(map[uint64]BeaconBlobSidecarResponse),
	}
	for _, sidecar := range data.Data {
		index, err := strconv.ParseUint(sidecar.Index, 0, hex.BitSize64)
		if err != nil {
			return nil, fmt.Errorf("error parsing Index: %v", err)
		}
		//common.Hex2Bytes(sidecar.Blob)
		if has0xPrefix(sidecar.Blob) {
			sidecar.Blob = sidecar.Blob[2:]
		}
		blob, err := hex.DecodeHex(sidecar.Blob)
		if err != nil {
			return nil, fmt.Errorf("error decoding Blob: %v", err)
		}
		response.Sidecars[index] = BeaconBlobSidecarResponse{
			Index:         index,
			KzgCommitment: sidecar.KzgCommitment,
			Blob:          blob,
		}
	}
	return &response, nil
}

// BeaconBlobSidecars fetches the blob sidecars for a given blockID
func (c *BeaconAPIClient) BeaconBlobSidecars(ctx context.Context, blockID uint64) (*BeaconBlobSidecarsResponse, error) {
	response, err := JSONRPCBeaconCall(ctx, c.urlBase, beaconBlobSidecarsPath+fmt.Sprintf("%d", blockID))
	if err != nil {
		return nil, err
	}

	internalStruct, err := unserializeGenericResponse[beaconBlobSidecarsResponseInternal](response)
	if err != nil {
		return nil, err
	}

	responseData, err := convertBeaconBlobSidecarsResponseInternal(internalStruct)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
