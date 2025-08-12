package lib

import (
	"fmt"
	"log"
	"sync"
)

func RunLoad(load HTTPLoadRequest) (string, error) {
	request, host := BuildHttpRequest(load.HttpHead, load.Body)
	requests := make([]string, load.Count)
	for i := range load.Count {
		requests[i] = request
	}
	var wg sync.WaitGroup
	for _, req := range requests {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			SendHttpRequest(r, host)
			log.Print(r, host)
		}(req)
	}
	wg.Wait()

	return fmt.Sprintf("%d", len(requests)), nil
}
