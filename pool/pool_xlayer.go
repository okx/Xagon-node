package pool

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// BridgeClaimMethodSignature for tracking BridgeClaimMethodSignature method
	BridgeClaimMethodSignature = "0xccaa2d11"
	// BridgeClaimMessageMethodSignature for tracking BridgeClaimMethodSignature method
	BridgeClaimMessageMethodSignature = "0xf5efcd79"
	//ExWithdrawalMethodSignature erc20 contract transfer(address recipient, uint256 amount)
	ExWithdrawalMethodSignature = "0xa9059cbb"

	// TestnetChainID the chain id of xlayer testnet
	TestnetChainID = 195
	// MainnetChainID the chain id of xlayer mainnet
	MainnetChainID = 196
	// TestnetBridgeURL the bridge url od xlayer testnet
	TestnetBridgeURL = "https://www.okx.com/xlayer/bridge-test"
	// MainnetBridgeURL the bridge url od xlayer mainnet
	MainnetBridgeURL = "https://www.okx.com/xlayer/bridge"
)

// Contains returns if string[] contains ele string
func Contains(s []string, ele common.Address) bool {
	for _, e := range s {
		if common.HexToAddress(e) == ele {
			return true
		}
	}
	return false
}

// ContainsMethod returns if data has prefix of method sig
func ContainsMethod(data string, methods []string) bool {
	for _, m := range methods {
		if strings.HasPrefix(data, m) {
			return true
		}
	}
	return false
}

// StartRefreshingWhiteAddressesPeriodically will make this instance of the pool
// to check periodically(accordingly to the configuration) for updates regarding
// the white address and update the in memory blocked addresses
func (p *Pool) StartRefreshingWhiteAddressesPeriodically() {
	interval := p.cfg.IntervalToRefreshWhiteAddresses.Duration
	if interval.Nanoseconds() <= 0 {
		interval = 20 * time.Second //nolint:gomnd
	}

	p.refreshWhitelistedAddresses()
	go func(p *Pool) {
		for {
			time.Sleep(interval)
			p.refreshWhitelistedAddresses()
		}
	}(p)
}

// refreshWhitelistedAddresses refreshes the list of whitelisted addresses for the provided instance of pool
func (p *Pool) refreshWhitelistedAddresses() {
	whitelistedAddresses, err := p.storage.GetAllAddressesWhitelisted(context.Background())
	if err != nil {
		log.Error("failed to load whitelisted addresses")
		return
	}

	whitelistedAddressesMap := sync.Map{}
	for _, whitelistedAddress := range whitelistedAddresses {
		whitelistedAddressesMap.Store(whitelistedAddress.String(), 1)
		p.whitelistedAddresses.Store(whitelistedAddress.String(), 1)
	}

	nonWhitelistedAddresses := []string{}
	p.whitelistedAddresses.Range(func(key, value any) bool {
		addrHex := key.(string)
		_, found := whitelistedAddressesMap.Load(addrHex)
		if found {
			return true
		}

		nonWhitelistedAddresses = append(nonWhitelistedAddresses, addrHex)
		return true
	})

	for _, nonWhitelistedAddress := range nonWhitelistedAddresses {
		p.whitelistedAddresses.Delete(nonWhitelistedAddress)
	}
}

// GetMinSuggestedGasPriceWithDelta gets the min suggested gas price
func (p *Pool) GetMinSuggestedGasPriceWithDelta(ctx context.Context, delta time.Duration) (uint64, error) {
	fromTimestamp := time.Now().UTC().Add(-p.cfg.MinAllowedGasPriceInterval.Duration)
	fromTimestamp = fromTimestamp.Add(delta)
	if fromTimestamp.Before(p.startTimestamp) {
		fromTimestamp = p.startTimestamp
	}

	return p.storage.MinL2GasPriceSince(ctx, fromTimestamp)
}

// GetDynamicGasPrice returns the current L2 dynamic gas price
func (p *Pool) GetDynamicGasPrice() *big.Int {
	p.dgpMux.RLock()
	dgp := p.dynamicGasPrice
	p.dgpMux.RUnlock()
	if dgp == nil || dgp.Cmp(big.NewInt(0)) == 0 {
		_, l2Gp := p.GetL1AndL2GasPrice()
		dgp = new(big.Int).SetUint64(l2Gp)
	}
	return dgp
}

