package handler

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/pkg/query"
	"aprilpollo/internal/core/ports/input"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	svc input.BookService
}

func NewBookHandler(svc input.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// GET /api/v1/books
func (h *BookHandler) Gets(c *fiber.Ctx) error {
	opts, err := query.Parse("books", c.Queries())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	books, total, err := h.svc.List(c.Context(), opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  books,
		"total": total,
	})
}

// GET /api/v1/books/:id
func (h *BookHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	book, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if book == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}

	return c.JSON(fiber.Map{"data": book})
}

// POST /api/v1/books
func (h *BookHandler) Create(c *fiber.Ctx) error {
	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.svc.Create(c.Context(), &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": book})
}

// PUT /api/v1/books/:id
func (h *BookHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	existing, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}

	if err := c.BodyParser(existing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.svc.Update(c.Context(), existing); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": existing})
}

// DELETE /api/v1/books/:id
func (h *BookHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
