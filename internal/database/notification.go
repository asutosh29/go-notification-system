package database

import "gorm.io/gorm"

type Notification struct {
	gorm.Model

	NotificationResp
}

type NotificationResp struct {
	Title       string `gorm:"not null;unique" json:"title"`
	Description string `gorm:"not null" json:"description"`
}
