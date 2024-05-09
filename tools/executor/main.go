package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/google/uuid"
)

const (
	forkID           = 4
	waitForDBSeconds = 3
	vectorDir        = "./vectors/"
	genesisDir       = "./genesis/"
	executorURL      = "localhost:50071"
)

func main() {
	// Start containers
	defer func() {
		cmd := exec.Command("docker", "compose", "down", "--remove-orphans")
		if err := cmd.Run(); err != nil {
			log.Errorf("Failed stop containers: %v", err)
			return
		}
	}()
	log.Info("Starting DB and prover")
	cmd := exec.Command("docker", "compose", "up", "-d", "executor-tool-db")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("Failed to star DB: %w. %v", err, out)
		return
	}
	time.Sleep(time.Second * waitForDBSeconds)
	cmd = exec.Command("docker", "compose", "up", "-d", "executor-tool-prover")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("Failed to star prover: %v. %v", err, out)
		return
	}
	log.Info("DONE starting DB and prover")

	// Load vector file names
	files, err := os.ReadDir(vectorDir)
	if err != nil {
		log.Errorf("Error reading directory: %v", err)
		return
	}
	for _, file := range files {
		genesis, test, err := loadCase(vectorDir + file.Name())
		if err != nil {
			log.Errorf("Failed loading case: %v", err)
			return
		}
		if test.Skip {
			log.Infof("Case %v skipped", test.Title)
			continue
		}
		log.Infof("Running case %v\n", test.Title)
		err = runTestCase(context.Background(), genesis, test)
		if err != nil {
			log.Errorf("Failed running case: %v", err)
			return
		}
		log.Infof("Done running case %v\n\n\n\n\n", test.Title)
	}
}

func runTestCase(ctx context.Context, genesis []genesisItem, tc testCase) error {
	// DB connection
	dbConfig := db.Config{
		User:      testutils.GetEnv("PGUSER", "prover_user"),
		Password:  testutils.GetEnv("PGPASSWORD", "prover_pass"),
		Name:      testutils.GetEnv("PGDATABASE", "prover_db"),
		Host:      testutils.GetEnv("PGHOST", "localhost"),
		Port:      testutils.GetEnv("PGPORT", "5432"),
		EnableLog: false,
		MaxConns:  1,
	}
	sqlDB, err := db.NewSQLDB(dbConfig)
	if err != nil {
		return err
	}

	// Clean DB
	_, err = sqlDB.Exec(context.Background(), "DELETE FROM state.merkletree")
	if err != nil {
		return err
	}

	// Insert genesis
	for _, item := range genesis {
		_, err = sqlDB.Exec(
			context.Background(),
			"INSERT INTO state.merkletree (hash, data) VALUES ($1, $2)",
			item.Hash, item.Data,
		)
		if err != nil {
			return err
		}
	}

	// Executor connection
	executorClient, _, _ := executor.NewExecutorClient(ctx, executor.Config{URI: executorURL, MaxGRPCMessageSize: 100000000}) //nolint:gomnd
	// Execute batches
	for i := 0; i < len(tc.Requests); i++ {
		pbr := executor.ProcessBatchRequest(tc.Requests[i]) //nolint
		res, err := executorClient.ProcessBatch(ctx, &pbr)
		if err != nil {
			return err
		}
		log.Infof("**********              BATCH %d              **********", tc.Requests[i].OldBatchNum)
		txs, _, _, err := state.DecodeTxs(tc.Requests[i].BatchL2Data, forkID)
		if err != nil {
			log.Warnf("Txs are not correctly encoded")
		}
		lastTxWithResponse := 0
		log.Infof("CumulativeGasUsed: %v", res.CumulativeGasUsed)
		log.Infof("NewStateRoot: %v", hex.EncodeToString(res.NewStateRoot))
		log.Infof("NewLocalExitRoot: %v", hex.EncodeToString(res.NewLocalExitRoot))
		log.Infof("CntKeccakHashes: %v", res.CntKeccakHashes)
		log.Infof("CntPoseidonHashes: %v", res.CntPoseidonHashes)
		log.Infof("CntPoseidonPaddings: %v", res.CntPoseidonPaddings)
		log.Infof("CntMemAligns: %v", res.CntMemAligns)
		log.Infof("CntArithmetics: %v", res.CntArithmetics)
		log.Infof("CntBinaries: %v", res.CntBinaries)
		log.Infof("CntSteps: %v", res.CntSteps)
		for i, txRes := range res.Responses {
			log.Infof("=====> TX #%d", i)
			if "0x"+hex.EncodeToString(txRes.TxHash) != txs[i].Hash().Hex() {
				log.Warnf("TxHash missmatch:\nexecutor: %s\ndecoded: %s", "0x"+hex.EncodeToString(txRes.TxHash), txs[i].Hash().Hex())
			} else {
				log.Infof("= TxHash: %v", txs[i].Hash().Hex())
			}
			log.Infof("= Nonce: %d", txs[i].Nonce())
			log.Infof("= To: %v", txs[i].To())
			log.Infof("= Gas: %v", txs[i].Gas())
			log.Infof("= GasPrice: %v", txs[i].GasPrice())
			log.Infof("= StateRoot: %v", hex.EncodeToString(txRes.StateRoot))
			log.Infof("= Error: %v", txRes.Error)
			log.Infof("= GasUsed: %v", txRes.GasUsed)
			log.Infof("= GasLeft: %v", txRes.GasLeft)
			log.Infof("= GasRefunded: %v", txRes.GasRefunded)
			log.Infof("<======")
			lastTxWithResponse = i + 1
		}
		if lastTxWithResponse != len(txs) {
			log.Warnf("%d txs sent to the executor, but only got %d responses", len(txs), lastTxWithResponse)
			log.Info("Txs without response:")
			for i := lastTxWithResponse; i < len(txs); i++ {
				log.Info("--------------------------------")
				log.Infof("TxHash (decoded): %v", txs[i].Hash().Hex())
				log.Infof("Nonce: %d", txs[i].Nonce())
				log.Infof("To: %v", txs[i].To())
				log.Infof("Gas: %v", txs[i].Gas())
				log.Infof("GasPrice: %v", txs[i].GasPrice())
				log.Info("--------------------------------")
			}
		}
		log.Info("*******************************************************")
	}
	return nil
}

