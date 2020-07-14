package models

// import uuid "github.com/satori/go.uuid"

type SysUser struct {
	Model
	Email     string `json:"email" gorm:"type:varchar(100);unique_index"`
	Handphone string `json:"handphone" gorm:"type:varchar(20)"`
	Password  string `json:"password"`
	FullName  string `json:"full_name" gorm:"size:100"`
	// FileID    uuid.UUID `json:"file_id" gorm:"type:uuid"`
	IsAdmin bool `json:"is_admin"`
}

type AddUser struct {
	Email     string `json:"email" valid:"Required"`
	Handphone string `json:"handphone"`
	Password  string `json:"password"`
	FullName  string `json:"full_name" valid:"Required"`
	IsAdmin   bool   `json:"is_admin"`
}
