package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Title       string `gorm:"not null;unique" json:"title"`
	Description string `gorm:"not null" json:"description"`
}

type NotificationResp struct {
	Title       string `gorm:"not null;unique" json:"title"`
	Description string `gorm:"not null" json:"description"`
}

func AddNotification(notification Notification) (Notification, error) {
	tx := Db.Create(&notification)

	return notification, tx.Error
}

func AllNotification() ([]Notification, error) {
	var notifications []Notification
	tx := Db.Find(&notifications)

	return notifications, tx.Error
}

func GetNotification(id int) (Notification, error) {
	var notification Notification
	tx := Db.First(&notification, id)

	return notification, tx.Error
}

func DeleteNotification(id int) error {
	tx := Db.First(&Notification{}, id)
	if tx.Error != nil {
		return fmt.Errorf("Error deleting notification: %v", tx.Error)
	}
	tx = Db.Unscoped().Delete(&Notification{}, id)
	return tx.Error
}
