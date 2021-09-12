package model

import (
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	//ID       uuid.UUID `gorm:"type:uuid:default:uuid_generate_v4()"`
	ID       uuid.UUID `gorm:"type:char(36);primary_key"`
	Ic       string    `gorm:"not null"`
	Name     string    `gorm:"not null"`
	Phone    string    `gorm:"not null"`
	Email    string
	Role     string `gorm:"not null"`
	Password string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

// Struct to Login
type Login struct {
	IC       string `json:"ic"`
	Password string `json:"password"`
}

// Gorm hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	u.ID = uuid
	return
}

// sql wrapper

func (u *User) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserById(gorm *gorm.DB, id string) error {
	if res := gorm.Debug().Select("id", "ic", "name", "phone", "email", "role", "created_at", "deleted_at").Where("id = ?", id).First(u); res.Error != nil {
		return res.Error
	}
	return nil
}

func (u *User) GetUserByIc(gorm *gorm.DB, ic string) error {
	if res := gorm.Debug().Where("ic = ?", ic).First(u); res.Error != nil {
		return res.Error
	}
	return nil
}

func (u *User) Get(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(u).Error; err != nil {
		return err
	}
	return nil
}

// Helpers

func (u *User) IsEmailExist(gorm *gorm.DB) bool {
	if res := gorm.Debug().Select("email").Where("email = ?", u.Email).First(u); res != nil && res.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func (u *User) IsICExist(gorm *gorm.DB) bool {
	if res := gorm.Debug().Select("ic").Where("ic = ?", u.Ic).First(u); res != nil && res.RowsAffected == 0 {
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

func (u *User) RolesToArray() []string {
	return strings.Split(u.Role, ",")
}

func (u *User) RolesToString(roles []string) {
	// use this to serialize slice string to string
	u.Role = strings.Join(roles, ", ")
}
