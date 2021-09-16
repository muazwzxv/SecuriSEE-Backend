package controller

import (
	"Oracle-Hackathon-BE/model"
	"Oracle-Hackathon-BE/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReportRepository struct {
	gorm *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{gorm: db}
}

func (reportRepository *ReportRepository) Create(ctx *fiber.Ctx) error {
	// validate role
	claim := util.GetClaims(ctx)
	var user model.User
	user.GetUserById(reportRepository.gorm, claim["ID"].(string))

	// Check permissions
	if isUser := user.IsRoleExist("user"); !isUser {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"Success": false,
			"Message": "Not Allowed",
		})
	}

	// parse json
	var report model.Report
	if err := ctx.BodyParser(&report); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": "Cannot parse JSON",
		})
	}

	// Assign reference id
	report.UserID = user.ID

	if err := report.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	if err := report.Create(reportRepository.gorm); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"Success": false,
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"Success": true,
		"Message": "Report created",
		"Report":  report,
	})
}

func (reportRepository *ReportRepository) GetAll(ctx *fiber.Ctx) error {
	return nil
}

func (reportRepository *ReportRepository) GetById(ctx *fiber.Ctx) error {
	return nil
}
