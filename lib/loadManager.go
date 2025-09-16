package lib

import (
	"log"
	"sync"
	"time"
)

func RunLoad(load HTTPLoadRequest, mc int, endpoint string) (map[int]int, float64, float64, error) {
	request, host := BuildHttpRequest(load.HttpHead, load.Body)
	requests := make([]string, load.Count)

	var report Report
	report.StatusCodes = make(map[int]int)

	for i := 0; i < load.Count; i++ {
		requests[i] = request
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorCount int
	var completed int
	statusCodes := make(map[int]int)

	maxConcurrent := mc
	sem := make(chan struct{}, maxConcurrent)
	start := time.Now()

	// горутина-репортер
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				partial := Report{
					StatusCodes:  make(map[int]int),
					LoadDuration: time.Since(start).Seconds(),
					RPS:          float64(completed) / time.Since(start).Seconds(),
				}
				for k, v := range report.StatusCodes {
					partial.StatusCodes[k] = v
				}
				mu.Unlock()

				if err := sendReport(endpoint, partial); err != nil {
					log.Printf("Ошибка отправки отчета: %v", err)
				}

			case <-done:
				return
			}
		}
	}()

	for _, req := range requests {
		wg.Add(1)
		sem <- struct{}{}

		go func(r string) {
			defer wg.Done()
			defer func() { <-sem }()

			code, err := SendHttpRequest(r, host)

			mu.Lock()
			if err != nil {
				errorCount++
				log.Printf("Ошибка запроса: %v", err)
			} else {
				statusCodes[code]++
				report.StatusCodes[code]++
			}
			completed++
			mu.Unlock()
		}(req)
	}

	wg.Wait()
	close(done) // останавливаем репортинг

	elapsed := time.Since(start).Seconds()
	rps := float64(len(requests)) / elapsed
	report.LoadDuration = elapsed
	report.RPS = rps

	// отправляем финальный отчёт
	if err := sendReport(endpoint, report); err != nil {
		log.Printf("Ошибка отправки финального отчета: %v", err)
	}

	return statusCodes, elapsed, rps, nil
}
