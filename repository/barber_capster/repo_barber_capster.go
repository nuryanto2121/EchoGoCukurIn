package repobarbercapster

import (
	"fmt"
	ibarbercapster "nuryanto2121/dynamic_rest_api_go/interface/barber_capster"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBarberCapster struct {
	Conn *gorm.DB
}

func NewRepoBarberCapster(Conn *gorm.DB) ibarbercapster.Repository {
	return &repoBarberCapster{Conn}
}

func (db *repoBarberCapster) GetDataBy(ID int) (result *models.BarberCapster, err error) {
	var (
		logger         = logging.Logger{}
		mBarberCapster = &models.BarberCapster{}
	)
	query := db.Conn.Where("barber_id = ? ", ID).Find(mBarberCapster)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarberCapster, nil
}
func (db *repoBarberCapster) GetList(queryparam models.ParamList) (result []*models.BarberCapster, err error) {

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
	query := db.Conn.Table("ss_user").Select("ss_user.user_id as barber_capster_id,ss_user.name,ss_user.is_active,sa_file_upload.file_id,sa_file_upload.file_name,sa_file_upload.file_path,sa_file_upload.file_type, '0' as rating").Joins("left join sa_file_upload ON sa_file_upload.file_id = ss_user.file_id").Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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
func (db *repoBarberCapster) Create(data *models.BarberCapster) error {
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
func (db *repoBarberCapster) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.BarberCapster{}).Where("barber_capster_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_capster_id = ?", ID).Delete(&models.BarberCapster{})
	query := db.Conn.Exec("Delete From barber_capster_collection WHERE barber_capster_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) Count(queryparam models.ParamList) (result int, err error) {
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

	// query := db.Conn.Model(&models.BarberCapster{}).Where(sWhere).Count(&result)
	query := db.Conn.Table("ss_user").Select("ss_user.user_id as barber_capster_id,ss_user.name,ss_user.is_active, 0 as rating").Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
