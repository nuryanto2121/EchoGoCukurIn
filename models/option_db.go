package models

import "time"

type OptionDB struct {
	OptionID  int       `json:"option_id" db:"option_id"`
	OptionUrl string    `json:"option_url" db:"option_url"`
	MethodApi string    `json:"method_api" db:"method_api"`
	SP        string    `json:"sp" db:"sp"`
	LineNo    int       `json:"line_no" db:"line_no"`
	TableName string    `json:"table_name" db:"table_name"`
	UserInput string    `json:"user_input" db:"user_input"`
	UserEdit  string    `json:"user_edit" db:"user_edit"`
	TimeInput time.Time `json:"time_input" db:"time_input"`
	TimeEdit  time.Time `json:"time_edit" db:"time_edit"`
}

type DefineColumn struct {
	ColumnField string `json:"column_field" db:"column_field"`
}
