package config

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

const l2BridgeSCName = "PolygonZkEVMBridge proxy"

func checkAndSetBridgeAddress(cfg *NetworkConfig, genesis []genesisAccountFromJSON) {
	for _, account := range genesis {
		if account.ContractName == l2BridgeSCName {
			cfg.L2BridgeAddr = common.HexToAddress(account.Address)
			log.Infof("Get L2 bridge address from genesis: %s", cfg.L2BridgeAddr.String())
			break
		}
	}
}
