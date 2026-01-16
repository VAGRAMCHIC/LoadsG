package model

import "encoding/json"

type ConstantHttpLoad struct {
	Id        string            `json:"id" binding:"required"`
	LoadJobId string            `json:"load_job_id" binding:"required"`
	Count     int               `json:"count" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
}

var DB_TABLE_CONSTANT_HTTP_LOAD = DB_TABLE{
	Name: "constant_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"count":       "INTEGER NOT NULL",
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
