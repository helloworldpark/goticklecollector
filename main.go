package main

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/helloworldpark/goticklecollector/api"

	"github.com/helloworldpark/goticklecollector/collector"
	"github.com/helloworldpark/goticklecollector/holder"
)

func main() {
	// Parse flags
	credPath := flag.String("credential", "", "Credential for DB access")
	flag.Parse()

	if credPath == nil || *credPath == "" {
		panic("No -credential provided")
	}

	// DB Holder
	dbHolder := holder.NewDBHolder(*credPath, 100)
	dbHolder.Init()

	dbChannels := generateDBChannels()

	// Setup Coin Collectors
	currencies := []string{"eos", "btc"}
	holderMap := make(map[string][]holder.Holder)
	holderMap[api.Coinone.Name] = generateHolders(
		api.Coinone,
		currencies,
		20,
		10,
		&dbChannels)
	holderMap[api.Gopax.Name] = generateHolders(
		api.Gopax,
		currencies,
		3,
		10,
		&dbChannels)

	// Feed to DB
	dbHolder.ConnectDBChannel(dbChannels)

	// Setup API
	router := gin.Default()
	router.GET("/coins/last", func(c *gin.Context) {
		v := c.Query("vendor")
		cur := c.Query("currency")
		lastSeconds, _ := strconv.ParseInt(c.Query("seconds"), 10, 64)
		coins := make([]collector.Coin, 0)
		holders, ok := holderMap[v]
		if !ok {
			c.String(http.StatusBadRequest, "Bad request: vendor %s not supported", v)
		}
		for i := 0; i < len(holders); i++ {
			if (&holders[i]).Currency == cur {
				coins = (&holders[i]).ProvideLast(lastSeconds)
			}
		}
		c.JSON(http.StatusOK, coins)
	})
	router.GET("/coins/from/:from/to/:to", func(c *gin.Context) {
		v := c.Query("vendor")
		cur := c.Query("currency")

		from, err := strconv.ParseInt(c.Param("from"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request: invalid format of parameter 'from'")
			return
		}
		to, err := strconv.ParseInt(c.Param("to"), 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request: invalid format of parameter 'to'")
			return
		}

		if from >= to {
			c.String(http.StatusBadRequest, "Bad request: from must be smaller than to")
			return
		}

		coins := make([]collector.Coin, 0)
		holders, ok := holderMap[v]
		if !ok {
			c.String(http.StatusBadRequest, "Bad request: vendor %s not supported", v)
		}
		for i := 0; i < len(holders); i++ {
			if (&holders[i]).Currency == cur {
				coins = (&holders[i]).Provide(from, to)
			}
		}
		c.JSON(http.StatusOK, coins)
	})
	router.Run(":50001")
}

func generateDBChannels() []<-chan collector.Coin {
	return make([]<-chan collector.Coin, 0)
}

func generateHolders(vendor api.Vendor, currencyList []string, period time.Duration, capacity int, dbChannels *[]<-chan collector.Coin) []holder.Holder {
	if len(currencyList) == 0 {
		return make([]holder.Holder, 0)
	}

	collectors := collector.NewCollectors(vendor, currencyList)
	holders := make([]holder.Holder, 0)
	for _, col := range collectors {
		h := holder.New(vendor.Name, col.Currency(), capacity)
		worker := collector.GiveWork(col, period*time.Second)
		dbChannel := h.UpdatingPipe(worker)

		holders = append(holders, h)
		*dbChannels = append(*dbChannels, dbChannel)
	}
	return holders
}
