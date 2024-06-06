package eigenda

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrConvertFromABIInterface = errors.New("conversion from abi interface error")

type BlobData struct {
	BlobHeader            BlobHeader            `abi:"blobHeader"`
	BlobVerificationProof BlobVerificationProof `abi:"blobVerificationProof"`
}

type BlobHeader struct {
	Commitment       Commitment        `abi:"commitment"`
	DataLength       uint32            `abi:"dataLength"`
	QuorumBlobParams []QuorumBlobParam `abi:"quorumBlobParams"`
}

type BlobVerificationProof struct {
	BatchId        uint32        `abi:"batchId"`
	BlobIndex      uint32        `abi:"blobIndex"`
	BatchMetadata  BatchMetadata `abi:"batchMetadata"`
	InclusionProof []byte        `abi:"inclusionProof"`
	QuorumIndices  []byte        `abi:"quorumIndices"`
}

type Commitment struct {
	X *big.Int `abi:"X"`
	Y *big.Int `abi:"Y"`
}

type QuorumBlobParam struct {
	QuorumNumber                    uint8  `abi:"quorumNumber"`
	AdversaryThresholdPercentage    uint8  `abi:"adversaryThresholdPercentage"`
	ConfirmationThresholdPercentage uint8  `abi:"confirmationThresholdPercentage"`
	ChunkLength                     uint32 `abi:"chunkLength"`
}

type BatchMetadata struct {
	// The header of the data store
	BatchHeader BatchHeader `abi:"batchHeader"`
	// The hash of the signatory record
	SignatoryRecordHash common.Hash `abi:"signatoryRecordHash"`
	// The block number at which the batch was confirmed
	ConfirmationBlockNumber uint32 `abi:"confirmationBlockNumber"`
}

type BatchHeader struct {
	BlobHeadersRoot common.Hash `abi:"blobHeadersRoot"`
	// Each byte is a different quorum number
	QuorumNumbers []byte `abi:"quorumNumbers"`
	// Every bytes is an amount less than 100 specifying the percentage of stake
	// that has signed in the corresponding quorum in `quorumNumbers`
	SignedStakeForQuorums []byte `abi:"signedStakeForQuorums"`
	ReferenceBlockNumber  uint32 `abi:"referenceBlockNumber"`
}

type ReducedBatchHeader struct {
	BlobHeadersRoot      common.Hash `abi:"blobHeadersRoot"`
	ReferenceBlockNumber uint32      `abi:"referenceBlockNumber"`
}

func GetBlobData(info *disperser_rpc.BlobInfo) (BlobData, error) {
	header := GetBlobHeader(info.GetBlobHeader())
	proof, err := GetBlobVerificationProof(info.GetBlobVerificationProof())
	if err != nil {
		return BlobData{}, err
	}
	return BlobData{
		BlobHeader:            header,
		BlobVerificationProof: proof,
	}, nil
}

func GetBlobHeader(header *disperser_rpc.BlobHeader) BlobHeader {
	quorums := []QuorumBlobParam{}
	for _, quorum := range header.GetBlobQuorumParams() {
		q := QuorumBlobParam{
			QuorumNumber:                    uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage:    uint8(quorum.GetAdversaryThresholdPercentage()),
			ConfirmationThresholdPercentage: uint8(quorum.GetConfirmationThresholdPercentage()),
			ChunkLength:                     quorum.GetChunkLength(),
		}
		quorums = append(quorums, q)
	}

	return BlobHeader{
		Commitment: Commitment{
			X: new(big.Int).SetBytes(header.GetCommitment().GetX()),
			Y: new(big.Int).SetBytes(header.GetCommitment().GetY()),
		},
		DataLength:       header.GetDataLength(),
		QuorumBlobParams: quorums,
	}
}

func GetBlobVerificationProof(proof *disperser_rpc.BlobVerificationProof) (BlobVerificationProof, error) {
	if len(proof.BatchMetadata.BatchHeader.BatchRoot) != 32 {
		return BlobVerificationProof{}, fmt.Errorf("BlobHeadersRoot not type bytes32")
	}

	if len(proof.BatchMetadata.SignatoryRecordHash) != 32 {
		return BlobVerificationProof{}, fmt.Errorf("SignatoryRecordHash not type bytes32")
	}

	return BlobVerificationProof{
		BatchId:   proof.BatchId,
		BlobIndex: proof.BlobIndex,
		BatchMetadata: BatchMetadata{
			BatchHeader: BatchHeader{
				BlobHeadersRoot:       common.BytesToHash(proof.BatchMetadata.BatchHeader.BatchRoot),
				QuorumNumbers:         proof.BatchMetadata.BatchHeader.QuorumNumbers,
				SignedStakeForQuorums: proof.BatchMetadata.BatchHeader.QuorumSignedPercentages,
				ReferenceBlockNumber:  proof.BatchMetadata.BatchHeader.ReferenceBlockNumber,
			},
			SignatoryRecordHash:     common.BytesToHash(proof.BatchMetadata.SignatoryRecordHash),
			ConfirmationBlockNumber: proof.BatchMetadata.ConfirmationBlockNumber,
		},
		InclusionProof: proof.InclusionProof,
		QuorumIndices:  proof.QuorumIndexes,
	}, nil
}

