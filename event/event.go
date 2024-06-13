package event

import (
	"math/big"
	"time"
)

// EventID is the ID of the event
type EventID string

// Source is the source of the event
type Source string

// Component is the component that triggered the event
type Component string

// Level is the level of the event
type Level string

const (
	// EventID_NodeComponentStarted is triggered when the node starts
	EventID_NodeComponentStarted = "NODE COMPONENT STARTED"
	// EventID_PreexecutionOOC is triggered when an OOC error is detected during the preexecution
	EventID_PreexecutionOOC EventID = "PRE EXECUTION OOC"
	// EventID_PreexecutionOOG is triggered when an OOG error is detected during the preexecution
	EventID_PreexecutionOOG EventID = "PRE EXECUTION OOG"
	// EventID_ExecutorError is triggered when an error is detected during the execution
	EventID_ExecutorError EventID = "EXECUTOR ERROR"
	// EventID_ReprocessFullBatchOOC is triggered when an OOC error is detected during the reprocessing of a full batch
	EventID_ReprocessFullBatchOOC EventID = "REPROCESS FULL BATCH OOC"
	// EventID_ExecutorRLPError is triggered when an RLP error is detected during the execution
	EventID_ExecutorRLPError EventID = "EXECUTOR RLP ERROR"
	// EventID_FinalizerHalt is triggered when the finalizer halts
	EventID_FinalizerHalt EventID = "FINALIZER HALT"
	// EventID_FinalizerRestart is triggered when the finalizer restarts
	EventID_FinalizerRestart EventID = "FINALIZER RESTART"
	// EventID_FinalizerBreakEvenGasPriceBigDifference is triggered when the finalizer recalculates the break even gas price and detects a big difference
	EventID_FinalizerBreakEvenGasPriceBigDifference EventID = "FINALIZER BREAK EVEN GAS PRICE BIG DIFFERENCE"
	// EventID_SynchronizerRestart is triggered when the Synchonizer restarts
	EventID_SynchronizerRestart EventID = "SYNCHRONIZER RESTART"
	// EventID_SynchronizerHalt is triggered when the synchronizer halts
	EventID_SynchronizerHalt EventID = "SYNCHRONIZER HALT"
	// EventID_SequenceSenderHalt is triggered when the SequenceSender halts
	EventID_SequenceSenderHalt EventID = "SEQUENCESENDER HALT"
	// EventID_NodeOOC is triggered when an OOC at node level is detected
	EventID_NodeOOC EventID = "NODE OOC"
	// EventID_UsedZKCountersOverflow is triggered when used ZK counters exceeds remaining batch ZK counters
	EventID_UsedZKCountersOverflow EventID = "USED ZKCOUNTERS OVERFLOW"
	// EventID_ReservedZKCountersOverflow is triggered when reserved ZK counters exceeds remaining batch ZK counters
	EventID_ReservedZKCountersOverflow EventID = "RESERVED ZKCOUNTERS OVERFLOW"
	// EventID_InvalidInfoRoot is triggered when an invalid l1InfoRoot was synced
	EventID_InvalidInfoRoot EventID = "INVALID INFOROOT"
	// EventID_L2BlockReorg is triggered when a L2 block reorg has happened in the sequencer
	EventID_L2BlockReorg EventID = "L2 BLOCK REORG"
	// Source_Node is the source of the event
	Source_Node Source = "node"

	// Component_RPC is the component that triggered the event
	Component_RPC Component = "rpc"
	// Component_Pool is the component that triggered the event
	Component_Pool Component = "pool"
	// Component_Sequencer is the component that triggered the event
	Component_Sequencer Component = "sequencer"
	// Component_Synchronizer is the component that triggered the event
	Component_Synchronizer Component = "synchronizer"
	// Component_Aggregator is the component that triggered the event
	Component_Aggregator Component = "aggregator"
	// Component_EthTxManager is the component that triggered the event
	Component_EthTxManager Component = "ethtxmanager"
	// Component_GasPricer is the component that triggered the event
	Component_GasPricer Component = "gaspricer"
	// Component_Executor is the component that triggered the event
	Component_Executor Component = "executor"
	// Component_Broadcast is the component that triggered the event
	Component_Broadcast Component = "broadcast"
	// Component_Sequence_Sender is the component that triggered the event
	Component_Sequence_Sender = "seqsender"

	// Level_Emergency is the most severe level
	Level_Emergency Level = "emerg"
	// Level_Alert is the second most severe level
	Level_Alert Level = "alert"
	// Level_Critical is the third most severe level
	Level_Critical Level = "crit"
	// Level_Error is the fourth most severe level
	Level_Error Level = "err"
	// Level_Warning is the fifth most severe level
	Level_Warning Level = "warning"
	// Level_Notice is the sixth most severe level
	Level_Notice Level = "notice"
	// Level_Info is the seventh most severe level
	Level_Info Level = "info"
	// Level_Debug is the least severe level
	Level_Debug Level = "debug"
)

// Event represents a event that may be investigated
type Event struct {
	Id          big.Int
	ReceivedAt  time.Time
	IPAddress   string
	Source      Source
	Component   Component
	Level       Level
	EventID     EventID
	Description string
	Data        []byte
	Json        interface{}
}
