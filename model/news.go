package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
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
