package etherman

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	beaconclient "github.com/0xPolygonHermez/zkevm-node/beacon_client"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/eip4844"
	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
	"github.com/0xPolygonHermez/zkevm-node/etherman/ethgasstation"
	"github.com/0xPolygonHermez/zkevm-node/etherman/metrics"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/elderberrypolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/etrogpolygonzkevmglobalexitroot"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/pol"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/preetrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/preetrogpolygonzkevmglobalexitroot"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

const (
	// ETrogUpgradeVersion is the version of the LxLy upgrade
	ETrogUpgradeVersion = 2
)

var (
	// Events EtrogRollupManager
	setBatchFeeSignatureHash                       = crypto.Keccak256Hash([]byte("SetBatchFee(uint256)"))
	setTrustedAggregatorSignatureHash              = crypto.Keccak256Hash([]byte("SetTrustedAggregator(address)"))       // Used in oldZkEvm as well
	setVerifyBatchTimeTargetSignatureHash          = crypto.Keccak256Hash([]byte("SetVerifyBatchTimeTarget(uint64)"))    // Used in oldZkEvm as well
	setMultiplierBatchFeeSignatureHash             = crypto.Keccak256Hash([]byte("SetMultiplierBatchFee(uint16)"))       // Used in oldZkEvm as well
	setPendingStateTimeoutSignatureHash            = crypto.Keccak256Hash([]byte("SetPendingStateTimeout(uint64)"))      // Used in oldZkEvm as well
	setTrustedAggregatorTimeoutSignatureHash       = crypto.Keccak256Hash([]byte("SetTrustedAggregatorTimeout(uint64)")) // Used in oldZkEvm as well
	overridePendingStateSignatureHash              = crypto.Keccak256Hash([]byte("OverridePendingState(uint32,uint64,bytes32,bytes32,address)"))
	proveNonDeterministicPendingStateSignatureHash = crypto.Keccak256Hash([]byte("ProveNonDeterministicPendingState(bytes32,bytes32)")) // Used in oldZkEvm as well
	consolidatePendingStateSignatureHash           = crypto.Keccak256Hash([]byte("ConsolidatePendingState(uint32,uint64,bytes32,bytes32,uint64)"))
	verifyBatchesTrustedAggregatorSignatureHash    = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)"))
	rollupManagerVerifyBatchesSignatureHash        = crypto.Keccak256Hash([]byte("VerifyBatches(uint32,uint64,bytes32,bytes32,address)"))
	onSequenceBatchesSignatureHash                 = crypto.Keccak256Hash([]byte("OnSequenceBatches(uint32,uint64)"))
	updateRollupSignatureHash                      = crypto.Keccak256Hash([]byte("UpdateRollup(uint32,uint32,uint64)"))
	addExistingRollupSignatureHash                 = crypto.Keccak256Hash([]byte("AddExistingRollup(uint32,uint64,address,uint64,uint8,uint64)"))
	createNewRollupSignatureHash                   = crypto.Keccak256Hash([]byte("CreateNewRollup(uint32,uint32,address,uint64,address)"))
	obsoleteRollupTypeSignatureHash                = crypto.Keccak256Hash([]byte("ObsoleteRollupType(uint32)"))
	addNewRollupTypeSignatureHash                  = crypto.Keccak256Hash([]byte("AddNewRollupType(uint32,address,address,uint64,uint8,bytes32,string)"))

	// Events new ZkEvm/RollupBase
	acceptAdminRoleSignatureHash        = crypto.Keccak256Hash([]byte("AcceptAdminRole(address)"))                 // Used in oldZkEvm as well
	transferAdminRoleSignatureHash      = crypto.Keccak256Hash([]byte("TransferAdminRole(address)"))               // Used in oldZkEvm as well
	setForceBatchAddressSignatureHash   = crypto.Keccak256Hash([]byte("SetForceBatchAddress(address)"))            // Used in oldZkEvm as well
	setForceBatchTimeoutSignatureHash   = crypto.Keccak256Hash([]byte("SetForceBatchTimeout(uint64)"))             // Used in oldZkEvm as well
	setTrustedSequencerURLSignatureHash = crypto.Keccak256Hash([]byte("SetTrustedSequencerURL(string)"))           // Used in oldZkEvm as well
	setTrustedSequencerSignatureHash    = crypto.Keccak256Hash([]byte("SetTrustedSequencer(address)"))             // Used in oldZkEvm as well
	verifyBatchesSignatureHash          = crypto.Keccak256Hash([]byte("VerifyBatches(uint64,bytes32,address)"))    // Used in oldZkEvm as well
	sequenceForceBatchesSignatureHash   = crypto.Keccak256Hash([]byte("SequenceForceBatches(uint64)"))             // Used in oldZkEvm as well
	forceBatchSignatureHash             = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)")) // Used in oldZkEvm as well
	sequenceBatchesSignatureHash        = crypto.Keccak256Hash([]byte("SequenceBatches(uint64,bytes32)"))          // Used in oldZkEvm as well
	initialSequenceBatchesSignatureHash = crypto.Keccak256Hash([]byte("InitialSequenceBatches(bytes,bytes32,address)"))
	updateEtrogSequenceSignatureHash    = crypto.Keccak256Hash([]byte("UpdateEtrogSequence(uint64,bytes,bytes32,address)"))

	// Extra RollupManager
	initializedSignatureHash               = crypto.Keccak256Hash([]byte("Initialized(uint64)"))                       // Initializable. Used in RollupBase as well
	roleAdminChangedSignatureHash          = crypto.Keccak256Hash([]byte("RoleAdminChanged(bytes32,bytes32,bytes32)")) // IAccessControlUpgradeable
	roleGrantedSignatureHash               = crypto.Keccak256Hash([]byte("RoleGranted(bytes32,address,address)"))      // IAccessControlUpgradeable
	roleRevokedSignatureHash               = crypto.Keccak256Hash([]byte("RoleRevoked(bytes32,address,address)"))      // IAccessControlUpgradeable
	emergencyStateActivatedSignatureHash   = crypto.Keccak256Hash([]byte("EmergencyStateActivated()"))                 // EmergencyManager. Used in oldZkEvm as well
	emergencyStateDeactivatedSignatureHash = crypto.Keccak256Hash([]byte("EmergencyStateDeactivated()"))               // EmergencyManager. Used in oldZkEvm as well

	// New GER event Etrog
	updateL1InfoTreeSignatureHash = crypto.Keccak256Hash([]byte("UpdateL1InfoTree(bytes32,bytes32)"))

	// PreLxLy events
	updateGlobalExitRootSignatureHash                   = crypto.Keccak256Hash([]byte("UpdateGlobalExitRoot(bytes32,bytes32)"))
	preEtrogVerifyBatchesTrustedAggregatorSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	transferOwnershipSignatureHash                      = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))
	updateZkEVMVersionSignatureHash                     = crypto.Keccak256Hash([]byte("UpdateZkEVMVersion(uint64,uint64,string)"))
	preEtrogConsolidatePendingStateSignatureHash        = crypto.Keccak256Hash([]byte("ConsolidatePendingState(uint64,bytes32,uint64)"))
	preEtrogOverridePendingStateSignatureHash           = crypto.Keccak256Hash([]byte("OverridePendingState(uint64,bytes32,address)"))
	sequenceBatchesPreEtrogSignatureHash                = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))

	// Proxy events
	initializedProxySignatureHash = crypto.Keccak256Hash([]byte("Initialized(uint8)"))
	adminChangedSignatureHash     = crypto.Keccak256Hash([]byte("AdminChanged(address,address)"))
	beaconUpgradedSignatureHash   = crypto.Keccak256Hash([]byte("BeaconUpgraded(address)"))
	upgradedSignatureHash         = crypto.Keccak256Hash([]byte("Upgraded(address)"))

	// methodIDSequenceBatchesEtrog: MethodID for sequenceBatches in Etrog
	methodIDSequenceBatchesEtrog = []byte{0xec, 0xef, 0x3f, 0x99} // 0xecef3f99
	// methodIDSequenceBatchesElderberry: MethodID for sequenceBatches in Elderberry
	methodIDSequenceBatchesElderberry = []byte{0xde, 0xf5, 0x7e, 0x54} // 0xdef57e54 sequenceBatches((bytes,bytes32,uint64,bytes32)[],uint64,uint64,address)

	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("not found")
	// ErrIsReadOnlyMode is used when the EtherMan client is in read-only mode.
	ErrIsReadOnlyMode = errors.New("etherman client in read-only mode: no account configured to send transactions to L1. " +
		"please check the [Etherman] PrivateKeyPath and PrivateKeyPassword configuration")
	// ErrPrivateKeyNotFound used when the provided sender does not have a private key registered to be used
	ErrPrivateKeyNotFound = errors.New("can't find sender private key to sign tx")
)

// SequencedBatchesSigHash returns the hash for the `SequenceBatches` event.
func SequencedBatchesSigHash() common.Hash { return sequenceBatchesSignatureHash }

// TrustedVerifyBatchesSigHash returns the hash for the `TrustedVerifyBatches` event.
func TrustedVerifyBatchesSigHash() common.Hash { return verifyBatchesTrustedAggregatorSignatureHash }

