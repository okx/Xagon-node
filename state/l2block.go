package state

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type gethHeader struct {
	*types.Header
}
type gethBlock struct {
	*types.Block
}

// L2Header represents a block header in the L2.
type L2Header struct {
	*gethHeader
	GlobalExitRoot common.Hash `json:"globalExitRoot"`
	BlockInfoRoot  common.Hash `json:"blockInfoRoot"`
}

// NewL2Header creates an instance of L2Header from a types.Header
func NewL2Header(h *types.Header) *L2Header {
	return &L2Header{gethHeader: &gethHeader{types.CopyHeader(h)}}
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *L2Header) Hash() common.Hash {
	return h.gethHeader.Hash()
}

// MarshalJSON encodes a json object
func (h *L2Header) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	if h.gethHeader != nil && h.gethHeader.Header != nil {
		b, err := json.Marshal(h.gethHeader.Header)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
	}

	m["globalExitRoot"] = h.GlobalExitRoot.String()
	m["blockInfoRoot"] = h.BlockInfoRoot.String()

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalJSON decodes a json object
func (h *L2Header) UnmarshalJSON(input []byte) error {
	str := strings.Trim(string(input), "\"")
	if strings.ToLower(strings.TrimSpace(str)) == "null" {
		return nil
	}

	var header *types.Header
	err := json.Unmarshal(input, &header)
	if err != nil {
		return err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(input, &m)
	if err != nil {
		return err
	}

	h.gethHeader = &gethHeader{header}
	if globalExitRoot, found := m["globalExitRoot"]; found {
		h.GlobalExitRoot = common.HexToHash(globalExitRoot.(string))
	}
	if blockInfoRoot, found := m["blockInfoRoot"]; found {
		h.BlockInfoRoot = common.HexToHash(blockInfoRoot.(string))
	}

	return nil
}

// L2Block represents a block from L2
type L2Block struct {
	*gethBlock
	header *L2Header
	uncles []*L2Header

	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// GlobalExitRoot returns the header GlobalExitRoot
func (b *L2Block) GlobalExitRoot() common.Hash {
	return b.Header().GlobalExitRoot
}

// BlockInfoRoot returns the header BlockInfoRoot
func (b *L2Block) BlockInfoRoot() common.Hash {
	return b.Header().BlockInfoRoot
}

// Header returns the block header (as a copy).
func (b *L2Block) Header() *L2Header {
	return CopyHeader(b.header)
}

// Number returns the block header number.
func (b *L2Block) Number() *big.Int {
	return b.header.Number
}

// NewL2Block creates a new block. The input data is copied, changes to header and to the
// field values will not affect the block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewL2Block(h *L2Header, txs []*types.Transaction, uncles []*L2Header, receipts []*types.Receipt, hasher types.TrieHasher) *L2Block {
	l2Uncles := make([]*L2Header, 0, len(uncles))
	gethUncles := make([]*types.Header, 0, len(uncles))
	for _, uncle := range uncles {
		l2Uncles = append(l2Uncles, CopyHeader(uncle))
		gethUncles = append(gethUncles, types.CopyHeader(uncle.gethHeader.Header))
	}

	cpy := CopyHeader(h)
	b := types.NewBlock(h.gethHeader.Header, txs, gethUncles, receipts, hasher)
	cpy.gethHeader = &gethHeader{b.Header()}
	return &L2Block{
		header:    cpy,
		gethBlock: &gethBlock{b},
		uncles:    l2Uncles,
	}
}

// NewL2BlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewL2BlockWithHeader(h *L2Header) *L2Block {
	b := types.NewBlockWithHeader(h.gethHeader.Header)
	return &L2Block{
		header:    CopyHeader(h),
		gethBlock: &gethBlock{b},
	}
}

// WithBody returns a copy of the block with the given transaction and uncle contents.
func (b *L2Block) WithBody(transactions []*types.Transaction, uncles []*L2Header) *L2Block {
	l2Uncles := make([]*L2Header, 0, len(uncles))
	gethUncles := make([]*types.Header, 0, len(uncles))
	for _, uncle := range uncles {
		l2Uncles = append(l2Uncles, CopyHeader(uncle))
		gethUncles = append(gethUncles, types.CopyHeader(uncle.gethHeader.Header))
	}

	return &L2Block{
		header:    b.header,
		gethBlock: &gethBlock{b.gethBlock.WithBody(transactions, gethUncles)},
		uncles:    l2Uncles,
	}
}

// CopyHeader creates a deep copy of a block header.
func CopyHeader(h *L2Header) *L2Header {
	if h == nil {
		return nil
	}
	cpy := *h
	cpy.gethHeader = &gethHeader{types.CopyHeader(h.gethHeader.Header)}
	cpy.GlobalExitRoot = h.GlobalExitRoot
	cpy.BlockInfoRoot = h.BlockInfoRoot
	return &cpy
}

const newL2BlocksCheckInterval = 200 * time.Millisecond

