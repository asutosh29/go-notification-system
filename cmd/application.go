package main

import (
	"github.com/asutosh29/go-gin/internal/config"
	"gorm.io/gorm"
)

type Application struct {
	db     *gorm.DB
	config config.Config
}
