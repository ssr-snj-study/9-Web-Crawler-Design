package config

import (
	"context"
	"crawler/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB
var cache *redis.Client

var e error

func init() {
	databaseInit()
	cacheInit()
	setRedisKey()
}

func databaseInit() {
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbName := os.Getenv("DB_NAME")
	//port := os.Getenv("DB_PORT")
	host := "127.0.0.1"
	user := "snj"
	password := "snj"
	dbName := "snj_db"
	port := 5432

	connectInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(connectInfo), &gorm.Config{})

	if e != nil {
		panic(e)
	}
}

func cacheInit() {
	//connectInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	connectInfo := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
	cache = redis.NewClient(&redis.Options{
		Addr: connectInfo, // Redis 서버 주소
		//Password: os.Getenv("REDIS_PASSWORD"), // 비밀번호가 없다면 빈 문자열
		Password: "snj", // 비밀번호가 없다면 빈 문자열
	})

}

func DB() *gorm.DB {
	return database
}

func Cache() *redis.Client {
	return cache
}

func setRedisKey() {
	key := "counter"

	rdb := Cache()
	ctx := context.Background()

	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to check key existence: %v", err)
	}
	if exists == 0 {
		url := &model.Url{}
		db := DB()

		var maxId int
		db.Model(url).Select("MAX(url_id)").Scan(&maxId)

		err = rdb.Set(ctx, key, maxId, 0).Err()
		if err != nil {
			log.Fatalf("Failed to set initial value: %v", err)
		}
	}

}
