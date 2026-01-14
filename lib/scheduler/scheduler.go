package scheduler

import (
	"context"
	"loadsg/lib/generators"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
	"time"
)

type Handler func(context.Context, model.LoadJob) error

type Scheduler struct {
	repo        repository.LoadRepository
	registry    *generators.Registry
	interval    time.Duration
	workerLimit chan struct{}
}

func New(
	repo repository.LoadRepository,
	registry *generators.Registry,
	interval time.Duration,
	workers int,
) *Scheduler {
	return &Scheduler{
		repo:        repo,
		registry:    registry,
		interval:    interval,
		workerLimit: make(chan struct{}, workers),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.tick(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Scheduler) tick(ctx context.Context) {
	jobs, err := s.repo.LockDueJobs(ctx, 10)
	if err != nil {
		log.Println("scheduler error:", err)
		return
	}

	for _, job := range jobs {
		s.workerLimit <- struct{}{}
		go s.processJob(ctx, job)
	}
}
