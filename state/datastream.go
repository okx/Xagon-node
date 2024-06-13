package state

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/datastream"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
)

const (
	// StreamTypeSequencer represents a Sequencer stream
	StreamTypeSequencer datastreamer.StreamType = 1
	// EntryTypeBookMark represents a bookmark entry
	EntryTypeBookMark datastreamer.EntryType = datastreamer.EtBookmark
	// SystemSC is the system smart contract address
	SystemSC = "0x000000000000000000000000000000005ca1ab1e"
	// posConstant is the constant used to compute the position of the intermediate state root
	posConstant = 1
)

// DSBatch represents a data stream batch
type DSBatch struct {
	Batch
	ForkID         uint64
	EtrogTimestamp *time.Time
}

// DSFullBatch represents a data stream batch ant its L2 blocks
type DSFullBatch struct {
	DSBatch
	L2Blocks []DSL2FullBlock
}

// DSL2FullBlock represents a data stream L2 full block and its transactions
type DSL2FullBlock struct {
	DSL2Block
	Txs []DSL2Transaction
}

// DSL2Block is a full l2 block
type DSL2Block struct {
	BatchNumber     uint64
	L2BlockNumber   uint64
	Timestamp       uint64
	MinTimestamp    uint64
	L1InfoTreeIndex uint32
	L1BlockHash     common.Hash
	GlobalExitRoot  common.Hash
	Coinbase        common.Address
	ForkID          uint64
	ChainID         uint64
	BlockHash       common.Hash
	StateRoot       common.Hash
	BlockGasLimit   uint64
	BlockInfoRoot   common.Hash
}

// DSL2Transaction represents a data stream L2 transaction
type DSL2Transaction struct {
	L2BlockNumber               uint64
	ImStateRoot                 common.Hash
	EffectiveGasPricePercentage uint8
	IsValid                     uint8
	Index                       uint64
	StateRoot                   common.Hash
	EncodedLength               uint32
	Encoded                     []byte
}

