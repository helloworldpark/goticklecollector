package main

import (
	"fmt"
	"time"

	"github.com/helloworldpark/goticklecollector/collector"
	"github.com/helloworldpark/goticklecollector/gatherer"
)

func main() {
	coinoneGW := gatherer.GiveWork(collector.CoinoneCollector{}, 3*time.Second)
	dfBundle, dbBundle := gatherer.Gather(coinoneGW)

	fmt.Println("Main!")

	for _, b := range dfBundle {
		go func(bundle gatherer.CoinGateway) {
			for coin := range bundle.Channel() {
				fmt.Println(fmt.Sprintf("DF %v", coin))
			}
		}(b)
	}

	func(bundle gatherer.CoinGateway) {
		for coin := range bundle.Channel() {
			fmt.Println(fmt.Sprintf("DB %v", coin))
		}
	}(dbBundle)

}
