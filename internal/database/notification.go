package database

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model

	NotificationResp
}

type NotificationResp struct {
	Title       string `gorm:"not null;unique" json:"title"`
	Description string `gorm:"not null" json:"description"`
}

func AddNotification(notification NotificationResp) (NotificationResp, error) {
	tx := Db.Create(&Notification{
		NotificationResp: notification,
	})

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
	tx := Db.Unscoped().Delete(&Notification{}, id)

	return tx.Error
}
