package ipaket

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Paket, err error)
	GetList(queryparam models.ParamList) (result []*models.Paket, err error)
	Create(data *models.Paket) (err error)
	Update(ID int, data map[string]interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result *models.Paket, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.DataPaket) error
	Update(ctx context.Context, Claims util.Claims, ID int, data *models.DataPaket) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
