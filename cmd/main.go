package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/controller"
	"Oracle-Hackathon-BE/service"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Application boot up")

	// load configuration file
	config.New()

	// Connect to database
	gorm := service.ConnectDatabase()

	// Load all app configs
	app := setup(gorm.Orm)

	// Recover after program panic
	app.Use(recover.New())

	app.Listen(":8080")
}

func setup(gorm *gorm.DB) *fiber.App {
	app := fiber.New()
	v1 := app.Group("/api")

	// Unauthenticated routes go here
	userRepository := controller.NewUserController(gorm)
	v1.Post("/user", userRepository.CreateUser)
	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
			"Success": true,
			"Message": "Welcome to our endpoint bitch",
		})
	})

	// Authenticated routes go here

	return app
}

func JwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Message": "Unauthorized",
				"Error":   e.Error(),
			})
		},
		SigningKey: []byte(config.CFG.GetJWTSecret()),
	})
}
