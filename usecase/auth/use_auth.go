package useauth

import (
	"context"
	"errors"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	iusers "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"nuryanto2121/dynamic_rest_api_go/redisdb"
	useemailauth "nuryanto2121/dynamic_rest_api_go/usecase/email_auth"
	"time"
)

type useAuht struct {
	repoAuth       iusers.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserAuth(a iusers.Repository, b ifileupload.Repository, timeout time.Duration) iauth.Usecase {
	return &useAuht{repoAuth: a, repoFile: b, contextTimeOut: timeout}
}
func (u *useAuht) Login(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		DataOwner        = models.SsUser{}
		DataCapster      = models.LoginCapster{}
		isBarber    bool = true
		response         = map[string]interface{}{}
	)

	DataOwner, err = u.repoAuth.GetByAccount(dataLogin.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if DataOwner.UserType == "" && err == models.ErrNotFound {
		return nil, errors.New("Your Account not valid.")
	} else {
		if DataOwner.UserType == "capster" {
			DataCapster, err = u.repoAuth.GetByCapster(dataLogin.Account)
			if err != nil {
				return nil, errors.New("Your Account not valid.")
			}
			isBarber = false
		}
	}
	// if err != nil {
	// 	if err == models.ErrNotFound {
	// 		DataCapster, err = u.repoAuth.GetByCapster(dataLogin.Account)
	// 		if err != nil {
	// 			return nil, errors.New("Your Account not valid.")
	// 		}
	// 		isBarber = false
	// 	} else {
	// 		return nil, errors.New("Your Account not valid.")
	// 	}
	// }

	if isBarber {
		if !util.ComparePassword(DataOwner.Password, util.GetPassword(dataLogin.Password)) {
			return nil, errors.New("Your Password not valid.")
		}
		DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataOwner.FileID)

		token, err := util.GenerateToken(DataOwner.UserID, dataLogin.Account, DataOwner.UserType)
		if err != nil {
			return nil, err
		}

		redisdb.AddSession(token, DataOwner.UserID, 0)

		restUser := map[string]interface{}{
			"owner_id":   DataOwner.UserID,
			"owner_name": DataOwner.Name,
			"email":      DataOwner.Email,
			"telp":       DataOwner.Telp,
			"file_id":    DataOwner.FileID,
			"file_name":  DataFile.FileName,
			"file_path":  DataFile.FilePath,
		}
		response = map[string]interface{}{
			"token":      token,
			"data_owner": restUser,
			"user_type":  "barber",
		}

	} else {
		if !util.ComparePassword(DataCapster.Password, util.GetPassword(dataLogin.Password)) {
			return nil, errors.New("Your Password not valid.")
		}

		token, err := util.GenerateTokenCapster(DataCapster.CapsterID, DataCapster.OwnerID, DataCapster.BarberID)
		if err != nil {
			return nil, err
		}
		redisdb.AddSession(token, DataCapster.CapsterID, 0)

		restUser := map[string]interface{}{
			"owner_id":     DataCapster.OwnerID,
			"owner_name":   DataCapster.OwnerName,
			"barber_id":    DataCapster.BarberID,
			"barber_name":  DataCapster.BarberName,
			"capster_id":   DataCapster.CapsterID,
			"email":        DataCapster.Email,
			"telp":         DataCapster.Telp,
			"capster_name": DataCapster.CapsterName,
			"file_id":      DataCapster.FileID,
			"file_name":    DataCapster.FileName,
			"file_path":    DataCapster.FilePath,
		}
		response = map[string]interface{}{
			"token":        token,
			"data_capster": restUser,
			"user_type":    "capster",
		}
	}

	return response, nil
}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetByAccount(dataForgot.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
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

	DataUser, err := u.repoAuth.GetByAccount(dataReset.Account)
	if err != nil {
		return err
	}

	DataUser.Password, _ = util.Hash(dataReset.Passwd)
	// email, err := util.ParseEmailToken(dataReset.TokenEmail)
	// if err != nil {
	// 	email = dataReset.TokenEmail
	// }

	// dataUser, err := u.repoUser.GetByEmailSaUser(email)

	// dataUser.Password, _ = util.Hash(dataReset.Passwd)
	var data = map[string]interface{}{
		"password": DataUser.Password,
	}

	err = u.repoAuth.Update(DataUser.UserID, data)
	if err != nil {
		return err
	}

	return nil
}

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	CekData, err := u.repoAuth.GetByAccount(dataRegister.EmailAddr)

	if CekData.Email == dataRegister.EmailAddr {
		return output, errors.New("email sudah terdaftar.")
	}

	var User models.SsUser
	// GenCode := util.GenerateNumber(4)
	GenPassword := util.GenerateCode(4)

	User.Name = ""

	User.UserType = "owner"
	User.UserEdit = "cukur_in"
	User.UserInput = "cukur_in"
	User.Email = dataRegister.EmailAddr
	User.IsActive = true

	User.Password, _ = util.Hash(GenPassword)
	//check email or telp
	if !util.CheckEmail(dataRegister.EmailAddr) {
		return output, errors.New("email not valid")
	}

	err = u.repoAuth.Create(&User)
	if err != nil {
		return output, err
	}

	// send generate code
	mailService := &useemailauth.Register{
		Email:      User.Email,
		Name:       User.Name,
		PasswordCd: GenPassword,
	}

	err = mailService.SendRegister()
	if err != nil {
		u.repoAuth.Delete(User.UserID)
		return output, err
	}

	//store to redis
	err = redisdb.AddSession(dataRegister.EmailAddr, GenPassword, 2)
	if err != nil {
		u.repoAuth.Delete(User.UserID)
		return output, err
	}
	out := map[string]interface{}{
		"gen_password": GenPassword,
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