// Method to calculate the BatchHeaderHash.
// Ref: https://github.com/Layr-Labs/eigenda/blob/5aecf5c2b679e69d363824513ba227388edcad82/contracts/src/libraries/EigenDAHasher.sol#L50
func (data BatchMetadata) GetBatchHeaderHash() ([]byte, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(batchHeaderABI)))
	if err != nil {
		return nil, err
	}

	// Encode the data
	method, exist := parsedABI.Methods["ReducedBatchHeader"]
	if !exist {
		return nil, fmt.Errorf("abi error, BatchHeader method not found")
	}

	ReducedBatchHeader := ReducedBatchHeader{
		BlobHeadersRoot:      data.BatchHeader.BlobHeadersRoot,
		ReferenceBlockNumber: data.BatchHeader.ReferenceBlockNumber,
	}
	encoded, err := method.Inputs.Pack(ReducedBatchHeader)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(encoded), nil
}

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
		case "BlobHeader":
			blobData.BlobHeader, err = convertBlobHeader(value)
			if err != nil {
				return BlobData{}, ErrConvertFromABIInterface
			}
		case "BlobVerificationProof":
			blobData.BlobVerificationProof, err = convertBlobVerificationProof(value)
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
func convertBlobHeader(val reflect.Value) (BlobHeader, error) {
	blobHeader := BlobHeader{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		switch field.Name {
		case "Commitment":
			blobHeader.Commitment = Commitment{
				X: value.FieldByName("X").Interface().(*big.Int),
				Y: value.FieldByName("Y").Interface().(*big.Int),
			}
		case "DataLength":
			blobHeader.DataLength = uint32(value.Uint())
		case "QuorumBlobParams":
			params := make([]QuorumBlobParam, value.Len())
			for j := 0; j < value.Len(); j++ {
				param := value.Index(j)
				params[j] = QuorumBlobParam{
					QuorumNumber:                    uint8(param.FieldByName("QuorumNumber").Uint()),
					AdversaryThresholdPercentage:    uint8(param.FieldByName("AdversaryThresholdPercentage").Uint()),
					ConfirmationThresholdPercentage: uint8(param.FieldByName("ConfirmationThresholdPercentage").Uint()),
					ChunkLength:                     uint32(param.FieldByName("ChunkLength").Uint()),
				}
			}
			blobHeader.QuorumBlobParams = params
		default:
			return BlobHeader{}, ErrConvertFromABIInterface
		}
	}

	return blobHeader, nil
}

func convertBlobVerificationProof(val reflect.Value) (BlobVerificationProof, error) {
	proof := BlobVerificationProof{}
	var err error

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		switch field.Name {
		case "BatchId":
			proof.BatchId = uint32(value.Uint())
		case "BlobIndex":
			proof.BlobIndex = uint32(value.Uint())
		case "BatchMetadata":
			proof.BatchMetadata, err = convertBatchMetadata(value)
			if err != nil {
				return BlobVerificationProof{}, ErrConvertFromABIInterface
			}
		case "InclusionProof":
			proof.InclusionProof = value.Interface().([]byte)
		case "QuorumIndices":
			proof.QuorumIndices = value.Interface().([]byte)
		default:
			return BlobVerificationProof{}, ErrConvertFromABIInterface
		}
	}

	return proof, nil
}

func convertBatchMetadata(val reflect.Value) (BatchMetadata, error) {
	metadata := BatchMetadata{}
	var err error

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		switch field.Name {
		case "BatchHeader":
			metadata.BatchHeader, err = convertBatchHeader(value)
			if err != nil {
				return BatchMetadata{}, ErrConvertFromABIInterface
			}
		case "SignatoryRecordHash":
			metadata.SignatoryRecordHash = value.Interface().([32]byte)
		case "ConfirmationBlockNumber":
			metadata.ConfirmationBlockNumber = uint32(value.Uint())
		default:
			return BatchMetadata{}, ErrConvertFromABIInterface
		}
	}

	return metadata, nil
}

func convertBatchHeader(val reflect.Value) (BatchHeader, error) {
	header := BatchHeader{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		switch field.Name {
		case "BlobHeadersRoot":
			header.BlobHeadersRoot = value.Interface().([32]byte)
		case "QuorumNumbers":
			header.QuorumNumbers = value.Interface().([]uint8)
		case "SignedStakeForQuorums":
			header.SignedStakeForQuorums = value.Interface().([]uint8)
		case "ReferenceBlockNumber":
			header.ReferenceBlockNumber = uint32(value.Uint())
		default:
			return BatchHeader{}, ErrConvertFromABIInterface
		}
	}

	return header, nil
}

