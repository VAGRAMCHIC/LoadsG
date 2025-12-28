package tokenb

import (
	"net/http"
	"time"
)



func NewExecutor(workers int, metrics chan Metric) *Executor {
	tr := &http.Transport{
		MaxIdleConnsPerHost: workers,
		DisableCompression: true,
	}

	e := &Executor{
		client: &http.Client{Transport: tr},
		workers: workers,
		jobs:    make(chan *http.Request, workers*2),
		metrics: metrics,
	}

	for i := 0; i < workers; i++ {
		go e.worker()
	}

	return e
}

func (e *Executor) worker() {
	for req := range e.jobs {
		start := time.Now()

		resp, err := e.client.Do(req)

		lat := time.Since(start)
		e.metrics <- Metric{
			Latency: lat,
			Error:   err != nil,
		}

		if err == nil {
			resp.Body.Close()
		}
	}
}

