package models

// Model :
type Model struct {
	ID         int `json:"id" gorm:"primary_key"`
	CreatedOn  int `json:"created_on,omitempty"`
	ModifiedOn int `json:"modified_on,omitempty"`
	DeletedOn  int `json:"deleted_on,omitempty"`
	// IDCreated  int `json:"id_created"`
}
