package generators

import (
	"context"
	"net/http"
	"time"

	"loadsg/lib/model"
)

type ConstantHttp struct{}

func (h *ConstantHttp) Name() string {
	return "http_flood"
}

func (h *ConstantHttp) Run(ctx context.Context, job model.LoadJob) error {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, _ = client.Get("https://example.com")
		}
	}
}
