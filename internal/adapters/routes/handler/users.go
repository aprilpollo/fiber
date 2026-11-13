package handler

import (
	"aprilpollo/internal/core/ports/input"

	"github.com/gofiber/fiber/v3"
)

// UserHandler is the HTTP input adapter for users: it reads the request,
// calls the input port, and writes the response. It holds no business rules.
type UserHandler struct {
	svc input.UserService
}

func NewUserHandler(svc input.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GET /api/v1/users
func (h *UserHandler) Gets(c fiber.Ctx) error {
	q, err := parseListQuery(c.Queries())
	if err != nil {
		return badRequest(c, err)
	}

	users, total, err := h.svc.List(c.Context(), q)
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, mapSlice(users, toUserResponse), &total, &q)
}

// GET /api/v1/users/:id
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	user, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, toUserResponse(user), nil, nil)
}

// POST /api/v1/users
func (h *UserHandler) Create(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, err)
	}

	user, err := h.svc.Create(c.Context(), req.toInput())
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusCreated, toUserResponse(user), nil, nil)
}

// PUT /api/v1/users/:id
func (h *UserHandler) Update(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	var req updateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, err)
	}

	user, err := h.svc.Update(c.Context(), id, req.toInput())
	if err != nil {
		return respondError(c, err)
	}

	return ResOk(c, fiber.StatusOK, toUserResponse(user), nil, nil)
}

// DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c fiber.Ctx) error {
	id, err := bindID(c)
	if err != nil {
		return badRequest(c, err)
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return respondError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
