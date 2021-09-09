package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid:default:uuid_generate_v4()"`
	Name  string    `gorm:"not null"`
	Email string    `gorm:"not null"`
	Phone string    `gorm:"not null"`
	Ic    string    `gorm:"not null"`
	Role  string    `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (u *User) isRoleExist(role string) bool {

	split := strings.Split(u.Role, ",")
	for _, val := range split {
		if val == role {
			return true
		}
	}

	return false
}

func (u *User) rolesToString(roles []string) string {
	// use this to serialize slice string to string
	return strings.Join(roles, ", ")
}
