package main

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/controller"
	_ "Oracle-Hackathon-BE/docs/swagger"
	"Oracle-Hackathon-BE/service"
	"Oracle-Hackathon-BE/util"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
)

// @title CrimeNow Backend APi
// @version 1.0
// @description This is the first version of this API service.
// @termsOfService http://swagger.io/terms/

// @contact.name Muaz terkacak
// @contact.email

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	fmt.Println("Application boot up")
	app := fiber.New()

	// Connect to database
	// gorm := service.ConnectDatabase()
	if err := service.GetGormInstance().Migrate(); err != nil {
		panic(err)
	}

	// Setup middleware
	setupMiddleware(app)

	// Load Routers

	setupRouter(app)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", config.GetInstance().FetchServerPort())); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)

	// Notify channel if interrup or termination signal is sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // Block main thread until interrupt is received
	fmt.Println("Gracefully shutting down ....")
	_ = app.Shutdown()

	fmt.Println("Cleaning up task ........")

	service.GetGormInstance().GormShutDown()

	fmt.Println("Shitdown Complete !")

}

func setupMiddleware(app *fiber.App) {
	// Recover after program panic
	app.Use(recover.New())

	// Logger
	app.Use(logger.New())

	app.Use(cors.New())

	app.Server().MaxConnsPerIP = 1
	app.Static("/cdn/image/muazkacak", "./images")

	// Setup swagger
	app.Get("/swagger/*", swagger.Handler)

}

func setupRouter(app *fiber.App) {
	v1 := app.Group("/api")

	userRepository := controller.NewUserController()
	v1.Post("/user", userRepository.CreateUser)
	v1.Post("/user/admin", userRepository.CreateAdminOrCamera)
	v1.Get("/user/:id", JwtMiddleware(), userRepository.GetByID)
	v1.Get("/user", JwtMiddleware(), userRepository.GetAll)
	v1.Get("/user/:id/reports", JwtMiddleware(), userRepository.GetUserReports)

	// Auth
	authRepository := controller.NewAuthController()
	v1.Post("/login/user", authRepository.LoginUser)
	v1.Post("/login/admin", authRepository.LoginAdminAndCamera)
	v1.Get("/me", JwtMiddleware(), userRepository.Me)

	carEntryrepository := controller.NewCarEntryController()
	v1.Post("/car", JwtMiddleware(), carEntryrepository.CreateEntry)
	v1.Get("/car", JwtMiddleware(), carEntryrepository.GetAll)
	v1.Get("/car/:id", JwtMiddleware(), carEntryrepository.GetById)
	v1.Get("/car/:plate/plate", JwtMiddleware(), carEntryrepository.GetByPlate)

	newsRepository := controller.NewNewsRepository()
	v1.Post("/news", JwtMiddleware(), newsRepository.Create)
	v1.Get("/news/:id", JwtMiddleware(), newsRepository.GetById)
	v1.Get("/news", JwtMiddleware(), newsRepository.GetAll)

	imageRepository := controller.NewImageRepository()
	v1.Post("/image/upload/:reportId", JwtMiddleware(), imageRepository.Upload)
	v1.Get("/image/download/:imageId", JwtMiddleware(), imageRepository.Download)

	reportRepository := controller.NewReportRepository()
	v1.Post("/report", JwtMiddleware(), reportRepository.Create)
	v1.Get("/report", JwtMiddleware(), reportRepository.GetAll)
	v1.Get("/report/:id", JwtMiddleware(), reportRepository.GetById)
	v1.Get("/report/:id/image", JwtMiddleware(), reportRepository.GetImageFromReport)
}

// Jwt middleware
func JwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: func(c *fiber.Ctx) error {
			claims := util.GetClaims(c)

			c.Locals("userId", claims["ID"].(string))
			c.Locals("claims", claims)

			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Message": "Unauthorized",
				"Error":   e.Error(),
			})
		},
		SigningKey: []byte(config.GetInstance().GetJWTSecret()),
	})
}
