package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/helloworldpark/goticklecollector/api"

	"github.com/helloworldpark/goticklecollector/collector"
	"github.com/helloworldpark/goticklecollector/holder"
)

func main() {
	fmt.Println("Main!")

	// Setting up Coin Collectors...
	coll := collector.CoinoneCollector{}
	holders := make([]holder.Holder, 0)

	for _, currency := range coll.Currencies() {
		h := holder.New(api.Coinone.Name, currency, 10)
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

	go func(bundle collector.CoinGateway) {
		for coin := range bundle.Channel() {
			_ = coin.Timestamp
		}
	}(dbBundle)

	fmt.Println("Start running")

	// Setup API
	router := gin.Default()
	router.GET("/coins/last", func(c *gin.Context) {
		// v := c.Query("vendor")
		cur := c.Query("currency")
		lastSeconds, _ := strconv.ParseInt(c.Query("seconds"), 10, 64)
		coins := make([]collector.Coin, 0)
		for i := 0; i < len(holders); i++ {
			if (&holders[i]).Currency == cur {
				coins = (&holders[i]).ProvideLast(lastSeconds)
			}
		}

		if len(coins) > 0 {
			c.JSON(http.StatusOK, coins)
		} else {
			c.String(http.StatusBadRequest, "Invalid request: currency=%s lastseconds=%d", cur, lastSeconds)
		}
	})
	router.Run(":50001")
}
