package scheduler

import (
	"context"
	"log"
	"time"

	"loadsg/lib/model"
)

func (s *Scheduler) processJob(
	parentCtx context.Context,
	job model.LoadJob,
) {
	defer func() { <-s.workerLimit }()

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Minute)
	defer cancel()

	algo, err := s.registry.Get(job.Type)
	if err != nil {
		log.Println("unknown job type:", job.Type)
		_ = s.repo.MarkFailed(ctx, job.Id)
		return
	}

	if err := algo.Run(ctx, job); err != nil && err != context.Canceled {
		log.Println("job failed:", job.Id, err)
		_ = s.repo.MarkFailed(ctx, job.Id)
		return
	}

	_ = s.repo.MarkDone(ctx, job.Id)
}
