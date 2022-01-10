package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"ranking/config"
	"ranking/middleware/mongodb"
	"time"
)

// GiftRecRecord 送礼流水记录mongo实体
type GiftRecRecord struct {
	AnchorId   int       `json:"anchorId" bson:"anchorId"`     // 主播id
	Uid        int       `json:"uid" bson:"uid"`               // 用户id
	GiftValue  int       `json:"giftValue" bson:"giftValue"`   // 礼物价值
	CreateTime time.Time `json:"createTime" bson:"createTime"` // 送礼时间
}

// GiftRecDto 送礼入参接收实体
type GiftRecDto struct {
	AnchorId  int `json:"anchorId"`  // 主播id
	Uid       int `json:"uid"`       // 用户id
	GiftValue int `json:"giftValue"` // 礼物价值
}

const COLLECTION = "gift_rec_record"

func SaveGiftRecRecord(record GiftRecRecord) error {
	session := mongodb.GetMongodbSession()
	defer mongodb.CloseMongodbSession(session)
	collection := session.DB(config.AppConfig.MongoDb.Database).C(COLLECTION)
	err := collection.Insert(record)
	if err != nil {
		log.Println("failed to save gift receive record failed", err)
		return err
	}
	return nil
}

var rateLimit = make(chan struct{}, 10)

func GetGroupedGiftValue() ([]bson.M, error) {
	rateLimit <- struct{}{}
	session := mongodb.GetMongodbSession()
	defer func() {
		mongodb.CloseMongodbSession(session)
		<-rateLimit
	}()
	collection := session.DB(config.AppConfig.MongoDb.Database).C(COLLECTION)
	pipe := collection.Pipe([]bson.M{
		{"$group": bson.M{"_id": "$anchorId", "Score": bson.M{"$sum": "$giftValue"}}},
		{"$sort": bson.M{"Score": -1}},
	})
	var res []bson.M
	err := pipe.All(&res)
	if err != nil {
		log.Println("failed to get grouped gift value ", err)
	}
	return res, err
}

func GetGiftRecRecordListByAnchorId(anchorId, page, limit int) ([]GiftRecRecord, error) {
	session := mongodb.GetMongodbSession()
	defer mongodb.CloseMongodbSession(session)
	collection := session.DB(config.AppConfig.MongoDb.Database).C(COLLECTION)
	iter := collection.Find(bson.M{"anchorId": anchorId}).Sort("-createTime").Skip((page - 1) * limit).Limit(limit).Iter()
	var result []GiftRecRecord
	err := iter.All(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get gift receive record list by anchorId : %v", err)
	}
	return result, nil
}
