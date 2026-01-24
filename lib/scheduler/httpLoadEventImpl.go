package scheduler

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
)

type event struct {
	loadJobRepo repository.LoadRepository
	eventRepo   repository.EventRepository
}

func NewEvent(lr repository.LoadRepository, ev repository.EventRepository) *event {
	return &event{loadJobRepo: lr, eventRepo: ev}
}

func (e *event) CreateHttpLoad(ctx context.Context) (event *model.Event, err error) {
	events, err := e.eventRepo.ScanEvents(ctx, "running")
	if err != nil {
		log.Printf("CreateHttpLoad - cant scan events: %s", err)
		return nil, err
	}
	ids := make([]string, 0)
	for _, v := range events {
		ids = append(ids, v.Id)
	}
	loads, err := e.loadJobRepo.GetAllMatchById(ctx, ids)

}
