package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
)

const (
	mockedStateRoot     = "0x090bcaf734c4f06c93954a827b45a6e8c67b8e0fd1e0a35a1c5982d6961828f9"
	mockedLocalExitRoot = "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e"

	ethTxManagerOwner = "aggregator"
	monitoredIDFormat = "proof-from-%v-to-%v"

	forkId9 = uint64(9)
)

type finalProofMsg struct {
	proverName     string
	proverID       string
	recursiveProof *state.Proof
	finalProof     *prover.FinalProof
}

// Aggregator represents an aggregator
type Aggregator struct {
	prover.UnimplementedAggregatorServiceServer

	cfg Config

	State                   stateInterface
	EthTxManager            ethTxManager
	Ethman                  etherman
	ProfitabilityChecker    aggregatorTxProfitabilityChecker
	TimeSendFinalProof      time.Time
	TimeCleanupLockedProofs types.Duration
	StateDBMutex            *sync.Mutex
	TimeSendFinalProofMutex *sync.RWMutex

	finalProof     chan finalProofMsg
	verifyingProof bool

	srv  *grpc.Server
	ctx  context.Context
	exit context.CancelFunc
}

// New creates a new aggregator.
func New(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
) (Aggregator, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	a := Aggregator{
		cfg: cfg,

		State:                   stateInterface,
		EthTxManager:            ethTxManager,
		Ethman:                  etherman,
		ProfitabilityChecker:    profitabilityChecker,
		StateDBMutex:            &sync.Mutex{},
		TimeSendFinalProofMutex: &sync.RWMutex{},
		TimeCleanupLockedProofs: cfg.CleanupLockedProofsInterval,

		finalProof: make(chan finalProofMsg),
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start(ctx context.Context) error {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	a.ctx = ctx
	a.exit = cancel

	metrics.Register()

	// process monitored batch verifications before starting
	a.EthTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
		a.handleMonitoredTxResult(result)
	}, nil)

	// Delete ungenerated recursive proofs
	err := a.State.DeleteUngeneratedBatchProofs(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to initialize proofs cache %w", err)
	}

	address := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	a.srv = grpc.NewServer()
	prover.RegisterAggregatorServiceServer(a.srv, a)

	healthService := newHealthChecker()
	grpchealth.RegisterHealthServer(a.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", a.cfg.Port)
		if err := a.srv.Serve(lis); err != nil {
			a.exit()
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	a.resetVerifyProofTime()

	go a.cleanupLockedProofs()
	go a.sendFinalProof()

	<-ctx.Done()
	return ctx.Err()
}

// Stop stops the Aggregator server.
func (a *Aggregator) Stop() {
	a.exit()
	a.srv.Stop()
}

// Channel implements the bi-directional communication channel between the
// Prover client and the Aggregator server.
func (a *Aggregator) Channel(stream prover.AggregatorService_ChannelServer) error {
	metrics.ConnectedProver()
	defer metrics.DisconnectedProver()

	ctx := stream.Context()
	var proverAddr net.Addr
	p, ok := peer.FromContext(ctx)
	if ok {
		proverAddr = p.Addr
	}
	prover, err := prover.New(stream, proverAddr, a.cfg.ProofStatePollingInterval)
	if err != nil {
		return err
	}

	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)
	log.Info("Establishing stream connection with prover")

	// Check if prover supports the required Fork ID
	if !prover.SupportsForkID(forkId9) {
		err := errors.New("prover does not support required fork ID")
		log.Warn(FirstToUpper(err.Error()))
		return err
	}

	for {
		select {
		case <-a.ctx.Done():
			// server disconnected
			return a.ctx.Err()
		case <-ctx.Done():
			// client disconnected
			return ctx.Err()

		default:
			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(a.cfg.RetryTime.Duration)
				continue
			}

			_, err = a.tryBuildFinalProof(ctx, prover, nil)
			if err != nil {
				log.Errorf("Error checking proofs to verify: %v", err)
			}

			proofGenerated, err := a.tryAggregateBlobOuterProofs(ctx, prover)
			if err != nil {
				log.Errorf("Error trying to aggregate blobOuter proofs: %v", err)
			}

			if !proofGenerated {
				proofGenerated, err = a.tryGenerateBlobOuterProof(ctx, prover)
				if err != nil {
					log.Errorf("Error trying to generate blobOuter proofs: %v", err)
				}
			}

			if !proofGenerated {
				proofGenerated, err = a.tryGenerateBlobInnerProof(ctx, prover)
				if err != nil {
					log.Errorf("Error trying to aggregate blobInner proofs: %v", err)
				}
			}

			if !proofGenerated {
				proofGenerated, err = a.tryAggregateBatchProofs(ctx, prover)
				if err != nil {
					log.Errorf("Error trying to aggregate batch proofs: %v", err)
				}
			}

			if !proofGenerated {
				proofGenerated, err = a.tryGenerateBatchProof(ctx, prover)
				if err != nil {
					log.Errorf("Error trying to generate batch proof: %v", err)
				}
			}
			if !proofGenerated {
				// if no proof was generated (aggregated or batch) wait some time before retry
				time.Sleep(a.cfg.RetryTime.Duration)
			} // if proof was generated we retry immediately as probably we have more proofs to process
		}
	}
}

