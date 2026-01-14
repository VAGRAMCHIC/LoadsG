package postgres

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoadRepository struct {
	db *pgxpool.Pool
}

func NewLoadRepository(db *pgxpool.Pool) repository.LoadRepository {
	return &LoadRepository{
		db: db,
	}
}

func (r *LoadRepository) GetById(ctx context.Context, id string) (*model.LoadJob, error) {
	var loadJob model.LoadJob
	err := r.db.QueryRow(ctx, "SELECT id, job_name, type, start_time, status FROM load_job where id=$1", id).Scan(&loadJob.Id, &loadJob.JobName, &loadJob.Type, &loadJob.StartTime, &loadJob.Status)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &loadJob, err
	}
	return &loadJob, nil
}

func (r *LoadRepository) Create(ctx context.Context, loadJob *model.LoadJob) (*model.LoadJob, error) {
	err := r.db.QueryRow(ctx, "INSERT INTO load_job (job_name, type, start_time) VALUES ($1, $2, $3) RETURNING id, job_name, type, start_time", loadJob.JobName, loadJob.Type, loadJob.StartTime).Scan(&loadJob.Id, &loadJob.JobName, &loadJob.Type, &loadJob.StartTime)
	if err != nil {
		log.Printf("cant create load job: %s", err)
		return loadJob, err
	}
	return loadJob, nil
}

func (r *LoadRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM load_job WHERE id=$1`, id)
	if err != nil {
		log.Printf("cant delete load job: %s", err)
		return err
	}
	return nil
}

func (r *LoadRepository) LockDueJobs(
	ctx context.Context,
	limit int,
) ([]model.LoadJob, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, job_name, type, start_time
		FROM load_job
		WHERE status = 'pending'
		  AND start_time <= now()
		ORDER BY start_time
		FOR UPDATE SKIP LOCKED
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.LoadJob
	for rows.Next() {
		var j model.LoadJob
		if err := rows.Scan(
			&j.Id,
			&j.JobName,
			&j.Type,
			&j.StartTime,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	for _, j := range jobs {
		_, err := tx.Exec(ctx, `
			UPDATE load_job
			SET status = 'processing',
			    locked_at = now()
			WHERE id = $1
		`, j.Id)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *LoadRepository) MarkDone(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE load_job SET status = 'done' WHERE id = $1`, id)
	return err
}

func (r *LoadRepository) MarkFailed(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE load_job SET status = 'failed' WHERE id = $1`, id)
	return err
}
