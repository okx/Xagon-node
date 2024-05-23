package apollo

const (
	// RPC is the rpc unit,the content of the unit is the config for rpc with toml format
	RPC = "rpc"
	// State is the state unit,the content of the unit is the config for state with toml format
	State = "state"
	// Pool is the pool unit,the content of the unit is the config for pool with toml format
	Pool = "pool"
	// Etherman is the etherman unit,the content of the unit is the config for etherman with toml format
	Etherman = "etherman"
	// Synchronizer is the syncer unit,the content of the unit is the config for syncer with toml format
	Synchronizer = "synchronizer"
	// Sequencer is the sequencer unit,the content of the unit is the config for sequencer with toml format
	Sequencer = "sequencer"
	// SequenceSender is the sequence sender unit,the content of the unit is the config for sequence sender with toml format
	SequenceSender = "sequence_sender"
	// Aggregator is the aggregator unit,the content of the unit is the config for aggregator with toml format
	Aggregator = "aggregator"
	// EthTxManager is the eth tx manager unit,the content of the unit is the config for eth tx manager with toml format
	EthTxManager = "eth_tx_manager"
	// L2GasPriceSuggester is the l2 gas price suggester unit,the content of the unit is the config for l2 gas price suggester with toml format
	L2GasPriceSuggester = "l2_gas_price_suggester"
)
