package api

import (
	"fmt"

	"github.com/parnurzeal/gorequest"

	"encoding/json"
)

var requesterMap map[string]*gorequest.SuperAgent

func init() {
	requesterMap = map[string]*gorequest.SuperAgent{}
	requesterMap[Coinone.Name] = gorequest.New()
	requesterMap[Gopax.Name] = gorequest.New()
}

// Request makes a request to the vendor using the data from api.
func (api OutboundAPI) Request() (int, map[string]interface{}) {
	requester := requesterMap[api.vendor.Name]
	if requester == nil {
		panic(fmt.Sprintf("Not prepared for requesting to %s", api.vendor.Name))
	}

	targetURL := api.buildURL()
	resp, body, errs := requester.Get(targetURL).End()
	contents := make(map[string]interface{})
	if len(errs) == 0 {
		err := json.Unmarshal([]byte(body), &contents)
		if err != nil {
			contents["errorMsg"] = err.Error()
		}
	} else {
		errMsg := ""
		for _, err := range errs {
			errMsg += err.Error() + "\n"
		}
		contents["errorMsg"] = errMsg
	}
	return resp.StatusCode, contents
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
