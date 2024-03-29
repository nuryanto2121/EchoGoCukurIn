package usepaket

import (
	"context"
	"fmt"
	"math"
	ipaket "nuryanto2121/cukur_in_barber/interface/b_paket"
	inotification "nuryanto2121/cukur_in_barber/interface/notification"
	"nuryanto2121/cukur_in_barber/models"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	repofunction "nuryanto2121/cukur_in_barber/repository/function"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type usePaket struct {
	repoPaket        ipaket.Repository
	repoNotification inotification.Repository
	contextTimeOut   time.Duration
}

func NewUserMPaket(a ipaket.Repository, b inotification.Repository, timeout time.Duration) ipaket.Usecase {
	return &usePaket{repoPaket: a, repoNotification: b, contextTimeOut: timeout}
}

func (u *usePaket) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result *models.Paket, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoPaket.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *usePaket) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var tUser = models.Paket{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		// queryparam.Search = fmt.Sprintf("paket_name iLIKE '%%%s%%' OR descs iLIKE '%%%s%%'", queryparam.Search, queryparam.Search)
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += " AND owner_id = " + Claims.UserID
	} else {
		queryparam.InitSearch = " owner_id = " + Claims.UserID
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
func (u *usePaket) Create(ctx context.Context, Claims util.Claims, data *models.DataPaket) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mPaket models.Paket
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mPaket)
	if err != nil {
		return err
	}
	mPaket.OwnerID, _ = strconv.Atoi(Claims.UserID)
	mPaket.UserEdit = Claims.UserID
	mPaket.UserInput = Claims.UserID

	err = u.repoPaket.Create(&mPaket)
	if err != nil {
		return err
	}
	return nil

}
func (u *usePaket) Update(ctx context.Context, Claims util.Claims, ID int, data *models.DataPaket) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// myMap := util.ConvertStructToMap(data)
	myMap := structs.Map(data)
	// myMap := data.(map[string]interface{})
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoPaket.Update(ID, myMap)
	if err != nil {
		return err
	}

	// if paket inactive
	if !data.IsActive {
		fn := &repofunction.FN{
			Claims:    Claims,
			RepoNotif: u.repoNotification,
		}

		go fn.SendNotifNonAktifPaket(ID)
	}
	return nil
}
func (u *usePaket) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoPaket.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
