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