// EventOrder is the the type used to identify the events order
type EventOrder string

const (
	// GlobalExitRootsOrder identifies a GlobalExitRoot event
	GlobalExitRootsOrder EventOrder = "GlobalExitRoots"
	// L1InfoTreeOrder identifies a L1InTree event
	L1InfoTreeOrder EventOrder = "L1InfoTreeOrder"
	// SequenceBatchesOrder identifies a VerifyBatch event
	SequenceBatchesOrder EventOrder = "SequenceBatches"
	// UpdateEtrogSequenceOrder identifies a VerifyBatch event
	UpdateEtrogSequenceOrder EventOrder = "UpdateEtrogSequence"
	// ForcedBatchesOrder identifies a ForcedBatches event
	ForcedBatchesOrder EventOrder = "ForcedBatches"
	// TrustedVerifyBatchOrder identifies a TrustedVerifyBatch event
	TrustedVerifyBatchOrder EventOrder = "TrustedVerifyBatch"
	// VerifyBatchOrder identifies a VerifyBatch event
	VerifyBatchOrder EventOrder = "VerifyBatch"
	// SequenceForceBatchesOrder identifies a SequenceForceBatches event
	SequenceForceBatchesOrder EventOrder = "SequenceForceBatches"
	// ForkIDsOrder identifies an updateZkevmVersion event
	ForkIDsOrder EventOrder = "forkIDs"
	// InitialSequenceBatchesOrder identifies a VerifyBatch event
	InitialSequenceBatchesOrder EventOrder = "InitialSequenceBatches"
)

type ethereumClient interface {
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.LogFilterer
	ethereum.TransactionReader
	ethereum.TransactionSender
	ethereum.PendingStateReader

	bind.DeployBackend
}

// L1Config represents the configuration of the network used in L1
type L1Config struct {
	// Chain ID of the L1 network
	L1ChainID uint64 `json:"chainId"`
	// ZkEVMAddr Address of the L1 contract polygonZkEVMAddress
	ZkEVMAddr common.Address `json:"polygonZkEVMAddress"`
	// RollupManagerAddr Address of the L1 contract
	RollupManagerAddr common.Address `json:"polygonRollupManagerAddress"`
	// PolAddr Address of the L1 Pol token Contract
	PolAddr common.Address `json:"polTokenAddress"`
	// GlobalExitRootManagerAddr Address of the L1 GlobalExitRootManager contract
	GlobalExitRootManagerAddr common.Address `json:"polygonZkEVMGlobalExitRootAddress"`
}

type externalGasProviders struct {
	MultiGasProvider bool
	Providers        []ethereum.GasPricer
}

// Client is a simple implementation of EtherMan.
type Client struct {
	EthClient                     ethereumClient
	PreEtrogZkEVM                 *preetrogpolygonzkevm.Preetrogpolygonzkevm
	ElderberryZKEVM               *elderberrypolygonzkevm.Elderberrypolygonzkevm
	EtrogZkEVM                    *etrogpolygonzkevm.Etrogpolygonzkevm
	EtrogRollupManager            *etrogpolygonrollupmanager.Etrogpolygonrollupmanager
	EtrogGlobalExitRootManager    *etrogpolygonzkevmglobalexitroot.Etrogpolygonzkevmglobalexitroot
	PreEtrogGlobalExitRootManager *preetrogpolygonzkevmglobalexitroot.Preetrogpolygonzkevmglobalexitroot
	FeijoaContracts               *FeijoaContracts
	Pol                           *pol.Pol
	SCAddresses                   []common.Address

	RollupID uint32

	GasProviders externalGasProviders

	l1Cfg              L1Config
	cfg                Config
	auth               map[common.Address]bind.TransactOpts // empty in case of read-only client
	EIP4844            *eip4844.EthermanEIP4844
	eventFeijoaManager *EventManager
}

// NewClient creates a new etherman.
func NewClient(cfg Config, l1Config L1Config) (*Client, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	if cfg.ConsensusL1URL == "" {
		log.Warn("ConsensusL1URL is not set, so Feijoa is not going to work")
	}
	feijoaEnabled := true
	beaconClient := beaconclient.NewBeaconAPIClient(cfg.ConsensusL1URL)
	eip4844 := eip4844.NewEthermanEIP4844(beaconClient)
	if err := eip4844.Initialize(context.Background()); err != nil {
		// TODO: Must be mandatory to have a consensusL1URL configured, but
		// for maintain compatibility allow to disable Feijoa
		// so the log.Warnf must be an Errorf and must return  nil, err
		log.Warnf("error initializing EIP-4844,Feijoa is going to be disabled.  URL:%s : %+v", cfg.ConsensusL1URL, err)
		feijoaEnabled = false
	}
	// Create smc clients
	etrogZkevm, err := etrogpolygonzkevm.NewEtrogpolygonzkevm(l1Config.ZkEVMAddr, ethClient)
	if err != nil {
		log.Errorf("error creating Polygonzkevm client (%s). Error: %w", l1Config.ZkEVMAddr.String(), err)
		return nil, err
	}
	elderberryZkevm, err := elderberrypolygonzkevm.NewElderberrypolygonzkevm(l1Config.RollupManagerAddr, ethClient)
	if err != nil {
		log.Errorf("error creating NewElderberryPolygonzkevm client (%s). Error: %w", l1Config.RollupManagerAddr.String(), err)
		return nil, err
	}
	preEtrogZkevm, err := preetrogpolygonzkevm.NewPreetrogpolygonzkevm(l1Config.RollupManagerAddr, ethClient)
	if err != nil {
		log.Errorf("error creating Newpreetrogpolygonzkevm client (%s). Error: %w", l1Config.RollupManagerAddr.String(), err)
		return nil, err
	}
	etrogRollupManager, err := etrogpolygonrollupmanager.NewEtrogpolygonrollupmanager(l1Config.RollupManagerAddr, ethClient)
	if err != nil {
		log.Errorf("error creating NewPolygonrollupmanager client (%s). Error: %w", l1Config.RollupManagerAddr.String(), err)
		return nil, err
	}
	etrogGlobalExitRoot, err := etrogpolygonzkevmglobalexitroot.NewEtrogpolygonzkevmglobalexitroot(l1Config.GlobalExitRootManagerAddr, ethClient)
	if err != nil {
		log.Errorf("error creating NewPolygonzkevmglobalexitroot client (%s). Error: %w", l1Config.GlobalExitRootManagerAddr.String(), err)
		return nil, err
	}
	preEtrogGlobalExitRoot, err := preetrogpolygonzkevmglobalexitroot.NewPreetrogpolygonzkevmglobalexitroot(l1Config.GlobalExitRootManagerAddr, ethClient)
	if err != nil {
		log.Errorf("error creating Newpreetrogpolygonzkevmglobalexitroot client (%s). Error: %w", l1Config.GlobalExitRootManagerAddr.String(), err)
		return nil, err
	}
	pol, err := pol.NewPol(l1Config.PolAddr, ethClient)
	if err != nil {
		log.Errorf("error creating NewPol client (%s). Error: %w", l1Config.PolAddr.String(), err)
		return nil, err
	}
	feijoaContracts, err := NewFeijoaContracts(ethClient, l1Config)
	if err != nil {
		log.Errorf("error creating NewFeijoaContracts client (%s). Error: %w", l1Config.RollupManagerAddr.String(), err)
		return nil, err
	}
	scAddresses := feijoaContracts.GetAddresses()
	scAddresses = append(scAddresses, l1Config.ZkEVMAddr, l1Config.RollupManagerAddr, l1Config.GlobalExitRootManagerAddr)

	gProviders := []ethereum.GasPricer{ethClient}
	if cfg.MultiGasProvider {
		if cfg.Etherscan.ApiKey == "" {
			log.Info("No ApiKey provided for etherscan. Ignoring provider...")
		} else {
			log.Info("ApiKey detected for etherscan")
			gProviders = append(gProviders, etherscan.NewEtherscanService(cfg.Etherscan.ApiKey))
		}
		gProviders = append(gProviders, ethgasstation.NewEthGasStationService())
	}
	metrics.Register()
	// Get RollupID
	rollupID, err := etrogRollupManager.RollupAddressToID(&bind.CallOpts{Pending: false}, l1Config.ZkEVMAddr)
	if err != nil {
		log.Debugf("error rollupManager.RollupAddressToID(%s). Error: %w", l1Config.RollupManagerAddr, err)
		return nil, err
	}
	log.Debug("rollupID: ", rollupID)

	client := &Client{
		EthClient:                     ethClient,
		EtrogZkEVM:                    etrogZkevm,
		ElderberryZKEVM:               elderberryZkevm,
		PreEtrogZkEVM:                 preEtrogZkevm,
		EtrogRollupManager:            etrogRollupManager,
		Pol:                           pol,
		EtrogGlobalExitRootManager:    etrogGlobalExitRoot,
		PreEtrogGlobalExitRootManager: preEtrogGlobalExitRoot,
		SCAddresses:                   scAddresses,
		RollupID:                      rollupID,
		GasProviders: externalGasProviders{
			MultiGasProvider: cfg.MultiGasProvider,
			Providers:        gProviders,
		},
		l1Cfg:   l1Config,
		cfg:     cfg,
		auth:    map[common.Address]bind.TransactOpts{},
		EIP4844: eip4844,
	}
	if feijoaEnabled {
		eventFeijoaManager := NewEventManager(client, NewCallDataExtratorGeth(ethClient))
		eventFeijoaManager.AddProcessor(NewEventFeijoaSequenceBlobsProcessor(feijoaContracts))
		client.eventFeijoaManager = eventFeijoaManager
	}
	return client, nil
}

