package reposysusers

import (
	"fmt"
	iusers "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

type repoSysUser struct {
	Conn *gorm.DB
}

func NewRepoSysUser(Conn *gorm.DB) iusers.Repository {
	return &repoSysUser{Conn}
}

func (db *repoSysUser) GetByEmailSaUser(email string) (result models.SysUser, err error) {

	query := db.Conn.Where("email = ?", email).First(&result)
	log.Info(fmt.Sprintf("%v", query.QueryExpr()))
	// logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error

	if err != nil {
		//
		if err == gorm.ErrRecordNotFound {
			return result, models.ErrNotFound
		}
		return result, err
	}

	return result, err
}

func (db *repoSysUser) GetDataBy(ID int) (result *models.SysUser, err error) {
	var sysUser = &models.SysUser{}
	query := db.Conn.Where("id = ? and deleted_on = ?", ID, 0).Find(sysUser)
	log.Info(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		return nil, err
	}
	return sysUser, nil
}

func (db *repoSysUser) GetList(queryparam models.ParamList) (result []*models.SysUser, err error) {

	var (
		pageNum  = 0
		pageSize = setting.FileConfigSetting.App.PageSize
		sWhere   = ""
		// logger   = logging.Logger{}
		orderBy = "id desc"
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
	if pageNum >= 0 && pageSize > 0 {
		query := db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
		fmt.Printf("%v", query.QueryExpr()) //cath to log query string
		err = query.Error
	} else {
		query := db.Conn.Where(sWhere).Order(orderBy).Find(&result)
		fmt.Printf("%v", query.QueryExpr()) //cath to log query string
		err = query.Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoSysUser) Create(data *models.SysUser) (err error) {

	err = db.Conn.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Update(ID int, data interface{}) (err error) {
	err = db.Conn.Model(models.SysUser{}).Where("id = ?", ID).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Delete(ID int) (err error) {
	if err := db.Conn.Where("id = ?", ID).Delete(&models.SysUser{}).Error; err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Count(queryparam models.ParamList) (result int, err error) {
	var (
		sWhere = ""
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

	err = db.Conn.Model(&models.SysUser{}).Where(sWhere).Count(&result).Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
