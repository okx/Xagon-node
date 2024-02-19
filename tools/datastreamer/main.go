package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	nodeConfig "github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/tools/datastreamer/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
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

	genesisFileFlag = cli.StringFlag{
		Name:        config.FlagGenesis,
		Aliases:     []string{"g"},
		Usage:       "Genesis `FILE`",
		DefaultText: "./config/genesis.json",
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

	updateFileFlag = cli.BoolFlag{
		Name:     "update",
		Aliases:  []string{"u"},
		Usage:    "Update `FILE`",
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
			Name:    "reprocess",
			Aliases: []string{},
			Usage:   "Reprocess l2block since a given l2block number",
			Action:  reprocess,
			Flags: []cli.Flag{
				&configFileFlag,
				&genesisFileFlag,
				&l2blockFlag,
				&updateFileFlag,
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
			Name:    "truncate",
			Aliases: []string{},
			Usage:   "Truncates the stream file",
			Action:  truncate,
			Flags: []cli.Flag{
				&configFileFlag,
				&entryFlag,
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
	streamServer, err := datastreamer.NewServer(c.Offline.Port, c.Offline.Version, c.Offline.ChainID, state.StreamTypeSequencer, c.Offline.Filename, &c.Log)
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

	err = state.GenerateDataStreamerFile(cliCtx.Context, streamServer, stateDB, false, &imStateRoots, c.Offline.ChainID, c.Offline.UpgradeEtrogBatchNumber) // nolint:gomnd
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
		imStateRootMux.Lock()
		(*isStateRoots)[x] = imStateRoot.Bytes()
		imStateRootMux.Unlock()
	}
}

func reprocess(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Init(c.Log)

	ctx := cliCtx.Context

	genesisFileAsStr, err := nodeConfig.LoadGenesisFileAsString(cliCtx.String(config.FlagGenesis))
	if err != nil {
		fmt.Printf("failed to load genesis file. Error: %v", err)
		os.Exit(1)
	}

	networkConfig, err := nodeConfig.LoadGenesisFromJSONString(genesisFileAsStr)
	if err != nil {
		fmt.Printf("failed to load genesis configuration from file. Error: %v", err)
		os.Exit(1)
	}

	currentL2BlockNumber := cliCtx.Uint64("l2block")
	var stateRoot []byte

	streamServer, err := initializeStreamServer(c)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if currentL2BlockNumber == 0 {
		printColored(color.FgHiYellow, "\n\nSetting Genesis block\n\n")

		mtDBServerConfig := merkletree.Config{URI: c.MerkleTree.URI}
		var mtDBCancel context.CancelFunc
		mtDBServiceClient, mtDBClientConn, mtDBCancel := merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
		defer func() {
			mtDBCancel()
			mtDBClientConn.Close()
		}()

		stateTree := merkletree.NewStateTree(mtDBServiceClient)

		stateRoot, err = setGenesis(ctx, stateTree, networkConfig.Genesis)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Get Genesis block from the file and validate the state root
		bookMark := state.DSBookMark{
			Type:  state.BookMarkTypeL2Block,
			Value: 0,
		}

		firstEntry, err := streamServer.GetFirstEventAfterBookmark(bookMark.Encode())
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
		printEntry(secondEntry)

		if common.Bytes2Hex(stateRoot) != common.Bytes2Hex(secondEntry.Data[40:72]) {
			printColored(color.FgRed, "\nError: Genesis state root does not match\n\n")
			os.Exit(1)
		} else {
			printColored(color.FgGreen, "\nGenesis state root matches\n\n")
		}
		currentL2BlockNumber++
	}

	// Connect to the executor
	executorClient, executorClientConn, executorCancel := executor.NewExecutorClient(ctx, c.Executor)
	defer func() {
		executorCancel()
		executorClientConn.Close()
	}()

	bookMark := state.DSBookMark{
		Type:  state.BookMarkTypeL2Block,
		Value: currentL2BlockNumber,
	}

	startEntry, err := streamServer.GetFirstEventAfterBookmark(bookMark.Encode())
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var previousStateRoot = stateRoot
	var maxEntry = streamServer.GetHeader().TotalEntries

	for x := startEntry.Number; x < maxEntry; x++ {
		printColored(color.FgHiYellow, fmt.Sprintf("\nProcessing entity: %d\n", x))

		currentEntry, err := streamServer.GetEntry(x)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		var processBatchRequest *executor.ProcessBatchRequest
		var expectedNewRoot []byte
		var entryToUpdate *datastreamer.FileEntry

		switch currentEntry.Type {
		case state.EntryTypeBookMark:
			printEntry(currentEntry)
			entryToUpdate = nil
			continue
		case state.EntryTypeUpdateGER:
			printEntry(currentEntry)
			processBatchRequest = &executor.ProcessBatchRequest{
				OldBatchNum:      binary.BigEndian.Uint64(currentEntry.Data[0:8]) - 1,
				Coinbase:         common.Bytes2Hex(currentEntry.Data[48:68]),
				BatchL2Data:      nil,
				OldStateRoot:     previousStateRoot,
				GlobalExitRoot:   currentEntry.Data[16:48],
				OldAccInputHash:  []byte{},
				EthTimestamp:     binary.BigEndian.Uint64(currentEntry.Data[8:16]),
				UpdateMerkleTree: uint32(1),
				ChainId:          c.Offline.ChainID,
				ForkId:           uint64(binary.BigEndian.Uint16(currentEntry.Data[68:70])),
			}

			expectedNewRoot = currentEntry.Data[70:102]
			entryToUpdate = nil
		case state.EntryTypeL2BlockStart:
			startEntry = currentEntry
			printEntry(startEntry)

			txEntry, err := streamServer.GetEntry(startEntry.Number + 1)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			printEntry(txEntry)

			endEntry, err := streamServer.GetEntry(startEntry.Number + 2) //nolint:gomnd
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			printEntry(endEntry)

			forkID := uint64(binary.BigEndian.Uint16(startEntry.Data[76:78]))

			tx, err := state.DecodeTx(common.Bytes2Hex((txEntry.Data[6:])))
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// Get the old state root
			oldStateRoot := getOldStateRoot(startEntry.Number, streamServer)

			// RLP encode the transaction using the proper fork id
			batchL2Data, err := state.EncodeTransaction(*tx, txEntry.Data[0], forkID) //nolint:gomnd
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			processBatchRequest = &executor.ProcessBatchRequest{
				OldBatchNum:      binary.BigEndian.Uint64(startEntry.Data[0:8]) - 1,
				Coinbase:         common.Bytes2Hex(startEntry.Data[56:76]),
				BatchL2Data:      batchL2Data,
				OldStateRoot:     oldStateRoot,
				GlobalExitRoot:   startEntry.Data[24:56],
				OldAccInputHash:  []byte{},
				EthTimestamp:     binary.BigEndian.Uint64(startEntry.Data[16:24]),
				UpdateMerkleTree: uint32(1),
				ChainId:          c.Offline.ChainID,
				ForkId:           uint64(binary.BigEndian.Uint16(startEntry.Data[76:78])),
			}

			expectedNewRoot = endEntry.Data[40:72]
			entryToUpdate = &endEntry
			x += 2 //nolint:gomnd
		}

		// Process batch
		processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			fmt.Printf("Error: %v\n", processBatchResponse.Error)
			os.Exit(1)
		}

		if common.Bytes2Hex(processBatchResponse.NewStateRoot) != common.Bytes2Hex(expectedNewRoot) {
			printColored(color.FgRed, "\nNew state root does not match\n\n")
			printColored(color.FgRed, fmt.Sprintf("Old State Root.........: %s\n", "0x"+common.Bytes2Hex(processBatchRequest.GetOldStateRoot())))
			printColored(color.FgRed, fmt.Sprintf("New State Root.........: %s\n", "0x"+common.Bytes2Hex(processBatchResponse.NewStateRoot)))
			printColored(color.FgRed, fmt.Sprintf("Expected New State Root: %s\n", "0x"+common.Bytes2Hex(expectedNewRoot)))
			// Check if we must update the file with the new state root
			if cliCtx.Bool("update") {
				if entryToUpdate.Type != state.EntryTypeL2BlockEnd {
					printColored(color.FgRed, "Error: Entry to update is not a L2BlockEnd\n")
					os.Exit(1)
				}
				blockEnd := state.DSL2BlockEnd{}.Decode(entryToUpdate.Data)
				blockEnd.StateRoot = common.BytesToHash(processBatchResponse.NewStateRoot)
				err = streamServer.UpdateEntryData(entryToUpdate.Number, state.EntryTypeL2BlockEnd, blockEnd.Encode())
				if err != nil {
					printColored(color.FgRed, fmt.Sprintf("Error: %v\n", err))
					os.Exit(1)
				}
			} else {
				break
			}
		} else {
			printColored(color.FgGreen, "New state root matches\n")
			previousStateRoot = processBatchResponse.NewStateRoot
		}
	}

	return nil
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

	client.FromEntry = cliCtx.Uint64("entry")
	err = client.ExecCommand(datastreamer.CmdEntry)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	printEntry(client.Entry)
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

	bookMark := state.DSBookMark{
		Type:  state.BookMarkTypeL2Block,
		Value: l2BlockNumber,
	}

	client.FromBookmark = bookMark.Encode()
	err = client.ExecCommand(datastreamer.CmdBookmark)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	firstEntry := client.Entry
	printEntry(firstEntry)

	client.FromEntry = firstEntry.Number + 1
	err = client.ExecCommand(datastreamer.CmdEntry)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	secondEntry := client.Entry
	printEntry(secondEntry)

	i := uint64(2) //nolint:gomnd
	for secondEntry.Type == state.EntryTypeL2Tx {
		client.FromEntry = firstEntry.Number + i
		err = client.ExecCommand(datastreamer.CmdEntry)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		secondEntry = client.Entry
		printEntry(secondEntry)
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

	bookMark := state.DSBookMark{
		Type:  state.BookMarkTypeL2Block,
		Value: l2BlockNumber,
	}

	firstEntry, err := streamServer.GetFirstEventAfterBookmark(bookMark.Encode())
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
	printEntry(secondEntry)
	for secondEntry.Type == state.EntryTypeL2Tx {
		secondEntry, err = streamServer.GetEntry(firstEntry.Number + i)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		printEntry(secondEntry)
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

func printEntry(entry datastreamer.FileEntry) {
	var bookmarkTypeDesc = map[byte]string{
		state.BookMarkTypeL2Block: "L2 Block Number",
		state.BookMarkTypeBatch:   "Batch Number",
	}

	switch entry.Type {
	case state.EntryTypeBookMark:
		bookmark := state.DSBookMark{}.Decode(entry.Data)
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "BookMark\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Type............: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d (%s)\n", bookmark.Type, bookmarkTypeDesc[bookmark.Type]))
		printColored(color.FgGreen, "Value...........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", bookmark.Value))
	case state.EntryTypeL2BlockStart:
		blockStart := state.DSL2BlockStart{}.Decode(entry.Data)
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "L2 Block Start\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.BatchNumber))
		printColored(color.FgGreen, "L2 Block Number.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.L2BlockNumber))
		printColored(color.FgGreen, "Timestamp.......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%v (%d)\n", time.Unix(blockStart.Timestamp, 0), blockStart.Timestamp))
		printColored(color.FgGreen, "Delta Timestamp.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.DeltaTimestamp))
		printColored(color.FgGreen, "L1 InfoTree Idx.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.L1InfoTreeIndex))
		printColored(color.FgGreen, "L1 Block Hash...: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", blockStart.L1BlockHash))
		printColored(color.FgGreen, "Global Exit Root: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", blockStart.GlobalExitRoot))
		printColored(color.FgGreen, "Coinbase........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", blockStart.Coinbase))
		printColored(color.FgGreen, "Fork ID.........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.ForkID))
		printColored(color.FgGreen, "Chain ID........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockStart.ChainID))
	case state.EntryTypeL2Tx:
		dsTx := state.DSL2Transaction{}.Decode(entry.Data)
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "L2 Transaction\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Effec. Gas Price: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", dsTx.EffectiveGasPricePercentage))
		printColored(color.FgGreen, "Is Valid........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%t\n", dsTx.IsValid == 1))
		printColored(color.FgGreen, "IM State Root...: ")
		printColored(color.FgHiWhite, fmt.Sprint(dsTx.StateRoot.Hex()+"\n"))
		printColored(color.FgGreen, "Encoded Length..: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", dsTx.EncodedLength))
		printColored(color.FgGreen, "Encoded.........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", "0x"+common.Bytes2Hex(dsTx.Encoded)))

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
	case state.EntryTypeL2BlockEnd:
		blockEnd := state.DSL2BlockEnd{}.Decode(entry.Data)
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "L2 Block End\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "L2 Block Number.: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", blockEnd.L2BlockNumber))
		printColored(color.FgGreen, "L2 Block Hash...: ")
		printColored(color.FgHiWhite, fmt.Sprint(blockEnd.BlockHash.Hex()+"\n"))
		printColored(color.FgGreen, "State Root......: ")
		printColored(color.FgHiWhite, fmt.Sprint(blockEnd.StateRoot.Hex()+"\n"))
	case state.EntryTypeUpdateGER:
		updateGer := state.DSUpdateGER{}.Decode(entry.Data)
		printColored(color.FgGreen, "Entry Type......: ")
		printColored(color.FgHiYellow, "Update GER\n")
		printColored(color.FgGreen, "Entry Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", entry.Number))
		printColored(color.FgGreen, "Batch Number....: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.BatchNumber))
		printColored(color.FgGreen, "Timestamp.......: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%v (%d)\n", time.Unix(updateGer.Timestamp, 0), updateGer.Timestamp))
		printColored(color.FgGreen, "Global Exit Root: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", updateGer.GlobalExitRoot))
		printColored(color.FgGreen, "Coinbase........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%s\n", updateGer.Coinbase))
		printColored(color.FgGreen, "Fork ID.........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.ForkID))
		printColored(color.FgGreen, "Chain ID........: ")
		printColored(color.FgHiWhite, fmt.Sprintf("%d\n", updateGer.ChainID))
		printColored(color.FgGreen, "State Root......: ")
		printColored(color.FgHiWhite, fmt.Sprint(updateGer.StateRoot.Hex()+"\n"))
	}
}

