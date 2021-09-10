package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid:default:uuid_generate_v4()"`
	Ic       string    `gorm:"index"`
	Name     string    `gorm:"not null"`
	Phone    string    `gorm:"not null"`
	Email    string
	Role     string `gorm:"not null"`
	Password string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

// sql wrapper

func (u *User) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(u).Error; err != nil {
		return err
	}
	return nil
}

// Helpers

func (u *User) IsEmailExist(gorm *gorm.DB) bool {
	if res := gorm.Debug().Where("email = ?", u.Email).First(u); res != nil && res.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func (u *User) IsICExist(gorm *gorm.DB) bool {
	if res := gorm.Debug().Where("ic = ?", u.Ic).First(u); res != nil && res.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func (u *User) HashPassword(p string) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14); err != nil {
		panic(err)
	} else {
		u.Password = string(bytes)
	}
}

func (u *User) CheckHash(pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	return err == nil
}

func (u *User) IsRoleExist(role string) bool {

	split := strings.Split(u.Role, ",")
	for _, val := range split {
		if val == role {
			return true
		}
	}
	return false
}

func (u *User) RolesToString(roles []string) string {
	// use this to serialize slice string to string
	return strings.Join(roles, ", ")
}
