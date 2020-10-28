package reposysusers

import (
	"fmt"
	iusers "nuryanto2121/dynamic_rest_api_go/interface/user"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

type repoSysUser struct {
	Conn *gorm.DB
}

func NewRepoSysUser(Conn *gorm.DB) iusers.Repository {
	return &repoSysUser{Conn}
}

func (db *repoSysUser) GetByAccount(Account string, IsOwner bool) (result models.SsUser, err error) {
	var query *gorm.DB

	if IsOwner {
		query = db.Conn.Where("(email iLIKE ? OR telp = ?) AND user_type = 'owner' ", Account, Account).First(&result)
	} else {
		query = db.Conn.Where("(email iLIKE ? OR telp = ?)", Account, Account).First(&result)
	}

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
func (db *repoSysUser) GetByCapster(Account string) (result models.LoginCapster, err error) {

	// query := db.Conn.Where("email = ?", Account).Or("telp = ?", Account).First(&result)
	query := db.Conn.Table("ss_user su").Select(`su.user_id as capster_id, su."name" as capster_name,su."password",su.email ,
						su.telp ,su.file_id ,sf.file_name ,sf.file_path ,b.barber_id ,b.barber_name,su.user_input as owner_id ,so."name" as owner_name,su.is_active,
						b.is_active as barber_is_active
						`).Joins(`
						left join barber_capster bc on su.user_id = bc.capster_id`).Joins(`
						left join barber b on b.barber_id = bc.barber_id `).Joins(`
						left join sa_file_upload sf on sf.file_id =su.file_id`).Joins(`
						left join ss_user so on so.user_id = su.user_input::integer `).Where(`
						(su.email iLike ? OR su.telp = ?)`, Account, Account).First(&result)

	// query := db.Conn.Table("ss_user su").Select(`su.user_id as capster_id, su."name" as capster_name,su."password",su.email ,
	// su.telp ,su.file_id ,sf.file_name ,sf.file_path ,b.barber_id ,b.barber_name,b.owner_id ,so."name" as owner_name`).Joins(`
	// inner join barber_capster bc on su.user_id = bc.capster_id`).Joins(`
	// inner join barber b on b.barber_id = bc.barber_id `).Joins(`
	// left join sa_file_upload sf on sf.file_id =su.file_id`).Joins(`
	// inner join ss_user so on so.user_id = b.owner_id `).Where(`
	// su.email = ?`, Account).Or(`su.telp = ?`, Account).First(&result)
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
func (db *repoSysUser) GetDataBy(ID int) (result *models.SsUser, err error) {
	var sysUser = &models.SsUser{}
	query := db.Conn.Where("user_id = ? ", ID).Find(sysUser)
	log.Info(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return sysUser, nil
}

func (db *repoSysUser) GetList(queryparam models.ParamList) (result []*models.UserList, err error) {

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
	query := db.Conn.Table("ss_user sa ").Select(`
				sa.user_id ,sa.user_name ,
				sa."name" ,sa.telp ,
				sa.email ,sa.is_active ,
				sa.join_date ,sa.user_type ,
				sa.file_id ,sf.file_name ,
				sf.file_path ,sf.file_type
	`).Joins(`
	left join sa_file_upload sf
	on sf.file_id = sa.file_id 
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
func (db *repoSysUser) Create(data *models.SsUser) error {
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
func (db *repoSysUser) Update(ID int, data map[string]interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
		tUser  models.SsUser
	)
	query := db.Conn.Model(&tUser).Where("user_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) UpdatePasswordByEmail(Email string, Password string) error {
	var (
		logger = logging.Logger{}
		err    error
		// tUser  models.SsUser
	)
	// query := db.Conn.Model(&tUser).Where("user_id = ?", ID).Updates(data)
	query := db.Conn.Exec(`UPDATE ss_user
							set password = ?
						  where user_type IN ('capster','owner')
						  AND email = ?`, Password, Email)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("user_id = ?", ID).Delete(&models.SsUser{})
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoSysUser) Count(queryparam models.ParamList) (result int, err error) {
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

	query := db.Conn.Model(&models.SsUser{}).Where(sWhere).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
func (db *repoSysUser) GenUserCapster() (string, error) {
	result := ""
	row := db.Conn.Table("ss_user").Select("max(user_name)").Row()
	//.Where("to_timestamp(created_on)::date = now()::date").Row()
	row.Scan(&result)
	// if result == "" {
	// 	return "error"
	// }
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return result
	// }
	// return result, nil
	var (
		currentTime     = time.Now()
		year            = currentTime.Year()
		month       int = int(currentTime.Month())
		day             = currentTime.Day()
	)

	// result := u.repoUser.GetMaxUserCapster()
	sYear := strconv.Itoa(year)[2:]
	var sMonth string = strconv.Itoa(month)
	if len(sMonth) == 1 {
		sMonth = fmt.Sprintf("0%s", sMonth)
	}
	var sDay string = strconv.Itoa(day)
	if len(sDay) == 1 {
		sDay = fmt.Sprintf("0%s", sDay)
	}
	seqNo := "0001"
	if result == "" {
		result = fmt.Sprintf("CP%s%s%s%v", sYear, sMonth, sDay, seqNo)
	} else {
		seqNo = fmt.Sprintf("1%s", result[9:])
		iSeqno, err := strconv.Atoi(seqNo)
		if err != nil {
			return "", err
		}
		iSeqno += 1
		seqNo = strconv.Itoa(iSeqno)[1:]
		result = fmt.Sprintf("CP%s%s%s%v", sYear, sMonth, sDay, seqNo)
	}
	return result, nil
}
