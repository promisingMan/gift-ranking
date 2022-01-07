package service

import (
	"github.com/go-redis/redis"
	"ranking/constant"
	"ranking/middleware/redisutil"
	"ranking/model"
	"strconv"
	"time"
)

func GiveGift(giftRecDto model.GiftRecDto) {
	// 转换成mongo实体并入库
	record := model.GiftRecRecord{
		AnchorId:   giftRecDto.AnchorId,
		Uid:        giftRecDto.Uid,
		GiftValue:  giftRecDto.GiftValue,
		CreateTime: time.Now(),
	}
	model.SaveGiftRecRecord(record)

	// 更新redis排行榜缓存
	redisClient := redisutil.GetRedisClient()
	defer redisutil.Close(redisClient)
	exists := redisClient.Exists(constant.RankRedisKey)
	if exists.Val() == 0 {
		res := model.GetGroupedGiftValue()
		redisData := make([]redis.Z, 0, cap(res))
		for _, v := range res {
			redisData = append(redisData, redis.Z{Score: float64(v["Score"].(int)), Member: v["_id"]})
		}
		redisClient.ZAdd(constant.RankRedisKey, redisData...)
		redisClient.Expire(constant.RankRedisKey, time.Minute*30)
	} else {
		redisClient.ZIncrBy(constant.RankRedisKey, float64(giftRecDto.GiftValue), strconv.Itoa(giftRecDto.AnchorId))
		redisClient.Expire(constant.RankRedisKey, time.Minute*30)
	}
}
