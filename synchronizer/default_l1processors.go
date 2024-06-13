package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/elderberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/feijoa"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer, l2Blockchecker *actions.CheckL2BlockHash) *processor_manager.L1EventProcessors {
	p := processor_manager.NewL1EventProcessorsBuilder()
	p.Register(incaberry.NewProcessorL1GlobalExitRoot(sync.state))
	p.Register(actions.NewCheckL2BlockDecorator(incaberry.NewProcessorL1SequenceBatches(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync), l2Blockchecker))
	p.Register(actions.NewCheckL2BlockDecorator(incaberry.NewProcessL1ForcedBatches(sync.state), l2Blockchecker))
	p.Register(actions.NewCheckL2BlockDecorator(incaberry.NewProcessL1SequenceForcedBatches(sync.state, sync), l2Blockchecker))
	p.Register(incaberry.NewProcessorForkId(sync.state, sync))
	p.Register(etrog.NewProcessorL1InfoTreeUpdate(sync.state))
	p.Register(feijoa.NewProcessorL1InfoTreeUpdate(sync.state))
	sequenceBatchesProcessor := etrog.NewProcessorL1SequenceBatches(sync.state, sync, common.DefaultTimeProvider{}, sync.halter)
	p.Register(actions.NewCheckL2BlockDecorator(sequenceBatchesProcessor, l2Blockchecker))
	p.Register(incaberry.NewProcessorL1VerifyBatch(sync.state))
	p.Register(etrog.NewProcessorL1UpdateEtrogSequence(sync.state, sync, common.DefaultTimeProvider{}))
	p.Register(actions.NewCheckL2BlockDecorator(elderberry.NewProcessorL1SequenceBatchesElderberry(sequenceBatchesProcessor, sync.state), l2Blockchecker))
	// intialSequence is process in ETROG by the same class, this is just a wrapper to pass directly to ETROG
	p.Register(elderberry.NewProcessorL1InitialSequenceBatchesElderberry(sequenceBatchesProcessor))
	p.Register(feijoa.NewProcessorSequenceBlobs(sync.state, sync.state, nil))
	return p.Build()
}
