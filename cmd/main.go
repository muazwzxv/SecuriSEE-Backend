package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/controller"
	"Oracle-Hackathon-BE/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Application boot up")

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovering")
			log.Printf("\n Error was %v", r)
		}
	}()

	// load configuration file
	config.New()

	// Connect to database
	gorm := service.ConnectDatabase()

	// Load all app configs
	app := setup(gorm.Orm)

	app.Listen(":8080")

	fmt.Println("Hello what the fuck")
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

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.CFG.GetJWTSecret()),
	}))

	// Authenticated routes go here

	return app
}
