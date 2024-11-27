package utils

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

// FetchHTML: URL에서 HTML 응답을 가져옵니다.
func FetchHTML(link string) (*http.Response, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML from %s: %v", link, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close() // HTTP 상태가 비정상이면 즉시 리소스를 정리
		return nil, fmt.Errorf("failed to fetch HTML from %s: HTTP %d", link, resp.StatusCode)
	}

	return resp, nil
}

// ParseHTML: FetchHTML을 통해 HTML 문서를 가져오고, 파싱한 결과를 반환합니다.
func ParseHTML(link string) (*html.Node, error) {
	resp, err := FetchHTML(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // 응답 바디 닫기

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return doc, nil
}
