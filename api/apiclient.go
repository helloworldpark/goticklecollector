package api

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

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

var requesterMap map[string]*gorequest.SuperAgent

func init() {
	requesterMap = map[string]*gorequest.SuperAgent{}
	requesterMap[Coinone.Name] = gorequest.New()
	requesterMap[Gopax.Name] = gorequest.New()
}

// Request makes a request to the vendor using the data from api.
// Returns: response code int, body string, slice of error errs
func (api OutboundAPI) Request() (int, string, []error) {
	requester := requesterMap[api.vendor.Name]
	if requester == nil {
		panic(fmt.Sprintf("Not prepared for requesting to %s", api.vendor.Name))
	}

	targetURL := api.buildURL()
	resp, body, errs := requester.Get(targetURL).End()
	return resp.StatusCode, body, errs
}

func (api OutboundAPI) buildURL() string {
	base := api.vendor.url
	resource := ""
	for _, r := range api.restList {
		resource += r.name + "/"
		if r.value != "" {
			resource += r.value + "/"
		}
	}
	resource = resource[:(len(resource) - 1)]
	queryString := "?"
	for k, v := range api.params {
		queryString += k + "=" + v + "&"
	}
	queryString = queryString[:(len(queryString) - 1)]
	return base + resource + queryString
}
