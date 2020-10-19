package iberandabarber

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Usecase interface {
	GetStatusOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList) (interface{}, error)
	GetListOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamDynamicList) (interface{}, error)
}

type Repository interface {
	GetStatusOrder(ParamView string, OwnerID int) (result models.Beranda, err error)
	GetListOrder(queryparam models.ParamDynamicList) (result []*models.BerandaList, err error)
	Count(queryparam models.ParamDynamicList) (result int, err error)
}
