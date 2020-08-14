package models

import "time"

type Barber struct {
	BarberID       int       `json:"barber_id" gorm:"primary_key;auto_increment:true"`
	OwnerID        int       `json:"owner_id" gorm:"PRIMARY_KEY;type:integer"`
	BarberName     string    `json:"barber_name" gorm:"type:varchar(60)"`
	Address        string    `json:"address" gorm:"type:varchar(150)"`
	FileID         int       `json:"file_id" gorm:"type:integer"`
	PinMap         string    `json:"pin_map" gorm:"type:varchar(100)"`
	Starts         float32   `json:"starts" gorm:"type:numeric(5,2)"`
	OperationStart int       `json:"operation_start" gorm:"type:integer"`
	OperationEnd   int       `json:"operation_end" gorm:"type:integer"`
	IsActive       bool      `json:"is_active" gorm:"type:boolean"`
	UserInput      string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit       string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput      time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit       time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type BarbersUpdate struct {
	BarberName     string  `json:"barber_name" valid:"Required"`
	Address        string  `json:"address" valid:"Required"`
	FileID         int     `json:"file_id,omitempty"`
	PinMap         string  `json:"pin_map,omitempty"`
	Starts         float32 `json:"starts,omitempty"`
	OperationStart int     `json:"operation_start" valid:"Required"`
	OperationEnd   int     `json:"operation_end" valid:"Required"`
	IsActive       bool    `json:"is_active" valid:"Required"`
}

type BarbersPost struct {
	BarberName     string              `json:"barber_name" valid:"Required"`
	Address        string              `json:"address" valid:"Required"`
	FileID         int                 `json:"file_id,omitempty"`
	PinMap         string              `json:"pin_map,omitempty"`
	Starts         float32             `json:"starts,omitempty"`
	OperationStart int                 `json:"operation_start" valid:"Required"`
	OperationEnd   int                 `json:"operation_end" valid:"Required"`
	IsActive       bool                `json:"is_active" valid:"Required"`
	BarberPaket    []BarberPaketPost   `json:"barber_paket"`
	BarberCapster  []BarberCapsterPost `json:"barber_capster"`
}

type BarberPaketPost struct {
	PaketID int `json:"paket_id"`
}

type BarberPaket struct {
	BarberID  int       `json:"barber_id" gorm:"primary_key;type:integer"`
	PaketID   int       `json:"paket_id" gorm:"primary_key;type:integer"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type BarberCapsterPost struct {
	CapsterID int `json:"capster_jd"`
}

type BarberCapster struct {
	BarberID  int       `json:"barber_id" gorm:"primary_key;type:integer"`
	CapsterID int       `json:"capster_jd" gorm:"primary_key;type:integer"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}
