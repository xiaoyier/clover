package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrorVoteExpire = errors.New("vote expired")
var ErrorVoted = errors.New("already voted")

func Vote(postId int64, value float64) error {

	id := strconv.FormatInt(postId, 10)
	t := redisClient.ZScore(redisClient.Context(), RedisKeyPostTime, id).Val()
	if (time.Now().UnixNano() - int64(t)) > int64(VoteExpireTime) {
		return ErrorVoteExpire
	}

	currentValue := redisClient.ZScore(redisClient.Context(), RedisKeyPostVotedPrefix+id, id).Val()
	diffAbs := math.Abs(value - currentValue)

	// update
	pipeline := redisClient.TxPipeline()
	pipeline.ZIncrBy(redisClient.Context(), RedisKeyPostScore, ScoreOneVote*diffAbs*value, id)

	pipeline.ZAdd(redisClient.Context(), RedisKeyPostVotedPrefix+id, &redis.Z{
		Member: id,
		Score:  value,
	})

	switch math.Abs(currentValue) - math.Abs(value) {
	case 1:
		// 取消投票 ov=1/-1 v=0
		// 投票数-1
		pipeline.HIncrBy(redisClient.Context(), RedisKeyPostInfoPrefix+id, "vote", -1)
	case 0:
		// 反转投票 ov=-1/1 v=1/-1
		// 投票数不用更新
	case -1:
		// 新增投票 ov=0 v=1/-1
		// 投票数+1
		pipeline.HIncrBy(redisClient.Context(), RedisKeyPostInfoPrefix+id, "votes", 1)
	default:
		// 已经投过票了
		return ErrorVoted
	}

	_, err := pipeline.Exec(redisClient.Context())
	return err
}
