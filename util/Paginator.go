package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func Paginate(ctx *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(ctx.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Helper
func GetClaims(ctx *fiber.Ctx) jwt.MapClaims {
	token := ctx.Locals("user").(*jwt.Token)
	return token.Claims.(jwt.MapClaims)
}

// db.Scopes(Paginate(r)).Find(&users)
// db.Scopes(Paginate(r)).Find(&articles)
