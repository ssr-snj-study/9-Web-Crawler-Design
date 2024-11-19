package main

import (
	"crawler/cmd/crawler"
)

func main() {

	// 미수집 저장소
	// 도메인별
	crawler.MakeUrlQueue()
	crawler.RunCrawlerMaxRoutine()
}
