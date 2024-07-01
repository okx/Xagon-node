package jsonrpc

import (
	"crypto/md5"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/lru_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

func getCallKey(blockNumber *uint64, to *common.Address, input []byte) (string, string) {
	baseKey := fmt.Sprintf("%d-%s-%s", blockNumber, to.Hex(), md5.Sum(input))
	return baseKey + "ret", baseKey + "err"
}

func getCallResultFromLRU(blockNumber *uint64, to *common.Address, input []byte) (interface{}, types.Error, bool) {
	retKey, errKey := getCallKey(blockNumber, to, input)
	value, ok := lru_xlayer.GetLRU().Get(retKey)
	if !ok {
		return nil, nil, false
	}
	errValue, ok := lru_xlayer.GetLRU().Get(errKey)
	if !ok {
		return nil, nil, false
	}
	v, ok := errValue.(types.Error)
	if !ok {
		return nil, nil, false
	}

	return value, v, true
}

func setCallResultToLRU(blockNumber *uint64, to *common.Address, input []byte, value interface{}, errValue types.Error) {
	retKey, errKey := getCallKey(blockNumber, to, input)
	err := lru_xlayer.GetLRU().Set(retKey, value)
	if err != nil {
		log.Debugf("Failed to set value to LRU cache call ret: %v", err)
	}
	err = lru_xlayer.GetLRU().Set(errKey, errValue)
	if err != nil {
		log.Debugf("Failed to set value to LRU cache call err: %v", err)
	}
}
