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
	err := r.db.QueryRow(ctx,
		`SELECT id, job_name, type, start_time
		FROM load_job where id=$1`, id).
		Scan(&loadJob.Id, &loadJob.JobName, &loadJob.Type, &loadJob.StartTime)
	if err != nil {
		log.Printf("cant get load job by id: %s", err)
		return &loadJob, err
	}
	return &loadJob, nil
}

func (r *LoadRepository) GetAllMatchById(ctx context.Context, ids []string) ([]model.LoadJob, error) {
	// Если список ID пустой, возвращаем пустой результат
	if len(ids) == 0 {
		return []model.LoadJob{}, nil
	}

	// Формируем запрос с параметром
	const query = `
        SELECT id, job_name, type, start_time
        FROM load_job
        WHERE id = ANY($1::uuid[])
        ORDER BY id`

	// Выполняем запрос с передачей массива UUID
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		log.Printf("query failed: %s", err)
		return nil, err
	}
	defer rows.Close()

	loadJobs := make([]model.LoadJob, 0)

	for rows.Next() {
		var job model.LoadJob
		if err := rows.Scan(
			&job.Id,
			&job.JobName,
			&job.Type,
			&job.StartTime,
		); err != nil {
			log.Printf("scan failed: %s", err)
			return nil, err
		}
		loadJobs = append(loadJobs, job)
	}
	if err := rows.Err(); err != nil {
		log.Printf("rows error: %s", err)
		return nil, err
	}
	return loadJobs, nil
}

func (r *LoadRepository) Create(ctx context.Context, loadJob *model.LoadJob) (*model.LoadJob, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO load_job (job_name, type, start_time)
		 VALUES ($1, $2, $3) RETURNING id, job_name, type, start_time`,
		loadJob.JobName, loadJob.Type, loadJob.StartTime).
		Scan(&loadJob.Id, &loadJob.JobName, &loadJob.Type, &loadJob.StartTime)
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

func (r *LoadRepository) ScanLJob(ctx context.Context) ([]model.LoadJob, error) {
	const query = `
	SELECT id, job_name, type, strart_time
	FROM load_job`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loadJobs := make([]model.LoadJob, 0)

	for rows.Next() {
		var job model.LoadJob
		if err := rows.Scan(
			&job.Id,
			&job.JobName,
			&job.Type,
			&job.StartTime,
		); err != nil {
			return nil, err
		}
		loadJobs = append(loadJobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return loadJobs, nil
}

func (r *LoadRepository) ScanClosestLJob(ctx context.Context) ([]model.LoadJob, error) {
	const query = `
	SELECT id, job_name, type, strart_time
	FROM load_job
	WHERE start_time >= NOW() - INTERVAL '10 seconds'`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loadJobs := make([]model.LoadJob, 0)

	for rows.Next() {
		var job model.LoadJob
		if err := rows.Scan(
			&job.Id,
			&job.JobName,
			&job.Type,
			&job.StartTime,
		); err != nil {
			return nil, err
		}
		loadJobs = append(loadJobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return loadJobs, nil
}
