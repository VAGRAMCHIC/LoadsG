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
	var httoLoad model.FixedHttpLoad
	err := r.db.QueryRow(context.Background(), "SELECT id, load_job_id, request_count, payload FROM fixed_http_load where id=$1").Scan(&httoLoad.Id, &httoLoad.LoadJobId, &httoLoad.RequestCount, &httoLoad.Payload)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &httoLoad, err
	}
	return &httoLoad, nil
}

func (r *HttpLoadRepository) CreateFixed(ctx context.Context, httpLoad *model.FixedHttpLoad) (*model.FixedHttpLoad, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO fixed_http_load (load_job_id, request_count, payload) VALUES ($1, $2, $3) RETURNING id, load_job_id, request_count", httpLoad.LoadJobId, httpLoad.RequestCount, httpLoad.Payload).Scan(&httpLoad.Id, &httpLoad.LoadJobId, &httpLoad.RequestCount)
	if err != nil {
		log.Printf("cant create load job: %s", err)
		return httpLoad, err
	}
	return httpLoad, nil
}

func (r *HttpLoadRepository) DeleteFixed(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM fixed_http_load WHERE id=$1`, id)
	if err != nil {
		log.Printf("cant delete load job: %s", err)
		return err
	}
	return nil
}
