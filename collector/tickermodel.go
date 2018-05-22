package collector

import (
	"encoding/json"
	"errors"
	"fmt"

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

type tickerModelCoinone struct {
	Price     int
	Qty       float32
	Timestamp int64
	Currency  string
}

type tickerModelGopax struct {
	Price     int     `json:"price"`
	Qty       float32 `json:"volume"`
	Timestamp string  `json:"time"`
	Currency  string  `json:"-"`
}

type coinConvertable interface {
	convert() Coin
}

// JSONToCoin converts JSON string to slice of Coin depending on the vendor
func JSONToCoin(vendor api.Vendor, s string) []Coin {
	var convertables []coinConvertable
	var err error
	if vendor == api.Coinone {
		convertables, err = jsonToCoinoneTicker(s)
	} else if vendor == api.Gopax {
		convertables, err = jsonToGopaxTicker(s)
	}
	if err != nil {
		panic(err)
	}
	coins := make([]Coin, 0, len(convertables))
	for _, convertable := range convertables {
		coins = append(coins, convertable.convert())
	}
	return coins
}

func jsonToCoinoneTicker(s string) ([]coinConvertable, error) {
	bucket := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &bucket); err != nil {
		return make([]coinConvertable, 0), err
	}

	timestamp, ok := utils.ExtractInt64(bucket, "timestamp")
	if !ok {
		return make([]coinConvertable, 0), errors.New("Timestamp not existing in the response")
	}
	models := make([]coinConvertable, 0, len(bucket))
	for _, v := range bucket {
		fmt.Println(fmt.Sprintf("Type = %v", v))
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		model, err := mapToCoinoneTicker(m)
		if err != nil {
			return make([]coinConvertable, 0), err
		}
		model.Timestamp = timestamp
		models = append(models, model)
	}
	return models, nil
}

func mapToCoinoneTicker(m map[string]interface{}) (tickerModelCoinone, error) {
	model := tickerModelCoinone{}
	model.Price, _ = utils.ExtractInt32(m, "last")
	model.Qty, _ = utils.ExtractFloat32(m, "volume")
	model.Currency, _ = m["currency"].(string)
	return model, nil
}

func jsonToGopaxTicker(s string) ([]coinConvertable, error) {
	model := tickerModelGopax{}
	if err := json.Unmarshal([]byte(s), &model); err != nil {
		return make([]coinConvertable, 0), err
	}
	return []coinConvertable{model}, nil
}

func (model tickerModelCoinone) convert() Coin {
	coin := Coin{}
	coin.Vendor = api.Coinone.Name
	coin.Currency = model.Currency
	coin.Price = model.Price
	coin.Qty = model.Qty
	coin.Timestamp = model.Timestamp

	return coin
}

func (model tickerModelGopax) convert() Coin {
	coin := Coin{}
	coin.Vendor = api.Gopax.Name
	coin.Currency = model.Currency
	coin.Price = model.Price
	coin.Qty = model.Qty
	timestamp, err := utils.TimestampAsInt64(model.Timestamp)
	if err != nil {
		panic(err)
	}
	coin.Timestamp = timestamp
	return coin
}
