package ibarberpaket

import (
	"nuryanto2121/dynamic_rest_api_go/models"
)

type Repository interface {
	GetDataBy(ID int) (result *models.BarberPaket, err error)
	GetList(queryparam models.ParamList) (result []*models.BarberPaket, err error)
	Create(data *models.BarberPaket) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}
