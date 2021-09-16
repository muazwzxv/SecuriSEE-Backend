package model

import (
	"Oracle-Hackathon-BE/util"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type News struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	Headline    string `gorm:"not null" json:"headline"`
	Description string `gorm:"not null" json:"description"`
	Image       string `json:"image"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (n News) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.Headline, validation.Required),
		validation.Field(&n.Description, validation.Required),
		validation.Field(&n.Image, validation.Required),
	)
}

// Gorm hooks
func (n *News) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	n.ID = uuid
	return
}

// CRUD Queries
func (n *News) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(n).Error; err != nil {
		return err
	}
	return nil
}

func (n *News) GetById(gorm *gorm.DB, id string) error {
	if err := gorm.Debug().Where("id = ?", id).First(n).Error; err != nil {
		return err
	}
	return nil
}

func (n *News) GetAll(gorm *gorm.DB, ctx *fiber.Ctx) ([]News, error) {
	var news []News
	if err := gorm.Debug().Scopes(util.Paginate(ctx)).Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}
