package pool

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")

	// ErrAlreadyKnown is returned if the transactions is already contained
	// within the pool.
	ErrAlreadyKnown = errors.New("already known")

	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")

	// ErrEffectiveGasPriceGasPriceTooLow the tx gas price is lower than breakEvenGasPrice and lower than L2GasPrice
	ErrEffectiveGasPriceGasPriceTooLow = errors.New("effective gas price: gas price too low")
)

// Pool is an implementation of the Pool interface
// that uses a postgres database to store the data
type Pool struct {
	storage
	state                   stateInterface
	chainID                 uint64
	cfg                     Config
	batchConstraintsCfg     state.BatchConstraintsCfg
	blockedAddresses        sync.Map
	minSuggestedGasPrice    *big.Int
	minSuggestedGasPriceMux *sync.RWMutex
	eventLog                *event.EventLog
	startTimestamp          time.Time
	gasPrices               GasPrices
	gasPricesMux            *sync.RWMutex
	effectiveGasPrice       *EffectiveGasPrice
}

type preExecutionResponse struct {
	usedZKCounters       state.ZKCounters
	reservedZKCounters   state.ZKCounters
	isExecutorLevelError bool
	OOCError             error
	OOGError             error
	isReverted           bool
	txResponse           *state.ProcessTransactionResponse
}

// GasPrices contains the gas prices for L2 and L1
type GasPrices struct {
	L2GasPrice uint64
	L1GasPrice uint64
}

// NewPool creates and initializes an instance of Pool
func NewPool(cfg Config, batchConstraintsCfg state.BatchConstraintsCfg, s storage, st stateInterface, chainID uint64, eventLog *event.EventLog) *Pool {
	startTimestamp := time.Now()
	p := &Pool{
		cfg:                     cfg,
		batchConstraintsCfg:     batchConstraintsCfg,
		startTimestamp:          startTimestamp,
		storage:                 s,
		state:                   st,
		chainID:                 chainID,
		blockedAddresses:        sync.Map{},
		minSuggestedGasPriceMux: new(sync.RWMutex),
		minSuggestedGasPrice:    big.NewInt(int64(cfg.DefaultMinGasPriceAllowed)),
		eventLog:                eventLog,
		gasPrices:               GasPrices{0, 0},
		gasPricesMux:            new(sync.RWMutex),
		effectiveGasPrice:       NewEffectiveGasPrice(cfg.EffectiveGasPrice),
	}
	p.refreshGasPrices()
	go func(cfg *Config, p *Pool) {
		for {
			p.refreshGasPrices()
			time.Sleep(cfg.IntervalToRefreshGasPrices.Duration)
		}
	}(&cfg, p)

	return p
}

// refresGasPRices refreshes the gas price
func (p *Pool) refreshGasPrices() {
	gasPrices, err := p.GetGasPrices(context.Background())
	if err != nil {
		log.Error("failed to load gas prices")
		return
	}

	p.gasPricesMux.Lock()
	p.gasPrices = gasPrices
	p.gasPricesMux.Unlock()
}

// StartRefreshingBlockedAddressesPeriodically will make this instance of the pool
// to check periodically(accordingly to the configuration) for updates regarding
// the blocked address and update the in memory blocked addresses
func (p *Pool) StartRefreshingBlockedAddressesPeriodically() {
	p.refreshBlockedAddresses()
	go func(p *Pool) {
		for {
			time.Sleep(p.cfg.IntervalToRefreshBlockedAddresses.Duration)
			p.refreshBlockedAddresses()
		}
	}(p)
}

