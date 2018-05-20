package api

import (
	"fmt"

	"github.com/parnurzeal/gorequest"

	"encoding/json"
)

var requesterMap map[string]*gorequest.SuperAgent

func init() {
	requesterMap = map[string]*gorequest.SuperAgent{}
	requesterMap[Coinone.name] = gorequest.New()
	requesterMap[Gopax.name] = gorequest.New()
}

func (api OutboundAPI) request() (int, map[string]interface{}) {
	requester := requesterMap[api.vendor.name]
	if requester == nil {
		panic(fmt.Sprintf("Not prepared for requesting to %s", api.vendor.name))
	}

	targetURL := api.buildURL()
	resp, body, errs := requester.Get(targetURL).End()
	if len(errs) == 0 {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(body), &jsonMap)
		if err != nil {
			jsonMap["errorMsg"] = err.Error()
		}
		return resp.StatusCode, jsonMap
	}
	errMsg := ""
	for _, err := range errs {
		errMsg += err.Error() + "\n"
	}
	errContents := map[string]interface{}{"errorMsg": errMsg}
	return resp.StatusCode, errContents
}
