package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	"context"
)

var ctx = context.Background()

var RedisClient *redis.Client // Redis 연결 관리

func InitRedis() { // Redis 클라이언트 초기화
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 서버 주소
		Password: "password",       // 비밀번호 (없으면 빈 문자열)
		DB:       0,                // 사용할 DB 번호 (기본: 0)
	})

	// Redis 연결 테스트
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}

func SetRedis(key string, value string, expiration time.Duration) error {
	InitRedis()
	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %v", key, err)
	}
	return nil
}

func GetRedis(key string) (string, error) {
	InitRedis()
	value, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' does not exist: %v", key, err)
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch redis key: %v", err)
	}

	return value, nil
}