// refreshBlockedAddresses refreshes the list of blocked addresses for the provided instance of pool
func (p *Pool) refreshBlockedAddresses() {
	blockedAddresses, err := p.storage.GetAllAddressesBlocked(context.Background())
	if err != nil {
		log.Error("failed to load blocked addresses")
		return
	}

	blockedAddressesMap := sync.Map{}
	for _, blockedAddress := range blockedAddresses {
		blockedAddressesMap.Store(blockedAddress.String(), 1)
		p.blockedAddresses.Store(blockedAddress.String(), 1)
	}

	unblockedAddresses := []string{}
	p.blockedAddresses.Range(func(key, value any) bool {
		addrHex := key.(string)
		_, found := blockedAddressesMap.Load(addrHex)
		if found {
			return true
		}

		unblockedAddresses = append(unblockedAddresses, addrHex)
		return true
	})

	for _, unblockedAddress := range unblockedAddresses {
		p.blockedAddresses.Delete(unblockedAddress)
	}
}

// StartPollingMinSuggestedGasPrice starts polling the minimum suggested gas price
func (p *Pool) StartPollingMinSuggestedGasPrice(ctx context.Context) {
	p.tryUpdateMinSuggestedGasPrice(p.cfg.DefaultMinGasPriceAllowed)
	p.pollMinSuggestedGasPrice(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.cfg.PollMinAllowedGasPriceInterval.Duration):
				p.pollMinSuggestedGasPrice(ctx)
			}
		}
	}()
}

// AddTx adds a transaction to the pool with the pending state
func (p *Pool) AddTx(ctx context.Context, tx types.Transaction, ip string) error {
	poolTx := NewTransaction(tx, ip, false)
	if err := p.validateTx(ctx, *poolTx); err != nil {
		return err
	}

	return p.StoreTx(ctx, tx, ip, false)
}

// StoreTx adds a transaction to the pool with the pending state
func (p *Pool) StoreTx(ctx context.Context, tx types.Transaction, ip string, isWIP bool) error {
	// Execute transaction to calculate its zkCounters
	preExecutionResponse, err := p.preExecuteTx(ctx, tx)
	if errors.Is(err, runtime.ErrIntrinsicInvalidBatchGasLimit) {
		return ErrGasLimit
	} else if preExecutionResponse.isExecutorLevelError {
		// Do not add tx to the pool
		return err
	} else if err != nil {
		log.Errorf("Pre execution error: %v", err)
		return err
	}

	if preExecutionResponse.OOCError != nil {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			IPAddress:   ip,
			Source:      event.Source_Node,
			Component:   event.Component_Pool,
			Level:       event.Level_Warning,
			EventID:     event.EventID_PreexecutionOOC,
			Description: tx.Hash().String(),
		}

		err := p.eventLog.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error adding event: %v", err)
		}
		// Do not add tx to the pool
		return fmt.Errorf("failed to add tx to the pool: %w", preExecutionResponse.OOCError)
	} else if preExecutionResponse.OOGError != nil {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			IPAddress:   ip,
			Source:      event.Source_Node,
			Component:   event.Component_Pool,
			Level:       event.Level_Warning,
			EventID:     event.EventID_PreexecutionOOG,
			Description: tx.Hash().String(),
		}

		err := p.eventLog.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error adding event: %v", err)
		}
	}

	gasPrices, err := p.GetGasPrices(ctx)
	if err != nil {
		return err
	}

	err = p.ValidateBreakEvenGasPrice(ctx, tx, preExecutionResponse.txResponse.GasUsed, gasPrices)
	if err != nil {
		return err
	}

	poolTx := NewTransaction(tx, ip, isWIP)
	poolTx.GasUsed = preExecutionResponse.txResponse.GasUsed
	poolTx.ZKCounters = preExecutionResponse.usedZKCounters
	poolTx.ReservedZKCounters = preExecutionResponse.reservedZKCounters

	return p.storage.AddTx(ctx, *poolTx)
}

