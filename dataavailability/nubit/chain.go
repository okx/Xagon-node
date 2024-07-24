package nubit

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// MarshalBatchData packs the batch data into ABI-encoded byte array
func MarshalBatchData(batchData [][]byte) ([]byte, error) {
	byteArrayType, _ := abi.NewType("bytes[]", "", nil)
	args := abi.Arguments{
		{Type: byteArrayType, Name: "batchData"},
	}
	res, err := args.Pack(&batchData)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("cannot pack batchData:%w", err)
	}
	return res, nil
}

// UnmarshalBatchData unpacks the ABI-encoded byte array into batch data
func UnmarshalBatchData(encodedData []byte) ([][]byte, error) {
	byteArrayType, _ := abi.NewType("bytes[]", "", nil)
	args := abi.Arguments{
		{Type: byteArrayType, Name: "batchData"},
	}
	res, err := args.Unpack(encodedData)
	if err != nil {
		return nil, fmt.Errorf("cannot unpack batchData: %w", err)
	}
	batchData := make([][]byte, len(res))
	for i, v := range res {
		byteSlice, ok := v.([]byte)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a []byte", i)
		}
		batchData[i] = byteSlice
	}
	return batchData, nil
}
