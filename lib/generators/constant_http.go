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
	"loadsg/lib/report"
)

type ConstantHttp struct {
	repo   repository.HttpLoadRepository
	client *http.Client
	report report.Config
}

func NewConstantHttp(repo repository.HttpLoadRepository, reportCfg report.Config) *ConstantHttp {
	return &ConstantHttp{
		repo: repo,
		report: reportCfg,
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
	col := report.NewCollector()
	defer func() {
		col.Finish()
		sum := col.Summary(job.Id, job.JobName, job.Type)
		path, err := report.Write(h.report, sum)
		if err != nil {
			log.Printf("constant_http report write failed for job %s: %v", job.Id, err)
			return
		}
		if path != "" {
			log.Printf("constant_http report written: %s", path)
		}
	}()

	for sent := 0; sent < profile.Count; sent++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			h.send(ctx, profile, col)
		}
	}
	return nil
}

func (h *ConstantHttp) send(ctx context.Context, profile *model.ConstantHttpLoad, col *report.Collector) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(profile.Method), profile.URL, bodyReader(profile.Body))
	if err != nil {
		col.Record(time.Since(start), 0, err)
		return
	}
	for k, v := range profile.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	lat := time.Since(start)
	if err != nil {
		col.Record(lat, 0, err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	col.Record(lat, resp.StatusCode, nil)
}