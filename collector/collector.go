package collector

import (
	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/utils"
)

type coin struct {
	vendor    string
	currency  string
	timestamp int64
	price     int
	qty       float32
}

type collector interface {
	collect() []coin
	report(coins []coin)
}

type coinoneCollector struct {
	currency string
}

type gopaxCollector struct {
	currency string
}

var (
	coinoneCurrencies utils.StringSet
	gopaxCurrencies   utils.StringSet
)

func init() {
	coinoneCurrencies = make(utils.StringSet)
	coinoneList := []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
	for _, e := range coinoneList {
		coinoneCurrencies.Add(e)
	}

	gopaxCurrencies = make(utils.StringSet)
	gopaxList := []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
	for _, e := range gopaxList {
		gopaxCurrencies.Add(e)
	}
}

func (collector coinoneCollector) collect() []coin {
	status, contents := api.Coinone.TickerAPI("all").Request()
	if status != 200 {
		return make([]coin, 0)
	}
	if contents["errorMsg"] != nil {
		return make([]coin, 0)
	}

	timestamp, ok := utils.ExtractInt64(contents, "timestamp")
	if !ok {
		return make([]coin, 0)
	}

	coins := make([]coin, 0, len(contents))
	for currency, obj := range contents {
		if coinoneCurrencies.Contains(currency) == false {
			continue
		}
		info, ok := obj.(map[string]interface{})
		if !ok {
			continue
		}

		price, ok := utils.ExtractInt32(info, "last")
		if !ok {
			continue
		}

		qty, ok := utils.ExtractFloat32(info, "volume")
		if !ok {
			continue
		}

		coin := coin{}
		coin.vendor = api.Coinone.Name
		coin.currency = currency
		coin.timestamp = timestamp
		coin.price = price
		coin.qty = qty
		coins = append(coins, coin)
	}
	return coins
}

func (collector gopaxCollector) collect() []coin {
	status, contents := api.Gopax.TickerAPI(collector.currency).Request()
	if status != 200 {
		return make([]coin, 0)
	}
	if contents["errorMsg"] != nil {
		return make([]coin, 0)
	}

	if gopaxCurrencies.Contains(collector.currency) == false {
		return make([]coin, 0)
	}

	timestamp, ok := utils.ExtractTimestamp(contents, "time")
	if !ok {
		return make([]coin, 0)
	}

	price, ok := utils.ExtractFloat64(contents, "price")
	if !ok {
		return make([]coin, 0)
	}

	qty, ok := utils.ExtractFloat64(contents, "volume")
	if !ok {
		return make([]coin, 0)
	}

	c := coin{}
	c.vendor = api.Coinone.Name
	c.currency = collector.currency
	c.timestamp = timestamp
	c.price = int(price)
	c.qty = float32(qty)
	coins := []coin{c}
	return coins
}
