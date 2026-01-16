package postgres

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HttpLoadRepository struct {
	db *pgxpool.Pool
}

func NewHttpLoadRepository(db *pgxpool.Pool) repository.HttpLoadRepository {
	return &HttpLoadRepository{
		db: db,
	}
}

func (r *HttpLoadRepository) GetFixedById(ctx context.Context, id string) (*model.FixedHttpLoad, error) {
	var httpLoad model.FixedHttpLoad
	err := r.db.QueryRow(context.Background(), "SELECT load_job_id, rps, duration, payload FROM fixed_http_load where load_job_id=$1", id).
		Scan(&httpLoad.LoadJobId, &httpLoad.RPS, &httpLoad.Duration, &httpLoad.Payload)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &httpLoad, err
	}
	return &httpLoad, nil
}

func (r *HttpLoadRepository) CreateFixed(ctx context.Context, httpLoad *model.FixedHttpLoad) (*model.FixedHttpLoad, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO fixed_http_load (load_job_id, rps, duration, payload) VALUES ($1, $2, $3, $4) RETURNING load_job_id, rps, duration",
		httpLoad.LoadJobId, httpLoad.RPS, httpLoad.Duration, httpLoad.Payload).
		Scan(&httpLoad.LoadJobId, &httpLoad.RPS, &httpLoad.Duration)
	if err != nil {
		log.Printf("cant create load job: %s", err)
		return httpLoad, err
	}
	return httpLoad, nil
}

func (r *HttpLoadRepository) DeleteFixed(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM fixed_http_load WHERE load_job_id=$1`, id)
	if err != nil {
		log.Printf("cant delete load job: %s", err)
		return err
	}
	return nil
}

func (r *HttpLoadRepository) GetConstantById(ctx context.Context, id string) (*model.ConstantHttpLoad, error) {
	var httpLoad model.ConstantHttpLoad
	err := r.db.QueryRow(ctx, "SELECT load_job_id, count, url, method, headers, body FROM constant_http_load where load_job_id=$1", id).
		Scan(&httpLoad.LoadJobId, &httpLoad.Count, &httpLoad.URL, &httpLoad.Method, &httpLoad.Headers, &httpLoad.Body)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &httpLoad, err
	}
	return &httpLoad, nil
}

func (r *HttpLoadRepository) CreateConstant(ctx context.Context, httpLoad *model.ConstantHttpLoad) (*model.ConstantHttpLoad, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO constant_http_load (load_job_id, count, url, method, headers, body) VALUES ($1, $2, $3, $4, $5, $6) RETURNING load_job_id, count, url",
		httpLoad.LoadJobId, httpLoad.Count, httpLoad.URL, httpLoad.Method, httpLoad.Headers, httpLoad.Body).
		Scan(&httpLoad.LoadJobId, &httpLoad.Count, &httpLoad.URL)
	if err != nil {
		log.Printf("cant create load job: %s", err)
		return httpLoad, err
	}
	return httpLoad, nil
}

func (r *HttpLoadRepository) DeleteConstant(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM constant_http_load WHERE load_job_id=$1`, id)
	if err != nil {
		log.Printf("cant delete load job: %s", err)
		return err
	}
	return nil
}

func (r *HttpLoadRepository) GetRampUpById(ctx context.Context, id string) (*model.RampUpHttpLoad, error) {
	var httpLoad model.RampUpHttpLoad
	err := r.db.QueryRow(context.Background(), "SELECT load_job_id, rps_s, rps_f, duration, payload FROM ramp_up_http_load where load_job_id=$1", id).
		Scan(&httpLoad.LoadJobId, &httpLoad.RPS_S, &httpLoad.RPS_F, &httpLoad.Duration, &httpLoad.Payload)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &httpLoad, err
	}
	return &httpLoad, nil
}

func (r *HttpLoadRepository) CreateRampUp(ctx context.Context, httpLoad *model.RampUpHttpLoad) (*model.RampUpHttpLoad, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO ramp_up_http_load (load_job_id, rps_s, rps_f, duration, payload) VALUES ($1, $2, $3, $4, $5) RETURNING load_job_id, rps_s, rps_f",
		httpLoad.LoadJobId, httpLoad.RPS_S, httpLoad.RPS_F, httpLoad.Duration, httpLoad.Payload).
		Scan(&httpLoad.LoadJobId, &httpLoad.RPS_S, &httpLoad.RPS_F)
	if err != nil {
		log.Printf("cant create load job: %s", err)
		return httpLoad, err
	}
	return httpLoad, nil
}

func (r *HttpLoadRepository) DeleteRampUp(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ramp_up_http_load WHERE load_job_id=$1`, id)
	if err != nil {
		log.Printf("cant delete load job: %s", err)
		return err
	}
	return nil
}
