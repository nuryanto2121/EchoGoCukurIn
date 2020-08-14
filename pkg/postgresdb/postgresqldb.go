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
	)

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
