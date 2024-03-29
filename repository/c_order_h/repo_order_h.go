package repoorderh

import (
	"fmt"
	iorder_h "nuryanto2121/cukur_in_barber/interface/c_order_h"
	"nuryanto2121/cukur_in_barber/models"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	"nuryanto2121/cukur_in_barber/pkg/setting"
	"strings"

	"github.com/jinzhu/gorm"
)

type repoOrderH struct {
	Conn *gorm.DB
}

func NewRepoOrderH(Conn *gorm.DB) iorder_h.Repository {
	return &repoOrderH{Conn}
}

func (db *repoOrderH) GetDataBy(ID int) (result models.OrderHGet, err error) {
	var (
		logger = logging.Logger{}
		data   models.OrderHGet
	)
	query := db.Conn.Raw(`select barber.barber_id,barber.barber_cd,barber.barber_name ,order_h.capster_id ,ss_user."name" as capster_name,
								sa_file_upload.file_id ,sa_file_upload.file_name,sa_file_upload.file_path ,
								sum(order_d.price) as price ,order_h.order_date ,order_h.from_apps,order_h.status
							from order_h inner join order_d 
							on order_h.order_id = order_d.order_id 
							left join barber on barber.barber_id =order_h.barber_id 
							left join ss_user on ss_user.user_id = order_h.capster_id
							left join sa_file_upload on sa_file_upload.file_id = ss_user.file_id
							where order_h.order_id = ?
							group by barber.barber_name ,barber.barber_cd,order_h.capster_id ,ss_user."name",
								sa_file_upload.file_id ,sa_file_upload.file_name,sa_file_upload.file_path,
								order_h.order_date ,barber.barber_id,order_h.from_apps,order_h.status

				`, ID).Scan(&data) //Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, models.ErrNotFound
		}
		return result, err
	}
	return data, nil
}
func (db *repoOrderH) GetList(queryparam models.ParamList) (result []*models.OrderList, err error) {

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
		query = db.Conn.Table("v_order_h").Select(`
		owner_id,barber_id ,
		barber_name ,order_id ,
		status ,from_apps ,
		capster_id ,order_date ,
		capster_name,
		file_id ,file_name,
		file_path , price ,
		weeks,years,months
	`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Table("v_order_h").Select(`
		owner_id,barber_id ,
		barber_name ,order_id ,
		status ,from_apps ,
		capster_id ,order_date ,
		capster_name,
		file_id ,file_name,
		file_path , price ,
		weeks,years,months
	`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	}

	// end where
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
func (db *repoOrderH) Create(data *models.OrderH) error {
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
func (db *repoOrderH) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.OrderH{}).Where("order_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoOrderH) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("order_id = ?", ID).Delete(&models.OrderH{})
	query := db.Conn.Exec("Delete From order_h WHERE order_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoOrderH) Count(queryparam models.ParamList) (result int, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
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
		query = db.Conn.Table("v_order_h").Select(`
		owner_id,barber_id ,
		barber_name ,order_id ,
		status ,from_apps ,
		capster_id ,order_date ,
		capster_name,
		file_id ,file_name,
		file_path , price ,
		weeks,years,months
	`).Where(sWhere, queryparam.Search).Count(&result)

	} else {
		query = db.Conn.Table("v_order_h").Select(`
		owner_id,barber_id ,
		barber_name ,order_id ,
		status ,from_apps ,
		capster_id ,order_date ,
		capster_name,
		file_id ,file_name,
		file_path , price ,
		weeks,years,months
	`).Where(sWhere).Count(&result)
	}

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (db *repoOrderH) SumPriceDetail(queryparam models.ParamList) (result float32, err error) {
	type Results struct {
		Price float32 `json:"price"`
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
	sWhere = strings.ReplaceAll(sWhere, "barber_id", "v_order_h.barber_id")

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Table("v_order_h").Select(`
		coalesce(sum(order_d.price ),0) as price
		`).Joins(`inner join order_d 
		on v_order_h.order_id = order_d.order_id
		`).Where(sWhere, queryparam.Search).First(&op)
	} else {
		query = db.Conn.Table("v_order_h").Select(`
		coalesce(sum(order_d.price ),0) as price
		`).Joins(`inner join order_d 
		on v_order_h.order_id = order_d.order_id
		`).Where(sWhere).First(&op)
	}

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op.Price, nil
}
