package usedynamic

import (
	"context"
	"fmt"
	"math"
	idynamic "nuryanto2121/dynamic_rest_api_go/interface/dynamic"
	"nuryanto2121/dynamic_rest_api_go/models"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type useOptionTemplate struct {
	repoOption     idynamic.Repository
	contextTimeOut time.Duration
	claims         util.Claims
}

func NewUserSysUser(a idynamic.Repository, timeout time.Duration) idynamic.Usecase {
	return &useOptionTemplate{repoOption: a, contextTimeOut: timeout}
}

func (u *useOptionTemplate) Execute(ctx context.Context, claims util.Claims, data map[string]interface{}) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// parameter wajib
	OptionUrl := fmt.Sprintf("%v", data["option_url"])
	Method := fmt.Sprintf("%v", data["method"])
	LineNo, err := strconv.Atoi(fmt.Sprintf("%v", data["line_no"])) //data["line_no"].(int)
	if err != nil {
		return nil, err
	}

	if _, ok := data["option_url"]; ok {
		delete(data, "option_url")
	}
	if _, ok := data["method"]; ok {
		delete(data, "method")
	}
	if _, ok := data["line_no"]; ok {
		delete(data, "line_no")
	}

	OptionDbList, err := u.repoOption.GetOptionByUrl(ctx, OptionUrl)
	if err != nil {
		return nil, err
	}
	var DataOption = tool.FilterOptionList(OptionDbList, LineNo, Method)[0]
	fmt.Printf("%v", DataOption)

	SpName := DataOption.SP

	DataParameter, err := u.repoOption.GetParamFunction(ctx, SpName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", DataParameter)
	DataPostSP, err := tool.SetParameterSP(DataParameter, data, claims)
	if err != nil {
		return nil, err
	}

	sQuery := tool.QueryFunction(SpName, DataParameter)
	fmt.Printf(sQuery)
	resultPost, err := u.repoOption.CRUD(ctx, sQuery, DataPostSP)
	if err != nil {
		return nil, err
	}

	return resultPost, nil
}

func (u *useOptionTemplate) Delete(ctx context.Context, claims util.Claims, ParamGet models.ParamGet) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	OptionDbList, err := u.repoOption.GetOptionByUrl(ctx, ParamGet.MenuUrl)
	if err != nil {
		return err
	}
	var DataOption = tool.FilterOptionList(OptionDbList, ParamGet.LineNo, "DELETE")[0]
	fmt.Printf("%v", DataOption)
	SpName := DataOption.SP

	DataParameter, err := u.repoOption.GetParamFunction(ctx, SpName)
	if err != nil {
		return err
	}
	fmt.Printf("%v", DataParameter)
	DataPostSP := make(map[string]interface{}, 0)
	DataPostSP["p_row_id"] = ParamGet.ID
	DataPostSP["p_lastupdatestamp"] = ParamGet.Lastupdatestamp

	sQuery := tool.QueryFunction(SpName, DataParameter)
	fmt.Printf(sQuery)
	resultPost, err := u.repoOption.CRUD(ctx, sQuery, DataPostSP)
	if err != nil {
		return err
	}
	fmt.Printf("%v", resultPost)
	return nil
}
func (u *useOptionTemplate) GetDataBy(ctx context.Context, claims util.Claims, ParamGet models.ParamGet) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	OptionDbList, err := u.repoOption.GetOptionByUrl(ctx, ParamGet.MenuUrl)
	if err != nil {
		return nil, err
	}
	var DataOption = tool.FilterOptionList(OptionDbList, ParamGet.LineNo, "GETBYID")[0]
	fmt.Printf("%v", DataOption)
	SpName := DataOption.SP

	DataParameter, err := u.repoOption.GetParamFunction(ctx, SpName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", DataParameter)
	DataPostSP := make(map[string]interface{}, 0)
	DataPostSP["p_row_id"] = ParamGet.ID
	DataPostSP["p_lastupdatestamp"] = ParamGet.Lastupdatestamp

	sQuery := tool.QueryFunctionByID(SpName, DataParameter)
	fmt.Printf(sQuery)
	resultPost, err := u.repoOption.CRUD(ctx, sQuery, DataPostSP)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", resultPost)

	return resultPost, nil
}