// DSState gathers the methods required to interact with the data stream state.
type DSState interface {
	GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*DSL2Block, error)
	GetDSBatches(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, readWIPBatch bool, dbTx pgx.Tx) ([]*DSBatch, error)
	GetDSL2Blocks(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*DSL2Block, error)
	GetDSL2Transactions(ctx context.Context, firstL2Block, lastL2Block uint64, dbTx pgx.Tx) ([]*DSL2Transaction, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root common.Hash) (*big.Int, error)
	GetVirtualBatchParentHash(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetForcedBatchParentHash(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (L1InfoTreeExitRootStorageEntry, error)
}

// GenerateDataStreamFile generates or resumes a data stream file
func GenerateDataStreamFile(ctx context.Context, streamServer *datastreamer.StreamServer, stateDB DSState, readWIPBatch bool, imStateRoots *map[uint64][]byte, chainID uint64, upgradeEtrogBatchNumber uint64) error {
	header := streamServer.GetHeader()

	var currentBatchNumber uint64 = 0
	var lastAddedL2BlockNumber uint64 = 0
	var lastAddedBatchNumber uint64 = 0
	var previousTimestamp uint64 = 0

	if header.TotalEntries == 0 {
		// Get Genesis block
		genesisL2Block, err := stateDB.GetDSGenesisBlock(ctx, nil)
		if err != nil {
			return err
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			return err
		}

		bookMark := &datastream.BookMark{
			Type:  datastream.BookmarkType_BOOKMARK_TYPE_BATCH,
			Value: genesisL2Block.BatchNumber,
		}

		marshalledBookMark, err := proto.Marshal(bookMark)
		if err != nil {
			return err
		}

		_, err = streamServer.AddStreamBookmark(marshalledBookMark)
		if err != nil {
			return err
		}

		genesisBatchStart := &datastream.BatchStart{
			Number:  genesisL2Block.BatchNumber,
			Type:    datastream.BatchType_BATCH_TYPE_UNSPECIFIED,
			ForkId:  genesisL2Block.ForkID,
			ChainId: chainID,
		}

		marshalledGenesisBatchStart, err := proto.Marshal(genesisBatchStart)
		if err != nil {
			return err
		}

		_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_START), marshalledGenesisBatchStart)
		if err != nil {
			return err
		}

		bookMark = &datastream.BookMark{
			Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
			Value: genesisL2Block.L2BlockNumber,
		}

		marshalledBookMark, err = proto.Marshal(bookMark)
		if err != nil {
			return err
		}

		_, err = streamServer.AddStreamBookmark(marshalledBookMark)
		if err != nil {
			return err
		}

		genesisBlock := &datastream.L2Block{
			Number:          genesisL2Block.L2BlockNumber,
			DeltaTimestamp:  0,
			MinTimestamp:    0,
			L1InfotreeIndex: 0,
			Hash:            genesisL2Block.BlockHash.Bytes(),
			StateRoot:       genesisL2Block.StateRoot.Bytes(),
			GlobalExitRoot:  genesisL2Block.GlobalExitRoot.Bytes(),
			Coinbase:        genesisL2Block.Coinbase.Bytes(),
		}

		log.Debugf("Genesis block: %+v", genesisBlock)

		marshalledGenesisBlock, err := proto.Marshal(genesisBlock)
		if err != nil {
			return err
		}

		_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK), marshalledGenesisBlock)
		if err != nil {
			return err
		}

		genesisBatchEnd := &datastream.BatchEnd{
			Number:        genesisL2Block.BatchNumber,
			LocalExitRoot: common.Hash{}.Bytes(),
			StateRoot:     genesisL2Block.StateRoot.Bytes(),
		}

		marshalledGenesisBatchEnd, err := proto.Marshal(genesisBatchEnd)
		if err != nil {
			return err
		}

		_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_END), marshalledGenesisBatchEnd)
		if err != nil {
			return err
		}

		err = streamServer.CommitAtomicOp()
		if err != nil {
			return err
		}
		currentBatchNumber++
	} else {
		latestEntry, err := streamServer.GetEntry(header.TotalEntries - 1)
		if err != nil {
			return err
		}

		log.Infof("Latest entry: %+v", latestEntry)

		switch latestEntry.Type {
		case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_START):
			log.Info("Latest entry type is Batch Start")

			batchStart := &datastream.BatchStart{}
			if err := proto.Unmarshal(latestEntry.Data, batchStart); err != nil {
				return err
			}

			currentBatchNumber = batchStart.Number
		case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_END):
			log.Info("Latest entry type is Batch End")

			batchEnd := &datastream.BatchStart{}
			if err := proto.Unmarshal(latestEntry.Data, batchEnd); err != nil {
				return err
			}

			currentBatchNumber = batchEnd.Number
			currentBatchNumber++
		case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_UPDATE_GER):
			log.Info("Latest entry type is UpdateGER")

			updateGer := &datastream.UpdateGER{}
			if err := proto.Unmarshal(latestEntry.Data, updateGer); err != nil {
				return err
			}

			currentBatchNumber = updateGer.BatchNumber
			currentBatchNumber++
		case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK):
			log.Info("Latest entry type is L2Block")

			l2Block := &datastream.L2Block{}

			if err := proto.Unmarshal(latestEntry.Data, l2Block); err != nil {
				return err
			}

			currentL2BlockNumber := l2Block.Number
			currentBatchNumber = l2Block.BatchNumber
			previousTimestamp = l2Block.Timestamp
			lastAddedL2BlockNumber = currentL2BlockNumber
		case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION):
			log.Info("Latest entry type is Transaction")

			transaction := &datastream.Transaction{}
			if err := proto.Unmarshal(latestEntry.Data, transaction); err != nil {
				return err
			}

			currentL2BlockNumber := transaction.L2BlockNumber
			currentBatchNumber = transaction.L2BlockNumber
			lastAddedL2BlockNumber = currentL2BlockNumber

			// Get Previous l2block timestamp
			bookMark := &datastream.BookMark{
				Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
				Value: currentL2BlockNumber - 1,
			}

			marshalledBookMark, err := proto.Marshal(bookMark)
			if err != nil {
				return err
			}

			prevL2BlockEntryNumber, err := streamServer.GetBookmark(marshalledBookMark)
			if err != nil {
				return err
			}

			prevL2BlockEntry, err := streamServer.GetEntry(prevL2BlockEntryNumber)
			if err != nil {
				return err
			}

			prevL2Block := &datastream.L2Block{}
			if err := proto.Unmarshal(prevL2BlockEntry.Data, prevL2Block); err != nil {
				return err
			}

			previousTimestamp = prevL2Block.Timestamp

		case EntryTypeBookMark:
			log.Info("Latest entry type is BookMark")

			bookMark := &datastream.BookMark{}
			if err := proto.Unmarshal(latestEntry.Data, bookMark); err != nil {
				return err
			}

			if bookMark.Type == datastream.BookmarkType_BOOKMARK_TYPE_BATCH {
				currentBatchNumber = bookMark.Value
			} else {
				log.Fatalf("Latest entry type is an unexpected bookmark type: %v", bookMark.Type)
			}
		default:
			log.Fatalf("Latest entry type is not an expected one: %v", latestEntry.Type)
		}
	}

	var entry uint64 = header.TotalEntries
	var currentGER = common.Hash{}

	if entry > 0 {
		entry--
	}

	var err error
	const limit = 10000

	log.Infof("Current entry number: %d", entry)
	log.Infof("Current batch number: %d", currentBatchNumber)
	log.Infof("Last added L2 block number: %d", lastAddedL2BlockNumber)

	for err == nil {
		// Get Next Batch
		batches, err := stateDB.GetDSBatches(ctx, currentBatchNumber, currentBatchNumber+limit, readWIPBatch, nil)
		if err != nil {
			if err == ErrStateNotSynchronized {
				break
			}
			log.Errorf("Error getting batch %d: %s", currentBatchNumber, err.Error())
			return err
		}

		// Finished?
		if len(batches) == 0 {
			break
		}

		l2Blocks, err := stateDB.GetDSL2Blocks(ctx, batches[0].BatchNumber, batches[len(batches)-1].BatchNumber, nil)
		if err != nil {
			log.Errorf("Error getting L2 blocks for batches starting at %d: %s", batches[0].BatchNumber, err.Error())
			return err
		}

		l2Txs := make([]*DSL2Transaction, 0)
		if len(l2Blocks) > 0 {
			l2Txs, err = stateDB.GetDSL2Transactions(ctx, l2Blocks[0].L2BlockNumber, l2Blocks[len(l2Blocks)-1].L2BlockNumber, nil)
			if err != nil {
				log.Errorf("Error getting L2 transactions for blocks starting at %d: %s", l2Blocks[0].L2BlockNumber, err.Error())
				return err
			}
		}

		// Generate full batches
		fullBatches := computeFullBatches(batches, l2Blocks, l2Txs, lastAddedL2BlockNumber)
		currentBatchNumber += limit

		for b, batch := range fullBatches {
			if batch.BatchNumber <= lastAddedBatchNumber && lastAddedBatchNumber != 0 {
				continue
			} else {
				lastAddedBatchNumber = batch.BatchNumber
			}

			err = streamServer.StartAtomicOp()
			if err != nil {
				return err
			}

			bookMark := &datastream.BookMark{
				Type:  datastream.BookmarkType_BOOKMARK_TYPE_BATCH,
				Value: batch.BatchNumber,
			}

			marshalledBookMark, err := proto.Marshal(bookMark)
			if err != nil {
				return err
			}

			missingBatchBookMark := true
			if b == 0 {
				_, err = streamServer.GetBookmark(marshalledBookMark)
				if err == nil {
					missingBatchBookMark = false
				}
			}

			if missingBatchBookMark {
				_, err = streamServer.AddStreamBookmark(marshalledBookMark)
				if err != nil {
					return err
				}

				batchStart := &datastream.BatchStart{
					Number:  batch.BatchNumber,
					Type:    datastream.BatchType_BATCH_TYPE_REGULAR,
					ForkId:  batch.ForkID,
					ChainId: chainID,
				}

				if batch.ForkID >= FORKID_ETROG && (batch.BatchNumber == 1 || (upgradeEtrogBatchNumber != 0 && batch.BatchNumber == upgradeEtrogBatchNumber)) {
					batchStart.Type = datastream.BatchType_BATCH_TYPE_INJECTED
				}

				if batch.ForcedBatchNum != nil {
					batchStart.Type = datastream.BatchType_BATCH_TYPE_FORCED
				}

				marshalledBatchStart, err := proto.Marshal(batchStart)
				if err != nil {
					return err
				}

				_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_START), marshalledBatchStart)
				if err != nil {
					return err
				}
			}

			if len(batch.L2Blocks) == 0 {
				// Empty batch
				// Check if there is a GER update
				if batch.GlobalExitRoot != currentGER && batch.GlobalExitRoot != (common.Hash{}) {
					updateGER := &datastream.UpdateGER{
						BatchNumber:    batch.BatchNumber,
						Timestamp:      uint64(batch.Timestamp.Unix()),
						GlobalExitRoot: batch.GlobalExitRoot.Bytes(),
						Coinbase:       batch.Coinbase.Bytes(),
						ForkId:         batch.ForkID,
						ChainId:        chainID,
						StateRoot:      batch.StateRoot.Bytes(),
					}

					marshalledUpdateGER, err := proto.Marshal(updateGER)
					if err != nil {
						return err
					}

					_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_UPDATE_GER), marshalledUpdateGER)
					if err != nil {
						return err
					}
					currentGER = batch.GlobalExitRoot
				}
			} else {
				for blockIndex, l2Block := range batch.L2Blocks {
					if l2Block.L2BlockNumber <= lastAddedL2BlockNumber && lastAddedL2BlockNumber != 0 {
						continue
					} else {
						lastAddedL2BlockNumber = l2Block.L2BlockNumber
					}

					l1BlockHash := common.Hash{}
					l1InfoTreeIndex := uint32(0)

					// Get L1 block hash
					if l2Block.ForkID >= FORKID_ETROG {
						isForcedBatch := false
						batchRawData := &BatchRawV2{}

						if batch.BatchNumber == 1 || (upgradeEtrogBatchNumber != 0 && batch.BatchNumber == upgradeEtrogBatchNumber) || batch.ForcedBatchNum != nil {
							isForcedBatch = true
						} else {
							batchRawData, err = DecodeBatchV2(batch.BatchL2Data)
							if err != nil {
								log.Errorf("Failed to decode batch data, err: %v", err)
								return err
							}
						}

						if !isForcedBatch {
							// Get current block by index
							l2blockRaw := batchRawData.Blocks[blockIndex]
							l1InfoTreeIndex = l2blockRaw.IndexL1InfoTree
							if l2blockRaw.IndexL1InfoTree != 0 {
								l1InfoTreeExitRootStorageEntry, err := stateDB.GetL1InfoRootLeafByIndex(ctx, l2blockRaw.IndexL1InfoTree, nil)
								if err != nil {
									return err
								}
								l1BlockHash = l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.PreviousBlockHash
							}
						} else {
							// Initial batch must be handled differently
							if batch.BatchNumber == 1 || (upgradeEtrogBatchNumber != 0 && batch.BatchNumber == upgradeEtrogBatchNumber) {
								l1BlockHash, err = stateDB.GetVirtualBatchParentHash(ctx, batch.BatchNumber, nil)
								if err != nil {
									return err
								}
							} else {
								l1BlockHash, err = stateDB.GetForcedBatchParentHash(ctx, *batch.ForcedBatchNum, nil)
								if err != nil {
									return err
								}
							}
						}
					}

					streamL2Block := &datastream.L2Block{
						Number:          l2Block.L2BlockNumber,
						BatchNumber:     l2Block.BatchNumber,
						Timestamp:       l2Block.Timestamp,
						DeltaTimestamp:  uint32(l2Block.Timestamp - previousTimestamp),
						MinTimestamp:    uint64(batch.Timestamp.Unix()),
						L1Blockhash:     l1BlockHash.Bytes(),
						L1InfotreeIndex: l1InfoTreeIndex,
						Hash:            l2Block.BlockHash.Bytes(),
						StateRoot:       l2Block.StateRoot.Bytes(),
						GlobalExitRoot:  l2Block.GlobalExitRoot.Bytes(),
						Coinbase:        l2Block.Coinbase.Bytes(),
						BlockInfoRoot:   l2Block.BlockInfoRoot.Bytes(),
						BlockGasLimit:   l2Block.BlockGasLimit,
					}

					// Keep the l2 block hash as it is, as the state root can be found in the StateRoot field
					// So disable this
					/*
						if l2Block.ForkID >= FORKID_ETROG {
							streamL2Block.Hash = l2Block.StateRoot.Bytes()
						}
					*/

					if l2Block.ForkID == FORKID_ETROG && batch.EtrogTimestamp != nil {
						streamL2Block.MinTimestamp = uint64(batch.EtrogTimestamp.Unix())
					}

					if l2Block.ForkID >= FORKID_ETROG && l2Block.L1InfoTreeIndex == 0 {
						streamL2Block.MinTimestamp = 0
					}

					previousTimestamp = l2Block.Timestamp

					bookMark := &datastream.BookMark{
						Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
						Value: streamL2Block.Number,
					}

					marshalledBookMark, err := proto.Marshal(bookMark)
					if err != nil {
						return err
					}

					// Check if l2 block was already added
					_, err = streamServer.GetBookmark(marshalledBookMark)
					if err == nil {
						continue
					}

					_, err = streamServer.AddStreamBookmark(marshalledBookMark)
					if err != nil {
						return err
					}

					marshalledL2Block, err := proto.Marshal(streamL2Block)
					if err != nil {
						return err
					}

					_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK), marshalledL2Block)
					if err != nil {
						return err
					}

					for _, tx := range l2Block.Txs {
						// < ETROG => IM State root is retrieved from the system SC (using cache is available)
						// = ETROG => IM State root is retrieved from the receipt.post_state => Do nothing
						// > ETROG => IM State root is retrieved from the receipt.im_state_root
						if l2Block.ForkID < FORKID_ETROG {
							// Populate intermediate state root with information from the system SC (or cache if available)
							if imStateRoots == nil || (*imStateRoots)[streamL2Block.Number] == nil {
								position := GetSystemSCPosition(l2Block.L2BlockNumber)
								imStateRoot, err := stateDB.GetStorageAt(ctx, common.HexToAddress(SystemSC), big.NewInt(0).SetBytes(position), l2Block.StateRoot)
								if err != nil {
									return err
								}
								tx.StateRoot = common.BigToHash(imStateRoot)
							} else {
								tx.StateRoot = common.BytesToHash((*imStateRoots)[streamL2Block.Number])
							}
						} else if l2Block.ForkID > FORKID_ETROG {
							tx.StateRoot = tx.ImStateRoot
						}

						transaction := &datastream.Transaction{
							L2BlockNumber:               tx.L2BlockNumber,
							Index:                       tx.Index,
							IsValid:                     tx.IsValid != 0,
							Encoded:                     tx.Encoded,
							EffectiveGasPricePercentage: uint32(tx.EffectiveGasPricePercentage),
							ImStateRoot:                 tx.StateRoot.Bytes(),
						}

						// Clear the state root if the ForkID is > ETROG
						if l2Block.ForkID > FORKID_ETROG {
							transaction.ImStateRoot = common.Hash{}.Bytes()
						}

						marshalledTransaction, err := proto.Marshal(transaction)
						if err != nil {
							return err
						}

						_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION), marshalledTransaction)
						if err != nil {
							return err
						}
					}

					currentGER = l2Block.GlobalExitRoot
				}
			}

			batchEnd := &datastream.BatchEnd{
				Number:        batch.BatchNumber,
				LocalExitRoot: batch.LocalExitRoot.Bytes(),
				StateRoot:     batch.StateRoot.Bytes(),
			}

			marshalledBatch, err := proto.Marshal(batchEnd)
			if err != nil {
				return err
			}

			_, err = streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_END), marshalledBatch)
			if err != nil {
				return err
			}

			// Commit at the end of each batch group
			err = streamServer.CommitAtomicOp()
			if err != nil {
				return err
			}
		}
	}

	return err
}

