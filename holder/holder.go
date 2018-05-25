package holder

import (
	"time"

	"github.com/helloworldpark/goticklecollector/collector"
)

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
