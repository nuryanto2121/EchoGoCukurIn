package models

import "time"

// import uuid "github.com/satori/go.uuid"

type SysUser struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Telp      string    `json:"telp" db:"telp"`
	Email     string    `json:"email" db:"email"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	JoinDate  time.Time `json:"join_date" db:"join_date"`
	Password  string    `json:"pwd" db:"pwd"`
	FileID    int       `json:"file_id" db:"file_id"`
	UserType  string    `json:"user_type" db:"user_type"`
	UserInput string    `json:"user_input" db:"user_input"`
	UserEdit  string    `json:"user_edit" db:"user_edit"`
	TimeInput time.Time `json:"time_input" db:"time_input"`
	TimeEdit  time.Time `json:"time_edit" db:"time_edit"`
}

type AddUser struct {
	Email     string `json:"email" valid:"Required"`
	Handphone string `json:"handphone"`
	Password  string `json:"password"`
	FullName  string `json:"full_name" valid:"Required"`
	IsAdmin   bool   `json:"is_admin"`
}
