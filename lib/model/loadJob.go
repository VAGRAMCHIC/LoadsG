package model

import (
	"time"
)

type LoadJob struct {
	Id        string    `json:"id" binding:"required"`
	JobName   string    `json:"jobName" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time" binding:"required"`
}

var DB_TABLE_LOAD_JOB = DB_TABLE{
	Name: "load_job",
	Params: map[string]string{
		"id":         "UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"job_name":   "TEXT NOT NULL",
		"type":       "TEXT NOT NULL",
		"status":     "TEXT NOT NULL DEFAULT 'pending'",
		"start_time": "TIMESTAMPTZ",
		"locked_at":  "TIMESTAMPTZ",
	},
	InlinePrimaryKey: true,
}
