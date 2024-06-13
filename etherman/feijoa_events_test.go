package etherman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	_, err := NewFeijoaContracts(nil, L1Config{})
	require.NoError(t, err)
}

func TestFeijoaEventsSignature(t *testing.T) {
	// Signature extracted from https://sepolia.etherscan.io/tx/0x644699c839d34a61c531d7ecf12390bf38c06a62715ca4edce978b9213ce3cd1#eventlog
	require.Equal(t, "0x470f4ca4b003755c839b80ab00c3efbeb69d6eafec00e1a3677482933ec1fd0c", eventSequenceBlobsSignatureHash.String())
}
