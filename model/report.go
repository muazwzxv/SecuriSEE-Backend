package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Report struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	Description string  `json:"description" json:"description"`
	UserID      string  `gorm:"column:user_id" json:"user_id"`
	Lat         float64 `gorm:"type:decimal(10,8)" json:"lat"`
	Lng         float64 `gorm:"type:decimal(11,8)" json:"lng"`
	FileName    string  `json:"fileName"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// relationship
	User User
}

func (r Report) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.Lat, validation.Required),
		validation.Field(&r.Lng, validation.Required),
	)
}
