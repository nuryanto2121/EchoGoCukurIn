package usecapster

import (
	"context"
	"errors"
	"fmt"
	"math"
	ibarbercapster "nuryanto2121/cukur_in_barber/interface/b_barber_capster"
	icapster "nuryanto2121/cukur_in_barber/interface/b_capster"
	ifileupload "nuryanto2121/cukur_in_barber/interface/fileupload"
	inotification "nuryanto2121/cukur_in_barber/interface/notification"
	iuser "nuryanto2121/cukur_in_barber/interface/user"
	"nuryanto2121/cukur_in_barber/models"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	util "nuryanto2121/cukur_in_barber/pkg/utils"

	repofunction "nuryanto2121/cukur_in_barber/repository/function"

	// useemailcapster "nuryanto2121/cukur_in_barber/usecase/email_capster"
	useemailauth "nuryanto2121/cukur_in_barber/usecase/email_auth"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useCapster struct {
	repoCapster       icapster.Repository
	repoUser          iuser.Repository
	repoBarberCapster ibarbercapster.Repository
	repoFile          ifileupload.Repository
	repoNotification  inotification.Repository
	contextTimeOut    time.Duration
}

func NewUserMCapster(a icapster.Repository, b iuser.Repository, c ibarbercapster.Repository, d ifileupload.Repository, e inotification.Repository, timeout time.Duration) icapster.Usecase {
	return &useCapster{
		repoCapster:       a,
		repoUser:          b,
		repoBarberCapster: c,
		repoFile:          d,
		repoNotification:  e,
		contextTimeOut:    timeout,
	}
}

func (u *useCapster) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataCapster, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return result, err
	}

	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, dataCapster.FileID)
	if err != nil {
		if err != models.ErrNotFound {
			return result, err
		}
	}

	dataCollection, err := u.repoCapster.GetListFileCapter(ID)
	if err != nil {
		if err != models.ErrNotFound {
			return result, err
		}

	}
	response := map[string]interface{}{
		"capster_id":     dataCapster.UserID,
		"email":          dataCapster.Email,
		"name":           dataCapster.Name,
		"join_date":      dataCapster.JoinDate,
		"user_type":      dataCapster.UserType,
		"is_active":      dataCapster.IsActive,
		"file_id":        dataCapster.FileID,
		"file_name":      dataFile.FileName,
		"file_path":      dataFile.FilePath,
		"top_collection": dataCollection,
	}

	return response, nil
}
func (u *useCapster) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var tUser = models.Capster{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {

		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		// queryparam.InitSearch = strings.ReplaceAll(queryparam.InitSearch, "capster_id", "ss_user.user_id")
		queryparam.InitSearch += fmt.Sprintf(" AND user_type='capster' and user_input = '%s'", Claims.UserID)
	} else {
		queryparam.InitSearch = fmt.Sprintf("user_type='capster' and user_input = '%s'", Claims.UserID)
	}
	result.Data, err = u.repoCapster.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoCapster.Count(queryparam)
	if err != nil {
		result.Data = nil
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useCapster) Create(ctx context.Context, Claims util.Claims, data *models.Capster) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		mUser               = models.SsUser{}
		logger              = logging.Logger{}
		OwnerCapster bool   = true
		GenPassword  string = ""
	)
	fn := &repofunction.FN{
		Claims: Claims,
	}
	//insert user
	err = mapstructure.Decode(data, &mUser)
	if err != nil {
		return err
	}

	dataOwner, err := fn.GetOwnerData()
	if err != nil {
		return err
	}

	if dataOwner.Email != data.Email {
		OwnerCapster = false
		// if Claims.UserType != "owner" {
		dataCapster, err := u.repoUser.GetByAccount(data.Email, false)
		if dataCapster.Email != "" {
			return errors.New("Email Capster sudah terdaftar.")
		}
		if err != nil {
			logger.Error(err)
		}

		// fmt.Printf("%v", err)
		// }
		GenPassword = util.GenerateCode(4)
		mUser.Password, _ = util.Hash(GenPassword)
	} else {
		mUser.Password = dataOwner.Password
	}

	// gen Password

	// mUser.UserName, err = u.repoUser.GenUserCapster()
	// if err != nil {
	// 	return err
	// }
	mUser.JoinDate = data.JoinDate

	mUser.UserEdit = Claims.UserID
	mUser.UserInput = Claims.UserID
	err = u.repoUser.Create(&mUser)
	if err != nil {
		return err
	}

	for _, dataCollection := range data.TopCollection {
		if dataCollection.FileID > 0 {
			var capsterCollection = models.CapsterCollection{}
			capsterCollection.CapsterID = mUser.UserID
			capsterCollection.FileID = dataCollection.FileID
			capsterCollection.UserInput = Claims.UserID
			capsterCollection.UserEdit = Claims.UserID
			err = u.repoCapster.Create(&capsterCollection)
			if err != nil {
				return err
			}
		}

	}

	if !OwnerCapster {
		// send Password
		mailService := &useemailauth.Register{
			Email:      mUser.Email,
			Name:       mUser.Name,
			PasswordCd: GenPassword,
		}

		go mailService.SendRegister()

	}

	return nil

}
func (u *useCapster) Update(ctx context.Context, Claims util.Claims, ID int, data *models.Capster) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		dataUser = &models.CapsterUpdate{}
		logger   = logging.Logger{}
	)

	fn := &repofunction.FN{
		Claims:    Claims,
		RepoNotif: u.repoNotification,
	}

	dataOwner, err := fn.GetOwnerData()
	if err != nil {
		return err
	}
	err = mapstructure.Decode(data, &dataUser)
	if err != nil {
		return err
	}
	if dataOwner.Email != data.Email {
		dataCapster, err := u.repoUser.GetByAccount(dataUser.Email, false)
		if dataCapster.UserID != ID {
			return errors.New("Email Capster sudah terdaftar.")
		}
		if err != nil {
			logger.Error(err)
		}
	}

	dataUser.JoinDate = data.JoinDate

	//if status not active then delete relasi in barber_capster
	if !dataUser.IsActive {
		cnt := fn.GetCountTrxCapster(ID)
		if cnt > 0 {
			return errors.New("Mohon maaf anda tidak dapat non-aktifkan Kapster, sedang ada transaksi yang berlangsung")
		}

		go fn.SendNotifNonAktifCapster(ID)

	}

	datas := structs.Map(dataUser)
	datas["user_edit"] = Claims.UserID
	fmt.Println(datas)

	err = u.repoUser.Update(ID, datas)
	if err != nil {
		return err
	}

	err = u.repoCapster.Delete(ID)
	if err != nil {
		return err
	}
	for _, dataCollection := range data.TopCollection {
		var capsterCollection = models.CapsterCollection{}
		capsterCollection.CapsterID = ID
		capsterCollection.FileID = dataCollection.FileID
		capsterCollection.UserInput = Claims.UserID
		capsterCollection.UserEdit = Claims.UserID

		err = u.repoCapster.Create(&capsterCollection)
		if err != nil {
			return err
		}
	}

	return nil
}
func (u *useCapster) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoCapster.Delete(ID)
	if err != nil {
		return err
	}

	err = u.repoUser.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}

