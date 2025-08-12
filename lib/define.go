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
	Method       _httpMethod       `json:"method"`
	URL          string            `json:"url"`
	ProtoVersion string            `json:"proto_version"`
	Length       int               `json:"length"`
	Headers      map[string]string `json:"headers"`
}

type HTTPLoadRequest struct {
	Id       int      `json:"id" binding:"required"`
	HttpHead HttpHead `json:"httpHead" binding:"required"`
	Body     string   `json:"body" binding:"required"`
	Count    int      `json:"count" binding:"required"`
}
