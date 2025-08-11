package main

import (
	"fmt"
	"loadsg/lib"
	"loadsg/utils"
	//"sync"
	//"time"
)

func main() {
	config := utils.Config{}

	config, err := utils.ReadConfig("config.json")
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	fmt.Print(config)

	lib.Server([]byte(config.JwtKey))

	/* ======================= TEST LOAD_REQUESTS ============================

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "LoadsG/1.0",
	}
	fmt.Print("\n")

	start := time.Now() // Засекаем время

	head := lib.CreateHttpHead("GET", "http://test.customlabs.ru/test2/", "HTTP/1.1", headers)

	request, host := lib.BuildHttpRequest(head, "")
	fmt.Print(request)

	requests := []string{request, request, request, request, request}
	var wg sync.WaitGroup
	for _, req := range requests {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			lib.SendHttpRequest(r, host)
		}(req)
	}
	wg.Wait()

	elapsed := time.Since(start) // Считаем разницу
	fmt.Printf("Время выполнения: %v\n", elapsed)
	*/
}
