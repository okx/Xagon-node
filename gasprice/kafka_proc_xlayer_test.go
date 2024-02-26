package gasprice

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCoinPrice(t *testing.T) {
	testcases := []struct {
		coinIds []int
		msg     string
		check   func(prices map[int]float64, err error)
	}{
		{
			// param err
			coinIds: []int{},
			msg:     "{\"topic\":\"middle_coinPrice_push\"}",
			check: func(prices map[int]float64, err error) {
				require.Error(t, err)
			},
		},
		{
			// param err
			coinIds: []int{ethcoinId, okbcoinId},
			msg:     "{\"topic\":\"middle_coinPrice_push\"}",
			check: func(prices map[int]float64, err error) {
				require.Error(t, err)
			},
		},
		{
			// not find all, find one
			coinIds: []int{ethcoinId, okbcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, prices[ethcoinId], 0.02)
			},
		},
		{
			// not find all
			coinIds: []int{ethcoinId, okbcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", okbcoinId),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, prices[okbcoinId], 0.02)
			},
		},
		{
			// correct
			coinIds: []int{ethcoinId, okbcoinId, okbcoinId + 1},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.001}, {\"coinId\":%d,\"price\":0.002}, {\"coinId\":%d,\"price\":0.003}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId, okbcoinId, okbcoinId+1),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, len(prices), 3)
				require.Equal(t, prices[ethcoinId], 0.001)
				require.Equal(t, prices[okbcoinId], 0.002)
				require.Equal(t, prices[okbcoinId+1], 0.003)
			},
		},
		{
			// correct
			coinIds: []int{ethcoinId, okbcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}, {\"coinId\":%d,\"price\":0.002}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId, okbcoinId),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, len(prices), 2)
				require.Equal(t, prices[ethcoinId], 0.02)
				require.Equal(t, prices[okbcoinId], 0.002)
			},
		},
		{
			// correct
			coinIds: []int{ethcoinId, okbcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}, {\"coinId\":%d,\"price\":0.002}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", okbcoinId, ethcoinId),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, len(prices), 2)
				require.Equal(t, prices[ethcoinId], 0.002)
				require.Equal(t, prices[okbcoinId], 0.02)
			},
		},
		{
			// correct
			coinIds: []int{okbcoinId, ethcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.02}, {\"coinId\":%d,\"price\":0.002}, {\"coinId\":%d,\"price\":0.003}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", okbcoinId, ethcoinId, ethcoinId+1),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, len(prices), 2)
				require.Equal(t, prices[okbcoinId], 0.02)
				require.Equal(t, prices[ethcoinId], 0.002)
			},
		},
		{
			// correct
			coinIds: []int{okbcoinId},
			msg:     fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}, {\"coinId\":%d,\"price\":0.002}, {\"coinId\":123,\"price\":0.005}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId, okbcoinId),
			check: func(prices map[int]float64, err error) {
				require.NoError(t, err)
				require.Equal(t, len(prices), 1)
				require.Equal(t, prices[okbcoinId], 0.002)
			},
		},
	}

	for _, tc := range testcases {
		rp := newKafkaProcessor(Config{Topic: "middle_coinPrice_push"}, context.Background())
		rt, err := rp.parseCoinPrice([]byte(tc.msg), tc.coinIds)
		tc.check(rt, err)
	}
}

