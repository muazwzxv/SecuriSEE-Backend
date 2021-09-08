package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Application boot up")

	// load configuration file
	config.New()

	// Connect to database
	database.Connect()

	fmt.Println("Hello what the fuck")
}

func setup() *fiber.App {
	app := fiber.New()

	return app
}
