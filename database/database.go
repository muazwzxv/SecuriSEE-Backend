package database

import "gorm.io/gorm"

type DatabaseConfig struct {
	User         string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

type GormInstance struct {
	orm *gorm.DB
}

var (
	GORM = &GormInstance{}
)

func (g *GormInstance) Connect() *GormInstance {
	return &GormInstance{}
}
