package api

type Params map[string]string

type RESTResource struct {
	name  string
	value string
}

type OutboundAPI struct {
	vendor   Vendor
	restList []RESTResource
	params   Params
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
