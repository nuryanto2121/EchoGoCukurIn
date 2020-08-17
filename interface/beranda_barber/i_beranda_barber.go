package iberandabarber

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Usecase interface {
	GetStatusOrder(ctx context.Context, Claims util.Claims) (interface{}, error)
	GetListOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (interface{}, error)
}

type Repository interface {
	GetStatusOrder(BarberID int) (result models.Beranda, err error)
	GetListOrder(queryparam models.ParamList) (result []*models.BerandaList, err error)
	Count(queryparam models.ParamList) (result int, err error)
}
