package repoorderh

import (
	"fmt"
	iorder_h "nuryanto2121/dynamic_rest_api_go/interface/c_order_h"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoOrderH struct {
	Conn *gorm.DB
}

func NewRepoOrderH(Conn *gorm.DB) iorder_h.Repository {
	return &repoOrderH{Conn}
}

func (db *repoOrderH) GetDataBy(ID int) (result *models.OrderH, err error) {
	var (
		logger  = logging.Logger{}
		mOrderH = &models.OrderH{}
	)
	query := db.Conn.Where("order_id = ? ", ID).Find(mOrderH)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mOrderH, nil
}
func (db *repoOrderH) GetList(queryparam models.ParamList) (result []*models.OrderH, err error) {

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

	query := db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	// query := db.Conn.Table("ss_user").Select("ss_user.user_id as order_id,ss_user.name,ss_user.is_active,sa_file_upload.file_id,sa_file_upload.file_name,sa_file_upload.file_path,sa_file_upload.file_type, '0' as rating").Joins("left join sa_file_upload ON sa_file_upload.file_id = ss_user.file_id").Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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

	query := db.Conn.Model(&models.OrderH{}).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