// GetSystemSCPosition computes the position of the intermediate state root for the system smart contract
func GetSystemSCPosition(blockNumber uint64) []byte {
	v1 := big.NewInt(0).SetUint64(blockNumber).Bytes()
	v2 := big.NewInt(0).SetUint64(uint64(posConstant)).Bytes()

	// Add 0s to make v1 and v2 32 bytes long
	for len(v1) < 32 {
		v1 = append([]byte{0}, v1...)
	}
	for len(v2) < 32 {
		v2 = append([]byte{0}, v2...)
	}

	return keccak256.Hash(v1, v2)
}

// computeFullBatches computes the full batches
func computeFullBatches(batches []*DSBatch, l2Blocks []*DSL2Block, l2Txs []*DSL2Transaction, prevL2BlockNumber uint64) []*DSFullBatch {
	currentL2Tx := 0
	currentL2Block := uint64(0)

	fullBatches := make([]*DSFullBatch, 0)

	for _, batch := range batches {
		fullBatch := &DSFullBatch{
			DSBatch: *batch,
		}

		for i := currentL2Block; i < uint64(len(l2Blocks)); i++ {
			l2Block := l2Blocks[i]

			if prevL2BlockNumber != 0 && l2Block.L2BlockNumber <= prevL2BlockNumber {
				continue
			}

			if l2Block.BatchNumber == batch.BatchNumber {
				fullBlock := DSL2FullBlock{
					DSL2Block: *l2Block,
				}

				for j := currentL2Tx; j < len(l2Txs); j++ {
					l2Tx := l2Txs[j]
					if l2Tx.L2BlockNumber == l2Block.L2BlockNumber {
						fullBlock.Txs = append(fullBlock.Txs, *l2Tx)
						currentL2Tx++
					}
					if l2Tx.L2BlockNumber > l2Block.L2BlockNumber {
						break
					}
				}

				fullBatch.L2Blocks = append(fullBatch.L2Blocks, fullBlock)
				prevL2BlockNumber = l2Block.L2BlockNumber
				currentL2Block++
			} else if l2Block.BatchNumber > batch.BatchNumber {
				break
			}
		}
		fullBatches = append(fullBatches, fullBatch)
	}

	return fullBatches
}
