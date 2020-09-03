package usesysuser

import (
	"context"
	"math"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	iusers "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	querywhere "nuryanto2121/dynamic_rest_api_go/pkg/query"
	"reflect"
	"time"
)

type useSysUser struct {
	repoUser       iusers.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserSysUser(a iusers.Repository, b ifileupload.Repository, timeout time.Duration) iusers.Usecase {
	return &useSysUser{repoUser: a, repoFile: b, contextTimeOut: timeout}
}

func (u *useSysUser) GetByEmailSaUser(ctx context.Context, email string) (result models.SsUser, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	a := models.SsUser{}
	result, err = u.repoUser.GetByAccount(email)
	if err != nil {
		return a, err
	}
	return result, nil
}

func (u *useSysUser) GetDataBy(ctx context.Context, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataOwner, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataOwner.FileID)
	if err != nil {
		return result, err
	}
	response := map[string]interface{}{
		"owner_id":   DataOwner.UserID,
		"owner_name": DataOwner.Name,
		"email":      DataOwner.Email,
		"telp":       DataOwner.Telp,
		"file_id":    DataOwner.FileID,
		"file_name":  DataFile.FileName,
		"file_path":  DataFile.FilePath,
	}
	return response, nil
}
func (u *useSysUser) GetList(ctx context.Context, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var tUser = models.SsUser{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		value := reflect.ValueOf(tUser)
		types := reflect.TypeOf(&tUser)
		queryparam.Search = querywhere.GetWhereLikeStruct(value, types, queryparam.Search, "")
	}
	result.Data, err = u.repoUser.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUser.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useSysUser) Create(ctx context.Context, data *models.SsUser) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Create(data)
	if err != nil {
		return err
	}
	return nil

}
func (u *useSysUser) Update(ctx context.Context, ID int, data models.UpdateUser) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataOwner, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return err
	}
	DataOwner.FileID = data.FileID
	DataOwner.Name = data.Name
	DataOwner.Telp = data.Telp
	DataOwner.Email = data.Email

	err = u.repoUser.Update(ID, DataOwner)
	if err != nil {
		return err
	}
	return nil
}
func (u *useSysUser) Delete(ctx context.Context, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
