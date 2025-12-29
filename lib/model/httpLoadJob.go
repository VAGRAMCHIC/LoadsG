package model

import (
	"time"
)

type HTTPLoadJob struct{
	Id 					string 							`json:"id" binding:"required"`
	JobName 		string 							`json:"jobName" binding:"required"`
	Duration 		float32 						`json:"duration" binding:"required"`
	Type 				string 							`json:"type" binding:"required"`
	Payload 		map[string]string 	`json:"payload" binding:"required"`
	StartTime 	time.Time 					`json:"start_time" binding:"required"`
}	

var DB_TABLE_HTTP_LOAD_JOB = DB_TABLE{
	Name: "http_load_job",
	Params: map[string]string{
		"id":       	"UUID NOT NULL",
		"job_name": 	"TEXT NOT NULL",
		"duration": 	"REAL",
		"type":				"TEXT NOT NULL",
		"payload":  	"JSONB NOT NULL",
		"start_time": "TIMESTAMPTZ",
	},
	PrimaryKey: []string{"id"},
}


