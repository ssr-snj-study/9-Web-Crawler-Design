package internal

import (
	"github.com/go-rod/rod"
	"log"
	"strings"
)

type UrlInfo struct {
	Page *rod.Page
}

func (UrlInfo *UrlInfo) GetInfoFromUrl(url string) {

	// Rod 브라우저를 시작하고 연결
	browser := rod.New().MustConnect()

	// 원하는 페이지로 이동
	UrlInfo.Page = browser.MustPage(url)
	UrlInfo.Page.MustWaitLoad() // 페이지가 완전히 로드될 때까지 대기
}

func (UrlInfo *UrlInfo) GetHtml() string {
	htmlContent, err := UrlInfo.Page.HTML()
	if err != nil {
		log.Fatalf("Failed to get HTML: %v", err)
	}
	return htmlContent
}

func (UrlInfo *UrlInfo) GetUrls() []string {
	var urls []string
	// 모든 <a> 태그에서 href 속성 추출
	elements := UrlInfo.Page.MustElements("a")
	var aTags []string

	for _, el := range elements {
		href, _ := el.Attribute("href")
		if href != nil {
			aTags = append(aTags, *href)
		}
	}

	// URL 리스트 출력
	for _, url := range aTags {
		if strings.HasPrefix(url, "https://") {
			urls = append(urls, url)
		}
	}
	return urls
}
