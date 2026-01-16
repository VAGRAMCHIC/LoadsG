package model

import "encoding/json"

type FixedHttpLoad struct {
	Id        string            `json:"id" binding:"required"`
	LoadJobId string            `json:"load_job_id" binding:"required"`
	RPS       int               `json:"rps" binding:"required"`
	Duration  float32           `json:"duration" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
}

var DB_TABLE_FIXED_HTTP_LOAD = DB_TABLE{
	Name: "fixed_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"duration":    "REAL NOT NULL",
		"rps":         "INTEGER NOT NULL",
		"url":         "TEXT NOT NULL",
		"method":      "TEXT NOT NULL",
		"headers":     "JSONB",
		"body":        "JSONB",
	},
	PrimaryKey: []string{"load_job_id"},
	ForeignKeys: []ForeignKey{
		{
			Column:    "load_job_id",
			RefTable:  "load_job",
			RefColumn: "id",
			OnDelete:  "CASCADE",
		},
	},
}
