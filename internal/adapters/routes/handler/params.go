package handler

import "github.com/gofiber/fiber/v3"

// idParam is the :id route param; the StructValidator rejects non-UUID values.
type idParam struct {
	ID string `uri:"id" validate:"required,uuid"`
}

// bindID binds and validates :id with Fiber v3's c.Bind().URI().
func bindID(c fiber.Ctx) (string, error) {
	var p idParam
	if err := c.Bind().URI(&p); err != nil {
		return "", err
	}
	return p.ID, nil
}