// ValidateBreakEvenGasPrice validates the effective gas price
func (p *Pool) ValidateBreakEvenGasPrice(ctx context.Context, tx types.Transaction, preExecutionGasUsed uint64, gasPrices GasPrices) error {
	// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price and l2 gas price
	txGasPrice, l2GasPrice := p.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice(), gasPrices.L1GasPrice, gasPrices.L2GasPrice)

	breakEvenGasPrice, err := p.effectiveGasPrice.CalculateBreakEvenGasPrice(tx.Data(), txGasPrice, preExecutionGasUsed, gasPrices.L1GasPrice)
	if err != nil {
		if p.cfg.EffectiveGasPrice.Enabled {
			log.Errorf("error calculating BreakEvenGasPrice: %v", err)
			return err
		} else {
			log.Warnf("EffectiveGasPrice is disabled, but failed to calculate BreakEvenGasPrice: %s", err)
			return nil
		}
	}

	reject := false
	loss := new(big.Int).SetUint64(0)

	tmpFactor := new(big.Float).Mul(new(big.Float).SetInt(breakEvenGasPrice), new(big.Float).SetFloat64(p.cfg.EffectiveGasPrice.BreakEvenFactor))
	breakEvenGasPriceWithFactor := new(big.Int)
	tmpFactor.Int(breakEvenGasPriceWithFactor)

	if breakEvenGasPriceWithFactor.Cmp(txGasPrice) == 1 { // breakEvenGasPriceWithMargin > txGasPrice
		// check against l2GasPrice now
		biL2GasPrice := big.NewInt(0).SetUint64(l2GasPrice)
		if txGasPrice.Cmp(biL2GasPrice) == -1 { // txGasPrice < l2GasPrice
			// reject tx
			reject = true
		} else {
			// accept loss
			loss = loss.Sub(breakEvenGasPriceWithFactor, txGasPrice)
		}
	}

	log.Infof("egp-log: txGasPrice(): %v, breakEvenGasPrice: %v, breakEvenGasPriceWithFactor: %v, gasUsed: %v, reject: %t, loss: %v, L1GasPrice: %d, L2GasPrice: %d, Enabled: %t, tx: %s",
		txGasPrice, breakEvenGasPrice, breakEvenGasPriceWithFactor, preExecutionGasUsed, reject, loss, gasPrices.L1GasPrice, l2GasPrice, p.cfg.EffectiveGasPrice.Enabled, tx.Hash().String())

	// Reject transaction if EffectiveGasPrice is enabled
	if p.cfg.EffectiveGasPrice.Enabled && reject {
		log.Infof("reject tx with gasPrice lower than L2GasPrice, tx: %s", tx.Hash().String())
		return ErrEffectiveGasPriceGasPriceTooLow
	}

	return nil
}

// preExecuteTx executes a transaction to calculate its zkCounters
func (p *Pool) preExecuteTx(ctx context.Context, tx types.Transaction) (preExecutionResponse, error) {
	response := preExecutionResponse{usedZKCounters: state.ZKCounters{}, reservedZKCounters: state.ZKCounters{}, OOCError: nil, OOGError: nil, isReverted: false}

	// TODO: Add effectivePercentage = 0xFF to the request (factor of 1) when gRPC message is updated
	processBatchResponse, err := p.state.PreProcessTransaction(ctx, &tx, nil)
	if err != nil {
		isOOC := executor.IsROMOutOfCountersError(executor.RomErrorCode(err))
		isOOG := errors.Is(err, runtime.ErrOutOfGas)
		if !isOOC && !isOOG {
			return response, err
		} else {
			if isOOC {
				response.OOCError = err
			}
			if isOOG {
				response.OOGError = err
			}
			if processBatchResponse != nil && processBatchResponse.BlockResponses != nil && len(processBatchResponse.BlockResponses) > 0 {
				response.usedZKCounters = processBatchResponse.UsedZkCounters
				response.reservedZKCounters = processBatchResponse.ReservedZkCounters
				response.txResponse = processBatchResponse.BlockResponses[0].TransactionResponses[0]
			}
			return response, nil
		}
	}

	if processBatchResponse.BlockResponses != nil && len(processBatchResponse.BlockResponses) > 0 {
		errorToCheck := processBatchResponse.BlockResponses[0].TransactionResponses[0].RomError
		response.isExecutorLevelError = processBatchResponse.IsExecutorLevelError
		if errorToCheck != nil {
			response.isReverted = errors.Is(errorToCheck, runtime.ErrExecutionReverted)
			if executor.IsROMOutOfCountersError(executor.RomErrorCode(errorToCheck)) {
				response.OOCError = err
			}
			if errors.Is(errorToCheck, runtime.ErrOutOfGas) {
				response.OOGError = err
			}
		} else {
			if !p.batchConstraintsCfg.IsWithinConstraints(processBatchResponse.UsedZkCounters) {
				err := fmt.Errorf("OutOfCounters Error (Node level) for tx: %s", tx.Hash().String())
				response.OOCError = err
				log.Error(err.Error())
			}
		}

		response.usedZKCounters = processBatchResponse.UsedZkCounters
		response.reservedZKCounters = processBatchResponse.ReservedZkCounters
		response.txResponse = processBatchResponse.BlockResponses[0].TransactionResponses[0]
	}

	return response, nil
}

