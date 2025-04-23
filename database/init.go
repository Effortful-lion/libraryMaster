package database

import (
        "log"
        "time"

        "librarysystem/config"
)

// InitTables 初始化数据库表
func InitTables() {
        log.Println("初始化数据库表...")
        
        db := config.GetDB()
        
        // 创建用户表
        _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(200) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
        `)
        if err != nil {
                log.Fatalf("创建用户表失败: %v", err)
        }
        
        // 创建图书表
        _, err = db.Exec(`
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL,
    isbn VARCHAR(20) NOT NULL UNIQUE,
    published_year INTEGER NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    cover_url TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
        `)
        if err != nil {
                log.Fatalf("创建图书表失败: %v", err)
        }
        
        // 创建借阅记录表
        _, err = db.Exec(`
CREATE TABLE IF NOT EXISTS borrow_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    book_id INTEGER NOT NULL REFERENCES books(id),
    borrow_date TIMESTAMP WITH TIME ZONE NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    return_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
        `)
        if err != nil {
                log.Fatalf("创建借阅记录表失败: %v", err)
        }
        
        log.Println("数据库表初始化完成")
}

// ClearData 清空数据库数据（仅用于开发和测试）
func ClearData() {
        log.Println("清空所有数据表...")
        
        db := config.GetDB()
        
        // 清空借阅记录表
        _, err := db.Exec("DELETE FROM borrow_records")
        if err != nil {
                log.Printf("清空借阅记录表失败: %v", err)
        }
        
        // 清空图书表
        _, err = db.Exec("DELETE FROM books")
        if err != nil {
                log.Printf("清空图书表失败: %v", err)
        }
        
        // 清空用户表
        _, err = db.Exec("DELETE FROM users")
        if err != nil {
                log.Printf("清空用户表失败: %v", err)
        }
        
        log.Println("所有数据表已清空")
}

// SeedData 填充初始数据
func SeedData() {
        log.Println("填充初始数据...")
        
        // 检查是否已有数据
        db := config.GetDB()
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
        if err != nil {
                log.Printf("检查用户数据失败: %v", err)
        } else if count > 0 {
                log.Println("数据库已包含数据，跳过初始化")
                return
        }
        
        // 添加用户
        InitUsers()
        
        // 添加图书
        InitBooks()
        
        // 添加借阅记录
        InitBorrowRecords()
        
        log.Println("初始数据填充完成")
}

// InitUsers 初始化用户数据
func InitUsers() {
        log.Println("初始化用户数据...")
        db := config.GetDB()
        
        // 创建密码哈希（在实际应用中应使用bcrypt）
        adminPassword := "admin123"
        librarianPassword := "librarian123"
        readerPassword := "reader123"
        
        // 插入用户
        _, err := db.Exec(`
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@example.com', $1, 'admin'),
('librarian', 'librarian@example.com', $2, 'librarian'),
('reader', 'reader@example.com', $3, 'reader')
        `, adminPassword, librarianPassword, readerPassword)
        
        if err != nil {
                log.Fatalf("初始化用户数据失败: %v", err)
        }
        
        log.Println("用户数据初始化完成")
}

// InitBooks 初始化图书数据
func InitBooks() {
        log.Println("初始化图书数据...")
        db := config.GetDB()
        
        // 插入图书数据
        _, err := db.Exec(`
INSERT INTO books (title, author, isbn, published_year, category, description, cover_url, quantity) VALUES 
('Python编程：从入门到实践', '埃里克·马瑟斯', '9787115428028', 2016, '编程', 
 '本书是一本针对所有层次的Python读者而作的Python入门书。全书分两部分：第一部分介绍用Python编程所必须了解的基本概念，第二部分将理论付诸实践，讲解如何开发三个项目。', 
 'https://img3.doubanio.com/view/subject/l/public/s29424065.jpg', 5),
 
('Go语言实战', '威廉·肯尼迪', '9787115445353', 2017, '编程', 
 '本书首先介绍Go语言的独特之处，然后讲解如何编写地道的Go代码并使用其特有的特性和工具包编写代码。后续章节会介绍测试、Web编程以及与其他主流语言的集成。', 
 'https://img9.doubanio.com/view/subject/l/public/s29446435.jpg', 3),
 
('明朝那些事儿', '当年明月', '9787807023630', 2009, '历史', 
 '《明朝那些事儿》讲述从1344年到1644年，明朝三百年间的历史。以史料为基础，以年代和具体人物为主线，运用小说的手法，对明朝十七帝和其他王公权贵和小人物的命运进行全景展示。', 
 'https://img1.doubanio.com/view/subject/l/public/s27131114.jpg', 4),
 
('三体', '刘慈欣', '9787536692930', 2008, '科幻', 
 '文化大革命如火如荼进行的同时，军方探寻外星文明的绝秘计划"红岸工程"取得了突破性进展。但在按下发射键的那一刻，历经劫难的叶文洁没有意识到，她彻底改变了人类的命运。', 
 'https://img2.doubanio.com/view/subject/l/public/s2768378.jpg', 2),
 
('围城', '钱钟书', '9787020090006', 1991, '文学', 
 '《围城》是钱钟书所著的长篇小说，自问世以来，就以它的犀利的语言、巧妙的结构和象征性的意义在中国文学史上占据重要地位。', 
 'https://img2.doubanio.com/view/subject/l/public/s1070222.jpg', 3)
        `)
        
        if err != nil {
                log.Fatalf("初始化图书数据失败: %v", err)
        }
        
        log.Println("图书数据初始化完成")
}

// InitBorrowRecords 初始化借阅记录
func InitBorrowRecords() {
        log.Println("初始化借阅记录...")
        db := config.GetDB()
        
        // 设置时间
        now := time.Now()
        oneWeekAgo := now.AddDate(0, 0, -7)
        twoWeeksAgo := now.AddDate(0, 0, -14)
        inOneWeek := now.AddDate(0, 0, 7)
        inTwoWeeks := now.AddDate(0, 0, 14)
        overdueDate := now.AddDate(0, 0, -1) // 昨天到期
        
        // 插入借阅记录
        _, err := db.Exec(`
INSERT INTO borrow_records (user_id, book_id, borrow_date, due_date, return_date) VALUES 
(1, 1, $1, $2, $3), -- 管理员已归还记录
(2, 2, $4, $5, NULL), -- 图书管理员当前借阅
(3, 3, $6, $7, NULL), -- 读者当前借阅
(3, 4, $8, $9, NULL)  -- 读者逾期未归还
        `, twoWeeksAgo, twoWeeksAgo.AddDate(0, 0, 14), now, // 管理员已归还
           oneWeekAgo, inOneWeek, // 图书管理员当前借阅
           oneWeekAgo, inTwoWeeks, // 读者当前借阅
           twoWeeksAgo, overdueDate) // 读者逾期未归还
        
        if err != nil {
                log.Fatalf("初始化借阅记录失败: %v", err)
        }
        
        log.Println("借阅记录初始化完成")
}