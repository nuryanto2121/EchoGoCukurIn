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
	  
	  CREATE OR REPLACE FUNCTION public.fbarber_beranda_s(p_status varchar, p_date varchar)
	  RETURNS 
	  TABLE(
		  owner_id integer,
		  barber_id integer,
		  barber_name varchar,
		  file_id integer,
		  file_name varchar,
		  file_path varchar,
		  file_type varchar,
		  price numeric
	  )
	  LANGUAGE plpgsql
	 AS $function$
	 DECLARE v_id INTEGER; 
	 BEGIN 	
		   RETURN QUERY                
			   select 	barber.owner_id ,
						 barber.barber_id ,
						 barber.barber_name ,
						 barber.file_id ,
						 sa_file_upload.file_name ,
						 sa_file_upload.file_path ,
						 sa_file_upload.file_type ,
						 (
							 select coalesce(sum(od.price),0) from order_d od join order_h oh
								 on oh.order_id = od.order_id 
							 where oh.barber_id = barber.barber_id 
							 and oh.order_date::date = p_date::date 
							 and oh.status =p_status
						 ) as price
						 from barber
						 left join sa_file_upload  on sa_file_upload.file_id = barber.file_id
	 ;
				
	 END;
	 $function$
	 ;

	 
		CREATE OR replace FUNCTION public.fbarber_beranda_status(p_status varchar, p_date varchar)
		RETURNS 
		TABLE(
			owner_id integer,
			progress_status integer,
			finish_status integer,
			cancel_status integer
			,income_price numeric
		)
		LANGUAGE plpgsql
		AS $function$
		DECLARE v_id INTEGER; 
		BEGIN 	
			RETURN QUERY                
			select 
				ss.owner_id ,
				ss.progress_status,
				ss.finish_status,
				ss.cancel_status
				,
				(
						select coalesce(sum(od.price ),0)::numeric 
					from order_h oh join order_d od
					on oh.order_id = od.order_id 
					where oh.order_date::date= p_date::date 
					and oh.status = p_status
					and oh.barber_id in(
						select b2.barber_id from barber b2 
						where b2.owner_id = ss.owner_id
					)
				)::numeric as income_price
			from (
				select b.owner_id,
				count(case a.status when 'P' then 1 else null end)::integer as progress_status,
							count(case a.status when 'F' then 1 else null end)::integer  as finish_status,
							count(case a.status when 'C' then 1 else null end)::integer  as cancel_status
				from order_h a join barber b
				on a.barber_id = b.barber_id 
				where a.order_date::date=p_date::date
				group by b.owner_id
			) ss		

		;
				
		END;
		$function$
		;

	  

		CREATE OR REPLACE VIEW public.v_order_h
		AS SELECT barber.owner_id,
			barber.barber_id,
			barber.barber_name,
			order_h.order_id,
			order_h.status,
			order_h.from_apps,
			order_h.capster_id,
			order_h.order_date,
			ss_user.name AS capster_name,
			ss_user.file_id,
			sa_file_upload.file_name,
			sa_file_upload.file_path,
			( SELECT sum(order_d.price) AS sum
				FROM order_d
				WHERE order_d.order_id = order_h.order_id) AS price,
			week_of_month(order_h.order_date::date, 1) AS weeks,
			date_part('year'::text, order_h.order_date) AS years,
			date_part('month'::text, order_h.order_date) AS months,
			order_h.customer_name,
			order_h.order_no
		FROM barber
			JOIN order_h ON order_h.barber_id = barber.barber_id
			JOIN ss_user ON ss_user.user_id = order_h.capster_id
			LEFT JOIN sa_file_upload ON sa_file_upload.file_id = ss_user.file_id;
	 
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
