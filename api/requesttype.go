package api

type requestType int

var requestTypes []string

func (e requestType) String() string {
	return requestTypes[int(e)]
}

func riota(s string) requestType {
	requestTypes = append(requestTypes, s)
	return requestType(len(requestTypes) - 1)
}

var (
	GET    = riota("GET")
	POST   = riota("POST")
	PUT    = riota("PUT")
	DELETE = riota("DELETE")
)
