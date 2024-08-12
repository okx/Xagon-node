package pool

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// apolloConfig is the apollo pool dynamic config
type apolloConfig struct {
	EnableApollo       bool
	FreeGasAddresses   []string
	GlobalQueue        uint64
	AccountQueue       uint64
	EnableWhitelist    bool
	BridgeClaimMethods []string
	EnablePendingStat  bool

	EnableFreeGasByNonce bool
	FreeGasExAddress     []string
	FreeGasCountPerAddr  uint64
	FreeGasLimit         uint64

	// for special project
	EnableFreeGasList  bool
	FreeGasFromNameMap map[string]string       // map[from]projectName
	FreeGasList        map[string]*FreeGasInfo // map[projectName]FreeGasInfo

	BlockedList []string

	sync.RWMutex
}

var apolloConf = &apolloConfig{}

// getApolloConfig returns the pool singleton instance
func getApolloConfig() *apolloConfig {
	return apolloConf
}

// Enable returns true if apollo is enabled
func (c *apolloConfig) enable() bool {
	if c == nil || !c.EnableApollo {
		return false
	}
	c.RLock()
	defer c.RUnlock()
	return c.EnableApollo
}

func (c *apolloConfig) setFreeGasList(freeGasList []FreeGasInfo) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.FreeGasFromNameMap = make(map[string]string)
	c.FreeGasList = make(map[string]*FreeGasInfo, len(freeGasList))
	for _, info := range freeGasList {
		for _, from := range info.FromList {
			c.FreeGasFromNameMap[from] = info.Name
		}
		infoCopy := info
		c.FreeGasList[info.Name] = &infoCopy
	}
}

func (c *apolloConfig) setFreeGasAddresses(freeGasAddrs []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.FreeGasAddresses = make([]string, len(freeGasAddrs))
	copy(c.FreeGasAddresses, freeGasAddrs)
}

func (c *apolloConfig) setBlockedList(blockedAddrs []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.BlockedList = make([]string, len(blockedAddrs))
	copy(c.BlockedList, blockedAddrs)
}

func (c *apolloConfig) setFreeGasExAddresses(freeGasExAddrs []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.FreeGasExAddress = make([]string, len(freeGasExAddrs))
	copy(c.FreeGasExAddress, freeGasExAddrs)
}

func (c *apolloConfig) setBridgeClaimMethods(bridgeClaimMethods []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.BridgeClaimMethods = make([]string, len(bridgeClaimMethods))
	copy(c.BridgeClaimMethods, bridgeClaimMethods)
}

// UpdateConfig updates the apollo config
// GlobalQueue
// AccountQueue
// FreeGasAddress
// EnableWhitelist
// EnablePendingStat
func UpdateConfig(apolloConfig Config) {
	getApolloConfig().Lock()
	getApolloConfig().EnableApollo = true
	getApolloConfig().GlobalQueue = apolloConfig.GlobalQueue
	getApolloConfig().AccountQueue = apolloConfig.AccountQueue
	getApolloConfig().setFreeGasAddresses(apolloConfig.FreeGasAddress)
	getApolloConfig().setBlockedList(apolloConfig.BlockedList)
	getApolloConfig().EnableWhitelist = apolloConfig.EnableWhitelist
	getApolloConfig().setBridgeClaimMethods(apolloConfig.BridgeClaimMethodSigs)

	getApolloConfig().EnableFreeGasByNonce = apolloConfig.EnableFreeGasByNonce
	getApolloConfig().setFreeGasExAddresses(apolloConfig.FreeGasExAddress)
	getApolloConfig().FreeGasCountPerAddr = apolloConfig.FreeGasCountPerAddr
	getApolloConfig().FreeGasLimit = apolloConfig.FreeGasLimit
	getApolloConfig().EnableFreeGasList = apolloConfig.EnableFreeGasList
	getApolloConfig().setFreeGasList(apolloConfig.FreeGasList)

	getApolloConfig().Unlock()
}

func getClaimMethod(localBridgeClaimMethods []string) []string {
	var methods []string
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		methods = getApolloConfig().BridgeClaimMethods
	} else {
		methods = localBridgeClaimMethods
	}
	if len(methods) == 0 {
		methods = append(methods, BridgeClaimMethodSignature, BridgeClaimMessageMethodSignature)
	}

	return methods
}

func isFreeGasAddress(localFreeGasAddrs []string, address common.Address) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return Contains(getApolloConfig().FreeGasAddresses, address)
	}

	return Contains(localFreeGasAddrs, address)
}

func getEnableFreeGasByNonce(enableFreeGasByNonce bool) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().EnableFreeGasByNonce
	}

	return enableFreeGasByNonce
}

func isFreeGasExAddress(localFreeGasExAddrs []string, address common.Address) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return Contains(getApolloConfig().FreeGasExAddress, address)
	}

	return Contains(localFreeGasExAddrs, address)
}

func getFreeGasCountPerAddr(localFreeGasCountPerAddr uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().FreeGasCountPerAddr
	}
	return localFreeGasCountPerAddr
}

func getFreeGasLimit(localFreeGasLimit uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().FreeGasLimit
	}
	return localFreeGasLimit
}

func getGlobalQueue(globalQueue uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().GlobalQueue
	}

	return globalQueue
}

func getAccountQueue(accountQueue uint64) uint64 {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().AccountQueue
	}

	return accountQueue
}

func getEnableWhitelist(enableWhitelist bool) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().EnableWhitelist
	}

	return enableWhitelist
}

func isBlockedAddress(localBlockedList []string, address common.Address) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return Contains(getApolloConfig().BlockedList, address)
	}

	return Contains(localBlockedList, address)
}

// GetEnableSpecialFreeGasList returns enable flag of FreeGasList
func GetEnableSpecialFreeGasList(enableFreeGasList bool) bool {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().EnableFreeGasList
	}
	return enableFreeGasList
}

// GetSpecialFreeGasList returns the special project in XLayer for free gas
func GetSpecialFreeGasList(freeGasList []FreeGasInfo) (map[string]string, map[string]*FreeGasInfo) {
	if getApolloConfig().enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().FreeGasFromNameMap, getApolloConfig().FreeGasList
	}

	freeGasFromNameMap := make(map[string]string)
	freeGasMap := make(map[string]*FreeGasInfo, len(freeGasList))
	for _, info := range freeGasList {
		for _, from := range info.FromList {
			freeGasFromNameMap[from] = info.Name
		}
		infoCopy := info
		freeGasMap[info.Name] = &infoCopy
	}

	return freeGasFromNameMap, freeGasMap
}
