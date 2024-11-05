package main

import (
	"context"
	"crawler/config"
	"crawler/config/cmd/crawler"
	"fmt"
	"log"
	"sync"
)

func main() {
	rdb := config.NonStoredUrlCache()
	ctx := context.Background()
	var wg sync.WaitGroup
	maxRoutine := 3

	// 채널 구독
	pubsub := rdb.Subscribe(ctx, "urls")
	defer pubsub.Close()

	fmt.Println("Subscribed to channel. Waiting for messages...")

	ch := make(chan struct{}, maxRoutine)
	// 메시지 대기
	for {
		wg.Add(1)
		ch <- struct{}{}
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		// 메시지가 들어오면 processMessage 함수 호출
		fmt.Printf("Received message: %s\n", msg.Payload)
		go func() {
			defer wg.Done()
			defer func() { <-ch }()
			crawler.StartCrawler(msg.Payload)
		}()

		wg.Wait()

	}

}
