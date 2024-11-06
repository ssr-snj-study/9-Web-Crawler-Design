package model

import "time"

type Url struct {
	LongUrl  string    `json:"longUrl"`
	Ip       string    `json:"ip"`
	IsEnable int       `json:"isEnable"`
	RegDate  time.Time `json:"regDate"`
	UrlId    int       `json:"urlId" gorm:"primaryKey"`
}

func (Url) TableName() string {
	return "url"
}

type Contents struct {
	ContentId int       `json:"content_id"`
	UrlId     int       `json:"urlId"`
	LongUrl   string    `json:"longUrl"`
	Size      int       `json:"size"`
	Type      string    `json:"type"`
	Html      string    `json:"html"`
	Hash      string    `json:"hash"`
	IsEnable  int       `json:"isEnable"`
	RegDate   time.Time `json:"regDate"`
}

func (Contents) TableName() string {
	return "contents"
}
