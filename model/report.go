package model

import (
	"Oracle-Hackathon-BE/util"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Report struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	Description string  `json:"description" json:"description"`
	UserID      string  `gorm:"column:user_id" json:"user_id"`
	Lat         float64 `gorm:"type:decimal(10,8)" json:"lat"`
	Lng         float64 `gorm:"type:decimal(11,8)" json:"lng"`
	Status      string  `gorm:"not null" json:"status"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// Relations

	// HasOne
	Image Image `gorm:"foreignKey:ReportID"`
}

const (
	OPEN        = "open"
	ACKNOWLEDGE = "acknowledge"
)

// Gorm hooks
func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	r.ID = uuid
	return
}

// Validator
func (r Report) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.Lat, validation.Required),
		validation.Field(&r.Lng, validation.Required),
		//validation.Field(&r.FileName, validation.Required),
	)
}

// CRUD queries
func (r *Report) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (r *Report) GetAll(gorm *gorm.DB, ctx *fiber.Ctx) ([]Report, error) {
	var report []Report

	if err := gorm.Debug().Scopes(util.Paginate(ctx)).Find(&report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (r *Report) GetById(gorm *gorm.DB, id string) error {
	if err := gorm.Debug().Where("id = ?", id).Find(r).Error; err != nil {
		return err
	}
	return nil
}

// func (r *Report) GetAssociateImage(gorm *gorm.DB, image *Image) error {
// 	if err := gorm.Debug().Model(r).Association("Image").Find(image).Error; err != nil {
// 		return nil
// 	}
// 	return nil
// }

func (r *Report) GetAssociateImage(gorm *gorm.DB, image *Image) error {
	if err := gorm.Debug().Model(r).Joins("left join images on images.report_id  = reports.id").Error; err != nil {
		return nil
	}
	return nil
}

//db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id
