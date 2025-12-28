package tokenb

import (
	"time"
	"net/http"
)

type RateFunc func(elapsed time.Duration) float64

// -------- Token Bucket ------------

type Generator interface {
	Next() *http.Request
}

type Metric struct {
	Latency time.Duration
	Error   bool
}

// --------- Executor --------------

type Executor struct {
	client  *http.Client
	workers int
	jobs    chan *http.Request
	metrics chan Metric
}


// --------- Scheduler --------------

type Scheduler struct {
	bucket    *TokenBucket
	generator Generator
	executor  *Executor
	duration  time.Duration
}

