package crawler

import (
	"crawler/config"
	"crawler/internal"
	"fmt"
	"net/url"
	"strings"
)

var dnsQueueList QueueList
var priQueueList QueueList

type QueueList struct {
	First  internal.StringQueue
	Second internal.StringQueue
	Third  internal.StringQueue
}

func init() {
	dnsQueueList = QueueList{}
	priQueueList = QueueList{}
}

func MakeUrlQueue() {
	cache := getCache()
	cache.SetDnsList()
	go func() {
		pubsub := cache.Subscribe("urls")
		for {
			fmt.Println("Finding Url in Cache")
			msg, err := pubsub.ReceiveMessage(config.Ctx)
			if err != nil {
				fmt.Println("err: ", err)
			} else {
				parsedURL, err := url.Parse(msg.Payload)
				if err != nil {
					fmt.Println("Error parsing URL:", err)
					return
				}
				host := parsedURL.Hostname()
				// 전면 큐 -> 우선순위 설정
				switch getRootDomain(host) {
				case "naver.com":
					priQueueList.First.Enqueue(msg.Payload)
				default:
					priQueueList.Second.Enqueue(msg.Payload)
				}
			}
		}
	}()

	go func() {
		for {
			for i := 0; i < 3; i++ {
				value, ok := priQueueList.First.Dequeue()
				if ok {
					makeDnsQueue(value)
				}
			}
			for i := 0; i < 2; i++ {
				value, ok := priQueueList.Second.Dequeue()
				if ok {
					makeDnsQueue(value)
				}
			}
			for i := 0; i < 1; i++ {
				value, ok := priQueueList.Third.Dequeue()
				if ok {
					makeDnsQueue(value)
				}
			}
		}
	}()

}

func makeDnsQueue(rawUrl string) {
	cache := getCache()
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	host := parsedURL.Hostname()
	// dns 별 큐 생성 -> 후면 큐
	for {
		switch cache.FindQueue(host) {
		case "q1":
			dnsQueueList.First.Enqueue(rawUrl)
		case "q2":
			dnsQueueList.Second.Enqueue(rawUrl)
		case "q3":
			dnsQueueList.Third.Enqueue(rawUrl)
		}
	}
}

func getRootDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".") // 마지막 두 부분 조합
	}
	return host
}
