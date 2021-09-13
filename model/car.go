package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CarDetect struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	PlateNumber string `gorm:"not null" json:"plate_number"`
	City        string `gorm:"not null" json:"city"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (c CarDetect) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.PlateNumber, validation.Required),
		validation.Field(&c.City, validation.Required),
	)
}
