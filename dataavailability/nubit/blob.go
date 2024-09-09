package nubit

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ErrConvertFromABIInterface is used when there is a decoding error
var ErrConvertFromABIInterface = errors.New("conversion from abi interface error")

// BlobData is the NubitDA blob data
type BlobData struct {
	BlobID    []byte `abi:"blobID"`
	Signature []byte `abi:"signature"`
}

// TryEncodeToDataAvailabilityMessage is a fallible encoding method to encode
// Nubit blob data into data availability message represented as byte array.
func TryEncodeToDataAvailabilityMessage(blobData BlobData) ([]byte, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(blobDataABI)))
	if err != nil {
		return nil, err
	}

	// Encode the data
	method, exist := parsedABI.Methods["BlobData"]
	if !exist {
		return nil, fmt.Errorf("abi error, BlobData method not found")
	}

	encoded, err := method.Inputs.Pack(blobData)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

// TryDecodeFromDataAvailabilityMessage is a fallible decoding method to
// decode data availability message into Nubit blob data.
func TryDecodeFromDataAvailabilityMessage(msg []byte) (BlobData, error) {
	// Parse the ABI
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(blobDataABI)))
	if err != nil {
		return BlobData{}, err
	}

	// Decode the data
	method, exist := parsedABI.Methods["BlobData"]
	if !exist {
		return BlobData{}, fmt.Errorf("abi error, BlobData method not found")
	}

	unpackedMap := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(unpackedMap, msg)
	if err != nil {
		return BlobData{}, err
	}
	unpacked, ok := unpackedMap["blobData"]
	if !ok {
		return BlobData{}, fmt.Errorf("abi error, failed to unpack to BlobData")
	}

	val := reflect.ValueOf(unpacked)
	typ := reflect.TypeOf(unpacked)

	blobData := BlobData{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		switch field.Name {
		case "BlobID":
			blobData.BlobID, err = convertBlobID(value)
			if err != nil {
				return BlobData{}, ErrConvertFromABIInterface
			}
		case "Signature":
			blobData.Signature, err = convertSignature(value)
			if err != nil {
				return BlobData{}, ErrConvertFromABIInterface
			}
		default:
			return BlobData{}, ErrConvertFromABIInterface
		}
	}

	return blobData, nil
}

// -------- Helper fallible conversion methods --------
func convertBlobID(val reflect.Value) ([]byte, error) {
	return val.Interface().([]byte), nil
}

func convertSignature(val reflect.Value) ([]byte, error) {
	return val.Interface().([]byte), nil
}
