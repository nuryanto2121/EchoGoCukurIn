package models

import "time"

type Paket struct {
	PaketID     int       `json:"paket_id" gorm:"primary_key;auto_increment:true"`
	OwnerID     int       `json:"owner_id" gorm:"PRIMARY_KEY;type:integer"`
	PaketName   string    `json:"paket_name" gorm:"type:varchar(60)"`
	Descs       string    `json:"descs" gorm:"type:varchar(200)"`
	DurasiStart int       `json:"durasi_start" gorm:"type:integer"`
	DurasiEnd   int       `json:"durasi_end" gorm:"type:integer"`
	Price       float32   `json:"price" gorm:"type:numeric(20,4)"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean"`
	IsPromo     bool      `json:"is_promo" gorm:"type:boolean"`
	PromoPrice  float32   `json:"promo_price" gorm:"type:numeric(20,4)"`
	PromoStart  int       `json:"promo_start" gorm:"type:integer"`
	PromoEnd    int       `json:"promo_end" gorm:"type:integer"`
	UserInput   string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit    string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput   time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit    time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type DataPaket struct {
	OwnerID     int       `json:"owner_id" valid:"Required"`
	PaketName   string    `json:"paket_name" valid:"Required"`
	Descs       string    `json:"descs,omitempty"`
	DurasiStart int       `json:"durasi_start,omitempty"`
	DurasiEnd   int       `json:"durasi_end,omitempty"`
	Price       float32   `json:"price" valid:"Required"`
	IsActive    bool      `json:"is_active"`
	IsPromo     bool      `json:"is_promo,omitempty"`
	PromoPrice  float32   `json:"promo_price,omitempty"`
	PromoStart  int       `json:"promo_start,omitempty"`
	PromoEnd    int       `json:"promo_end,omitempty"`
	UserEdit    string    `json:"user_edit,omitempty"`
	TimeEdit    time.Time `json:"time_Edit,omitempty"`
}
