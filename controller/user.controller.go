package controller

import (
	"Oracle-Hackathon-BE/model"
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
	if isAdmin := user.IsRoleExist("admin"); !isAdmin {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	if users, err := user.GetAll(userRepository.gorm, ctx); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": true,
			"Message": err.Error(),
		})
	} else {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"Success": true,
			"User":    users,
		})
	}
}

func (userRepository *UserRepository) GetByID(ctx *fiber.Ctx) error {
	var user model.User
	if err := user.GetUserById(userRepository.gorm, ctx.Params("id")); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
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

func (userRepository *UserRepository) CreateAdminOrCamera(ctx *fiber.Ctx) error {
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
		user.RolesToString([]string{"camera"})
	}

	if err := ValidateAdminCamera(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
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
	if err := user.Create(userRepository.gorm); err != nil {
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
	if err := user.ValidateCreate(); err != nil {
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
