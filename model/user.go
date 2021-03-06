package model

import (
	"Oracle-Hackathon-BE/util"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	//ID       uuid.UUID `gorm:"type:uuid:default:uuid_generate_v4()"`
	ID       uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Ic       string    `gorm:"not null" json:"ic"`
	Name     string    `gorm:"not null" json:"name"`
	Phone    string    `gorm:"not null" json:"phone"`
	Email    string    `json:"email"`
	Role     string    `gorm:"not null" json:"role"`
	Password string    `gorm:"not null" json:"password"`

	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// Relations

	// Has Many
	Reports []Report `gorm:"foreignKey:UserID"`
}

const (
	ADMIN  = "admin"
	USER   = "user"
	CAMERA = "camera"
)

// Struct to Login
type LoginUser struct {
	IC       string `json:"ic"`
	Password string `json:"password"`
}

type LoginAdminAndCamera struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validator
func (u User) ValidateCreate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Ic, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}

// Gorm hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	u.ID = uuid
	return
}

// CRUD Queries
func (u *User) Update(gorm *gorm.DB, toUpdate User) error {
	if err := gorm.Debug().Model(u).Select("*").Omit("ic", "role", "name", "password").Updates(toUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) Create(gorm *gorm.DB) error {
	if err := gorm.Debug().Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserById(gorm *gorm.DB, id string) error {
	if res := gorm.Debug().Scopes(selectUser).Where("id = ?", id).First(u); res.Error != nil {
		return res.Error
	}
	return nil
}

func (u *User) GetUserByIc(gorm *gorm.DB, ic string) error {
	if err := gorm.Debug().Where("ic = ?", ic).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByEmail(gorm *gorm.DB, email string) error {
	if err := gorm.Debug().Where("email = ?", email).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetAll(gorm *gorm.DB, ctx *fiber.Ctx) ([]User, error) {
	var user []User
	if err := gorm.Debug().Scopes(util.Paginate(ctx), selectUser).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetAssociateReports(gorm *gorm.DB) ([]Report, error) {
	var reports []Report

	if err := gorm.Debug().Model(&Report{}).Preload(clause.Associations).Where("user_id = ?", u.ID).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil

}

// Scope wrapper
func selectUser(db *gorm.DB) *gorm.DB {
	return db.Select("id", "ic", "name", "phone", "email", "role", "created_at", "deleted_at")
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

func (u *User) IsRoleCamera() bool {
	return u.IsRoleExist("camera")
}

func (u *User) IsRoleUser() bool {
	return u.IsRoleExist("user")
}

func (u *User) IsRoleAdmin() bool {
	return u.IsRoleExist("admin")
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
