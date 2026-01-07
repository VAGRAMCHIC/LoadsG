package service

import (
	"context"
	"loadsg/lib/dto"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
	"strconv"
	"time"
)

type loadManagerService struct {
	lrepo repository.LoadRepository
	hrepo repository.HttpLoadRepository
}

func NewLoadManagerService(lr repository.LoadRepository, hr repository.HttpLoadRepository) *loadManagerService {
	return &loadManagerService{lrepo: lr, hrepo: hr}
}

func (s *loadManagerService) CreateFixedHTTPLoadJob(ctx context.Context,
	req dto.CreateFixedHTTPLoadRequest) (*model.LoadJob, error) {

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("cant read date: %s", err)
	}
	requestCount, err := strconv.Atoi(req.RequestCount)
	if err != nil {
		log.Printf("cant read request count: %s", err)
	}

	loadJob := &model.LoadJob{
		JobName:   req.JobName,
		Type:      req.Type,
		StartTime: startTime,
	}

	lj, err := s.lrepo.Create(ctx, loadJob)
	if err != nil {
		log.Printf("cant create loadJob: %s", err)
	}

	fixedHttpload := &model.FixedHttpLoad{
		LoadJobId:    lj.Id,
		Payload:      req.Payload,
		RequestCount: requestCount,
	}
	hl, err := s.hrepo.CreateFixed(ctx, fixedHttpload)
	if err != nil {
		log.Printf("cant create http load: %s", err)
		log.Print(hl)
		return lj, err
	}
	return lj, nil
}

func (s *loadManagerService) DeleteFixedHTTPLoadJob(ctx context.Context, loadJobId string) (string, error) {
	log.Printf("load job id: %s", loadJobId)

	fJob, err := s.lrepo.GetById(ctx, loadJobId)
	if err != nil {
		log.Printf("cant find job: %s", err)
		return fJob.Id, err
	}
	err = s.lrepo.Delete(ctx, fJob.Id)
	if err != nil {
		log.Printf("cant delete job: %s", err)
		return fJob.Id, err
	}
	return fJob.Id, nil
}
