package jsonrpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/lru_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

func getCallKey(blockNumber *uint64, sender common.Address, tx *ethtypes.Transaction) (string, string) {
	baseKey := fmt.Sprintf("%d-%s-%s", blockNumber, sender.String(), tx.Hash().String())
	return baseKey + "ret", baseKey + "err"
}

func getCallResultFromLRU(blockNumber *uint64, sender common.Address, tx *ethtypes.Transaction) (interface{}, types.Error, bool) {
	retKey, errKey := getCallKey(blockNumber, sender, tx)
	value, ok := lru_xlayer.GetLRU().Get(retKey)
	if !ok {
		return nil, nil, false
	}
	errValue, ok := lru_xlayer.GetLRU().Get(errKey)
	if !ok {
		return nil, nil, false
	}
	if errValue == nil {
		return value, nil, true
	}
	v, ok := errValue.(types.Error)
	if !ok {
		return nil, nil, false
	}

	return value, v, true
}

func setCallResultToLRU(blockNumber *uint64, sender common.Address, tx *ethtypes.Transaction, value interface{}, errValue types.Error) {
	retKey, errKey := getCallKey(blockNumber, sender, tx)
	err := lru_xlayer.GetLRU().Set(retKey, value)
	if err != nil {
		log.Debugf("Failed to set value to LRU cache call ret: %v", err)
		return
	}
	err = lru_xlayer.GetLRU().Set(errKey, errValue)
	if err != nil {
		log.Debugf("Failed to set value to LRU cache call err: %v", err)
	}
}
