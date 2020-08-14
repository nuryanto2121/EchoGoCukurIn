package repobarberpaket

import (
	"fmt"
	ibarberpaket "nuryanto2121/dynamic_rest_api_go/interface/b_barber_paket"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBarberPaket struct {
	Conn *gorm.DB
}

func NewRepoBarberPaket(Conn *gorm.DB) ibarberpaket.Repository {
	return &repoBarberPaket{Conn}
}

func (db *repoBarberPaket) GetDataBy(ID int) (result *models.BarberPaket, err error) {
	var (
		logger       = logging.Logger{}
		mBarberPaket = &models.BarberPaket{}
	)
	query := db.Conn.Where("barber_id = ? ", ID).Find(mBarberPaket)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarberPaket, nil
}
func (db *repoBarberPaket) GetList(queryparam models.ParamList) (result []*models.Paket, err error) {

	var (
		pageNum  = 0
		pageSize = setting.FileConfigSetting.App.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
	)
	// pagination
	if queryparam.Page > 0 {
		pageNum = (queryparam.Page - 1) * queryparam.PerPage
	}
	if queryparam.PerPage > 0 {
		pageSize = queryparam.PerPage
	}
	//end pagination

	// Order
	if queryparam.SortField != "" {
		orderBy = queryparam.SortField
	}
	//end Order by

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and " + queryparam.Search
		} else {
			sWhere += queryparam.Search
		}
	}

	// end where

	// query := db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	query := db.Conn.Table("barber_paket").Select("paket.paket_id,paket.owner_id,paket.paket_name,paket.descs,paket.durasi_start,paket.durasi_end,paket.price,paket.is_active,paket.is_promo,paket.promo_price,paket.promo_start,paket.promo_end").Joins("left join paket ON paket.paket_id = barber_paket.paket_id").Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoBarberPaket) Create(data *models.BarberPaket) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberPaket) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.BarberPaket{}).Where("barber_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberPaket) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_id = ?", ID).Delete(&models.BarberPaket{})
	query := db.Conn.Exec("Delete From barber_paket WHERE barber_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberPaket) Count(queryparam models.ParamList) (result int, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and " + queryparam.Search
		}
	}
	// end where

	// query := db.Conn.Model(&models.BarberPaket{}).Where(sWhere).Count(&result)
	query := db.Conn.Table("ss_user").Select("ss_user.user_id as barber_id,ss_user.name,ss_user.is_active, 0 as rating").Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
