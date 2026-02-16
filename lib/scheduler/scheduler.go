package scheduler

import (
	"context"
)

type Scheduler interface {
	Create(ctx context.Context) (event_ids []string, err error)
	RunEvents(ctx context.Context) error
}
