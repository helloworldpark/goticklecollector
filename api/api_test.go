package api

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestAPICall(t *testing.T) {
	vendor := Coinone
	rest := RestResource{"ticker", ""}
	restList := []RestResource{rest}
	params := make(Querys)
	params["currency"] = "eos"
	outbound := OutboundAPI{vendor, restList, params}

	status, contents := outbound.Request()

	if status != 200 || contents["errorMsg"] != nil {
		t.Errorf("Expected %d -> but got %d, %s", 200, status, contents["errorMsg"])
	} else {
		spewed := spew.Sdump(contents)
		t.Errorf("Success!\n %s", spewed)
	}
}

func TestAPIResourceGOPAX(t *testing.T) {
	gopaxReq := Gopax.TickerAPI("eos")

	status, contents := gopaxReq.Request()
	if status != 200 || contents["errorMsg"] != nil {
		t.Errorf("Failed GOPAX: %d, %s", status, contents["errorMsg"])
	}
}

func TestAPIResourceCOINONE(t *testing.T) {
	coinoneReq := Coinone.TickerAPI("all")

	status, contents := coinoneReq.Request()
	if status != 200 || contents["errorMsg"] != nil {
		t.Errorf("Failed COINONE: %d, %s", status, contents["errorMsg"])
	}
}

func TestBuildURL(t *testing.T) {
	vendor := Coinone
	rest1 := RestResource{"test1", ""}
	rest2 := RestResource{"test2", "myname"}
	restList := []RestResource{rest1, rest2}
	params := make(Querys)
	params["abc"] = "123"
	params["cde"] = "kkk"

	outbound := OutboundAPI{vendor, restList, params}
	url := outbound.buildURL()

	answer := "https://api.coinone.co.kr/test1/test2/myname?abc=123&cde=kkk"
	if url != answer {
		t.Errorf("Expected %s -> but got %s", answer, url)
	}
}
