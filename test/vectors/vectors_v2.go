package vectors

import (
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// StateTransitionTestCaseV2 holds the metadata needed to run a state transition test
type StateTransitionTestCaseV2 struct {
	Description          string          `json:"Description"`
	Genesis              []GenesisEntity `json:"genesis"`
	ExpectedOldStateRoot string          `json:"expectedOldStateRoot"`
	ExpectedNewStateRoot string          `json:"expectedNewStateRoot"`
	ExpectedNewLeafs     []LeafV2        `json:"expectedNewLeafs"`
	Receipts             []TestReceipt   `json:"receipts"`
	GlobalExitRoot       string          `json:"globalExitRoot"`
	BatchL2Data          string          `json:"batchL2Data"`
}

// LeafV2 represents the state of a leaf in the merkle tree
type LeafV2 struct {
	Address         string            `json:"address"`
	Balance         argBigInt         `json:"balance"`
	Nonce           string            `json:"nonce"`
	Storage         map[string]string `json:"storage"`
	IsSmartContract bool              `json:"isSmartContract"`
	Bytecode        string            `json:"bytecode"`
	HashBytecode    string            `json:"hashBytecode"`
	BytecodeLength  int               `json:"bytecodeLength"`
}

// GenesisEntity represents the state of an account or smart contract when the network
// starts
type GenesisEntity struct {
	Address         string            `json:"address"`
	PvtKey          *string           `json:"pvtKey"`
	Balance         argBigInt         `json:"balance"`
	Nonce           string            `json:"nonce"`
	Storage         map[string]string `json:"storage"`
	IsSmartContract bool              `json:"isSmartContract"`
	Bytecode        *string           `json:"bytecode"`
}

func GenerateGenesisActions(genesis []GenesisEntity) []*state.GenesisAction {
	var genesisActions []*state.GenesisAction
	for _, genesisEntity := range genesis {

		if genesisEntity.Balance.String() != "0" {
			action := &state.GenesisAction{
				Address: genesisEntity.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   genesisEntity.Balance.String(),
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.Nonce != "" && genesisEntity.Nonce != "0" {
			action := &state.GenesisAction{
				Address: genesisEntity.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   genesisEntity.Nonce,
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.IsSmartContract && genesisEntity.Bytecode != nil && *genesisEntity.Bytecode != "0x" {
			action := &state.GenesisAction{
				Address:  genesisEntity.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: *genesisEntity.Bytecode,
			}
			genesisActions = append(genesisActions, action)
		}

		if genesisEntity.IsSmartContract && len(genesisEntity.Storage) > 0 {
			for storageKey, storageValue := range genesisEntity.Storage {
				action := &state.GenesisAction{
					Address:         genesisEntity.Address,
					Type:            int(merkletree.LeafTypeStorage),
					StoragePosition: storageKey,
					Value:           storageValue,
				}
				genesisActions = append(genesisActions, action)
			}
		}
	}

	return genesisActions
}
