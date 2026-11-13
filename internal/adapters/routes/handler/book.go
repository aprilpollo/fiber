package handler

import (
	"aprilpollo/internal/core/ports/input"

	"github.com/gofiber/fiber/v3"
)

// BookHandler is the HTTP input adapter for books: it reads the request,
// calls the input port, and writes the response. It holds no business rules.
type BookHandler struct {
	svc input.BookService
}

func NewBookHandler(svc input.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// GET /api/v1/books
func (h *BookHandler) Gets(c fiber.Ctx) error {
	q, err := parseListQuery(c.Queries())
	if err != nil {
		return badRequest(c, err)
	}

	books, total, err := h.svc.List(c.Context(), q)
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, mapSlice(books, toBookResponse), &total, &q)
}

// GET /api/v1/books/:id
func (h *BookHandler) GetByID(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	book, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, toBookResponse(book), nil, nil)
}

// POST /api/v1/books
func (h *BookHandler) Create(c fiber.Ctx) error {
	var req createBookRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, err)
	}

	book, err := h.svc.Create(c.Context(), req.toInput())
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusCreated, toBookResponse(book), nil, nil)
}

// PUT /api/v1/books/:id
func (h *BookHandler) Update(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	var req updateBookRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, err)
	}

	book, err := h.svc.Update(c.Context(), id, req.toInput())
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, toBookResponse(book), nil, nil)
}

// DELETE /api/v1/books/:id
func (h *BookHandler) Delete(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