func (p *Pool) checkFreeGp(ctx context.Context, poolTx Transaction, from common.Address) (bool, error) {
	// claim tx
	if isFreeGasAddress(p.cfg.FreeGasAddress, from) && poolTx.IsClaims {
		return true, nil
	}

	// special project
	if GetEnableSpecialFreeGasList(p.cfg.EnableFreeGasList) {
		fromToName, freeGpList := GetSpecialFreeGasList(p.cfg.FreeGasList)
		info := freeGpList[fromToName[from.String()]]
		if info != nil &&
			Contains(info.ToList, *poolTx.To()) &&
			ContainsMethod("0x"+common.Bytes2Hex(poolTx.Data()), info.MethodSigs) {
			return true, nil
		}
	}

	// new bridge address
	if getEnableFreeGasByNonce(p.cfg.EnableFreeGasByNonce) && poolTx.GasPrice().Cmp(big.NewInt(0)) == 0 { // free-gas tx by count
		isFreeAddr, err := p.storage.IsFreeGasAddr(ctx, from)
		if err != nil {
			log.Errorf("failed to check free gas address from storage: %v", err)
			return false, err
		}

		freeGasCountPerAddrConfig := getFreeGasCountPerAddr(p.cfg.FreeGasCountPerAddr)
		if isFreeAddr {
			if poolTx.Nonce() < freeGasCountPerAddrConfig {
				if poolTx.Gas() > getFreeGasLimit(p.cfg.FreeGasLimit) {
					return false, fmt.Errorf("gas-free transaction with too high gas limit")
				}
				return true, nil
			} else {
				return false, fmt.Errorf("You are no longer eligible for gas-free transactions because for each new address. Only the first %d transactions(address nonce less than %d) can be gas-free",
					freeGasCountPerAddrConfig,
					freeGasCountPerAddrConfig)
			}
		} else {
			bridgeURL := ""
			switch p.chainID {
			case TestnetChainID:
				bridgeURL = TestnetBridgeURL
			case MainnetChainID:
				bridgeURL = MainnetBridgeURL
			}
			return false, fmt.Errorf("You are unable to initiate a gas-free transaction from this address unless you have previously transferred funds to this address via the X Layer Bridge (%s) or the OKX Exchange. Only the first %d transactions (address nonce must be less than %d) can be gas-free",
				bridgeURL,
				freeGasCountPerAddrConfig,
				freeGasCountPerAddrConfig)
		}
	}
	return false, nil
}

func (p *Pool) checkAndUpdateFreeGasAddr(ctx context.Context, poolTx Transaction, from common.Address, root common.Hash) error {
	// check and store the free gas address
	var freeGpAddr common.Address
	inputHex := hex.EncodeToHex(poolTx.Data())
	// hard code
	if isFreeGasExAddress(p.cfg.FreeGasExAddress, from) {
		if strings.HasPrefix(inputHex, ExWithdrawalMethodSignature) && len(inputHex) > 74 { // erc20 contract transfer
			addrHex := "0x" + inputHex[10:74]
			freeGpAddr = common.HexToAddress(addrHex)
		} else {
			// the to address of any Ex withdrawal okb tx will be considered as a gas-free address
			// even if this address is a contract, it will not affect the gas-free
			freeGpAddr = *poolTx.To()
		}
	} else if poolTx.IsClaims && len(inputHex) > 4554 { // bridge contract claim
		addrHex := "0x" + inputHex[4490:4554]
		freeGpAddr = common.HexToAddress(addrHex)
	}

	if freeGpAddr.Cmp(common.Address{}) != 0 {
		nonce, err := p.state.GetNonce(ctx, freeGpAddr, root)
		if err != nil {
			log.Errorf("failed to get nonce while adding tx to the pool", err)
			return err
		}
		if nonce < getFreeGasCountPerAddr(p.cfg.FreeGasCountPerAddr) {
			if err = p.storage.AddFreeGasAddr(ctx, freeGpAddr); err != nil {
				log.Errorf("failed to save free gas address to the storage", err)
				return err
			}
		}
	}
	return nil
}

// AddDynamicGp cache the dynamic gas price of L2
func (p *Pool) AddDynamicGp(dgp *big.Int) {
	_, l2Gp := p.GetL1AndL2GasPrice()
	result := new(big.Int).SetUint64(l2Gp)
	if result.Cmp(dgp) < 0 {
		result = new(big.Int).Set(dgp)
	}
	p.dgpMux.Lock()
	p.dynamicGasPrice = result
	p.dgpMux.Unlock()
}

func (p *Pool) checkBlockedAddr(address common.Address) bool {
	// check from db
	if _, blocked := p.blockedAddresses.Load(address.String()); blocked {
		return true
	}

	// check from dynamic config
	return isBlockedAddress(p.cfg.BlockedList, address)
}