// VerifyGenBlockNumber verifies if the genesis Block Number is valid
func (etherMan *Client) VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error) {
	start := time.Now()
	log.Info("Verifying genesis blockNumber: ", genBlockNumber)
	// Filter query
	genBlock := new(big.Int).SetUint64(genBlockNumber)
	query := ethereum.FilterQuery{
		FromBlock: genBlock,
		ToBlock:   genBlock,
		Addresses: etherMan.SCAddresses,
		Topics:    [][]common.Hash{{updateZkEVMVersionSignatureHash, createNewRollupSignatureHash}},
	}
	logs, err := etherMan.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return false, err
	}
	if len(logs) == 0 {
		return false, fmt.Errorf("the specified genBlockNumber in config file does not contain any forkID event. Please use the proper blockNumber.")
	}
	var zkevmVersion preetrogpolygonzkevm.PreetrogpolygonzkevmUpdateZkEVMVersion
	switch logs[0].Topics[0] {
	case updateZkEVMVersionSignatureHash:
		log.Debug("UpdateZkEVMVersion event detected during the Verification of the GenBlockNumber")
		zkevmV, err := etherMan.PreEtrogZkEVM.ParseUpdateZkEVMVersion(logs[0])
		if err != nil {
			return false, err
		}
		if zkevmV != nil {
			zkevmVersion = *zkevmV
		}
	case createNewRollupSignatureHash:
		log.Debug("CreateNewRollup event detected during the Verification of the GenBlockNumber")
		createNewRollupEvent, err := etherMan.EtrogRollupManager.ParseCreateNewRollup(logs[0])
		if err != nil {
			return false, err
		}
		// Query to get the forkID
		rollupType, err := etherMan.EtrogRollupManager.RollupTypeMap(&bind.CallOpts{Pending: false}, createNewRollupEvent.RollupTypeID)
		if err != nil {
			log.Error(err)
			return false, err
		}
		zkevmVersion.ForkID = rollupType.ForkID
		zkevmVersion.NumBatch = 0
	}
	if zkevmVersion.NumBatch != 0 {
		return false, fmt.Errorf("the specified genBlockNumber in config file does not contain the initial forkID event (BatchNum: %d). Please use the proper blockNumber.", zkevmVersion.NumBatch)
	}
	metrics.VerifyGenBlockTime(time.Since(start))
	return true, nil
}

// GetL1BlockUpgradeLxLy It returns the block genesis for LxLy before genesisBlock or error
// TODO: Check if all RPC providers support this range of blocks
func (etherMan *Client) GetL1BlockUpgradeLxLy(ctx context.Context, genesisBlock uint64) (uint64, error) {
	it, err := etherMan.EtrogRollupManager.FilterInitialized(&bind.FilterOpts{
		Start:   1,
		End:     &genesisBlock,
		Context: ctx,
	})
	if err != nil {
		return uint64(0), err
	}
	for it.Next() {
		log.Debugf("BlockNumber: %d Topics:Initialized(%d)", it.Event.Raw.BlockNumber, it.Event.Version)
		if it.Event.Version == ETrogUpgradeVersion { // 2 is ETROG (LxLy upgrade)
			log.Infof("LxLy upgrade found at blockNumber: %d", it.Event.Raw.BlockNumber)
			return it.Event.Raw.BlockNumber, nil
		}
	}
	return uint64(0), ErrNotFound
}

// GetForks returns fork information
func (etherMan *Client) GetForks(ctx context.Context, genBlockNumber uint64, lastL1BlockSynced uint64) ([]state.ForkIDInterval, error) {
	log.Debug("Getting forkIDs from blockNumber: ", genBlockNumber)
	start := time.Now()
	var logs []types.Log
	// At minimum it checks the GenesisBlock
	if lastL1BlockSynced < genBlockNumber {
		lastL1BlockSynced = genBlockNumber
	}
	log.Debug("Using ForkIDChunkSize: ", etherMan.cfg.ForkIDChunkSize)
	for i := genBlockNumber; i <= lastL1BlockSynced; i = i + etherMan.cfg.ForkIDChunkSize + 1 {
		final := i + etherMan.cfg.ForkIDChunkSize
		if final > lastL1BlockSynced {
			// Limit the query to the last l1BlockSynced
			final = lastL1BlockSynced
		}
		log.Debug("INTERVAL. Initial: ", i, ". Final: ", final)
		// Filter query
		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(i),
			ToBlock:   new(big.Int).SetUint64(final),
			Addresses: etherMan.SCAddresses,
			Topics:    [][]common.Hash{{updateZkEVMVersionSignatureHash, updateRollupSignatureHash, addExistingRollupSignatureHash, createNewRollupSignatureHash}},
		}
		l, err := etherMan.EthClient.FilterLogs(ctx, query)
		if err != nil {
			return []state.ForkIDInterval{}, err
		}
		logs = append(logs, l...)
	}

	var forks []state.ForkIDInterval
	for i, l := range logs {
		var zkevmVersion preetrogpolygonzkevm.PreetrogpolygonzkevmUpdateZkEVMVersion
		switch l.Topics[0] {
		case updateZkEVMVersionSignatureHash:
			log.Debug("updateZkEVMVersion Event received")
			zkevmV, err := etherMan.PreEtrogZkEVM.ParseUpdateZkEVMVersion(l)
			if err != nil {
				return []state.ForkIDInterval{}, err
			}
			if zkevmV != nil {
				zkevmVersion = *zkevmV
			}
		case updateRollupSignatureHash:
			log.Debug("updateRollup Event received")
			updateRollupEvent, err := etherMan.EtrogRollupManager.ParseUpdateRollup(l)
			if err != nil {
				return []state.ForkIDInterval{}, err
			}
			if etherMan.RollupID != updateRollupEvent.RollupID {
				continue
			}
			// Query to get the forkID
			rollupType, err := etherMan.EtrogRollupManager.RollupTypeMap(&bind.CallOpts{Pending: false}, updateRollupEvent.NewRollupTypeID)
			if err != nil {
				return []state.ForkIDInterval{}, err
			}
			zkevmVersion.ForkID = rollupType.ForkID
			zkevmVersion.NumBatch = updateRollupEvent.LastVerifiedBatchBeforeUpgrade

		case addExistingRollupSignatureHash:
			log.Debug("addExistingRollup Event received")
			addExistingRollupEvent, err := etherMan.EtrogRollupManager.ParseAddExistingRollup(l)
			if err != nil {
				return []state.ForkIDInterval{}, err
			}
			if etherMan.RollupID != addExistingRollupEvent.RollupID {
				continue
			}
			zkevmVersion.ForkID = addExistingRollupEvent.ForkID
			zkevmVersion.NumBatch = addExistingRollupEvent.LastVerifiedBatchBeforeUpgrade

		case createNewRollupSignatureHash:
			log.Debug("createNewRollup Event received")
			createNewRollupEvent, err := etherMan.EtrogRollupManager.ParseCreateNewRollup(l)
			if err != nil {
				return []state.ForkIDInterval{}, err
			}
			if etherMan.RollupID != createNewRollupEvent.RollupID {
				continue
			}
			// Query to get the forkID
			rollupType, err := etherMan.EtrogRollupManager.RollupTypeMap(&bind.CallOpts{Pending: false}, createNewRollupEvent.RollupTypeID)
			if err != nil {
				log.Error(err)
				return []state.ForkIDInterval{}, err
			}
			zkevmVersion.ForkID = rollupType.ForkID
			zkevmVersion.NumBatch = 0
		}
		var fork state.ForkIDInterval
		if i == 0 {
			fork = state.ForkIDInterval{
				FromBatchNumber: zkevmVersion.NumBatch + 1,
				ToBatchNumber:   math.MaxUint64,
				ForkId:          zkevmVersion.ForkID,
				Version:         zkevmVersion.Version,
				BlockNumber:     l.BlockNumber,
			}
		} else {
			forks[len(forks)-1].ToBatchNumber = zkevmVersion.NumBatch
			fork = state.ForkIDInterval{
				FromBatchNumber: zkevmVersion.NumBatch + 1,
				ToBatchNumber:   math.MaxUint64,
				ForkId:          zkevmVersion.ForkID,
				Version:         zkevmVersion.Version,
				BlockNumber:     l.BlockNumber,
			}
		}
		forks = append(forks, fork)
	}
	metrics.GetForksTime(time.Since(start))
	log.Debugf("ForkIDs found: %+v", forks)
	return forks, nil
}

