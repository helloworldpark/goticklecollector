package api

type vendor struct {
	name string
	url  string
}

var vendors []vendor

func addVendor(name string, url string) vendor {
	v := vendor{name, url}
	vendors = append(vendors, v)
	return v
}

var (
	coinone = addVendor("Coinone", "https://api.coinone.co.kr/")
	gopax   = addVendor("Gopax", "https://api.gopax.co.kr/")
)
