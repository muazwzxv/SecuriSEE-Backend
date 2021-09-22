package controller

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/service"

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

func (r *AuthRepository) LoginAdminAndCamera(ctx *fiber.Ctx) error {
	var login model.LoginAdminAndCamera
	if err := ctx.BodyParser(&login); err != nil {
    return BadRequest(ctx, err.Error(), err)
	}

	// Get User by Email
	var user model.User
	if err := user.GetUserByEmail(r.gorm, login.Email); err != nil {
    return BadRequest(ctx, err.Error(), err)
	}

	// Check password
	if isMatch := user.CheckHash(login.Password); !isMatch {
    return BadRequest(ctx, "Password does not match", nil)
	}

	// Generate jwt token
	jwt := service.JwtWrapper{
		SecretKey:    config.GetInstance().GetJWTSecret(),
		Issuer:       "CrimeNow Backend",
		ExpiredHours: 24,
	}

	payload := &config.UserJwt{
		ID:   user.ID.String(),
		IC:   user.Ic,
		Role: user.RolesToArray(),
	}

	if token, err := jwt.GenerateToken(payload); err != nil {
    return BadRequest(ctx, err.Error(), err)
	} else {
    return Ok(ctx, "Successfully logged in", token)
	}
}

func (r *AuthRepository) LoginUser(ctx *fiber.Ctx) error {

	var login model.LoginUser
	if err := ctx.BodyParser(&login); err != nil {
    return BadRequest(ctx, err.Error(), nil)
	}

	// Get User by IC
	var user model.User
	if err := user.GetUserByIc(r.gorm, login.IC); err != nil {
    return BadRequest(ctx, err.Error(), nil)
	}

	// Check password
	if isMatch := user.CheckHash(login.Password); !isMatch {
    return BadRequest(ctx, "Password does not match", nil)
	}

	// Generate jwt token
	jwt := service.JwtWrapper{
		SecretKey:    config.GetInstance().GetJWTSecret(),
		Issuer:       "CrimeNow Backend",
		ExpiredHours: 24,
	}

	payload := &config.UserJwt{
		ID:   user.ID.String(),
		IC:   user.Ic,
		Role: user.RolesToArray(),
	}

	if token, err := jwt.GenerateToken(payload); err != nil {
    return BadRequest(ctx, "Failed to generated token", nil)
	} else {
    return Ok(ctx, "Successfully logged in", token)
	}
}
