package redis

import (
	"clover/model/mysql"
	"clover/pkg/log"
	"time"

	"github.com/go-redis/redis/v8"
)

type PostInfo struct {
	PostID        string
	AuthorID      string
	CommunityID   string
	AuthorName    string
	CommunityName string
	Title         string
	Summary       string
}

const (
	ScoreOneVote   = 432
	VoteExpireTime = time.Hour * 24 * 7
)

func CreatePost(info *PostInfo) error {

	now := time.Now().UnixNano()
	pipeline := redisClient.TxPipeline()
	pipeline.ZAdd(redisClient.Context(), RedisKeyPostTime, &redis.Z{
		Member: info.PostID,
		Score:  float64(now),
	})
	pipeline.ZAdd(redisClient.Context(), RedisKeyPostScore, &redis.Z{
		Member: info.PostID,
		Score:  float64(now) + ScoreOneVote,
	})
	pipeline.HMSet(redisClient.Context(), RedisKeyPostInfoPrefix+info.PostID, map[string]interface{}{
		"post:id":        info.PostID,
		"author:id":      info.AuthorID,
		"author_name":    info.AuthorName,
		"community_name": info.CommunityName,
		"vote":           1,
		"title":          info.Title,
		"summary":        info.Summary,
	})

	pipeline.ZAdd(redisClient.Context(), RedisKeyPostVotedPrefix+info.PostID, &redis.Z{
		Member: info.AuthorID,
		Score:  1,
	})
	pipeline.Expire(redisClient.Context(), RedisKeyPostVotedPrefix+info.PostID, VoteExpireTime)

	pipeline.SAdd(redisClient.Context(), RedisKeyCommunityPostPrefix+info.CommunityID, info.PostID)
	_, err := pipeline.Exec(redisClient.Context())
	return err
}

func GetPostDetail(id string) (map[string]string, error) {

	key := RedisKeyCommunityPostPrefix + id
	cmd := redisClient.HGetAll(redisClient.Context(), key)
	return cmd.Val(), cmd.Err()
}

func GetPostList(order string, page *mysql.PostListPage) []map[string]string {
	key := RedisKeyPostTime
	if order == "score" {
		key = RedisKeyPostScore
	}

	start := page.PageNumber * page.PageSize
	end := (page.PageNumber+1)*page.PageSize - 1
	ids := redisClient.ZRevRange(redisClient.Context(), key, int64(start), int64(end)).Val()

	result := make([]map[string]string, len(ids))
	for i, id := range ids {
		infoKey := RedisKeyPostInfoPrefix + id
		info := redisClient.HGetAll(redisClient.Context(), infoKey).Val()
		result[i] = info
	}

	return result
}

func GetPostListByKey(key string, page *mysql.PostListPage) []map[string]string {
	start := page.PageNumber * page.PageSize
	end := (page.PageNumber+1)*page.PageSize - 1
	ids := redisClient.ZRevRange(redisClient.Context(), key, int64(start), int64(end)).Val()

	result := make([]map[string]string, len(ids))
	for i, id := range ids {
		infoKey := RedisKeyPostInfoPrefix + id
		info := redisClient.HGetAll(redisClient.Context(), infoKey).Val()
		result[i] = info
	}

	return result
}

func GetCommunityPostList(order string, communityId string, page *mysql.PostListPage) []map[string]string {

	//redisClient.Z
	communityKey := RedisKeyCommunityPostPrefix + communityId
	orderKey := RedisKeyPostTime
	if order == "score" {
		orderKey = RedisKeyPostScore
	}
	//redisClient
	key := communityId + order
	if redisClient.Exists(redisClient.Context(), key).Val() < 1 {
		redisClient.ZInterStore(redisClient.Context(), key, &redis.ZStore{
			Keys:      []string{communityKey, orderKey},
			Aggregate: "SUM",
		})
	}
	redisClient.Expire(redisClient.Context(), key, time.Second*60)
	members := redisClient.SMembers(redisClient.Context(), key).Val()
	log.WithCategory("redis.post").Debugf("GetCommunityPostList: members is %v", members)

	return GetPostListByKey(key, page)
}
