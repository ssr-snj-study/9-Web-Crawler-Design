package crawler

import (
	"crawler/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var stCache *Cache

type Cache struct {
	NonStoredUrlCache *redis.Client
	StoredUrlCache    *redis.Client
	DnsCache          *redis.Client
}

func init() {
	urlInit()
}

func getCache() *Cache {
	return stCache
}

func urlInit() {
	cache := &Cache{}
	//connectInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	connectInfo := fmt.Sprintf("%s:%d", "127.0.0.1", 6380)
	cache.StoredUrlCache = redis.NewClient(&redis.Options{
		Addr: connectInfo, // Redis 서버 주소
		//Password: os.Getenv("REDIS_PASSWORD"), // 비밀번호가 없다면 빈 문자열
		Password: "snj", // 비밀번호가 없다면 빈 문자열
	})

	connectInfo = fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
	cache.NonStoredUrlCache = redis.NewClient(&redis.Options{
		Addr: connectInfo, // Redis 서버 주소
		//Password: os.Getenv("REDIS_PASSWORD"), // 비밀번호가 없다면 빈 문자열
		Password: "snj", // 비밀번호가 없다면 빈 문자열
	})

	connectInfo = fmt.Sprintf("%s:%d", "127.0.0.1", 6381)
	cache.DnsCache = redis.NewClient(&redis.Options{
		Addr: connectInfo, // Redis 서버 주소
		//Password: os.Getenv("REDIS_PASSWORD"), // 비밀번호가 없다면 빈 문자열
		Password: "snj", // 비밀번호가 없다면 빈 문자열
	})

	stCache = cache
}

func (c *Cache) alreadyCheckUrl(url string) bool {
	storedUrlCache := c.StoredUrlCache
	if exists, _ := storedUrlCache.Exists(config.Ctx, url).Result(); exists > 0 {
		return true
	}
	return false
}

func (c *Cache) Subscribe(channel string) *redis.PubSub {
	return c.NonStoredUrlCache.Subscribe(config.Ctx, channel)
}

func (c *Cache) Publish(channel, value string) {
	err := c.NonStoredUrlCache.Publish(config.Ctx, channel, value).Err()
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}

func (c *Cache) StoreAlreadyCheckUrl(url string) {
	err := c.StoredUrlCache.Set(config.Ctx, url, "1", 0).Err()
	if err != nil {
		fmt.Println("Error setting value:", err)
	}
}