// NewL2BlockEventHandler represent a func that will be called by the
// state when a NewL2BlockEvent is triggered
type NewL2BlockEventHandler func(e NewL2BlockEvent)

// NewL2BlockEvent is a struct provided from the state to the NewL2BlockEventHandler
// when a new l2 block is detected with data related to this new l2 block.
type NewL2BlockEvent struct {
	Block L2Block
	Logs  []*types.Log
}

// StartToMonitorNewL2Blocks starts 2 go routines that will
// monitor new blocks and execute handlers registered to be executed
// when a new l2 block is detected. This is used by the RPC WebSocket
// filter subscription but can be used by any other component that
// needs to react to a new L2 block added to the state.
func (s *State) StartToMonitorNewL2Blocks() {
	go InfiniteSafeRun(s.monitorNewL2Blocks, "fail to monitor new l2 blocks: %v:", time.Second)
	go InfiniteSafeRun(s.handleEvents, "fail to handle events: %v", time.Second)
}

// RegisterNewL2BlockEventHandler add the provided handler to the list of handlers
// that will be triggered when a new l2 block event is triggered
func (s *State) RegisterNewL2BlockEventHandler(h NewL2BlockEventHandler) {
	log.Info("new l2 block event handler registered")
	s.newL2BlockEventHandlers = append(s.newL2BlockEventHandlers, h)
}

func (s *State) monitorNewL2Blocks() {
	waitNextCycle := func() {
		time.Sleep(newL2BlocksCheckInterval)
	}

	lastL2BlockNumber, err := s.GetLastL2BlockNumber(context.Background(), nil)
	if errors.Is(err, ErrStateNotSynchronized) {
		lastL2BlockNumber = 0
	} else if err != nil {
		log.Fatalf("failed to load the last l2 block: %v", err)
	}
	lastL2BlockNumberSeen := lastL2BlockNumber

	for {
		if len(s.newL2BlockEventHandlers) == 0 {
			waitNextCycle()
			continue
		}

		lastL2BlockNumber, err := s.GetLastL2BlockNumber(context.Background(), nil)
		if errors.Is(err, ErrStateNotSynchronized) {
			waitNextCycle()
			continue
		} else if err != nil {
			log.Errorf("failed to get last l2 block while monitoring new blocks: %v", err)
			waitNextCycle()
			continue
		}

		// not updates until now
		if lastL2BlockNumber == 0 || lastL2BlockNumberSeen >= lastL2BlockNumber {
			waitNextCycle()
			continue
		}

		fromBlockNumber := lastL2BlockNumberSeen + uint64(1)
		toBlockNumber := lastL2BlockNumber
		log.Debugf("[monitorNewL2Blocks] new l2 block detected from block %v to %v", fromBlockNumber, toBlockNumber)

		for bn := fromBlockNumber; bn <= toBlockNumber; bn++ {
			block, err := s.GetL2BlockByNumber(context.Background(), bn, nil)
			if err != nil {
				log.Errorf("failed to get l2 block while monitoring new blocks: %v", err)
				break
			}
			logs, err := s.GetLogsByBlockNumber(context.Background(), bn, nil)
			if err != nil {
				log.Errorf("failed to get l2 block while monitoring new blocks: %v", err)
				break
			}

			log.Debugf("[monitorNewL2Blocks] sending NewL2BlockEvent for block %v", block.NumberU64())
			start := time.Now()
			s.newL2BlockEvents <- NewL2BlockEvent{
				Block: *block,
				Logs:  logs,
			}
			lastL2BlockNumberSeen = block.NumberU64()
			log.Debugf("[monitorNewL2Blocks] NewL2BlockEvent for block %v took %v to be sent", block.NumberU64(), time.Since(start))
			log.Infof("new l2 block detected: number %v, hash %v", block.NumberU64(), block.Hash().String())
		}

		// interval to check for new l2 blocks
		waitNextCycle()
	}
}

func (s *State) handleEvents() {
	for newL2BlockEvent := range s.newL2BlockEvents {
		log.Debugf("[handleEvents] new l2 block event detected for block: %v", newL2BlockEvent.Block.NumberU64())
		if len(s.newL2BlockEventHandlers) == 0 {
			continue
		}

		wg := sync.WaitGroup{}
		for _, handler := range s.newL2BlockEventHandlers {
			wg.Add(1)
			go func(h NewL2BlockEventHandler, e NewL2BlockEvent) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {
						log.Errorf("failed and recovered in NewL2BlockEventHandler: %v", r)
					}
				}()
				log.Debugf("[handleEvents] triggering new l2 block event handler for block: %v", e.Block.NumberU64())
				start := time.Now()
				h(e)
				log.Debugf("[handleEvents] new l2 block event handler for block %v took %v to be executed", e.Block.NumberU64(), time.Since(start))
			}(handler, newL2BlockEvent)
		}
		wg.Wait()
	}
}
