package lib

import "time"

type _httpMethod string


const (
	GET     _httpMethod = "GET"
	POST    _httpMethod = "POST"
	PUT     _httpMethod = "PUT"
	DELETE  _httpMethod = "DELETE"
	CONNECT _httpMethod = "CONNECT"
	PATCH   _httpMethod = "PATCH"
	OPTIONS _httpMethod = "OPTIONS"
	HEAD    _httpMethod = "HEAD"
)

type HttpHead struct {
	Method       _httpMethod       `json:"method"`
	URL          string            `json:"url"`
	ProtoVersion string            `json:"proto_version"`
	Length       int               `json:"length"`
	Headers      map[string]string `json:"headers"`
}

type Loader struct {
	Id              int    `json:"id" binding:"required"`
	Token           string `json:"token" binding:"required"`
	Status          bool   `json:"status"`
	Max_concurrency int    `json:"max_concurrency" binding:"required"`
	Download_link   string `json:"download_link"`
}

type HTTPLoadRequest struct {
	Id       int      `json:"id" binding:"required"`
	HttpHead HttpHead `json:"httpHead" binding:"required"`
	Body     string   `json:"body" binding:"required"`
	Count    int      `json:"count" binding:"required"`
}

type User struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Report struct {
	LoadDuration float64           `json:"loadDuration" binding:"required"`
	RPS          float64           `json:"rps" binding:"required"`
	StatusCodes  map[int]int       `json:"statusCodes" binding:"required"`
	Errors       map[string]string `json:"errors" binding:"required"`
}

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
		"username": "TEXT NOT NULL",
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


