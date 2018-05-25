package holder

import (
	"time"

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
	Vendor               string
	Currency             string
	records              []collector.Coin
	capacity             int
	lastBaseTime         int64
	lastCumulativeVolume float64
	isFirst              bool
}

// New is a constructor of Holder
func New(vendor, currency string, capacity int) Holder {
	h := Holder{}

	h.Vendor = vendor
	h.Currency = currency
	h.capacity = capacity
	h.records = make([]collector.Coin, 0, capacity)
	h.isFirst = true

	return h
}

// StartUpdate starts to update coin information
func (h Holder) StartUpdate(gateway collector.CoinGateway) {
	for coin := range gateway.Channel() {
		h.update(coin)
	}
}

// CanProvide checks if the provider can provide data of the query
// Params:
//    from: int64, UNIX timestamp, should be smaller than to
//    to: int64, UNIX timestamp, should be bigger than from
func (h Holder) CanProvide(from, to int64) bool {
	records := h.records
	if len(records) == 0 {
		return false
	}
	headInbounds := from >= records[0].Timestamp
	tailInbounds := to <= records[len(records)-1].Timestamp
	return headInbounds && tailInbounds
}

// CanProvideLast checks if the provider can provide data of the query
// Params:
//    seconds: int64, UNIX timestamp, in seconds.
func (h Holder) CanProvideLast(seconds int64) bool {
	records := h.records
	if len(records) == 0 {
		return false
	}
	headTime := records[0].Timestamp
	tailTime := records[len(records)-1].Timestamp
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
	for _, coin := range records {
		if isIn(coin.Timestamp) {
			result = append(result, coin)
		}
	}
	return result
}

// ProvideLast provides data as much as it can in the query interval.
// Params:
//    seconds: int64, UNIX timestamp, in seconds.
func (h Holder) ProvideLast(seconds int64) []collector.Coin {
	records := h.records
	lastTimestamp := records[len(records)-1].Timestamp
	result := make([]collector.Coin, 0)
	for i := len(records) - 1; i >= 0; i-- {
		coin := records[i]
		if coin.Timestamp > lastTimestamp-seconds {
			break
		}
		result = append(result, coin)
	}
	return result
}

func (h *Holder) updateBaseTime() bool {
	current := time.Now()
	timeDiff := time.Duration(current.Unix() - h.lastBaseTime)
	if timeDiff > (time.Hour*24)/time.Second {
		y, m, d := current.Date()
		newBase := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		h.lastBaseTime = newBase.Unix()

		return true
	}
	return false
}

func (h *Holder) update(coin collector.Coin) {
	// Update time
	baseTimeChanged := h.updateBaseTime()
	if baseTimeChanged {
		h.lastCumulativeVolume = coin.Qty
	}
	if h.isFirst {
		h.isFirst = false
	} else {
		for len(h.records) >= h.capacity {
			h.records = h.records[1:]
		}
		newCoin := coin
		if baseTimeChanged {
			newCoin.Qty = coin.Qty
		} else {
			newCoin.Qty = coin.Qty - h.lastCumulativeVolume
		}
		h.lastCumulativeVolume = coin.Qty
		h.records = append(h.records, newCoin)
	}
}
