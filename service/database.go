package service

import (
	"Oracle-Hackathon-BE/config"
	"Oracle-Hackathon-BE/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db GormInstance

type GormInstance struct {
	Orm *gorm.DB
}

func newGormInstance() GormInstance {

	config := config.GetInstance().FetchDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)

	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Skip default transaction to improve speed
		SkipDefaultTransaction: true,
		// Cache statements
		PrepareStmt: true,
	}); err != nil {
		// Panic
		panic(fmt.Sprintf("Failed to connect to database: \n %v", err))
	} else {
		return GormInstance{db}
	}
}

func GetGormInstance() *GormInstance {
	if !db.isInstantiated() {
		db = newGormInstance()
	}
	return &db
}

func (g *GormInstance) isInstantiated() bool {
	return g.Orm != nil
}

func (g *GormInstance) GormShutDown() error {

	if sql, err := g.Orm.DB(); err != nil {
		return err
	} else {
		sql.Close()
		return nil
	}
}

func (g *GormInstance) Migrate() error {

	if err := g.Orm.Debug().AutoMigrate(
		&model.User{},
		&model.Report{},
		&model.News{},
		&model.CarEntry{},
		&model.Image{},
	); err != nil {
		return err
	}
	return nil
}

// func ConnectDatabase() *GormInstance {
// 	config := config.GetInstance().FetchDatabaseConfig()
// 	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
// 		config.User,
// 		config.Password,
// 		config.Host,
// 		config.Port,
// 		config.DatabaseName,
// 	)

// 	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		// Skip default transaction to improve speed
// 		SkipDefaultTransaction: true,
// 		// Cache statements
// 		PrepareStmt: true,
// 	}); err != nil {
// 		panic(fmt.Sprintf("Failed to connect to database: \n %v", err))
// 	} else {
// 		// Migrate all tables
// 		db.Debug().AutoMigrate(
// 			&model.User{},
// 			&model.Report{},
// 			&model.News{},
// 			&model.CarEntry{},
// 			&model.Image{},
// 		)

// 		return &GormInstance{orm: db}
// 	}

// }
