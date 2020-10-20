package contcorder

import (
	"context"
	"fmt"
	"net/http"
	ibarbers "nuryanto2121/dynamic_rest_api_go/interface/b_barber"
	icorder "nuryanto2121/dynamic_rest_api_go/interface/c_order_h"
	midd "nuryanto2121/dynamic_rest_api_go/middleware"
	"nuryanto2121/dynamic_rest_api_go/models"
	app "nuryanto2121/dynamic_rest_api_go/pkg"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContOrder struct {
	useOrder  icorder.Usecase
	useBarber ibarbers.Usecase
}

func NewContOrder(e *echo.Echo, a icorder.Usecase, b ibarbers.Usecase) {
	controller := &ContOrder{
		useOrder:  a,
		useBarber: b,
	}

	r := e.Group("/barber/order")
	r.Use(midd.JWT)
	r.GET("/:id", controller.GetDataBy)
	r.GET("", controller.GetList)
	r.POST("", controller.Create)
	r.PUT("/:id", controller.Update)
	r.DELETE("/:id", controller.Delete)
}

// GetDataByID :
// @Summary GetById
// @Security ApiKeyAuth
// @Tags Order
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/order/{id} [get]
func (u *ContOrder) GetDataBy(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	data, err := u.useOrder.GetDataBy(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}

// GetList :
// @Summary GetList Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param barber_id query int true "BarberId"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Success 200 {object} models.ResponseModelList
// @Router /barber/order [get]
func (u *ContOrder) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	idUser := e.QueryParam("barber_id")
	fmt.Printf(idUser)
	var (
		logger = logging.Logger{}
		appE   = tool.Res{R: e} // wajib
		//valid      validation.Validation // wajib
		paramquery   = models.ParamListOrder{} // ini untuk list
		responseList = models.ResponseModelList{}
		err          error
		dataBarber   interface{}
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseErrorList(http.StatusBadRequest, errMsg, responseList)
	}
	if paramquery.BarberId == 0 {
		paramquery.BarberId, _ = strconv.Atoi(idUser)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseErrorList(http.StatusBadRequest, fmt.Sprintf("%v", err), responseList)
	}

	if paramquery.BarberId > 0 {
		dataBarber, err = u.useBarber.GetDataFirst(ctx, claims, paramquery.BarberId)
	}

	responseList, err = u.useOrder.GetList(ctx, claims, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	totalPrice, err := u.useOrder.GetSumPrice(ctx, claims, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}
	result := map[string]interface{}{
		"data_barber": dataBarber,
		"data_list":   responseList,
		"total_price": totalPrice,
	}

	// return e.JSON(http.StatusOK, ListOrderPost)
	return appE.Response(http.StatusOK, "", result)
}

// CreateSaOrder :
// @Summary Add Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param req body models.OrderPost true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/order [post]
func (u *ContOrder) Create(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		form   models.OrderPost
	)

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

	err = u.useOrder.Create(ctx, claims, &form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", nil)
}

// UpdateSaOrder :
// @Summary Rubah Profile
// @Security ApiKeyAuth
// @Tags Order
// @Produce json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Param req body models.OrderPost true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/order/{id} [put]
func (u *ContOrder) Update(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		logger = logging.Logger{} // wajib
		appE   = tool.Res{R: e}   // wajib
		err    error
		// valid  validation.Validation                 // wajib
		id   = e.Param("id") //kalo bukan int => 0
		form = models.OrderPost{}
		// form    models.OrderPost
	)

	OrderID, _ := strconv.Atoi(id)
	// logger.Info(id)

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

	// form.UpdatedBy = claims.OrderName
	err = u.useOrder.Update(ctx, claims, OrderID, form)
	if err != nil {
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusCreated, "Ok", nil)
}

// DeleteSaOrder :
// @Summary Delete Order
// @Security ApiKeyAuth
// @Tags Order
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param id path string true "ID"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/order/{id} [delete]
func (u *ContOrder) Delete(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		id   = e.Param("id")  //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	ID, err := strconv.Atoi(id)
	if err != nil {
		return appE.Response(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	err = u.useOrder.Delete(ctx, claims, ID)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", nil)
}
