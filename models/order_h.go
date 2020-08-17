package models

import "time"

type OrderH struct {
	OrderID      int       `json:"order_id" gorm:"primary_key;auto_increment:true"`
	BarberID     int       `json:"barber_id" gorm:"type:integer"`
	CapsterID    int       `json:"capster_id" gorm:"type:integer"`
	OrderDate    int       `json:"order_date" gorm:"type:integer"`
	UserID       int       `json:"user_id" gorm:"type:integer"`
	CustomerName string    `json:"customer_name" gorm:"type:varchar(60);not null"`
	Telp         string    `json:"telp" gorm:"type:varchar(20)"`
	Status       string    `json:"status" gorm:"type:varchar(1)"`
	FromApps     bool      `json:"from_apps" gorm:"type:boolean"`
	UserInput    string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit     string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput    time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit     time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type OrderPost struct {
	BarberID     int          `json:"barber_id" valid:"Required"`
	CapsterID    int          `json:"capster_id,omitempty"`
	OrderDate    int          `json:"order_date" valid:"Required"`
	UserID       int          `json:"user_id,omitempty"`
	CustomerName string       `json:"customer_name" valid:"Required"`
	Telp         string       `json:"telp,omitempty"`
	Pakets       []OrderDPost `json:"paket_ids"`
}

type OrderList struct {
	BarberID    int     `json:"barber_id" valid:"Required"`
	BarberName  string  `json:"barber_name"`
	OrderID     int     `json:"order_id"`
	Status      string  `json:"status" `
	FromApps    bool    `json:"from_apps"`
	CapsterID   int     `json:"capster_id,omitempty"`
	CapsterName string  `json:"capster_name"`
	OrderDate   int     `json:"order_date" valid:"Required"`
	FileID      int     `json:"file_id" `
	FileName    string  `json:"file_name"`
	FilePath    string  `json:"file_path"`
	Price       float32 `json:"price" `
}
