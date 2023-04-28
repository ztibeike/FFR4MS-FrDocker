package db

import (
	"context"
	"fmt"

	"gitee.com/zengtao321/frdocker/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB() (*mongo.Database, error) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(getMongoURL()))
	if err != nil {
		return nil, err
	}
	return client.Database(config.MONGO_DB), nil
}

func CloseMongo(db *mongo.Database) {
	ctx := context.TODO()
	err := db.Client().Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func getMongoURL() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		config.MONGO_USER,
		config.MONGO_PASS,
		config.MONGO_HOST,
		config.MONGO_PORT,
		config.MONGO_DB,
	)
}
