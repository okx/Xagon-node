package pool

import (
	"context"
	"errors"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"sync"
)

// Pool is an implementation of the Pool interface
// that uses a postgres database to store the data
type Pool struct {
	storage
	cfg pool.Config

	levelGp      []int64
	toAddrLevel  map[common.Address]int
	txsWithLevel map[int][]customTransaction
	signers      map[common.Address]bind.SignerFn

	sigMutex sync.Mutex
	txsMutex sync.Mutex

	prepareWG sync.WaitGroup
	speedWG   sync.WaitGroup
}

// NewPool creates and initializes an instance of Pool
func NewPool(cfg pool.Config, s storage) *Pool {
	p := &Pool{
		cfg:     cfg,
		storage: s,

		toAddrLevel:  make(map[common.Address]int),
		txsWithLevel: make(map[int][]customTransaction),
		signers:      make(map[common.Address]bind.SignerFn),
	}
	return p
}

func (p *Pool) PrepareTx(ctx context.Context, levelCount, addrPerLevel, txPerAddr int) {
	baseGP := int64(10)

	for i := 0; i < levelCount; i++ {
		p.levelGp = append(p.levelGp, baseGP*int64(i+1))

		// gen to addr
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
		toAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
		p.toAddrLevel[toAddr] = i

		p.prepareWG.Add(1)
		// gen from addr
		go func(level int, toAddr common.Address) {
			for j := 0; j < addrPerLevel; j++ {
				privateKey, err := crypto.GenerateKey()
				if err != nil {
					log.Fatal(err)
				}
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
				if err != nil {
					log.Fatal(err)
				}
				p.sigMutex.Lock()
				p.signers[auth.From] = auth.Signer
				p.sigMutex.Unlock()

				// gen tx
				for k := 0; k < txPerAddr; k++ {
					tx := types.NewTx(&types.LegacyTx{
						Nonce:    uint64(k),
						GasPrice: big.NewInt(0).SetInt64(baseGP),
						Gas:      24000,
						To:       &toAddr,
						Value:    big.NewInt(0),
					})

					signedTx, err := auth.Signer(auth.From, tx)
					if err != nil {
						log.Fatal(err)
					}
					customTx := customTransaction{
						tx:   *signedTx,
						from: auth.From,
					}
					err = p.AddTx(ctx, customTx, "127.0.0.1")
					if err != nil {
						log.Warn(err)
					}

					p.txsMutex.Lock()
					levelTxs := []customTransaction{customTx}
					if txs, ok := p.txsWithLevel[level]; ok {
						levelTxs = append(txs, levelTxs...)
					}
					p.txsWithLevel[level] = levelTxs
					p.txsMutex.Unlock()
				}
			}
			p.prepareWG.Done()
		}(i, toAddr)
	}
	p.prepareWG.Wait()
}

func (p *Pool) Speed(ctx context.Context) {
	for level, txs := range p.txsWithLevel {
		speedGP := p.levelGp[level]

		p.speedWG.Add(1)
		go func(speedGP int64, txs []customTransaction) {
			for _, tx := range txs {
				speedTx := types.NewTx(&types.LegacyTx{
					Nonce:    tx.tx.Nonce(),
					GasPrice: big.NewInt(0).SetInt64(speedGP),
					Gas:      tx.tx.Gas(),
					To:       tx.tx.To(),
					Value:    tx.tx.Value(),
					Data:     tx.tx.Data(),
				})

				p.sigMutex.Lock()
				signedTx, err := p.signers[tx.from](tx.from, speedTx)
				if err != nil {
					log.Fatal(err)
				}
				p.sigMutex.Unlock()

				customTx := customTransaction{
					tx:   *signedTx,
					from: tx.from,
				}
				err = p.AddTx(ctx, customTx, "127.0.0.1")
				if err != nil {
					log.Warn(err)
				}
			}
			p.speedWG.Done()
		}(speedGP, txs)
	}
	p.speedWG.Wait()
}

type customTransaction struct {
	tx   types.Transaction
	from common.Address
}

// AddTx adds a transaction to the pool with the pending state
func (p *Pool) AddTx(ctx context.Context, tx customTransaction, ip string) error {
	if err := p.validateTx(ctx, tx); err != nil {
		return err
	}

	return p.StoreTx(ctx, tx.tx, ip, false)
}

func (p *Pool) validateTx(ctx context.Context, tx customTransaction) error {
	processedCount, err := p.storage.CountTransactionsByFromAndStatus(ctx, tx.from, pool.TxStatusSelected, pool.TxStatusFailed)
	if err != nil {
		return err
	}

	if tx.tx.Nonce() < processedCount {
		return errors.New("invalid nonce")
	}

	return nil
}

// StoreTx adds a transaction to the pool with the pending state
func (p *Pool) StoreTx(ctx context.Context, tx types.Transaction, ip string, isWIP bool) error {
	// ignore preExecuteTx
	poolTx := pool.NewTransaction(tx, ip, isWIP, &pool.Pool{})

	return p.storage.AddTx(ctx, *poolTx)
}

// UpdateTxStatus updates a transaction state accordingly to the
// provided state and hash
func (p *Pool) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool, failedReason *string) error {
	return p.storage.UpdateTxStatus(ctx, pool.TxStatusUpdateInfo{
		Hash:         hash,
		NewStatus:    newStatus,
		IsWIP:        isWIP,
		FailedReason: failedReason,
	})
}

// UpdateTxWIPStatus updates a transaction wip status accordingly to the
// provided WIP status and hash
func (p *Pool) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	return p.storage.UpdateTxWIPStatus(ctx, hash, isWIP)
}

// GetDefaultMinGasPriceAllowed return the configured DefaultMinGasPriceAllowed value
func (p *Pool) GetDefaultMinGasPriceAllowed() uint64 {
	return 0
}

// GetL1AndL2GasPrice returns the L1 and L2 gas price from memory struct
func (p *Pool) GetL1AndL2GasPrice() (uint64, uint64) {
	return 0, 0
}

// CountPendingTransactions get number of pending transactions
// used in bench tests
func (p *Pool) CountPendingTransactions(ctx context.Context) (uint64, error) {
	return p.storage.CountTransactionsByStatus(ctx, pool.TxStatusPending)
}

// GetGasPrices returns the current L2 Gas Price and L1 Gas Price
func (p *Pool) GetGasPrices(ctx context.Context) (pool.GasPrices, error) {
	return pool.GasPrices{}, nil
}

func (p *Pool) GetLevelByAddr(addr common.Address) int {
	return p.toAddrLevel[addr]
}

func (p *Pool) CountTransactionsByFromAndStatus(ctx context.Context, from common.Address, status ...pool.TxStatus) (uint64, error) {
	return p.storage.CountTransactionsByFromAndStatus(ctx, from, status...)
}
