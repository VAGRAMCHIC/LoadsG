package main

import (
	"context"
	"fmt"
	"loadsg/lib"
	"loadsg/lib/dto"
	"loadsg/lib/generators"
	"loadsg/lib/handler"
	"loadsg/lib/repository/postgres"
	"loadsg/lib/router"
	"loadsg/lib/scheduler"
	"loadsg/lib/security"
	"loadsg/lib/service"
	"loadsg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
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
	eventRepo := postgres.NewEventRepository(dbpool)

	// Регистрируем генераторы
	registry := generators.NewRegistry()
	registry.Register(&generators.FakeHttpLoad{})
	constantGenerator := generators.NewConstantHttp(httpLoadRepo)
    registry.Register(constantGenerator)
    registry.RegisterName("constant_http", constantGenerator)
	rampUpGenerator := generators.NewRampUpHttp(httpLoadRepo)
	registry.Register(rampUpGenerator)
	registry.RegisterName("ramp_up_http_load", rampUpGenerator)
	registry.RegisterName("ramp-up", rampUpGenerator)

	// Создаём планировщик с интервалом 5 секунд
	sched := scheduler.NewScheduler(
		loadJobRepo,
		eventRepo,
		httpLoadRepo,
		registry,
		5*time.Second,
	)

	// Запускаем планировщик в фоновой горутине
	ctxSched, cancelSched := context.WithCancel(context.Background())
	go func() {
		if err := sched.Run(ctxSched); err != nil && err != context.Canceled {
			log.Printf("Scheduler stopped with error: %v", err)
		}
	}()

	jwtManager := security.NewJWTManager(config.JwtKey, config.JwtRefreshKey, config.AppName, 300)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtRepo, *jwtManager)
	loadManagerService := service.NewLoadManagerService(loadJobRepo, httpLoadRepo)

	handler := handler.NewHandler(userService, authService, loadManagerService)

	router.RegisterRoutes(r, handler, jwtManager)

	// Создаём root пользователя
	rootUser, err := userRepo.GetByComment(ctx, config.RootUser)
	if err != nil {
		var createUser dto.CreateUserRequest
		createUser.Comment = config.RootUser
		createUser.Token = config.RootToken
		_, err := userService.Create(ctx, createUser)
		if err != nil {
			log.Fatalf("cant create root user: %s", err)
		}
		// После создания получаем пользователя заново, чтобы вывести его UID
		rootUser, err = userRepo.GetByComment(ctx, config.RootUser)
		if err != nil {
			log.Printf("warning: root user created but cannot retrieve UID: %v", err)
		} else {
			log.Printf("Root user created with UID: %s", rootUser.UID)
		}
	} else {
		log.Printf("Root user already exists with UID: %s", rootUser.UID)
	}

	// Запускаем HTTP сервер
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Останавливаем планировщик
	cancelSched()
	// Даём время на завершение генераторов (можно добавить ожидание)
	time.Sleep(2 * time.Second)
	log.Println("Server stopped")
}

