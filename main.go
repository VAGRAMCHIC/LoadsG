package main

import (
	"fmt"
	"loadsg/lib"
	"loadsg/utils"
)

func main() {
	config := utils.Config{}

	//config, err := utils.ReadConfig("config.json")
	config, err := utils.ReadOSENV()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	fmt.Print(config)

	lib.Server([]byte(config.JwtKey), config.MaxConcurrent)

}
