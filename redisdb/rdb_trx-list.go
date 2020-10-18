package redisdb

import (
	"fmt"
	"time"
)

// GetList :
func GetList(key string) ([]string, error) {
	list, err := rdb.SMembers(key).Result()
	return list, err
}

// RemoveList :
func RemoveList(key string, val interface{}) error {
	_, err := rdb.SRem(key, val).Result()
	if err != nil {
		return err
	}
	return nil
}

// AddList :
func AddList(key, val string) error {
	_, err := rdb.SAdd(key, val).Result()
	if err != nil {
		return err
	}
	return nil
}

// TurncateList :
func TurncateList(key string) error {
	_, err := rdb.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

// AddSession :
func AddSession(key string, val interface{}, mn time.Duration) error {
	set, err := rdb.Set(key, val, mn).Result()
	if err != nil {
		return err
	}
	fmt.Println(set)
	return nil
}

// GetSession :
func GetSession(key string) interface{} {
	value := rdb.Get(key).Val()
	fmt.Println(value)
	return value
}
