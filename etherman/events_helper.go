package etherman

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockRetriever is the interface required from etherman main object
type BlockRetriever interface {
	RetrieveFullBlockForEvent(ctx context.Context, vLog types.Log) (*Block, error)
}

// GenericEventProcessor is the interface that a processor must implement
type GenericEventProcessor interface {
	// EventSignature returns the signature of the event supported
	// evaluate if make sens to support multiples signatures
	EventSignature() common.Hash
	AddEventDataToBlock(ctx context.Context, vLog types.Log, block *Block, callData *CallData) (*Order, error)
}

// CallDataExtractor is the interface required to extract the call data from a transaction
type CallDataExtractor interface {
	ExtractCallData(ctx context.Context, blockHash, txHash common.Hash, txIndex uint) (*CallData, error)
}

// EventManager is a struct that manages the L1 events
// The way of using this is create and add Processor
// A processor only need to code the specific part of adding specific data
// to the block
type EventManager struct {
	blockRetriever    BlockRetriever
	callDataExtractor CallDataExtractor

	processors []GenericEventProcessor
}

// NewEventManager creates a new EventManager
func NewEventManager(blockRetriever BlockRetriever, callDataExtractor CallDataExtractor) *EventManager {
	return &EventManager{
		blockRetriever:    blockRetriever,
		callDataExtractor: callDataExtractor,
		processors:        []GenericEventProcessor{},
	}
}

// AddProcessor adds a new processor to the EventManager
func (e *EventManager) AddProcessor(processor GenericEventProcessor) {
	e.processors = append(e.processors, processor)
}

// ProcessEvent processes an event
// this is the interface with etherman
// it returns true if this event belong to this processor
func (e *EventManager) ProcessEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) (bool, error) {
	for idx := range e.processors {
		processor := e.processors[idx]
		if len(vLog.Topics) > 0 && vLog.Topics[0] == processor.EventSignature() {
			return true, e.processGenericEvent(ctx, vLog, blocks, blocksOrder, processor)
		}
	}
	return false, nil
}

func (e *EventManager) processGenericEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order, processor GenericEventProcessor) error {
	callData, err := e.callDataExtractor.ExtractCallData(ctx, vLog.BlockHash, vLog.TxHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	block, err := e.addNewBlockToResult(ctx, vLog, blocks, blocksOrder)
	if err != nil {
		return err
	}
	order, err := processor.AddEventDataToBlock(ctx, vLog, block, callData)
	if err != nil {
		return err
	}
	addNewOrder(order, block.BlockHash, blocksOrder)
	return nil
}

func addNewOrder(order *Order, blockHash common.Hash, blocksOrder *map[common.Hash][]Order) {
	(*blocksOrder)[blockHash] = append((*blocksOrder)[blockHash], *order)
}

// addNewEvent adds a new event to the blocks array and order array.
// it returns the block that must be filled with event data
func (e *EventManager) addNewBlockToResult(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) (*Block, error) {
	var block *Block
	var err error
	if !isheadBlockInArray(blocks, vLog.BlockHash, vLog.BlockNumber) {
		// Need to add the block, doesnt mind if inside the blocks because I have to respect the order so insert at end
		//TODO: Check if the block is already in the blocks array and copy it instead of retrieve it again
		block, err = e.blockRetriever.RetrieveFullBlockForEvent(ctx, vLog)
		if err != nil {
			return nil, err
		}
		*blocks = append(*blocks, *block)
	}
	block = &(*blocks)[len(*blocks)-1]
	return block, nil
}

// CallData is a struct that contains the calldata of a transaction
type CallData struct {
	data  []byte
	nonce uint64
	from  common.Address
}

// NewCallData creates a new CallData struct
func NewCallData(data []byte, nonce uint64, from common.Address) *CallData {
	return &CallData{
		data:  data,
		nonce: nonce,
		from:  from,
	}
}

// MethodID returns the method ID of the transaction
func (c *CallData) MethodID() []byte {
	return c.data[:4]
}

// InputData returns the input data of the transaction
func (c *CallData) InputData() []byte {
	return c.data[4:]
}

// Nonce returns the nonce of the transaction
func (c *CallData) Nonce() uint64 {
	return c.nonce
}

// From returns the address of the sender of the transaction
func (c *CallData) From() common.Address {
	return c.from
}

// CallDataExtratorGeth is a CallDataExtractor based on Geth
type CallDataExtratorGeth struct {
	ethClient ethereum.ChainReader
}

// NewCallDataExtratorGeth creates a new CallDataExtrator based on Geth
func NewCallDataExtratorGeth(ethClient ethereum.ChainReader) *CallDataExtratorGeth {
	return &CallDataExtratorGeth{
		ethClient: ethClient,
	}
}

// ExtractCallData get the call data from a transaction
func (e *CallDataExtratorGeth) ExtractCallData(ctx context.Context, blockHash, txHash common.Hash, txIndex uint) (*CallData, error) {
	// Read the tx for this event.
	tx, err := e.ethClient.TransactionInBlock(ctx, blockHash, txIndex)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("error: tx not found in block %s at index %d", blockHash.String(), txIndex)
	}
	//log.Debug("tx: ", tx2string(tx))
	if tx.Hash() != txHash {
		return nil, fmt.Errorf("error: tx hash mismatch. want: %s have: %s", txHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return &CallData{
		data:  tx.Data(),
		nonce: msg.Nonce,
		from:  msg.From,
	}, nil
}

// Function used to convert a transaction to a string to used as input data for unittest
func tx2string(tx *types.Transaction) string { //nolint:unused
	writer := new(bytes.Buffer)
	err := tx.EncodeRLP(writer)
	if err != nil {
		return "error:" + err.Error()
	}
	return common.Bytes2Hex(writer.Bytes())
}
