package usecapster

import (
	"context"
	"fmt"
	"math"
	icapster "nuryanto2121/dynamic_rest_api_go/interface/b_capster"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	iuser "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	useemailcapster "nuryanto2121/dynamic_rest_api_go/usecase/email_capster"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useCapster struct {
	repoCapster    icapster.Repository
	repoUser       iuser.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserMCapster(a icapster.Repository, b iuser.Repository, c ifileupload.Repository, timeout time.Duration) icapster.Usecase {
	return &useCapster{
		repoCapster:    a,
		repoUser:       b,
		repoFile:       c,
		contextTimeOut: timeout,
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
		return result, err
	}

	dataCollection, err := u.repoCapster.GetListFileCapter(ID)
	if err != nil {
		return result, err
	}
	response := map[string]interface{}{
		"capster_id":     dataCapster.UserID,
		"name":           dataCapster.Name,
		"join_date":      dataCapster.JoinDate,
		"user_type":      dataCapster.UserType,
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

		queryparam.Search = fmt.Sprintf("lower(ss_user.name) LIKE '%%%s%%' ", queryparam.Search)
	}
	queryparam.InitSearch = fmt.Sprintf("ss_user.user_type='capster' and ss_user.user_input = '%s'", Claims.UserID)

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
		mUser = models.SsUser{}
	)
	//insert user
	err = mapstructure.Decode(data, &mUser)
	if err != nil {
		return err
	}
	// gen Password
	GenPassword := util.GenerateCode(4)
	// mUser.UserName, err = u.repoUser.GenUserCapster()
	// if err != nil {
	// 	return err
	// }
	mUser.Password, _ = util.Hash(GenPassword)
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

	// send Password
	mailService := &useemailcapster.RegisterCapster{
		Email:    mUser.Email,
		Name:     mUser.Name,
		Password: GenPassword,
	}

	err = mailService.SendRegisterCapster()
	if err != nil {
		return err
	}

	return nil

}
func (u *useCapster) Update(ctx context.Context, Claims util.Claims, ID int, data *models.Capster) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	// var dataUser = &models.SsUser{}

	dataUser, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return err
	}
	dataUser.Name = data.Name
	dataUser.IsActive = data.IsActive
	dataUser.FileID = data.FileID
	dataUser.UserEdit = Claims.UserID

	err = u.repoUser.Update(ID, dataUser)
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
