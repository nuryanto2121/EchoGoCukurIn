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

		queryparam.Search = fmt.Sprintf("lower(name) LIKE '%%%s%%' ", queryparam.Search)
	}
	queryparam.InitSearch = "user_type='capster'"

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

	mUser.UserEdit = Claims.UserName
	mUser.UserInput = Claims.UserName
	err = u.repoUser.Create(&mUser)
	if err != nil {
		return err
	}

	for _, dataCollection := range data.TopCollection {
		var capsterCollection = models.CapsterCollection{}
		capsterCollection.CapsterID = mUser.UserID
		capsterCollection.FileID = dataCollection.FileID
		capsterCollection.UserInput = Claims.UserName
		capsterCollection.UserEdit = Claims.UserName
		err = u.repoCapster.Create(&capsterCollection)
		if err != nil {
			return err
		}
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
	dataUser.UserEdit = Claims.UserName

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
		capsterCollection.UserInput = Claims.UserName
		capsterCollection.UserEdit = Claims.UserName

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
