package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"golang.org/x/sync/semaphore"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

// GiftRecRecord 送礼流水记录结构体
type GiftRecRecord struct {
	AnchorId   int       `json:"anchorId" bson:"anchorId"`     // 主播id
	Uid        int       `json:"uid" bson:"uid"`               // 用户id
	GiftValue  int       `json:"giftValue" bson:"giftValue"`   // 礼物价值
	CreateTime time.Time `json:"createTime" bson:"createTime"` // 送礼时间
}

// redis
var rdb *redis.Client

// mongo
var mongoSession *mgo.Session

// 信号量控制mongo数据库并发访问的协程数
var sema = semaphore.NewWeighted(10)

// RedisKey 常量
var RedisKey = "rank"

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
	saveToMongoDB(giftRecRecord)

	// 更新redis排行榜缓存
	rdb.ZIncrBy("rank", float64(giftRecRecord.GiftValue), strconv.Itoa(giftRecRecord.AnchorId))
}

func saveToMongoDB(giftRecRecord GiftRecRecord) {
	// 信号量控制mongo并发访问数
	err := sema.Acquire(context.TODO(), 1)
	if err != nil {
		panic(err)
	}
	defer sema.Release(1)

	collection := mongoSession.DB("test").C("gift_rec_record")
	err = collection.Insert(giftRecRecord)
	if err != nil {
		panic(err)
	}
}

// Rank 获取排行榜
func Rank(w http.ResponseWriter, req *http.Request) {
	if handleDataNotInCache(w) {
		return
	}
	scores := rdb.ZRangeByScoreWithScores(RedisKey, redis.ZRangeBy{Min: "-inf", Max: "+inf"})
	result, err := scores.Result()
	if err != nil {
		panic(err)
	}
	// 将切片倒序
	reverse(result)
	resp, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	respJson(w, resp)
}

func respJson(w http.ResponseWriter, resp []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(resp)
	if err != nil {
		panic(err)
	}
}

func reverse(result []redis.Z) {
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
}

func handleDataNotInCache(w http.ResponseWriter) bool {
	//rdb.Del(RedisKey)
	exists := rdb.Exists(RedisKey)
	if exists.Val() == 0 {
		// 只有一个协程能够set成功，防止并发操作mongo写缓存
		lock := rdb.SetNX("rank_nx", 1, time.Second)
		if lock.Val() {
			err := sema.Acquire(context.TODO(), 1)
			if err != nil {
				panic(err)
			}
			defer sema.Release(1)
			collection := mongoSession.DB("test").C("gift_rec_record")
			pipe := collection.Pipe([]bson.M{
				{"$group": bson.M{"_id": "$anchorId", "Score": bson.M{"$sum": "$giftValue"}}},
				{"$sort": bson.M{"Score": -1}},
			})
			var res []bson.M
			err = pipe.All(&res)
			if err != nil {
				panic(err)
			}
			resp := make([]redis.Z, 0, cap(res))
			for _, v := range res {
				resp = append(resp, redis.Z{Score: float64(v["Score"].(int)), Member: v["_id"]})
			}
			rdb.ZAdd(RedisKey, resp...)
			rdb.Expire(RedisKey, time.Minute*30)
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}
			respJson(w, jsonResp)
			return true
		}
	}
	return false
}

// GiftRecRecordList 收礼流水查询
func GiftRecRecordList(w http.ResponseWriter, req *http.Request) {
	// 处理查询参数
	anchorId, page, limit := handleArgs(req)

	// 信号量控制mongo并发访问数
	err := sema.Acquire(context.TODO(), 1)
	if err != nil {
		panic(err)
	}
	defer sema.Release(1)

	collection := mongoSession.DB("test").C("gift_rec_record")
	iter := collection.Find(bson.M{"anchorId": anchorId}).Sort("createTime").Skip((page - 1) * limit).Limit(limit).Iter()
	var result []GiftRecRecord
	err = iter.All(&result)
	if err != nil {
		panic(err)
	}
	resp, err := json.Marshal(result)
	respJson(w, resp)
}

func handleArgs(req *http.Request) (int, int, int) {
	values := req.URL.Query()
	var anchorIdStr = values.Get("anchorId")
	anchorId, err := strconv.Atoi(anchorIdStr)
	if err != nil {
		panic("anchorId解析失败")
	}
	var pageStr = values.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		panic("page解析失败")
	}
	var limitStr = values.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		panic("limit解析失败")
	}
	return anchorId, page, limit
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var err error
	mongoSession, err = mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/giveGift", GiveGift)
	http.HandleFunc("/rank", Rank)
	http.HandleFunc("/giftRecRecordList", GiftRecRecordList)
	http.ListenAndServe(":8090", nil)
}
