package scheduler

import (
	"context"
	"loadsg/lib/model"
)

type HttpLoadEvent interface {
	CreateHttpLoad(ctx context.Context) (event *model.Event, load *model.LoadJob, err error)
	CreateHttpLoadReport(ctx context.Context) (event *model.Event, err error)
}
