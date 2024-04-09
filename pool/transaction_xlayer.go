package pool

import (
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

var l2BridgeAddr common.Address

// SetL2BridgeAddr sets the L2 bridge address
func SetL2BridgeAddr(value common.Address) {
	log.Infof("Set L2 bridge address: %s", value.String())
	l2BridgeAddr = value
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(freeClaimGasLimit uint64, bridgeClaimMethods []string) bool {
	if tx.To() == nil {
		return false
	}

	txGas := tx.Gas()
	if txGas > freeClaimGasLimit {
		return false
	}

	if *tx.To() != l2BridgeAddr {
		return false
	}

	methods := getClaimMethod(bridgeClaimMethods)
	for _, method := range methods {
		if strings.HasPrefix("0x"+common.Bytes2Hex(tx.Data()), method) {
			log.Infof("Transaction %s is a claim tx", tx.Hash().String())
			return true
		}
	}

	return false
}
