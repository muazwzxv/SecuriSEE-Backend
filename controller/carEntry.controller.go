package controller

import (
	"Oracle-Hackathon-BE/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CarEntryRepository struct {
	gorm *gorm.DB
}

func NewCarEntryController(db *gorm.DB) *CarEntryRepository {
	return &CarEntryRepository{gorm: db}
}

func (r *CarEntryRepository) GetByPlate(ctx *fiber.Ctx) error {
	var car model.CarEntry

	if cars, err := car.GetEntryByPlate(r.gorm, ctx.Params("plateId")); err != nil {
		return NotFound(ctx, err.Error(), nil)
	} else {
		return Ok(ctx, "car found", cars)
	}
}

func (r *CarEntryRepository) GetById(ctx *fiber.Ctx) error {
	var car model.CarEntry
	if err := car.GetEntryById(r.gorm, ctx.Params("id")); err != nil {
		return NotFound(ctx, err.Error(), nil)
	}

	return Ok(ctx, "Found car", car)
}

func (r *CarEntryRepository) GetAll(ctx *fiber.Ctx) error {
	// userId := ctx.Locals("userId").(string)
	// var user model.User
	// user.GetUserById(r.gorm, claim["ID"].(string))

	// // Check permissions
	// if !user.IsRoleAdmin() {
	// 	return Forbidden(ctx, "Not allowed", nil)
	// }

	var car model.CarEntry
	if cars, err := car.GetAll(r.gorm, ctx); err != nil {
		return Conflict(ctx, err.Error(), nil)
	} else {
		return Ok(ctx, "Successfully get all cars", cars)
	}
}

func (r *CarEntryRepository) CreateEntry(ctx *fiber.Ctx) error {
	// validate role
	// userId := ctx.Locals("userId").(string)
	// var user model.User
	// user.GetUserById(r.gorm, userId)

	// // Check permissions
	// if !user.IsRoleCamera() {
	// 	return Forbidden(ctx, "Not Allowed", nil)
	// }

	var entry model.CarEntry

	// parse json
	if err := ctx.BodyParser(&entry); err != nil {
		return BadRequest(ctx, "Cannot parse JSON", err)
	}

	// Validate input
	if err := entry.Validate(); err != nil {
		return BadRequest(ctx, err.Error(), err)
	}

	// // validate status provided
	// if err := entry.CheckStatus(); err != nil {
	// 	return BadRequest(ctx, err.Error(), err)
	// }

	// Create
	if err := entry.Create(r.gorm); err != nil {
		return Conflict(ctx, err.Error(), err)
	}

	return Created(ctx, "Car entry created", entry)
}
