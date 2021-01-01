package repofunction

import (
	"fmt"
	"log"
	"nuryanto2121/cukur_in_barber/models"
	"nuryanto2121/cukur_in_barber/pkg/logging"
	"nuryanto2121/cukur_in_barber/pkg/postgresdb"
	util "nuryanto2121/cukur_in_barber/pkg/utils"
	"nuryanto2121/cukur_in_barber/redisdb"
	"time"

	inotification "nuryanto2121/cukur_in_barber/interface/notification"
	fcmgetway "nuryanto2121/cukur_in_barber/pkg/fcm"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

type FN struct {
	Claims    util.Claims
	RepoNotif inotification.Repository
}

func (fn *FN) GenBarberCode() (string, error) {
	var (
		result string
		conn   *gorm.DB
		logger = logging.Logger{}
		mSeqNo = &models.SsSequenceNo{}
		abjad  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	)
	conn = postgresdb.Conn

	prefixArr := strings.Split(abjad, "")
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

		pref := strings.Split(mSeqNo.Prefix, "")
		i1 := strings.Index(abjad, pref[0])
		i2 := strings.Index(abjad, pref[1])
		aa := len(prefixArr) - 1
		if i2 == aa {
			i2 = 0
			i1 += 1
		} else {
			i2 += 1
		}

		mSeqNo.Prefix = fmt.Sprintf("%s%s", prefixArr[i1], prefixArr[i2])

	} else {
		mSeqNo.SeqNo += 1
	}

	sNo := mSeqNo.SeqNo + 100

	runes := []rune(strconv.Itoa(sNo))
	no := string(runes[1:])
	query = conn.Model(models.SsSequenceNo{}).Where("sequence_id = ?", mSeqNo.SequenceID).Updates(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%s%s", mSeqNo.Prefix, no)

	return result, nil
}

func (fn *FN) GetCountTrxCapster(ID int) int {
	var (
		logger = logging.Logger{}
		result = 0
		conn   *gorm.DB
	)
	OwnerID, _ := strconv.Atoi(fn.Claims.UserID)

	conn = postgresdb.Conn
	query := conn.Table("order_h oh").Select(`
		oh.order_id ,oh.order_no ,oh.barber_id ,b.barber_name ,b.owner_id ,oh.order_date ,oh.status ,oh.capster_id
	`).Joins(`
		join barber b on b.barber_id = oh.barber_id 
	`).Where(`
	oh.status in('P','N') AND b.owner_id = ? AND oh.capster_id = ?
		`, OwnerID, ID).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err := query.Error
	if err != nil {
		if err == models.ErrNotFound {
			return 0
		}
		logger.Error(err)
	}
	return result
}

func (fn *FN) GetCountTrxBarber(ID int) int {

	var (
		logger = logging.Logger{}
		result = 0
		conn   *gorm.DB
	)
	OwnerID, _ := strconv.Atoi(fn.Claims.UserID)

	conn = postgresdb.Conn
	query := conn.Table("order_h oh").Select(`
		oh.order_id ,oh.order_no ,oh.barber_id ,b.barber_name ,b.owner_id ,oh.order_date ,oh.status ,oh.capster_id
	`).Joins(`
		join barber b on b.barber_id = oh.barber_id 
	`).Where(`
	oh.status in('P','N') AND b.owner_id = ? AND oh.barber_id = ?
		`, OwnerID, ID).Count(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err := query.Error
	if err != nil {
		if err == models.ErrNotFound {
			return 0
		}
		logger.Error(err)
	}
	return result
}

func (fn *FN) GetOwnerData() (result *models.SsUser, err error) {
	var (
		logger   = logging.Logger{}
		mCapster = &models.SsUser{}
		conn     *gorm.DB
	)
	conn = postgresdb.Conn
	query := conn.Where("user_id = ? ", fn.Claims.UserID).Find(mCapster)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mCapster, nil
}

func (fn *FN) SendNotifNonAktifPaket(PaketID int) {
	type DataLoop struct {
		OrderID int `json:"order_id"`
		UserID  int `json:"user_id"`
	}
	var (
		mNotif     = models.Notification{}
		queryparam models.ParamList
		logger     = logging.Logger{}
		sDate      = time.Now().Format("2006-01-02")
		dataLoop   = []DataLoop{}
	)
	conn := postgresdb.Conn

	sSql := `
	select distinct oh.order_id ,oh.user_id
	from order_h oh join order_d od 
		on oh.order_id = od.order_id 
	join barber b2 
		on b2.barber_id = oh.barber_id 
	join barber_paket bp 
		on bp.barber_id = b2.barber_id 
		and bp.paket_id = od.paket_id 
	where oh.user_id > 0
	and oh.status = 'N'
	and oh.order_date::date = ?
	and bp.paket_id = ?
	`
	query := conn.Raw(sSql, sDate, PaketID).Find(&dataLoop)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err := query.Error
	if err != nil {
		log.Printf(" err : %v", err)
		panic(err)
	}

	mNotif.Title = "Paket Tidak Tersedia"
	mNotif.Descs = "Paket potong rambut kamu di nonaktifkan, silahkan ganti paket yang lain atau membatalkan pesanan"
	mNotif.NotificationStatus = "N"
	mNotif.NotificationType = "I"
	mNotif.UserEdit = fn.Claims.UserID
	mNotif.UserInput = fn.Claims.UserID
	mNotif.NotificationDate = time.Now()
	//Loping User yg booking dengan paket yg dibatalkan
	for _, data := range dataLoop {

		go func(dUser DataLoop) {

			mNotif.UserId = dUser.UserID
			mNotif.LinkId = dUser.OrderID

			TokenFCM := fmt.Sprintf("%v", redisdb.GetSession(strconv.Itoa(dUser.UserID)+"_fcm"))
			if TokenFCM != "" {
				TokenFCMArr := []string{TokenFCM}

				fcm := &fcmgetway.SendFCM{
					Title:       mNotif.Title,
					Body:        mNotif.Descs,
					JumlahNotif: 0,
					DeviceToken: TokenFCMArr,
				}
				go fcm.SendPushNotification()
			}
			queryparam.InitSearch = fmt.Sprintf("notification_status = 'N' AND user_id = %d and title = 'Paket Tidak Tersedia' ", mNotif.UserId)
			cnt, err := fn.RepoNotif.Count(queryparam)

			if cnt == 0 {
				err = fn.RepoNotif.Create(&mNotif)
				if err != nil {
					log.Printf(" err : %v", err)
					panic(err)
				}
			}

		}(data)
	}

}
