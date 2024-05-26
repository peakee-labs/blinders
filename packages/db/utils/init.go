package dbutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// username:password@host:port/database
const MongoURLTemplate = "mongodb://%s:%s@%s:%s/%s"

func InitMongoClient(url string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to mongo, err: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("cannot ping to mongo, err: %v", err)
	}

	return client, nil
}

func InitMongoDatabaseFromEnv(prefix ...string) (*mongo.Database, error) {
	info := GetMongoInfoFromEnv(prefix...)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(info.URL))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to mongo, err: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("cannot ping to mongo, err: %v", err)
	}

	return client.Database(info.DBName), nil
}

// return a Database channel, leverage goroutines to optimize aws lambda cold-start
// temporarily use log.Fatalf if error happen
func InitMongoDatabaseChanFromEnv(prefix ...string) chan *mongo.Database {
	ch := make(chan *mongo.Database)
	go func() {
		DB, err := InitMongoDatabaseFromEnv(prefix...)
		if err != nil {
			log.Fatal(err)
		}

		ch <- DB
	}()

	return ch
}

type MongoInfo struct {
	URL    string
	DBName string
}

func GetMongoInfoFromEnv(prefix ...string) MongoInfo {
	var url string
	var dbName string
	var pf string
	if len(prefix) > 0 && prefix[0] != "" {
		pf = prefix[0] + "_"
	}

	url = os.Getenv(pf + "MONGO_DATABASE_URL")
	dbName = os.Getenv(pf + "MONGO_DATABASE")

	if url == "" {
		url = fmt.Sprintf(
			MongoURLTemplate,
			os.Getenv(pf+"MONGO_USERNAME"),
			os.Getenv(pf+"MONGO_PASSWORD"),
			os.Getenv(pf+"MONGO_HOST"),
			os.Getenv(pf+"MONGO_PORT"),
			dbName,
		)
	}

	return MongoInfo{
		URL:    url,
		DBName: dbName,
	}
}
