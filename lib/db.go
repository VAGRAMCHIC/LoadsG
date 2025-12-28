package lib

import (
	"context"
	"log"
	"fmt"
	"sort"
	"strings"
	"github.com/jackc/pgx/v5"

	"loadsg/lib/model"
)

func InitDB(conn *pgx.Conn) (bool, error) {
	code, err := initTable(conn, DB_TABLE_USERS)
	code, err = initTable(conn, DB_TABLE_HTTP_LOAD_JOB)
	if err != nil {
		return code, err
	}
	return code, err
}

func initTable(conn *pgx.Conn, table DB_TABLE) (bool, error) {
	var init_table_line string

	var keys []string
	for k := range table.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table.Name))

	for i, k := range keys {
		b.WriteString(fmt.Sprintf("    %s %s", k, table.Params[k]))
		if i < len(keys)-1 {
			b.WriteString(",\n") // запятая только между полями
		} else {
			b.WriteString("\n") // без запятой на последней строке
		}
	}
	b.WriteString(");")

	init_table_line = b.String()
	fmt.Printf("Create table %s: %s", table.Name, init_table_line)
	_, err := conn.Exec(context.Background(), init_table_line)
	//defer conn.Close(context.Background())
	if err != nil {
		log.Printf("cant create table: %s", err.Error())
		return false, err
	}
	return true, err
}



func Connect(pgConn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), pgConn)
	if err != nil {
		log.Printf("cant create connection: %s", err.Error())
	}
	return conn
}



// --------------- DB_TABLE_HTTP_LOAD_JOB ------------------

func InsertHTTPLoadJob(conn *pgx.Conn, loadJob model.HTTPLoadJob) {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO http_load_job (job_name, duration, type, payload, start_time) VALUES ($1, $2, $3, $4, $5)", 
													loadJob.JobName, loadJob.Duration, loadJob.Type, loadJob.Payload, loadJob.StartTime)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
		defer conn.Close(context.Background())
	}
	defer conn.Close(context.Background())
}


func GetHTTPLoadJob(conn *pgx.Conn, id string) (model.HTTPLoadJob, error) {
	var loadJob model.HTTPLoadJob
	err := conn.QueryRow(context.Background(), "SELECT * FROM http_load_job where id=$1", id).Scan(&loadJob.JobName, &loadJob.Duration, &loadJob.Type, &loadJob.Payload, &loadJob.StartTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			// load job not found
			return loadJob, nil
		}
		return loadJob, err
	}
	defer conn.Close(context.Background())
	return loadJob, err
}