// GetRollupInfoByBlockRange function retrieves the Rollup information that are included in all this ethereum blocks
// from block x to block y.
func (etherMan *Client) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]Block, map[common.Hash][]Order, error) {
	// Filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: etherMan.SCAddresses,
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, blocksOrder, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return blocks, blocksOrder, nil
}

// GetRollupInfoByBlockRangePreviousRollupGenesis function retrieves the Rollup information that are included in all this ethereum blocks
// but it only retrieves the information from the previous rollup genesis block to the current block.
func (etherMan *Client) GetRollupInfoByBlockRangePreviousRollupGenesis(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]Block, map[common.Hash][]Order, error) {
	// Filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: []common.Address{etherMan.l1Cfg.GlobalExitRootManagerAddr},
		Topics:    [][]common.Hash{{updateL1InfoTreeSignatureHash}},
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, blocksOrder, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return blocks, blocksOrder, nil
}

// Order contains the event order to let the synchronizer store the information following this order.
type Order struct {
	Name EventOrder
	Pos  int
}

func (etherMan *Client) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]Block, map[common.Hash][]Order, error) {
	start := time.Now()
	logs, err := etherMan.EthClient.FilterLogs(ctx, query)
	metrics.GetEventsTime(time.Since(start))
	if err != nil {
		return nil, nil, err
	}
	var blocks []Block
	blocksOrder := make(map[common.Hash][]Order)
	startProcess := time.Now()
	for _, vLog := range logs {
		startProcessSingleEvent := time.Now()
		err := etherMan.processEvent(ctx, vLog, &blocks, &blocksOrder)
		metrics.ProcessSingleEventTime(time.Since(startProcessSingleEvent))
		metrics.EventCounter()
		if err != nil {
			log.Warnf("error processing event. Retrying... Error: %s. vLog: %+v", err.Error(), vLog)
			return nil, nil, err
		}
	}
	metrics.ProcessAllEventTime(time.Since(startProcess))
	metrics.ReadAndProcessAllEventsTime(time.Since(start))
	return blocks, blocksOrder, nil
}
func (etherMan *Client) processEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	if etherMan.eventFeijoaManager != nil {
		processed, err := etherMan.eventFeijoaManager.ProcessEvent(ctx, vLog, blocks, blocksOrder)
		if processed || err != nil {
			return err
		}
	}
	return etherMan.processEventLegacy(ctx, vLog, blocks, blocksOrder)
}

func (etherMan *Client) processEventLegacy(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	switch vLog.Topics[0] {
	case sequenceBatchesSignatureHash:
		return etherMan.sequencedBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case sequenceBatchesPreEtrogSignatureHash:
		return etherMan.sequencedBatchesPreEtrogEvent(ctx, vLog, blocks, blocksOrder)
	case updateGlobalExitRootSignatureHash:
		return etherMan.updateGlobalExitRootEvent(ctx, vLog, blocks, blocksOrder)
	case updateL1InfoTreeSignatureHash:
		return etherMan.updateL1InfoTreeEvent(ctx, vLog, blocks, blocksOrder)
	case forceBatchSignatureHash:
		return etherMan.forcedBatchEvent(ctx, vLog, blocks, blocksOrder)
	case initialSequenceBatchesSignatureHash:
		return etherMan.initialSequenceBatches(ctx, vLog, blocks, blocksOrder)
	case updateEtrogSequenceSignatureHash:
		return etherMan.updateEtrogSequence(ctx, vLog, blocks, blocksOrder)
	case verifyBatchesTrustedAggregatorSignatureHash:
		log.Debug("VerifyBatchesTrustedAggregator event detected. Ignoring...")
		return nil
	case rollupManagerVerifyBatchesSignatureHash:
		log.Debug("RollupManagerVerifyBatches event detected. Ignoring...")
		return nil
	case preEtrogVerifyBatchesTrustedAggregatorSignatureHash:
		return etherMan.preEtrogVerifyBatchesTrustedAggregatorEvent(ctx, vLog, blocks, blocksOrder)
	case verifyBatchesSignatureHash:
		return etherMan.verifyBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case sequenceForceBatchesSignatureHash:
		return etherMan.forceSequencedBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case setTrustedSequencerURLSignatureHash:
		log.Debug("SetTrustedSequencerURL event detected. Ignoring...")
		return nil
	case setTrustedSequencerSignatureHash:
		log.Debug("SetTrustedSequencer event detected. Ignoring...")
		return nil
	case initializedSignatureHash:
		log.Debug("Initialized event detected. Ignoring...")
		return nil
	case initializedProxySignatureHash:
		log.Debug("InitializedProxy event detected. Ignoring...")
		return nil
	case adminChangedSignatureHash:
		log.Debug("AdminChanged event detected. Ignoring...")
		return nil
	case beaconUpgradedSignatureHash:
		log.Debug("BeaconUpgraded event detected. Ignoring...")
		return nil
	case upgradedSignatureHash:
		log.Debug("Upgraded event detected. Ignoring...")
		return nil
	case transferOwnershipSignatureHash:
		log.Debug("TransferOwnership event detected. Ignoring...")
		return nil
	case emergencyStateActivatedSignatureHash:
		log.Debug("EmergencyStateActivated event detected. Ignoring...")
		return nil
	case emergencyStateDeactivatedSignatureHash:
		log.Debug("EmergencyStateDeactivated event detected. Ignoring...")
		return nil
	case updateZkEVMVersionSignatureHash:
		return etherMan.updateZkevmVersion(ctx, vLog, blocks, blocksOrder)
	case consolidatePendingStateSignatureHash:
		log.Debug("ConsolidatePendingState event detected. Ignoring...")
		return nil
	case preEtrogConsolidatePendingStateSignatureHash:
		log.Debug("PreEtrogConsolidatePendingState event detected. Ignoring...")
		return nil
	case setTrustedAggregatorTimeoutSignatureHash:
		log.Debug("SetTrustedAggregatorTimeout event detected. Ignoring...")
		return nil
	case setTrustedAggregatorSignatureHash:
		log.Debug("SetTrustedAggregator event detected. Ignoring...")
		return nil
	case setPendingStateTimeoutSignatureHash:
		log.Debug("SetPendingStateTimeout event detected. Ignoring...")
		return nil
	case setMultiplierBatchFeeSignatureHash:
		log.Debug("SetMultiplierBatchFee event detected. Ignoring...")
		return nil
	case setVerifyBatchTimeTargetSignatureHash:
		log.Debug("SetVerifyBatchTimeTarget event detected. Ignoring...")
		return nil
	case setForceBatchTimeoutSignatureHash:
		log.Debug("SetForceBatchTimeout event detected. Ignoring...")
		return nil
	case setForceBatchAddressSignatureHash:
		log.Debug("SetForceBatchAddress event detected. Ignoring...")
		return nil
	case transferAdminRoleSignatureHash:
		log.Debug("TransferAdminRole event detected. Ignoring...")
		return nil
	case acceptAdminRoleSignatureHash:
		log.Debug("AcceptAdminRole event detected. Ignoring...")
		return nil
	case proveNonDeterministicPendingStateSignatureHash:
		log.Debug("ProveNonDeterministicPendingState event detected. Ignoring...")
		return nil
	case overridePendingStateSignatureHash:
		log.Debug("OverridePendingState event detected. Ignoring...")
		return nil
	case preEtrogOverridePendingStateSignatureHash:
		log.Debug("PreEtrogOverridePendingState event detected. Ignoring...")
		return nil
	case roleAdminChangedSignatureHash:
		log.Debug("RoleAdminChanged event detected. Ignoring...")
		return nil
	case roleGrantedSignatureHash:
		log.Debug("RoleGranted event detected. Ignoring...")
		return nil
	case roleRevokedSignatureHash:
		log.Debug("RoleRevoked event detected. Ignoring...")
		return nil
	case onSequenceBatchesSignatureHash:
		log.Debug("OnSequenceBatches event detected. Ignoring...")
		return nil
	case updateRollupSignatureHash:
		return etherMan.updateRollup(ctx, vLog, blocks, blocksOrder)
	case addExistingRollupSignatureHash:
		return etherMan.addExistingRollup(ctx, vLog, blocks, blocksOrder)
	case createNewRollupSignatureHash:
		return etherMan.createNewRollup(ctx, vLog, blocks, blocksOrder)
	case obsoleteRollupTypeSignatureHash:
		log.Debug("ObsoleteRollupType event detected. Ignoring...")
		return nil
	case addNewRollupTypeSignatureHash:
		log.Debug("addNewRollupType event detected but not implemented. Ignoring...")
		return nil
	case setBatchFeeSignatureHash:
		log.Debug("SetBatchFee event detected. Ignoring...")
		return nil
	}
	log.Warnf("Event not registered: %+v", vLog)
	return nil
}

func (etherMan *Client) updateZkevmVersion(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateZkEVMVersion event detected")
	zkevmVersion, err := etherMan.PreEtrogZkEVM.ParseUpdateZkEVMVersion(vLog)
	if err != nil {
		log.Error("error parsing UpdateZkEVMVersion event. Error: ", err)
		return err
	}
	return etherMan.updateForkId(ctx, vLog, blocks, blocksOrder, zkevmVersion.NumBatch, zkevmVersion.ForkID, zkevmVersion.Version, etherMan.RollupID)
}

