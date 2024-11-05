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
	rdb := config.Cache()
	ctx := context.Background()

	// 채널 구독
	pubsub := rdb.Subscribe(ctx, "urls")
	defer pubsub.Close()

	fmt.Println("Subscribed to channel. Waiting for messages...")
	var wg sync.WaitGroup
	// 메시지 대기
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		// 메시지가 들어오면 processMessage 함수 호출
		fmt.Printf("Received message: %s\n", msg.Payload)
		wg.Add(1)
		go crawler.StartCrawler(msg.Payload)
	}

}
