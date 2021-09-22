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


func (r *UserRepository) GetAll(ctx *fiber.Ctx) error {
	claim := util.GetClaims(ctx)

	var user model.User
	user.GetUserById(r.gorm, claim["ID"].(string))

	// Check permissions
	if !user.IsAdmin() {
    return Error(ctx, "Not Allowed", nil, http.StatusForbidden)
	}

	if users, err := user.GetAll(r.gorm, ctx); err != nil {
    return Error(ctx, err.Error(), nil, http.StatusConflict)
	} else {
    return Success(ctx, "Successfully get all users", users, http.StatusOK)
	}
}

func (r *UserRepository) GetByID(ctx *fiber.Ctx) error {
	var user model.User
	if err := user.GetUserById(r.gorm, ctx.Params("id")); err != nil {
    return Error(ctx, "User not found", err, http.StatusNotFound)
	}

  return Success(ctx, "User found", user, http.StatusOK)
}

func (r *UserRepository) Me(ctx *fiber.Ctx) error {
	var user model.User

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	err := user.GetUserById(r.gorm, claims["ID"].(string))
	if err != nil {
    return Error(ctx, "User not found", err, http.StatusNotFound)
	}

  return Success(ctx, "User found", user, http.StatusOK)
}

func (r *UserRepository) CreateAdminOrCamera(ctx *fiber.Ctx) error {
	var user model.User
	// Parse json
	if err := ctx.BodyParser(&user); err != nil {
    return Error(ctx, "Cannot parse JSON", err, http.StatusBadRequest)
	}

	// Set Role
	if user.Role == "" {
		// If role not provided, default to user
		user.RolesToString([]string{"camera"})
	}

	if err := ValidateAdminCamera(user); err != nil {
    return Error(ctx, err.Error(), err, http.StatusBadRequest)
	}

	// check email exists
	if cond := user.IsEmailExist(r.gorm); cond == true {
    return Error(ctx, "Email already exists", nil, http.StatusBadRequest)
	}

	// Hash password
	user.HashPassword(user.Password)

	// Create user
	if err := user.Create(r.gorm); err != nil {
    return Error(ctx, err.Error(), nil, http.StatusConflict)
	}

  return Success(ctx, "User created", user, http.StatusOK)
}

func (r *UserRepository) CreateUser(ctx *fiber.Ctx) error {
	var user model.User

	// Parse json
	if err := ctx.BodyParser(&user); err != nil {
    return Error(ctx, "Cannot parse JSON", err, http.StatusBadRequest)
	}

	// Set Role
	if user.Role == "" {
		// If role not provided, default to user
		user.RolesToString([]string{"user"})
	}

	// validate payload
	if err := user.ValidateCreate(); err != nil {
    return Error(ctx, err.Error(), nil, http.StatusBadRequest)
	}

	// Check ic exist
	if cond := user.IsICExist(r.gorm); cond == true {
    return Error(ctx, "IC already exists", nil, http.StatusBadRequest)
	}

	// check email exists
	if cond := user.IsEmailExist(r.gorm); cond == true {
    return Error(ctx, "Email already exists", nil, http.StatusBadRequest)
	}

	// Hash password
	user.HashPassword(user.Password)

	// Create user
	err := user.Create(r.gorm)
	if err != nil {
    return Error(ctx, err.Error(), nil, http.StatusConflict)
	}

  return Success(ctx, "User created", user, http.StatusCreated)
}

func (r *UserRepository) GetUserReports(ctx *fiber.Ctx) error {
	var user model.User

	if err := user.GetUserById(r.gorm, ctx.Params("id")); err != nil {
    return Error(ctx, err.Error(), err, http.StatusNotFound)
	}

	if reports, err := user.GetAssociateReports(r.gorm); err != nil {
    return Error(ctx, err.Error(), err, http.StatusNotFound)
	} else {
    return Success(ctx, "Reports found", reports, http.StatusOK)
	}
}