func (etherMan *Client) updateRollup(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateRollup event detected")
	updateRollup, err := etherMan.EtrogRollupManager.ParseUpdateRollup(vLog)
	if err != nil {
		log.Error("error parsing UpdateRollup event. Error: ", err)
		return err
	}
	rollupType, err := etherMan.EtrogRollupManager.RollupTypeMap(&bind.CallOpts{Pending: false}, updateRollup.NewRollupTypeID)
	if err != nil {
		return err
	}
	return etherMan.updateForkId(ctx, vLog, blocks, blocksOrder, updateRollup.LastVerifiedBatchBeforeUpgrade, rollupType.ForkID, "", updateRollup.RollupID)
}

func (etherMan *Client) createNewRollup(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("createNewRollup event detected")
	createRollup, err := etherMan.EtrogRollupManager.ParseCreateNewRollup(vLog)
	if err != nil {
		log.Error("error parsing createNewRollup event. Error: ", err)
		return err
	}
	rollupType, err := etherMan.EtrogRollupManager.RollupTypeMap(&bind.CallOpts{Pending: false}, createRollup.RollupTypeID)
	if err != nil {
		return err
	}
	return etherMan.updateForkId(ctx, vLog, blocks, blocksOrder, 0, rollupType.ForkID, "", createRollup.RollupID)
}

func (etherMan *Client) addExistingRollup(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("addExistingRollup event detected")
	addExistingRollup, err := etherMan.EtrogRollupManager.ParseAddExistingRollup(vLog)
	if err != nil {
		log.Error("error parsing createNewRollup event. Error: ", err)
		return err
	}

	return etherMan.updateForkId(ctx, vLog, blocks, blocksOrder, addExistingRollup.LastVerifiedBatchBeforeUpgrade, addExistingRollup.ForkID, "", addExistingRollup.RollupID)
}

func (etherMan *Client) updateEtrogSequence(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("updateEtrogSequence event detected")
	updateEtrogSequence, err := etherMan.ElderberryZKEVM.ParseUpdateEtrogSequence(vLog)
	if err != nil {
		log.Error("error parsing updateEtrogSequence event. Error: ", err)
		return err
	}

	// Read the tx for this event.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting fullBlockInfo. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}

	log.Info("update Etrog transaction sequence...")
	sequence := UpdateEtrogSequence{
		BatchNumber:   updateEtrogSequence.NumBatch,
		SequencerAddr: updateEtrogSequence.Sequencer,
		TxHash:        vLog.TxHash,
		Nonce:         msg.Nonce,
		PolygonRollupBaseEtrogBatchData: &etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions:         updateEtrogSequence.Transactions,
			ForcedGlobalExitRoot: updateEtrogSequence.LastGlobalExitRoot,
			ForcedTimestamp:      fullBlock.Time(),
			ForcedBlockHashL1:    fullBlock.ParentHash(),
		},
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.UpdateEtrogSequence = sequence
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].UpdateEtrogSequence = sequence
	} else {
		log.Error("Error processing UpdateEtrogSequence event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing UpdateEtrogSequence event")
	}
	or := Order{
		Name: UpdateEtrogSequenceOrder,
		Pos:  0,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) initialSequenceBatches(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("initialSequenceBatches event detected")
	initialSequenceBatches, err := etherMan.EtrogZkEVM.ParseInitialSequenceBatches(vLog)
	if err != nil {
		log.Error("error parsing initialSequenceBatches event. Error: ", err)
		return err
	}

	// Read the tx for this event.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting fullBlockInfo. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}

	var sequences []SequencedBatch
	log.Info("initial transaction sequence...")
	sequences = append(sequences, SequencedBatch{
		BatchNumber:   1,
		SequencerAddr: initialSequenceBatches.Sequencer,
		TxHash:        vLog.TxHash,
		Nonce:         msg.Nonce,
		PolygonRollupBaseEtrogBatchData: &etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions:         initialSequenceBatches.Transactions,
			ForcedGlobalExitRoot: initialSequenceBatches.LastGlobalExitRoot,
			ForcedTimestamp:      fullBlock.Time(),
			ForcedBlockHashL1:    fullBlock.ParentHash(),
		},
	})

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedBatches = append(block.SequencedBatches, sequences)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedBatches = append((*blocks)[len(*blocks)-1].SequencedBatches, sequences)
	} else {
		log.Error("Error processing SequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing SequencedBatches event")
	}
	or := Order{
		Name: InitialSequenceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}
func (etherMan *Client) updateForkId(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order, batchNum, forkID uint64, version string, affectedRollupID uint32) error {
	if etherMan.RollupID != affectedRollupID {
		log.Debug("ignoring this event because it is related to another rollup %d, we are rollupID %d", affectedRollupID, etherMan.RollupID)
		return nil
	}
	fork := ForkID{
		BatchNumber: batchNum,
		ForkID:      forkID,
		Version:     version,
	}
	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		t := time.Unix(int64(fullBlock.Time()), 0)
		block := prepareBlock(vLog, t, fullBlock)
		block.ForkIDs = append(block.ForkIDs, fork)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].ForkIDs = append((*blocks)[len(*blocks)-1].ForkIDs, fork)
	} else {
		log.Error("Error processing updateZkevmVersion event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing updateZkevmVersion event")
	}
	or := Order{
		Name: ForkIDsOrder,
		Pos:  len((*blocks)[len(*blocks)-1].ForkIDs) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) updateL1InfoTreeEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateL1InfoTree event detected")
	etrogGlobalExitRootL1InfoTree, err := etherMan.EtrogGlobalExitRootManager.ParseUpdateL1InfoTree(vLog)
	if err != nil {
		return err
	}

	var gExitRoot GlobalExitRoot
	gExitRoot.MainnetExitRoot = etrogGlobalExitRootL1InfoTree.MainnetExitRoot
	gExitRoot.RollupExitRoot = etrogGlobalExitRootL1InfoTree.RollupExitRoot
	gExitRoot.BlockNumber = vLog.BlockNumber
	gExitRoot.GlobalExitRoot = hash(etrogGlobalExitRootL1InfoTree.MainnetExitRoot, etrogGlobalExitRootL1InfoTree.RollupExitRoot)
	var block *Block
	if !isheadBlockInArray(blocks, vLog.BlockHash, vLog.BlockNumber) {
		// Need to add the block, doesnt mind if inside the blocks because I have to respect the order so insert at end
		block, err = etherMan.RetrieveFullBlockForEvent(ctx, vLog)
		if err != nil {
			return err
		}
		*blocks = append(*blocks, *block)
	}
	// Get the block in the HEAD of the array that contain the current block
	block = &(*blocks)[len(*blocks)-1]
	gExitRoot.PreviousBlockHash = block.ParentHash
	gExitRoot.Timestamp = block.ReceivedAt
	// Add the event to the block
	block.L1InfoTree = append(block.L1InfoTree, gExitRoot)
	order := Order{
		Name: L1InfoTreeOrder,
		Pos:  len(block.L1InfoTree) - 1,
	}
	(*blocksOrder)[block.BlockHash] = append((*blocksOrder)[block.BlockHash], order)
	return nil
}

// RetrieveFullBlockForEvent retrieves the full block for a given event
func (etherMan *Client) RetrieveFullBlockForEvent(ctx context.Context, vLog types.Log) (*Block, error) {
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	t := time.Unix(int64(fullBlock.Time()), 0)
	block := prepareBlock(vLog, t, fullBlock)
	return &block, nil
}

// Check if head block in blocks array is the same as blockHash / blockNumber
func isheadBlockInArray(blocks *[]Block, blockHash common.Hash, blockNumber uint64) bool {
	// Check last item on array blocks if match Hash and Number
	headBlockIsNotExpected := len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != blockHash || (*blocks)[len(*blocks)-1].BlockNumber != blockNumber)
	return !headBlockIsNotExpected
}

func (etherMan *Client) updateGlobalExitRootEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateGlobalExitRoot event detected")
	preEtrogGlobalExitRoot, err := etherMan.PreEtrogGlobalExitRootManager.ParseUpdateGlobalExitRoot(vLog)
	if err != nil {
		return err
	}
	return etherMan.processUpdateGlobalExitRootEvent(ctx, preEtrogGlobalExitRoot.MainnetExitRoot, preEtrogGlobalExitRoot.RollupExitRoot, vLog, blocks, blocksOrder)
}

