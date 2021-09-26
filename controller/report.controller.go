package controller

import (
	"Oracle-Hackathon-BE/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReportRepository struct {
	gorm *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{gorm: db}
}

func (r *ReportRepository) Create(ctx *fiber.Ctx) error {
	// validate role
	userId := ctx.Locals("userId").(string)
	var user model.User
	user.GetUserById(r.gorm, userId)

	// Check permissions
	if !user.IsRoleUser() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	// parse json
	var report model.Report
	if err := ctx.BodyParser(&report); err != nil {
		return BadRequest(ctx, "Cannot parse JSON", err)
	}

	// Assign reference id
	report.UserID = user.ID.String()

	if err := report.Validate(); err != nil {
		return BadRequest(ctx, err.Error(), err)
	}

	if err := report.Create(r.gorm); err != nil {
		return Conflict(ctx, err.Error(), err)
	}

	return Created(ctx, "Report created", report)
}

func (r *ReportRepository) GetAll(ctx *fiber.Ctx) error {

	// validate role
	userId := ctx.Locals("userId").(string)
	var user model.User
	user.GetUserById(r.gorm, userId)

	// Check permissions
	if !user.IsRoleAdmin() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	var report model.Report

	if reports, err := report.GetAll(r.gorm, ctx); err != nil {
		return Forbidden(ctx, err.Error(), err)
	} else {
		return Ok(ctx, "Sucessfully get all reports", reports)
	}
}

func (r *ReportRepository) GetById(ctx *fiber.Ctx) error {
	// validate role
	userId := ctx.Locals("userId").(string)
	var user model.User
	user.GetUserById(r.gorm, userId)

	// Check permissions
	if !user.IsRoleAdmin() {
		return Forbidden(ctx, "Not allowed", nil)
	}

	var report model.Report
	if err := report.GetById(r.gorm, ctx.Params("id")); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	return Ok(ctx, "Report found", report)
}

func (r *ReportRepository) GetImageFromReport(ctx *fiber.Ctx) error {
	// validate role
	// claim := util.GetClaims(ctx)
	// var user model.User
	// user.GetUserById(r.gorm, claim["ID"].(string))

	// // Check permissions
	// if !user.IsRoleAdmin() {
	// 	return Forbidden(ctx, "Not allowed", nil)
	// }

	var report model.Report
	var image model.Image

	if err := report.GetById(r.gorm, ctx.Params("id")); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	if err := report.GetAssociateImage(r.gorm, &image); err != nil {
		return Conflict(ctx, err.Error(), nil)
	}

	fmt.Println(image)
	return Ok(ctx, "Report with image", report)
}
