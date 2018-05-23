package holder

import (
	"fmt"

	"github.com/helloworldpark/goticklecollector/collector"
)

type Holder struct {
	vendor   string
	currency string
	records  []collector.Coin
	capacity int
}

func New(vendor, currency string, capacity int) Holder {
	h := Holder{}

	h.vendor = vendor
	h.currency = currency
	h.capacity = capacity
	h.records = make([]collector.Coin, 0, capacity)

	return h
}

func (h Holder) StartUpdate(gateway collector.CoinGateway) {
	for coin := range gateway.Channel() {
		h.update(coin)
	}
}

func (h *Holder) update(coin collector.Coin) {
	for len(h.records) >= h.capacity {
		h.records = h.records[1:]
	}
	h.records = append(h.records, coin)
	fmt.Println(fmt.Sprintf("NOW %s:%s: %d", h.vendor, h.currency, len(h.records)))
}
