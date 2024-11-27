package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// HashContent creates a SHA256 hash of the given content
func HashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash) // Convert hash to hexadecimal string
}

// SaveHashToFile saves the hash to a file
func SaveHashToFile(filename, hash string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(hash)
	if err != nil {
		return fmt.Errorf("failed to write hash to file: %v", err)
	}

	return nil
}

func Difference(a, b []string) []string {
	bSet := make(map[string]struct{})
	for _, item := range b {
		bSet[item] = struct{}{}
	}

	var result []string
	for _, item := range a {
		if _, found := bSet[item]; !found {
			result = append(result, item)
		}
	}
	return result
}
