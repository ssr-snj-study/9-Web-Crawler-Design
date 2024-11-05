package internal

import "strings"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encodeBase62(num int64) string {
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