func (etherMan *Client) processUpdateGlobalExitRootEvent(ctx context.Context, mainnetExitRoot, rollupExitRoot common.Hash, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	var gExitRoot GlobalExitRoot
	gExitRoot.MainnetExitRoot = mainnetExitRoot
	gExitRoot.RollupExitRoot = rollupExitRoot
	gExitRoot.BlockNumber = vLog.BlockNumber
	gExitRoot.GlobalExitRoot = hash(mainnetExitRoot, rollupExitRoot)

	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	t := time.Unix(int64(fullBlock.Time()), 0)
	gExitRoot.Timestamp = t

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, t, fullBlock)
		block.GlobalExitRoots = append(block.GlobalExitRoots, gExitRoot)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].GlobalExitRoots = append((*blocks)[len(*blocks)-1].GlobalExitRoots, gExitRoot)
	} else {
		log.Error("Error processing UpdateGlobalExitRoot event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing UpdateGlobalExitRoot event")
	}
	or := Order{
		Name: GlobalExitRootsOrder,
		Pos:  len((*blocks)[len(*blocks)-1].GlobalExitRoots) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

// WaitTxToBeMined waits for an L1 tx to be mined. It will return error if the tx is reverted or timeout is exceeded
func (etherMan *Client) WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) (bool, error) {
	err := operations.WaitTxToBeMined(ctx, etherMan.EthClient, tx, timeout)
	if errors.Is(err, context.DeadlineExceeded) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// EstimateGasSequenceBatches estimates gas for sending batches
func (etherMan *Client) EstimateGasSequenceBatches(sender common.Address, sequences []ethmanTypes.Sequence, maxSequenceTimestamp uint64, lastSequencedBatchNumber uint64, l2Coinbase common.Address) (*types.Transaction, error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, ErrPrivateKeyNotFound
	}
	opts.NoSend = true

	tx, err := etherMan.sequenceBatches(opts, sequences, maxSequenceTimestamp, lastSequencedBatchNumber, l2Coinbase)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// BuildSequenceBatchesTxData builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence, maxSequenceTimestamp uint64, lastSequencedBatchNumber uint64, l2Coinbase common.Address) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceBatches(opts, sequences, maxSequenceTimestamp, lastSequencedBatchNumber, l2Coinbase)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceBatches(opts bind.TransactOpts, sequences []ethmanTypes.Sequence, maxSequenceTimestamp uint64, lastSequencedBatchNumber uint64, l2Coinbase common.Address) (*types.Transaction, error) {
	var batches []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData
	for _, seq := range sequences {
		var ger common.Hash
		if seq.ForcedBatchTimestamp > 0 {
			ger = seq.GlobalExitRoot
		}
		batch := etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions:         seq.BatchL2Data,
			ForcedGlobalExitRoot: ger,
			ForcedTimestamp:      uint64(seq.ForcedBatchTimestamp),
			ForcedBlockHashL1:    seq.PrevBlockHash,
		}

		batches = append(batches, batch)
	}

	tx, err := etherMan.EtrogZkEVM.SequenceBatches(&opts, batches, maxSequenceTimestamp, lastSequencedBatchNumber, l2Coinbase)
	if err != nil {
		log.Debugf("Batches to send: %+v", batches)
		log.Debug("l2CoinBase: ", l2Coinbase)
		log.Debug("Sequencer address: ", opts.From)
		a, err2 := etrogpolygonzkevm.EtrogpolygonzkevmMetaData.GetAbi()
		if err2 != nil {
			log.Error("error getting abi. Error: ", err2)
		}
		input, err3 := a.Pack("sequenceBatches", batches, maxSequenceTimestamp, lastSequencedBatchNumber, l2Coinbase)
		if err3 != nil {
			log.Error("error packing call. Error: ", err3)
		}
		ctx := context.Background()
		var b string
		block, err4 := etherMan.EthClient.BlockByNumber(ctx, nil)
		if err4 != nil {
			log.Error("error getting blockNumber. Error: ", err4)
			b = "latest"
		} else {
			b = fmt.Sprintf("%x", block.Number())
		}
		log.Warnf(`Use the next command to debug it manually.
		curl --location --request POST 'http://localhost:8545' \
		--header 'Content-Type: application/json' \
		--data-raw '{
			"jsonrpc": "2.0",
			"method": "eth_call",
			"params": [{"from": "%s","to":"%s","data":"0x%s"},"0x%s"],
			"id": 1
		}'`, opts.From, &etherMan.SCAddresses[0], common.Bytes2Hex(input), b)
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
	}

	return tx, err
}

// BuildTrustedVerifyBatchesTxData builds a []bytes to be sent to the PoE SC method TrustedVerifyBatches.
func (etherMan *Client) BuildTrustedVerifyBatchesTxData(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs, beneficiary common.Address) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.generateRandomAuth()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build trusted verify batches, err: %w", err)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	var newLocalExitRoot [32]byte
	copy(newLocalExitRoot[:], inputs.NewLocalExitRoot)

	var newStateRoot [32]byte
	copy(newStateRoot[:], inputs.NewStateRoot)

	proof, err := convertProof(inputs.FinalProof.Proof)
	if err != nil {
		log.Errorf("error converting proof. Error: %v, Proof: %s", err, inputs.FinalProof.Proof)
		return nil, nil, err
	}

	const pendStateNum = 0 // TODO hardcoded for now until we implement the pending state feature

	tx, err := etherMan.EtrogRollupManager.VerifyBatchesTrustedAggregator(
		&opts,
		etherMan.RollupID,
		pendStateNum,
		lastVerifiedBatch,
		newVerifiedBatch,
		newLocalExitRoot,
		newStateRoot,
		beneficiary,
		proof,
	)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func convertProof(p string) ([24][32]byte, error) {
	if len(p) != 24*32*2+2 {
		return [24][32]byte{}, fmt.Errorf("invalid proof length. Length: %d", len(p))
	}
	p = strings.TrimPrefix(p, "0x")
	proof := [24][32]byte{}
	for i := 0; i < 24; i++ {
		data := p[i*64 : (i+1)*64]
		p, err := encoding.DecodeBytes(&data)
		if err != nil {
			return [24][32]byte{}, fmt.Errorf("failed to decode proof, err: %w", err)
		}
		var aux [32]byte
		copy(aux[:], p)
		proof[i] = aux
	}
	return proof, nil
}

// GetSendSequenceFee get super/trusted sequencer fee
func (etherMan *Client) GetSendSequenceFee(numBatches uint64) (*big.Int, error) {
	f, err := etherMan.EtrogRollupManager.GetBatchFee(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, err
	}
	fee := new(big.Int).Mul(f, new(big.Int).SetUint64(numBatches))
	return fee, nil
}

// TrustedSequencer gets trusted sequencer address
func (etherMan *Client) TrustedSequencer() (common.Address, error) {
	return etherMan.EtrogZkEVM.TrustedSequencer(&bind.CallOpts{Pending: false})
}

func (etherMan *Client) forcedBatchEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("ForceBatch event detected")
	fb, err := etherMan.EtrogZkEVM.ParseForceBatch(vLog)
	if err != nil {
		return err
	}
	var forcedBatch ForcedBatch
	forcedBatch.BlockNumber = vLog.BlockNumber
	forcedBatch.ForcedBatchNumber = fb.ForceBatchNum
	forcedBatch.GlobalExitRoot = fb.LastGlobalExitRoot

	// Read the tx for this batch.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}

	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	if fb.Sequencer == msg.From {
		txData := tx.Data()
		// Extract coded txs.
		// Load contract ABI
		abi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
		if err != nil {
			return err
		}

		// Recover Method from signature and ABI
		method, err := abi.MethodById(txData[:4])
		if err != nil {
			return err
		}

		// Unpack method inputs
		data, err := method.Inputs.Unpack(txData[4:])
		if err != nil {
			return err
		}
		bytedata := data[0].([]byte)
		forcedBatch.RawTxsData = bytedata
	} else {
		forcedBatch.RawTxsData = fb.Transactions
	}
	forcedBatch.Sequencer = fb.Sequencer
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	t := time.Unix(int64(fullBlock.Time()), 0)
	forcedBatch.ForcedAt = t
	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, t, fullBlock)
		block.ForcedBatches = append(block.ForcedBatches, forcedBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].ForcedBatches = append((*blocks)[len(*blocks)-1].ForcedBatches, forcedBatch)
	} else {
		log.Error("Error processing ForceBatch event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing ForceBatch event")
	}
	or := Order{
		Name: ForcedBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].ForcedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) sequencedBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debugf("SequenceBatches event detected: txHash: %s", common.Bytes2Hex(vLog.TxHash[:]))

	sb, err := etherMan.EtrogZkEVM.ParseSequenceBatches(vLog)
	if err != nil {
		return err
	}

	// Read the tx for this event.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewCancunSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}

	var sequences []SequencedBatch
	if sb.NumBatch != 1 {
		methodId := tx.Data()[:4]
		log.Debugf("MethodId: %s", common.Bytes2Hex(methodId))
		if bytes.Equal(methodId, methodIDSequenceBatchesEtrog) {
			sequences, err = decodeSequencesEtrog(tx.Data(), sb.NumBatch, msg.From, vLog.TxHash, msg.Nonce, sb.L1InfoRoot)
			if err != nil {
				return fmt.Errorf("error decoding the sequences (etrog): %v", err)
			}
		} else if bytes.Equal(methodId, methodIDSequenceBatchesElderberry) {
			sequences, err = decodeSequencesElderberry(tx.Data(), sb.NumBatch, msg.From, vLog.TxHash, msg.Nonce, sb.L1InfoRoot)
			if err != nil {
				return fmt.Errorf("error decoding the sequences (elderberry): %v", err)
			}
		} else {
			return fmt.Errorf("error decoding the sequences: methodId %s unknown", common.Bytes2Hex(methodId))
		}
	} else {
		log.Info("initial transaction sequence...")
		sequences = append(sequences, SequencedBatch{
			BatchNumber:   1,
			SequencerAddr: msg.From,
			TxHash:        vLog.TxHash,
			Nonce:         msg.Nonce,
		})
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedBatches = append(block.SequencedBatches, sequences)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedBatches = append((*blocks)[len(*blocks)-1].SequencedBatches, sequences)
	} else {
		log.Error("Error processing SequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing SequencedBatches event")
	}
	or := Order{
		Name: SequenceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) sequencedBatchesPreEtrogEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("Pre etrog SequenceBatches event detected")
	sb, err := etherMan.PreEtrogZkEVM.ParseSequenceBatches(vLog)
	if err != nil {
		return err
	}

	// Read the tx for this event.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}

	sequences, err := decodeSequencesPreEtrog(tx.Data(), sb.NumBatch, msg.From, vLog.TxHash, msg.Nonce)
	if err != nil {
		return fmt.Errorf("error decoding the sequences: %v", err)
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedBatches = append(block.SequencedBatches, sequences)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedBatches = append((*blocks)[len(*blocks)-1].SequencedBatches, sequences)
	} else {
		log.Error("Error processing SequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing SequencedBatches event")
	}
	or := Order{
		Name: SequenceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func decodeSequencesElderberry(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := smcAbi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}
	var sequences []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}
	maxSequenceTimestamp := data[1].(uint64)
	initSequencedBatchNumber := data[2].(uint64)
	coinbase := (data[3]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))

	for i, seq := range sequences {
		elderberry := SequencedBatchElderberryData{
			MaxSequenceTimestamp:     maxSequenceTimestamp,
			InitSequencedBatchNumber: initSequencedBatchNumber,
		}
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := seq
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:                     bn,
			L1InfoRoot:                      &l1InfoRoot,
			SequencerAddr:                   sequencer,
			TxHash:                          txHash,
			Nonce:                           nonce,
			Coinbase:                        coinbase,
			PolygonRollupBaseEtrogBatchData: &s,
			SequencedBatchElderberryData:    &elderberry,
		}
	}

	return sequencedBatches, nil
}

