package lib

import (
	"log"
	"sync"
	"time"
)

func RunLoad(load HTTPLoadRequest, mc int) (map[int]int, float64, float64, error) {
	request, host := BuildHttpRequest(load.HttpHead, load.Body)
	requests := make([]string, load.Count)

	for i := range load.Count {
		requests[i] = request
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorCount int
	statusCodes := make(map[int]int)

	maxConcurrent := mc
	sem := make(chan struct{}, maxConcurrent)
	start := time.Now()
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
			}
			mu.Unlock()
		}(req)
	}
	wg.Wait()

	elapsed := time.Since(start).Seconds()
	rps := float64(len(requests)) / elapsed

	return statusCodes, elapsed, rps, nil
}
