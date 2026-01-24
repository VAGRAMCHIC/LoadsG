package dto

import "encoding/json"

type CreateFixedHTTPLoadRequest struct {
	JobName   string            `json:"job_name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	StartTime string            `json:"start_time" binding:"required"`
	RPS       string            `json:"rps" binding:"required"`
	Duration  float32           `json:"duration" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
}

type CreateFixedHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type CreateConstantHTTPLoadRequest struct {
	JobName   string            `json:"job_name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	StartTime string            `json:"start_time" binding:"required"`
	Count     int               `json:"count" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
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
	Duration  float32           `json:"duration" binding:"required"`
	URL       string            `json:"url" binding:"required"`
	Method    string            `json:"method" binding:"required"`
	Headers   map[string]string `json:"headers"`
	Body      json.RawMessage   `json:"body"`
}

type CreateRampUpHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type CreateFakeHTTPLoadRequest struct {
	JobName   string  `json:"job_name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	StartTime string  `json:"start_time" binding:"required"`
	Duration  float32 `json:"duration" binding:"required"`
}

type CreateFakeHTTPLoadResponse struct {
	JobID string `json:"id" binding:"required"`
}

type DeleteLoadJobRequest struct {
	JobID string `json:"id" binding:"required"`
}

type DeleteLoadJobResponse struct {
	JobID string `json:"id" binding:"required"`
}
