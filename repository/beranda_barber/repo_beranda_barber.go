package repoorderh

import (
	"fmt"
	iberanda_barber "nuryanto2121/cukur_in_barber/interface/beranda_barber"
	"nuryanto2121/cukur_in_barber/models"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	"nuryanto2121/cukur_in_barber/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBerandaBarber struct {
	Conn *gorm.DB
}

func NewRepoBerandaBarber(Conn *gorm.DB) iberanda_barber.Repository {
	return &repoBerandaBarber{Conn}
}

func (db *repoBerandaBarber) GetStatusOrder(ParamView string, ID int) (result models.Beranda, err error) {
	var (
		logger = logging.Logger{}
		data   models.Beranda
	)

	sQuery := fmt.Sprintf(`
		SELECT *
		FROM fbarber_beranda_status(%s)
		WHERE owner_id = ?
	`, ParamView)
	query := db.Conn.Raw(sQuery, ID).Scan(&data) //Find(&result)
	// sSourceFrom := fmt.Sprintf("fbarber_beranda_status(%s)", ParamView)
	// query := db.Conn.Table(sSourceFrom).Select(`
	// *
	// `).Where("owner_id = ?", ID).Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		return data, err
	}

	return data, nil
}

func (db *repoBerandaBarber) GetListOrder(queryparam models.ParamDynamicList) (result []*models.BerandaList, err error) {
	var (
		pageNum  = 0
		pageSize = setting.FileConfigSetting.App.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		query    *gorm.DB
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
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sQuery := fmt.Sprintf(`
		SELECT *
		FROM fbarber_beranda_s(%s)
		WHERE %s
		`, queryparam.ParamView, sWhere)
		query = db.Conn.Raw(sQuery, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result) //Find(&result)
	} else {
		sQuery := fmt.Sprintf(`
		SELECT *
		FROM fbarber_beranda_s(%s)
		WHERE %s
		`, queryparam.ParamView, sWhere)
		query = db.Conn.Raw(sQuery).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result) //Find(&result)
	}

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

func (db *repoBerandaBarber) Count(queryparam models.ParamDynamicList) (result int, err error) {

	type Results struct {
		Cnt int `json:"cnt"`
	}

	var (
		sWhere = ""
		logger = logging.Logger{}
		op     = &Results{}
		query  *gorm.DB
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sQuery := fmt.Sprintf(`
		SELECT count(*) as cnt
		FROM fbarber_beranda_s(%s)
		WHERE %s
	`, queryparam.ParamView, sWhere)
		query = db.Conn.Raw(sQuery, queryparam.Search).First(&op)
	} else {
		sQuery := fmt.Sprintf(`
		SELECT count(*) as cnt
		FROM fbarber_beranda_s(%s)
		WHERE %s
	`, queryparam.ParamView, sWhere)
		query = db.Conn.Raw(sQuery).First(&op)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	return op.Cnt, nil
}
