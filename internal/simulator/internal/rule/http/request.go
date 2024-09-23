package http

type Request struct {
	Method      string
	Host        string
	Path        string
	QueryString string
	Body        []byte
}