func printColored(color color.Attribute, text string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
	fmt.Print(colored)
}

// setGenesis populates state with genesis information
func setGenesis(ctx context.Context, tree *merkletree.StateTree, genesis state.Genesis) ([]byte, error) {
	var (
		root    common.Hash
		newRoot []byte
		err     error
	)

	if tree == nil {
		return newRoot, fmt.Errorf("state tree is nil")
	}

	uuid := uuid.New().String()

	for _, action := range genesis.Actions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = tree.SetBalance(ctx, address, balance, newRoot, uuid)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = tree.SetNonce(ctx, address, nonce, newRoot, uuid)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return newRoot, fmt.Errorf("could not decode SC bytecode for address %q: %v", address, err)
			}
			newRoot, _, err = tree.SetCode(ctx, address, code, newRoot, uuid)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeStorage):
			// Parse position and value
			positionBI, err := encoding.DecodeBigIntHexOrDecimal(action.StoragePosition)
			if err != nil {
				return newRoot, err
			}
			valueBI, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			// Store
			newRoot, _, err = tree.SetStorageAt(ctx, address, positionBI, valueBI, newRoot, uuid)
			if err != nil {
				return newRoot, err
			}
		default:
			return newRoot, fmt.Errorf("unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(newRoot)

	// flush state db
	err = tree.Flush(ctx, root, uuid)
	if err != nil {
		fmt.Printf("error flushing state tree after genesis: %v", err)
		return newRoot, err
	}

	return newRoot, nil
}

func getOldStateRoot(entityNumber uint64, streamServer *datastreamer.StreamServer) []byte {
	var found = false
	var entry datastreamer.FileEntry
	var err error

	for !found && entityNumber > 1 {
		entityNumber--
		entry, err = streamServer.GetEntry(entityNumber)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if entry.Type == state.EntryTypeL2BlockEnd || entry.Type == state.EntryTypeUpdateGER {
			found = true
		}
	}

	if !found {
		fmt.Printf("Error: Could not find old state root")
		os.Exit(1)
	}

	printColored(color.FgHiYellow, "Getting Old State Root from\n")
	printEntry(entry)

	if entry.Type == state.EntryTypeUpdateGER {
		return entry.Data[70:102]
	}

	return entry.Data[40:72]
}
