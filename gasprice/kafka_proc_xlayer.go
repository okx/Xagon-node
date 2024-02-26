package gasprice

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

const (
	okbcoinId      = 7184
	ethcoinId      = 15756
	defaultTime    = 10
	defaultMaxData = 10e6 // 10M
)

var (
	// ErrNotFindCoinPrice not find a correct coin price
	ErrNotFindCoinPrice = errors.New("not find a correct coin price")
)

// MsgInfo msg info
type MsgInfo struct {
	Topic string `json:"topic"`
	Data  *Body  `json:"data"`
}

// Body msg body
type Body struct {
	Id        string   `json:"id"`
	PriceList []*Price `json:"priceList"`
}

// Price coin price
type Price struct {
	CoinId                   int     `json:"coinId"`
	Symbol                   string  `json:"symbol"`
	FullName                 string  `json:"fullName"`
	Price                    float64 `json:"price"`
	PriceStatus              int     `json:"priceStatus"`
	MaxPrice24H              float64 `json:"maxPrice24H"`
	MinPrice24H              float64 `json:"minPrice24H"`
	MarketCap                float64 `json:"marketCap"`
	Timestamp                int64   `json:"timestamp"`
	Vol24H                   float64 `json:"vol24h"`
	CirculatingSupply        float64 `json:"circulatingSupply"`
	MaxSupply                float64 `json:"maxSupply"`
	TotalSupply              float64 `json:"totalSupply"`
	PriceChange24H           float64 `json:"priceChange24H"`
	PriceChangeRate24H       float64 `json:"priceChangeRate24H"`
	CirculatingMarketCap     float64 `json:"circulatingMarketCap"`
	PriceChange7D            float64 `json:"priceChange7D"`
	PriceChangeRate7D        float64 `json:"priceChangeRate7D"`
	PriceChange30D           float64 `json:"priceChange30D"`
	PriceChangeRate30D       float64 `json:"priceChangeRate30D"`
	PriceChangeYearStart     float64 `json:"priceChangeYearStart"`
	PriceChangeRateYearStart float64 `json:"priceChangeRateYearStart"`
	ExceptionStatus          int     `json:"exceptionStatus"`
	Source                   int     `json:"source"`
	Type                     string  `json:"type"`
	Id                       string  `json:"id"`
}

// L1L2PriceRecord l1 l2 coin price record
type L1L2PriceRecord struct {
	l1Price  float64
	l2Price  float64
	l1Update bool
	l2Update bool
}

// KafkaProcessor kafka processor
type KafkaProcessor struct {
	cfg       Config
	kreader   *kafka.Reader
	ctx       context.Context
	rwLock    sync.RWMutex
	l1CoinId  int
	l2CoinId  int
	l1Price   float64
	l2Price   float64
	tmpPrices L1L2PriceRecord
}

func newKafkaProcessor(cfg Config, ctx context.Context) *KafkaProcessor {
	rp := &KafkaProcessor{
		cfg:      cfg,
		kreader:  getKafkaReader(cfg),
		l1Price:  cfg.DefaultL1CoinPrice,
		l2Price:  cfg.DefaultL2CoinPrice,
		ctx:      ctx,
		l2CoinId: okbcoinId,
		l1CoinId: ethcoinId,
	}
	if cfg.L2CoinId != 0 {
		rp.l2CoinId = cfg.L2CoinId
	}
	if cfg.L1CoinId != 0 {
		rp.l1CoinId = cfg.L1CoinId
	}

	go rp.processor()
	return rp
}

func getKafkaReader(cfg Config) *kafka.Reader {
	brokers := strings.Split(cfg.KafkaURL, ",")

	var dialer *kafka.Dialer
	if cfg.Password != "" && cfg.Username != "" && cfg.RootCAPath != "" {
		rootCA, err := os.ReadFile(cfg.RootCAPath)
		if err != nil {
			panic("kafka read root ca fail")
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(rootCA); !ok {
			panic("caCertPool.AppendCertsFromPEM")
		}
		dialer = &kafka.Dialer{
			Timeout:       defaultTime * time.Second,
			DualStack:     true,
			SASLMechanism: plain.Mechanism{Username: cfg.Username, Password: cfg.Password},
		}
		{ // #nosec G402
			dialer.TLS = &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: true}
		}
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     cfg.GroupID,
		Topic:       cfg.Topic,
		MinBytes:    1, // 1
		MaxBytes:    defaultMaxData,
		Dialer:      dialer,
		StartOffset: kafka.LastOffset, // read data from new message
	})
}

