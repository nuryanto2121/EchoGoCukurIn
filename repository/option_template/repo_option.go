package repooption

import (
	"context"
	"fmt"
	idynamic "nuryanto2121/dynamic_rest_api_go/interface/dynamic"
	"nuryanto2121/dynamic_rest_api_go/models"
	tool "nuryanto2121/dynamic_rest_api_go/pkg/tools"
	queryoption "nuryanto2121/dynamic_rest_api_go/query/option"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type repoOptionDB struct {
	DB *sqlx.DB
}

func NewRepoOptionDB(Conn *sqlx.DB) idynamic.Repository {
	return &repoOptionDB{Conn}
}

func (m *repoOptionDB) GetOptionByUrl(ctx context.Context, Url string) (result []models.OptionDB, err error) {
	// fmt.Printf(queryoption.QueryGetListOption)
	errs := m.DB.SelectContext(ctx, &result, queryoption.QueryGetListOption, Url)
	if errs != nil {
		return nil, errs
	}
	return result, nil
}
func (m *repoOptionDB) GetParamFunction(ctx context.Context, SpName string) (result []models.ParamFunction, err error) {
	err = m.DB.SelectContext(ctx, &result, queryoption.QueryGetListParamFunction, SpName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *repoOptionDB) CRUD(ctx context.Context, sQuery string, data interface{}) (result interface{}, err error) {
	rows, err := m.DB.NamedQueryContext(ctx, sQuery, data)
	if err != nil {
		return nil, err
	}
	result, err = tool.ResultQuery(rows)
	return result, nil
}
func (m *repoOptionDB) GetDataList(ctx context.Context, sQuery string, Limit int, Offset int) (result interface{}, err error) {
	// fmt.Printf(queryoption.QueryGetListOption)
	rows, err := m.DB.QueryxContext(ctx, sQuery, Limit, Offset)
	if err != nil {
		return nil, err
	}
	result, err = tool.ResultQuery(rows)
	return result, nil
}
func (m *repoOptionDB) GetDefineColumn(ctx context.Context, MenuUrl string, LineNo int) (result models.DefineColumn, err error) {
	// errs := m.DB.SelectContext(ctx, &result, queryoption.QueryDefineColumn, MenuUrl, LineNo)
	errs := m.DB.GetContext(ctx, &result, queryoption.QueryDefineColumn, MenuUrl, LineNo)
	if errs != nil {
		return result, errs
	}
	return result, nil
}

func (m *repoOptionDB) GetFieldType(ctx context.Context, SourceFrom string, isViewFunction bool) (result []models.ParamFunction, err error) {
	var Query string
	if isViewFunction {
		Query = queryoption.QueryResultFunctionType
	} else {
		Query = queryoption.QueryGetListFieldType
	}

	err = m.DB.SelectContext(ctx, &result, Query, SourceFrom)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *repoOptionDB) CountList(ctx context.Context, ViewName string, sWhere string) (int, error) {
	var count int
	sQueryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", ViewName, sWhere)
	err := m.DB.QueryRow(sQueryCount).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
