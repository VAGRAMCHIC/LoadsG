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
	rps, err := strconv.Atoi(req.RPS)

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
		LoadJobId: lj.Id,
		URL:       req.URL,
		Method:    req.Method,
		Headers:   req.Headers,
		Body:      req.Body,
		RPS:       rps,
		Duration:  req.Duration,
	}
	hl, err := s.hrepo.CreateFixed(ctx, fixedHttpload)
	if err != nil {
		log.Printf("cant create http load: %s", err)
		log.Print(hl)
		return lj, err
	}
	return lj, nil
}

func (s *loadManagerService) CreateConstantHTTPLoadJob(ctx context.Context,
	req dto.CreateConstantHTTPLoadRequest) (*model.LoadJob, error) {

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("cant read date: %s", err)
	}
	//count, err := strconv.Atoi(req.Count)
	//if err != nil {
	//	log.Printf("cant read request count: %s", err)
	//}

	loadJob := &model.LoadJob{
		JobName:   req.JobName,
		Type:      req.Type,
		StartTime: startTime,
	}

	lj, err := s.lrepo.Create(ctx, loadJob)
	if err != nil {
		log.Printf("cant create loadJob: %s", err)
	}

	constantHttpload := &model.ConstantHttpLoad{
		LoadJobId: lj.Id,
		Count:     req.Count,
		URL:       req.URL,
		Method:    req.Method,
		Headers:   req.Headers,
		Body:      req.Body,
	}
	hl, err := s.hrepo.CreateConstant(ctx, constantHttpload)
	if err != nil {
		log.Printf("cant create http load: %s", err)
		log.Print(hl)
		return lj, err
	}
	return lj, nil
}

func (s *loadManagerService) CreateRampUpHTTPLoadJob(ctx context.Context,
	req dto.CreateRampUpHTTPLoadRequest) (*model.LoadJob, error) {

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("cant read date: %s", err)
	}
	rps_s, err := strconv.Atoi(req.RPS_S)
	rps_f, err := strconv.Atoi(req.RPS_F)
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

	rampUpHttpload := &model.RampUpHttpLoad{
		LoadJobId: lj.Id,
		URL:       req.URL,
		Method:    req.Method,
		Headers:   req.Headers,
		Body:      req.Body,
		RPS_S:     rps_s,
		RPS_F:     rps_f,
		Duration:  req.Duration,
	}
	hl, err := s.hrepo.CreateRampUp(ctx, rampUpHttpload)
	if err != nil {
		log.Printf("cant create http load: %s", err)
		log.Print(hl)
		return lj, err
	}
	return lj, nil
}

func (s *loadManagerService) CreateFakeHTTPLoadJob(ctx context.Context,
	req dto.CreateFakeHTTPLoadRequest) (*model.LoadJob, error) {

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("cant read date: %s", err)
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

	fakeHttpload := &model.FakeHttpLoad{
		LoadJobId: lj.Id,
		Duration:  req.Duration,
	}
	hl, err := s.hrepo.CreateFake(ctx, fakeHttpload)
	if err != nil {
		log.Printf("cant create http load: %s", err)
		log.Print(hl)
		return lj, err
	}
	return lj, nil
}

func (s *loadManagerService) DeleteLoadJob(ctx context.Context, loadJobId string) (string, error) {
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
