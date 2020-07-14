package useauth

import (
	"context"
	"errors"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"nuryanto2121/dynamic_rest_api_go/redisdb"
	"time"
)

type useAuht struct {
	repoAuth       iauth.Repository
	contextTimeOut time.Duration
}

func NewUserAuth(a iauth.Repository, timeout time.Duration) iauth.Usecase {
	return &useAuht{repoAuth: a, contextTimeOut: timeout}
}
func (u *useAuht) Login(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetDataLogin(ctx, dataLogin.UserName) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return nil, errors.New("Your Account not valid.")
	}

	if !util.ComparePassword(DataUser.Password, util.GetPassword(dataLogin.Password)) {
		return nil, errors.New("Your Password not valid.")
	}

	token, err := util.GenerateToken(DataUser.UserID, dataLogin.UserName, DataUser.UserType)
	if err != nil {
		return nil, err
	}

	redisdb.AddSession(token, DataUser.UserID)

	restUser := map[string]interface{}{
		"id":        DataUser.UserID,
		"email":     DataUser.Email,
		"telp":      DataUser.Telp,
		"user_name": DataUser.Name,
		"user_type": DataUser.UserType,
		"file_id":   DataUser.FileID.Int64,
		"file_name": DataUser.FileName.String,
		"file_path": DataUser.FilePath.String,
	}
	response := map[string]interface{}{
		"token":     token,
		"data_user": restUser,
	}

	return response, nil
}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error) {
	return nil
}

func (u *useAuht) ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return errors.New("Password and Confirm Password not same.")
	}

	DataUser, err := u.repoAuth.GetDataLogin(ctx, dataReset.Account)
	if err != nil {
		return err
	}
	DataChange := make(map[string]interface{}, 0)
	DataChange["user_id"] = DataUser.UserID
	DataChange["pwd"], _ = util.Hash(dataReset.Passwd)
	// email, err := util.ParseEmailToken(dataReset.TokenEmail)
	// if err != nil {
	// 	email = dataReset.TokenEmail
	// }

	// dataUser, err := u.repoUser.GetByEmailSaUser(email)

	// dataUser.Password, _ = util.Hash(dataReset.Passwd)

	err = u.repoAuth.ChangePassword(ctx, DataChange) //u.repoUser.Update(dataUser.ID, &dataUser)
	if err != nil {
		return err
	}

	return nil
}
