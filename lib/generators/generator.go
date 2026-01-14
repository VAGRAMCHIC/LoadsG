package generators

import (
	"context"

	"loadsg/lib/model"
)

type Generator interface {
	Name() string
	Run(ctx context.Context, job model.LoadJob) error
}
