package internal

import (
	"fmt"
	"net"
	"net/url"
)

func GetIpFromUrl(rawURL string) (string, error) {
	var ip string
	// DNS를 통해 IP 주소 조회
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %s: %v", rawURL, err)
	}

	// Lookup the IP address using the hostname
	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return "", fmt.Errorf("failed to lookup IP for %s: %v", parsedURL.Hostname(), err)
	}

	// IP 주소 출력
	fmt.Printf("IP addresses for %s:\n", ips)
	for _, urlIp := range ips {
		ip = fmt.Sprintf("%s", urlIp)
		break
	}

	return ip, nil
}
