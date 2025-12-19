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
	fmt.Print(config.PgConn)	
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
	user:= lib.User{
		Id: config.Id,
		Password: config.Key,
	}
	lib.InsertUser(conn, user)

	lib.Server([]byte(config.JwtKey), config.MaxConcurrent, conn)

}