// -------- Debugging methods --------
func (blobData BlobData) DebugLogBlobData() {
	fmt.Println("Logging blob data...")

	fmt.Println("--- Blob header ---")
	fmt.Println("Blob header commitment x: ", blobData.BlobHeader.Commitment.X)
	fmt.Println("Blob header commitment y: ", blobData.BlobHeader.Commitment.Y)
	fmt.Println("Blob header data length: ", blobData.BlobHeader.DataLength)
	for idx, q := range blobData.BlobHeader.QuorumBlobParams {
		fmt.Printf("Blob header quorum %v quorum number: %v\n", idx, q.QuorumNumber)
		fmt.Printf("Blob header quorum %v quorum adversary threshold percentage: %v\n", idx, q.AdversaryThresholdPercentage)
		fmt.Printf("Blob header quorum %v quorum confirmation threshold percentage: %v\n", idx, q.ConfirmationThresholdPercentage)
		fmt.Printf("Blob header quorum %v quorum chunk length: %v\n", idx, q.ChunkLength)
	}

	fmt.Println("--- Blob verification proof ---")
	fmt.Println("Blob verification proof batch id: ", blobData.BlobVerificationProof.BatchId)
	fmt.Println("Blob verification proof blob idx: ", blobData.BlobVerificationProof.BlobIndex)

	fmt.Println("Blob verification proof batch metadata:")
	fmt.Println("Blob verification proof batch metadata batch header batch headers root: ", blobData.BlobVerificationProof.BatchMetadata.BatchHeader.BlobHeadersRoot.Bytes())
	fmt.Println("Blob verification proof batch metadata batch header quorum numbers: ", blobData.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	fmt.Println("Blob verification proof batch metadata batch header signed stake for quorums: ", blobData.BlobVerificationProof.BatchMetadata.BatchHeader.SignedStakeForQuorums)
	fmt.Println("Blob verification proof batch metadata batch header reference block number: ", blobData.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber)
	fmt.Println("Blob verification proof batch metadata signature record hash: ", blobData.BlobVerificationProof.BatchMetadata.SignatoryRecordHash.Bytes())
	fmt.Println("Blob verification proof batch metadata confirmation block number: ", blobData.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber)

	fmt.Println("Blob verification proof inclusion proof: ", blobData.BlobVerificationProof.InclusionProof)
	fmt.Println("Blob verification proof quorum indices: ", blobData.BlobVerificationProof.QuorumIndices)
}

func DebugLogBlobInfoResponse(info *disperser_rpc.BlobInfo) {
	fmt.Println("Logging blob data...")

	fmt.Println("--- Blob header ---")
	fmt.Println("Blob header commitment x: ", info.BlobHeader.Commitment.X)
	fmt.Println("Blob header commitment y: ", info.BlobHeader.Commitment.Y)
	fmt.Println("Blob header data length: ", info.BlobHeader.DataLength)
	for idx, q := range info.BlobHeader.BlobQuorumParams {
		fmt.Printf("Blob header quorum %v quorum number: %v\n", idx, q.QuorumNumber)
		fmt.Printf("Blob header quorum %v quorum adversary threshold percentage: %v\n", idx, q.AdversaryThresholdPercentage)
		fmt.Printf("Blob header quorum %v quorum confirmation threshold percentage: %v\n", idx, q.ConfirmationThresholdPercentage)
		fmt.Printf("Blob header quorum %v quorum chunk length: %v\n", idx, q.ChunkLength)
	}

	fmt.Println("--- Blob verification proof ---")
	fmt.Println("Blob verification proof batch id: ", info.BlobVerificationProof.BatchId)
	fmt.Println("Blob verification proof blob idx: ", info.BlobVerificationProof.BlobIndex)

	fmt.Println("Blob verification proof batch metadata:")
	fmt.Println("Blob verification proof batch metadata batch header batch headers root: ", info.BlobVerificationProof.BatchMetadata.BatchHeader.BatchRoot)
	fmt.Println("Blob verification proof batch metadata batch header quorum numbers: ", info.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumNumbers)
	fmt.Println("Blob verification proof batch metadata batch header signed stake for quorums: ", info.BlobVerificationProof.BatchMetadata.BatchHeader.QuorumSignedPercentages)
	fmt.Println("Blob verification proof batch metadata batch header reference block number: ", info.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber)
	fmt.Println("Blob verification proof batch metadata signature record hash: ", info.BlobVerificationProof.BatchMetadata.SignatoryRecordHash)
	fmt.Println("Blob verification proof batch metadata confirmation block number: ", info.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber)

	fmt.Println("Blob verification proof inclusion proof: ", info.BlobVerificationProof.InclusionProof)
	fmt.Println("Blob verification proof quorum indices: ", info.BlobVerificationProof.QuorumIndexes)
}
