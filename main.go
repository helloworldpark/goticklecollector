package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/helloworldpark/goticklecollector/api"

	"github.com/helloworldpark/goticklecollector/collector"
	"github.com/helloworldpark/goticklecollector/holder"
)

func main() {
	fmt.Println("Main!")

	coll := collector.CoinoneCollector{}

	holders := make([]holder.Holder, 0)

	for _, currency := range coll.Currencies() {
		h := holder.New(api.Coinone.Name, currency, 3)
		holders = append(holders, h)
	}
	sort.Slice(holders, func(i, j int) bool {
		return strings.Compare(holders[i].Currency, holders[j].Currency) <= 0
	})

	coinoneGW := collector.GiveWork(coll, 3*time.Second)
	dfBundle, dbBundle := collector.Gather(coinoneGW)

	for i, bundle := range dfBundle {
		go func(idx int, g collector.CoinGateway) {
			holders[idx].StartUpdate(g)
		}(i, bundle)
	}

	func(bundle collector.CoinGateway) {
		for coin := range bundle.Channel() {
			_ = coin.Timestamp
		}
	}(dbBundle)
}
