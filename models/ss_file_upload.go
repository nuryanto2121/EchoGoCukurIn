package models

type SaFileUpload struct {
	FileID   int    `json:"file_id" gorm:"primary_key;auto_increment:true"`
	FileName string `json:"file_name" gorm:"type:varchar(60)"`
	FilePath string `json:"file_path" gorm:"type:varchar(150)"`
	FileType string `json:"file_type" gorm:"type:varchar(10)"`
	Model
}

type SaFileOutput struct {
	FileID   int    `json:"file_id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
