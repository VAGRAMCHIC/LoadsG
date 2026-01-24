package scheduler

import (
	"loadsg/lib/repository"
)

type event struct {
	loadJobRepo repository.LoadRepository
	eventRepo   repository.EventRepository
}

func NewEvent(lr repository.LoadRepository, ev repository.EventRepository) *event {
	return &event{loadJobRepo: lr, eventRepo: ev}
}
