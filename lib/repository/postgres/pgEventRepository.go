package postgres

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PENDING_E = `
	SELECT id, load_job_id, status
	FROM events
	WHERE status='pending'`
	RUNNING_E = `
	SELECT id, load_job_id, status
	FROM events
	WHERE status='running'`
	PROCESSING_E = `
	SELECT id, load_job_id, status
	FROM events
	WHERE status='processing'`
	DONE_E = `
	SELECT id, load_job_id, status
	FROM events
	WHERE status='done'`
	FAILED_E = `
	SELECT id, load_job_id, status
	FROM events
	WHERE status='failed'`
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) repository.EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GetById(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.QueryRow(ctx,
		`SELECT id, load_job_id, status
		 FROM events where id=$1`, id).
		Scan(&event.Id, &event.LoadJobId, &event.Status)
	if err != nil {
		log.Printf("cant get event by id: %s", err)
		return &event, err
	}
	return &event, nil
}

func (r *EventRepository) Create(ctx context.Context, event *model.Event) (string, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO events (load_job_id, status)
		 VALUES ($1, $2) RETURNING id`,
		event.LoadJobId, event.Status).
		Scan(&event.Id)
	if err != nil {
		log.Printf("cant create event: %s", err)
		return event.Id, err
	}
	return event.Id, nil
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		log.Printf("cant delete event: %s", err)
		return err
	}
	return nil
}

func (r *EventRepository) UpdateEvent(
	ctx context.Context,
	event *model.Event,
) error {
	_, err := r.db.Exec(ctx, `UPDATE events SET status = $1 WHERE id=$2`, event.Status, event.Id)
	if err != nil {
		log.Printf("cant mark event as runnung: %s", err)
		return err
	}
	return nil
}

func (r *EventRepository) MarkDone(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE events SET status = 'done' WHERE id = $1`, id)
	return err
}

func (r *EventRepository) MarkFailed(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE events SET status = 'failed' WHERE id = $1`, id)
	return err
}

func (r *EventRepository) ScanEvents(ctx context.Context, status string) ([]model.Event, error) {
	query := ""
	switch status {
	case "pending":
		{
			query = PENDING_E
		}
	case "running":
		{
			query = RUNNING_E
		}
	case "processing":
		{
			query = PROCESSING_E
		}
	case "done":
		{
			query = DONE_E
		}
	case "failed":
		{
			query = FAILED_E
		}
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.Event, 0)

	for rows.Next() {
		var event model.Event
		if err := rows.Scan(
			&event.Id,
			&event.LoadJobId,
			&event.Status,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
