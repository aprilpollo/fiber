package routes

import (
	"github.com/gofiber/fiber/v2"
)

func ResData(ctx *fiber.Ctx, status int, message string, errorText string, data any, optional ...int) error {
	if len(optional) > 0 {
		total := optional[0]
		page  := optional[1]
		limit := optional[2]

		rsp := fiber.Map{"code": status, "message": message, "error": errorText, "payload": data, "pagination": fiber.Map{"total": total, "page": page, "limit": limit}}
		return ctx.Status(status).JSON(rsp)
	} else {
		rsp := fiber.Map{"code": status, "message": message, "error": errorText, "payload": data}
		return ctx.Status(status).JSON(rsp)
	}
}
