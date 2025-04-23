package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // 添加MySQL驱动
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDatabase() {
	once.Do(func() {
		// 获取数据库连接参数（添加空值校验）
		dbhost := mustGetEnv("DB_HOST")
		dbport := mustGetEnv("DB_PORT")
		dbuser := mustGetEnv("DB_USER")
		dbpassword := mustGetEnv("DB_PASS")
		dbname := mustGetEnv("DB_NAME")

		// 添加端口默认值处理
		if dbport == "" {
			dbport = "3306"
		}

		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbuser, dbpassword, dbhost, dbport, dbname) // 构建MySQL连接字符串

		// 连接数据库
		var err error
		db, err = sql.Open("mysql", dbURL) // 修改驱动类型为mysql
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

// 新增环境变量校验函数
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func GetDB() *sql.DB {
	return db
}
