package api

import (
	"fmt"
	"testing"
)

func TestBuildURL(t *testing.T) {
	vendor := Coinone
	rest1 := RESTResource{"test1", ""}
	rest2 := RESTResource{"test2", "myname"}
	restList := []RESTResource{rest1, rest2}
	params := make(Params)
	params["abc"] = string(123)
	params["cde"] = "kkk"

	outbound := OutboundAPI{vendor}
	url := outbound.buildURL(restList, params)

	fmt.Println(url)
	fmt.Println(url)
	// Output:
	// https://api.coinone.co.kr/test1/test2/myname?abc=123&cde=kkk
}
