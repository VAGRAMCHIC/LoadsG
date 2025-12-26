package lib

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

// --------- DB_TABLES ----------------

type DB_TABLE struct {
	Name   string            `json:"name" binding:"required"`
	Params map[string]string `json:"params" binding:"required"`
}

var DB_TABLE_USERS = DB_TABLE{
	Name: "users",
	Params: map[string]string{
		"id":       "SERIAL PRIMARY KEY",
		"uid": 			"TEXT NOT NULL",
		"password": "TEXT NOT NULL",
	},
}

var DB_TABLE_HTTP_LOAD_JOB = DB_TABLE{
	Name: "http_load_job",
	Params: map[string]string{
		"id":       	"SERIAL PRIMARY KEY",
		"job_name": 	"TEXT NOT NULL",
		"duration": 	"REAL",
		"type":				"TEXT NOT NULL",
		"payload":  	"JSONB NOT NULL",
		"start_time": "TIMESTAMPTZ",
	},
}


