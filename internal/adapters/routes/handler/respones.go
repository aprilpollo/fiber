package handler

import (
	"aprilpollo/internal/pkg/query"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func ResOk(ctx *fiber.Ctx, status int, payload any, total *int64, opts *query.QueryOptions) error {
	if total != nil && opts != nil {
		page := uint(1)
		if opts.Limit > 0 {
			page = uint((opts.Offset / opts.Limit) + 1)
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
				"limit": opts.Limit,
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

func ResError(ctx *fiber.Ctx, status int, message string, errorText string) error {
	rsp := fiber.Map{
		"code":    status,
		"message": message,
		"error":   errorText,
		"payload": nil,
	}
	return ctx.Status(status).JSON(rsp)
}