// GetPendingTxs from the pool
// limit parameter is used to limit amount of pending txs from the db,
// if limit = 0, then there is no limit
func (p *Pool) GetPendingTxs(ctx context.Context, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusPending, limit)
}

// GetNonWIPPendingTxs from the pool
func (p *Pool) GetNonWIPPendingTxs(ctx context.Context) ([]Transaction, error) {
	return p.storage.GetNonWIPPendingTxs(ctx)
}

// GetSelectedTxs gets selected txs from the pool db
func (p *Pool) GetSelectedTxs(ctx context.Context, limit uint64) ([]Transaction, error) {
	return p.storage.GetTxsByStatus(ctx, TxStatusSelected, limit)
}

// GetPendingTxHashesSince returns the hashes of pending tx since the given date.
func (p *Pool) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	return p.storage.GetPendingTxHashesSince(ctx, since)
}

// UpdateTxStatus updates a transaction state accordingly to the
// provided state and hash
func (p *Pool) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus TxStatus, isWIP bool, failedReason *string) error {
	return p.storage.UpdateTxStatus(ctx, TxStatusUpdateInfo{
		Hash:         hash,
		NewStatus:    newStatus,
		IsWIP:        isWIP,
		FailedReason: failedReason,
	})
}

// SetGasPrices sets the current L2 Gas Price and L1 Gas Price
func (p *Pool) SetGasPrices(ctx context.Context, l2GasPrice uint64, l1GasPrice uint64) error {
	return p.storage.SetGasPrices(ctx, l2GasPrice, l1GasPrice)
}

// DeleteGasPricesHistoryOlderThan deletes gas prices older than a given date except the most recent one
func (p *Pool) DeleteGasPricesHistoryOlderThan(ctx context.Context, date time.Time) error {
	return p.storage.DeleteGasPricesHistoryOlderThan(ctx, date)
}

// GetGasPrices returns the current L2 Gas Price and L1 Gas Price
func (p *Pool) GetGasPrices(ctx context.Context) (GasPrices, error) {
	l2GasPrice, l1GasPrice, err := p.storage.GetGasPrices(ctx)
	return GasPrices{L1GasPrice: l1GasPrice, L2GasPrice: l2GasPrice}, err
}

// CountPendingTransactions get number of pending transactions
// used in bench tests
func (p *Pool) CountPendingTransactions(ctx context.Context) (uint64, error) {
	return p.storage.CountTransactionsByStatus(ctx, TxStatusPending)
}

// IsTxPending check if tx is still pending
func (p *Pool) IsTxPending(ctx context.Context, hash common.Hash) (bool, error) {
	return p.storage.IsTxPending(ctx, hash)
}

