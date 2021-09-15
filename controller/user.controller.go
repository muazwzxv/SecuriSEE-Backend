package controller

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"
	"Oracle-Hackathon-BE/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type UserRepository struct {
	gorm *gorm.DB
}

func NewUserController(db *gorm.DB) *UserRepository {
	return &UserRepository{gorm: db}
}

func (userRepository *UserRepository) GetAll(ctx *fiber.Ctx) error {
	claim := util.GetClaims(ctx)

	var user model.User
	user.GetUserById(userRepository.gorm, claim["ID"].(string))

	// Check permissions
	isAdmin := user.IsRoleExist("admin")
	if !isAdmin {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	users, err := user.GetAll(userRepository.gorm, ctx)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Success": true,
		"Message": "User found",
		"User":    users,
	})
}

func (userRepository *UserRepository) GetByID(ctx *fiber.Ctx) error {
	var user model.User

	err := user.GetUserById(userRepository.gorm, ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"Success": false,
			"Message": "User not found",
			"Error":   err,
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Success": true,
		"Message": "User found",
		"User":    user,
	})
}

func (userRepository *UserRepository) Me(ctx *fiber.Ctx) error {
	var user model.User

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	err := user.GetUserById(userRepository.gorm, claims["ID"].(string))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"Success": false,
			"Message": "User not found",
			"Error":   err,
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Success": true,
		"Message": "User found",
		"User":    user,
	})

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

	// Set Role
	if user.Role == "" {
		// If role not provided, default to user
		user.RolesToString([]string{"user"})
	}

	// validate payload
	if err := user.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
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

	// Create user
	err := user.Create(userRepository.gorm)
	if err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"Success": true,
		"Message": "User created",
		"User":    user,
	})

}
