package main

import (
	"context"
	"fmt"
	"loadsg/lib"
	"loadsg/lib/dto"
	"loadsg/lib/generators"
	"loadsg/lib/handler"
	"loadsg/lib/model"
	"loadsg/lib/repository/postgres"
	"loadsg/lib/router"
	"loadsg/lib/scheduler"
	"loadsg/lib/security"
	"loadsg/lib/service"
	"loadsg/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func handleLoadJob(ctx context.Context, job model.LoadJob) error {
	// здесь может быть генератор нагрузки,
	// запуск тестов, external calls и т.д.
	time.Sleep(500 * time.Millisecond)
	return nil
}

func main() {
	config := utils.Config{}

	//config, err := utils.ReadConfig("config.json")
	config, err := utils.ReadOSENV()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	ctx := context.Background()

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
	loadJobRepo := postgres.NewLoadRepository(dbpool)
	httpLoadRepo := postgres.NewHttpLoadRepository(dbpool)

	registry := generators.NewRegistry()
	registry.Register(&generators.ConstantHttp{})

	sch := scheduler.New(
		loadJobRepo,
		registry,
		1*time.Second,
		3,
	)
	sch.Start(ctx)

	jwtManager := security.NewJWTManager(config.JwtKey, config.JwtRefreshKey, config.AppName, 300)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtRepo, *jwtManager)
	loadManagerService := service.NewLoadManagerService(loadJobRepo, httpLoadRepo)

	handler := handler.NewHandler(userService, authService, loadManagerService)

	router.RegisterRoutes(r, handler, jwtManager)
	log.Print("test\n")
	rootUser, err := userRepo.GetByComment(ctx, config.RootUser)
	log.Print(rootUser)
	if err != nil {
		var createUser dto.CreateUserRequest
		createUser.Comment = config.RootUser
		createUser.Token = config.RootToken
		_, err := userService.Create(ctx, createUser)
		if err != nil {
			log.Fatalf("cant create root user: %s", err)
		}
	}
	log.Fatal(r.Run(":8080"))
}
