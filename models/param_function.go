package models

import "database/sql"

type ParamFunction struct {
	RoutimeName     string         `json:"routine_name" db:"routine_name"`
	ParameterName   string         `json:"parameter_name" db:"parameter_name"`
	DataType        string         `json:"data_type" db:"data_type"`
	OrdinalPosition int            `json:"ordinal_position" db:"ordinal_position"`
	MaxLength       sql.NullInt64  `json:"max_length" db:"max_length"`
	IsNullable      string         `json:"is_nullable" db:"is_nullable"`
	Precision       sql.NullInt64  `json:"precision" db:"precision"`
	Scale           sql.NullInt64  `json:"scale" db:"scale"`
	DefaultValue    sql.NullString `json:"default_value" db:"default_value"`
}
