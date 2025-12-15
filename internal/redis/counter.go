package redis

import "github.com/go-redis/redis/v8"

func IncrementCounter(rdb *redis.Client, counterType string, id string) (int64, error) {
	key := counterType + ":" + id
	return rdb.Incr(ctx, key).Result()
}

func DecrementCounter(rdb *redis.Client, counterType string, id string) (int64, error) {
	key := counterType + ":" + id
	return rdb.Decr(ctx, key).Result()
}
