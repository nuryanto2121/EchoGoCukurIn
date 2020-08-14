package sqlxposgresdb

import (
	"fmt"
	"log"
	"net/url"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DbCon *sqlx.DB

func Setup() {
	dbHost := setting.FileConfigSetting.Database.Host     //`34.87.141.90` //viper.GetString(`database.host`)
	dbPort := setting.FileConfigSetting.Database.Port     //`1401`                                  //viper.GetString(`database.port`)
	dbUser := setting.FileConfigSetting.Database.User     //`postgres`                              //viper.GetString(`database.user`)
	dbPass := setting.FileConfigSetting.Database.Password //`postgres_dev`                          //viper.GetString(`database.password`)
	dbName := setting.FileConfigSetting.Database.Name     //`cukur_in`                                  //viper.GetString(`database.name`)
	// connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	connection := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		dbUser,
		dbPass,
		dbName,
		dbHost,
		dbPort)
	fmt.Printf("%s", connection)
	var err error
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	log.Printf(" %v", dsn)
	DbCon, err = sqlx.Open("postgres", connection)
	if err != nil && setting.FileConfigSetting.Debug {
		fmt.Println(err)
	}

	err = DbCon.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// defer func() {
	// 	err := DbCon.Close()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	timeSpent := time.Since(time.Now())
	log.Printf("Config database sqlx is ready in %v", timeSpent)

}
