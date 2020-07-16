package logging

import (
	"fmt"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.FileConfigSetting.App.RuntimeRootPath, setting.FileConfigSetting.App.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s.%s",
		time.Now().Format(setting.FileConfigSetting.App.TimeFormat),
		setting.FileConfigSetting.App.LogFileExt,
	)
}
