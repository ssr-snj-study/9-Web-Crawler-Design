package crawler

import (
	"crawler/internal"
	"fmt"
	"sync"
	"time"
)

func RunCrawlerMaxRoutine() {
	var wg sync.WaitGroup
	maxRoutine := 3
	ch := make(chan struct{}, maxRoutine)

	// 메시지 대기
	for {
		time.Sleep(1 * time.Second)
		wg.Add(1)
		fmt.Println("Subscribed to channel. Waiting for messages...")
		for i := 0; i < 3; i++ {
			fmt.Println("start dns first queue")
			value, ok := dnsQueueList.First.Dequeue()
			if ok {
				fmt.Println("start dns first ok")
				wg.Add(1)
				go func() {
					defer wg.Done()
					ch <- struct{}{}
					defer func() { <-ch }()
					StartCrawler(value)
				}()
			}
		}
		for i := 0; i < 2; i++ {
			value, ok := dnsQueueList.Second.Dequeue()
			fmt.Println("start dns Second queue")
			if ok {
				fmt.Println("start dns Second ok")
				wg.Add(1)
				go func() {
					defer wg.Done()
					ch <- struct{}{}
					defer func() { <-ch }()
					StartCrawler(value)
				}()
			}
		}
		for i := 0; i < 1; i++ {
			value, ok := dnsQueueList.Third.Dequeue()
			fmt.Println("start dns Third queue")
			if ok {
				fmt.Println("start dns Third ok")
				wg.Add(1)
				go func() {
					defer wg.Done()
					ch <- struct{}{}
					defer func() { <-ch }()
					StartCrawler(value)
				}()
			}
		}
		wg.Done()
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
	SetContents(url, &page)

	// 채널 구독
	for _, htmlUrl := range page.GetUrls() {
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

func SetContents(url string, page *internal.UrlInfo) {
	html := page.GetHtml()
	hash := internal.GetHtmlHash(html)

	if !IsDuplicateHtml(hash) {
		go func() {
			_ = InsertContent(url, hash, html)
		}()
	}
}
