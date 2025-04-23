package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // 替换为MySQL驱动
)

var DB *sql.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
        log.Println("初始化数据库连接...")
        
        // 获取数据库连接URL
        dbhost := os.Getenv("DB_HOST")
        dbport := os.Getenv("DB_PORT")
        dbuser := os.Getenv("DB_USER")
        dbpassword := os.Getenv("DB_PASSWORD")
        dbname := os.Getenv("DB_NAME")
        dburl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpassword, dbhost, dbport, dbname)
        if dburl == "" {
                log.Fatal("mysql的url为空")
        }
        
        // 连接到数据库
        var err error
        DB, err = sql.Open("mysql", dburl)
        if err != nil {
                log.Fatalf("连接数据库失败: %v", err)
        }
        
        // 验证连接
        err = DB.Ping()
        if err != nil {
                log.Fatalf("数据库连接验证失败: %v", err)
        }
        
        log.Println("数据库连接初始化完成")
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
        return DB
}