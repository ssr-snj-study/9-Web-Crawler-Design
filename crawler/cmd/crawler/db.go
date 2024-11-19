package crawler

import (
	"crawler/config"
	"crawler/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func InsertUrl(htmlUrl string) error {
	db := config.DB()
	url := &model.Url{}

	if res := db.Where("long_url = ?", htmlUrl).Find(url); res.Error != nil {
		return res.Error
	}
	if url.UrlId == 0 {
		url = &model.Url{
			LongUrl:  htmlUrl,
			IsEnable: 1,
			RegDate:  time.Now(),
		}
		if err := db.Create(&url).Error; err != nil {
			return err
		}
	}

	return nil
}

func IsDuplicateHtml(hash string) bool {
	db := config.DB()
	contents := &model.Contents{}

	result := db.Where("hash = ?", hash).First(&contents)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			fmt.Println("select Error:", result.Error)
		}
	}
	return true
}

func InsertContent(url, hash, htmlUrl string) error {
	db := config.DB()
	// 존재하지 않으면 Insert
	contents := &model.Contents{LongUrl: url, Hash: hash, IsEnable: 1, RegDate: time.Now()}
	if err := db.Create(contents).Error; err != nil {
		fmt.Println("Error inserting record:", err)
	} else {
		fmt.Println("Record inserted:", contents)
		return err
	}
	return nil
}
