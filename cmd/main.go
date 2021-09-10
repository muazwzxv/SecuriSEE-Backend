package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/controller"
	"Oracle-Hackathon-BE/service"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Application boot up")
	app := fiber.New()

	// load configuration file
	config.New()

	// Connect to database
	gorm := service.ConnectDatabase()

	// Setup middleware
	setupMiddleware(app)

	// Load Routers
	setupRouter(gorm.Orm, app)

	app.Listen(":8080")
}

func setupMiddleware(app *fiber.App) {
	// Recover after program panic
	app.Use(recover.New())
	// Logger
	app.Use(logger.New())
}

func setupRouter(gorm *gorm.DB, app *fiber.App) {
	v1 := app.Group("/api")

	userRepository := controller.NewUserController(gorm)
	v1.Post("/user", userRepository.CreateUser)
	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusAccepted).JSON(fiber.Map{
			"Success": true,
			"Message": "Welcome to our endpoint bitch",
		})
	})

	v1.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Panic testing")
	})
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
