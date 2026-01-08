package dto

type CreateFixedHTTPLoadRequest struct {
	JobName   string            `json:"job_name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	StartTime string            `json:"start_time" binding:"required"`
	RPS       string            `json:"rps" binding:"required"`
	Duration  string            `json:"duration" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

type CreateFixedHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type CreateConstantHTTPLoadRequest struct {
	JobName   string            `json:"job_name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	StartTime string            `json:"start_time" binding:"required"`
	Count     string            `json:"count" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

type CreateConstantHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type CreateRampUpHTTPLoadRequest struct {
	JobName   string            `json:"job_name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	StartTime string            `json:"start_time" binding:"required"`
	RPS_S     string            `json:"rps_s" binding:"required"`
	RPS_F     string            `json:"rps_f" binding:"required"`
	Duration  string            `json:"duration" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

type CreateRampUpHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type DeleteLoadJobRequest struct {
	JobID string `json:"id" binding:"required"`
}

type DeleteLoadJobResponse struct {
	JobID string `json:"id" binding:"required"`
}
