package aggregator

import "context"

func (a *Aggregator) tryGenerateBlobOuterProof(ctx context.Context, prover proverInterface) (bool, error) {
	return false, nil
}

func (a *Aggregator) tryAggregateBlobOuterProofs(ctx context.Context, prover proverInterface) (bool, error) {
	return false, nil
}
