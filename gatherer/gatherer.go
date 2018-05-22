package gatherer

import (
	"time"

	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/collector"
)

type Gatherer struct {
	vendor       api.Vendor
	lastQty      float64
	lastBaseTime time.Time
}

func (g Gatherer) Gather(coins []collector.Coin) {

}

func (g Gatherer) flush() {

}

func (g Gatherer) update() {

}
