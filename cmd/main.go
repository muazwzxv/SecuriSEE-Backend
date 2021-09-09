package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	fmt.Println("Application boot up")

	// load configuration file
	config.New()

	// Connect to database
	service.ConnectDatabase()

	// Load all app configs
	app := setup()

	app.Listen(":8080")

	fmt.Println("Hello what the fuck")
}

func setup() *fiber.App {
	app := fiber.New()

	// Unauthenticated routes go here

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.CFG.GetJWTSecret()),
	}))

	// Authenticated routes go here

	return app
}
