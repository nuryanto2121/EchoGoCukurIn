package contnotification

import (
	"context"
	"fmt"
	"net/http"
	inotification "nuryanto2121/cukur_in_barber/interface/notification"
	"nuryanto2121/cukur_in_barber/models"
	app "nuryanto2121/cukur_in_barber/pkg"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	tool "nuryanto2121/cukur_in_barber/pkg/tools"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type contNotification struct {
	useNotification inotification.Usecase
}

func NewContNotification(e *echo.Echo, a inotification.Usecase) {
	controller := &contNotification{
		useNotification: a,
	}

	r := e.Group("/barber/notification")
	// r.Use(midd.JWT)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.GET("/beranda", controller.Beranda)
	// r.POST("", controller.Create)
	r.PUT("/:id", controller.Update)
	// r.DELETE("/:id", controller.Delete)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Notification
// @Produce  json
// @Param app_id header string true "AppID Device"
// @Param version header string true "Version Device"
// @Param language header string true "Language"
// @Param id path string true "ID"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/notification/{id} [get]
func (u *contNotification) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{}
		appE   = tool.Res{R: e} // wajib
		id     = e.Param("id")  //kalo bukan int => 0
	)
	ID, err := strconv.Atoi(id)
	logger.Info(ID)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	data, err := u.useNotification.GetDataBy(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// Beranda :
// @Summary Jumlah Notif yg belum dibuka
// @Security ApiKeyAuth
// @Tags Notification
// @Produce  json
// @Param app_id header string true "AppID Device"
// @Param version header string true "Version Device"
// @Param language header string true "Language"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/notification/beranda [get]
func (u *contNotification) Beranda(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
	)

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	data, err := u.useNotification.GetCountNotif(ctx, claims)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Notification
// @Security ApiKeyAuth
// @Tags Notification
// @Produce  json
// @Param app_id header string true "AppID Device"
// @Param version header string true "Version Device"
// @Param language header string true "Language"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /barber/notification [get]
func (u *contNotification) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE         = tool.Res{R: e}     // wajib
		paramquery   = models.ParamList{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	// logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	}

	responseList, err = u.useNotification.GetList(ctx, claims, paramquery)
	if err != nil {
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	return appE.ResponseList(http.StatusOK, "", responseList)
}

// UpdateNotification :
// @Summary Rubah Status Notification
// @Security ApiKeyAuth
// @Tags Notification
// @Produce json
// @Param app_id header string true "AppID Device"
// @Param version header string true "Version Device"
// @Param language header string true "Language"
// @Param id path string true "ID"
// @Param req body models.StatusNotification true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/notification/{id} [put]
func (u *contNotification) Update(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		err    error

		id   = e.Param("id") //kalo bukan int => 0
		form = models.StatusNotification{}
	)

	SchoolID, err := strconv.Atoi(id)
	logger.Info(SchoolID)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	// validasi and bind to struct
	httpCode, errMsg := app.BindAndValid(e, &form)
	logger.Info(util.Stringify(form))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	// form.UpdatedBy = claims.NotificationName
	err = u.useNotification.Update(ctx, claims, SchoolID, &form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusCreated, "Ok", nil)
}
