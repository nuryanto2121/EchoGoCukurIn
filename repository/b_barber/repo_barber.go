package repobarber

import (
	"fmt"
	ibarber "nuryanto2121/dynamic_rest_api_go/interface/b_barber"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBarber struct {
	Conn *gorm.DB
}

func NewRepoBarber(Conn *gorm.DB) ibarber.Repository {
	return &repoBarber{Conn}
}

func (db *repoBarber) GetDataBy(ID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
	)
	query := db.Conn.Where("barber_id = ? ", ID).Find(mBarber)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarber, nil
}
func (db *repoBarber) GetDataFirst(OwnerID int, BarberID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
		sQuery  = ""
	)
	// query := db.Conn.First(&mBarber)
	if BarberID == 0 {
		sQuery = `SELECT * FROM barber where is_active = true and owner_id = ? 
		order by barber_id 
		limit 1`

		query := db.Conn.Raw(sQuery, OwnerID).Scan(&mBarber)
		logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
		err = query.Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
	} else {
		sQuery = `SELECT * FROM barber where is_active = true and owner_id = ? AND barber_id = ?
		limit 1`

		query := db.Conn.Raw(sQuery, OwnerID, BarberID).Scan(&mBarber)
		logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
		err = query.Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
	}

	return mBarber, nil
}
func (db *repoBarber) GetList(queryparam models.ParamList) (result []*models.BarbersList, err error) {

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
		query = db.Conn.Table("barber b ").Select(`
		b.barber_id,b.barber_cd,b.barber_name,
		b.address,b.latitude,b.longitude,
		b.operation_start,b.operation_end,
		b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type
		`).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
	`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Table("barber b ").Select(`
		b.barber_id,b.barber_cd,b.barber_name,
		b.address,b.latitude,b.longitude,
		b.operation_start,b.operation_end,
		b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type
		`).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
	`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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
func (db *repoBarber) Create(data *models.Barber) error {
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
func (db *repoBarber) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Barber{}).Where("barber_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarber) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_id = ?", ID).Delete(&models.Barber{})
	query := db.Conn.Exec("Delete From barber_collection WHERE barber_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarber) Count(queryparam models.ParamList) (result int, err error) {
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

		query = db.Conn.Model(&models.Barber{}).Where(sWhere, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.Model(&models.Barber{}).Where(sWhere).Count(&result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
