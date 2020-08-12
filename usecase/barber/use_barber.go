package usebarber

import (
	"context"
	"math"
	ibarber "nuryanto2121/dynamic_rest_api_go/interface/barber"
	ibarbercapster "nuryanto2121/dynamic_rest_api_go/interface/barber_capster"
	ibarberpaket "nuryanto2121/dynamic_rest_api_go/interface/barber_paket"
	"nuryanto2121/dynamic_rest_api_go/models"
	querywhere "nuryanto2121/dynamic_rest_api_go/pkg/query"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useBarber struct {
	repoBarber        ibarber.Repository
	repoBarberPaket   ibarberpaket.Repository
	repoBarberCapster ibarbercapster.Repository
	contextTimeOut    time.Duration
}

func NewUserMBarber(a ibarber.Repository, b ibarberpaket.Repository, c ibarbercapster.Repository, timeout time.Duration) ibarber.Usecase {
	return &useBarber{
		repoBarber:        a,
		repoBarberPaket:   b,
		repoBarberCapster: c,
		contextTimeOut:    timeout}
}

func (u *useBarber) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoBarber.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useBarber) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var tUser = models.Barber{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		value := reflect.ValueOf(tUser)
		types := reflect.TypeOf(&tUser)
		queryparam.Search = querywhere.GetWhereLikeStruct(value, types, queryparam.Search, "")
	}
	result.Data, err = u.repoBarber.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBarber.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useBarber) Create(ctx context.Context, Claims util.Claims, data *models.BarbersPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mBarber models.Barber
	)
	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mBarber)
	if err != nil {
		return err
	}
	mBarber.OwnerID, _ = strconv.Atoi(Claims.UserID)
	mBarber.UserInput = Claims.UserID
	mBarber.UserEdit = Claims.UserID
	err = u.repoBarber.Create(&mBarber)
	if err != nil {
		return err
	}

	for _, dataCapster := range data.BarberCapster {
		var BCapster = models.BarberCapster{}
		BCapster.BarberID = mBarber.BarberID
		BCapster.CapsterID = dataCapster.CapsterID
		BCapster.UserInput = Claims.UserID
		BCapster.UserEdit = Claims.UserID
		err = u.repoBarberCapster.Create(&BCapster)
		if err != nil {
			return err
		}
	}

	for _, dataCapster := range data.BarberPaket {
		var BPaket = models.BarberPaket{}
		BPaket.BarberID = mBarber.BarberID
		BPaket.PaketID = dataCapster.PaketID
		BPaket.UserInput = Claims.UserID
		BPaket.UserEdit = Claims.UserID
		err = u.repoBarberPaket.Create(&BPaket)
		if err != nil {
			return err
		}
	}

	return nil

}
func (u *useBarber) Update(ctx context.Context, Claims util.Claims, ID int, data interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoBarber.Update(ID, data)
	return nil
}
func (u *useBarber) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoBarber.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
