package postgresdb

import (
	"fmt"
	"log"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // add database driver bridge
)

var Conn *gorm.DB

func Setup() {
	now := time.Now()
	var err error
	fmt.Print(setting.FileConfigSetting.Database)
	connectionstring := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		setting.FileConfigSetting.Database.User,
		setting.FileConfigSetting.Database.Password,
		setting.FileConfigSetting.Database.Name,
		setting.FileConfigSetting.Database.Host,
		setting.FileConfigSetting.Database.Port)
	fmt.Printf("%s", connectionstring)
	Conn, err = gorm.Open(setting.FileConfigSetting.Database.Type, connectionstring)
	if err != nil {
		log.Printf("connection.setup err : %v", err)
		panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.FileConfigSetting.Database.TablePrefix + defaultTableName
	}
	Conn.SingularTable(true)
	Conn.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	Conn.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	Conn.Callback().Delete().Replace("gorm:delete", deleteCallback)

	Conn.DB().SetMaxIdleConns(10)
	Conn.DB().SetMaxOpenConns(100)

	go autoMigrate()

	timeSpent := time.Since(now)
	log.Printf("Config database is ready in %v", timeSpent)
}

// autoMigrate : create or alter table from struct
func autoMigrate() {
	// Add auto migrate bellow this line
	Conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	log.Println("STARTING AUTO MIGRATE ")
	Conn.AutoMigrate(
		models.SsUser{},
		models.Paket{},
		models.CapsterCollection{},
		models.Barber{},
		models.BarberPaket{},
		models.BarberCapster{},
		models.SaFileUpload{},
		models.OrderH{},
		models.OrderD{},
		models.SsSequenceNo{},
	)

	Conn.Exec(`
	CREATE OR REPLACE FUNCTION public.last_day(date)
	RETURNS date AS
	$$
  		SELECT (date_trunc('MONTH', $1) + INTERVAL '1 MONTH - 1 day')::date;
	$$ LANGUAGE 'sql' IMMUTABLE STRICT;
			
	CREATE OR REPLACE FUNCTION public.week_of_month(
		p_date        DATE,
		p_direction   INT -- DEFAULT 1 -- for 8.4 and above
	  ) RETURNS INT AS
	  $$
		SELECT CASE WHEN $2 >= 0 THEN
		  CEIL(EXTRACT(DAY FROM $1) / 7)::int
		ELSE 
		  0 - CEIL(
			(EXTRACT(DAY FROM last_day($1)) - EXTRACT(DAY FROM $1) + 1) / 7
		  )::int
		END
	  $$ LANGUAGE 'sql' IMMUTABLE;
	  
	 
	  create or replace view v_order_h
	  as 
	  SELECT 	barber.owner_id,barber.barber_id ,
			  barber.barber_name ,order_h.order_id ,
			  order_h.status ,order_h.from_apps ,
			  order_h.capster_id ,order_h.order_date ,
			  ss_user."name" as capster_name,
			  ss_user.file_id ,sa_file_upload.file_name,
			  sa_file_upload.file_path ,
			  (select sum(order_d.price ) from order_d where order_d.order_id = order_h.order_id ) as price ,
			  week_of_month(order_h.order_date::date,1) as weeks,
			  extract (year from order_h.order_date ) as years,
			  extract (month from order_h.order_date ) as months
	  FROM "barber" inner join order_h 
			  on order_h.barber_id = barber.barber_id 
		  inner join ss_user 
			  on ss_user.user_id = order_h.capster_id 
		  left join sa_file_upload 
			  on sa_file_upload.file_id = ss_user.file_id ;
	 
	  `)

	log.Println("FINISHING AUTO MIGRATE ")
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		// nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("TimeInput"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(util.GetTimeNow())
			}
		}

		if modifyTimeField, ok := scope.FieldByName("TimeEdit"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(util.GetTimeNow())
			}
		}

	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("TimeEdit", util.GetTimeNow())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
