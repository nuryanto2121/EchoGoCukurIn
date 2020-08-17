package useorder

import (
	"context"
	"fmt"
	"math"
	iorderd "nuryanto2121/dynamic_rest_api_go/interface/c_order_d"
	iorderh "nuryanto2121/dynamic_rest_api_go/interface/c_order_h"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useOrder struct {
	repoOrderH     iorderh.Repository
	repoOrderD     iorderd.Repository
	contextTimeOut time.Duration
}

func NewUserMOrder(a iorderh.Repository, b iorderd.Repository, timeout time.Duration) iorderh.Usecase {
	return &useOrder{
		repoOrderH:     a,
		repoOrderD:     b,
		contextTimeOut: timeout}
}

func (u *useOrder) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err := u.repoOrderH.GetDataBy(ID)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (u *useOrder) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = fmt.Sprintf("lower(order_name) LIKE '%%%s%%' ", queryparam.Search)
	}

	queryparam.InitSearch = fmt.Sprintf("barber.owner_id = %s", Claims.UserID)
	result.Data, err = u.repoOrderH.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoOrderH.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useOrder) Create(ctx context.Context, Claims util.Claims, data *models.OrderPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mOrder models.OrderH
	)
	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOrder)
	if err != nil {
		return err
	}
	mOrder.Status = "N"
	mOrder.FromApps = false
	if mOrder.CapsterID == 0 {
		mOrder.CapsterID, _ = strconv.Atoi(Claims.UserID)
	}

	mOrder.UserInput = Claims.UserID
	mOrder.UserEdit = Claims.UserID
	err = u.repoOrderH.Create(&mOrder)
	if err != nil {
		return err
	}

	for _, dataDetail := range data.Pakets {
		var mDetail models.OrderD
		err = mapstructure.Decode(dataDetail, &mDetail)
		if err != nil {
			return err
		}
		mDetail.BarberID = mOrder.BarberID
		mDetail.OrderID = mOrder.OrderID
		mDetail.UserEdit = Claims.UserID
		mDetail.UserInput = Claims.UserID
		err = u.repoOrderD.Create(&mDetail)
		if err != nil {
			return err
		}
	}

	return nil

}
func (u *useOrder) Update(ctx context.Context, Claims util.Claims, ID int, data models.OrderPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		mOrder models.OrderH
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOrder)
	if err != nil {
		return err
	}
	err = u.repoOrderH.Update(ID, mOrder)
	if err != nil {
		return err
	}

	//delete then insert detail

	err = u.repoOrderD.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
func (u *useOrder) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoOrderH.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
