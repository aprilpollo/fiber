package routes

import (
	"context"
	"log"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/routes/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// NewApp wires the Fiber v3 app: config, middleware and routes.
// ready is the readiness probe behind /readyz (DB + Redis in production).
func NewApp(cfg *config.App, userHandler *handler.UserHandler, bookHandler *handler.BookHandler, ready func(context.Context) error) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
		// c.Bind().Body()/URI() run `validate:"..."` tags through this automatically
		StructValidator: NewStructValidator(),
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

	if origins := splitCSV(cfg.AllowedCredentialOrigins); len(origins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins: origins,
			// Fiber v3 panics on AllowCredentials with "*", so only enable it for explicit origins
			AllowCredentials: !slices.Contains(origins, "*"),
		}))
	}

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     cfg.AppName,
			"version": cfg.AppVersion,
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
		Max:        cfg.RateLimitMax,
		Expiration: time.Duration(cfg.RateLimitWindowSeconds) * time.Second,
		LimitReached: func(c fiber.Ctx) error {
			return handler.ResError(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded")
		},
	}))

	RegisterUsersRoutes(app, userHandler)
	RegisterBookRoutes(app, bookHandler)

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
