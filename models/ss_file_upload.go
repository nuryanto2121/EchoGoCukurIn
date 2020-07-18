package models

import "time"

type SaFileUpload struct {
	FileID    int       `json:"file_id" db:"file_id"`
	FileName  string    `json:"file_name" db:"file_name"`
	FilePath  string    `json:"file_path" db:"file_path"`
	FileType  string    `json:"file_type" db:"file_type"`
	UserInput string    `json:"user_input" db:"user_input"`
	TimeInput time.Time `json:"time_input" db:"time_input"`
	UserEdit  string    `json:"user_edit" db:"user_edit"`
	TimeEdit  time.Time `json:"time_edit" db:"time_edit"`
}

type SaFileOutput struct {
	FileID   int    `json:"file_id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
