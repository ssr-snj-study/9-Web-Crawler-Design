package crawler

import (
	"context"
	"crawler/config"
	"crawler/internal"
	"fmt"
	"log"
)

func StartCrawler(url string) {
	go func() {
		originIp, err := internal.GetIpFromUrl(url)
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println(originIp)
	}()

	urls, html := internal.GetUrlsFromUrl(url)
	nonStoredUrlCache := config.NonStoredUrlCache()
	storedUrlCache := config.StoredUrlCache()
	ctx := context.Background()
	hash := internal.GetHtmlHash(html)

	fmt.Println("==================================")
	fmt.Println(hash)

	// 채널 구독
	for _, htmlUrl := range urls {
		if exists, _ := storedUrlCache.Exists(ctx, htmlUrl).Result(); exists > 0 {
			continue
		}
		err := nonStoredUrlCache.Publish(ctx, "urls", htmlUrl).Err()
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
		err = storedUrlCache.Set(ctx, htmlUrl, "1", 0).Err()
		if err != nil {
			fmt.Println("Error setting value:", err)
		}
		fmt.Println(htmlUrl)
		go func() {
			_ = Create(htmlUrl)
		}()
	}
}
