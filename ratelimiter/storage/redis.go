package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage() *redisStorage {
	redisAddr := viper.GetString("REDIS_ADDR")
	redisDb := viper.GetInt("REDIS_DB")
	redisPassword := viper.GetString("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       redisDb,
		Password: redisPassword,
	})

	return &redisStorage{
		client: client,
	}
}

func (rs *redisStorage) GetBlock(ctx context.Context, key string) (*time.Time, error) {
	redisKey := fmt.Sprintf("block-%s", key)

	value, err := rs.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	blockTime, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return nil, err
	}

	return &blockTime, nil
}

func (rs *redisStorage) Increment(ctx context.Context, key string, maxAccess int64) (bool, error) {
	redisKey := fmt.Sprintf("count-%s", key)
	now := time.Now()
	clearBefore := now.Add(-time.Second)

	pipe := rs.client.Pipeline()

	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(clearBefore.UnixMicro(), 10))
	count := pipe.ZCard(ctx, redisKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if count.Val() >= maxAccess {
		return false, nil
	}

	pipe = rs.client.Pipeline()

	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now.UnixMicro()),
		Member: now.Format(time.RFC3339Nano),
	})
	pipe.Expire(ctx, redisKey, time.Second)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (rs *redisStorage) AddBlock(ctx context.Context, key string, blockInMilliseconds int64) (*time.Time, error) {
	redisKey := fmt.Sprintf("block-%s", key)
	expiration := time.Duration(int64(time.Millisecond) * blockInMilliseconds)
	blockedUntil := time.Now().Add(expiration)

	err := rs.client.Set(ctx, redisKey, blockedUntil.Format(time.RFC3339Nano), expiration).Err()
	if err != nil {
		return nil, err
	}

	return &blockedUntil, nil
}
