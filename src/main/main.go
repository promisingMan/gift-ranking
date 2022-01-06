package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GiftRecRecord 送礼流水记录结构体
type GiftRecRecord struct {
	AnchorId   int       `json:"anchorId" bson:"anchorId"` // 主播id
	Uid        int       `json:"uid"`                      // 用户id
	GiftValue  int       `json:"giftValue"`                // 礼物价值
	CreateTime time.Time `json:"createTime"`               // 送礼时间
}

var mdb = ConnectToDB("mongodb://localhost:27017", "test", 10, 10)
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var sema = semaphore.NewWeighted(10)

// GiveGift 送礼
func GiveGift(w http.ResponseWriter, req *http.Request) {
	// 获取json入参，并更新送礼时间为当前时间
	decoder := json.NewDecoder(req.Body)
	var giftRecRecord GiftRecRecord
	err := decoder.Decode(&giftRecRecord)
	if err != nil {
		panic(err)
	}
	giftRecRecord.CreateTime = time.Now()

	// 入库mongo
	err = sema.Acquire(context.TODO(), 1)
	if err != nil {
		SaveToMongoDB(err, giftRecRecord)
		sema.Release(1)
	}

	rdb.ZIncrBy("rank", float64(giftRecRecord.GiftValue), strconv.Itoa(giftRecRecord.AnchorId))
}

func SaveToMongoDB(err error, giftRecRecord GiftRecRecord) {
	collection := mdb.Collection("gift_rec_record")
	res, err := collection.InsertOne(context.TODO(), giftRecRecord)
	if err != nil {
		panic(err)
	}
	log.Printf("%v", res)
}

func ConnectToDB(uri, name string, timeout time.Duration, num uint64) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(num)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil
	}
	return client.Database(name)
}

// Rank 获取排行榜
func Rank(w http.ResponseWriter, req *http.Request) {
	scores := rdb.ZRangeByScoreWithScores("rank", redis.ZRangeBy{})
	log.Println(scores)
}

// GiftRecRecordList 收礼流水查询
func GiftRecRecordList(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	var anchorId = values.Get("anchorId")
	anchorIdInt, err := strconv.ParseInt(anchorId, 10, 64)
	if err != nil {
		panic("解析失败")
	}
	collection := mdb.Collection("gift_rec_record")
	var findOption *options.FindOptions
	if limit > 0 {
		findOption.SetLimit(limit)
		findOption.SetSkip(limit * index)
	}
	collection.Find(context.TODO(), bson.D{{"anchorId", anchorId}}, findOption)
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	http.HandleFunc("/giveGift", GiveGift)
	http.HandleFunc("/rank", Rank)
	http.HandleFunc("/giftRecRecordList", GiftRecRecordList)
	http.ListenAndServe(":8090", nil)
}
