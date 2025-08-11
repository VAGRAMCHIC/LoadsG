package lib

type _httpMethod string

const (
	GET     _httpMethod = "GET"
	POST    _httpMethod = "POST"
	PUT     _httpMethod = "PUT"
	DELETE  _httpMethod = "DELETE"
	CONNECT _httpMethod = "CONNECT"
	PATCH   _httpMethod = "PATCH"
	OPTIONS _httpMethod = "OPTIONS"
	HEAD    _httpMethod = "HEAD"
)

type HttpHead struct {
	method        _httpMethod
	url           string
	proto_version string
	length        int
	headers       map[string]string
}

type LoadRrequest struct {
	Id      string `json:"id" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
	Request string `json:"request" binding:"required"`
}

type HTTPLoadRequest struct {
	Id       int      `json:"id" binding:"required"`
	HttpHead HttpHead `json:"httpHead" binding:"required"`
	Body     string   `json:"body" binding:"required"`
}