// func (u *useCapster) GenUserCapster() (string, error) {
// 	var (
// 		currentTime     = time.Now()
// 		year            = currentTime.Year()
// 		month       int = int(currentTime.Month())
// 		day             = currentTime.Day()
// 	)

// 	result := u.repoUser.GetMaxUserCapster()
// 	sYear := strconv.Itoa(year)[2:]
// 	var sMonth string = strconv.Itoa(month)
// 	if len(sMonth) == 1 {
// 		sMonth = fmt.Sprintf("0%s", sMonth)
// 	}
// 	var sDay string = strconv.Itoa(day)
// 	if len(sDay) == 1 {
// 		sDay = fmt.Sprintf("0%s", sDay)
// 	}
// 	seqNo := "0001"
// 	if result == "" {
// 		result = fmt.Sprintf("CP%s%s%s%v", sYear, sMonth, sDay, seqNo)
// 	} else {
// 		seqNo = fmt.Sprintf("1%s", result[9:])
// 		iSeqno, err := strconv.Atoi(seqNo)
// 		if err != nil {
// 			return "", err
// 		}
// 		iSeqno += 1
// 		seqNo = strconv.Itoa(iSeqno)[1:]
// 		result = fmt.Sprintf("CP%s%s%s%v", sYear, sMonth, sDay, seqNo)
// 	}
// 	return result, nil

// }
