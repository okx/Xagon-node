package aggregator

// SettlementBackend is the type of the settlement backend
type SettlementBackend string

const (
	// AggLayer settlement backend
	AggLayer SettlementBackend = "agglayer"

	// L1 settlement backend
	L1 SettlementBackend = "l1"
)
