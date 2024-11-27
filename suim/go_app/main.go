package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go_app/crawler"
	"go_app/utils"
	"golang.org/x/net/html"
	"net"
	"slices"
	"strings"
	"time"
)

var msgQueue = make(map[string][]string)
var domainIPMap = make(map[string]string)

func createDomainTopic(domain string, domainPath string) {
	topic := crawler.DomainToCamel(domain)
	if _, exists := msgQueue[topic]; !exists {
		//msgQueue[topic] = make(map[string]struct{})
		msgQueue[topic] = []string{}
	}
	//msgQueue[topic][strings.Split(value, "#")[0]] = struct{}{}
	domainPath = strings.Split(domainPath, "#")[0]
	if domainPath == "" {
		domainPath = "/"
	} else if !strings.HasSuffix(domainPath, "/") {
		domainPath += "/"
	}

	if slices.Contains(msgQueue[topic], domainPath) {
		return
	}
	msgQueue[topic] = append(msgQueue[topic], domainPath)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(domainPath),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func startDownloadURL(c chan string) {
	link := <-c

	doc, err := utils.ParseHTML(link)
	if err != nil {
		fmt.Println(err)
		return
	}

	mainDomain, path := crawler.SplitDomain(link)
	createDomainTopic(mainDomain, path)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if strings.HasPrefix(href, "/") {
						createDomainTopic(mainDomain, href)
					} else if strings.HasPrefix(href, "http") {
						domain, path := crawler.SplitDomain(href)
						createDomainTopic(domain, path)
					} else {
						fmt.Println(href)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
}

func DomainTransfer(domain string) {
	ipv4, err := net.LookupIP(domain)
	if err != nil || len(ipv4) == 0 {
		fmt.Printf("failed to lookup IP for domain %s: %v\n", domain, err)
		return
	}
	ipStr := ipv4[0].String()

	domainIPMap[domain] = ipStr
}

// saveURL saves the URL to the message queue and updates the domain-IP map
func saveURL(domain, path string) {
	ipv4, err := net.LookupIP(domain)
	if err != nil || len(ipv4) == 0 {
		fmt.Printf("failed to lookup IP for domain %s: %v\n", domain, err)
		return
	}
	ipStr := ipv4[0].String()

	domainIPMap[domain] = ipStr
	if _, exists := msgQueue[ipStr]; !exists {
		//msgQueue[ipStr] = make(map[string]struct{})
		msgQueue[ipStr] = []string{}
	}
	//msgQueue[ipStr][domain+strings.Split(path, "#")[0]] = struct{}{}
}

func GetKafkaTopicsName() (topics []string) {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}
	defer admin.Close()

	meta, err := admin.GetMetadata(nil, false, 5000)
	if err != nil {
		panic(err)
	}

	for topic := range meta.Topics {
		if strings.HasPrefix(topic, "domain") {
			topics = append(topics, topic)
		}
	}

	return topics
}

func SubscribeKafkaTopics(c chan string, topic []string) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	defer consumer.Close()

	err = consumer.SubscribeTopics(topic, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			if strings.HasPrefix(*msg.TopicPartition.Topic, "domain") {
				c <- "https://" + crawler.CamelToDomain((*msg.TopicPartition.Topic)[6:]) + string(msg.Value)
			} else {
				c <- string(msg.Value)
			}
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func downloadAllURL(c chan string) {
	link := <-c

	isMatch, robotsTxt := crawler.CompareRobotsTxt(link)
	if isMatch || robotsTxt == "" {
		return
	}

	//_ = utils.SaveFile("contents/"+key, robotsContent) // 일단 캐쉬만으로도 충분
	sitemaps, _ := crawler.ExtractSitemap(robotsTxt)
	if len(sitemaps) > 0 {
		for _, sitemap := range sitemaps {
			crawler.CrawlingSitemap(link, sitemap)
		}
	}
}

func main() {
	// https://pypi.org/
	// https://www.krcert.or.kr/
	// https://www.cloudflare.com/
	var c = make(chan string)
	go SubscribeKafkaTopics(c, []string{"startURL"})
	go downloadAllURL(c)
	//go SubscribeKafkaTopics(c, GetKafkaTopicsName())

	var s string
	fmt.Scanln(&s)
}