func (u *useOptionTemplate) GetList(ctx context.Context, claims util.Claims, queryparam models.ParamDynamicList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		iStart         int
		iPerpage       int
		isViewFunction bool
		ViewName       string
		FieldList      []models.ParamFunction
		DefineColumns  string
		// DefineColumnFormat string
	)
	OptionDbList, err := u.repoOption.GetOptionByUrl(ctx, queryparam.MenuUrl)
	if err != nil {
		return result, err
	}

	iStart = queryparam.Page
	iPerpage = queryparam.PerPage
	isViewFunction = queryparam.ParamView != ""
	MenuUrl := queryparam.MenuUrl
	LineNo := queryparam.LineNo
	ParamWhere := queryparam.Search
	InitialWhere := queryparam.InitSearch
	sSortField := queryparam.SortField
	if sSortField == "" {
		sSortField = "ORDER BY time_edit desc"
	} else {
		sSortField = "ORDER BY " + sSortField
	}

	var DataOption = tool.FilterOptionList(OptionDbList, LineNo, "LIST")[0]
	fmt.Printf("%v", DataOption)
	ViewName = DataOption.SP

	DefineColumn, err := u.repoOption.GetDefineColumn(ctx, MenuUrl, LineNo)
	if err != nil {
		return result, err
	}

	FieldList, err = u.repoOption.GetFieldType(ctx, ViewName, isViewFunction)
	if err != nil {
		return result, err
	}
	AllColumnQuery, _, DefineSize, FieldWhere := tool.SetFieldList(FieldList, DefineColumn, 20, true)

	_, AllColumn, _, _ := tool.SetFieldList(FieldList, DefineColumn, 0, true)

	if DefineColumn.ColumnField != "" {
		DefineColumns = DefineColumn.ColumnField
		// DefineColumnFormat =
	} else {
		SpName := "fss_define_column_i"
		DefineColumns = "no," + AllColumn
		DataParameter, err := u.repoOption.GetParamFunction(ctx, SpName)
		if err != nil {
			return result, err
		}
		fmt.Printf("%v", DataParameter)
		DataPostSP := make(map[string]interface{}, 0)
		DataPostSP["p_option_url"] = MenuUrl
		DataPostSP["p_user_id"] = claims.UserID
		DataPostSP["p_subportfolio_id"] = 1
		DataPostSP["p_line_no"] = LineNo
		DataPostSP["p_column_field"] = DefineColumns

		sQuery := tool.QueryFunctionByID(SpName, DataParameter)
		fmt.Printf(sQuery)
		_, err = u.repoOption.CRUD(ctx, sQuery, DataPostSP)
		if err != nil {
			return result, err
		}
	}

	if InitialWhere != "" {
		InitialWhere = "WHERE " + InitialWhere
	}
	sWhere := strings.Replace(InitialWhere, "claims.user_id", claims.UserName, -1)
	sWhereLike := tool.SetWhereLikeList(FieldWhere, ParamWhere)

	if ParamWhere != "" {
		if sWhere != "" {
			sWhere += " AND " + sWhereLike
		} else {
			sWhere += " WHERE " + sWhereLike
		}
	}

	if queryparam.ParamView != "" {
		ViewName = fmt.Sprintf("%s(%s)", ViewName, queryparam.ParamView)
	}
	iOffset := (iStart * iPerpage) - iPerpage
	// DataList := make(map[string]interface{}, 0)
	// DataList["Limit"] = iPerpage
	// DataList["Offset"] = iOffset

	sQuery := tool.QueryFunctionList(ViewName, sSortField, AllColumnQuery, sWhere)
	result.Data, err = u.repoOption.GetDataList(ctx, sQuery, iPerpage, iOffset)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoOption.CountList(ctx, ViewName, sWhere)
	if err != nil {
		return result, err
	}
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page
	result.DefineSize = DefineSize
	result.DefineColumn = DefineColumns
	result.AllColumn = AllColumn
	return result, err
}
