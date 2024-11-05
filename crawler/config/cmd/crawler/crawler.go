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
	rdb := config.Cache()
	ctx := context.Background()
	hash := internal.GetHtmlHash(html)

	fmt.Println("==================================")
	fmt.Println(hash)

	// 채널 구독
	for _, htmlUrl := range urls {
		err := rdb.Publish(ctx, "urls", htmlUrl).Err()
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
		fmt.Println(urls)
	}

}
