package logging

import (
	"context"
	"log"
	"nuryanto2121/dynamic_rest_api_go/pkg/monggodb"
	util "nuryanto2121/dynamic_rest_api_go/pkg/utils"
	"time"
)

type auditLog struct {
	ID        int64     `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Level     string    `json:"level"`
	UUID      string    `json:"uuid"`
	FuncName  string    `json:"func_name"`
	FileName  string    `json:"file_name"`
	Line      int       `json:"line"`
	Time      string    `json:"time"`
	Message   string    `json:"message"`
}

func (a *auditLog) saveAudit() {
	db, err := monggodb.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	a.ID = util.GetTimeNow().Unix()
	_, err = db.Collection("auditlogs").InsertOne(context.TODO(), a)
	if err != nil {
		log.Fatal(err.Error())
	}

}
