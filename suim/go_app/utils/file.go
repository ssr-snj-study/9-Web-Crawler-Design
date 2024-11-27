package utils

import (
	"fmt"
	"io"
	"os"
)

func SaveFile(filePath string, contents string) error {
	//fileName := strings.ReplaceAll(strings.TrimPrefix(domain, "https://"), "/", "_") + "_robots.txt"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// 파일에 데이터 쓰기
	_, err = file.WriteString(contents)
	if err != nil {
		return fmt.Errorf("Failed to save file %s: %v", filePath, err)
	}
	return nil
}

func ReadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	return string(content), nil
}

func CheckFile(filePath string, contents string) error {
	return nil
}
