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
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
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

	// --- Fiber app ---
	app := fiber.New(fiber.Config{
		AppName: cfg.App.AppName,
	})

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.Printf("[PANIC] path=%s method=%s error=%v\n%s",
				c.Path(),
				c.Method(),
				e,
				debug.Stack(),
			)
		},
		// v3: the response is written here; whatever this returns goes to the ErrorHandler
		PanicHandler: func(c fiber.Ctx, _ any) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     cfg.App.AppName,
			"version": cfg.App.AppVersion,
		})
	})

	routes.RegisterUsersRoutes(app, userHandler)
	routes.RegisterBookRoutes(app, bookHandler)

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
