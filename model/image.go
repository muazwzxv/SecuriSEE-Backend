package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID       uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	FileName string    `gorm:"not null" json:"file_name"`
	Path     string    `gorm:"not null" json:"path"`
	UserID   uuid.UUID `gorm:"column:user_id" json:"user_id"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// relationship
	User User
}

func (i Image) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.FileName, validation.Required),
		validation.Field(&i.Path, validation.Required),
	)
}
