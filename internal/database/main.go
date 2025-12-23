package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db    *gorm.DB
	DbMux sync.Mutex
	DbErr error
)

func InitDatabase(config config.Config) {
	Db, DbErr = ConnectDb(config)
	if DbErr != nil {
		log.Fatal("Error connecting to Database: ", DbErr)
	}
}

func ConnectDb(config config.Config) (*gorm.DB, error) {
	// DB Setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

func Migrate() {
	log.Println("Migrating tables...")
	DbErr = Db.AutoMigrate(&Notification{})
	if DbErr != nil {
		log.Fatal("Error migrating tables: ", DbErr)
	}
}
