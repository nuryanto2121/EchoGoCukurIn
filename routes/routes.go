package routes

import (
	"nuryanto2121/cukur_in_barber/pkg/postgresdb"
	// sqlxposgresdb "nuryanto2121/cukur_in_barber/pkg/postgresqlxdb"
	"nuryanto2121/cukur_in_barber/pkg/setting"

	_saauthcont "nuryanto2121/cukur_in_barber/controllers/auth"
	_authuse "nuryanto2121/cukur_in_barber/usecase/auth"

	_saFilecont "nuryanto2121/cukur_in_barber/controllers/fileupload"
	_repoFile "nuryanto2121/cukur_in_barber/repository/ss_fileupload"
	_useFile "nuryanto2121/cukur_in_barber/usecase/ss_fileupload"

	_contUser "nuryanto2121/cukur_in_barber/controllers/user"
	_repoUser "nuryanto2121/cukur_in_barber/repository/ss_user"
	_useUser "nuryanto2121/cukur_in_barber/usecase/ss_user"

	_contPaket "nuryanto2121/cukur_in_barber/controllers/b_paket"
	_repoPaket "nuryanto2121/cukur_in_barber/repository/b_paket"
	_usePaket "nuryanto2121/cukur_in_barber/usecase/b_paket"

	_contCapster "nuryanto2121/cukur_in_barber/controllers/b_capster"
	_repoCapster "nuryanto2121/cukur_in_barber/repository/b_capster"
	_useCapster "nuryanto2121/cukur_in_barber/usecase/b_capster"

	_contBerandaBarber "nuryanto2121/cukur_in_barber/controllers/beranda_barber"
	_repoBerandaBarber "nuryanto2121/cukur_in_barber/repository/beranda_barber"
	_useBerandaBarber "nuryanto2121/cukur_in_barber/usecase/beranda_barber"

	_contBarber "nuryanto2121/cukur_in_barber/controllers/b_barber"
	_repoBarber "nuryanto2121/cukur_in_barber/repository/b_barber"
	_repoBarberCapster "nuryanto2121/cukur_in_barber/repository/b_barber_capster"
	_repoBarberPaket "nuryanto2121/cukur_in_barber/repository/b_barber_paket"
	_useBarber "nuryanto2121/cukur_in_barber/usecase/b_barber"

	_contOrder "nuryanto2121/cukur_in_barber/controllers/c_order"
	_repoOrderd "nuryanto2121/cukur_in_barber/repository/c_order_d"
	_repoOrder "nuryanto2121/cukur_in_barber/repository/c_order_h"
	_useOrder "nuryanto2121/cukur_in_barber/usecase/c_order"

	_contNotification "nuryanto2121/cukur_in_barber/controllers/notification"
	_repoNotification "nuryanto2121/cukur_in_barber/repository/notification"
	_useNotification "nuryanto2121/cukur_in_barber/usecase/notification"

	"time"

	"github.com/labstack/echo/v4"
)

//Echo :
type EchoRoutes struct {
	E *echo.Echo
}

func (e *EchoRoutes) InitialRouter() {
	timeoutContext := time.Duration(setting.FileConfigSetting.Server.ReadTimeout) * time.Second

	repoNotif := _repoNotification.NewRepoNotification(postgresdb.Conn)
	useNotif := _useNotification.NewUseNotification(repoNotif, timeoutContext)
	_contNotification.NewContNotification(e.E, useNotif)

	repoFile := _repoFile.NewRepoFileUpload(postgresdb.Conn)
	useFile := _useFile.NewSaFileUpload(repoFile, timeoutContext)
	_saFilecont.NewContFileUpload(e.E, useFile)

	repoUser := _repoUser.NewRepoSysUser(postgresdb.Conn)
	useUser := _useUser.NewUserSysUser(repoUser, repoFile, timeoutContext)
	_contUser.NewContUser(e.E, useUser)

	repoPaket := _repoPaket.NewRepoPaket(postgresdb.Conn)
	usePaket := _usePaket.NewUserMPaket(repoPaket, repoNotif, timeoutContext)
	_contPaket.NewContPaket(e.E, usePaket)

	repoBarberPaket := _repoBarberPaket.NewRepoBarberPaket(postgresdb.Conn)
	repoBarberCapster := _repoBarberCapster.NewRepoBarberCapster(postgresdb.Conn)
	repoBarber := _repoBarber.NewRepoBarber(postgresdb.Conn)
	useBarber := _useBarber.NewUserMBarber(repoBarber, repoBarberPaket, repoBarberCapster, repoFile, timeoutContext)
	_contBarber.NewContBarber(e.E, useBarber)

	repoCapster := _repoCapster.NewRepoCapsterCollection(postgresdb.Conn)
	useCapster := _useCapster.NewUserMCapster(repoCapster, repoUser, repoBarberCapster, repoFile, repoNotif, timeoutContext)
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
