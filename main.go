package main

import (
	"os"
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

	conn:= lib.Connect(config.PgConn)
	db_status, err := lib.InitDB(conn)
	
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}

	if db_status != true {
		fmt.Println("Ошибка инициализации базы данных:", db_status)
		return
	}

	lib.Server([]byte(config.JwtKey), config.MaxConcurrent, conn)

}
