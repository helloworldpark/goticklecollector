package collector

import (
	"fmt"
	"log"

	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/utils"
)

// Collector is an interface a collector should do.
// Collect() collects and wraps in a Coin slice.
type Collector interface {
	Collect() []Coin
	Currency() string
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

// NewCollector is a constructor of Collector interface
func NewCollector(v api.Vendor, currency string) Collector {
	if v == api.Gopax {
		return GopaxCollector{currency: currency}
	}
	if v == api.Coinone {
		return CoinoneCollector{currency: currency}
	}
	panic(fmt.Sprintf("Not prepared for %s", v.Name))
}

// NewCollectors is a constructor of Collector interface
func NewCollectors(v api.Vendor, currencies []string) []Collector {
	if v == api.Gopax {
		collectors := make([]Collector, len(currencies))
		for i := 0; i < len(collectors); i++ {
			collectors[i] = GopaxCollector{currency: currencies[i]}
		}
		return collectors
	}
	if v == api.Coinone {
		collectors := make([]Collector, len(currencies))
		for i := 0; i < len(collectors); i++ {
			collectors[i] = CoinoneCollector{currency: currencies[i]}
		}
		return collectors
	}
	panic(fmt.Sprintf("Not prepared for %s", v.Name))
}

// Collect collects coin data from Coinone.
// Returns all coin data collected by ticker.
func (collector CoinoneCollector) Collect() []Coin {
	status, contents, errs := api.Coinone.TradesAPI(collector.currency).Request()
	if status != 200 || len(errs) > 0 {
		log.Printf("Status: %d, Errs: %v", status, errs)
		return make([]Coin, 0)
	}

	coins := JSONToCoinTrades(api.Coinone, contents)
	for i := 0; i < len(coins); i++ {
		coins[i].Currency = collector.currency
		coins[i].Vendor = api.Coinone.Name
	}

	return coins
}

// Currency returns the currency the collector is now collecting
func (collector CoinoneCollector) Currency() string {
	return collector.currency
}

// Collect collects coin data from Gopax.
// Returns coin data collected by ticker, by specified currency of the collector.
func (collector GopaxCollector) Collect() []Coin {
	status, contents, errs := api.Gopax.TickerAPI(collector.currency).Request()
	if status != 200 || len(errs) > 0 {
		log.Printf("Status: %d, Errs: %v", status, errs)
		return make([]Coin, 0)
	}

	coins := JSONToCoinTrades(api.Gopax, contents)
	for i := 0; i < len(coins); i++ {
		coins[i].Currency = collector.currency
		coins[i].Vendor = api.Gopax.Name
	}
	// Reverse it, since Coinone is ascending
	for i := len(coins)/2 - 1; i >= 0; i-- {
		opp := len(coins) - 1 - i
		coins[i], coins[opp] = coins[opp], coins[i]
	}

	return coins
}

// Currency returns the currency the collector is now collecting
func (collector GopaxCollector) Currency() string {
	return collector.currency
}
