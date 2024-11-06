package crawler

import (
	"crawler/config"
	"crawler/internal"
	"fmt"
	"log"
	"sync"
)

func RunCrawlerMaxRoutine() {
	var wg sync.WaitGroup
	maxRoutine := 3
	cache := getCache()

	pubsub := cache.Subscribe("urls")
	ch := make(chan struct{}, maxRoutine)

	// 메시지 대기
	for {
		wg.Add(1)
		ch <- struct{}{}
		fmt.Println("Subscribed to channel. Waiting for messages...")
		msg, err := pubsub.ReceiveMessage(config.Ctx)
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		// 메시지가 들어오면 processMessage 함수 호출
		fmt.Printf("Received message: %s\n", msg.Payload)
		go func() {
			defer wg.Done()
			defer func() { <-ch }()
			StartCrawler(msg.Payload)
		}()

		wg.Wait()

	}
}

func StartCrawler(url string) {
	cache := getCache()
	// DNS 세팅 go Routine
	go func() {
		originIp, err := GetIpFromUrl(url)
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println(originIp)
	}()
	var page internal.UrlInfo
	page.GetInfoFromUrl(url)

	html := page.GetHtml()
	hash := internal.GetHtmlHash(html)

	if CheckDuplicateHtml(hash) {
		go func() {
			_ = InsertContent(url, hash, html)
		}()
	}

	urls := page.GetUrls()
	// 채널 구독
	for _, htmlUrl := range urls {
		if cache.alreadyCheckUrl(htmlUrl) {
			continue
		}
		cache.Publish("urls", htmlUrl)
		cache.StoreAlreadyCheckUrl(htmlUrl)

		fmt.Println(htmlUrl)

		go func() {
			_ = InsertUrl(htmlUrl)
		}()
	}
}
