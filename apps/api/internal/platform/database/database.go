package database

import (
	"link-nest/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	db, err := gorm.Open(postgres.Open(config.Current.DatabaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
