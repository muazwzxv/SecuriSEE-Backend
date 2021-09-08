package database

import (
	"Oracle-Hackathon-BE/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormInstance struct {
	orm *gorm.DB
}

var (
	GORM = &GormInstance{}
)

func Connect() (*GormInstance, error) {
	config := config.CFG.FetchDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)

	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		// Migrate all tables
		db.Debug().AutoMigrate()
		GORM = &GormInstance{
			orm: db,
		}

		return GORM, nil
	}
}
