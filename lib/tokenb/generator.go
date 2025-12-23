package tokenb

import (
	"net/http"
	"context"
)

type StaticGenerator struct {
	req *http.Request
}

func NewStaticGenerator(url string) *StaticGenerator {
	req, _ := http.NewRequest("GET", url, nil)
	return &StaticGenerator{req: req}
}

func (g *StaticGenerator) Next() *http.Request {
	return g.req.Clone(context.Background())
}

