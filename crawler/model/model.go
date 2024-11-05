package model

import "time"

type Url struct {
	ShortUrl string    `json:"shortUrl"`
	LongUrl  string    `json:"longUrl"`
	IsEnable int       `json:"isEnable"`
	RegDate  time.Time `json:"regDate"`
	UrlId    int       `json:"urlId" gorm:"primaryKey"`
}

func (Url) TableName() string {
	return "url"
}

type Contents struct {
	ContentId string    `json:"content_id"`
	ShortUrl  string    `json:"shortUrl"`
	Size      string    `json:"size"`
	Type      string    `json:"type"`
	Html      string    `json:"html"`
	Hash      string    `json:"hash"`
	IsEnable  int       `json:"isEnable"`
	RegDate   time.Time `json:"regDate"`
}

func (Contents) TableName() string {
	return "contents"
}
