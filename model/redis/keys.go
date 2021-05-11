package redis

const (
	RedisKeyPostInfoPrefix      = "clover:post:Info:"      //hash
	RedisKeyPostTime            = "clover:post:time"       //zset
	RedisKeyPostScore           = "clover:post:score"      //zset
	RedisKeyPostVotedPrefix     = "clover:post:vote:"      //zset
	RedisKeyCommunityPostPrefix = "clover:community:post:" //set
)
