package main

import (
	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/routes"
	"aprilpollo/internal/adapters/routes/handler"
	"aprilpollo/internal/adapters/storage/cache"
	"aprilpollo/internal/adapters/storage/orm"
	"aprilpollo/internal/adapters/storage/repository"
	"aprilpollo/internal/core/services"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("[PANIC] path=%s method=%s error=%v\n%s",
				c.Path(),
				c.Method(),
				e,
				debug.Stack(),
			)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     cfg.App.AppName,
			"version": cfg.App.AppVersion,
		})
	})

	routes.RegisterUsersRoutes(app, userHandler)
	routes.RegisterBookRoutes(app, bookHandler)

	if err := app.Listen(fmt.Sprintf(":%s", cfg.App.ApiPort)); err != nil {
		log.Println(err)
	}
}
