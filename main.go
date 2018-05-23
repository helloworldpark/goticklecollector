package main

import (
	"fmt"
	"time"

	"github.com/helloworldpark/goticklecollector/collector"
)

func main() {
	coinoneGW := collector.GiveWork(collector.CoinoneCollector{}, 3*time.Second)
	dfBundle, dbBundle := collector.Gather(coinoneGW)

	fmt.Println("Main!")

	for _, b := range dfBundle {
		go func(bundle collector.CoinGateway) {
			for coin := range bundle.Channel() {
				fmt.Println(fmt.Sprintf("DF %v", coin))
			}
		}(b)
	}

	func(bundle collector.CoinGateway) {
		for coin := range bundle.Channel() {
			fmt.Println(fmt.Sprintf("DB %v", coin))
		}
	}(dbBundle)

}
