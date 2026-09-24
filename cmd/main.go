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
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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
	app := newApp(cfg, userHandler, bookHandler, func(ctx context.Context) error {
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

// newApp wires the Fiber v3 app: config, middleware and routes.
// ready is the readiness probe behind /readyz (DB + Redis in production).
func newApp(cfg *config.Config, userHandler *handler.UserHandler, bookHandler *handler.BookHandler, ready func(context.Context) error) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.App.AppName,
		// c.Bind().Body()/URI() run `validate:"..."` tags through this automatically
		StructValidator: routes.NewStructValidator(),
	})

	// requestid must be created before the logger so ${requestid} renders
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time} | ${requestid} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.Printf("[PANIC] request_id=%s path=%s method=%s error=%v\n%s",
				requestid.FromContext(c),
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

	if origins := splitCSV(cfg.App.AllowedCredentialOrigins); len(origins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins: origins,
			// Fiber v3 panics on AllowCredentials with "*", so only enable it for explicit origins
			AllowCredentials: !slices.Contains(origins, "*"),
		}))
	}

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     cfg.App.AppName,
			"version": cfg.App.AppVersion,
		})
	})
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
			defer cancel()
			return ready(ctx) == nil
		},
	}))

	// Rate limit only the API, not the health endpoints
	app.Use("/api", limiter.New(limiter.Config{
		Max:        cfg.App.RateLimitMax,
		Expiration: time.Duration(cfg.App.RateLimitWindowSeconds) * time.Second,
		LimitReached: func(c fiber.Ctx) error {
			return handler.ResError(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded")
		},
	}))

	routes.RegisterUsersRoutes(app, userHandler)
	routes.RegisterBookRoutes(app, bookHandler)

	return app
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
