package handler

import (
	"aprilpollo/internal/core/domain"
	"reflect"

	"github.com/gofiber/fiber/v3"
)

func ResOk(ctx fiber.Ctx, status int, payload any, total *int64, q *domain.ListQuery) error {
	if total != nil && q != nil {
		page := uint(1)
		if q.Limit > 0 {
			page = uint((q.Offset / q.Limit) + 1)
		}

		count := 0
		if v := reflect.ValueOf(payload); v.Kind() == reflect.Slice {
			count = v.Len()
		}

		rsp := fiber.Map{
			"code":    status,
			"message": "OK",
			"error":   nil,
			"payload": payload,
			"pagination": fiber.Map{
				"total": total,
				"count": count,
				"page":  page,
				"limit": q.Limit,
			},
		}
		return ctx.Status(status).JSON(rsp)
	}

	rsp := fiber.Map{
		"code":    status,
		"message": "OK",
		"error":   nil,
		"payload": payload,
	}
	return ctx.Status(status).JSON(rsp)
}

func ResError(ctx fiber.Ctx, status int, message string, errorText string) error {
	rsp := fiber.Map{
		"code":    status,
		"message": message,
		"error":   errorText,
		"payload": nil,
	}
	return ctx.Status(status).JSON(rsp)
}
