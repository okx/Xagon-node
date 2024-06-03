package test

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

const (
	Ether155V = 27
)

var (
	stateTree                          *merkletree.StateTree
	stateDb                            *pgxpool.Pool
	err                                error
	StateDBCfg                         = dbutils.NewStateConfigFromEnv()
	ctx                                = context.Background()
	ExecutorClient                     executor.ExecutorServiceClient
	mtDBServiceClient                  hashdb.HashDBServiceClient
	executorClientConn, mtDBClientConn *grpc.ClientConn
	executorCancel, mtDBCancel         context.CancelFunc
	Genesis                            = state.Genesis{}
)

func CloseTestState() {
	stateDb.Close()
	executorCancel()
	executorClientConn.Close()
	mtDBCancel()
	mtDBClientConn.Close()
}

func InitTestState(stateCfg state.Config) *state.State {
	InitOrResetDB(StateDBCfg)

	stateDb, err = db.NewSQLDB(StateDBCfg)
	if err != nil {
		panic(err)
	}

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}
	ExecutorClient, executorClientConn, executorCancel = executor.NewExecutorClient(ctx, executorServerConfig)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())

	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s = mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())

	stateTree = merkletree.NewStateTree(mtDBServiceClient)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	return state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), ExecutorClient, stateTree, eventLog, mt)
}

func InitOrResetDB(cfg db.Config) {
	if err := dbutils.InitOrResetState(cfg); err != nil {
		panic(err)
	}
}
