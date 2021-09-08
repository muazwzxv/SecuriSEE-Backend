package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID    uint64 `gorm:"primary_key:auto_increment"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Ic    string `gorm:"not null"`
	Role  string

	CreatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}
