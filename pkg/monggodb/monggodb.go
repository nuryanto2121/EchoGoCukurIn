package monggodb

import (
	"context"
	"fmt"
	"log"
	"nuryanto2121/dynamic_rest_api_go/pkg/setting"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var CTX, _ = context.WithTimeout(context.Background(), 60*time.Second)
var MCon *mongo.Database

type sequence struct {
	Name          string `bson:"name"`
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
	fmt.Printf(connectionString)
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

func Setup() {
	DB := setting.FileConfigSetting.MongoDBSetting.Name
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		setting.FileConfigSetting.MongoDBSetting.User,
		setting.FileConfigSetting.MongoDBSetting.Password,
		setting.FileConfigSetting.MongoDBSetting.Host,
		setting.FileConfigSetting.MongoDBSetting.Port)
	// clientOptions := options.Client()
	fmt.Println(connectionString)
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// err = client.Disconnect(context.TODO())

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connection to MongoDB closed.")

	MCon = client.Database(DB)
}
func NextSequence() int {
	// db,  := MCon
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// names, err := db.ListCollectionNames(context.TODO(),{})
	var result sequence
	var filter sequence
	filter.Name = "next"
	// filter := bson.D{{"name", "next"}}
	_ = MCon.Collection("sequence_logs").FindOne(context.TODO(), filter).Decode(&result)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	if result.SequenceValue == 0 { //= +1 //= result.SequenceValue + 1
		filter.SequenceValue = 1
		_, err := MCon.Collection("sequence_logs").InsertOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err.Error())
		}
		return 1
	}
	result.SequenceValue += 1
	updateResult, err := MCon.Collection("sequence_logs").UpdateOne(context.TODO(), filter, result) //db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return result.SequenceValue
}
