package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CarEntryRepository struct {
	gorm *gorm.DB
}

func NewCarEntryController(db *gorm.DB) *CarEntryRepository {
	return &CarEntryRepository{gorm: db}
}

func (carEntryRepository *CarEntryRepository) CreateEntry(ctx *fiber.Ctx) error {
	// validate role
	claim := util.GetClaims(ctx)
	var user model.User
	user.GetUserById(carEntryRepository.gorm, claim["ID"].(string))

	// Check permissions
	if isCamera := user.IsRoleExist("camera"); !isCamera {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	var entry model.CarEntry

	// parse json
	if err := ctx.BodyParser(&entry); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Cannot parse JSON",
		})
	}

	// Validate input
	if err := entry.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	// validate status provided
	if err := entry.CheckStatus(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	// Create
	err := entry.Create(carEntryRepository.gorm)
	if err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"Success": true,
		"Message": "Entry created",
		"entry":   entry,
	})
}
