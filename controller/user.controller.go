package controller

import (
	"Oracle-Hackathon-BE/model"
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

func (r *UserRepository) GetAll(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	var user model.User
	user.GetUserById(r.gorm, userId)

	// Check permissions
	if !user.IsRoleAdmin() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	if users, err := user.GetAll(r.gorm, ctx); err != nil {
		return Error(ctx, err.Error(), nil, http.StatusConflict)
	} else {
		return Ok(ctx, "Successfully get all users", users)
	}
}

func (r *UserRepository) GetByID(ctx *fiber.Ctx) error {
	var user model.User
	if err := user.GetUserById(r.gorm, ctx.Params("id")); err != nil {
		return NotFound(ctx, "User not found", err)
	}

	return Ok(ctx, "User found", user)
}

func (r *UserRepository) Me(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	var user model.User

	if err := user.GetUserById(r.gorm, userId); err != nil {
		return NotFound(ctx, "User not found", err)
	}

	return Ok(ctx, "User found", user)
}

func (r *UserRepository) CreateAdminOrCamera(ctx *fiber.Ctx) error {
	var user model.User
	// Parse json
	if err := ctx.BodyParser(&user); err != nil {
		return BadRequest(ctx, "Cannot parse JSON", err)
	}

	// Set Role
	if user.Role == "" {
		// If role not provided, default to user
		user.RolesToString([]string{"camera"})
	}

	if err := ValidateAdminCamera(user); err != nil {
		return BadRequest(ctx, err.Error(), err)
	}

	// check email exists
	if cond := user.IsEmailExist(r.gorm); cond == true {
		return BadRequest(ctx, "Email already exists", nil)
	}

	// Hash password
	user.HashPassword(user.Password)

	// Create user
	if err := user.Create(r.gorm); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	return Created(ctx, "User created", user)
}

func (r *UserRepository) CreateUser(ctx *fiber.Ctx) error {
	var user model.User

	// Parse json
	if err := ctx.BodyParser(&user); err != nil {
		return BadRequest(ctx, "Cannot parse JSON", err)
	}

	// Set Role
	if user.Role == "" {
		// If role not provided, default to user
		user.RolesToString([]string{"user"})
	}

	// validate payload
	if err := user.ValidateCreate(); err != nil {
		return BadRequest(ctx, err.Error(), nil)
	}

	// Check ic exist
	if cond := user.IsICExist(r.gorm); cond == true {
		return BadRequest(ctx, "IC already exists", nil)
	}

	// check email exists
	if cond := user.IsEmailExist(r.gorm); cond == true {
		return BadRequest(ctx, "Email already exists", nil)
	}

	// Hash password
	user.HashPassword(user.Password)

	// Create user
	err := user.Create(r.gorm)
	if err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	return Created(ctx, "User created", user)
}

func (r *UserRepository) GetUserReports(ctx *fiber.Ctx) error {
	var user model.User
	if err := user.GetUserById(r.gorm, ctx.Params("id")); err != nil {
		return NotFound(ctx, err.Error(), err)
	}

	if reports, err := user.GetAssociateReports(r.gorm); err != nil {
		return NotFound(ctx, err.Error(), err)
	} else {
		return Ok(ctx, "Reports found", reports)
	}
}
