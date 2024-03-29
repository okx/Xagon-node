package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

// NetworkConfig is the configuration struct for the different environments
type NetworkConfig struct {
	// L1: Configuration related to L1
	L1Config etherman.L1Config `json:"l1Config"`
	// L1: Genesis of the rollup, first block number and root

	// L2: address of the `PolygonZkEVMBridge proxy` smart contract
	L2BridgeAddr common.Address

	Genesis state.Genesis
}

type network string

const mainnet network = "mainnet"
const testnet network = "testnet"
const cardona network = "cardona"
const custom network = "custom"

// GenesisFromJSON is the config file for network_custom
type GenesisFromJSON struct {
	// L1: root hash of the genesis block
	Root string `json:"root"`
	// L1: block number in which the rollup was created
	RollupCreationBlockNum uint64 `json:"rollupCreationBlockNumber"`
	// L1: block number in which the rollup manager was created
	RollupManagerCreationBlockNum uint64 `json:"rollupManagerCreationBlockNumber"`
	// L2:  List of states contracts used to populate merkle tree at initial state
	Genesis []genesisAccountFromJSON `json:"genesis"`
	// L1: configuration of the network
	L1Config etherman.L1Config
}

type genesisAccountFromJSON struct {
	// Address of the account
	Balance string `json:"balance"`
	// Nonce of the account
	Nonce string `json:"nonce"`
	// Address of the contract
	Address string `json:"address"`
	// Byte code of the contract
	Bytecode string `json:"bytecode"`
	// Initial storage of the contract
	Storage map[string]string `json:"storage"`
	// Name of the contract in L1 (e.g. "PolygonZkEVMDeployer", "PolygonZkEVMBridge",...)
	ContractName string `json:"contractName"`
}

func (cfg *Config) loadNetworkConfig(ctx *cli.Context) {
	var networkJSON string
	switch ctx.String(FlagNetwork) {
	case string(mainnet):
		networkJSON = MainnetNetworkConfigJSON
	case string(testnet):
		networkJSON = TestnetNetworkConfigJSON
	case string(cardona):
		networkJSON = CardonaNetworkConfigJSON
	case string(custom):
		var err error
		cfgPath := ctx.String(FlagCustomNetwork)
		networkJSON, err = LoadGenesisFileAsString(cfgPath)
		if err != nil {
			panic(err.Error())
		}
	default:
		log.Fatalf("unsupported --network value. Must be one of: [%s, %s, %s]", mainnet, testnet, cardona, custom)
	}
	config, err := LoadGenesisFromJSONString(networkJSON)
	if err != nil {
		panic(fmt.Errorf("failed to load genesis configuration from file. Error: %v", err))
	}
	cfg.NetworkConfig = config
}

// LoadGenesisFileAsString loads the genesis file as a string
func LoadGenesisFileAsString(cfgPath string) (string, error) {
	if cfgPath != "" {
		f, err := os.Open(cfgPath) //nolint:gosec
		if err != nil {
			return "", err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Error(err)
			}
		}()

		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	} else {
		return "", errors.New("custom netwrork file not provided. Please use the custom-network-file flag")
	}
}

// LoadGenesisFromJSONString loads the genesis file from JSON string
func LoadGenesisFromJSONString(jsonStr string) (NetworkConfig, error) {
	var cfg NetworkConfig

	var cfgJSON GenesisFromJSON
	if err := json.Unmarshal([]byte(jsonStr), &cfgJSON); err != nil {
		return NetworkConfig{}, err
	}

	if len(cfgJSON.Genesis) == 0 {
		return cfg, nil
	}

	cfg.L1Config = cfgJSON.L1Config
	cfg.Genesis = state.Genesis{
		RollupBlockNumber:        cfgJSON.RollupCreationBlockNum,
		RollupManagerBlockNumber: cfgJSON.RollupManagerCreationBlockNum,
		Root:                     common.HexToHash(cfgJSON.Root),
		Actions:                  []*state.GenesisAction{},
	}

	checkAndSetBridgeAddress(&cfg, cfgJSON.Genesis)

	for _, account := range cfgJSON.Genesis {
		if account.Balance != "" && account.Balance != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   account.Balance,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if account.Nonce != "" && account.Nonce != "0" {
			action := &state.GenesisAction{
				Address: account.Address,
				Type:    int(merkletree.LeafTypeNonce),
				Value:   account.Nonce,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if account.Bytecode != "" {
			action := &state.GenesisAction{
				Address:  account.Address,
				Type:     int(merkletree.LeafTypeCode),
				Bytecode: account.Bytecode,
			}
			cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
		}
		if len(account.Storage) > 0 {
			for storageKey, storageValue := range account.Storage {
				action := &state.GenesisAction{
					Address:         account.Address,
					Type:            int(merkletree.LeafTypeStorage),
					StoragePosition: storageKey,
					Value:           storageValue,
				}
				cfg.Genesis.Actions = append(cfg.Genesis.Actions, action)
			}
		}
	}

	return cfg, nil
}
