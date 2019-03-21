package collector

import (
	"encoding/json"

	"github.com/helloworldpark/goticklecollector/logger"

	"github.com/helloworldpark/goticklecollector/api"
)

// Coin defines what to collect and save.
type Coin struct {
	Vendor    string  `json:"vendor"`
	Currency  string  `json:"currency"`
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
	Qty       float64 `json:"qty"`
}

type tradesModelCoinone struct {
	Result         string              `json:"result"`
	ErrorCode      string              `json:"errorCode"`
	Timestamp      int64               `json:"timestamp,string"`
	CompleteOrders []tradeModelCoinone `json:"completeOrders"`
}

type tradeModelCoinone struct {
	Price     float64 `json:"price,string"`
	Qty       float64 `json:"qty,string"`
	Timestamp int64   `json:"timestamp,string"`
}

type tradesModelGopax struct {
	Price     float64 `json:"price"`
	Qty       float64 `json:"amount"`
	Timestamp int64   `json:"date"`
}

type coinConvertable interface {
	convert() Coin
}

// JSONToCoinTrades converts JSON string to slice of Coin depending on the vendor
func JSONToCoinTrades(vendor api.Vendor, jsonString string) []Coin {
	var convertables []coinConvertable
	var err error
	if vendor == api.Coinone {
		convertables, err = jsonToCoinoneTrades(jsonString)
	} else if vendor == api.Gopax {
		convertables, err = jsonToGopaxTrades(jsonString)
	}
	if err != nil {
		logger.Error("[Collector] %v", err)
		return make([]Coin, 0)
	}
	coins := make([]Coin, 0, len(convertables))
	for _, convertable := range convertables {
		coins = append(coins, convertable.convert())
	}
	return coins
}

func jsonToCoinoneTrades(jsonString string) ([]coinConvertable, error) {
	trades := tradesModelCoinone{}
	if err := json.Unmarshal([]byte(jsonString), &trades); err != nil {
		return make([]coinConvertable, 0), err
	}

	convertables := make([]coinConvertable, 0, len(trades.CompleteOrders))
	for _, trade := range trades.CompleteOrders {
		convertables = append(convertables, trade)
	}
	return convertables, nil
}

func jsonToGopaxTrades(jsonString string) ([]coinConvertable, error) {
	trades := make([]tradesModelGopax, 0)
	if err := json.Unmarshal([]byte(jsonString), &trades); err != nil {
		return make([]coinConvertable, 0), err
	}
	convertables := make([]coinConvertable, 0, len(trades))
	for _, trade := range trades {
		convertables = append(convertables, trade)
	}
	return convertables, nil
}

func (m tradeModelCoinone) convert() Coin {
	coin := Coin{}
	coin.Price = m.Price
	coin.Qty = m.Qty
	coin.Timestamp = m.Timestamp
	coin.Vendor = api.Coinone.Name
	return coin
}

func (m tradesModelGopax) convert() Coin {
	coin := Coin{}
	coin.Price = m.Price
	coin.Qty = m.Qty
	coin.Timestamp = m.Timestamp
	coin.Vendor = api.Gopax.Name
	return coin
}
