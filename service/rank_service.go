package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"ranking/constant"
	"ranking/middleware/redisutil"
	"ranking/model/gift"
)

const FormatError = "failed to get rank : %v"

// GetRank 获取排行榜
func GetRank() ([]redis.Z, error) {
	redisClient := redisutil.GetRedisClient()
	defer redisutil.Close(redisClient)
	scores := redisClient.ZRevRangeByScoreWithScores(constant.RankRedisKey+"1", redis.ZRangeBy{Min: "-inf", Max: "+inf"})
	result, err := scores.Result()
	if err != nil {
		err := fmt.Errorf(FormatError, err)
		return nil, err
	}
	// 如果key不存在
	if len(result) == 0 {
		res, err := gift.GetGroupedGiftValue()
		if err != nil {
			err := fmt.Errorf(FormatError, err)
			return nil, err
		}
		redisData := make([]redis.Z, 0, cap(res))
		for _, v := range res {
			redisData = append(redisData, redis.Z{Score: float64(v["Score"].(int)), Member: v["_id"]})
		}
		redisClient.ZAdd(constant.RankRedisKey, redisData...)
		redisClient.Expire(constant.RankRedisKey, constant.HalfHour)
		return redisData, nil
	} else {
		return result, nil
	}
}
