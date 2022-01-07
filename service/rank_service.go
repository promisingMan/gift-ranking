package service

import (
	"github.com/go-redis/redis"
	"ranking/constant"
	"ranking/middleware/redisutil"
	"ranking/model"
	"time"
)

// GetRank 获取排行榜
func GetRank() ([]redis.Z, error) {
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
		return redisData, nil
	} else {
		scores := redisClient.ZRangeByScoreWithScores(constant.RankRedisKey, redis.ZRangeBy{Min: "-inf", Max: "+inf"})
		result, err := scores.Result()
		if err != nil {
			return nil, err
		}
		redisClient.Expire(constant.RankRedisKey, time.Minute*30)
		// 将切片倒序
		reverse(result)
		return result, nil
	}
}

func reverse(result []redis.Z) {
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
}
