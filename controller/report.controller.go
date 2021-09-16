package controller

import (
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
	return nil
}

func (reportRepository *ReportRepository) GetAll(ctx *fiber.Ctx) error {
	return nil
}

func (reportRepository *ReportRepository) GetById(ctx *fiber.Ctx) error {
	return nil
}
