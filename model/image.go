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
	UserID   uuid.UUID `gorm:"column:user_id" json:"user_id"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// relationship
	User User
}

func (i Image) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.FileName, validation.Required),
	)
}

// Gorm hooks
func (i *Image) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	i.ID = uuid
	return
}

// CRUD Queries
func (i *Image) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(i).Error; err != nil {
		return err
	}
	return nil
}

func (i *Image) GetById(gorm *gorm.DB, id string) error {
	if err := gorm.Debug().Where("id = ?", id).First(i).Error; err != nil {
		return err
	}
	return nil
}
