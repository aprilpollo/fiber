package middleware

import (
	config "aprilpollo/internal/adapter/config"
	"aprilpollo/internal/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type MiddlewareApp struct {
	app       *fiber.App
	appConfig *config.App
	jwtConfig *config.JWT
}

func NewMiddlewareHandler(app *fiber.App, appConfig *config.App, jwtConfig *config.JWT) *MiddlewareApp {
	return &MiddlewareApp{
		app:       app,
		appConfig: appConfig,
		jwtConfig: jwtConfig,
	}
}

func (m *MiddlewareApp) SetupGlobalMiddleware() {
	m.app.Use(cors.New(cors.Config{
		AllowOrigins:     m.appConfig.AllowedCredentialOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))
	m.app.Use(recover.New())
	m.app.Use(requestid.New())

	if m.appConfig.Development {
		m.app.Use(logger.New(logger.Config{
			Format: "[${time}] [${ip}] [${method} ${status}] - ${path}\n",
		}))
	} else {
		multi, err := util.RotateLogs()
		if err != nil {
			log.Fatalf("failed to create log file: %v", err)
		}

		m.app.Use(logger.New(logger.Config{
			Format: "[${time}] [${ip}] [${method} ${status}] - ${path}\n",
			Output: multi,
		}))
	}

}

func (m *MiddlewareApp) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func (m *MiddlewareApp) RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
