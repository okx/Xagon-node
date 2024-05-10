package etherman

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/feijoapolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// SequenceBlobsOrder identifies a SequenceBlobs order
	SequenceBlobsOrder EventOrder = "SequenceBlobs"
)

var (
	// Events Feijoa Signatures
	// Events new ZkEvm/RollupBase
	// lastBlobSequenced is the count of blob sequenced after process this event
	//  if the first event have 1 blob -> lastBlobSequenced=1
	eventSequenceBlobsSignatureHash = crypto.Keccak256Hash([]byte("SequenceBlobs(uint64)"))
)

// EventFeijoaSequenceBlobsProcessor is the processor for event SequenceBlobs(uint64)
type EventFeijoaSequenceBlobsProcessor struct {
	contracts *FeijoaContracts
}

// NewEventFeijoaSequenceBlobsProcessor creates a new EventFeijoaSequenceBlobsProcessor
func NewEventFeijoaSequenceBlobsProcessor(contracts *FeijoaContracts) *EventFeijoaSequenceBlobsProcessor {
	return &EventFeijoaSequenceBlobsProcessor{
		contracts: contracts,
	}
}

// EventSignature returns the event signature supported
func (e *EventFeijoaSequenceBlobsProcessor) EventSignature() common.Hash {
	return eventSequenceBlobsSignatureHash
}

// AddEventDataToBlock adds the event data to the block and returns the Order
func (e *EventFeijoaSequenceBlobsProcessor) AddEventDataToBlock(ctx context.Context, vLog types.Log, block *Block, callData *CallData) (*Order, error) {
	//err := contract.UnpackLog(&event, "SequenceBlobs", vLog.Data)
	eventData, err := e.contracts.FeijoaZKEVM.ParseSequenceBlobs(vLog)
	if err != nil {
		return nil, err
	}
	for idx := range vLog.Topics {
		log.Debugf("vlog.Topics[%d]: %s ", idx, vLog.Topics[idx].Hex())
	}
	log.Debugf("LastBlobSequenced: %d", eventData.LastBlobSequenced)
	// decode Data
	inputData, err := e.parseCallData(callData)
	if err != nil {
		return nil, err
	}
	inputData.EventData = &SequenceBlobsEventData{
		LastBlobSequenced: eventData.LastBlobSequenced,
	}

	if inputData.thereIsAnyBlobType() {
		// TODO:Retrieve blobs
		return nil, fmt.Errorf("data-availability in blobs: not supported yet")
	}
	// Add the blobs to the block list
	block.SequenceBlobs = append(block.SequenceBlobs, *inputData)
	order := Order{
		Name: SequenceBatchesOrder,
		Pos:  len(block.SequenceBlobs) - 1,
	}

	return &order, nil
	// Extract Calldata
}

func (e *EventFeijoaSequenceBlobsProcessor) parseCallData(callData *CallData) (*SequenceBlobs, error) {
	//smcAbi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
	smcAbi, err := abi.JSON(strings.NewReader(feijoapolygonzkevm.FeijoapolygonzkevmABI))
	if err != nil {
		return nil, err
	}
	method, err := smcAbi.MethodById(callData.MethodID())
	if err != nil {
		return nil, err
	}
	// Unpack method inputs
	data, err := method.Inputs.Unpack(callData.InputData())
	if err != nil {
		return nil, err
	}
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	// Solidity: function sequenceBlobs((uint8,bytes)[] blobsRaw, address l2Coinbase, bytes32 finalAccInputHash) returns()
	var blobsRaw []feijoapolygonzkevm.PolygonRollupBaseFeijoaBlobData
	err = json.Unmarshal(bytedata, &blobsRaw)
	if err != nil {
		return nil, err
	}
	blobs := make([]SequenceBlob, 0)

	for i := range blobsRaw {
		//log.Debugf("BlobType: %d", blobs[i].BlobType)
		var blobBlobTypeParams *BlobBlobTypeParams
		var blobTypeParams *BlobCommonParams
		var txData []byte
		switch BlobType(blobsRaw[i].BlobType) {
		case TypeCallData:
			blobTypeParams, txData, err = parseBlobCallDataTypeParams(blobsRaw[i].BlobTypeParams)
			if err != nil {
				return nil, err
			}
		case TypeBlobTransaction:
			return nil, fmt.Errorf("blobType 'BlobTransaction' not supported yet")
		default:
			return nil, fmt.Errorf("blobType not supported")
		}
		blobs = append(blobs, SequenceBlob{
			Type:               BlobType(blobsRaw[i].BlobType),
			Params:             *blobTypeParams,
			Data:               txData,
			BlobBlobTypeParams: blobBlobTypeParams,
		})
	}
	l1CoinBase := (data[1]).(common.Address)
	finalAccInputHashRaw := (data[2]).([32]byte)
	finalAccInputHash := common.Hash(finalAccInputHashRaw)

	return &SequenceBlobs{
		Blobs:             blobs,
		L2Coinbase:        l1CoinBase,
		FinalAccInputHash: finalAccInputHash,
	}, nil
}

// returns data and the transtaction_data
func parseBlobCallDataTypeParams(data []byte) (*BlobCommonParams, []byte, error) {
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/feature/feijoa/contracts/v2/lib/PolygonRollupBaseFeijoa.sol
	// case: if (currentBlob.blobType == CALLDATA_BLOB_TYPE)
	//
	//		maxSequenceTimestamp uint64
	//		zkGasLimit           uint64
	//		l1InfoLeafIndex      uint32
	//	    transactions         []byte

	// Prepare blob params using ABI encoder
	uint64Ty, _ := abi.NewType("uint64", "", nil)
	uint32Ty, _ := abi.NewType("uint32", "", nil)
	bytesTy, _ := abi.NewType("bytes", "", nil)
	arguments := abi.Arguments{
		{Type: uint64Ty},
		{Type: uint64Ty},
		{Type: uint32Ty},
		{Type: bytesTy},
	}
	unpacked, err := arguments.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	result := &BlobCommonParams{}
	result.MaxSequenceTimestamp = unpacked[0].(uint64)
	result.ZkGasLimit = unpacked[1].(uint64)
	result.L1InfoLeafIndex = unpacked[2].(uint32)
	transactionData := unpacked[3].([]byte)
	return result, transactionData, nil
}
