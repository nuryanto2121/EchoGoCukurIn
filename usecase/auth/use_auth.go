package useauth

import (
	"context"
	"errors"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"nuryanto2121/dynamic_rest_api_go/redisdb"
	useemailauth "nuryanto2121/dynamic_rest_api_go/usecase/email_auth"
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

	DataUser, err := u.repoAuth.GetDataLogin(ctx, dataLogin.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return nil, errors.New("Your Account not valid.")
	}

	if !util.ComparePassword(DataUser.Password, util.GetPassword(dataLogin.Password)) {
		return nil, errors.New("Your Password not valid.")
	}

	token, err := util.GenerateToken(DataUser.UserID, dataLogin.Account, DataUser.UserType)
	if err != nil {
		return nil, err
	}

	redisdb.AddSession(token, DataUser.UserID, 0)

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
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetDataLogin(ctx, dataForgot.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return errors.New("Your Account not valid.")
	}
	if DataUser.Name == "" {
		return errors.New("Your Account not valid.")
	}

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

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var User models.SysUser

	User.Name = dataRegister.Name
	if dataRegister.Passwd != dataRegister.ConfirmPasswd {
		return output, errors.New("Password and Confirm Password Not Valid")
	}
	User.Password, _ = util.Hash(dataRegister.Passwd)
	User.UserType = dataRegister.UserType
	User.UserEdit = dataRegister.Name
	User.UserInput = dataRegister.Name

	//check email or telp
	if util.CheckEmail(dataRegister.Account) {
		User.Email = dataRegister.Account
	} else {
		User.Telp = dataRegister.Account
	}
	err = u.repoAuth.Register(ctx, User)
	if err != nil {
		return output, err
	}

	GenCode := util.GenerateNumber(4)

	// send generate code
	mailService := &useemailauth.Register{
		Email:      User.Email,
		Name:       User.Name,
		GenerateNo: GenCode,
	}

	err = mailService.SendRegister()
	if err != nil {
		return output, err
	}

	//store to redis
	err = redisdb.AddSession(dataRegister.Account, GenCode, 2)
	if err != nil {
		return output, err
	}
	out := map[string]interface{}{
		"gen_code": GenCode,
	}
	return out, nil
}

func (u *useAuht) Verify(ctx context.Context, dataVeriry models.VerifyForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	data := redisdb.GetSession(dataVeriry.Account)
	if data == "" {
		return errors.New("Please Resend Code")
	}

	if data != dataVeriry.VerifyCode {
		return errors.New("Invalid Code.")
	}

	return nil
}
