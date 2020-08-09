package models

import "time"

type Capster struct {
	Name string `json:"name" valid:"Required"`
	// Telp      string    `json:"telp"`
	// Email     string    `json:"email"`
	UserType      string    `json:"user_type" valid:"Required"`
	IsActive      bool      `json:"is_active" valid:"Required"`
	JoinDate      time.Time `json:"join_date,omitempty"`
	FileID        int       `json:"file_id,omitempty"`
	TopCollection []Collections
}

type CapsterCollection struct {
	CapsterID int       `json:"capster_id" gorm:"primary_key;type:integer"`
	FileID    int       `json:"file_id" gorm:"primary_key;type:integer"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

// type CapsterCollectionPath struct {
// 	CapsterCollection
// 	FileName string `json:"file_name"`
// 	FilePath string `json:"file_path"`
// }

type Collections struct {
	FileID int `json:"file_id"`
}

type CapsterList struct {
	CapsterID int    `json:"capster_id" valid:"Required"`
	Name      string `json:"name" valid:"Required"`
	IsActive  bool   `json:"is_active" valid:"Required"`
	SaFileOutput
	Rating float32 `json:"rating,omitempty"`
}
