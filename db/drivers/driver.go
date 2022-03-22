package drivers

import (
	"context"
	"fmt"
	"frdocker/settings"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MgoClient *mongo.Client

func GetMongoClient() *mongo.Client {
	if MgoClient == nil {
		MgoClient = Connect()
	}
	return MgoClient
}

// 连接
func Connect() *mongo.Client {

	// 设置客户端参数
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", settings.DB_HOST, settings.DB_PORT))
	clientOptions.SetAuth(options.Credential{
		AuthSource: settings.DB_NAME,
		Username:   settings.DB_USER,
		Password:   settings.DB_PASS,
	})

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//defer client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// 检查链接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

// 关闭
func Close() {
	if MgoClient == nil {
		return
	}

	err := MgoClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
