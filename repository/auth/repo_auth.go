package repoauth

import (
	"context"
	iauth "nuryanto2121/dynamic_rest_api_go/interface/auth"
	"nuryanto2121/dynamic_rest_api_go/models"

	// "nuryanto2121/dynamic_rest_api_go/pkg/logging"
	// queryauth "nuryanto2121/dynamic_rest_api_go/query/auth"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type repoAuth struct {
	DB *gorm.DB
}

func NewRepoOptionDB(Conn *gorm.DB) iauth.Repository {
	return &repoAuth{Conn}
}

func (m *repoAuth) GetDataLogin(ctx context.Context, Account string) (result models.DataLogin, err error) {

	// var logger = logging.Logger{}
	// logger.Query(queryauth.QueryAuthLogin, Account, Account)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (m *repoAuth) ChangePassword(ctx context.Context, data interface{}) (err error) {
	// var logger = logging.Logger{}
	// logger.Query(queryauth.QueryUpdatePassword, data)
	// _, errs := m.DB.NamedQueryContext(ctx, queryauth.QueryUpdatePassword, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *repoAuth) Register(ctx context.Context, dataUser models.SsUser) (err error) {
	// var logger = logging.Logger{}
	// logger.Query(queryauth.QueryRegister, dataUser)
	// _, err := m.DB.NamedExecContext(ctx, queryauth.QueryRegister, dataUser)
	if err != nil {
		return err
	}
	return nil
}
