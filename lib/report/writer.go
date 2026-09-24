package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Enabled bool
	Format  string // json | text | csv
	Dir     string // empty → os.TempDir()
}

func Write(cfg Config, sum Summary) (string, error) {
	if !cfg.Enabled {
		return "", nil
	}

	format := strings.ToLower(cfg.Format)
	if format == "" {
		format = "json"
	}

	dir := cfg.Dir
	if dir == "" {
		dir = os.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("report dir: %w", err)
	}

	ts := time.Now().UTC().Format("20060102-150405")
	base := fmt.Sprintf("loadsg-%s-%s", safe(sum.JobID), ts)

	var path string
	var err error

	switch format {
	case "json":
		path = filepath.Join(dir, base+".json")
		err = writeJSON(path, sum)
	case "csv":
		path = filepath.Join(dir, base+".csv")
		err = writeCSV(path, sum)
	case "text", "txt":
		path = filepath.Join(dir, base+".txt")
		err = writeText(path, sum)
	default:
		return "", fmt.Errorf("unknown report format: %s", format)
	}
	if err != nil {
		return "", err
	}
	return path, nil
}

func safe(s string) string {
	if len(s) > 8 {
		s = s[:8]
	}
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, s)
}

func writeJSON(path string, sum Summary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(sum)
}

func writeText(path string, sum Summary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, `LoadsG Load Report
==================
Job ID:      %s
Job Name:    %s
Type:        %s
Started:     %s
Ended:       %s
Duration:    %.3fs

Requests:    %d
  Success:   %d
  Errors:    %d
  2xx/3xx:   %d
  4xx/5xx:   %d
RPS:         %.2f

Latency (ms)
  min:  %.2f
  avg:  %.2f
  p50:  %.2f
  p95:  %.2f
  p99:  %.2f
  max:  %.2f

Status codes:
`,
		sum.JobID, sum.JobName, sum.JobType,
		sum.StartedAt.Format(time.RFC3339),
		sum.EndedAt.Format(time.RFC3339),
		sum.Duration,
		sum.Total, sum.Success, sum.Errors,
		sum.HTTP2xxOr3xx, sum.HTTP4xxOr5xx,
		sum.RPS,
		sum.Latency.MinMs, sum.Latency.AvgMs,
		sum.Latency.P50Ms, sum.Latency.P95Ms, sum.Latency.P99Ms,
		sum.Latency.MaxMs,
	)
	if err != nil {
		return err
	}

	codes := make([]int, 0, len(sum.StatusCodes))
	for c := range sum.StatusCodes {
		codes = append(codes, c)
	}
	sort.Ints(codes)
	for _, c := range codes {
		if _, err := fmt.Fprintf(f, "  %d: %d\n", c, sum.StatusCodes[c]); err != nil {
			return err
		}
	}
	return nil
}

func writeCSV(path string, sum Summary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{"metric", "value"})
	rows := [][]string{
		{"job_id", sum.JobID},
		{"job_name", sum.JobName},
		{"job_type", sum.JobType},
		{"started_at", sum.StartedAt.Format(time.RFC3339)},
		{"ended_at", sum.EndedAt.Format(time.RFC3339)},
		{"duration_sec", fmt.Sprintf("%.3f", sum.Duration)},
		{"total_requests", fmt.Sprintf("%d", sum.Total)},
		{"success", fmt.Sprintf("%d", sum.Success)},
		{"errors", fmt.Sprintf("%d", sum.Errors)},
		{"http_2xx_3xx", fmt.Sprintf("%d", sum.HTTP2xxOr3xx)},
		{"http_4xx_5xx", fmt.Sprintf("%d", sum.HTTP4xxOr5xx)},
		{"rps", fmt.Sprintf("%.2f", sum.RPS)},
		{"latency_min_ms", fmt.Sprintf("%.2f", sum.Latency.MinMs)},
		{"latency_avg_ms", fmt.Sprintf("%.2f", sum.Latency.AvgMs)},
		{"latency_p50_ms", fmt.Sprintf("%.2f", sum.Latency.P50Ms)},
		{"latency_p95_ms", fmt.Sprintf("%.2f", sum.Latency.P95Ms)},
		{"latency_p99_ms", fmt.Sprintf("%.2f", sum.Latency.P99Ms)},
		{"latency_max_ms", fmt.Sprintf("%.2f", sum.Latency.MaxMs)},
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	codes := make([]int, 0, len(sum.StatusCodes))
	for c := range sum.StatusCodes {
		codes = append(codes, c)
	}
	sort.Ints(codes)
	for _, c := range codes {
		if err := w.Write([]string{fmt.Sprintf("status_%d", c), fmt.Sprintf("%d", sum.StatusCodes[c])}); err != nil {
			return err
		}
	}
	return w.Error()
}