func decodeSequencesEtrog(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi, err := abi.JSON(strings.NewReader(elderberrypolygonzkevm.ElderberrypolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := smcAbi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}
	var sequences []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}
	coinbase := (data[1]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))
	for i, seq := range sequences {
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := seq
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:                     bn,
			L1InfoRoot:                      &l1InfoRoot,
			SequencerAddr:                   sequencer,
			TxHash:                          txHash,
			Nonce:                           nonce,
			Coinbase:                        coinbase,
			PolygonRollupBaseEtrogBatchData: &s,
		}
	}

	return sequencedBatches, nil
}

func decodeSequencesPreEtrog(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi, err := abi.JSON(strings.NewReader(preetrogpolygonzkevm.PreetrogpolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := smcAbi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}
	var sequences []preetrogpolygonzkevm.PolygonZkEVMBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}
	coinbase := (data[1]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))
	for i, seq := range sequences {
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := seq
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:           bn,
			SequencerAddr:         sequencer,
			TxHash:                txHash,
			Nonce:                 nonce,
			Coinbase:              coinbase,
			PolygonZkEVMBatchData: &s,
		}
	}

	return sequencedBatches, nil
}

func (etherMan *Client) preEtrogVerifyBatchesTrustedAggregatorEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("TrustedVerifyBatches event detected")
	var vb *preetrogpolygonzkevm.PreetrogpolygonzkevmVerifyBatchesTrustedAggregator
	vb, err := etherMan.PreEtrogZkEVM.ParseVerifyBatchesTrustedAggregator(vLog)
	if err != nil {
		log.Error("error parsing TrustedVerifyBatches event. Error: ", err)
		return err
	}
	return etherMan.verifyBatches(ctx, vLog, blocks, blocksOrder, vb.NumBatch, vb.StateRoot, vb.Aggregator, TrustedVerifyBatchOrder)
}

func (etherMan *Client) verifyBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("VerifyBatches event detected")
	vb, err := etherMan.EtrogZkEVM.ParseVerifyBatches(vLog)
	if err != nil {
		log.Error("error parsing VerifyBatches event. Error: ", err)
		return err
	}
	return etherMan.verifyBatches(ctx, vLog, blocks, blocksOrder, vb.NumBatch, vb.StateRoot, vb.Aggregator, VerifyBatchOrder)
}
func (etherMan *Client) verifyBatches(
	ctx context.Context,
	vLog types.Log,
	blocks *[]Block,
	blocksOrder *map[common.Hash][]Order,
	numBatch uint64,
	stateRoot common.Hash,
	aggregator common.Address,
	orderName EventOrder) error {
	var verifyBatch VerifiedBatch
	verifyBatch.BlockNumber = vLog.BlockNumber
	verifyBatch.BatchNumber = numBatch
	verifyBatch.TxHash = vLog.TxHash
	verifyBatch.StateRoot = stateRoot
	verifyBatch.Aggregator = aggregator

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.VerifiedBatches = append(block.VerifiedBatches, verifyBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].VerifiedBatches = append((*blocks)[len(*blocks)-1].VerifiedBatches, verifyBatch)
	} else {
		log.Error("Error processing verifyBatch event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing verifyBatch event")
	}
	or := Order{
		Name: orderName,
		Pos:  len((*blocks)[len(*blocks)-1].VerifiedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) forceSequencedBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("SequenceForceBatches event detect")
	fsb, err := etherMan.EtrogZkEVM.ParseSequenceForceBatches(vLog)
	if err != nil {
		return err
	}
	// TODO complete data forcedBlockHash, forcedGer y forcedTimestamp

	// Read the tx for this batch.
	tx, err := etherMan.EthClient.TransactionInBlock(ctx, vLog.BlockHash, vLog.TxIndex)
	if err != nil {
		return err
	}
	if tx.Hash() != vLog.TxHash {
		return fmt.Errorf("error: tx hash mismatch. want: %s have: %s", vLog.TxHash, tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	sequencedForceBatch, err := decodeSequencedForceBatches(tx.Data(), fsb.NumBatch, msg.From, vLog.TxHash, fullBlock, msg.Nonce)
	if err != nil {
		return err
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedForceBatches = append(block.SequencedForceBatches, sequencedForceBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedForceBatches = append((*blocks)[len(*blocks)-1].SequencedForceBatches, sequencedForceBatch)
	} else {
		log.Error("Error processing ForceSequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing ForceSequencedBatches event")
	}
	or := Order{
		Name: SequenceForceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedForceBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)

	return nil
}

func decodeSequencedForceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, block *types.Block, nonce uint64) ([]SequencedForceBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}

	var forceBatches []etrogpolygonzkevm.PolygonRollupBaseEtrogBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &forceBatches)
	if err != nil {
		return nil, err
	}

	sequencedForcedBatches := make([]SequencedForceBatch, len(forceBatches))
	for i, force := range forceBatches {
		bn := lastBatchNumber - uint64(len(forceBatches)-(i+1))
		sequencedForcedBatches[i] = SequencedForceBatch{
			BatchNumber:                     bn,
			Coinbase:                        sequencer,
			TxHash:                          txHash,
			Timestamp:                       time.Unix(int64(block.Time()), 0),
			Nonce:                           nonce,
			PolygonRollupBaseEtrogBatchData: force,
		}
	}
	return sequencedForcedBatches, nil
}

func prepareBlock(vLog types.Log, t time.Time, fullBlock *types.Block) Block {
	var block Block
	block.BlockNumber = vLog.BlockNumber
	block.BlockHash = vLog.BlockHash
	block.ParentHash = fullBlock.ParentHash()
	block.ReceivedAt = t
	return block
}

func hash(data ...[32]byte) [32]byte {
	var res [32]byte
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		hash.Write(d[:]) //nolint:errcheck,gosec
	}
	copy(res[:], hash.Sum(nil))
	return res
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (etherMan *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return etherMan.EthClient.HeaderByNumber(ctx, number)
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number.
func (etherMan *Client) EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := etherMan.EthClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) || err.Error() == "block does not exist in blockchain" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return block, nil
}

// GetLatestBatchNumber function allows to retrieve the latest proposed batch in the smc
func (etherMan *Client) GetLatestBatchNumber() (uint64, error) {
	rollupData, err := etherMan.EtrogRollupManager.RollupIDToRollupData(&bind.CallOpts{Pending: false}, etherMan.RollupID)
	if err != nil {
		return 0, err
	}
	return rollupData.LastBatchSequenced, nil
}

// GetLatestBlockHeader gets the latest block header from the ethereum
func (etherMan *Client) GetLatestBlockHeader(ctx context.Context) (*types.Header, error) {
	header, err := etherMan.EthClient.HeaderByNumber(ctx, big.NewInt(int64(rpc.LatestBlockNumber)))
	if err != nil || header == nil {
		return nil, err
	}
	return header, nil
}

