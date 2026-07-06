package generators

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
	"net/http"
	"strings"
	"time"
)

type RampUpHttp struct {
	repo   repository.HttpLoadRepository
	client *http.Client
}

func NewRampUpHttp(repo repository.HttpLoadRepository) *RampUpHttp {
	return &RampUpHttp{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (h *RampUpHttp) Name() string {
	return "ramp_up_http"
}

func (h *RampUpHttp) Run(ctx context.Context, job model.LoadJob) error {
	profile, err := h.repo.GetRampUpById(ctx, job.Id)
	if err != nil {
		return fmt.Errorf("load ramp-up profile: %w", err)
	}

	duration := time.Duration(float64(profile.Duration) * float64(time.Second))
	if duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	if profile.RPS_S < 0 || profile.RPS_F < 0 {
		return fmt.Errorf("rps values must be non-negative")
	}

	startedAt := time.Now()
	for {
		elapsed := time.Since(startedAt)
		if elapsed >= duration {
			return nil
		}

		rps := rampRPS(profile.RPS_S, profile.RPS_F, elapsed, duration)
		window := time.Second
		remaining := duration - elapsed
		if remaining < window {
			window = remaining
		}

		if rps <= 0 {
			if err := sleepContext(ctx, window); err != nil {
				return err
			}
			continue
		}

		if err := h.runWindow(ctx, profile, rps, window); err != nil {
			return err
		}
	}
}

func (h *RampUpHttp) runWindow(ctx context.Context, profile *model.RampUpHttpLoad, rps int, window time.Duration) error {
	interval := time.Second / time.Duration(rps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	deadline := time.NewTimer(window)
	defer deadline.Stop()

	for sent := 0; sent < rps; sent++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return nil
		case <-ticker.C:
			if err := h.send(ctx, profile); err != nil {
				log.Printf("ramp_up_http request failed for job %s: %v", profile.LoadJobId, err)
			}
		}
	}

	return nil
}

func (h *RampUpHttp) send(ctx context.Context, profile *model.RampUpHttpLoad) error {
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

func rampRPS(start, end int, elapsed, duration time.Duration) int {
	if elapsed >= duration {
		return end
	}
	progress := float64(elapsed) / float64(duration)
	return start + int(float64(end-start)*progress)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func bodyReader(body []byte) io.Reader {
	if len(body) == 0 {
		return nil
	}
	return bytes.NewReader(body)
}
