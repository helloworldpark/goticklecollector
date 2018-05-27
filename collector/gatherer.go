package collector

import (
	"time"
)

// CoinGateway is a wrapper of Coin channel
// It provides only a property to <-chan []Coin
type CoinGateway struct {
	gateway chan []Coin
}

// Channel is a property to <-chan Coin
func (cg CoinGateway) Channel() <-chan []Coin {
	return cg.gateway
}

// GiveWork gives work to Collector. It runs a goroutine and returns channels per coin type
// Params:
//    c Collector: Type of coin to collect.
//    period time.Duration: Time to sleep in nanoseconds.
// Returns:
//    CoinGateway: CoinGateway per coin type
func GiveWork(c Collector, period time.Duration) CoinGateway {
	gateway := CoinGateway{}
	gateway.gateway = make(chan []Coin)

	closeAllGateways := func() {
		close(gateway.gateway)
	}

	go func() {
		defer closeAllGateways()
		for true {
			gateway.gateway <- c.Collect()
			time.Sleep(period)
		}
	}()

	return gateway
}

// MergeGateways merges slices of CoinGateway to one slice of CoinGateway
func MergeGateways(gateways ...[]CoinGateway) []CoinGateway {
	total := 0
	for _, gateway := range gateways {
		total += len(gateway)
	}
	merged := make([]CoinGateway, 0, total)
	for _, gateway := range gateways {
		for _, g := range gateway {
			merged = append(merged, g)
		}
	}
	return merged
}
