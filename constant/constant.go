package constant

import "time"

// RankRedisKey 排行榜redis key, HalfHour 过期时间半小时
const (
	RankRedisKey = "rank"
	HalfHour     = time.Minute * 30
)
