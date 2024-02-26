package sequencesender

import (
	"context"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
)

type dataAbilitier interface {
	PostSequence(ctx context.Context, sequences []ethmanTypes.Sequence) ([]byte, error)
}
