package nubit

import (
	"encoding/json"

	"github.com/cometbft/cometbft/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// KzgOpen represents a cryptographic opening in the context of a KZG commitment.
// This struct is used to verify the integrity and validity of a specific index within a cryptographic proof.
//
// Fields:
// - Index: The specific index in the data set to which this opening applies.
// - Value: The byte representation of the value at the given index.
// - Proof: The byte array containing the proof data that can be used to verify the value at the specified index.
type KzgOpen struct {
	Index int32  `json:"index,omitempty"`
	Value []byte `json:"value,omitempty"`
	Proof []byte `json:"proof,omitempty"`
}

// NamespaceRangeProof represents a proof of the presence or absence of a namespace ID in a Namespaced Merkle Tree (NMT).
// It provides the necessary information to verify the inclusion or exclusion of a namespace ID in the tree.
// The proof includes details about the position and cryptographic openings at various indices.
//
// Fields:
// - Start: The start index of the proof range in the NMT.
// - End: The end index of the proof range in the NMT.
// - PreIndex: The index immediately before the start of the namespace in the NMT.
// - PostIndex: The index immediately after the end of the namespace in the NMT.
// - OpenStart: The cryptographic opening for the start of the proof range.
// - OpenEnd: The cryptographic opening for the end of the proof range.
// - OpenPreIndex: The cryptographic opening for the index preceding the namespace.
// - OpenPostIndex: The cryptographic opening for the index following the namespace.
// - InclusionOrAbsence: A boolean indicating whether the namespace is included (true) or absent (false) in the tree.
type NamespaceRangeProof struct {
	Start              int32    `json:"start"`
	End                int32    `json:"end"`
	PreIndex           int32    `json:"pre_index"`
	PostIndex          int32    `json:"post_index"`
	OpenStart          *KzgOpen `json:"open_start"`
	OpenEnd            *KzgOpen `json:"open_end"`
	OpenPreIndex       *KzgOpen `json:"open_pre_index"`
	OpenPostIndex      *KzgOpen `json:"open_post_index"`
	InclusionOrAbsence bool     `json:"inclusion_or_absence"`
}

// RowProof is a Merkle proof that verifies the existence of a set of rows in a Merkle tree with a given data root.
// It contains the necessary information to prove that specific rows are part of the Merkle tree.
//
// Fields:
// - RowRoots: A slice of byte arrays representing the roots of the rows being proven.
// - Proofs: A list of Merkle proofs, each of which demonstrates that a row exists in the Merkle tree with the given data root.
// - StartRow: The index of the first row in the range being proven.
// - EndRow: The index of the last row in the range being proven.
type RowProof struct {
	RowRoots []tmbytes.HexBytes `json:"row_roots"`
	Proofs   []*merkle.Proof    `json:"proofs"`
	StartRow uint32             `json:"start_row"`
	EndRow   uint32             `json:"end_row"`
}

// ShareProof is a combination of an NMT proof and a Merkle proof that verifies the existence of a set of shares
// within specific rows of a Namespaced Merkle Tree (NMT), and the existence of those rows within a Merkle tree
// with a given data root.
//
// Fields:
// - Data: The raw shares being proven, represented as a slice of byte slices.
// - ShareProofs: A list of NMT proofs that demonstrate the presence of shares within the rows. Each proof corresponds to a row.
// - NamespaceID: The namespace ID of the shares being proven, used for verifying the proof. The proof will fail if the namespace ID does not match the shares' namespace.
// - RowProof: A Merkle proof that verifies the rows containing the shares exist in the Merkle tree with the specified data root.
// - NamespaceVersion: The version of the namespace, which may be used to differentiate between different namespace formats or versions.
type ShareProof struct {
	Data             [][]byte               `json:"data"`
	ShareProofs      []*NamespaceRangeProof `json:"share_proofs"`
	NamespaceID      []byte                 `json:"namespace_id"`
	RowProof         RowProof               `json:"row_proof"`
	NamespaceVersion uint32                 `json:"namespace_version"`
}

func UnmarshalInclusionProofs(proofBytes [][]byte) ([]*NamespaceRangeProof, error) {
	proofs := make([]*NamespaceRangeProof, 0)
	for _, b := range proofBytes {
		var proof NamespaceRangeProof
		err := json.Unmarshal(b, &proof)
		if err != nil {
			return proofs, err
		}
		proofs = append(proofs, &proof)
	}
	return proofs, nil
}
