package main

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/helloworldpark/goticklecollector/api"
	"github.com/helloworldpark/goticklecollector/logger"

	"github.com/helloworldpark/goticklecollector/collector"
	"github.com/helloworldpark/goticklecollector/holder"
)

func main() {
	defer logger.Close()
	// Parse flags
	credPath := flag.String("credential", "", "Credential for DB access")
	flag.Parse()

	if credPath == nil || *credPath == "" {
		logger.Panic("No -credential provided")
	}

	// DB Holder
	dbHolder := holder.NewDBHolder(*credPath, 20)
	dbHolder.Init()

	dbChannels := generateDBChannels()

	// Setup Coin Collectors
	currencyCoinone := []string{"btc", "bch", "eth", "etc", "xrp", "qtum", "iota", "ltc", "btg", "omg", "eos"}
	currencyGopax := []string{"eth", "btc", "bch", "eos", "omg", "qtum", "xrp", "ltc"}
	holderMap := make(map[string][]holder.Holder)
	holderMap[api.Coinone.Name] = generateHolders(
		api.Coinone,
		currencyCoinone,
		20,
		10,
		&dbChannels)
	holderMap[api.Gopax.Name] = generateHolders(
		api.Gopax,
		currencyGopax,
		3,
		10,
		&dbChannels)

	// Feed to DB
	dbHolder.ConnectDBChannel(dbChannels)

	// Setup API
	router := gin.Default()
	if logger.IsLoggerGCE() {
		router.Use(ginLoggerWithCustomLogger())
	}
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

func ginLoggerWithCustomLogger(notlogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	handler := func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.Info("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment)
		}
	}
	return handler
}
