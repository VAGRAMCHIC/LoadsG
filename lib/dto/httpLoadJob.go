package dto

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
