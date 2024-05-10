package l2_sync_etrog

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// StateGERInteface interface that requires from State
type StateGERInteface interface {
	GetExitRootByGlobalExitRoot(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (*state.GlobalExitRoot, error)
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
}

// CheckSyncStatusToProcessBatch Implements CheckSyncStatusToProcessBatchInterface
type CheckSyncStatusToProcessBatch struct {
	zkevmRPCClient syncinterfaces.ZKEVMClientGlobalExitRootGetter
	state          StateGERInteface
}

// NewCheckSyncStatusToProcessBatch returns a new instance of CheckSyncStatusToProcessBatch
func NewCheckSyncStatusToProcessBatch(zkevmRPCClient syncinterfaces.ZKEVMClientGlobalExitRootGetter, state StateGERInteface) *CheckSyncStatusToProcessBatch {
	return &CheckSyncStatusToProcessBatch{
		zkevmRPCClient: zkevmRPCClient,
		state:          state,
	}
}

// CheckL1SyncGlobalExitRootEnoughToProcessBatch returns nil if the are sync and could process the batch
// if not:
//   - returns syncinterfaces.ErrMissingSyncFromL1 if we are behind the block number that contains the GlobalExitRoot
//   - returns l2_shared.NewDeSyncPermissionlessAndTrustedNodeError  if trusted and and permissionless are not in same page! pass also the discrepance point
func (c *CheckSyncStatusToProcessBatch) CheckL1SyncGlobalExitRootEnoughToProcessBatch(ctx context.Context, batchNumber uint64, globalExitRoot common.Hash, dbTx pgx.Tx) error {
	// Find out if this node have GlobalExitRoot
	// If not: ask to zkevm-RPC the block number of this GlobalExitRoot
	// If we are behind this block number returns ErrMissingSyncFromL1
	// If not we have a problem!
	if globalExitRoot == state.ZeroHash {
		// Special case that batch doesnt use any GlobalExitRoot
		return nil
	}
	debugStr := fmt.Sprintf("CheckL1SyncStatusEnoughToProcessBatch batchNumber:%d globalExitRoot: %s ", batchNumber, globalExitRoot.Hex())
	localGERInfo, err := c.state.GetExitRootByGlobalExitRoot(ctx, globalExitRoot, dbTx)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Errorf("error getting GetExitRootByGlobalExitRoot %s . Error: ", debugStr, err)
		return err
	}
	if err == nil {
		// We have this GlobalExitRoot, so we are synced from L1
		log.Infof("We have this GlobalExitRoot (%s) in L1block %d, so we are synced from L1 %s", globalExitRoot.String(), localGERInfo.BlockNumber, debugStr)
		return nil
	}
	// this means err != state.ErrNotFound -> so we have to ask to zkevm-RPC the block number of this GlobalExitRoot
	exitRoots, err := c.zkevmRPCClient.ExitRootsByGER(ctx, globalExitRoot)
	if err != nil || exitRoots == nil {
		log.Errorf("error getting blockNumber  from RPC. %s Error: ", debugStr, err)
		return err
	}
	// We have the L1 BlockNumber that contains this GlobalExitRoot check if we are behind
	lastL1BlockSynced, err := c.state.GetLastBlock(ctx, dbTx)
	if err != nil {
		log.Errorf("error getting last block from state. %s Error: ", debugStr, err)
		return err
	}
	if uint64(exitRoots.BlockNumber) > lastL1BlockSynced.BlockNumber {
		log.Warnf("We are behind this block number. GER=%s sync in block %d and we are in block %d %s", globalExitRoot, exitRoots.BlockNumber, lastL1BlockSynced.BlockNumber, debugStr)
		return syncinterfaces.ErrMissingSyncFromL1
	}
	// ??!?! We are desynced from L1!
	log.Errorf("We are desynced from L1! GER=%s sync in block %d and we are in block %d but dont have this GER!!! %s", globalExitRoot, exitRoots.BlockNumber, lastL1BlockSynced.BlockNumber, debugStr)
	return l2_shared.NewDeSyncPermissionlessAndTrustedNodeError(uint64(exitRoots.BlockNumber))
}
