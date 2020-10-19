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
			return data, models.ErrNotFound
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

	sQuery := fmt.Sprintf(`
		SELECT *
		FROM fbarber_beranda_s(%s)
		WHERE %s
	`, queryparam.ParamView, sWhere)
	query := db.Conn.Raw(sQuery).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result) //Find(&result)

	// sSourceFrom := fmt.Sprintf("fbarber_beranda_s(%s)", queryparam.ParamView)
	// // end where
	// query := db.Conn.Table(sSourceFrom).Select(`
	// barber_id,barber_name,file_id,file_name,file_path,file_type,price
	// `).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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
	// query := db.Conn.Table("barber").Select(`barber.barber_id ,barber.barber_name ,barber.file_id ,
	// 										sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type ,
	// 										sum(order_d.price ) as price`).Joins(`
	// 										left join sa_file_upload  on sa_file_upload.file_id = barber.file_id`).Joins(`
	// 										left join order_d on order_d.barber_id = barber.barber_id`).Group(`barber.barber_id ,barber.barber_name ,barber.file_id ,
	// 										sa_file_upload.file_name ,sa_file_upload.file_path ,sa_file_upload.file_type`).Where(sWhere).Count(&result)

	// sSourceFrom := fmt.Sprintf("fbarber_beranda_s(%s)", queryparam.ParamView)
	// query := db.Conn.Table(sSourceFrom).Select(`
	// barber_id,barber_name,file_id,file_name,file_path,file_type,price
	// `).Where(sWhere).Count(&result)

	sQuery := fmt.Sprintf(`
		SELECT count(*) as cnt
		FROM fbarber_beranda_s(%s)
		WHERE %s
	`, queryparam.ParamView, sWhere)
	query := db.Conn.Raw(sQuery).First(&op)

	// query := db.Conn.Table("v_order_h").Select(`
	// coalesce(sum(order_d.price ),0) as price
	// `).Joins(`inner join order_d
	// on v_order_h.order_id = order_d.order_id
	// `).Where(sWhere).First(&op)

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op.Cnt, nil
}
