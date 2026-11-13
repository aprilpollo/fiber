package main

import (
	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/routes"
	"aprilpollo/internal/adapters/routes/handler"
	"aprilpollo/internal/adapters/storage/cache"
	"aprilpollo/internal/adapters/storage/orm"
	"aprilpollo/internal/adapters/storage/repository"
	"aprilpollo/internal/core/services"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✔ [INFO] Loading Configuration")

	db, err := orm.NewGormDB(cfg.Database, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("✔ [INFO] Database Connection")

	redis, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	fmt.Println("✔ [INFO] Redis Connection")

	// --- Repositories (output adapters) ---
	userRepo := repository.NewUserRepository(db.GetDB())
	bookRepo := repository.NewBookRepository(db.GetDB())

	// --- Services (core / use cases) ---
	userSvc := services.NewUserService(userRepo)
	bookSvc := services.NewBookService(bookRepo)

	// --- Handlers (input adapters) ---
	userHandler := handler.NewUserHandler(userSvc)
	bookHandler := handler.NewBookHandler(bookSvc)

	// --- HTTP server (input adapter) ---
	app := routes.NewApp(cfg.App, userHandler, bookHandler, func(ctx context.Context) error {
		if err := db.Ping(ctx); err != nil {
			return err
		}
		return redis.Ping(ctx)
	})

	// --- Graceful shutdown ---
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(fmt.Sprintf(":%s", cfg.App.ApiPort))
	}()

	select {
	case err := <-listenErr:
		if err != nil {
			log.Println(err)
		}
	case <-ctx.Done():
		fmt.Println("✔ [INFO] Shutting down server")
		timeout := time.Duration(cfg.App.ShutdownTimeout) * time.Second
		if err := app.ShutdownWithTimeout(timeout); err != nil {
			log.Println(err)
		}
	}
}
