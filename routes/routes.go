package routes

import (
	"nuryanto2121/dynamic_rest_api_go/pkg/postgresdb"
	sqlxposgresdb "nuryanto2121/dynamic_rest_api_go/pkg/postgresqlxdb"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	_contDynamic "nuryanto2121/dynamic_rest_api_go/controllers/dynamic"
	_repoDynamic "nuryanto2121/dynamic_rest_api_go/repository/option_template"
	_useDynamic "nuryanto2121/dynamic_rest_api_go/usecase/dynamic"

	_saauthcont "nuryanto2121/dynamic_rest_api_go/controllers/auth"
	_repoAuth "nuryanto2121/dynamic_rest_api_go/repository/auth"
	_authuse "nuryanto2121/dynamic_rest_api_go/usecase/auth"

	_saFilecont "nuryanto2121/dynamic_rest_api_go/controllers/fileupload"
	_repoFile "nuryanto2121/dynamic_rest_api_go/repository/ss_fileupload"
	_useFile "nuryanto2121/dynamic_rest_api_go/usecase/ss_fileupload"

	_contUser "nuryanto2121/dynamic_rest_api_go/controllers/user"
	_repoUser "nuryanto2121/dynamic_rest_api_go/repository/ss_user"
	_useUser "nuryanto2121/dynamic_rest_api_go/usecase/ss_user"

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

	repoUser := _repoUser.NewRepoSysUser(postgresdb.Conn)
	useUser := _useUser.NewUserSysUser(repoUser, timeoutContext)
	_contUser.NewContUser(e.E, useUser)

	repoFile := _repoFile.NewRepoFileUpload(sqlxposgresdb.DbCon)
	useFile := _useFile.NewSaFileUpload(repoFile, timeoutContext)
	_saFilecont.NewContFileUpload(e.E, useFile)

	//_saauthcont
	repoAuth := _repoAuth.NewRepoOptionDB(sqlxposgresdb.DbCon)
	useAuth := _authuse.NewUserAuth(repoAuth, timeoutContext)
	_saauthcont.NewContAuth(e.E, useAuth)

}
