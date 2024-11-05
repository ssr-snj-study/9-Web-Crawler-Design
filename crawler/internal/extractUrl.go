package internal

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"strings"
)

func GetUrlsFromUrl(url string) ([]string, string) {
	var urls []string
	// Rod 브라우저를 시작하고 연결
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// 원하는 페이지로 이동
	page := browser.MustPage(url)
	page.MustWaitLoad() // 페이지가 완전히 로드될 때까지 대기

	htmlContent, err := page.HTML()
	if err != nil {
		log.Fatalf("Failed to get HTML: %v", err)
	}
	fmt.Println(htmlContent)
	// 모든 <a> 태그에서 href 속성 추출
	elements := page.MustElements("a")
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

	return urls, htmlContent
}
