package repoorderh

import (
	"fmt"
	iberanda_barber "nuryanto2121/dynamic_rest_api_go/interface/beranda_barber"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBerandaBarber struct {
	Conn *gorm.DB
}

func NewRepoBerandaBarber(Conn *gorm.DB) iberanda_barber.Repository {
	return &repoBerandaBarber{Conn}
}

func (db *repoBerandaBarber) GetStatusOrder(ID int) (result models.Beranda, err error) {
	var (
		logger = logging.Logger{}
		data   models.Beranda
	)

	query := db.Conn.Raw(`select 
			count(case a.status when 'P' then 1 else null end) as progress_status,
			count(case a.status when 'F' then 1 else null end) as finish_status,
			count(case a.status when 'C' then 1 else null end) as cancel_status,
			sum(d.price ) as income_price
		from order_h a inner join barber b on a.barber_id = b.barber_id 
		inner join order_d d on d.order_id = a.order_id 
			and d.barber_id = b.barber_id 
		where b.owner_id = ?`, ID).Scan(&data) //Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, models.ErrNotFound
		}
		return data, err
	}

	return data, nil
}

func (db *repoBerandaBarber) GetListOrder(queryparam models.ParamList) (result []*models.BerandaList, err error) {
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
	query := db.Conn.Table("barber").Select(`barber.barber_id ,barber.barber_name ,barber.file_id ,
											sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type ,
											sum(order_d.price ) as price`).Joins(`
											left join sa_file_upload  on sa_file_upload.file_id = barber.file_id`).Joins(`
											left join order_d on order_d.barber_id = barber.barber_id`).Group(`barber.barber_id ,barber.barber_name ,barber.file_id ,
											sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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

func (db *repoBerandaBarber) Count(queryparam models.ParamList) (result int, err error) {
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

	// query := db.Conn.Model(&models.Barber{}).Where(sWhere).Count(&result)
	query := db.Conn.Table("barber").Select(`barber.barber_id ,barber.barber_name ,barber.file_id ,
											sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type ,
											sum(order_d.price ) as price`).Joins(`
											left join sa_file_upload  on sa_file_upload.file_id = barber.file_id`).Joins(`
											left join order_d on order_d.barber_id = barber.barber_id`).Group(`barber.barber_id ,barber.barber_name ,barber.file_id ,
											sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type`).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
