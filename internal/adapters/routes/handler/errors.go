package handler

import (
	"errors"
	"log"

	"aprilpollo/internal/core/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// respondError answers with the status a core error stands for. Anything the
// core did not name is internal: it is logged and answered generically, so a
// driver or SQL error never reaches the client.
func respondError(c fiber.Ctx, err error) error {
	var verr *domain.ValidationError
	switch {
	case errors.As(err, &verr):
		return ResError(c, fiber.StatusBadRequest, "BAD_REQUEST", verr.Error())
	case errors.Is(err, domain.ErrNotFound):
		return ResError(c, fiber.StatusNotFound, "NOT_FOUND", domain.ErrNotFound.Error())
	case errors.Is(err, domain.ErrConflict):
		return ResError(c, fiber.StatusConflict, "CONFLICT", domain.ErrConflict.Error())
	}
	log.Printf("[ERROR] request_id=%s %s %s: %v", requestid.FromContext(c), c.Method(), c.Path(), err)
	return ResError(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
}

// badRequest answers a request the adapter could not read (malformed JSON,
// a failed `validate` tag, a bad query param). Those messages describe the
// request, not the server, so they are safe to return.
func badRequest(c fiber.Ctx, err error) error {
	return ResError(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error())
}
