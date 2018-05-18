package api

type RequestType int

var requestTypes []string

func (e RequestType) String() string {
	return requestTypes[int(e)]
}

func riota(s string) RequestType {
	requestTypes = append(requestTypes, s)
	return RequestType(len(requestTypes) - 1)
}

var (
	GET    = riota("GET")
	POST   = riota("POST")
	PUT    = riota("PUT")
	DELETE = riota("DELETE")
)

type Vendor struct {
	name string
	url  string
}

var vendors []Vendor

func addVendor(name string, url string) Vendor {
	v := Vendor{name, url}
	vendors = append(vendors, v)
	return v
}

var (
	Coinone = addVendor("Coinone", "https://www.coinone.co.kr")
)

type OutboundAPI struct {
	vendor      Vendor
	requestType RequestType
}
