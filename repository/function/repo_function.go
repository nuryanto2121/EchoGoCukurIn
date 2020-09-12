package repofunction

import (
	"fmt"
	"nuryanto2121/dynamic_rest_api_go/models"
	"nuryanto2121/dynamic_rest_api_go/pkg/logging"
	"nuryanto2121/dynamic_rest_api_go/pkg/postgresdb"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"strings"

	"github.com/jinzhu/gorm"
)

type FN struct {
	Claims util.Claims
}

func (fn *FN) GenBarberCode() (string, error) {
	var (
		result string
		conn   *gorm.DB
		logger = logging.Logger{}
		mSeqNo = &models.SsSequenceNo{}
		// []prefix = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	)
	conn = postgresdb.Conn

	prefixArr := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	fmt.Printf("%v", prefixArr)
	// ss := prefixArr[0]
	// query := conn.Table("barber").Select("max(barber_cd)") //
	query := conn.Where("sequence_cd = ?", "seq_barber").Find(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err := query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			mSeqNo.Prefix = "AA"
			mSeqNo.SeqNo = 1
			mSeqNo.SequenceCd = "seq_barber"
			mSeqNo.UserInput = fn.Claims.UserID
			mSeqNo.UserEdit = fn.Claims.UserID
			queryC := conn.Create(mSeqNo)
			logger.Query(fmt.Sprintf("%v", queryC.QueryExpr()))
			err = queryC.Error
			if err != nil {
				return "", err
			}
			result = "AA01"
			return result, nil
		}
		return "", err
	}

	if mSeqNo.SeqNo == 99 {
		mSeqNo.SeqNo = 1

	}

	return result, nil
}