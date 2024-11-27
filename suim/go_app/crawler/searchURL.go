package crawler

import (
	"encoding/xml"
	"fmt"
	"go_app/config"
	"go_app/utils"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// Sitemap represents a single sitemap entry
type Sitemap struct {
	Loc string `xml:"loc"`
}

// SitemapIndex represents the root sitemapindex element
type SitemapIndex struct {
	Sitemaps []Sitemap `xml:"sitemap"`
}

// URL represents a single <url> entry in the sitemap
type URL struct {
	Loc string `xml:"loc"`
}

// URLSet represents the root <urlset> element
type URLSet struct {
	URLs []URL `xml:"url"`
}

func FetchHTML(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("failed to fetch HTML from %s: %v", link, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func ParseHTML(link string) (*html.Node, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML from %s: %v", link, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return doc, nil
}

func FetchXml(link string) (index SitemapIndex) {
	body, err := FetchHTML(link)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode the XML data
	decoder := xml.NewDecoder(strings.NewReader(body))
	decoder.DefaultSpace = "http://www.sitemaps.org/schemas/sitemap/0.9" // Handle default namespace
	err = decoder.Decode(&index)
	if err != nil {
		fmt.Println("error decoding XML:", err)
		return
	}

	return index
}

func FetchUrl(link string) (urls []string) {
	body, err := FetchHTML(link)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the XML response body
	var urlSet URLSet
	decoder := xml.NewDecoder(strings.NewReader(body))
	err = decoder.Decode(&urlSet)
	if err != nil {
		fmt.Println("error decoding XML:", err)
		return
	}

	// Extract the URLs into a string slice
	for _, url := range urlSet.URLs {
		urls = append(urls, url.Loc)
	}
	return urls
}

func CrawlingSitemap(mainLink string, sitemapLink string) {
	key := DomainToCamel(mainLink)
	xmlArr := FetchXml(sitemapLink)
	db, _ := config.CreateConnection()

	defer config.CloseDB(db)

	config.InsertDomain(db, key, "https://"+CamelToDomain(key))

	for _, sitemap := range xmlArr.Sitemaps {
		dbUrlArr := config.GetAllUrl(db, key)
		urlArr := FetchUrl(sitemap.Loc)
		diff := utils.Difference(urlArr, dbUrlArr)
		config.InsertDomainPath(db, key, diff)
	}
}
