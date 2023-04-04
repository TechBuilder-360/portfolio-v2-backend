package database

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database. %s", err.Error())
	}

	setDB(db)

	return db
}

func setDB(db *gorm.DB) {
	gormDB = db
}

func DB() *gorm.DB {
	return gormDB
}

func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		model.Account{},
		model.User{},
		model.Education{},
		model.Experience{},
		model.Position{},
		model.Project{},
		model.Social{},
		model.Skill{},
		model.Certificate{},
	)
}
