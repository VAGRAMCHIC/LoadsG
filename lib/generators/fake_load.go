package generators

import (
	"context"
	"loadsg/lib/model"
	"log"
	"time"
)

type FakeHttpLoad struct{}

func (h *FakeHttpLoad) Name() string {
	return "fake_load"
}

func (h *FakeHttpLoad) Run(ctx context.Context, job model.LoadJob) error {

	timer := time.NewTimer(10)
	log.Printf("fake_load")
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
