package main

import (
	"context"
	"fmt"
	"loadsg/lib"
	"loadsg/lib/handler"
	"loadsg/lib/repository/postgres"
	"loadsg/lib/router"
	"loadsg/lib/security"
	"loadsg/lib/service"
	"loadsg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := utils.Config{}

	//config, err := utils.ReadConfig("config.json")
	config, err := utils.ReadOSENV()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	
	ctx:= context.Background()
	dbpool, err:= pgxpool.New(ctx, config.PgConn)
	if err != nil {
		log.Fatal(err)
	}

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
	r := gin.Default()

	userRepo := postgres.NewUserRepository(dbpool)
	jwtManager := security.NewJWTManager(config.JwtKey, 3600)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, *jwtManager)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	router.RegisterRoutes(r, userHandler, authHandler)

	log.Fatal(r.Run(":8080"))

}
