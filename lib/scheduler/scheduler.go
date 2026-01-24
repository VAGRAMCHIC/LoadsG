package scheduler

import (
	"context"
	"loadsg/lib/model"
)

type Scheduler interface {
	Create(ctx context.Context) (event *model.Event, err error)
	AddToQueue(ctx context.Context) (event *model.Event, job *model.LoadJob, err error)
}
