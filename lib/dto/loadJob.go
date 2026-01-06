package dto

type CreateFixedHTTPLoadRequest struct {
	JobName				string 						`json:"job_name" binding:"required"`
	Type					string 						`json:"type" binding:"required"`
	StartTime 		string 						`json:"start_time" binding:"required"`
	RequestCount 	string						`json:"request_count" binding:"required"`
	Payload				map[string]string `json:"payload" binding:"required"`
}

type CreateFixedHTTPLoadResponse struct {
	JobID	string	`json:"id" binding:"required"`

}
 
