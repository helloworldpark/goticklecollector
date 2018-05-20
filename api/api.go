package api

import (
	"fmt"
	"strings"
)

// Vendor is a struct representing cryptocurrency trading service company's open API.
type Vendor struct {
	Name string
	url  string
}

var vendors []Vendor

func addVendor(name string, url string) Vendor {
	v := Vendor{name, url}
	vendors = append(vendors, v)
	return v
}

var (
	// Coinone is a const struct representing Coinone(https://www.coinone.co.kr/)
	Coinone = addVendor("Coinone", "https://api.coinone.co.kr/")
	// Gopax is a const struct representing Gopax(https://www.gopax.co.kr/)
	Gopax = addVendor("Gopax", "https://api.gopax.co.kr/")
)

// Querys is a alias of map[string]string, a map for query strings.
type Querys map[string]string

// RestResource is a struct for a single REST API resource.
// name: name of the REST API.
// value: (Optional) value of the REST API.
// For example, if we want API such as '/user/shinji'
// rest := RestResource{"user", "shinji"}
type RestResource struct {
	name  string
	value string
}

// OutboundAPI is a struct wrapper for outbount API call.
// Note that only GET is supported.
// vendor: Vendor of the cryptocurrency trading service.
// restList: []RestResource for constructing REST API.
// params: Querys for constructing query string.
type OutboundAPI struct {
	vendor   Vendor
	restList []RestResource
	params   Querys
}

// TickerAPI constructs an OutboundAPI struct for calling /ticker for each vendor.
// currency: string representing type of currency whose data we collect.
func (v Vendor) TickerAPI(currency string) OutboundAPI {
	if v == Coinone {
		rest := RestResource{name: "ticker", value: ""}
		restList := []RestResource{rest}
		query := make(Querys)
		query["currency"] = currency
		outbound := OutboundAPI{vendor: v, restList: restList, params: query}
		return outbound
	}
	if v == Gopax {
		tradingName := strings.ToUpper(currency) + "-KRW"
		rest1 := RestResource{name: "trading-pairs", value: tradingName}
		rest2 := RestResource{name: "ticker", value: ""}
		restList := []RestResource{rest1, rest2}
		query := make(Querys)
		outbound := OutboundAPI{vendor: v, restList: restList, params: query}
		return outbound
	}
	panic(fmt.Sprintf("Not prepared for %s", v.Name))
}
