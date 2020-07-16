package monggodb

import (
	"context"
	"fmt"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sequence struct {
	ID            string `bson:"_id"`
	SequenceValue int    `bson:"sequence_value"`
}

func Connect() (*mongo.Database, error) {
	DB := setting.FileConfigSetting.MongoDBSetting.Name
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		setting.FileConfigSetting.MongoDBSetting.User,
		setting.FileConfigSetting.MongoDBSetting.Password,
		setting.FileConfigSetting.MongoDBSetting.Host,
		setting.FileConfigSetting.MongoDBSetting.Port)
	clientOptions := options.Client()
	// clientOptions.ApplyURI("mongodb://mongoadmin_dev:mongo_dev@34.101.133.247:1300")
	clientOptions.ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	return client.Database(DB), nil
}

// func NextSequence() int {
// 	db, err := Connect()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	db.Collection("sequence_logs").FindOneAndUpdate()
// 	_, err = db.Collection("sequence_logs").InsertOne(context.TODO(), sequence{"next", 0})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	return 0
// }
