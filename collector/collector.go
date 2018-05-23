package collector

import (
	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/utils"
)

// Collector is an interface a collector should do.
// Collect() collects and wraps in a Coin slice.
// Count() is a number how many types of Coin will be in Collect()
type Collector interface {
	Collect() []Coin
	Count() int
	Currencies() []string
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
	status, contents, errs := api.Coinone.TickerAPI("all").Request()
	if status != 200 {
		return make([]Coin, 0)
	}
	if len(errs) > 0 {
		return make([]Coin, 0)
	}

	coins := JSONToCoin(api.Coinone, contents)

	return coins
}

// Count returns how many types of Coin collector will collect
// In case of CoinoneCollector, 11
func (collector CoinoneCollector) Count() int {
	return 11
}

func (collector CoinoneCollector) Currencies() []string {
	return []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
}

// Collect collects coin data from Gopax.
// Returns coin data collected by ticker, by specified currency of the collector.
func (collector GopaxCollector) Collect() []Coin {
	status, contents, errs := api.Gopax.TickerAPI(collector.currency).Request()
	if status != 200 {
		return make([]Coin, 0)
	}
	if len(errs) > 0 {
		return make([]Coin, 0)
	}

	coins := JSONToCoin(api.Gopax, contents)
	for i := 0; i < len(coins); i++ {
		coins[i].Currency = collector.currency
	}

	return coins
}

// Count returns how many types of Coin collector will collect
// In case of GopaxCollector, 1
func (collector GopaxCollector) Count() int {
	return 1
}

func (collector GopaxCollector) Currencies() []string {
	return []string{collector.currency}
}
