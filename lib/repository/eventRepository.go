package repository

import (
	"context"
	"errors"
	"loadsg/lib/model"
)

var ErrEventNotFound = errors.New("Event not found")

type EventRepository interface {
	GetById(ctx context.Context, id string) (*model.Event, error)
	Create(ctx context.Context, event *model.Event) (string, error)
	Delete(ctx context.Context, id string) error
	MarkDone(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string) error
	UpdateEvent(ctx context.Context, event *model.Event) error
	ScanEvents(ctx context.Context, status string) ([]model.Event, error)
}
