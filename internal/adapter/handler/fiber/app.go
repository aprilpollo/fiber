package httpfiber

import (
	"aprilpollo/internal/adapter/handler/fiber/middleware"
	"aprilpollo/internal/adapter/handler/fiber/routes"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	app           *fiber.App
	middlewareApp *middleware.MiddlewareApp
}

func NewApp(app *fiber.App, middlewareApp *middleware.MiddlewareApp) *App {
	return &App{
		app:           app,
		middlewareApp: middlewareApp,
	}
}

func (r *App) MainRoutes() {
	r.app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version": "1.0.0",
		})
	})
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})
}

func (r *App) AuthRoutes(authHandler *routes.AuthHandler) {

}

func (r *App) Serve(port string) error {
	return r.app.Listen(port)
}
