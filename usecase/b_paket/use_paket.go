package usepaket

import (
	"context"
	"math"
	ipaket "nuryanto2121/dynamic_rest_api_go/interface/b_paket"
	"nuryanto2121/dynamic_rest_api_go/models"
	querywhere "nuryanto2121/dynamic_rest_api_go/pkg/query"
	"reflect"
	"time"
)

type usePaket struct {
	repoPaket      ipaket.Repository
	contextTimeOut time.Duration
}

func NewUserMPaket(a ipaket.Repository, timeout time.Duration) ipaket.Usecase {
	return &usePaket{repoPaket: a, contextTimeOut: timeout}
}

func (u *usePaket) GetDataBy(ctx context.Context, ID int) (result *models.Paket, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoPaket.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *usePaket) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var tUser = models.Paket{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		value := reflect.ValueOf(tUser)
		types := reflect.TypeOf(&tUser)
		queryparam.Search = querywhere.GetWhereLikeStruct(value, types, queryparam.Search, "")
	}
	result.Data, err = u.repoPaket.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoPaket.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *usePaket) Create(ctx context.Context, data *models.Paket) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoPaket.Create(data)
	if err != nil {
		return err
	}
	return nil

}
func (u *usePaket) Update(ctx context.Context, ID int, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoPaket.Update(ID, data)
	return nil
}
func (u *usePaket) Delete(ctx context.Context, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoPaket.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
