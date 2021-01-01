package icapster

import (
	"context"
	"nuryanto2121/cukur_in_barber/models"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.CapsterCollection, err error)
	GetListFileCapter(ID int) (result []*models.SaFileOutput, err error)
	GetList(queryparam models.ParamList) (result []*models.CapsterList, err error)
	Create(data *models.CapsterCollection) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.Capster) error
	Update(ctx context.Context, Claims util.Claims, ID int, data *models.Capster) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
