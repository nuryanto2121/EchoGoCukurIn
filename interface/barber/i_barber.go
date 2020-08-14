package ibarber

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Barber, err error)
	GetList(queryparam models.ParamList) (result []*models.Barber, err error)
	Create(data *models.Barber) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.BarbersPost) error
	Update(ctx context.Context, Claims util.Claims, ID int, data models.BarbersPost) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
