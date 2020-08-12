package models

import "time"

type SaFileUpload struct {
	FileID    int       `json:"file_id" gorm:"primary_key;auto_increment:true"`
	FileName  string    `json:"file_name" gorm:"type:varchar(60)"`
	FilePath  string    `json:"file_path" gorm:"type:varchar(150)"`
	FileType  string    `json:"file_type" gorm:"type:varchar(10)"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

type SaFileOutput struct {
	FileID   int    `json:"file_id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
