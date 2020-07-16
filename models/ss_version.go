package models

type VersionApps struct {
	OS      string `json:"os" db:"os"`
	Version int    `json:"version" db:"version"`
}
