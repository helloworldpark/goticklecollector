package api

import (
	"fmt"

	"github.com/parnurzeal/gorequest"

	"encoding/json"
)

var requesterMap map[string]*gorequest.SuperAgent

func init() {
	requesterMap = map[string]*gorequest.SuperAgent{}
	requesterMap[coinone.name] = gorequest.New()
	requesterMap[gopax.name] = gorequest.New()
}

func (api outboundAPI) request() (int, map[string]interface{}) {
	requester := requesterMap[api.vendor.name]
	if requester == nil {
		panic(fmt.Sprintf("Not prepared for requesting to %s", api.vendor.name))
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

func (api outboundAPI) buildURL() string {
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
