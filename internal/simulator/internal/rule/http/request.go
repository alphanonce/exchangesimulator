package http

type Request struct {
	Method      string
	Host        string
	Path        string
	QueryString string
	Header      map[string][]string
	Body        []byte
}
