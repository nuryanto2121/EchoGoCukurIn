package contberandabarber

import (
	"context"
	"fmt"
	"net/http"
	iberandabarber "nuryanto2121/dynamic_rest_api_go/interface/beranda_barber"
	midd "nuryanto2121/dynamic_rest_api_go/middleware"
	"nuryanto2121/dynamic_rest_api_go/models"
	app "nuryanto2121/dynamic_rest_api_go/pkg"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"

	"github.com/labstack/echo/v4"
)

type ContBerandaBarber struct {
	useBeranda iberandabarber.Usecase
}

func NewContBeranda(e *echo.Echo, a iberandabarber.Usecase) {
	controller := &ContBerandaBarber{
		useBeranda: a,
	}
	r := e.Group("/barber/beranda")
	r.Use(midd.JWT)
	// r.GET("/status_order", controller.GetStatusOrder)
	r.GET("", controller.GetList)
}

/*
// GetStatusOrder :
// @Summary GetStatusOrder
// @Security ApiKeyAuth
// @Tags Barber Beranda
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/beranda/status_order [get]
func (u *ContBerandaBarber) GetStatusOrder(e echo.Context) error {
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

	data, err := u.useBeranda.GetStatusOrder(ctx, claims)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusOK, "Ok", data)
}
*/

// GetList :
// @Summary GetList Barber Beranda
// @Security ApiKeyAuth
// @Tags Barber Beranda
// @Produce  json
// @Param OS header string true "OS Device"
// @Param Version header string true "OS Device"
// @Param page query int true "Page"
// @Param perpage query int true "PerPage"
// @Param search query string false "Search"
// @Param initsearch query string false "InitSearch"
// @Param sortfield query string false "SortField"
// @Param paramview query string false "ParamView"
// @Success 200 {object} models.ResponseModelList
// @Router /barber/beranda [get]
func (u *ContBerandaBarber) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		//valid      validation.Validation // wajib
		paramquery = models.ParamDynamicList{} // ini untuk list
		// responseList = models.ResponseModelList{}
		// err          error
	)

	httpCode, errMsg := app.BindAndValid(e, &paramquery)
	// logger.Info(util.Stringify(paramquery))
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	// if !claims.IsAdmin {
	// 	paramquery.InitSearch = " id_created = " + strconv.Itoa(claims.BarberID)
	// }
	data, err := u.useBeranda.GetStatusOrder(ctx, claims, paramquery)
	if err != nil {
		return appE.Response(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	responseList, err := u.useBeranda.GetListOrder(ctx, claims, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseError(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	result := map[string]interface{}{
		"status_order": data,
		"data_list":    responseList,
	}

	// return e.JSON(http.StatusOK, ListBarbersPost)
	return appE.Response(http.StatusOK, "", result)
}
