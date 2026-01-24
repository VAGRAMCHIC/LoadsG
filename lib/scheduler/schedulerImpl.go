package scheduler

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
)

type scheduler struct {
	loadJobRepo repository.LoadRepository
	eventRepo   repository.EventRepository
}

func NewScheduler(lr repository.LoadRepository, ev repository.EventRepository) *scheduler {
	return &scheduler{loadJobRepo: lr, eventRepo: ev}
}

func (s *scheduler) Create(ctx context.Context) (event_ids []string, err error) {
	loadJobs, err := s.loadJobRepo.ScanLJob(ctx)
	if err != nil {
		log.Printf("Create: cant read load jobs: %s", err)
		return nil, err
	}
	if len(loadJobs) == 0 {
		log.Print("pass")
		return nil, nil
	}

	events := make([]model.Event, 0)

	for _, v := range loadJobs {
		var event model.Event
		event.LoadJobId = v.Id
		event.Status = "pending"

		events = append(events, event)
	}
	for _, v := range events {
		id, err := s.eventRepo.Create(ctx, &v)
		if err != nil {
			log.Printf("Create: cant create event: %s", err)
			return nil, err
		}
		event_ids = append(event_ids, id)
	}
	return event_ids, err
}

func (s *scheduler) ProcessPendingEvents(ctx context.Context) error {
	// 1. Получаем ближайшие load_jobs
	loadJobs, err := s.loadJobRepo.ScanClosestLJob(ctx)
	if err != nil {
		log.Fatalf("scan closest load jobs: %s", err)
		return err
	}

	for _, lj := range loadJobs {
		// 2. Получаем pending events для load_job
		events, err := s.eventRepo.ScanEvents(ctx, "pending")
		if err != nil {
			log.Fatalf("scan pending events for load job %s: %s", lj.Id, err)
			return err
		}

		for _, event := range events {
			// 3. Меняем статус
			event.Status = "running"

			// 4. Сохраняем event
			if err := s.eventRepo.UpdateEvent(ctx, &event); err != nil {
				log.Fatalf("update event %s: %s", event.Id, err)
				return err
			}
		}
	}
	return nil
}

func (s *scheduler) ProcessRunningEvents(ctx context.Context) error {

}
