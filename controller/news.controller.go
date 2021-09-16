package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsRepository struct {
	gorm *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{gorm: db}
}

func (newsRepository *NewsRepository) Create(ctx *fiber.Ctx) error {
	// validate role
	claim := util.GetClaims(ctx)
	var user model.User
	user.GetUserById(newsRepository.gorm, claim["ID"].(string))

	// Check permissions
	if isAdmin := user.IsRoleExist("admin"); !isAdmin {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	// parse json
	var news model.News
	if err := ctx.BodyParser(&news); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Cannot parse JSON",
		})
	}

	// Validate json
	if err := news.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	// Create
	if err := news.Create(newsRepository.gorm); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"Success": true,
		"Message": "Entry created",
		"news":    news,
	})
}

func (newsRepository *NewsRepository) GetAll(ctx *fiber.Ctx) error {
	var n model.News

	if news, err := n.GetAll(newsRepository.gorm, ctx); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": true,
			"Message": err.Error(),
		})
	} else {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"Success": true,
			"News":    news,
		})
	}
}

func (newsRepository *NewsRepository) GetById(ctx *fiber.Ctx) error {
	var news model.News
	if err := news.GetById(newsRepository.gorm, ctx.Params("id")); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"Success": false,
			"Error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Success": true,
		"News":    news,
	})
}
