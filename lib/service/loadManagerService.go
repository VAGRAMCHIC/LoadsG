package service

import (
	"context"
	"loadsg/lib/dto"
	"loadsg/lib/model"
)

type LoadManagerService interface {
	CreateFixedHTTPLoadJob(ctx context.Context, 
		loadJob dto.CreateFixedHTTPLoadRequest) (job *model.LoadJob, err error)
}
