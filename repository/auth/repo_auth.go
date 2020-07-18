package repoauth

import (
	"context"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	"nuryanto2121/dynamic_rest_api_go/models"
	queryauth "nuryanto2121/dynamic_rest_api_go/query/auth"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type repoAuth struct {
	DB *sqlx.DB
}

func NewRepoOptionDB(Conn *sqlx.DB) iauth.Repository {
	return &repoAuth{Conn}
}

func (m *repoAuth) GetDataLogin(ctx context.Context, Account string) (result models.DataLogin, err error) {
	// fmt.Printf(queryoption.QueryGetListOption)
	// errs := m.DB.SelectContext(ctx, &result, queryauth.QueryAuthLogin, Account)
	errs := m.DB.GetContext(ctx, &result, queryauth.QueryAuthLogin, Account, Account)
	if errs != nil {
		return result, errs
	}
	return result, nil
}

func (m *repoAuth) ChangePassword(ctx context.Context, data interface{}) (err error) {
	_, errs := m.DB.NamedQueryContext(ctx, queryauth.QueryUpdatePassword, data)
	if errs != nil {
		return errs
	}
	return nil
}

func (m *repoAuth) Register(ctx context.Context, dataUser models.SysUser) error {
	_, err := m.DB.NamedExecContext(ctx, queryauth.QueryRegister, dataUser)
	if err != nil {
		return err
	}
	return nil
}