func (p *Pool) validateTx(ctx context.Context, poolTx Transaction) error {
	// Make sure the IP is valid.
	if poolTx.IP != "" && !IsValidIP(poolTx.IP) {
		return ErrInvalidIP
	}

	// Make sure the transaction is signed properly.
	if err := state.CheckSignature(poolTx.Transaction); err != nil {
		return ErrInvalidSender
	}

	// check chain id
	txChainID := poolTx.ChainId().Uint64()
	if txChainID != p.chainID && txChainID != 0 {
		return ErrInvalidChainID
	}

	// Accept only legacy transactions until EIP-2718/2930 activates.
	if poolTx.Type() != types.LegacyTxType {
		return ErrTxTypeNotSupported
	}

	// check Pre EIP155 txs signature
	if txChainID == 0 && !state.IsPreEIP155Tx(poolTx.Transaction) {
		return ErrInvalidSender
	}

	// gets tx sender for validations
	from, err := state.GetSender(poolTx.Transaction)
	if err != nil {
		return ErrInvalidSender
	}

	// Reject transactions over defined size to prevent DOS attacks
	decodedTx, err := state.EncodeTransaction(poolTx.Transaction, 0xFF, p.cfg.ForkID) //nolint: gomnd
	if err != nil {
		return ErrTxTypeNotSupported
	}

	if uint64(len(decodedTx)) > p.cfg.MaxTxBytesSize {
		log.Infof("%v: %v", ErrOversizedData.Error(), from.String())
		return ErrOversizedData
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if poolTx.Value().Sign() < 0 {
		return ErrNegativeValue
	}

	// check if sender is blocked
	_, blocked := p.blockedAddresses.Load(from.String())
	if blocked {
		log.Infof("%v: %v", ErrBlockedSender.Error(), from.String())
		return ErrBlockedSender
	}

	lastL2Block, err := p.state.GetLastL2Block(ctx, nil)
	if err != nil {
		log.Errorf("failed to load last l2 block while adding tx to the pool", err)
		return err
	}

	currentNonce, err := p.state.GetNonce(ctx, from, lastL2Block.Root())
	if err != nil {
		log.Errorf("failed to get nonce while adding tx to the pool", err)
		return err
	}
	// Ensure the transaction adheres to nonce ordering
	if poolTx.Nonce() < currentNonce {
		return ErrNonceTooLow
	}

	// check if sender has reached the limit of transactions in the pool
	if p.cfg.AccountQueue > 0 {
		// txCount, err := p.storage.CountTransactionsByFromAndStatus(ctx, from, TxStatusPending)
		// if err != nil {
		// 	return err
		// }
		// if txCount >= p.cfg.AccountQueue {
		// 	return ErrTxPoolAccountOverflow
		// }

		// Ensure the transaction does not jump out of the expected AccountQueue
		if poolTx.Nonce() > currentNonce+p.cfg.AccountQueue-1 {
			log.Infof("%v: %v", ErrNonceTooHigh.Error(), from.String())
			return ErrNonceTooHigh
		}
	}

	// check if the pool is full
	if p.cfg.GlobalQueue > 0 {
		txCount, err := p.storage.CountTransactionsByStatus(ctx, TxStatusPending)
		if err != nil {
			log.Errorf("failed to count pool txs by status pending while adding tx to the pool", err)
			return err
		}
		if txCount >= p.cfg.GlobalQueue {
			return ErrTxPoolOverflow
		}
	}

	// Reject transactions with a gas price lower than the minimum gas price
	p.minSuggestedGasPriceMux.RLock()
	gasPriceCmp := poolTx.GasPrice().Cmp(p.minSuggestedGasPrice)
	if gasPriceCmp == -1 {
		log.Debugf("low gas price: minSuggestedGasPrice %v got %v", p.minSuggestedGasPrice, poolTx.GasPrice())
	}
	p.minSuggestedGasPriceMux.RUnlock()
	if gasPriceCmp == -1 {
		return ErrGasPrice
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance, err := p.state.GetBalance(ctx, from, lastL2Block.Root())
	if err != nil {
		log.Errorf("failed to get balance for account %v while adding tx to the pool", from.String(), err)
		return err
	}

	if balance.Cmp(poolTx.Cost()) < 0 {
		return ErrInsufficientFunds
	}

	// Ensure the transaction has more gas than the basic poolTx fee.
	intrGas, err := IntrinsicGas(poolTx.Transaction)
	if err != nil {
		return err
	}
	if poolTx.Gas() < intrGas {
		return ErrIntrinsicGas
	}

	// try to get a transaction from the pool with the same nonce to check
	// if the new one has a price bump
	oldTxs, err := p.storage.GetTxsByFromAndNonce(ctx, from, poolTx.Nonce())
	if err != nil {
		log.Errorf("failed to txs for the same account and nonce while adding tx to the pool", err)
		return err
	}

	// check if the new transaction has more gas than all the other txs in the pool
	// with the same from and nonce to be able to replace the current txs by the new
	// when being selected
	for _, oldTx := range oldTxs {
		// discard invalid txs
		if oldTx.Status == TxStatusInvalid || oldTx.Status == TxStatusFailed {
			continue
		}

		oldTxPrice := new(big.Int).Mul(oldTx.GasPrice(), new(big.Int).SetUint64(oldTx.Gas()))
		txPrice := new(big.Int).Mul(poolTx.GasPrice(), new(big.Int).SetUint64(poolTx.Gas()))

		if oldTx.Hash() == poolTx.Hash() {
			return ErrAlreadyKnown
		}

		// if old Tx Price is higher than the new poolTx price, it returns an error
		if oldTxPrice.Cmp(txPrice) > 0 {
			return ErrReplaceUnderpriced
		}
	}

	// Executor field size requirements check
	if err := p.checkTxFieldCompatibilityWithExecutor(ctx, poolTx.Transaction); err != nil {
		return err
	}

	return nil
}

// pollMinSuggestedGasPrice polls the minimum L2 gas price since the previous
// check accordingly to the configured interval and tries to update it
func (p *Pool) pollMinSuggestedGasPrice(ctx context.Context) {
	fromTimestamp := time.Now().UTC().Add(-p.cfg.MinAllowedGasPriceInterval.Duration)
	// Ensuring we don't use a timestamp before the pool start as it may be using older L1 gas price factor
	if fromTimestamp.Before(p.startTimestamp) {
		fromTimestamp = p.startTimestamp
	}

	l2GasPrice, err := p.storage.MinL2GasPriceSince(ctx, fromTimestamp)
	if err != nil {
		if err == state.ErrNotFound {
			log.Warnf("No suggested min gas price since: %v", fromTimestamp)
		} else {
			log.Errorf("Error getting min gas price since: %v", fromTimestamp)
		}
	} else {
		p.tryUpdateMinSuggestedGasPrice(l2GasPrice)
	}
}

// tryUpdateMinSuggestedGasPrice tries to update the min suggested gas price
// with the provided minSuggestedGasPrice, it updates if the provided value
// is different from the value already store in p.minSuggestedGasPriceMux
func (p *Pool) tryUpdateMinSuggestedGasPrice(minSuggestedGasPrice uint64) {
	p.minSuggestedGasPriceMux.Lock()
	if p.minSuggestedGasPrice == nil || p.minSuggestedGasPrice.Uint64() != minSuggestedGasPrice {
		p.minSuggestedGasPrice = big.NewInt(0).SetUint64(minSuggestedGasPrice)
		log.Infof("Min suggested gas price updated to: %d", minSuggestedGasPrice)
	}
	p.minSuggestedGasPriceMux.Unlock()
}

// checkTxFieldCompatibilityWithExecutor checks the field sizes of the transaction to make sure
// they ar compatible with the Executor needs
// GasLimit: 256 bits
// GasPrice: 256 bits
// Value: 256 bits
// Data: 30000 bytes
// Nonce: 64 bits
// To: 160 bits
// ChainId: 64 bits
func (p *Pool) checkTxFieldCompatibilityWithExecutor(ctx context.Context, tx types.Transaction) error {
	maxUint64BigInt := big.NewInt(0).SetUint64(math.MaxUint64)

	// GasLimit, Nonce and To fields are limited by their types, no need to check
	// Gas Price and Value are checked against the balance, and the max balance allowed
	// by the merkletree service is uint256, in this case, if the transaction has a
	// gas price or value bigger than uint256, the check against the balance will
	// reject the transaction

	dataSize := len(tx.Data())
	if dataSize > p.cfg.MaxTxDataBytesSize {
		return fmt.Errorf("data size bigger than allowed, current size is %v bytes and max allowed is %v bytes", dataSize, p.cfg.MaxTxDataBytesSize)
	}

	if tx.ChainId().Cmp(maxUint64BigInt) == 1 {
		return fmt.Errorf("chain id higher than allowed, max allowed is %v", uint64(math.MaxUint64))
	}

	return nil
}

// DeleteReorgedTransactions deletes transactions from the pool
func (p *Pool) DeleteReorgedTransactions(ctx context.Context, transactions []*types.Transaction) error {
	hashes := []common.Hash{}

	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}

	return p.storage.DeleteTransactionsByHashes(ctx, hashes)
}

// UpdateTxWIPStatus updates a transaction wip status accordingly to the
// provided WIP status and hash
func (p *Pool) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	return p.storage.UpdateTxWIPStatus(ctx, hash, isWIP)
}

