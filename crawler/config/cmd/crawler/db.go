package crawler

import (
	"crawler/config"
	"crawler/internal"
	"crawler/model"
	"time"
)

func Create(htmlUrl string) error {
	db := config.DB()
	url := &model.Url{}

	if res := db.Where("long_url = ?", htmlUrl).Find(url); res.Error != nil {
		return res.Error
	}
	if url.UrlId == 0 {
		ShortUrl := internal.MakeShortUrl()
		url = &model.Url{
			ShortUrl: ShortUrl,
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
