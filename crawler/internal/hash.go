package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHtmlHash(html string) string {
	hash := sha256.New()
	hash.Write([]byte(html))

	// 해시 값을 16진수 문자열로 변환
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