func (rp *KafkaProcessor) processor() {
	log.Info("kafka processor start processor ")
	defer rp.kreader.Close()
	for {
		select {
		case <-rp.ctx.Done():
			return
		default:
			err := rp.ReadAndUpdate(rp.ctx)
			if err != nil && err != ErrNotFindCoinPrice {
				log.Warn("get the destion data fail ", err)
				time.Sleep(time.Second * defaultTime)
				continue
			}
		}
	}
}

// ReadAndUpdate read and update
func (rp *KafkaProcessor) ReadAndUpdate(ctx context.Context) error {
	m, err := rp.kreader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	return rp.Update(m.Value)
}

// Update update the coin price
func (rp *KafkaProcessor) Update(data []byte) error {
	if rp.cfg.Type == FixedType {
		price, err := rp.parseCoinPrice(data, []int{rp.l2CoinId})
		if err == nil {
			rp.updateL2CoinPrice(price[rp.l2CoinId])
		}
		return err
	} else if rp.cfg.Type == FollowerType {
		prices, err := rp.parseCoinPrice(data, []int{rp.l1CoinId, rp.l2CoinId})
		if err == nil {
			rp.updateL1L2CoinPrice(prices)
		}
		return err
	}
	return nil
}

func (rp *KafkaProcessor) updateL2CoinPrice(price float64) {
	rp.rwLock.Lock()
	defer rp.rwLock.Unlock()
	rp.l2Price = price
}

// GetL2CoinPrice get L2 coin price
func (rp *KafkaProcessor) GetL2CoinPrice() float64 {
	rp.rwLock.RLock()
	defer rp.rwLock.RUnlock()
	return rp.l2Price
}

func (rp *KafkaProcessor) updateL1L2CoinPrice(prices map[int]float64) {
	if len(prices) == 0 {
		return
	}
	rp.rwLock.Lock()
	defer rp.rwLock.Unlock()
	if v, ok := prices[rp.l1CoinId]; ok {
		rp.tmpPrices.l1Price = v
		rp.tmpPrices.l1Update = true
	}
	if v, ok := prices[rp.l2CoinId]; ok {
		rp.tmpPrices.l2Price = v
		rp.tmpPrices.l2Update = true
	}
	if rp.tmpPrices.l1Update && rp.tmpPrices.l2Update {
		rp.l1Price = rp.tmpPrices.l1Price
		rp.l2Price = rp.tmpPrices.l2Price
		rp.tmpPrices.l1Update = false
		rp.tmpPrices.l2Update = false
		return
	}
}

// GetL1L2CoinPrice get l1, L2 coin price
func (rp *KafkaProcessor) GetL1L2CoinPrice() (float64, float64) {
	rp.rwLock.RLock()
	defer rp.rwLock.RUnlock()
	return rp.l1Price, rp.l2Price
}

func (rp *KafkaProcessor) parseCoinPrice(value []byte, coinIds []int) (map[int]float64, error) {
	if len(coinIds) == 0 {
		return nil, fmt.Errorf("the params coinIds is empty")
	}
	msgI := &MsgInfo{}
	err := json.Unmarshal(value, &msgI)
	if err != nil {
		return nil, err
	}
	if msgI.Data == nil || len(msgI.Data.PriceList) == 0 {
		return nil, fmt.Errorf("the data PriceList is empty")
	}
	mp := make(map[int]*Price)
	for _, price := range msgI.Data.PriceList {
		mp[price.CoinId] = price
	}

	results := make(map[int]float64)
	for _, coinId := range coinIds {
		if coin, ok := mp[coinId]; ok {
			results[coinId] = coin.Price
		} else {
			log.Debugf("not find a correct coin price coin id is =%v", coinId)
		}
	}
	if len(results) == 0 {
		return results, ErrNotFindCoinPrice
	}
	return results, nil
}
