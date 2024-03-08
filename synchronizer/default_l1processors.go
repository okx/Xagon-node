package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/elderberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer, zkEVMClient syncinterfaces.ZKEVMClientInterface) *processor_manager.L1EventProcessors {
	p := processor_manager.NewL1EventProcessorsBuilder()
	p.Register(incaberry.NewProcessorL1GlobalExitRoot(sync.state))
	p.Register(incaberry.NewProcessorL1SequenceBatches(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync, zkEVMClient))
	p.Register(incaberry.NewProcessL1ForcedBatches(sync.state))
	p.Register(incaberry.NewProcessL1SequenceForcedBatches(sync.state, sync))
	p.Register(incaberry.NewProcessorForkId(sync.state, sync))
	p.Register(etrog.NewProcessorL1InfoTreeUpdate(sync.state))
	sequenceBatchesProcessor := etrog.NewProcessorL1SequenceBatches(sync.state, sync, common.DefaultTimeProvider{}, sync.halter)
	p.Register(sequenceBatchesProcessor)
	p.Register(incaberry.NewProcessorL1VerifyBatch(sync.state))
	p.Register(etrog.NewProcessorL1UpdateEtrogSequence(sync.state, sync, common.DefaultTimeProvider{}))
	p.Register(elderberry.NewProcessorL1SequenceBatchesElderberry(sequenceBatchesProcessor, sync.state))
	// intialSequence is process in ETROG by the same class, this is just a wrapper to pass directly to ETROG
	p.Register(elderberry.NewProcessorL1InitialSequenceBatchesElderberry(sequenceBatchesProcessor))
	return p.Build()
}
