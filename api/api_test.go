package api

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestAPICall(t *testing.T) {
	vendor := Coinone
	rest := RESTResource{"ticker", ""}
	restList := []RESTResource{rest}
	params := make(Params)
	params["currency"] = "eos"
	outbound := OutboundAPI{vendor, restList, params}

	status, contents := outbound.request()

	if status != 200 || contents["errorMsg"] != nil {
		t.Errorf("Expected %d -> but got %d, %s", 200, status, contents["errorMsg"])
	} else {
		spewed := spew.Sdump(contents)
		t.Errorf("Success!\n %s", spewed)
	}
}

func TestBuildURL(t *testing.T) {
	vendor := Coinone
	rest1 := RESTResource{"test1", ""}
	rest2 := RESTResource{"test2", "myname"}
	restList := []RESTResource{rest1, rest2}
	params := make(Params)
	params["abc"] = "123"
	params["cde"] = "kkk"

	outbound := OutboundAPI{vendor, restList, params}
	url := outbound.buildURL()

	answer := "https://api.coinone.co.kr/test1/test2/myname?abc=123&cde=kkk"
	if url != answer {
		t.Errorf("Expected %s -> but got %s", answer, url)
	}
}
