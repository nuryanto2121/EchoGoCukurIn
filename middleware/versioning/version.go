package version

import (
	"fmt"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"

	"github.com/jinzhu/gorm"
)

type SsVersion struct {
	VersionID int    `json:"version_id" gorm:"PRIMARY_KEY"`
	OS        string `json:"os" gorm:"type:varchar(20)"`
	Version   int    `json:"version" gorm:"type:integer"`
}

func (V *SsVersion) GetVersion(Conn *gorm.DB) (result SsVersion, err error) {
	var logger = logging.Logger{}
	query := Conn.Where("os = ? and apps = 'barber' ", V.OS).First(&result)

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string

	err = query.Error

	if err != nil {
		return result, err
	}
	return result, nil
}
