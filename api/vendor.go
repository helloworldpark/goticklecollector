package api

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
	Coinone = addVendor("Coinone", "https://api.coinone.co.kr/")
	Gopax   = addVendor("Gopax", "https://api.gopax.co.kr/")
)
