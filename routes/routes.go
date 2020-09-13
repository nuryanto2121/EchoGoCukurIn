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

	_contPaket "nuryanto2121/dynamic_rest_api_go/controllers/b_paket"
	_repoPaket "nuryanto2121/dynamic_rest_api_go/repository/b_paket"
	_usePaket "nuryanto2121/dynamic_rest_api_go/usecase/b_paket"

	_contCapster "nuryanto2121/dynamic_rest_api_go/controllers/b_capster"
	_repoCapster "nuryanto2121/dynamic_rest_api_go/repository/b_capster"
	_useCapster "nuryanto2121/dynamic_rest_api_go/usecase/b_capster"

	_contBerandaBarber "nuryanto2121/dynamic_rest_api_go/controllers/beranda_barber"
	_repoBerandaBarber "nuryanto2121/dynamic_rest_api_go/repository/beranda_barber"
	_useBerandaBarber "nuryanto2121/dynamic_rest_api_go/usecase/beranda_barber"

	_contBarber "nuryanto2121/dynamic_rest_api_go/controllers/b_barber"
	_repoBarber "nuryanto2121/dynamic_rest_api_go/repository/b_barber"
	_repoBarberCapster "nuryanto2121/dynamic_rest_api_go/repository/b_barber_capster"
	_repoBarberPaket "nuryanto2121/dynamic_rest_api_go/repository/b_barber_paket"
	_useBarber "nuryanto2121/dynamic_rest_api_go/usecase/b_barber"

	_contOrder "nuryanto2121/dynamic_rest_api_go/controllers/c_order"
	_repoOrderd "nuryanto2121/dynamic_rest_api_go/repository/c_order_d"
	_repoOrder "nuryanto2121/dynamic_rest_api_go/repository/c_order_h"
	_useOrder "nuryanto2121/dynamic_rest_api_go/usecase/c_order"

	"time"

	"github.com/labstack/echo/v4"
)

//Echo :
type EchoRoutes struct {
	E *echo.Echo
}

func (e *EchoRoutes) InitialRouter() {
	timeoutContext := time.Duration(setting.FileConfigSetting.Server.ReadTimeout) * time.Second

	repoFile := _repoFile.NewRepoFileUpload(postgresdb.Conn)
	useFile := _useFile.NewSaFileUpload(repoFile, timeoutContext)
	_saFilecont.NewContFileUpload(e.E, useFile)

	repoUser := _repoUser.NewRepoSysUser(postgresdb.Conn)
	useUser := _useUser.NewUserSysUser(repoUser, repoFile, timeoutContext)
	_contUser.NewContUser(e.E, useUser)

	repoPaket := _repoPaket.NewRepoPaket(postgresdb.Conn)
	usePaket := _usePaket.NewUserMPaket(repoPaket, timeoutContext)
	_contPaket.NewContPaket(e.E, usePaket)

	repoBarberPaket := _repoBarberPaket.NewRepoBarberPaket(postgresdb.Conn)
	repoBarberCapster := _repoBarberCapster.NewRepoBarberCapster(postgresdb.Conn)
	repoBarber := _repoBarber.NewRepoBarber(postgresdb.Conn)
	useBarber := _useBarber.NewUserMBarber(repoBarber, repoBarberPaket, repoBarberCapster, repoFile, timeoutContext)
	_contBarber.NewContBarber(e.E, useBarber)

	repoCapster := _repoCapster.NewRepoCapsterCollection(postgresdb.Conn)
	useCapster := _useCapster.NewUserMCapster(repoCapster, repoUser, repoBarberCapster, repoFile, timeoutContext)
	_contCapster.NewContCapster(e.E, useCapster)

	repoOrderD := _repoOrderd.NewRepoOrderD(postgresdb.Conn)
	repoOrder := _repoOrder.NewRepoOrderH(postgresdb.Conn)
	useOrder := _useOrder.NewUserMOrder(repoOrder, repoOrderD, timeoutContext)
	_contOrder.NewContOrder(e.E, useOrder, useBarber)

	repoBeranda := _repoBerandaBarber.NewRepoBerandaBarber(postgresdb.Conn)
	UseBeranda := _useBerandaBarber.NewUserMBarber(repoBeranda, repoFile, timeoutContext)
	_contBerandaBarber.NewContBeranda(e.E, UseBeranda)

	//_saauthcont
	// repoAuth := _repoAuth.NewRepoOptionDB(postgresdb.Conn)
	useAuth := _authuse.NewUserAuth(repoUser, repoFile, timeoutContext)
	_saauthcont.NewContAuth(e.E, useAuth)

}
