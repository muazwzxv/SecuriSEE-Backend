package service

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormInstance struct {
	Orm *gorm.DB
}

var GORM GormInstance

func ConnectDatabase() GormInstance {
	if GORM.Orm == nil {
		config := config.CFG.FetchDatabaseConfig()
		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DatabaseName,
		)

		if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
			panic(fmt.Sprintf("Failed to connect to database: \n %v", err))
		} else {
			// Migrate all tables
			db.Debug().AutoMigrate(
				&model.User{},
			)
			GORM = GormInstance{
				Orm: db,
			}

			return GORM
		}
	}

	return GORM
}
