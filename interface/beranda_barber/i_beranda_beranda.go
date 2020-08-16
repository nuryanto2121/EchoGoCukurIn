package iberandabarber

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Usecase interface {
	GetStatusOrder(ctx context.Context, Claims util.Claims) (interface{}, error)
	GetListOrder(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (models.ResponseModelList, error)
}
