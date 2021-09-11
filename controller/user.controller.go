package controller

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	gorm *gorm.DB
}

func NewUserController(db *gorm.DB) *UserRepository {
	return &UserRepository{gorm: db}
}

func (userRepository *UserRepository) Login(ctx *fiber.Ctx) error {

	var login model.Login
	if err := ctx.BodyParser(&login); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Cannot parse JSON",
			"Error":   err,
		})
	}

	// Get User by IC
	var user model.User
	if err := user.GetUserByIc(userRepository.gorm, login.IC); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "IC not found",
			"Error":   err,
		})
	}

	// Check password
	isMatch := user.CheckHash(login.Password)
	if !isMatch {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Password does not match",
		})
	}

	// Generate jwt token
	jwt := service.JwtWrapper{
		SecretKey:    config.CFG.GetJWTSecret(),
		Issuer:       "CrimeNow Backend",
		ExpiredHours: 24,
	}

	payload := &config.UserJwt{
		ID:   user.ID.String(),
		IC:   user.Ic,
		Role: user.RolesToArray(),
	}

	if token, err := jwt.GenerateToken(payload); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Failed to generate token",
		})
	} else {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"Success": true,
			"Token":   token,
		})
	}
}

func (userRepository *UserRepository) CreateUser(ctx *fiber.Ctx) error {
	var user model.User

	// Parse json
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Cannot parse JSON",
		})
	}

	// Check ic exist
	if cond := user.IsICExist(userRepository.gorm); cond == true {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "IC already exists",
		})
	}

	// check email exists
	if cond := user.IsEmailExist(userRepository.gorm); cond == true {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Email already exists",
		})
	}

	// Hash password
	user.HashPassword(user.Password)

	// Set Role
	user.RolesToString([]string{"user"})

	// Create user
	err := user.Create(userRepository.gorm)
	if err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": "Something wrong happened",
			"Error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"Success": true,
		"Message": "User created",
		"User":    user,
	})

}
