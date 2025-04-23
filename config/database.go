
package config

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDatabase() {
	once.Do(func() {
		// 获取数据库连接参数
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			log.Fatal("DATABASE_URL environment variable is not set")
		}

		// 连接数据库
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		// 测试连接
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}

		log.Println("Successfully connected to database")
	})
}

func GetDB() *sql.DB {
	return db
}
