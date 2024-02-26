package main

import (
	"crypto/ecdsa"
	"fmt"

	dataCommitteeClient "github.com/0xPolygon/cdk-data-availability/client"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/apollo"
	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
	"github.com/0xPolygonHermez/zkevm-node/dataavailability/datacommittee"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencesender"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func createSequenceSenderX1(cfg config.Config, pool *pool.Pool, etmStorage *ethtxmanager.PostgresStorage, st *state.State, eventLog *event.EventLog) *sequencesender.SequenceSender {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	da := setEthermanDaX1(cfg, st, etherman, true)

	_, privKey, err := etherman.LoadAuthFromKeyStoreX1(cfg.SequenceSender.DAPermitApiPrivateKey.Path, cfg.SequenceSender.DAPermitApiPrivateKey.Password)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.SequenceSender.SenderAddress.Cmp(common.Address{}) == 0 {
		log.Fatal("Sequence sender address not found")
	}
	if privKey == nil { //nolint:staticcheck
		log.Fatal("DA permit api private key not found")
	}
	log.Infof("from pk %s, from sender %s", crypto.PubkeyToAddress(privKey.PublicKey), cfg.SequenceSender.SenderAddress.String()) //nolint:staticcheck

	cfg.SequenceSender.ForkUpgradeBatchNumber = cfg.ForkUpgradeBatchNumber

	ethTxManager := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)

	seqSender, err := sequencesender.New(cfg.SequenceSender, st, etherman, ethTxManager, eventLog)
	if err != nil {
		log.Fatal(err)
	}
	seqSender.SetDataProvider(da)
	return seqSender
}

func initRunForX1(c *config.Config, components []string) {
	// Read configure from apollo
	apolloClient := apollo.NewClient(c)
	if apolloClient.LoadConfig() {
		log.Info("apollo config loaded")
	}

	pool.SetL2BridgeAddr(c.NetworkConfig.L2BridgeAddr)

	for _, component := range components {
		switch component {
		case RPC:
			jsonrpc.InitRateLimit(c.RPC.RateLimit)
		}
	}
}

func newDataAvailability(c config.Config, st *state.State, etherman *etherman.Client, isSequenceSender bool) (*dataavailability.DataAvailability, error) {
	var (
		trustedSequencerURL string
		err                 error
	)
	if !c.IsTrustedSequencer {
		if c.Synchronizer.TrustedSequencerURL != "" {
			trustedSequencerURL = c.Synchronizer.TrustedSequencerURL
		} else {
			log.Debug("getting trusted sequencer URL from smc")
			trustedSequencerURL, err = etherman.GetTrustedSequencerURL()
			if err != nil {
				return nil, fmt.Errorf("error getting trusted sequencer URI. Error: %v", err)
			}
		}
		log.Debug("trustedSequencerURL ", trustedSequencerURL)
	}
	zkEVMClient := client.NewClient(trustedSequencerURL)

	// Backend specific config
	daProtocolName, err := etherman.GetDAProtocolName()
	if err != nil {
		return nil, fmt.Errorf("error getting data availability protocol name: %v", err)
	}
	var daBackend dataavailability.DABackender
	switch daProtocolName {
	case string(dataavailability.DataAvailabilityCommittee):
		var (
			pk  *ecdsa.PrivateKey
			err error
		)
		if isSequenceSender {
			_, pk, err = etherman.LoadAuthFromKeyStoreX1(c.SequenceSender.DAPermitApiPrivateKey.Path, c.SequenceSender.DAPermitApiPrivateKey.Password)
			if err != nil {
				return nil, err
			}
			log.Infof("from pk %s", crypto.PubkeyToAddress(pk.PublicKey))
		}
		dacAddr, err := etherman.GetDAProtocolAddr()
		if err != nil {
			return nil, fmt.Errorf("error getting trusted sequencer URI. Error: %v", err)
		}
		daBackend, err = datacommittee.New(
			c.Etherman.URL,
			dacAddr,
			pk,
			&dataCommitteeClient.Factory{},
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected / unsupported DA protocol: %s", daProtocolName)
	}

	return dataavailability.New(
		c.IsTrustedSequencer,
		daBackend,
		st,
		zkEVMClient,
	)
}

func setEthermanDaX1(c config.Config, st *state.State, etherman *etherman.Client, isSequenceSender bool) *dataavailability.DataAvailability {
	da, err := newDataAvailability(c, st, etherman, isSequenceSender)
	if err != nil {
		log.Fatal(err)
	}
	etherman.SetDataProvider(da)
	return da
}
