package database

import (
	"fmt"

	"github.com/asutosh29/go-gin/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb(config config.Config) (*gorm.DB, error) {
	// DB Setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
