package ipaket

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Paket, err error)
	GetList(queryparam models.ParamList) (result []*models.Paket, err error)
	Create(data *models.Paket) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, ID int) (result *models.Paket, err error)
	GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, data *models.Paket) error
	Update(ctx context.Context, ID int, data interface{}) (err error)
	Delete(ctx context.Context, ID int) (err error)
}
