package api

type querys map[string]string

type restResource struct {
	name  string
	value string
}

type outboundAPI struct {
	vendor   vendor
	restList []restResource
	params   querys
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
