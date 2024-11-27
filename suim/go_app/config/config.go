package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// DBConfig 구조체: 데이터베이스 설정값을 저장
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDBConfig 함수: 환경 변수에서 데이터베이스 설정값 로드
func LoadDBConfig() DBConfig {
	// .env 파일 로드 (환경 변수에 값이 이미 있으면 덮어쓰지 않음)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or not loaded. Using system environment variables.")
	}

	// 환경 변수에서 설정값 읽기
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv 함수: 환경 변수 값 가져오기 (기본값 설정 가능)
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
