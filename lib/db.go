package lib

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"loadsg/lib/model"
)

func InitDB(ctx context.Context, pool *pgxpool.Pool) error {
	tables := []model.DB_TABLE{
		model.DB_TABLE_USERS,
		model.DB_TABLE_JWTTOKENS,
		model.DB_TABLE_LOAD_JOB,
		model.DB_TABLE_EVENTS,
		model.DB_TABLE_FIXED_HTTP_LOAD,
		model.DB_TABLE_CONSTANT_HTTP_LOAD,
		model.DB_TABLE_RAMPUP_HTTP_LOAD,
		model.DB_TABLE_FAKE_HTTP_LOAD,
	}

	for _, table := range tables {
		if err := initTable(ctx, pool, table); err != nil {
			return err
		}
	}
	if err := ensureEventsID(ctx, pool); err != nil {
		return err
	}
	return nil
}

func initTable(
	ctx context.Context,
	pool *pgxpool.Pool,
	table model.DB_TABLE,
) error {

	var keys []string
	for k := range table.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table.Name))

	// columns
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("    %s %s,\n", k, table.Params[k]))
	}

	// table-level primary key (если НЕ inline)
	if len(table.PrimaryKey) > 0 && !table.InlinePrimaryKey {
		b.WriteString(fmt.Sprintf(
			"    PRIMARY KEY (%s),\n",
			strings.Join(table.PrimaryKey, ", "),
		))
	}

	// foreign keys
	for _, fk := range table.ForeignKeys {
		b.WriteString(fmt.Sprintf(
			"    FOREIGN KEY (%s) REFERENCES %s(%s)",
			fk.Column,
			fk.RefTable,
			fk.RefColumn,
		))
		if fk.OnDelete != "" {
			b.WriteString(" ON DELETE " + fk.OnDelete)
		}
		b.WriteString(",\n")
	}

	sql := strings.TrimSuffix(b.String(), ",\n") + "\n);"

	_, err := pool.Exec(ctx, sql)
	return err
}

func ensureEventsID(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		ALTER TABLE events ADD COLUMN IF NOT EXISTS id UUID DEFAULT gen_random_uuid();
		UPDATE events SET id = gen_random_uuid() WHERE id IS NULL;
		ALTER TABLE events ALTER COLUMN id SET NOT NULL;
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conrelid = 'events'::regclass
				  AND contype = 'p'
			) THEN
				ALTER TABLE events ADD PRIMARY KEY (id);
			END IF;
		END $$;
	`)
	return err
}

func InitPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("database pool initialized")
	return pool, nil
}
