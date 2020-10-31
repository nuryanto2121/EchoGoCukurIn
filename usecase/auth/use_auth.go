package useauth

import (
	"context"
	"errors"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	ifileupload "nuryanto2121/dynamic_rest_api_go/interface/fileupload"
	iusers "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"nuryanto2121/dynamic_rest_api_go/redisdb"
	useemailauth "nuryanto2121/dynamic_rest_api_go/usecase/email_auth"
	"strconv"
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
		expireToken      = setting.FileConfigSetting.JWTExpired
		canChange   bool = false
	)

	if dataLogin.Type == "owner" {
		DataOwner, err = u.repoAuth.GetByAccount(dataLogin.Account, false) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
		if DataOwner.UserType == "" && err == models.ErrNotFound {
			return nil, errors.New("Email anda belum terdaftar.")
		} else {
			if DataOwner.UserType == "capster" {
				DataCapster, err = u.repoAuth.GetByCapster(dataLogin.Account)
				if err != nil {
					return nil, errors.New("Email anda belum terdaftar.")
				}
				if !DataCapster.IsActive {
					return nil, errors.New("Account anda belum aktif. Silahkan hubungi pemilik Barber")
				}

				if DataCapster.BarberID == 0 {
					return nil, errors.New("Anda belum terhubung dengan Barber, Silahkan hubungi pemilik Barber")
				}

				if !DataCapster.BarberIsActive {
					return nil, errors.New("Saat ini barber anda sedang tidak aktif.")
				}
				isBarber = false
			} else {
				DataCapster, err = u.repoAuth.GetByCapster(dataLogin.Account)
				if DataCapster.Email != "" && DataCapster.Email == dataLogin.Account {
					canChange = true
				}

			}

			if !DataOwner.IsActive {
				return nil, errors.New("Account andan belum aktif. Silahkan hubungi pemilik Barber")
			}
		}
	} else {
		isBarber = false
		DataCapster, err = u.repoAuth.GetByCapster(dataLogin.Account)
		if err != nil {
			return nil, errors.New("Email anda belum terdaftar.")
		}
		if !DataCapster.IsActive {
			return nil, errors.New("Account anda belum aktif. Silahkan hubungi pemilik Barber")
		}

		if DataCapster.BarberID == 0 {
			return nil, errors.New("Anda belum terhubung dengan Barber, Silahkan hubungi pemilik Barber")
		}

		if !DataCapster.BarberIsActive {
			return nil, errors.New("Saat ini barber anda sedang tidak aktif.")
		}
	}

	if isBarber {
		if !util.ComparePassword(DataOwner.Password, util.GetPassword(dataLogin.Password)) {
			return nil, errors.New("Password yang anda masukkan salah. Silahkan coba lagi.")
		}
		DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataOwner.FileID)

		token, err := util.GenerateToken(DataOwner.UserID, dataLogin.Account, DataOwner.UserType)
		if err != nil {
			return nil, err
		}

		redisdb.AddSession(token, DataOwner.UserID, time.Duration(expireToken)*time.Hour)

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
			"can_change": canChange,
		}

	} else {
		if !util.ComparePassword(DataCapster.Password, util.GetPassword(dataLogin.Password)) {
			return nil, errors.New("Password yang anda masukkan salah. Silahkan coba lagi.")
		}

		token, err := util.GenerateTokenCapster(DataCapster.CapsterID, DataCapster.OwnerID, DataCapster.BarberID)
		if err != nil {
			return nil, err
		}
		redisdb.AddSession(token, DataCapster.CapsterID, time.Duration(expireToken)*time.Hour)

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
			"can_change":   canChange,
		}
	}

	return response, nil
}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (result string, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetByAccount(dataForgot.Account, true) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return "", errors.New("Your Account not valid.")
	}
	if DataUser.Name == "" {
		return "", errors.New("Your AccountF not valid.")
	}
	GenOTP := util.GenerateNumber(4)
	// send generate password
	mailservice := &useemailauth.Forgot{
		Email: DataUser.Email,
		Name:  DataUser.Name,
		OTP:   GenOTP,
	}

	// check data redis
	data := redisdb.GetSession(dataForgot.Account + "_Forgot")
	if data != "" {
		redisdb.TurncateList(dataForgot.Account + "_Forgot")
	}
	//store to redis
	err = redisdb.AddSession(dataForgot.Account+"_Forgot", GenOTP, 24*time.Hour)
	if err != nil {
		return "", err
	}
	// out := map[string]interface{}{
	// 	"gen_password": GenPassword,
	// }
	go mailservice.SendForgot()
	// err = mailservice.SendForgot()
	// if err != nil {
	// 	return "", err
	// }

	return GenOTP, nil
}

func (u *useAuht) ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return errors.New("Password dan Confirm Password tidak boleh sama.")
	}

	DataUser, err := u.repoAuth.GetByAccount(dataReset.Account, true)
	if err != nil {
		return err
	}

	if util.ComparePassword(DataUser.Password, util.GetPassword(dataReset.Passwd)) {
		return errors.New("Password baru tidak boleh sama dengan yang lama.")
	}

	DataUser.Password, _ = util.Hash(dataReset.Passwd)
	// email, err := util.ParseEmailToken(dataReset.TokenEmail)
	// if err != nil {
	// 	email = dataReset.TokenEmail
	// }

	// dataUser, err := u.repoUser.GetByEmailSaUser(email)

	// dataUser.Password, _ = util.Hash(dataReset.Passwd)
	// var data = map[string]interface{}{
	// 	"password": DataUser.Password,
	// }

	err = u.repoAuth.UpdatePasswordByEmail(dataReset.Account, DataUser.Password) //u.repoAuth.Update(DataUser.UserID, data)
	if err != nil {
		return err
	}

	return nil
}

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	CekData, err := u.repoAuth.GetByAccount(dataRegister.EmailAddr, true)

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
	User.JoinDate = time.Now()

	User.Password, _ = util.Hash(GenPassword)
	//check email or telp
	if !util.CheckEmail(dataRegister.EmailAddr) {
		return output, errors.New("email not valid")
	}

	err = u.repoAuth.Create(&User)
	if err != nil {
		return output, err
	}

	// User.UserInput = strconv.Itoa(User.UserID)
	// User.UserEdit = strconv.Itoa(User.UserID)
	mUser := map[string]interface{}{
		"user_input": strconv.Itoa(User.UserID),
		"user_edit":  strconv.Itoa(User.UserID),
	}
	err = u.repoAuth.Update(User.UserID, mUser)
	if err != nil {
		return output, err
	}

	// send generate code
	mailService := &useemailauth.Register{
		Email:      User.Email,
		Name:       User.Email,
		PasswordCd: GenPassword,
	}

	go mailService.SendRegister()
	// err = mailService.SendRegister()
	// if err != nil {
	// 	u.repoAuth.Delete(User.UserID)
	// 	return output, err
	// }

	//store to redis
	err = redisdb.AddSession(dataRegister.EmailAddr, GenPassword, 24*time.Hour)
	if err != nil {
		u.repoAuth.Delete(User.UserID)
		return output, err
	}
	out := map[string]interface{}{
		"gen_password": GenPassword,
	}
	return out, nil
}

func (u *useAuht) Verify(ctx context.Context, dataVerify models.VerifyForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	data := redisdb.GetSession(dataVerify.Account + "_Forgot")
	if data == "" {
		return errors.New("Please Resend Code")
	}

	if data != dataVerify.VerifyCode {
		return errors.New("Invalid Code.")
	}
	redisdb.TurncateList(dataVerify.Account + "_Forgot")

	return nil
}
