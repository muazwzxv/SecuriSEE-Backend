package controller

import (
	"Oracle-Hackathon-BE/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsRepository struct {
	gorm *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{gorm: db}
}

func (r *NewsRepository) Create(ctx *fiber.Ctx) error {
	// validate role
	userId := ctx.Locals("userId").(string)
	var user model.User
	user.GetUserById(r.gorm, userId)

	// Check permissions
	if !user.IsRoleAdmin() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	// parse json
	var news model.News
	if err := ctx.BodyParser(&news); err != nil {
		return BadRequest(ctx, "Cannot parse JSON", err)
	}

	// Validate json
	if err := news.Validate(); err != nil {
		return BadRequest(ctx, err.Error(), err)
	}

	// Create
	if err := news.Create(r.gorm); err != nil {
		return Conflict(ctx, err.Error(), err)
	}

	return Created(ctx, "News entry created", news)
}

func (r *NewsRepository) GetAll(ctx *fiber.Ctx) error {
	var n model.News

	if news, err := n.GetAll(r.gorm, ctx); err != nil {
		return Conflict(ctx, err.Error(), err)
	} else {
		return Ok(ctx, "Successfully get all news", news)
	}
}

func (r *NewsRepository) GetById(ctx *fiber.Ctx) error {
	var news model.News
	if err := news.GetById(r.gorm, ctx.Params("id")); err != nil {
		return NotFound(ctx, err.Error(), err)
	}

	return Ok(ctx, "Found news", news)
}
