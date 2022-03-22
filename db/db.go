package db

import (
	"context"
	"fmt"
	"frdocker/db/drivers"
	"frdocker/settings"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TrainerModel
type Mgo struct {
	collection *mongo.Collection
}

func NewMgo(collectionName string) *Mgo {
	mgo := new(Mgo)
	mgo.collection = drivers.GetMongoClient().Database(settings.DB_NAME).Collection(collectionName)
	return mgo
}

// 插入单个
func (m *Mgo) InsertOne(document interface{}) (insertResult *mongo.InsertOneResult) {
	insertResult, err := m.collection.InsertOne(context.TODO(), document)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// 插入多个
func (m *Mgo) InsertMany(documents []interface{}) (insertManyResult *mongo.InsertManyResult) {
	insertManyResult, err := m.collection.InsertMany(context.TODO(), documents)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// 查询单个
func (m *Mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	filter := bson.D{{Key: key, Value: value}}
	singleResult := m.collection.FindOne(context.TODO(), filter)
	if singleResult != nil {
		fmt.Println(singleResult)
	}
	return singleResult
}

// 查询count总数
func (m *Mgo) Count() (string, int64) {
	name := m.collection.Name()
	size, _ := m.collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

// 按选项查询集合
// Skip 跳过
// Limit 读取数量
// Sort  排序   1 倒叙 ， -1 正序
func (m *Mgo) FindAll(Skip, Limit int64, sort int) *mongo.Cursor {
	SORT := bson.D{{Key: "_id", Value: sort}}
	filter := bson.D{{}}

	// where
	findOptions := options.Find()
	findOptions.SetSort(SORT)
	findOptions.SetLimit(Limit)
	findOptions.SetSkip(Skip)

	cur, err := m.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Println(err)
	}

	return cur
}

// 获取集合创建时间和编号
func (m *Mgo) ParsingId(result string) (time.Time, uint64) {
	temp1 := result[:8]
	timestamp, _ := strconv.ParseInt(temp1, 16, 64)
	dateTime := time.Unix(timestamp, 0) // 这是截获情报时间 时间格式 2019-04-24 09:23:39 +0800 CST
	temp2 := result[18:]
	count, _ := strconv.ParseUint(temp2, 16, 64) // 截获情报的编号
	return dateTime, count
}

// 删除
func (m *Mgo) Delete(key string, value interface{}) int64 {
	filter := bson.D{{Key: key, Value: value}}
	count, err := m.collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount

}

// 删除多个
func (m *Mgo) DeleteMany(key string, value interface{}) int64 {
	filter := bson.D{{Key: key, Value: value}}
	count, err := m.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount
}

// 更新一个
func (m *Mgo) UpdateOne(key string, value interface{}, update interface{}) (updateResult *mongo.UpdateResult) {
	filter := bson.D{{Key: key, Value: value}}
	updateResult, err := m.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return updateResult
}