// canVerifyProof returns true if we have reached the timeout to verify a proof
// and no other prover is verifying a proof (verifyingProof = false).
func (a *Aggregator) canVerifyProof() bool {
	a.TimeSendFinalProofMutex.RLock()
	defer a.TimeSendFinalProofMutex.RUnlock()
	return a.TimeSendFinalProof.Before(time.Now()) && !a.verifyingProof
}

// startProofVerification sets to true the verifyingProof variable to indicate that there is a proof verification in progress
func (a *Aggregator) startProofVerification() {
	a.TimeSendFinalProofMutex.Lock()
	defer a.TimeSendFinalProofMutex.Unlock()
	a.verifyingProof = true
}

// endProofVerification set verifyingProof to false to indicate that there is not proof verification in progress
func (a *Aggregator) endProofVerification() {
	a.TimeSendFinalProofMutex.Lock()
	defer a.TimeSendFinalProofMutex.Unlock()
	a.verifyingProof = false
}

// resetVerifyProofTime updates the timeout to verify a proof.
func (a *Aggregator) resetVerifyProofTime() {
	a.TimeSendFinalProofMutex.Lock()
	defer a.TimeSendFinalProofMutex.Unlock()
	a.TimeSendFinalProof = time.Now().Add(a.cfg.VerifyProofInterval.Duration)
}

// isSynced checks if the state is synchronized with L1. If a batch number is
// provided, it makes sure that the state is synced with that batch.
func (a *Aggregator) isSynced(ctx context.Context, batchNum *uint64) bool {
	// get latest verified batch as seen by the synchronizer
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err == state.ErrNotFound {
		return false
	}
	if err != nil {
		log.Warnf("Failed to get last consolidated batch: %v", err)
		return false
	}

	if lastVerifiedBatch == nil {
		return false
	}

	if batchNum != nil && lastVerifiedBatch.BatchNumber < *batchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, waiting for batch: %d", lastVerifiedBatch.BatchNumber, batchNum)
		return false
	}

	// latest verified batch in L1
	lastVerifiedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("Failed to get last eth batch, err: %v", err)
		return false
	}

	// check if L2 is synced with L1
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d, waiting for batch",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}

	return true
}

