package routes

import (
	"nuryanto2121/dynamic_rest_api_go/pkg/postgresdb"
	// sqlxposgresdb "nuryanto2121/dynamic_rest_api_go/pkg/postgresqlxdb"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	_saauthcont "nuryanto2121/dynamic_rest_api_go/controllers/auth"
	_authuse "nuryanto2121/dynamic_rest_api_go/usecase/auth"

	_saFilecont "nuryanto2121/dynamic_rest_api_go/controllers/fileupload"
	_repoFile "nuryanto2121/dynamic_rest_api_go/repository/ss_fileupload"
	_useFile "nuryanto2121/dynamic_rest_api_go/usecase/ss_fileupload"

	_contUser "nuryanto2121/dynamic_rest_api_go/controllers/user"
	_repoUser "nuryanto2121/dynamic_rest_api_go/repository/ss_user"
	_useUser "nuryanto2121/dynamic_rest_api_go/usecase/ss_user"

	_contPaket "nuryanto2121/dynamic_rest_api_go/controllers/paket"
	_repoPaket "nuryanto2121/dynamic_rest_api_go/repository/paket"
	_usePaket "nuryanto2121/dynamic_rest_api_go/usecase/paket"

	_contCapster "nuryanto2121/dynamic_rest_api_go/controllers/capster"
	_repoCapster "nuryanto2121/dynamic_rest_api_go/repository/capster"
	_useCapster "nuryanto2121/dynamic_rest_api_go/usecase/capster"

	_contBarber "nuryanto2121/dynamic_rest_api_go/controllers/barber"
	_repoBarber "nuryanto2121/dynamic_rest_api_go/repository/barber"
	_repoBarberCapster "nuryanto2121/dynamic_rest_api_go/repository/barber_capster"
	_repoBarberPaket "nuryanto2121/dynamic_rest_api_go/repository/barber_paket"
	_useBarber "nuryanto2121/dynamic_rest_api_go/usecase/barber"

	"time"

	"github.com/labstack/echo/v4"
)

//Echo :
type EchoRoutes struct {
	E *echo.Echo
}

func (e *EchoRoutes) InitialRouter() {
	timeoutContext := time.Duration(setting.FileConfigSetting.Server.ReadTimeout) * time.Second

	repoUser := _repoUser.NewRepoSysUser(postgresdb.Conn)
	useUser := _useUser.NewUserSysUser(repoUser, timeoutContext)
	_contUser.NewContUser(e.E, useUser)

	repoPaket := _repoPaket.NewRepoPaket(postgresdb.Conn)
	usePaket := _usePaket.NewUserMPaket(repoPaket, timeoutContext)
	_contPaket.NewContPaket(e.E, usePaket)

	repoFile := _repoFile.NewRepoFileUpload(postgresdb.Conn)
	useFile := _useFile.NewSaFileUpload(repoFile, timeoutContext)
	_saFilecont.NewContFileUpload(e.E, useFile)

	repoCapster := _repoCapster.NewRepoCapsterCollection(postgresdb.Conn)
	useCapster := _useCapster.NewUserMCapster(repoCapster, repoUser, repoFile, timeoutContext)
	_contCapster.NewContCapster(e.E, useCapster)

	repoBarberPaket := _repoBarberPaket.NewRepoBarberPaket(postgresdb.Conn)
	repoBarberCapster := _repoBarberCapster.NewRepoBarberCapster(postgresdb.Conn)
	repoBarber := _repoBarber.NewRepoBarber(postgresdb.Conn)
	useBarber := _useBarber.NewUserMBarber(repoBarber, repoBarberPaket, repoBarberCapster, timeoutContext)
	_contBarber.NewContBarber(e.E, useBarber)

	//_saauthcont
	// repoAuth := _repoAuth.NewRepoOptionDB(postgresdb.Conn)
	useAuth := _authuse.NewUserAuth(repoUser, repoFile, timeoutContext)
	_saauthcont.NewContAuth(e.E, useAuth)

}
