package redisdb

import (
	"encoding/json"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"github.com/mitchellh/mapstructure"
)

// Verify :
type Verify struct {
	Email      string
	UserName   string
	VerifyLink string
}

// StoreVerify :
func StoreVerify(data interface{}) error {
	var verify Verify

	err := mapstructure.Decode(data, &verify)
	if err != nil {
		return err
	}

	mVerify := map[string]interface{}{
		"mail_type": "verify",
		"data":      verify,
	}

	dVerify, err := json.Marshal(mVerify)
	if err != nil {
		return err
	}

	_, err = rdb.SAdd(setting.FileConfigSetting.RedisDBSetting.Key, string(dVerify)).Result()
	if err != nil {
		return err
	}

	return nil
}
