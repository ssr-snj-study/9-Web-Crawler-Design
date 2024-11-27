package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

// Domain 구조체 정의
type Domain struct {
	Name      string
	UserAgent string
	Priority  int
}

func CreateConnection() (*sql.DB, error) {
	dbConfig := LoadDBConfig()

	// PostgreSQL 연결 문자열 구성
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	// 데이터베이스 연결
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the database: %v", err)
	}

	// 데이터베이스 연결 확인
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Unable to ping the database: %v", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func GetAllUrl(db *sql.DB, domainKey string) []string {
	query := `
	SELECT concat(trim(TRAILING '/' FROM d.link), dp.path) as url
	FROM public.domain_path dp
	JOIN public.domain d ON dp.domain_id = d.id
	WHERE d.key = $1
	`

	rows, err := db.Query(query, domainKey)
	if err != nil {
		log.Fatalf("failed to execute domain query: %v\n", err)
		return nil
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
			return nil
		}
		paths = append(paths, path)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v\n", err)
		return nil
	}

	return paths
}

func GetDomainId(db *sql.DB, key string) (int64, error) {
	var existingID int64
	query := `SELECT id FROM domain WHERE key = $1;`
	err := db.QueryRow(query, key).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to check for existing key: %w", err)
	}
	return existingID, nil
}

func InsertDomain(db *sql.DB, key string, link string) error {
	if key == "" {
		return fmt.Errorf("key and link are required fields")
	}

	domainId, _ := GetDomainId(db, key)
	if domainId != 0 {
		return nil
	}

	query := `
		INSERT INTO domain (link, key, priority)
		VALUES ($1, $2, DEFAULT)
		RETURNING id;`

	var id int64

	// Execute the query
	err := db.QueryRow(query, link, key).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert into domain: %w", err)
	}

	log.Printf("Inserted row with ID: %d", id)
	return nil
}

func InsertDomainPath(db *sql.DB, domainKey string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	domainId, _ := GetDomainId(db, domainKey)
	query := `INSERT INTO public.domain_path (domain_id, path) VALUES `
	valueStrings := make([]string, 0, len(paths))
	valueArgs := make([]interface{}, 0, len(paths))

	for i, path := range paths {
		valueStrings = append(valueStrings, fmt.Sprintf("(%d, $%d)", domainId, i+1))
		valueArgs = append(valueArgs, path)
	}

	query += strings.Join(valueStrings, ",")
	query += " ON CONFLICT (path) DO NOTHING;" // Optional: Avoid duplicate insertions

	_, err := db.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to insert paths: %w", err)
	}
	return nil
}

//func IsExistURL(db *sql.DB, url string) bool {
//	domain := strings.Split(url, "/")[0]
//	query := `
//	SELECT EXISTS (
//		SELECT 1
//		FROM public.domain_path dp
//		JOIN public.domain d ON dp.domain_id = d.id
//		WHERE d.name = $1 AND dp.path = $2
//	)
//	`
//
//	var exists bool
//	err := db.QueryRow(query, domainName, path).Scan(&exists)
//	if err != nil {
//		return false, fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	return exists, nil
//}
//
//func InsertDomain(db *sql.DB, domain models.Domain) {
//	query := `
//	INSERT INTO public.domain (name, priority)
//	VALUES ($1, $2)
//	ON CONFLICT (name) DO NOTHING; -- 중복 방지
//	`
//
//	_, err := db.Exec(query, domain.Name, domain.Priority)
//	if err != nil {
//		log.Fatalf("Failed to insert domain: %v\n", err)
//	}
//}
//
//func UpdateDomainPriority(db *sql.DB, domain models.Domain) error {
//	query := `UPDATE public.domain SET priority=$1 WHERE `
//
//	_, err := db.Exec(query, domain.Name)
//	return err
//}
//
//func connectionRest(db *sql.DB) {
//
//	// 수집한 도메인 데이터
//	domains := []Domain{
//		{Name: "example.com", UserAgent: "Mozilla/5.0", Priority: 5},
//		{Name: "test.com", UserAgent: "Googlebot", Priority: 3},
//		{Name: "anotherdomain.com", UserAgent: "Bingbot", Priority: 7},
//	}
//
//	// 데이터를 데이터베이스에 삽입
//	for _, domain := range domains {
//		err := insertDomain(db, domain)
//		if err != nil {
//			log.Printf("Failed to insert domain %s: %v\n", domain.Name, err)
//		} else {
//			log.Printf("Successfully inserted domain: %s\n", domain.Name)
//		}
//	}
//}
//func execute(db *sql.DB, query, ){
//	_, err := db.Exec(query, domain.Name, domain.UserAgent, domain.Priority)
//
//}
//
//func Select
