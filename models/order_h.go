package models

import "time"

type OrderH struct {
	OrderID      int       `json:"order_id" gorm:"primary_key;auto_increment:true"`
	OrderDate    int       `json:"order_date" gorm:"type:integer"`
	CustomerName string    `json:"customer_name" gorm:"type:varchar(60);not null"`
	Telp         string    `json:"telp" gorm:"type:varchar(20)"`
	UserInput    string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit     string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput    time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit     time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}