func loadCase(vectorFileName string) ([]genesisItem, testCase, error) {
	tc := testCase{}
	gen := []genesisItem{}
	// Load and parse test case
	f, err := os.ReadFile(vectorFileName)
	if err != nil {
		return gen, tc, err
	}
	err = json.Unmarshal(f, &tc)
	if err != nil {
		return gen, tc, err
	}
	// Load and parse genesis
	f, err = os.ReadFile(genesisDir + tc.GenesisFile)
	if err != nil {
		return gen, tc, err
	}
	err = json.Unmarshal(f, &gen)
	if err != nil {
		return gen, tc, err
	}
	return gen, tc, err
}

type genesisItem struct {
	Hash []byte
	Data []byte
}

func (gi *genesisItem) UnmarshalJSON(src []byte) error {
	type jGenesis struct {
		Hash string `json:"hash"`
		Data string `json:"data"`
	}
	jg := jGenesis{}
	if err := json.Unmarshal(src, &jg); err != nil {
		return err
	}
	hash, err := hex.DecodeString(jg.Hash)
	if err != nil {
		return err
	}
	data, err := hex.DecodeString(jg.Data)
	if err != nil {
		return err
	}
	*gi = genesisItem{
		Hash: hash,
		Data: data,
	}
	return nil
}

type testCase struct {
	Title       string            `json:"title"`
	GenesisFile string            `json:"genesisFile"`
	Skip        bool              `json:"skip"`
	Requests    []executorRequest `json:"batches"`
}

type executorRequest executor.ProcessBatchRequest

func (er *executorRequest) UnmarshalJSON(data []byte) error {
	type jExecutorRequest struct {
		BatchL2Data     string `json:"batchL2Data"`
		GlobalExitRoot  string `json:"globalExitRoot"`
		OldBatchNum     uint64 `json:"oldBatchNum"`
		OldAccInputHash string `json:"oldAccInputHash"`
		OldStateRoot    string `json:"oldStateRoot"`
		SequencerAddr   string `json:"sequencerAddr"`
		Timestamp       uint64 `json:"timestamp"`
	}
	jer := jExecutorRequest{}
	if err := json.Unmarshal(data, &jer); err != nil {
		return err
	}
	batchL2Data, err := hex.DecodeString(strings.TrimPrefix(jer.BatchL2Data, "0x"))
	if err != nil {
		return err
	}
	globalExitRoot, err := hex.DecodeString(strings.TrimPrefix(jer.GlobalExitRoot, "0x"))
	if err != nil {
		return err
	}
	oldAccInputHash, err := hex.DecodeString(strings.TrimPrefix(jer.OldAccInputHash, "0x"))
	if err != nil {
		return err
	}
	oldStateRoot, err := hex.DecodeString(strings.TrimPrefix(jer.OldStateRoot, "0x"))
	if err != nil {
		return err
	}

	req := executor.ProcessBatchRequest{
		BatchL2Data:     batchL2Data,
		GlobalExitRoot:  globalExitRoot,
		OldBatchNum:     jer.OldBatchNum,
		OldAccInputHash: oldAccInputHash,
		OldStateRoot:    oldStateRoot,
		Coinbase:        jer.SequencerAddr,
		EthTimestamp:    jer.Timestamp,
		ContextId:       uuid.NewString(),
	}
	*er = executorRequest(req) //nolint
	return nil
}
