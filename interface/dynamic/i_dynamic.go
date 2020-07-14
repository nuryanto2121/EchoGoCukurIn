package idynamic

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
)

type Usecase interface {
	Execute(ctx context.Context, claims util.Claims, data map[string]interface{}) (result interface{}, err error)
	Delete(ctx context.Context, claims util.Claims, ParamGet models.ParamGet) error
	GetDataBy(ctx context.Context, claims util.Claims, ParamGet models.ParamGet) (result interface{}, err error)
	GetList(ctx context.Context, claims util.Claims, queryparam models.ParamDynamicList) (result models.ResponseModelList, err error)
}
type Repository interface {
	GetOptionByUrl(ctx context.Context, Url string) (result []models.OptionDB, err error)
	GetParamFunction(ctx context.Context, SpName string) (result []models.ParamFunction, err error)
	CRUD(ctx context.Context, sQuery string, data interface{}) (result interface{}, err error)
	GetDataList(ctx context.Context, sQuery string, Limit int, Offset int) (result interface{}, err error)
	GetDefineColumn(ctx context.Context, MenuUrl string, LineNo int) (result models.DefineColumn, err error)
	GetFieldType(ctx context.Context, SourceFrom string, isViewFunction bool) (result []models.ParamFunction, err error)
	CountList(ctx context.Context, ViewName string, sWhere string) (int, error)
}
