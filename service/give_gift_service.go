package service

import (
	"github.com/go-redis/redis"
	"ranking/constant"
	"ranking/middleware/redisutil"
	"ranking/model"
	"strconv"
	"time"
)

func GiveGift(giftRecDto model.GiftRecDto) error {
	// 转换成mongo实体并入库
	record := model.GiftRecRecord{
		AnchorId:   giftRecDto.AnchorId,
		Uid:        giftRecDto.Uid,
		GiftValue:  giftRecDto.GiftValue,
		CreateTime: time.Now(),
	}
	err := model.SaveGiftRecRecord(record)
	if err != nil {
		return err
	}

	// 更新redis排行榜缓存，先判断key是否存在，如果key不存在就直接使用zincrby命令，会自动生成key，导致脏数据
	// 或者此处删除缓存，不做更新缓存处理（根据业务来定）
	redisClient := redisutil.GetRedisClient()
	defer redisutil.Close(redisClient)
	exists := redisClient.Exists(constant.RankRedisKey)
	if exists.Val() == 0 {
		res, err := model.GetGroupedGiftValue()
		if err != nil {
			return err
		}
		redisData := make([]redis.Z, 0, cap(res))
		for _, v := range res {
			redisData = append(redisData, redis.Z{Score: float64(v["Score"].(int)), Member: v["_id"]})
		}
		redisClient.ZAdd(constant.RankRedisKey, redisData...)
		redisClient.Expire(constant.RankRedisKey, constant.HalfHour)
	} else {
		redisClient.ZIncrBy(constant.RankRedisKey, float64(giftRecDto.GiftValue), strconv.Itoa(giftRecDto.AnchorId))
	}
	return nil
}
