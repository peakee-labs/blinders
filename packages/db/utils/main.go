package dbutils

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// username:password@host:port/database
const MongoURLTemplate = "mongodb://%s:%s@%s:%s/%s"

func InitMongoClient(url string) (*mongo.Client, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*10)
	defer cal()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to mongo, err: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("cannot ping to mongo, err: %v", err)
	}
	return client, nil
}

type MongoInfo struct {
	URL    string
	DBName string
}

func GetMongoInfoFromEnvironment() MongoInfo {
	url := fmt.Sprintf(
		MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)
	dbName := os.Getenv("MONGO_DATABASE")

	return MongoInfo{
		URL:    url,
		DBName: dbName,
	}
}
