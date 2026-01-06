package service

import (
	"context"
	"loadsg/lib/dto"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
)

type loadManagerService struct{
	lrepo repository.LoadRepository
	hrepo repository.HttpLoadRepository
}

func NewLoadManagerService(lr repository.LoadRepository, hr repository.HttpLoadRepository) loadManagerService {
	return  &loadManagerService{lrepo: lr, hrepo: hr}
}

func (s *loadManagerService) CreateFixedHTTPLoadJob(ctx context.Context, 
	req dto.CreateFixedHTTPLoadRequest) (*model.LoadJob, error){

	loadJob := &model.LoadJob{
		JobName: req.JobName,
		Type: req.Type,
		StartTime: req.StartTime,
	}

	lj, err := s.lrepo.Create(ctx, loadJob)
	if err != nil {
		log.Printf("cant create loadJob: %s", err)
	}

	fixedHttpload := &model.FixedHttpLoad{
		LoadJobId: lj.Id,
		Payload: req.Payload,
		RequestCount: req.RequestCount,
	}
	hl, err := s.hrepo.CreateFixed(ctx, fixedHttpload)
	if err != nil{
		log.Printf("cant create http load: %s", err)
		return lj, err
	}
	return lj, nil
}
