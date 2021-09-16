package controller

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthRepository struct {
	gorm *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthRepository {
	return &AuthRepository{gorm: db}
}

func ValidateAdminCamera(u model.User) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}

func (authRepository *AuthRepository) LoginAdminAndCamera(ctx *fiber.Ctx) error {
	var login model.LoginAdminAndCamera
	if err := ctx.BodyParser(&login); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
		})
	}

	// Get User by Email
	var user model.User
	if err := user.GetUserByEmail(authRepository.gorm, login.Email); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
		})
	}

	// Check password
	if isMatch := user.CheckHash(login.Password); !isMatch {
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
			"Message": err.Error(),
		})
	} else {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"Success": true,
			"Token":   token,
		})
	}
}

func (authRepository *AuthRepository) LoginUser(ctx *fiber.Ctx) error {

	var login model.LoginUser
	if err := ctx.BodyParser(&login); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
		})
	}

	// Get User by IC
	var user model.User
	if err := user.GetUserByIc(authRepository.gorm, login.IC); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
		})
	}

	// Check password
	if isMatch := user.CheckHash(login.Password); !isMatch {
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
