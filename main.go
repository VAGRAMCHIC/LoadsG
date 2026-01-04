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

	
	dbpool, err := lib.InitPool(ctx, config.PgConn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	if err := lib.InitDB(ctx, dbpool); err != nil {
		log.Fatal(err)
	}


	r := gin.Default()

	userRepo := postgres.NewUserRepository(dbpool)
	jwtRepo := postgres.NewJWTRefreshRepository(dbpool)
	jwtManager := security.NewJWTManager(config.JwtKey, config.JwtRefreshKey, config.AppName, 300)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtRepo, *jwtManager)

	handler := handler.NewHandler(userService, authService)


	router.RegisterRoutes(r, handler, jwtManager)

	log.Fatal(r.Run(":8080"))

}
