package main

import (
        "log"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "github.com/gin-gonic/gin"
        "librarysystem/config"
        "librarysystem/database"
        "librarysystem/routes"
        "librarysystem/utils"
)

func main() {
        // 设置日志格式
        log.SetFlags(log.LstdFlags | log.Lshortfile)
        log.Println("启动图书管理系统...")

        // 设置模式（生产/开发）
        gin.SetMode(gin.ReleaseMode)

        // 初始化数据库连接
        log.Println("初始化数据库连接...")
        config.InitDatabase()
        
        // 初始化数据库表和示例数据
        log.Println("初始化数据库表...")
        database.InitTables()
        database.SeedData()

        // 创建路由
        log.Println("设置路由...")
        router := routes.SetupRouter()

        // 启动会话清理定时任务
        go sessionCleanupTask()

        // 启动服务器
        port := os.Getenv("PORT")
        if port == "" {
                port = "5000" // 默认端口，与Replit工作流一致
        }
        
        // 优雅关闭服务器
        srv := &http.Server{
                Addr:    "0.0.0.0:" + port,
                Handler: router,
        }

        go func() {
                log.Printf("服务器正在运行，地址为 http://0.0.0.0:%s\n", port)
                if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        log.Fatalf("监听失败: %v\n", err)
                }
        }()

        // 等待中断信号
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        log.Println("关闭服务器...")

        // 设置上下文超时
        time.Sleep(1 * time.Second)
        log.Println("服务器已关闭")
}

// 会话清理定时任务
func sessionCleanupTask() {
        ticker := time.NewTicker(1 * time.Hour)
        defer ticker.Stop()

        for range ticker.C {
                log.Println("清理过期会话...")
                utils.CleanupExpiredSessions()
        }
}