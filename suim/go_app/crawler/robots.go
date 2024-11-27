package crawler

import (
	"bufio"
	"fmt"
	"go_app/config"
	"go_app/utils"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

func DownloadRobots(link string) (string, error) {
	domain, _ := SplitDomain(link)
	robotsUrl := "https://" + domain + "/robots.txt"

	resp, err := utils.FetchHTML(robotsUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch html robots.txt: %v", err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read html robots.txt: %v", err)
	}
	return string(content), nil
}

func CompareRobotsTxt(link string) (bool, string) {
	robotsContent, err := DownloadRobots(link)
	if err != nil {
		fmt.Printf("failed to download robots.txt: %v\n", err)
		return false, ""
	}

	robotsKey := DomainToCamel(link)
	robotsHash := utils.HashContent(robotsContent)

	value, err := config.GetRedis(robotsKey)
	if err != nil {
		fmt.Printf("failed to get robots.txt: %v\n", err)
		config.SetRedis(robotsKey, robotsHash, 24*time.Hour)
		return false, robotsContent
	}

	if value == robotsHash {
		return true, ""
	} else {
		config.SetRedis(robotsKey, robotsHash, 24*time.Hour)
		return false, robotsContent
	}
}

func ExtractSitemap(contents string) ([]string, error) {
	// Regular expression to match "Sitemap" lines
	sitemapRegex := regexp.MustCompile(`(?i)^Sitemap:\s*(https://.+)$`)
	var sitemaps []string

	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := sitemapRegex.FindStringSubmatch(line); len(matches) == 2 {
			sitemaps = append(sitemaps, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning robots.txt content: %v", err)
	}

	return sitemaps, nil
}
