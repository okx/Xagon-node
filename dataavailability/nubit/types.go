package nubit

import (
	"encoding/json"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollkit/go-da"
)

// BatchDAData contains the batch information on NubitDA
type BatchDAData struct {
	ID []da.ID `json:"id,omitempty"`
}

// Encode encodes the BatchDAData into ABI-encoded bytes
func (b *BatchDAData) Encode() ([]byte, error) {
	return json.Marshal(b)
}

// Decode decodes the ABI-encoded bytes into BatchDAData
func (b *BatchDAData) Decode(data []byte) error {
	return json.Unmarshal(data, &b)
}

// IsEmpty checks i fthe BatchDAData is empty
func (b *BatchDAData) IsEmpty() bool {
	return reflect.DeepEqual(b, BatchDAData{})
}

// DataCommitteeMember represents a member of the Data Committee
type DataCommitteeMember struct {
	Addr common.Address
	URL  string
}

// DataCommittee represents a specific committee
type DataCommittee struct {
	AddressesHash      common.Hash
	Members            []DataCommitteeMember
	RequiredSignatures uint64
}
