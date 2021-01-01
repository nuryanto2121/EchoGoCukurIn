package ibarbercapster

import (
	"nuryanto2121/cukur_in_barber/models"
)

type Repository interface {
	GetDataBy(ID int) (result *models.BarberCapster, err error)
	GetList(queryparam models.ParamList) (result []*models.CapsterList, err error)
	Create(data *models.BarberCapster) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	DeleteByCapster(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}
