package models

type ParamGet struct {
	ID              int    `json:"id" valid:"Required"`
	Lastupdatestamp int    `json:"lastupdatestamp" query:"lastupdatestamp"  valid:"Required"`
	MenuUrl         string `json:"menu_url" query:"menu_url" valid:"Required"`
	LineNo          int    `json:"line_no" query:"line_no"`
}