func TestUpdateL1L2CoinPrice(t *testing.T) {
	testcases := []struct {
		check func()
	}{
		{
			check: func() {
				rp := newKafkaProcessor(Config{Topic: "middle_coinPrice_push"}, context.Background())
				prices := map[int]float64{ethcoinId: 1.5, okbcoinId: 0.5}
				rp.updateL1L2CoinPrice(prices)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 1.5)
				require.Equal(t, l2, 0.5)
			},
		},
		{
			check: func() {
				rp := newKafkaProcessor(Config{Topic: "middle_coinPrice_push"}, context.Background())
				prices := map[int]float64{ethcoinId: 1.5}
				rp.updateL1L2CoinPrice(prices)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 0.0)
				require.Equal(t, l2, 0.0)
				require.Equal(t, rp.tmpPrices.l1Update, true)
				require.Equal(t, rp.tmpPrices.l2Update, false)

				prices = map[int]float64{okbcoinId: 0.5}
				rp.updateL1L2CoinPrice(prices)
				l1, l2 = rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 1.5)
				require.Equal(t, l2, 0.5)
				require.Equal(t, rp.tmpPrices.l1Update, false)
				require.Equal(t, rp.tmpPrices.l2Update, false)
			},
		},
		{
			check: func() {
				rp := newKafkaProcessor(Config{Topic: "middle_coinPrice_push"}, context.Background())
				prices := map[int]float64{okbcoinId: 0.5}
				rp.updateL1L2CoinPrice(prices)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 0.0)
				require.Equal(t, l2, 0.0)
				require.Equal(t, rp.tmpPrices.l1Update, false)
				require.Equal(t, rp.tmpPrices.l2Update, true)

				prices = map[int]float64{ethcoinId: 1.5}
				rp.updateL1L2CoinPrice(prices)
				l1, l2 = rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 1.5)
				require.Equal(t, l2, 0.5)
				require.Equal(t, rp.tmpPrices.l1Update, false)
				require.Equal(t, rp.tmpPrices.l2Update, false)
			},
		},
	}
	for _, tc := range testcases {
		tc.check()
	}
}

func TestUpdate(t *testing.T) {
	testcases := []struct {
		msg   string
		cfg   Config
		check func(rp *KafkaProcessor, err error)
	}{
		// FixedType
		{ // correct
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}, {\"coinId\":%d,\"price\":0.002}, {\"coinId\":123,\"price\":0.005}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId, okbcoinId),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FixedType},
			check: func(rp *KafkaProcessor, err error) {
				require.NoError(t, err)
				require.Equal(t, rp.GetL2CoinPrice(), 0.002)
			},
		},
		{ // not find
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FixedType},
			check: func(rp *KafkaProcessor, err error) {
				require.Equal(t, err, ErrNotFindCoinPrice)
				require.Equal(t, rp.GetL2CoinPrice(), float64(0))
			},
		},
		{ // not find
			msg: "{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}",
			cfg: Config{Topic: "middle_coinPrice_push", Type: FixedType},
			check: func(rp *KafkaProcessor, err error) {
				require.EqualError(t, err, "the data PriceList is empty")
				require.Equal(t, rp.GetL2CoinPrice(), float64(0))
			},
		},

		// FollowerType
		{ // correct
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}, {\"coinId\":%d,\"price\":0.002}, {\"coinId\":123,\"price\":0.005}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId, okbcoinId),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FollowerType},
			check: func(rp *KafkaProcessor, err error) {
				require.NoError(t, err)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, 0.04)
				require.Equal(t, l2, 0.002)
			},
		},
		{ // not find
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId+1),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FollowerType},
			check: func(rp *KafkaProcessor, err error) {
				require.Equal(t, err, ErrNotFindCoinPrice)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, float64(0))
				require.Equal(t, l2, float64(0))
			},
		},
		{ // find one but not update
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", ethcoinId),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FollowerType},
			check: func(rp *KafkaProcessor, err error) {
				require.NoError(t, err)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, float64(0))
				require.Equal(t, l2, float64(0))
			},
		},
		{ // find one but not update
			msg: fmt.Sprintf("{\"topic\":\"middle_coinPrice_push\",\"source\":null,\"type\":null,\"data\":{\"priceList\":[{\"coinId\":%d,\"price\":0.04}],\"id\":\"98a797ce-f61b-4e90-87ac-445e77ad3599\"}}", okbcoinId),
			cfg: Config{Topic: "middle_coinPrice_push", Type: FollowerType},
			check: func(rp *KafkaProcessor, err error) {
				require.NoError(t, err)
				l1, l2 := rp.GetL1L2CoinPrice()
				require.Equal(t, l1, float64(0))
				require.Equal(t, l2, float64(0))
			},
		},
	}

	for _, tc := range testcases {
		rp := newKafkaProcessor(tc.cfg, context.Background())
		err := rp.Update([]byte(tc.msg))
		tc.check(rp, err)
	}
}
