package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"loadsg/lib/dto"
)

func newLoadCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a load job",
	}
	cmd.AddCommand(newLoadCreateConstantCmd())
	cmd.AddCommand(newLoadCreateRampUpCmd())
	return cmd
}

func newLoadCreateConstantCmd() *cobra.Command {
	var (
		jobName   string
		startTime string
		count     int
		url       string
		method    string
		headers   string
		body      string
	)

	cmd := &cobra.Command{
		Use:   "constant",
		Short: "Create constant HTTP load job",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagToken == "" {
				return fmt.Errorf("missing token: use --token or LOADSG_TOKEN")
			}

			req := dto.CreateConstantHTTPLoadRequest{
				JobName:   jobName,
				Type:      "constant_http",
				StartTime: startTime,
				Count:     count,
				URL:       url,
				Method:    method,
			}

			if headers != "" {
				var h map[string]string
				if err := json.Unmarshal([]byte(headers), &h); err != nil {
					return fmt.Errorf("invalid --headers: %w", err)
				}
				req.Headers = h
			}
			if body != "" {
				req.Body = json.RawMessage(body)
			}

			c := newClient()
			resp, err := c.CreateConstantLoad(req)
			if err != nil {
				return err
			}

			fmt.Printf("Job created: %s\n", resp.JobID)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobName, "name", "", "Job name")
	cmd.Flags().StringVar(&startTime, "start", time.Now().UTC().Format(time.RFC3339), "Start time (RFC3339)")
	cmd.Flags().IntVar(&count, "count", 0, "Number of requests")
	cmd.Flags().StringVar(&url, "url", "", "Target URL")
	cmd.Flags().StringVar(&method, "method", "GET", "HTTP method")
	cmd.Flags().StringVar(&headers, "headers", "", `Headers as JSON, e.g. '{"Accept":"application/json"}'`)
	cmd.Flags().StringVar(&body, "body", "", "Request body as JSON string")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("count")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func newLoadCreateRampUpCmd() *cobra.Command {
	var (
		jobName   string
		startTime string
		duration  float32
		rpsS      int
		rpsF      int
		url       string
		method    string
		headers   string
		body      string
	)

	cmd := &cobra.Command{
		Use:   "rampup",
		Short: "Create ramp-up HTTP load job",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagToken == "" {
				return fmt.Errorf("missing token: use --token or LOADSG_TOKEN")
			}

			// DTO ожидает RPS_S / RPS_F как string (см. dto.CreateRampUpHTTPLoadRequest)
			req := dto.CreateRampUpHTTPLoadRequest{
				JobName:   jobName,
				Type:      "ramp_up_http",
				StartTime: startTime,
				RPS_S:     strconv.Itoa(rpsS),
				RPS_F:     strconv.Itoa(rpsF),
				Duration:  duration,
				URL:       url,
				Method:    method,
			}

			if headers != "" {
				var h map[string]string
				if err := json.Unmarshal([]byte(headers), &h); err != nil {
					return fmt.Errorf("invalid --headers: %w", err)
				}
				req.Headers = h
			}
			if body != "" {
				req.Body = json.RawMessage(body)
			}

			c := newClient()
			resp, err := c.CreateRampUpLoad(req)
			if err != nil {
				return err
			}

			fmt.Printf("Job created: %s\n", resp.JobID)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobName, "name", "", "Job name")
	cmd.Flags().StringVar(&startTime, "start", time.Now().UTC().Format(time.RFC3339), "Start time (RFC3339)")
	cmd.Flags().Float32Var(&duration, "duration", 0.0, "Load duration (seconds)")
	cmd.Flags().IntVar(&rpsS, "rps_s", 0, "Start RPS")
	cmd.Flags().IntVar(&rpsF, "rps_f", 0, "Finish RPS")
	cmd.Flags().StringVar(&url, "url", "", "Target URL")
	cmd.Flags().StringVar(&method, "method", "GET", "HTTP method")
	cmd.Flags().StringVar(&headers, "headers", "", `Headers as JSON, e.g. '{"Accept":"application/json"}'`)
	cmd.Flags().StringVar(&body, "body", "", "Request body as JSON string")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("duration")
	_ = cmd.MarkFlagRequired("rps_s")
	_ = cmd.MarkFlagRequired("rps_f")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}