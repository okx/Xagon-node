package etherman

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/datacommittee"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/mockverifier"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmbridge"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func NewSimulatedEtherman(cfg Config, auth *bind.TransactOpts) (etherman *Client, ethBackend *backends.SimulatedBackend, maticAddr common.Address, br *polygonzkevmbridge.Polygonzkevmbridge, da *datacommittee.Datacommittee, err error) {
	if auth == nil {
		// read only client
		return &Client{}, nil, common.Address{}, nil, nil, nil
	}
	// 10000000 ETH in wei
	balance, _ := new(big.Int).SetString("10000000000000000000000000", 10) //nolint:gomnd
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(999999999999999999) //nolint:gomnd
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// DAC Setup
	dataCommitteeAddr, _, da, err := datacommittee.DeployDatacommittee(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	_, err = da.Initialize(auth)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	_, err = da.SetupCommittee(auth, big.NewInt(0), []string{}, []byte{})
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}

	// Deploy contracts
	const maticDecimalPlaces = 18
	totalSupply, _ := new(big.Int).SetString("10000000000000000000000000000", 10) //nolint:gomnd
	maticAddr, _, maticContract, err := matic.DeployMatic(auth, client, "Matic Token", "MATIC", maticDecimalPlaces, totalSupply)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	rollupVerifierAddr, _, _, err := mockverifier.DeployMockverifier(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	nonce, err := client.PendingNonceAt(context.TODO(), auth.From)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	const posBridge = 1
	calculatedBridgeAddr := crypto.CreateAddress(auth.From, nonce+posBridge)
	const posPoE = 2
	calculatedPoEAddr := crypto.CreateAddress(auth.From, nonce+posPoE)
	genesis := common.HexToHash("0xfd3434cd8f67e59d73488a2b8da242dd1f02849ea5dd99f0ca22c836c3d5b4a9") // Random value. Needs to be different to 0x0
	exitManagerAddr, _, globalExitRoot, err := polygonzkevmglobalexitroot.DeployPolygonzkevmglobalexitroot(auth, client, calculatedPoEAddr, calculatedBridgeAddr)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	bridgeAddr, _, br, err := polygonzkevmbridge.DeployPolygonzkevmbridge(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	poeAddr, _, poe, err := polygonzkevm.DeployPolygonzkevm(auth, client, exitManagerAddr, maticAddr, rollupVerifierAddr, bridgeAddr, dataCommitteeAddr, 1000, 1) //nolint
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	_, err = br.Initialize(auth, 0, exitManagerAddr, poeAddr)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}

	poeParams := polygonzkevm.PolygonZkEVMInitializePackedParameters{
		Admin:                    auth.From,
		TrustedSequencer:         auth.From,
		PendingStateTimeout:      10000, //nolint:gomnd
		TrustedAggregator:        auth.From,
		TrustedAggregatorTimeout: 10000, //nolint:gomnd
	}
	_, err = poe.Initialize(auth, poeParams, genesis, "http://localhost", "L2", "v1") //nolint:gomnd
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}

	if calculatedBridgeAddr != bridgeAddr {
		return nil, nil, common.Address{}, nil, nil, fmt.Errorf("bridgeAddr (%s) is different from the expected contract address (%s)",
			bridgeAddr.String(), calculatedBridgeAddr.String())
	}
	if calculatedPoEAddr != poeAddr {
		return nil, nil, common.Address{}, nil, nil, fmt.Errorf("poeAddr (%s) is different from the expected contract address (%s)",
			poeAddr.String(), calculatedPoEAddr.String())
	}

	// Approve the bridge and poe to spend 10000 matic tokens.
	approvedAmount, _ := new(big.Int).SetString("10000000000000000000000", 10) //nolint:gomnd
	_, err = maticContract.Approve(auth, bridgeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	_, err = maticContract.Approve(auth, poeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	_, err = poe.ActivateForceBatches(auth)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}

	client.Commit()
	c := &Client{
		EthClient:             client,
		ZkEVM:                 poe,
		Matic:                 maticContract,
		GlobalExitRootManager: globalExitRoot,
		DataCommittee:         da,
		SCAddresses:           []common.Address{poeAddr, exitManagerAddr, dataCommitteeAddr},
		auth:                  map[common.Address]bind.TransactOpts{},
		cfg:                   cfg,
	}
	err = c.AddOrReplaceAuth(*auth)
	if err != nil {
		return nil, nil, common.Address{}, nil, nil, err
	}
	return c, client, maticAddr, br, da, nil
}
