package database

import (
	"fmt"
	"log"

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

func Migrate() {
	// Migration functions
	config := config.InitConfig()
	db, err := ConnectDb(config)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}

	log.Println("Migrating tables...")
	err = db.AutoMigrate(&Notification{})
	if err != nil {
		log.Fatal("Error migrating tables: ", err)
	}
}
