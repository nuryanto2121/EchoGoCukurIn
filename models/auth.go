package models

import (
	"database/sql"
	"time"
)

//LoginForm :
type LoginForm struct {
	Account  string `json:"account" valid:"Required"`
	Password string `json:"pwd" valid:"Required"`
}

// RegisterForm :
type RegisterForm struct {
	Name     string `json:"name" valid:"Required"`
	UserType string `json:"user_type" valid:"Required"`
	ResetPasswd
	// Account          string    `json:"address,omitempty"`

	// PostCd           string    `json:"post_cd,omitempty"`
	// TelephoneNo      string    `json:"telephone_no,omitempty"`
	// EmailAddr        string    `json:"email_addr,omitempty" valid:"Email"`
	// ContactPerson    string    `json:"contact_person,omitempty"`
	// ClientType       string    `json:"client_type,omitempty"`
	// JoiningDate      time.Time `json:"joining_date,omitempty"`
	// StartBillingDate time.Time `json:"start_billing_date,omitempty"`
	// ExpiryDate       time.Time `json:"expiry_date,omitempty"`
	// CreatedBy        string    `json:"created_by" valid:"Required"`
}

// ForgotForm :
type ForgotForm struct {
	Account string `json:"account" valid:"Required"`
	// EmailAddr string `json:"email,omitempty" valid:"Required;Email"`
}

// ResetPasswd :
type ResetPasswd struct {
	Account       string `json:"account" valid:"Required"`
	Passwd        string `json:"pwd" valid:"Required"`
	ConfirmPasswd string `json:"confirm_pwd" valid:"Required"`
}

type VerifyForm struct {
	Account    string `json:"account" valid:"Required"`
	VerifyCode string `json:"verify_code" valid:"Required"`
}

type DataLogin struct {
	UserID   int            `json:"user_id" db:"user_id"`
	Password string         `json:"pwd" db:"pwd"`
	Name     string         `json:"name" db:"name"`
	Email    string         `json:"email" db:"email"`
	Telp     string         `json:"telp" db:"telp"`
	JoinDate time.Time      `json:"join_date" db:"join_date"`
	UserType string         `json:"user_type" db:"user_type"`
	FileID   sql.NullInt64  `json:"file_id" db:"file_id"`
	FileName sql.NullString `json:"file_name" db:"file_name"`
	FilePath sql.NullString `json:"file_path" db:"file_path"`
}
