package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

func MyIP() []string {
	myIp := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Local IP addresses:")
	for _, addr := range addrs {
		// IPv4 주소만 필터링
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String())
				myIp = append(myIp, ipNet.IP.String())
			}
		}
	}
	return myIp
}

func GetHtmlHash(html string) string {
	hash := sha256.New()
	hash.Write([]byte(html))

	// 해시 값을 16진수 문자열로 변환
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func encodeBase62(num int64) string {

	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if num == 0 {
		return "0"
	}

	var result strings.Builder
	base := int64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		result.WriteByte(base62Chars[remainder])
		num /= base
	}

	// 문자열을 뒤집어 최종 결과 생성
	encoded := result.String()
	// 뒤집힌 결과를 올바른 순서로 바꿔 반환
	return reverseString(encoded)
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
