package client

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/stretchr/testify/require"
)

func TestZkevmGetBatch(t *testing.T) {
	t.Skip("This test is exploratory")
	// Create a new client
	client := NewClient("https://zkevm-rpc.com/")
	lastTrustedStateBatchNumberSeen, err := client.BatchNumber(context.Background())
	require.NoError(t, err)
	log.Info("lastTrustedStateBatchNumberSeen: ", lastTrustedStateBatchNumberSeen)
	batch, err := client.BatchByNumber(context.Background(), big.NewInt(int64(lastTrustedStateBatchNumberSeen)))
	require.NoError(t, err)

	// Print the batch
	fmt.Println(batch)
}
