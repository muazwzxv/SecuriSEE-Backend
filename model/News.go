package model

import (
	"time"

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
