package controller

import "gorm.io/gorm"

type NewsRepository struct {
	gorm *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{gorm: db}
}
