package model

import "encoding/json"

type RampUpHttpLoad struct {
	Id        string            `json:"id" binding:"required"`
	LoadJobId string            `json:"load_job_id" binding:"required"`
	RPS_S     int               `json:"rps_s" binding:"required"`
	RPS_F     int               `json:"rps_f" binding:"required"`
	Duration  float32           `json:"duration" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
}

var DB_TABLE_RAMPUP_HTTP_LOAD = DB_TABLE{
	Name: "ramp_up_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"duration":    "REAL NOT NULL",
		"rps_s":       "INTEGER NOT NULL",
		"rps_f":       "INTEGER NOT NULL",
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
