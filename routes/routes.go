package routes

import (
	sqlxposgresdb "nuryanto2121/dynamic_rest_api_go/pkg/postgresqlxdb"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	_contDynamic "nuryanto2121/dynamic_rest_api_go/controllers/dynamic"
	_repoDynamic "nuryanto2121/dynamic_rest_api_go/repository/option_template"
	_useDynamic "nuryanto2121/dynamic_rest_api_go/usecase/dynamic"

	_saauthcont "nuryanto2121/dynamic_rest_api_go/controllers/auth"
	_repoAuth "nuryanto2121/dynamic_rest_api_go/repository/auth"
	_authuse "nuryanto2121/dynamic_rest_api_go/usecase/auth"

	"time"

	"github.com/labstack/echo/v4"
)

//Echo :
type EchoRoutes struct {
	E *echo.Echo
}

func (e *EchoRoutes) InitialRouter() {
	timeoutContext := time.Duration(setting.FileConfigSetting.Server.ReadTimeout) * time.Second

	repoDynamic := _repoDynamic.NewRepoOptionDB(sqlxposgresdb.DbCon)
	useDynamic := _useDynamic.NewUserSysUser(repoDynamic, timeoutContext)
	_contDynamic.NewContDynamic(e.E, useDynamic)

	// repoUser := _repoUser.NewRepoSysUser(postgresdb.Conn)
	// useUser := _useUser.NewUserSysUser(repoUser, timeoutContext)
	// _contUser.NewContUser(e.E, useUser)

	//_saauthcont
	repoAuth := _repoAuth.NewRepoOptionDB(sqlxposgresdb.DbCon)
	useAuth := _authuse.NewUserAuth(repoAuth, timeoutContext)
	_saauthcont.NewContAuth(e.E, useAuth)

}
