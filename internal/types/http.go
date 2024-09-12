package types

type Request struct {
	Method      string
	Host        string
	Path        string
	QueryString string
	Body        []byte
}

type Response struct {
	StatusCode int
	Body       []byte
}
