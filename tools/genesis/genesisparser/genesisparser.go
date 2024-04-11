package genesisparser

import (
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// GenesisAccountTest struct
type GenesisAccountTest struct {
	Balance  string
	Nonce    string
	Address  string
	Bytecode string
	Storage  map[string]string
}

// GenesisTest2Actions change format from testvector to the used internaly
func GenesisTest2Actions(accounts []GenesisAccountTest) []*state.GenesisAction {
	leaves := make([]*state.GenesisAction, 0)

	for _, acc := range accounts {
		if len(acc.Balance) != 0 && acc.Balance != "0" {
			leaves = append(leaves, &state.GenesisAction{
				Address: acc.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   acc.Balance,
			})
		}
		if len(acc.Nonce) != 0 && acc.Nonce != "0" {
			leaves = append(leaves, &state.GenesisAction{
				Address: acc.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   acc.Nonce,
			})
		}
		if len(acc.Bytecode) != 0 {
			leaves = append(leaves, &state.GenesisAction{
				Address:  acc.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: acc.Bytecode,
			})
		}
		for key, value := range acc.Storage {
			leaves = append(leaves, &state.GenesisAction{
				Address:         acc.Address,
				Type:            int(merkletree.LeafTypeStorage),
				StoragePosition: key,
				Value:           value,
			})
		}
	}
	return leaves
}
