package datacommittee

import (
	"math/big"
	"testing"

	polygondatacommittee "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygondatacommittee_x1"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateDataCommitteeEvent(t *testing.T) {
	// Set up testing environment
	dac, ethBackend, auth, da := newTestingEnv(t)

	// Update the committee
	requiredAmountOfSignatures := big.NewInt(2)
	URLs := []string{"1", "2", "3"}
	addrs := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
	}
	addrsBytes := []byte{}
	for _, addr := range addrs {
		addrsBytes = append(addrsBytes, addr.Bytes()...)
	}
	_, err := da.SetupCommittee(auth, requiredAmountOfSignatures, URLs, addrsBytes)
	require.NoError(t, err)
	ethBackend.Commit()

	// Assert the committee update
	actualSetup, err := dac.getCurrentDataCommittee()
	require.NoError(t, err)
	expectedMembers := []DataCommitteeMember{}
	expectedSetup := DataCommittee{
		RequiredSignatures: uint64(len(URLs) - 1),
		AddressesHash:      crypto.Keccak256Hash(addrsBytes),
	}
	for i, url := range URLs {
		expectedMembers = append(expectedMembers, DataCommitteeMember{
			URL:  url,
			Addr: addrs[i],
		})
	}
	expectedSetup.Members = expectedMembers
	assert.Equal(t, expectedSetup, *actualSetup)
}

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv(t *testing.T) (
	dac *DataCommitteeBackend,
	ethBackend *backends.SimulatedBackend,
	auth *bind.TransactOpts,
	da *polygondatacommittee.Polygondatacommittee,
) {
	t.Helper()
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	dac, ethBackend, da, err = newSimulatedDacman(t, auth)
	if err != nil {
		log.Fatal(err)
	}
	return dac, ethBackend, auth, da
}

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func newSimulatedDacman(t *testing.T, auth *bind.TransactOpts) (
	dacman *DataCommitteeBackend,
	ethBackend *backends.SimulatedBackend,
	da *polygondatacommittee.Polygondatacommittee,
	err error,
) {
	t.Helper()
	if auth == nil {
		// read only client
		return &DataCommitteeBackend{}, nil, nil, nil
	}
	// 10000000 ETH in wei
	balance, _ := new(big.Int).SetString("10000000000000000000000000", 10) //nolint:gomnd
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	// blockGasLimit := uint64(999999999999999999) //nolint:gomnd
	client := backends.NewSimulatedBackend(genesisAlloc, uint64(999999999999999999)) //nolint:staticcheck,gomnd

	// DAC Setup
	_, _, da, err = polygondatacommittee.DeployPolygondatacommittee(auth, client)
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	_, err = da.Initialize(auth)
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	_, err = da.SetupCommittee(auth, big.NewInt(0), []string{}, []byte{})
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	client.Commit()

	c := &DataCommitteeBackend{
		dataCommitteeContract: da,
	}
	return c, client, da, nil
}