// GetLatestBlockNumber gets the latest block number from the ethereum
func (etherMan *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return etherMan.getBlockNumber(ctx, rpc.LatestBlockNumber)
}

// GetSafeBlockNumber gets the safe block number from the ethereum
func (etherMan *Client) GetSafeBlockNumber(ctx context.Context) (uint64, error) {
	return etherMan.getBlockNumber(ctx, rpc.SafeBlockNumber)
}

// GetFinalizedBlockNumber gets the Finalized block number from the ethereum
func (etherMan *Client) GetFinalizedBlockNumber(ctx context.Context) (uint64, error) {
	return etherMan.getBlockNumber(ctx, rpc.FinalizedBlockNumber)
}

// getBlockNumber gets the block header by the provided block number from the ethereum
func (etherMan *Client) getBlockNumber(ctx context.Context, blockNumber rpc.BlockNumber) (uint64, error) {
	header, err := etherMan.EthClient.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil || header == nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// GetLatestBlockTimestamp gets the latest block timestamp from the ethereum
func (etherMan *Client) GetLatestBlockTimestamp(ctx context.Context) (uint64, error) {
	header, err := etherMan.EthClient.HeaderByNumber(ctx, nil)
	if err != nil || header == nil {
		return 0, err
	}
	return header.Time, nil
}

// GetLatestVerifiedBatchNum gets latest verified batch from ethereum
func (etherMan *Client) GetLatestVerifiedBatchNum() (uint64, error) {
	rollupData, err := etherMan.EtrogRollupManager.RollupIDToRollupData(&bind.CallOpts{Pending: false}, etherMan.RollupID)
	if err != nil {
		return 0, err
	}
	return rollupData.LastVerifiedBatch, nil
}

// GetTx function get ethereum tx
func (etherMan *Client) GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	return etherMan.EthClient.TransactionByHash(ctx, txHash)
}

// GetTxReceipt function gets ethereum tx receipt
func (etherMan *Client) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return etherMan.EthClient.TransactionReceipt(ctx, txHash)
}

// ApprovePol function allow to approve tokens in pol smc
func (etherMan *Client) ApprovePol(ctx context.Context, account common.Address, polAmount *big.Int, to common.Address) (*types.Transaction, error) {
	opts, err := etherMan.getAuthByAddress(account)
	if err == ErrNotFound {
		return nil, errors.New("can't find account private key to sign tx")
	}
	if etherMan.GasProviders.MultiGasProvider {
		opts.GasPrice = etherMan.GetL1GasPrice(ctx)
	}
	tx, err := etherMan.Pol.Approve(&opts, etherMan.l1Cfg.ZkEVMAddr, polAmount)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		return nil, fmt.Errorf("error approving balance to send the batch. Error: %w", err)
	}

	return tx, nil
}

// GetTrustedSequencerURL Gets the trusted sequencer url from rollup smc
func (etherMan *Client) GetTrustedSequencerURL() (string, error) {
	return etherMan.EtrogZkEVM.TrustedSequencerURL(&bind.CallOpts{Pending: false})
}

// GetL2ChainID returns L2 Chain ID
func (etherMan *Client) GetL2ChainID() (uint64, error) {
	chainID, err := etherMan.PreEtrogZkEVM.ChainID(&bind.CallOpts{Pending: false})
	log.Debug("chainID read from preEtrogZkevm: ", chainID)
	if err != nil || chainID == 0 {
		log.Debug("error from preEtrogZkevm: ", err)
		rollupData, err := etherMan.EtrogRollupManager.RollupIDToRollupData(&bind.CallOpts{Pending: false}, etherMan.RollupID)
		log.Debugf("ChainID read from EtrogRollupManager: %d using rollupID: %d", rollupData.ChainID, etherMan.RollupID)
		if err != nil {
			log.Debug("error from EtrogRollupManager: ", err)
			return 0, err
		} else if rollupData.ChainID == 0 {
			return rollupData.ChainID, fmt.Errorf("error: chainID received is 0!!")
		}
		return rollupData.ChainID, nil
	}
	return chainID, nil
}

// GetL1GasPrice gets the l1 gas price
func (etherMan *Client) GetL1GasPrice(ctx context.Context) *big.Int {
	// Get gasPrice from providers
	gasPrice := big.NewInt(0)
	for i, prov := range etherMan.GasProviders.Providers {
		gp, err := prov.SuggestGasPrice(ctx)
		if err != nil {
			log.Warnf("error getting gas price from provider %d. Error: %s", i+1, err.Error())
		} else if gasPrice.Cmp(gp) == -1 { // gasPrice < gp
			gasPrice = gp
		}
	}
	log.Debug("gasPrice chose: ", gasPrice)
	return gasPrice
}

// SendTx sends a tx to L1
func (etherMan *Client) SendTx(ctx context.Context, tx *types.Transaction) error {
	return etherMan.EthClient.SendTransaction(ctx, tx)
}

// PendingNonce returns the pending nonce for the provided account
func (etherMan *Client) PendingNonce(ctx context.Context, account common.Address) (uint64, error) {
	return etherMan.EthClient.PendingNonceAt(ctx, account)
}

// CurrentNonce returns the current nonce for the provided account
func (etherMan *Client) CurrentNonce(ctx context.Context, account common.Address) (uint64, error) {
	return etherMan.EthClient.NonceAt(ctx, account, nil)
}

// SuggestedGasPrice returns the suggest nonce for the network at the moment
func (etherMan *Client) SuggestedGasPrice(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice := etherMan.GetL1GasPrice(ctx)
	if suggestedGasPrice.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("failed to get the suggested gas price")
	}
	return suggestedGasPrice, nil
}

// EstimateGas returns the estimated gas for the tx
func (etherMan *Client) EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error) {
	return etherMan.EthClient.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	})
}

// DepositCount returns deposits count
func (etherman *Client) DepositCount(ctx context.Context, blockNumber *uint64) (*big.Int, error) {
	var opts *bind.CallOpts
	if blockNumber != nil {
		opts = new(bind.CallOpts)
		opts.BlockNumber = new(big.Int).SetUint64(*blockNumber)
	}

	return etherman.EtrogGlobalExitRootManager.DepositCount(opts)
}

// CheckTxWasMined check if a tx was already mined
func (etherMan *Client) CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error) {
	receipt, err := etherMan.EthClient.TransactionReceipt(ctx, txHash)
	if errors.Is(err, ethereum.NotFound) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}

	return true, receipt, nil
}

// SignTx tries to sign a transaction accordingly to the provided sender
func (etherMan *Client) SignTx(ctx context.Context, sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	auth, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, ErrPrivateKeyNotFound
	}
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// GetRevertMessage tries to get a revert message of a transaction
func (etherMan *Client) GetRevertMessage(ctx context.Context, tx *types.Transaction) (string, error) {
	if tx == nil {
		return "", nil
	}

	receipt, err := etherMan.GetTxReceipt(ctx, tx.Hash())
	if err != nil {
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		revertMessage, err := operations.RevertReason(ctx, etherMan.EthClient, tx, receipt.BlockNumber)
		if err != nil {
			return "", err
		}
		return revertMessage, nil
	}
	return "", nil
}

// AddOrReplaceAuth adds an authorization or replace an existent one to the same account
func (etherMan *Client) AddOrReplaceAuth(auth bind.TransactOpts) error {
	log.Infof("added or replaced authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return nil
}

// LoadAuthFromKeyStore loads an authorization from a key store file
func (etherMan *Client) LoadAuthFromKeyStore(path, password string) (*bind.TransactOpts, error) {
	auth, err := newAuthFromKeystore(path, password, etherMan.l1Cfg.L1ChainID)
	if err != nil {
		return nil, err
	}

	log.Infof("loaded authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return &auth, nil
}

// newKeyFromKeystore creates an instance of a keystore key from a keystore file
func newKeyFromKeystore(path, password string) (*keystore.Key, error) {
	if path == "" && password == "" {
		return nil, nil
	}
	keystoreEncrypted, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	log.Infof("decrypting key from: %v", path)
	key, err := keystore.DecryptKey(keystoreEncrypted, password)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// newAuthFromKeystore an authorization instance from a keystore file
func newAuthFromKeystore(path, password string, chainID uint64) (bind.TransactOpts, error) {
	log.Infof("reading key from: %v", path)
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		return bind.TransactOpts{}, err
	}
	if key == nil {
		return bind.TransactOpts{}, nil
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		return bind.TransactOpts{}, err
	}
	return *auth, nil
}

// getAuthByAddress tries to get an authorization from the authorizations map
func (etherMan *Client) getAuthByAddress(addr common.Address) (bind.TransactOpts, error) {
	auth, found := etherMan.auth[addr]
	if !found {
		return bind.TransactOpts{}, ErrNotFound
	}
	return auth, nil
}

// generateRandomAuth generates an authorization instance from a
// randomly generated private key to be used to estimate gas for PoE
// operations NOT restricted to the Trusted Sequencer
func (etherMan *Client) generateRandomAuth() (bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a private key to estimate L1 txs")
	}
	chainID := big.NewInt(0).SetUint64(etherMan.l1Cfg.L1ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a fake authorization to estimate L1 txs")
	}

	return *auth, nil
}
