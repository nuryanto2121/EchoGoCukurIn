package contauth

import (
	"context"
	"fmt"
	"net/http"
	iauth "nuryanto2121/cukur_in_barber/interface/auth"
	midd "nuryanto2121/cukur_in_barber/middleware"
	"nuryanto2121/cukur_in_barber/models"
	app "nuryanto2121/cukur_in_barber/pkg"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	tool "nuryanto2121/cukur_in_barber/pkg/tools"
	util "nuryanto2121/cukur_in_barber/pkg/utils"

	_ "nuryanto2121/cukur_in_barber/docs"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

type ContAuth struct {
	useAuth iauth.Usecase
}

func NewContAuth(e *echo.Echo, useAuth iauth.Usecase) {
	cont := &ContAuth{
		useAuth: useAuth,
		// useSaClient:     useSaClient,
		// useSaUser:       useSaUser,
		// useSaFileUpload: useSaFileUpload,
	}

	// e.POST("/barber/auth/register", cont.Register)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	r := e.Group("/barber-service/auth")
	r.Use(midd.Versioning)
	r.POST("/login", cont.Login)
	r.POST("/forgot", cont.ForgotPassword)
	r.POST("/change_password", cont.ChangePassword)
	r.POST("/verify", cont.Verify)
	r.POST("/register", cont.Register)

	L := e.Group("/barber-service/auth/logout")
	L.Use(midd.JWT)
	L.Use(midd.Versioning)
	L.POST("", cont.Logout)
}

// Logout :
// @Summary logout
// @Security ApiKeyAuth
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/logout [post]
func (u *ContAuth) Logout(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE = tool.Res{R: e} // wajib
	)

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	Token := e.Request().Header.Get("Authorization")
	err = u.useAuth.Logout(ctx, claims, Token)
	if err != nil {
		return appE.ResponseError(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", nil)
}

// Login :
// @Summary Login
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.LoginForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/login [post]
func (u *ContAuth) Login(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		// client sa_models.SaClient

		form = models.LoginForm{}
		// dataFiles = sa_models.SaFileOutput{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	out, err := u.useAuth.Login(ctx, &form)
	if err != nil {
		// return appE.Response(out)
		// return appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return appE.ResponseError(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", out)
}

// ChangePassword :
// @Summary Change Password
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.ResetPasswd true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/change_password [post]
func (u *ContAuth) ChangePassword(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		// client sa_models.SaClient

		form = models.ResetPasswd{}
	)
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}
	err := u.useAuth.ResetPassword(ctx, &form)
	if err != nil {
		return appE.ResponseError(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", "Please Login")
}

// Register :
// @Summary Register Owner
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.RegisterForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/register [post]
func (u *ContAuth) Register(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		// client sa_models.SaClient

		form = models.RegisterForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	data, err := u.useAuth.Register(ctx, form)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusOK, "Ok", data)
}

// Verify :
// @Summary Verify
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.VerifyForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/verify [post]
func (u *ContAuth) Verify(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		// client sa_models.SaClient

		form = models.VerifyForm{}
	)

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	err := u.useAuth.Verify(ctx, form)
	if err != nil {
		return appE.ResponseError(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusOK, "Ok", nil)
}

// ForgotPassword :
// @Summary Forgot Password
// @Tags Auth
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.ForgotForm true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/auth/forgot [post]
func (u *ContAuth) ForgotPassword(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		// client sa_models.SaClient

		form = models.ForgotForm{}
	)
	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	OTP, err := u.useAuth.ForgotPassword(ctx, &form)
	if err != nil {
		return appE.ResponseError(http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
	}
	result := map[string]interface{}{
		"otp": OTP,
	}
	return appE.Response(http.StatusOK, "Check Your Email", result)

}
