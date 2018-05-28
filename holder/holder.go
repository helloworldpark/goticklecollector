package holder

import (
	"github.com/helloworldpark/goticklecollector/logger"

	"github.com/helloworldpark/goticklecollector/collector"
)

// Provider is an interface any data provider must implement.
type Provider interface {
	// CanProvide checks if the provider can provide data of the query
	// Params:
	//    from: int64, UNIX timestamp, should be smaller than to
	//    to: int64, UNIX timestamp, should be bigger than from
	CanProvide(from, to int64) bool
	// Provide provides data as much as it can in the query interval.
	// Params:
	//    from: int64, UNIX timestamp, should be smaller than to
	//    to: int64, UNIX timestamp, should be bigger than from
	Provide(from, to int64) []collector.Coin
	CanProvideLast(seconds int64) bool
	ProvideLast(seconds int64) []collector.Coin
}

// Holder is a struct to manage accumulation and tracking of coin data.
type Holder struct {
	Vendor         string
	Currency       string
	records        *[]collector.Coin
	capacity       int
	lastTradeTime  int64
	lastTradeCount int
}

// New is a constructor of Holder
func New(vendor, currency string, capacity int) Holder {
	h := Holder{}

	h.Vendor = vendor
	h.Currency = currency
	h.capacity = capacity
	records := make([]collector.Coin, 0, capacity)
	h.records = &records

	return h
}

// CanProvide checks if the provider can provide data of the query
// Params:
//    from: int64, UNIX timestamp, should be smaller than to
//    to: int64, UNIX timestamp, should be bigger than from
func (h Holder) CanProvide(from, to int64) bool {
	records := h.records
	if len(*records) == 0 {
		return false
	}
	headInbounds := from >= (*records)[0].Timestamp
	tailInbounds := to <= (*records)[len(*records)-1].Timestamp
	return headInbounds && tailInbounds
}

// CanProvideLast checks if the provider can provide data of the query
// Params:
//    seconds: int64, UNIX timestamp, in seconds.
func (h Holder) CanProvideLast(seconds int64) bool {
	records := h.records
	if len(*records) == 0 {
		return false
	}
	headTime := (*records)[0].Timestamp
	tailTime := (*records)[len(*records)-1].Timestamp
	headInbounds := tailTime-seconds >= headTime
	return headInbounds
}

// Provide provides data as much as it can in the query interval.
// Params:
//    from: int64, UNIX timestamp, should be smaller than to
//    to: int64, UNIX timestamp, should be bigger than from
func (h Holder) Provide(from, to int64) []collector.Coin {
	records := h.records
	isIn := func(x int64) bool {
		return from <= x && x <= to
	}
	result := make([]collector.Coin, 0)
	for _, coin := range *records {
		if isIn(coin.Timestamp) {
			result = append(result, coin)
		}
	}
	return result
}

// ProvideLast provides data as much as it can in the query interval.
// Params:
//    seconds: int64, UNIX timestamp, in seconds.
func (h *Holder) ProvideLast(seconds int64) []collector.Coin {
	records := h.records
	logger.Info("[%s] Provide %d data", h.Currency, len(*records))
	if len(*records) == 0 {
		return *records
	}
	lastTimestamp := (*records)[len(*records)-1].Timestamp
	result := make([]collector.Coin, 0)
	for i := len(*records) - 1; i >= 0; i-- {
		coin := (*records)[i]
		if coin.Timestamp < lastTimestamp-seconds {
			break
		}
		result = append(result, coin)
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result
}

// UpdatingPipe starts to update coin information
func (h Holder) UpdatingPipe(gateway collector.CoinGateway) <-chan collector.Coin {
	updateChannel := make(chan collector.Coin)
	go func() {
		defer close(updateChannel)

		for coins := range gateway.Channel() {
			h.update(coins, updateChannel)
		}
	}()
	return updateChannel
}

func (h *Holder) update(coins []collector.Coin, output chan collector.Coin) {
	for len(*(h.records)) >= h.capacity {
		*(h.records) = (*(h.records))[1:]
	}

	// Skip data older than last time
	// if same to last time, skip for last same time count
	// else, append all

	lastTradeCounter := 0
	for _, coin := range coins {
		if h.lastTradeTime > coin.Timestamp {
			continue
		}
		if h.lastTradeTime == coin.Timestamp {
			if h.lastTradeCount > lastTradeCounter {
				lastTradeCounter++
				continue
			}
			length := len(*(h.records))
			last := (*(h.records))[length-1]
			(&last).Price = coin.Price
			(&last).Qty += coin.Qty

			lastTradeCounter++
		}
		if h.lastTradeTime < coin.Timestamp {
			lastTradeCounter = 1
			*(h.records) = append(*(h.records), coin)
			// Add not the latest but 1 index earlier to DB
			length := len(*(h.records))
			if length >= 2 {
				last := (*(h.records))[length-2]
				output <- last
			}
		}
		h.lastTradeTime = coin.Timestamp
		h.lastTradeCount = lastTradeCounter
	}
}
