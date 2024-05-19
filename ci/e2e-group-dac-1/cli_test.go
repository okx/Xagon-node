package e2e_group_dac_1

import (
	"context"
	"regexp"
	"testing"

	polygondatacommittee "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygondatacommittee_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestSetDataAvailabilityProtocol(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	err := operations.Teardown()
	require.NoError(t, err)

	opsCfg := operations.GetDefaultOperationsConfig()

	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)

	err = opsman.Setup()
	require.NoError(t, err)

	clientL1, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	require.NoError(t, err)

	auth, err := operations.GetAuth("0x815405dddb0e2a99b12af775fd2929e526704e1d1aea6a0b4e74dc33e2f7fcd2", operations.DefaultL1ChainID)
	require.NoError(t, err)

	// New DAC Setup
	_, tx, newDA, err := polygondatacommittee.DeployPolygondatacommitteeXlayer(auth, clientL1)
	require.NoError(t, err)
	require.NoError(t, operations.WaitTxToBeMined(ctx, clientL1, tx, operations.DefaultTimeoutTxToBeMined))

	tx, err = newDA.Initialize(auth)
	require.NoError(t, err)
	require.NoError(t, operations.WaitTxToBeMined(ctx, clientL1, tx, operations.DefaultTimeoutTxToBeMined))
}

func extractHexFromString(output string) string {
	re := regexp.MustCompile(`Transaction to set new data availability protocol sent. Hash: (0x[0-9a-fA-F]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
