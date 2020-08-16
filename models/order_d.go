package models

import "time"

type OrderD struct {
	OrderDID  int       `json:"order_d_id" gorm:"primary_key;auto_increment:true"`
	BarberID  int       `json:"barber_id" gorm:"type:integer"`
	OrderID   int       `json:"order_id" gorm:"primary_key;type:integer"`
	PaketID   int       `json:"paket_id" gorm:"type:integer;not null"`
	PaketName string    `json:"paket_name" gorm:"type:varchar(60)"`
	Price     float32   `json:"price" gorm:"type:numeric(20,4)"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type OrderDPost struct {
	PaketID   int     `json:"paket_id" gorm:"type:integer;not null"`
	PaketName string  `json:"paket_name" gorm:"type:varchar(60)"`
	Price     float32 `json:"price" gorm:"type:numeric(20,4)"`
}
