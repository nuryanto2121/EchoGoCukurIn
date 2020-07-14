package contdynamic

import (
	"context"
	"fmt"
	"net/http"
	idynamic "nuryanto2121/dynamic_rest_api_go/interface/dynamic"
	midd "nuryanto2121/dynamic_rest_api_go/middleware"
	"nuryanto2121/dynamic_rest_api_go/models"
	app "nuryanto2121/dynamic_rest_api_go/pkg"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"

	"github.com/labstack/echo/v4"
)

type ContDynamic struct {
	useOption idynamic.Usecase
}

func NewContDynamic(e *echo.Echo, a idynamic.Usecase) {
	controller := &ContDynamic{
		useOption: a,
	}

	r := e.Group("/barber/dynamic")
	r.Use(midd.JWT)
	r.POST("", controller.Save)
	r.PUT("", controller.Update)
	r.GET("/:id", controller.GetById)
	r.POST("/list", controller.GetList)
	r.DELETE("/:id", controller.Delete)
}

// Save :
// @Summary Add Dynamic
// @Security ApiKeyAuth
// @Tags Dynamic
// @Produce json
// @Param req body interface{} true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/dynamic [post]
func (c *ContDynamic) Save(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		// logger     = logging.Logger{} // wajib
		appE     = tool.Res{R: e} // wajib
		json_map = make(map[string]interface{})
	)

	err := e.Bind(&json_map)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("invalid request param: %v", err), nil)
	}

	if _, ok := json_map["option_url"]; !ok {
		return appE.ResponseError(http.StatusBadRequest, "invalid request param: option_url", nil)
	}

	if _, ok := json_map["line_no"]; !ok {
		return appE.ResponseError(http.StatusBadRequest, "invalid request param: line_no", nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	json_map["method"] = "Save"
	data, err := c.useOption.Execute(ctx, claims, json_map)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", data)
}

// Update :
// @Summary Edit Dynamic
// @Security ApiKeyAuth
// @Tags Dynamic
// @Produce json
// @Param req body interface{} true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/dynamic [put]
func (c *ContDynamic) Update(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		// logger     = logging.Logger{} // wajib
		appE     = tool.Res{R: e} // wajib
		json_map = make(map[string]interface{})
	)

	err := e.Bind(&json_map)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("invalid request param: %v", err), nil)
	}

	if _, ok := json_map["option_url"]; !ok {
		return appE.ResponseError(http.StatusBadRequest, "invalid request param: option_url", nil)
	}

	if _, ok := json_map["line_no"]; !ok {
		return appE.ResponseError(http.StatusBadRequest, "invalid request param: line_no", nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	json_map["method"] = "Update"
	data, err := c.useOption.Execute(ctx, claims, json_map)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusCreated, "Ok", data)
}

// GetById :
// @Summary GetById Dynamic
// @Security ApiKeyAuth
// @Tags Dynamic
// @Produce  json
// @Param id path int true "ID"
// @Param lastupdatestamp query int true "Lastupdatestamp"
// @Param menu_url query string true "MenuUrl"
// @Param line_no query int true "LineNo"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/dynamic/{id} [get]
func (c *ContDynamic) GetById(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE     = tool.Res{R: e} // wajib
		ParamGet = models.ParamGet{}
		// id   = e.Param("id")  //util.StrTo(e.Param("id")).MustInt() //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	httpCode, errMsg := app.BindAndValid(e, &ParamGet)
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}

	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	data, err := c.useOption.GetDataBy(ctx, claims, ParamGet)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}
	return appE.Response(http.StatusOK, "Ok", data)

}

// GetList :
// @Summary GetList Dynamic
// @Security ApiKeyAuth
// @Tags Dynamic
// @Produce  json
// @Param req body models.ParamDynamicList true "req param #changes are possible to adjust the form of the registration form from frontend"
// @Success 200 {object} models.ResponseModelList
// @Router /barber/dynamic/list [post]
func (u *ContDynamic) GetList(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		// logger = logging.Logger{}
		appE = tool.Res{R: e} // wajib
		//valid      validation.Validation // wajib
		paramquery   = models.ParamDynamicList{} // ini untuk list
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
	// if !claims.IsAdmin {
	// 	paramquery.InitSearch = " id_created = " + strconv.Itoa(claims.UserID)
	// }

	responseList, err = u.useOption.GetList(ctx, claims, paramquery) //.GetList(ctx, paramquery)
	if err != nil {
		// return e.JSON(http.StatusBadRequest, err.Error())
		return appE.ResponseErrorList(tool.GetStatusCode(err), fmt.Sprintf("%v", err), responseList)
	}

	// return e.JSON(http.StatusOK, ListDataUser)
	return appE.Response(http.StatusOK, "", responseList)
}

// Delete :
// @Summary Delete Dynamic
// @Security ApiKeyAuth
// @Tags Dynamic
// @Produce  json
// @Param id path int true "ID"
// @Param lastupdatestamp query int true "Lastupdatestamp"
// @Param menu_url query string true "MenuUrl"
// @Param line_no query int true "LineNo"
// @Success 200 {object} tool.ResponseModel
// @Router /barber/dynamic [delete]
func (c *ContDynamic) Delete(e echo.Context) error {
	ctx := e.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		appE     = tool.Res{R: e} // wajib
		ParamGet = models.ParamGet{}
		// id   = e.Param("id")  //util.StrTo(e.Param("id")).MustInt() //kalo bukan int => 0
		// valid  validation.Validation                 // wajib
	)
	httpCode, errMsg := app.BindAndValid(e, &ParamGet)
	if httpCode != 200 {
		return appE.ResponseError(http.StatusBadRequest, errMsg, nil)
	}
	claims, err := app.GetClaims(e)
	if err != nil {
		return appE.ResponseError(http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}
	err = c.useOption.Delete(ctx, claims, ParamGet)
	if err != nil {
		return appE.ResponseError(http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
	}

	return appE.Response(http.StatusNoContent, "Ok", nil)

}
