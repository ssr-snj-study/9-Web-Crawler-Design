package config

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

var Ctx context.Context = context.Background()

func init() {
	databaseInit()
}

func databaseInit() {
	var e error
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
	//sqlSet, err := database.DB()
	//if err != nil {
	//	panic("failed to get database")
	//}
	//sqlSet.SetConnMaxLifetime(time.Hour)
	//sqlSet.SetMaxOpenConns(50)

}

func DB() *gorm.DB {
	return database
}
