package models

import (
	"time"
)

type Notification struct {
	NotificationId int `json:"notification_id" gorm:"primary_key;auto_increment:true"`
	AddNotification
	Model
}

type AddNotification struct {
	NotificationDate   time.Time `json:"notification_date" gorm:"type:timestamp(0) without time zone;not null"`
	NotificationStatus string    `json:"notification_status" gorm:"type:varchar(1);not null"`
	NotificationType   string    `json:"notification_type" gorm:"type:varchar(1);not null"`
	UserId             int       `json:"user_id" gorm:"type:integer;not null"`
	Title              string    `json:"title" gorm:"type:varchar(150);not null"`
	Descs              string    `json:"descs" gorm:"type:varchar(255);not null"`
	LinkUrl            string    `json:"link_url" gorm:"type:varchar(60)"`
	LinkId             int       `json:"link_id" gorm:"type:integer"`
}

type StatusNotification struct {
	NotificationStatus string `json:"notification_status" gorm:"type:varchar(1);not null"`
}
