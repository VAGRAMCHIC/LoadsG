package report

import (
	"sort"
	"sync"
	"time"
)

type Sample struct {
	Latency    time.Duration
	StatusCode int
	Error      bool
	ErrMsg     string
	At         time.Time
}

type Collector struct {
	mu      sync.Mutex
	samples []Sample
	started time.Time
	ended   time.Time
}

func NewCollector() *Collector {
	return &Collector{
		samples: make([]Sample, 0, 1024),
		started: time.Now(),
	}
}

func (c *Collector) Record(latency time.Duration, statusCode int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := Sample{
		Latency:    latency,
		StatusCode: statusCode,
		At:         time.Now(),
	}
	if err != nil {
		s.Error = true
		s.ErrMsg = err.Error()
	}
	c.samples = append(c.samples, s)
}

func (c *Collector) Finish() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ended = time.Now()
}

func (c *Collector) Summary(jobID, jobName, jobType string) Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	sum := Summary{
		JobID:     jobID,
		JobName:   jobName,
		JobType:   jobType,
		StartedAt: c.started,
		EndedAt:   c.ended,
		StatusCodes: make(map[int]int),
	}
	if sum.EndedAt.IsZero() {
		sum.EndedAt = time.Now()
	}
	sum.Duration = sum.EndedAt.Sub(sum.StartedAt).Seconds()

	if len(c.samples) == 0 {
		return sum
	}

	latencies := make([]time.Duration, 0, len(c.samples))
	var totalLat time.Duration

	for _, s := range c.samples {
		sum.Total++
		if s.Error {
			sum.Errors++
		} else {
			sum.Success++
			sum.StatusCodes[s.StatusCode]++
			if s.StatusCode >= 200 && s.StatusCode < 400 {
				sum.HTTP2xxOr3xx++
			} else if s.StatusCode >= 400 {
				sum.HTTP4xxOr5xx++
			}
		}
		latencies = append(latencies, s.Latency)
		totalLat += s.Latency
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	sum.Latency.MinMs = float64(latencies[0]) / float64(time.Millisecond)
	sum.Latency.MaxMs = float64(latencies[len(latencies)-1]) / float64(time.Millisecond)
	sum.Latency.AvgMs = float64(totalLat) / float64(len(latencies)) / float64(time.Millisecond)
	sum.Latency.P50Ms = percentile(latencies, 50)
	sum.Latency.P95Ms = percentile(latencies, 95)
	sum.Latency.P99Ms = percentile(latencies, 99)

	if sum.Duration > 0 {
		sum.RPS = float64(sum.Total) / sum.Duration
	}

	return sum
}

func percentile(sorted []time.Duration, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return float64(sorted[0]) / float64(time.Millisecond)
	}
	rank := p / 100 * float64(len(sorted)-1)
	i := int(rank)
	frac := rank - float64(i)
	if i+1 >= len(sorted) {
		return float64(sorted[len(sorted)-1]) / float64(time.Millisecond)
	}
	v := float64(sorted[i])*(1-frac) + float64(sorted[i+1])*frac
	return v / float64(time.Millisecond)
}