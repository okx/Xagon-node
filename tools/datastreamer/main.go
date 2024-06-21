package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/datastream"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/tools/datastreamer/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
)

const (
	appName  = "zkevm-data-streamer-tool" //nolint:gosec
	appUsage = "zkevm datastream tool"
)

var (
	configFileFlag = cli.StringFlag{
		Name:        config.FlagCfg,
		Aliases:     []string{"c"},
		Usage:       "Configuration `FILE`",
		DefaultText: "./config/tool.config.toml",
		Required:    true,
	}

	entryFlag = cli.Uint64Flag{
		Name:     "entry",
		Aliases:  []string{"e"},
		Usage:    "Entry `NUMBER`",
		Required: true,
	}

	l2blockFlag = cli.Uint64Flag{
		Name:     "l2block",
		Aliases:  []string{"b"},
		Usage:    "L2Block `NUMBER`",
		Required: true,
	}

	batchFlag = cli.Uint64Flag{
		Name:     "batch",
		Aliases:  []string{"bn"},
		Usage:    "Batch `NUMBER`",
		Required: true,
	}

	dumpFlag = cli.BoolFlag{
		Name:     "dump",
		Aliases:  []string{"d"},
		Usage:    "Dump batch to file",
		Required: false,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{},
			Usage:   "Generate stream file from scratch",
			Action:  generate,
			Flags: []cli.Flag{
				&configFileFlag,
			},
		},
		{
			Name:    "decode-entry-offline",
			Aliases: []string{},
			Usage:   "Decodes an entry offline",
			Action:  decodeEntryOffline,
			Flags: []cli.Flag{
				&configFileFlag,
				&entryFlag,
			},
		},
		{
			Name:    "decode-l2block-offline",
			Aliases: []string{},
			Usage:   "Decodes a l2 block offline",
			Action:  decodeL2BlockOffline,
			Flags: []cli.Flag{
				&configFileFlag,
				&l2blockFlag,
			},
		},
		{
			Name:    "decode-batch-offline",
			Aliases: []string{},
			Usage:   "Decodes a batch offline",
			Action:  decodeBatchOffline,
			Flags: []cli.Flag{
				&configFileFlag,
				&batchFlag,
			},
		},
		{
			Name:    "decode-entry",
			Aliases: []string{},
			Usage:   "Decodes an entry",
			Action:  decodeEntry,
			Flags: []cli.Flag{
				&configFileFlag,
				&entryFlag,
			},
		},
		{
			Name:    "decode-l2block",
			Aliases: []string{},
			Usage:   "Decodes a l2 block",
			Action:  decodeL2Block,
			Flags: []cli.Flag{
				&configFileFlag,
				&l2blockFlag,
			},
		},
		{
			Name:    "decode-batch",
			Aliases: []string{},
			Usage:   "Decodes a batch",
			Action:  decodeBatch,
			Flags: []cli.Flag{
				&configFileFlag,
				&batchFlag,
			},
		},
		{
			Name:    "truncate",
			Aliases: []string{},
			Usage:   "Truncates the stream file",
			Action:  truncate,
			Flags: []cli.Flag{
				&configFileFlag,
				&entryFlag,
			},
		},
		{
			Name:    "dump-batch",
			Aliases: []string{},
			Usage:   "Dumps a batch to file",
			Action:  decodeBatch,
			Flags: []cli.Flag{
				&configFileFlag,
				&batchFlag,
				&dumpFlag,
			},
		},
		{
			Name:    "dump-batch-offline",
			Aliases: []string{},
			Usage:   "Dumps a batch to file offline",
			Action:  decodeBatchOffline,
			Flags: []cli.Flag{
				&configFileFlag,
				&batchFlag,
				&dumpFlag,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func initializeStreamServer(c *config.Config) (*datastreamer.StreamServer, error) {
	// Create a stream server
	streamServer, err := datastreamer.NewServer(c.Offline.Port, c.Offline.Version, c.Offline.ChainID, state.StreamTypeSequencer, c.Offline.Filename, c.Offline.WriteTimeout.Duration, &c.Log)
	if err != nil {
		return nil, err
	}

	err = streamServer.Start()
	if err != nil {
		return nil, err
	}

	return streamServer, nil
}

func generate(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Connect to the database
	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer stateSqlDB.Close()
	stateDBStorage := pgstatestorage.NewPostgresStorage(state.Config{}, stateSqlDB)
	log.Debug("Connected to the database")

	mtDBServerConfig := merkletree.Config{URI: c.MerkleTree.URI}
	var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, mtDBCancel := merkletree.NewMTDBServiceClient(cliCtx.Context, mtDBServerConfig)
	defer func() {
		mtDBCancel()
		mtDBClientConn.Close()
	}()
	stateTree := merkletree.NewStateTree(mtDBServiceClient)
	log.Debug("Connected to the merkle tree")

	stateDB := state.NewState(state.Config{}, stateDBStorage, nil, stateTree, nil, nil)

	// Calculate intermediate state roots
	var imStateRoots map[uint64][]byte
	var imStateRootsMux *sync.Mutex = new(sync.Mutex)
	var wg sync.WaitGroup

	lastL2BlockHeader, err := stateDB.GetLastL2BlockHeader(cliCtx.Context, nil)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	maxL2Block := lastL2BlockHeader.Number.Uint64()
	imStateRoots = make(map[uint64][]byte, maxL2Block)

	// Check if a cache file exists
	if c.MerkleTree.CacheFile != "" {
		// Check if the file exists
		if _, err := os.Stat(c.MerkleTree.CacheFile); os.IsNotExist(err) {
			log.Infof("Cache file %s does not exist", c.MerkleTree.CacheFile)
		} else {
			ReadFile, err := os.ReadFile(c.MerkleTree.CacheFile)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			err = json.Unmarshal(ReadFile, &imStateRoots)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			log.Infof("Cache file %s loaded", c.MerkleTree.CacheFile)
		}
	}

	cacheLength := len(imStateRoots)
	dif := int(maxL2Block) - cacheLength

	log.Infof("Cache length: %d, Max L2Block: %d, Dif: %d", cacheLength, maxL2Block, dif)

	for x := 0; dif > 0 && x < c.MerkleTree.MaxThreads && x < dif; x++ {
		start := uint64((x * dif / c.MerkleTree.MaxThreads) + cacheLength)
		end := uint64(((x + 1) * dif / c.MerkleTree.MaxThreads) + cacheLength - 1)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Debugf("Thread %d: Start: %d, End: %d, Total: %d", i, start, end, end-start)
			getImStateRoots(cliCtx.Context, start, end, &imStateRoots, imStateRootsMux, stateDB)
		}(x)
	}

	wg.Wait()

	// Convert imStateRoots to a json and save it to a file
	if c.MerkleTree.CacheFile != "" {
		jsonFile, _ := json.Marshal(imStateRoots)
		err = os.WriteFile(c.MerkleTree.CacheFile, jsonFile, 0644) // nolint:gosec, gomnd
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	err = state.GenerateDataStreamFile(cliCtx.Context, streamServer, stateDB, false, &imStateRoots, c.Offline.ChainID, c.Offline.UpgradeEtrogBatchNumber)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printColored(color.FgGreen, "Process finished\n")

	return nil
}

func getImStateRoots(ctx context.Context, start, end uint64, isStateRoots *map[uint64][]byte, imStateRootMux *sync.Mutex, stateDB *state.State) {
	for x := start; x <= end; x++ {
		l2Block, err := stateDB.GetL2BlockByNumber(ctx, x, nil)
		if err != nil {
			log.Errorf("Error: %v\n", err)
			os.Exit(1)
		}

		stateRoot := l2Block.Root()
		// Populate intermediate state root
		position := state.GetSystemSCPosition(x)
		imStateRoot, err := stateDB.GetStorageAt(ctx, common.HexToAddress(state.SystemSC), big.NewInt(0).SetBytes(position), stateRoot)
		if err != nil {
			log.Errorf("Error: %v\n", err)
			os.Exit(1)
		}

		if common.BytesToHash(imStateRoot.Bytes()) == state.ZeroHash && x != 0 {
			break
		}

		imStateRootMux.Lock()
		(*isStateRoots)[x] = imStateRoot.Bytes()
		imStateRootMux.Unlock()
	}
}

func decodeEntry(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	client, err := datastreamer.NewClient(c.Online.URI, c.Online.StreamType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = client.Start()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	entry, err := client.ExecCommandGetEntry(cliCtx.Uint64("entry"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printEntry(entry)
	return nil
}

func decodeL2Block(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	client, err := datastreamer.NewClient(c.Online.URI, c.Online.StreamType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = client.Start()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	l2BlockNumber := cliCtx.Uint64("l2block")

	bookMark := &datastream.BookMark{
		Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
		Value: l2BlockNumber,
	}

	marshalledBookMark, err := proto.Marshal(bookMark)
	if err != nil {
		return err
	}

	firstEntry, err := client.ExecCommandGetBookmark(marshalledBookMark)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	printEntry(firstEntry)

	secondEntry, err := client.ExecCommandGetEntry(firstEntry.Number + 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	i := uint64(2) //nolint:gomnd
	for secondEntry.Type == datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION) {
		printEntry(secondEntry)
		entry, err := client.ExecCommandGetEntry(firstEntry.Number + i)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		secondEntry = entry
		i++
	}

	return nil
}

func decodeEntryOffline(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	entry, err := streamServer.GetEntry(cliCtx.Uint64("entry"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printEntry(entry)

	return nil
}

func decodeL2BlockOffline(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	l2BlockNumber := cliCtx.Uint64("l2block")

	bookMark := &datastream.BookMark{
		Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
		Value: l2BlockNumber,
	}

	marshalledBookMark, err := proto.Marshal(bookMark)
	if err != nil {
		return err
	}

	firstEntry, err := streamServer.GetFirstEventAfterBookmark(marshalledBookMark)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	printEntry(firstEntry)

	secondEntry, err := streamServer.GetEntry(firstEntry.Number + 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	i := uint64(2) //nolint:gomnd

	for secondEntry.Type == datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION) {
		printEntry(secondEntry)
		secondEntry, err = streamServer.GetEntry(firstEntry.Number + i)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		i++
	}

	return nil
}

func truncate(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = streamServer.TruncateFile(cliCtx.Uint64("entry"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printColored(color.FgGreen, "File truncated\n")

	return nil
}

func decodeBatch(cliCtx *cli.Context) error {
	var batchData = []byte{}
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	client, err := datastreamer.NewClient(c.Online.URI, c.Online.StreamType)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = client.Start()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	batchNumber := cliCtx.Uint64("batch")

	bookMark := &datastream.BookMark{
		Type:  datastream.BookmarkType_BOOKMARK_TYPE_BATCH,
		Value: batchNumber,
	}

	marshalledBookMark, err := proto.Marshal(bookMark)
	if err != nil {
		return err
	}

	firstEntry, err := client.ExecCommandGetBookmark(marshalledBookMark)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	printEntry(firstEntry)

	batchData = append(batchData, firstEntry.Encode()...)

	secondEntry, err := client.ExecCommandGetEntry(firstEntry.Number + 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	printEntry(secondEntry)

	batchData = append(batchData, secondEntry.Encode()...)

	i := uint64(2) //nolint:gomnd
	for {
		entry, err := client.ExecCommandGetEntry(firstEntry.Number + i)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if entry.Type == state.EntryTypeBookMark {
			if err := proto.Unmarshal(entry.Data, bookMark); err != nil {
				return err
			}
			if bookMark.Type == datastream.BookmarkType_BOOKMARK_TYPE_BATCH {
				break
			}
		}

		secondEntry = entry
		printEntry(secondEntry)
		batchData = append(batchData, secondEntry.Encode()...)
		i++
	}

	// Dump batchdata to a file
	if cliCtx.Bool("dump") {
		err = os.WriteFile(fmt.Sprintf("batch_%d.bin", batchNumber), batchData, 0644) // nolint:gosec, gomnd
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	return nil
}

func decodeBatchOffline(cliCtx *cli.Context) error {
	var batchData = []byte{}
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	batchNumber := cliCtx.Uint64("batch")

	bookMark := &datastream.BookMark{
		Type:  datastream.BookmarkType_BOOKMARK_TYPE_BATCH,
		Value: batchNumber,
	}

	marshalledBookMark, err := proto.Marshal(bookMark)
	if err != nil {
		return err
	}

	firstEntry, err := streamServer.GetFirstEventAfterBookmark(marshalledBookMark)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	printEntry(firstEntry)
	batchData = append(batchData, firstEntry.Encode()...)

	secondEntry, err := streamServer.GetEntry(firstEntry.Number + 1)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	i := uint64(2) //nolint:gomnd
	printEntry(secondEntry)
	batchData = append(batchData, secondEntry.Encode()...)
	for {
		secondEntry, err = streamServer.GetEntry(firstEntry.Number + i)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if secondEntry.Type == state.EntryTypeBookMark {
			if err := proto.Unmarshal(secondEntry.Data, bookMark); err != nil {
				return err
			}
			if bookMark.Type == datastream.BookmarkType_BOOKMARK_TYPE_BATCH {
				break
			}
		}

		printEntry(secondEntry)
		batchData = append(batchData, secondEntry.Encode()...)
		i++
	}

	// Dump batchdata to a file
	if cliCtx.Bool("dump") {
		err = os.WriteFile(fmt.Sprintf("offline_batch_%d.bin", batchNumber), batchData, 0644) // nolint:gosec, gomnd
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	return nil
}

func printEntry(entry datastreamer.FileEntry) {
	switch entry.Type {
	case state.EntryTypeBookMark:
		bookmark := &datastream.BookMark{}
		err := proto.Unmarshal(entry.Data, bookmark)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "BookMark\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Type............: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d (%s)\n", bookmark.Type, datastream.BookmarkType_name[int32(bookmark.Type)]))
		printColored(color.FgGreen, "Value...........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", bookmark.Value))
	case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK):
		l2Block := &datastream.L2Block{}
		err := proto.Unmarshal(entry.Data, l2Block)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "L2 Block\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "L2 Block Number.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.BatchNumber))
		printColored(color.FgGreen, "Timestamp.......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d (%v)\n", l2Block.Timestamp, time.Unix(int64(l2Block.Timestamp), 0)))
		printColored(color.FgGreen, "Delta Timestamp.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.DeltaTimestamp))
		printColored(color.FgGreen, "Min. Timestamp..: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.MinTimestamp))
		printColored(color.FgGreen, "L1 Block Hash...: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToHash(l2Block.L1Blockhash)))
		printColored(color.FgGreen, "L1 InfoTree Idx.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.L1InfotreeIndex))
		printColored(color.FgGreen, "Block Hash......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToHash(l2Block.Hash)))
		printColored(color.FgGreen, "State Root......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToHash(l2Block.StateRoot)))
		printColored(color.FgGreen, "Global Exit Root: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToHash(l2Block.GlobalExitRoot)))
		printColored(color.FgGreen, "Coinbase........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToAddress(l2Block.Coinbase)))
		printColored(color.FgGreen, "Block Gas Limit.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", l2Block.BlockGasLimit))
		printColored(color.FgGreen, "Block Info Root.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToHash(l2Block.BlockInfoRoot)))

		if l2Block.Debug != nil && l2Block.Debug.Message != "" {
			printColored(color.FgGreen, "Debug...........: ")
			printColored(color.FgHiWhite, fmt.Sprintf("%s\n", l2Block.Debug))
		}

	case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_START):
		batch := &datastream.BatchStart{}
		err := proto.Unmarshal(entry.Data, batch)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "Batch Start\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", batch.Number))
		printColored(color.FgGreen, "Batch Type......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", datastream.BatchType_name[int32(batch.Type)]))
		printColored(color.FgGreen, "Fork ID.........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", batch.ForkId))
		printColored(color.FgGreen, "Chain ID........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", batch.ChainId))

		if batch.Debug != nil && batch.Debug.Message != "" {
			printColored(color.FgGreen, "Debug...........: ")
			printColored(color.FgHiWhite, fmt.Sprintf("%s\n", batch.Debug))
		}

	case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_END):
		batch := &datastream.BatchEnd{}
		err := proto.Unmarshal(entry.Data, batch)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "Batch End\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", batch.Number))
		printColored(color.FgGreen, "State Root......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", "0x"+common.Bytes2Hex(batch.StateRoot)))
		printColored(color.FgGreen, "Local Exit Root.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", "0x"+common.Bytes2Hex(batch.LocalExitRoot)))

		if batch.Debug != nil && batch.Debug.Message != "" {
			printColored(color.FgGreen, "Debug...........: ")
			printColored(color.FgHiWhite, fmt.Sprintf("%s\n", batch.Debug))
		}

	case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION):
		dsTx := &datastream.Transaction{}
		err := proto.Unmarshal(entry.Data, dsTx)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "L2 Transaction\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "L2 Block Number.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", dsTx.L2BlockNumber))
		printColored(color.FgGreen, "Index...........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", dsTx.Index))
		printColored(color.FgGreen, "Is Valid........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%t\n", dsTx.IsValid))
		printColored(color.FgGreen, "Data............: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", "0x"+common.Bytes2Hex(dsTx.Encoded)))
		printColored(color.FgGreen, "Effec. Gas Price: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", dsTx.EffectiveGasPricePercentage))
		printColored(color.FgGreen, "IM State Root...: ")
		printColored(color.FgHiWhite, fmt.Sprint("0x"+common.Bytes2Hex(dsTx.ImStateRoot)+"\n"))

		tx, err := state.DecodeTx(common.Bytes2Hex(dsTx.Encoded))
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		sender, err := state.GetSender(*tx)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		printColored(color.FgGreen, "Sender..........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", sender))
		nonce := tx.Nonce()
		printColored(color.FgGreen, "Nonce...........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", nonce))

		if dsTx.Debug != nil && dsTx.Debug.Message != "" {
			printColored(color.FgGreen, "Debug...........: ")
			printColored(color.FgHiWhite, fmt.Sprintf("%s\n", dsTx.Debug))
		}

	case datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_UPDATE_GER):
		updateGer := &datastream.UpdateGER{}
		err := proto.Unmarshal(entry.Data, updateGer)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "Update GER\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.BatchNumber))
		printColored(color.FgGreen, "Timestamp.......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%v (%d)\n", time.Unix(int64(updateGer.Timestamp), 0), updateGer.Timestamp))
		printColored(color.FgGreen, "Global Exit Root: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.Bytes2Hex(updateGer.GlobalExitRoot)))
		printColored(color.FgGreen, "Coinbase........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", common.BytesToAddress(updateGer.Coinbase)))
		printColored(color.FgGreen, "Fork ID.........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.ForkId))
		printColored(color.FgGreen, "Chain ID........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.ChainId))
		printColored(color.FgGreen, "State Root......: ")
		printColored(color.FgHiWhite, fmt.Sprint(common.Bytes2Hex(updateGer.StateRoot)+"\n"))

		if updateGer.Debug != nil && updateGer.Debug.Message != "" {
			printColored(color.FgGreen, "Debug...........: ")
			printColored(color.FgHiWhite, fmt.Sprintf("%s\n", updateGer.Debug))
		}
	}
}

func printColored(color color.Attribute, text string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
	fmt.Print(colored)
}
