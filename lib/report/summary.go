package report

import "time"

type LatencyStats struct {
	MinMs float64 `json:"min_ms"`
	AvgMs float64 `json:"avg_ms"`
	MaxMs float64 `json:"max_ms"`
	P50Ms float64 `json:"p50_ms"`
	P95Ms float64 `json:"p95_ms"`
	P99Ms float64 `json:"p99_ms"`
}

type Summary struct {
	JobID        string       `json:"job_id"`
	JobName      string       `json:"job_name"`
	JobType      string       `json:"job_type"`
	StartedAt    time.Time    `json:"started_at"`
	EndedAt      time.Time    `json:"ended_at"`
	Duration     float64      `json:"duration_sec"`
	Total        int          `json:"total_requests"`
	Success      int          `json:"success"`
	Errors       int          `json:"errors"`
	HTTP2xxOr3xx int          `json:"http_2xx_3xx"`
	HTTP4xxOr5xx int          `json:"http_4xx_5xx"`
	RPS          float64      `json:"rps"`
	Latency      LatencyStats `json:"latency"`
	StatusCodes  map[int]int  `json:"status_codes"`
	ReportPath   string       `json:"report_path,omitempty"`
}