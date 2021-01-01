package usenotification

import (
	"context"
	"fmt"
	"math"
	inotification "nuryanto2121/cukur_in_barber/interface/notification"
	"nuryanto2121/cukur_in_barber/models"
	fcmgetway "nuryanto2121/cukur_in_barber/pkg/fcm"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useNotification struct {
	repoNotification inotification.Repository
	contextTimeOut   time.Duration
}

func NewUseNotification(a inotification.Repository, timeout time.Duration) inotification.Usecase {
	return &useNotification{repoNotification: a, contextTimeOut: timeout}
}

func (u *useNotification) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result *models.Notification, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoNotification.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useNotification) GetCountNotif(ctx context.Context, Claims util.Claims) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		queryparam models.ParamList
	)

	queryparam.InitSearch = fmt.Sprintf(`notification_status = 'N' AND user_id = %s`, Claims.UserID)

	Total, err := u.repoNotification.Count(queryparam)
	if err != nil {
		return result, err
	}

	response := map[string]interface{}{
		"total_notif": Total,
	}
	return response, nil
}
func (u *useNotification) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {

	}
	result.Data, err = u.repoNotification.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoNotification.Count(queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useNotification) Create(ctx context.Context, Claims util.Claims, Token string, data *models.AddNotification) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mNotification = models.Notification{}
		queryParam    = models.ParamList{}
		TokenFCM      []string
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mNotification.AddNotification)
	if err != nil {
		return err
	}

	mNotification.UserEdit = Claims.UserID
	mNotification.UserInput = Claims.UserID

	err = u.repoNotification.Create(&mNotification)
	if err != nil {
		return err
	}
	// send notif to user
	TokenFCM = []string{Token}
	queryParam.InitSearch = fmt.Sprintf("user_id = %s and notification_status = 'N' ", strconv.Itoa(data.UserId))
	cntNotif, err := u.repoNotification.Count(queryParam)
	if err != nil {
		return err
	}

	fcm := &fcmgetway.SendFCM{
		Title:       mNotification.Title,
		Body:        mNotification.Descs,
		JumlahNotif: cntNotif,
		DeviceToken: TokenFCM,
	}

	go fcm.SendPushNotification()
	// err = fcm.SendPushNotification(ctx)
	// if err != nil {
	// 	return err
	// }
	//end send notif

	return nil

}

func (u *useNotification) Update(ctx context.Context, Claims util.Claims, ID int, data *models.StatusNotification) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoNotification.Update(ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useNotification) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoNotification.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
