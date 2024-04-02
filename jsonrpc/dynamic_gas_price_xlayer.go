package jsonrpc

import (
	"context"
	"math/big"
	"sort"
	"sync"
	"time"

	zktypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/core/types"
)

// DynamicGPConfig represents the configuration of the dynamic gas price
type DynamicGPConfig struct {

	// Enabled defines if the dynamic gas price is enabled or disabled
	Enabled bool `mapstructure:"Enabled"`

	// CongestionTxThreshold defines the tx threshold to measure whether there is congestion
	CongestionTxThreshold uint64 `mapstructure:"CongestionTxThreshold"`

	// CheckBatches defines the number of recent Batches used to sample gas price
	CheckBatches int `mapstructure:"CheckBatches"`

	// SampleTxNumer defines the number of sampled gas prices in each batch
	SampleNumber int `mapstructure:"SampleNumber"`

	// Percentile defines the sampling weight of all sampled gas prices
	Percentile int `mapstructure:"Percentile"`

	// MaxPrice defines the dynamic gas price upper limit
	MaxPrice uint64 `mapstructure:"MaxPrice"`

	// MinPrice defines the dynamic gas price lower limit
	MinPrice uint64 `mapstructure:"MinPrice"`

	//UpdatePeriod defines the time interval for updating dynamic gas price
	UpdatePeriod zktypes.Duration `mapstructure:"UpdatePeriod"`
}

// DynamicGPManager allows to update recommended gas price
type DynamicGPManager struct {
	lastL2BatchNumber uint64
	lastPrice         *big.Int
	cacheLock         sync.RWMutex
	fetchLock         sync.Mutex
}

// runDynamicSuggester init the routine for dynamic gas price updates
func (e *EthEndpoints) runDynamicGPSuggester() {
	ctx := context.Background()
	// initialization
	updateTimer := time.NewTimer(10 * time.Second) //nolint:gomnd
	for {
		select {
		case <-ctx.Done():
			log.Info("Finishing dynamic gas price suggester...")
			return
		case <-updateTimer.C:
			if getApolloConfig().Enable() {
				getApolloConfig().RLock()
				e.cfg.DynamicGP = getApolloConfig().DynamicGP
				getApolloConfig().RUnlock()
			}
			log.Info("Dynamic gas price update period is ", e.cfg.DynamicGP.UpdatePeriod.Duration.String())
			e.calcDynamicGP(ctx)
			updateTimer.Reset(e.cfg.DynamicGP.UpdatePeriod.Duration)
		}
	}
}

func (e *EthEndpoints) calcDynamicGP(ctx context.Context) {
	l2BatchNumber, err := e.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last l2 batch number, err: %v", err)
		return
	}

	e.dgpMan.cacheLock.RLock()
	lastL2BatchNumber, lastPrice := e.dgpMan.lastL2BatchNumber, new(big.Int).Set(e.dgpMan.lastPrice)
	e.dgpMan.cacheLock.RUnlock()
	if l2BatchNumber == lastL2BatchNumber {
		log.Debug("Batch is still the same, no need to update the gas price at the moment, lastL2BatchNumber: ", lastL2BatchNumber)
		return
	}

	e.dgpMan.fetchLock.Lock()
	defer e.dgpMan.fetchLock.Unlock()

	var (
		sent, exp int
		number    = lastL2BatchNumber
		result    = make(chan results, e.cfg.DynamicGP.CheckBatches)
		quit      = make(chan struct{})
		results   []*big.Int
	)

	for sent < e.cfg.DynamicGP.CheckBatches && number > 0 {
		go e.getL2BatchTxsTips(ctx, number, e.cfg.DynamicGP.SampleNumber, result, quit)
		sent++
		exp++
		number--
	}

	for exp > 0 {
		res := <-result
		if res.err != nil {
			close(quit)
			return
		}
		exp--

		if len(res.values) == 0 {
			res.values = []*big.Int{lastPrice}
		}
		results = append(results, res.values...)
	}

	price := lastPrice
	if len(results) > 0 {
		sort.Sort(bigIntArray(results))
		price = results[(len(results)-1)*e.cfg.DynamicGP.Percentile/100]
	}

	minGasPrice := big.NewInt(0).SetUint64(e.cfg.DynamicGP.MinPrice)
	if minGasPrice.Cmp(price) == 1 {
		price = minGasPrice
	}

	maxGasPrice := new(big.Int).SetUint64(e.cfg.DynamicGP.MaxPrice)
	if e.cfg.DynamicGP.MaxPrice > 0 && price.Cmp(maxGasPrice) == 1 {
		price = maxGasPrice
	}

	// judge if there is congestion
	isCongested, err := e.isCongested(ctx)
	if err != nil {
		log.Errorf("failed to count pool txs by status pending while judging if the pool is congested: ", err)
		return
	}

	if !isCongested {
		gasPrices, err := e.pool.GetGasPrices(ctx)
		if err != nil {
			log.Errorf("failed to get raw gas prices when it is not congested: ", err)
			return
		}
		rawGasPrice := new(big.Int).SetUint64(gasPrices.L2GasPrice)
		e.dgpMan.cacheLock.Lock()
		e.dgpMan.lastPrice = getAvgPrice(rawGasPrice, price)
		e.dgpMan.lastL2BatchNumber = l2BatchNumber
		e.dgpMan.cacheLock.Unlock()
		return
	}

	log.Debug("there is congestion for L2")
	e.dgpMan.cacheLock.Lock()
	e.dgpMan.lastPrice = price
	e.dgpMan.lastL2BatchNumber = l2BatchNumber
	e.dgpMan.cacheLock.Unlock()
}

func (e *EthEndpoints) getL2BatchTxsTips(ctx context.Context, l2BlockNumber uint64, limit int, result chan results, quit chan struct{}) {
	txs, _, err := e.state.GetTransactionsByBatchNumber(ctx, l2BlockNumber, nil)
	if txs == nil {
		select {
		case result <- results{nil, err}:
		case <-quit:
		}
		return
	}
	sorter := newSorter(txs)
	sort.Sort(sorter)

	var prices []*big.Int
	for _, tx := range sorter.txs {
		tip := tx.GasTipCap()
		prices = append(prices, tip)
		if len(prices) >= limit {
			break
		}
	}

	select {
	case result <- results{prices, nil}:
	case <-quit:
	}
}

func (e *EthEndpoints) isCongested(ctx context.Context) (bool, error) {
	txCount, err := e.pool.CountPendingTransactions(ctx)
	if err != nil {
		return false, err
	}
	if txCount >= e.cfg.DynamicGP.CongestionTxThreshold {
		return true, nil
	}
	return false, nil
}

type results struct {
	values []*big.Int
	err    error
}

type txSorter struct {
	txs []types.Transaction
}

func newSorter(txs []types.Transaction) *txSorter {
	return &txSorter{
		txs: txs,
	}
}

func (s *txSorter) Len() int { return len(s.txs) }
func (s *txSorter) Swap(i, j int) {
	s.txs[i], s.txs[j] = s.txs[j], s.txs[i]
}
func (s *txSorter) Less(i, j int) bool {
	tip1 := s.txs[i].GasTipCap()
	tip2 := s.txs[j].GasTipCap()
	return tip1.Cmp(tip2) < 0
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func getAvgPrice(low *big.Int, high *big.Int) *big.Int {
	avg := new(big.Int).Add(low, high)
	avg = avg.Quo(avg, big.NewInt(2)) //nolint:gomnd
	return avg
}
