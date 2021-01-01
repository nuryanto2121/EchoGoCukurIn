package inotification

import (
	"context"
	"nuryanto2121/cukur_in_barber/models"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Notification, err error)
	GetList(queryparam models.ParamList) (result []*models.Notification, err error)
	Create(data *models.Notification) (err error)
	Update(ID int, data map[string]interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result *models.Notification, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, TokenFCM string, data *models.AddNotification) (err error)
	Update(ctx context.Context, Claims util.Claims, ID int, data *models.StatusNotification) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
	GetCountNotif(ctx context.Context, Claims util.Claims) (result interface{}, err error)
}
