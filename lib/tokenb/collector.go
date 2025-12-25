package tokenb


import (
	"fmt"
	"time"
)

func Collect(metrics <-chan Metric) {
	var count, errors int
	var total time.Duration

	for m := range metrics {
		count++
		total += m.Latency
		if m.Error {
			errors++
		}

		if count%1000 == 0 {
			fmt.Printf(
				"req=%d avg=%v errors=%d\n",
				count,
				total/time.Duration(count),
				errors,
			)
		}
	}
}

