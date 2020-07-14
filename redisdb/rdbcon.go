package redisdb

import (
	"fmt"
	"log"
	"time"

	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// Setup :
func Setup() {
	now := time.Now()
	conString := fmt.Sprintf("%s:%d", setting.FileConfigSetting.RedisDBSetting.Host, setting.FileConfigSetting.RedisDBSetting.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     conString,
		Password: setting.FileConfigSetting.RedisDBSetting.Password,
		DB:       setting.FileConfigSetting.RedisDBSetting.DB,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		// logging.Error("0", err)
		// logging.Fatal("0", err)
	}
	// fmt.Println("Mem Cache is Ready...")

	timeSpent := time.Since(now)
	log.Printf("Config redis is ready in %v", timeSpent)
}
