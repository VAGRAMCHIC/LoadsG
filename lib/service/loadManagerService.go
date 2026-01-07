package service

import (
	"context"
	"loadsg/lib/dto"
	"loadsg/lib/model"
)

type LoadManagerService interface {
	CreateFixedHTTPLoadJob(ctx context.Context, loadJob dto.CreateFixedHTTPLoadRequest) (job *model.LoadJob, err error)
	DeleteFixedHTTPLoadJob(ctx context.Context, loadJobId string) (id string, err error)

	CreateConstantHTTPLoadJob(ctx context.Context, loadJob dto.CreateConstantHTTPLoadRequest) (job *model.LoadJob, err error)
	DeleteConstantHTTPLoadJob(ctx context.Context, loadJobId string) (id string, err error)

	CreateRampUpHTTPLoadJob(ctx context.Context, loadJob dto.CreateRampUpHTTPLoadRequest) (job *model.LoadJob, err error)
	DeleteRampUpHTTPLoadJob(ctx context.Context, loadJobId string) (id string, err error)
}