func (a *Aggregator) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*prover.InputProver, error) {
	previousBatch, err := a.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrNotFound {
		return nil, fmt.Errorf("failed to get previous batch, err: %v", err)
	}

	isForcedBatch := false
	batchRawData := &state.BatchRawV2{}
	if batchToVerify.BatchNumber == 1 || batchToVerify.ForcedBatchNum != nil || batchToVerify.BatchNumber == a.cfg.UpgradeEtrogBatchNumber {
		isForcedBatch = true
	} else {
		batchRawData, err = state.DecodeBatchV2(batchToVerify.BatchL2Data)
		if err != nil {
			log.Errorf("Failed to decode batch data, err: %v", err)
			return nil, err
		}
	}

	l1InfoTreeData := map[uint32]*prover.L1Data{}
	vb, err := a.State.GetVirtualBatch(ctx, batchToVerify.BatchNumber, nil)
	if err != nil {
		log.Errorf("Failed getting virtualBatch %d, err: %v", batchToVerify.BatchNumber, err)
		return nil, err
	}
	l1InfoRoot := vb.L1InfoRoot
	forcedBlockhashL1 := common.Hash{}

	if !isForcedBatch {
		tree, err := l1infotree.NewL1InfoTree(32, [][32]byte{}) // nolint:gomnd
		if err != nil {
			return nil, err
		}
		leaves, err := a.State.GetLeavesByL1InfoRoot(ctx, *l1InfoRoot, nil)
		if err != nil {
			return nil, err
		}

		aLeaves := make([][32]byte, len(leaves))
		for i, leaf := range leaves {
			aLeaves[i] = l1infotree.HashLeafData(leaf.GlobalExitRoot.GlobalExitRoot, leaf.PreviousBlockHash, uint64(leaf.Timestamp.Unix()))
		}

		for _, l2blockRaw := range batchRawData.Blocks {
			_, contained := l1InfoTreeData[l2blockRaw.IndexL1InfoTree]
			if !contained && l2blockRaw.IndexL1InfoTree != 0 {
				l1InfoTreeExitRootStorageEntry := state.L1InfoTreeExitRootStorageEntry{}
				l1InfoTreeExitRootStorageEntry.Timestamp = time.Unix(0, 0)
				if l2blockRaw.IndexL1InfoTree <= leaves[len(leaves)-1].L1InfoTreeIndex {
					l1InfoTreeExitRootStorageEntry, err = a.State.GetL1InfoRootLeafByIndex(ctx, l2blockRaw.IndexL1InfoTree, nil)
					if err != nil {
						return nil, err
					}
				}

				// Calculate smt proof
				smtProof, calculatedL1InfoRoot, err := tree.ComputeMerkleProof(l2blockRaw.IndexL1InfoTree, aLeaves)
				if err != nil {
					return nil, err
				}
				if l1InfoRoot != nil && *l1InfoRoot != calculatedL1InfoRoot {
					for i, l := range aLeaves {
						log.Infof("AllLeaves[%d]: %s", i, common.Bytes2Hex(l[:]))
					}
					for i, s := range smtProof {
						log.Infof("smtProof[%d]: %s", i, common.Bytes2Hex(s[:]))
					}
					return nil, fmt.Errorf("error: l1InfoRoot mismatch. L1InfoRoot: %s, calculatedL1InfoRoot: %s. l1InfoTreeIndex: %d", l1InfoRoot.String(), calculatedL1InfoRoot.String(), l2blockRaw.IndexL1InfoTree)
				}

				protoProof := make([][]byte, len(smtProof))
				for i, proof := range smtProof {
					tmpProof := proof
					protoProof[i] = tmpProof[:]
				}

				l1InfoTreeData[l2blockRaw.IndexL1InfoTree] = &prover.L1Data{
					GlobalExitRoot: l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.GlobalExitRoot.GlobalExitRoot.Bytes(),
					BlockhashL1:    l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.PreviousBlockHash.Bytes(),
					MinTimestamp:   uint32(l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.GlobalExitRoot.Timestamp.Unix()),
					SmtProof:       protoProof,
				}
			}
		}
	} else {
		// Initial batch must be handled differently
		if batchToVerify.BatchNumber == 1 || batchToVerify.BatchNumber == a.cfg.UpgradeEtrogBatchNumber {
			forcedBlockhashL1, err = a.State.GetVirtualBatchParentHash(ctx, batchToVerify.BatchNumber, nil)
			if err != nil {
				return nil, err
			}
		} else {
			forcedBlockhashL1, err = a.State.GetForcedBatchParentHash(ctx, *batchToVerify.ForcedBatchNum, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	inputProver := &prover.InputProver{
		PublicInputs: &prover.PublicInputs{
			OldStateRoot:      previousBatch.StateRoot.Bytes(),
			OldAccInputHash:   previousBatch.AccInputHash.Bytes(),
			OldBatchNum:       previousBatch.BatchNumber,
			ChainId:           a.cfg.ChainID,
			ForkId:            forkId9,
			BatchL2Data:       batchToVerify.BatchL2Data,
			L1InfoRoot:        l1InfoRoot.Bytes(),
			TimestampLimit:    uint64(batchToVerify.Timestamp.Unix()),
			SequencerAddr:     batchToVerify.Coinbase.String(),
			AggregatorAddr:    a.cfg.SenderAddress,
			L1InfoTreeData:    l1InfoTreeData,
			ForcedBlockhashL1: forcedBlockhashL1.Bytes(),
		},
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
	}

	printInputProver(inputProver)

	return inputProver, nil
}

func printInputProver(inputProver *prover.InputProver) {
	log.Debugf("OldStateRoot: %v", common.BytesToHash(inputProver.PublicInputs.OldStateRoot))
	log.Debugf("OldAccInputHash: %v", common.BytesToHash(inputProver.PublicInputs.OldAccInputHash))
	log.Debugf("OldBatchNum: %v", inputProver.PublicInputs.OldBatchNum)
	log.Debugf("ChainId: %v", inputProver.PublicInputs.ChainId)
	log.Debugf("ForkId: %v", inputProver.PublicInputs.ForkId)
	log.Debugf("BatchL2Data: %v", common.Bytes2Hex(inputProver.PublicInputs.BatchL2Data))
	log.Debugf("L1InfoRoot: %v", common.BytesToHash(inputProver.PublicInputs.L1InfoRoot))
	log.Debugf("TimestampLimit: %v", inputProver.PublicInputs.TimestampLimit)
	log.Debugf("SequencerAddr: %v", inputProver.PublicInputs.SequencerAddr)
	log.Debugf("AggregatorAddr: %v", inputProver.PublicInputs.AggregatorAddr)
	log.Debugf("L1InfoTreeData: %+v", inputProver.PublicInputs.L1InfoTreeData)
	log.Debugf("ForcedBlockhashL1: %v", common.Bytes2Hex(inputProver.PublicInputs.ForcedBlockhashL1))
}

// healthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// newHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Check(ctx context.Context, req *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Watch(req *grpchealth.HealthCheckRequest, server grpchealth.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	})
}

func (a *Aggregator) handleMonitoredTxResult(result ethtxmanager.MonitoredTxResult) {
	mTxResultLogger := ethtxmanager.CreateMonitoredTxResultLogger(ethTxManagerOwner, result)
	if result.Status == ethtxmanager.MonitoredTxStatusFailed {
		mTxResultLogger.Fatal("failed to send batch verification, TODO: review this fatal and define what to do in this case")
	}

	// monitoredIDFormat: "proof-from-%v-to-%v"
	idSlice := strings.Split(result.ID, "-")
	proofBatchNumberStr := idSlice[2]
	proofBatchNumber, err := strconv.ParseUint(proofBatchNumberStr, encoding.Base10, 0)
	if err != nil {
		mTxResultLogger.Errorf("failed to read final proof batch number from monitored tx: %v", err)
	}

	proofBatchNumberFinalStr := idSlice[4]
	proofBatchNumberFinal, err := strconv.ParseUint(proofBatchNumberFinalStr, encoding.Base10, 0)
	if err != nil {
		mTxResultLogger.Errorf("failed to read final proof batch number final from monitored tx: %v", err)
	}

	log := log.WithFields("txId", result.ID, "batches", fmt.Sprintf("%d-%d", proofBatchNumber, proofBatchNumberFinal))
	log.Info("Final proof verified")

	// wait for the synchronizer to catch up the verified batches
	log.Debug("A final proof has been sent, waiting for the network to be synced")
	for !a.isSynced(a.ctx, &proofBatchNumberFinal) {
		log.Info("Waiting for synchronizer to sync...")
		time.Sleep(a.cfg.RetryTime.Duration)
	}

	// network is synced with the final proof, we can safely delete all recursive
	// proofs up to the last synced batch
	err = a.State.CleanupBatchProofs(a.ctx, proofBatchNumberFinal, nil)
	if err != nil {
		log.Errorf("Failed to store proof aggregation result: %v", err)
	}
}

func buildMonitoredTxID(batchNumber, batchNumberFinal uint64) string {
	return fmt.Sprintf(monitoredIDFormat, batchNumber, batchNumberFinal)
}

func (a *Aggregator) cleanupLockedProofs() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-time.After(a.TimeCleanupLockedProofs.Duration):
			n, err := a.State.CleanupLockedBatchProofs(a.ctx, a.cfg.GeneratingProofCleanupThreshold, nil)
			if err != nil {
				log.Errorf("Failed to cleanup locked proofs: %v", err)
			}
			if n == 1 {
				log.Warn("Found a stale proof and removed form cache")
			} else if n > 1 {
				log.Warnf("Found %d stale proofs and removed from cache", n)
			}
		}
	}
}

// FirstToUpper returns the string passed as argument with the first letter in
// uppercase.
func FirstToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
