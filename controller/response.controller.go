package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Send a response with status and message.
func response(ctx *fiber.Ctx, status int, msg fiber.Map) error {
		return ctx.Status(status).JSON(msg)
}

func Error(ctx *fiber.Ctx, msg string, e error, status int) error {
  return response(ctx, status, fiber.Map{
    "success": false,
    "message": msg,
    "data": e,
  })
}

func Success(ctx *fiber.Ctx, msg string, data interface{}, status int) error {
  return response(ctx, http.StatusOK, fiber.Map{
    "success": true,
    "message": msg,
    "data": data,
  })
}

