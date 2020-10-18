package repocapstercollection

import (
	"fmt"
	icapster "nuryanto2121/dynamic_rest_api_go/interface/b_capster"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoCapsterCollection struct {
	Conn *gorm.DB
}

func NewRepoCapsterCollection(Conn *gorm.DB) icapster.Repository {
	return &repoCapsterCollection{Conn}
}

func (db *repoCapsterCollection) GetDataBy(ID int) (result *models.CapsterCollection, err error) {
	var (
		logger             = logging.Logger{}
		mCapsterCollection = &models.CapsterCollection{}
	)
	query := db.Conn.Where("capster_id = ? ", ID).Find(mCapsterCollection)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mCapsterCollection, nil
}
func (db *repoCapsterCollection) GetListFileCapter(ID int) (result []*models.SaFileOutput, err error) {
	var (
		logger = logging.Logger{}
	)
	query := db.Conn.Table("capster_collection").Select("capster_collection.file_id,sa_file_upload.file_name,sa_file_upload.file_path, sa_file_upload.file_type").Joins("Inner Join sa_file_upload ON sa_file_upload.file_id = capster_collection.file_id").Where("capster_id = ?", ID).Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db *repoCapsterCollection) GetList(queryparam models.ParamList) (result []*models.CapsterList, err error) {

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
	query := db.Conn.Table("ss_user").Select(`ss_user.user_id as capster_id,ss_user.user_name,ss_user.name,
								ss_user.is_active,sa_file_upload.file_id,sa_file_upload.file_name,
								sa_file_upload.file_path,sa_file_upload.file_type, 0 as rating,
								ss_user
	`).Joins(`
	left join sa_file_upload ON sa_file_upload.file_id = ss_user.file_id
	`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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
func (db *repoCapsterCollection) Create(data *models.CapsterCollection) error {
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
func (db *repoCapsterCollection) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.CapsterCollection{}).Where("capster_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoCapsterCollection) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("capster_id = ?", ID).Delete(&models.CapsterCollection{})
	query := db.Conn.Exec("Delete From capster_collection WHERE capster_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoCapsterCollection) Count(queryparam models.ParamList) (result int, err error) {
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

	// query := db.Conn.Model(&models.CapsterCollection{}).Where(sWhere).Count(&result)
	query := db.Conn.Table("ss_user").Select("ss_user.user_id as capster_id,ss_user.name,ss_user.is_active, 0 as rating").Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
