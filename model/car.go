package model

import (
	"Oracle-Hackathon-BE/util"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	// Car Entry status
	INBOUND  = "inbound"
	OUTBOUND = "outbound"
)

type CarEntry struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	PlateNumber string `gorm:"not null" json:"plate_number"`
	City        string `gorm:"not null" json:"city"`
	Status      string `gorm:"not null" json:"status"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (c CarEntry) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.PlateNumber, validation.Required),
		validation.Field(&c.City, validation.Required),
	)
}

// Gorm hooks
func (c *CarEntry) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	c.ID = uuid
	return
}

// Queries
func (c *CarEntry) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (c *CarEntry) GetAll(gorm *gorm.DB, ctx *fiber.Ctx) ([]CarEntry, error) {
	var entry []CarEntry
	if err := gorm.Debug().Scopes(util.Paginate(ctx)).Find(&entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (c *CarEntry) GetEntryById(gorm *gorm.DB, id string) error {
	if err := gorm.Debug().Where("id = ?", id).First(c).Error; err != nil {
		return err
	}
	return nil
}

// Helpers

func (c *CarEntry) CheckStatus() error {
	switch {
	case c.Status == "":
		return errors.New("include status")
	case c.Status != OUTBOUND && c.Status != INBOUND:
		return errors.New("status is not valid")
	}
	return nil
}
