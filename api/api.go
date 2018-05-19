package api

type Params map[string]string

type API interface {
	buildURL(rest []RESTResource, params Params) string
}

type RESTResource struct {
	name  string
	value string
}

type OutboundAPI struct {
	vendor Vendor
}

func (api OutboundAPI) buildURL(rest []RESTResource, params Params) string {
	base := api.vendor.url
	resource := ""
	for _, r := range rest {
		resource += r.name + "/"
		if r.value != "" {
			resource += r.value + "/"
		}
	}
	resource = resource[:(len(resource) - 1)]
	queryString := "?"
	for k, v := range params {
		queryString += k + "=" + v + "&"
	}
	queryString = queryString[:(len(queryString) - 1)]
	return base + resource + queryString
}
