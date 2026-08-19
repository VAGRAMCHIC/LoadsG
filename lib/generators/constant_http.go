package generators

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"loadsg/lib/model"
	"loadsg/lib/repository"
)

type ConstantHttp struct {
	repo   repository.HttpLoadRepository
	client *http.Client
}

func NewConstantHttp(repo repository.HttpLoadRepository) *ConstantHttp {
	return &ConstantHttp{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (h *ConstantHttp) Name() string {
	return "constant_http"
}

func (h *ConstantHttp) Run(ctx context.Context, job model.LoadJob) error {
	profile, err := h.repo.GetConstantById(ctx, job.Id)
	if err != nil {
		return fmt.Errorf("load constant profile: %w", err)
	}

	if profile.Count <= 0 {
		return fmt.Errorf("count must be positive")
	}
	if profile.URL == "" {
		return fmt.Errorf("url is required")
	}
	if profile.Method == "" {
		profile.Method = "GET"
	}

	for sent := 0; sent < profile.Count; sent++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := h.send(ctx, profile); err != nil {
				log.Printf("constant_http request failed for job %s (%d/%d): %v",
					profile.LoadJobId, sent+1, profile.Count, err)
			}
		}
	}

	return nil
}

func (h *ConstantHttp) send(ctx context.Context, profile *model.ConstantHttpLoad) error {
	req, err := http.NewRequestWithContext(
		ctx,
		strings.ToUpper(profile.Method),
		profile.URL,
		bodyReader(profile.Body),
	)
	if err != nil {
		return err
	}

	for key, value := range profile.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}