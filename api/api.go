package api

import (
	"fmt"
	"strings"
)

type querys map[string]string

type restResource struct {
	name  string
	value string
}

type outboundAPI struct {
	vendor   vendor
	restList []restResource
	params   querys
}

func (v vendor) tickerAPI(currency string) outboundAPI {
	if v == coinone {
		rest := restResource{name: "ticker", value: ""}
		restList := []restResource{rest}
		query := make(querys)
		query["currency"] = currency
		outbound := outboundAPI{vendor: v, restList: restList, params: query}
		return outbound
	}
	if v == gopax {
		tradingName := strings.ToUpper(currency) + "-KRW"
		rest1 := restResource{name: "trading-pairs", value: tradingName}
		rest2 := restResource{name: "ticker", value: ""}
		restList := []restResource{rest1, rest2}
		query := make(querys)
		outbound := outboundAPI{vendor: v, restList: restList, params: query}
		return outbound
	}
	panic(fmt.Sprintf("Not prepared for %s", v.name))
}
