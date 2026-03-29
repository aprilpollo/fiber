package routes

import (
	"aprilpollo/internal/adapters/routes/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterUsersRoutes(app *fiber.App, h *handler.UserHandler) {
	api := app.Group("/api/v1")

	users := api.Group("/users")
	users.Get("/", h.Gets)
	users.Get("/:id", h.GetByID)
	users.Post("/", h.Create)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}

func RegisterBookRoutes(app *fiber.App, h *handler.BookHandler) {
	api := app.Group("/api/v1")

	books := api.Group("/books")
	books.Get("/", h.Gets)
	books.Get("/:id", h.GetByID)
	books.Post("/", h.Create)
	books.Put("/:id", h.Update)
	books.Delete("/:id", h.Delete)
}
