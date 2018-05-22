package collector

import (
	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/utils"
)

// Coin defines what to collect and save.
type Coin struct {
	Vendor    string
	Currency  string
	Timestamp int64
	Price     int
	Qty       float32
}

// Collector is an interface a collector should do.
// Collect() collects and wraps in a Coin slice.
type Collector interface {
	Collect() []Coin
}

// CoinoneCollector collects coin data from Coinone.
// currency: string value for currency, but meaningless. Default is 'all'
type CoinoneCollector struct {
	currency string
}

// GopaxCollector collects coin data from Gopax.
// currency: string value for currency
type GopaxCollector struct {
	currency string
}

var (
	// CoinoneCurrencies defines available currency to collect.
	CoinoneCurrencies utils.StringSet
	// GopaxCurrencies defines available currency to collect.
	GopaxCurrencies utils.StringSet
)

func init() {
	CoinoneCurrencies = make(utils.StringSet)
	coinoneList := []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
	for _, e := range coinoneList {
		CoinoneCurrencies.Add(e)
	}

	GopaxCurrencies = make(utils.StringSet)
	gopaxList := []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
	for _, e := range gopaxList {
		GopaxCurrencies.Add(e)
	}
}

// Collect collects coin data from Coinone.
// Returns all coin data collected by ticker.
func (collector CoinoneCollector) Collect() []Coin {
	status, contents := api.Coinone.TickerAPI("all").Request()
	if status != 200 {
		return make([]Coin, 0)
	}
	if contents["errorMsg"] != nil {
		return make([]Coin, 0)
	}

	timestamp, ok := utils.ExtractInt64(contents, "timestamp")
	if !ok {
		return make([]Coin, 0)
	}

	coins := make([]Coin, 0, len(contents))
	for currency, obj := range contents {
		if CoinoneCurrencies.Contains(currency) == false {
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

		coin := Coin{}
		coin.Vendor = api.Coinone.Name
		coin.Currency = currency
		coin.Timestamp = timestamp
		coin.Price = price
		coin.Qty = qty
		coins = append(coins, coin)
	}
	return coins
}

// Collect collects coin data from Gopax.
// Returns coin data collected by ticker, by specified currency of the collector.
func (collector GopaxCollector) Collect() []Coin {
	status, contents := api.Gopax.TickerAPI(collector.currency).Request()
	if status != 200 {
		return make([]Coin, 0)
	}
	if contents["errorMsg"] != nil {
		return make([]Coin, 0)
	}

	if GopaxCurrencies.Contains(collector.currency) == false {
		return make([]Coin, 0)
	}

	timestamp, ok := utils.ExtractTimestamp(contents, "time")
	if !ok {
		return make([]Coin, 0)
	}

	price, ok := utils.ExtractFloat64(contents, "price")
	if !ok {
		return make([]Coin, 0)
	}

	qty, ok := utils.ExtractFloat64(contents, "volume")
	if !ok {
		return make([]Coin, 0)
	}

	c := Coin{}
	c.Vendor = api.Coinone.Name
	c.Currency = collector.currency
	c.Timestamp = timestamp
	c.Price = int(price)
	c.Qty = float32(qty)
	coins := []Coin{c}
	return coins
}
