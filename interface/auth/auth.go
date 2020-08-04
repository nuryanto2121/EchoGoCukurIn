package iauth

import (
	"context"
	"nuryanto2121/dynamic_rest_api_go/models"
)

type Repository interface {
	GetDataLogin(ctx context.Context, Account string) (models.DataLogin, error)
	ChangePassword(ctx context.Context, data interface{}) (err error)
	Register(ctx context.Context, dataUser models.SsUser) error
}
type Usecase interface {
	Login(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error)
	ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error)
	ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error)
	Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error)
	Verify(ctx context.Context, dataVeriry models.VerifyForm) (err error)
}