// GetDefaultMinGasPriceAllowed return the configured DefaultMinGasPriceAllowed value
func (p *Pool) GetDefaultMinGasPriceAllowed() uint64 {
	return p.cfg.DefaultMinGasPriceAllowed
}

// GetL1AndL2GasPrice returns the L1 and L2 gas price from memory struct
func (p *Pool) GetL1AndL2GasPrice() (uint64, uint64) {
	p.gasPricesMux.RLock()
	gasPrices := p.gasPrices
	p.gasPricesMux.RUnlock()

	return gasPrices.L1GasPrice, gasPrices.L2GasPrice
}

const (
	txDataNonZeroGas      uint64 = 16
	txGasContractCreation uint64 = 53000
	txGas                 uint64 = 21000
	txDataZeroGas         uint64 = 4
)

// CalculateEffectiveGasPrice calculates the final effective gas price for a tx
func (p *Pool) CalculateEffectiveGasPrice(rawTx []byte, txGasPrice *big.Int, txGasUsed uint64, l1GasPrice uint64, l2GasPrice uint64) (*big.Int, error) {
	return p.effectiveGasPrice.CalculateEffectiveGasPrice(rawTx, txGasPrice, txGasUsed, l1GasPrice, l2GasPrice)
}

// CalculateEffectiveGasPricePercentage calculates the gas price's effective percentage
func (p *Pool) CalculateEffectiveGasPricePercentage(gasPrice *big.Int, effectiveGasPrice *big.Int) (uint8, error) {
	return p.effectiveGasPrice.CalculateEffectiveGasPricePercentage(gasPrice, effectiveGasPrice)
}

// EffectiveGasPriceEnabled returns if effective gas price calculation is enabled or not
func (p *Pool) EffectiveGasPriceEnabled() bool {
	return p.effectiveGasPrice.IsEnabled()
}

// IntrinsicGas computes the 'intrinsic gas' for a given transaction.
func IntrinsicGas(tx types.Transaction) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if tx.To() == nil {
		gas = txGasContractCreation
	} else {
		gas = txGas
	}
	dataLen := uint64(len(tx.Data()))
	// Bump the required gas by the amount of transactional data
	if dataLen > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range tx.Data() {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := txDataNonZeroGas
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, ErrGasUintOverflow
		}
		gas += nz * nonZeroGas

		z := dataLen - nz
		if (math.MaxUint64-gas)/txDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * txDataZeroGas
	}
	return gas, nil
}
