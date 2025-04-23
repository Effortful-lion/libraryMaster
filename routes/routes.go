package routes

import (
	"librarysystem/controllers"
	"librarysystem/middleware"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
    
    // 添加静态文件路由
    r.Static("/static", "./static")
    
    r.SetFuncMap(template.FuncMap{
		"isNil": func(i interface{}) bool {
			return i == nil || reflect.ValueOf(i).IsNil()
		},
		"isOverDue": func(t time.Time) bool {
			return time.Now().After(t)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02") // Go 的日期格式化模板
		},
		"daysSince": func(t time.Time) int {
			return int(time.Since(t).Hours() / 24)
		},
		"daysUntil": func(t time.Time) int {
			return int(time.Until(t).Hours() / 24)
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i + 1
			}
			return s
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04") // 包含日期和时间
		},
		"url_for": func(routeName string, pairs ...string) string {
			routes := map[string]string{
				"static": "/static/%s",  // 新增静态文件路由
				"book_detail": "/books/%s",  // 修改为格式化字符串
				"index":           "/",
				"books":           "/books",
				"reader_books":    "/reader/books",
				"reader_borrowed": "/reader/borrowed",
				// 添加更多路由映射...
			}
			
			if len(pairs)%2 != 0 {
				return ""
			}
			
			path := routes[routeName]
			for i := 0; i < len(pairs); i += 2 {
				path = strings.Replace(path, ":"+pairs[i], pairs[i+1], 1)
			}
			return path
		},
	})

	r.LoadHTMLGlob("templates/*.html")
	r.LoadHTMLGlob("templates/**/*.html")

	// 公共路由
	r.GET("/", controllers.IndexGet)
	r.GET("/books", controllers.BooksGet)
	r.GET("/books/:id", controllers.BookDetailGet)
	r.GET("/login", controllers.LoginGet)
	r.POST("/login", controllers.LoginPost)
	r.GET("/register", controllers.RegisterGet)
	r.POST("/register", controllers.RegisterPost)
	r.GET("/logout", controllers.Logout)

	// API路由
	api := r.Group("/api")
	{
		api.GET("/books", controllers.APIBooksGet)
		api.GET("/books/:id", controllers.APIBookGet)
		api.GET("/categories", controllers.APICategoriesGet)
		api.GET("/borrow-records", controllers.APIBorrowRecordsGet)
	}

	// 需要登录的路由
	auth := r.Group("")
	auth.Use(middleware.RequireAuth())
	{
		auth.GET("/dashboard", controllers.Dashboard)
	}

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	{
		admin.GET("/books", controllers.AdminBooksGet)
		admin.GET("/users", controllers.AdminUsersGet)
		admin.GET("/add-book", controllers.AdminAddBookGet)
		admin.POST("/add-book", controllers.AdminAddBookPost)
		admin.GET("/edit-book/:id", controllers.AdminEditBookGet)
		admin.POST("/edit-book/:id", controllers.AdminEditBookPost)
		admin.GET("/delete-book/:id", controllers.AdminDeleteBookGet)
		admin.GET("/change-user-role/:id/:role", controllers.AdminChangeUserRoleGet)
	}

	// 图书管理员路由
	librarian := r.Group("/librarian")
	librarian.Use(middleware.RequireLibrarian())
	{
		librarian.GET("/books", controllers.LibrarianBooksGet)
		librarian.GET("/borrow", controllers.LibrarianBorrowGet)
		librarian.POST("/create-borrow", controllers.LibrarianCreateBorrowPost)
		librarian.GET("/return-book/:id", controllers.LibrarianReturnBookGet)
	}

	// 读者路由
	reader := r.Group("/reader")
	reader.Use(middleware.RequireReader())
	{
		reader.GET("/books", controllers.ReaderBooksGet)
		reader.GET("/borrow/:id", controllers.ReaderBorrowGet)
		reader.GET("/borrowed", controllers.ReaderBorrowedGet)
		reader.GET("/return-book/:id", controllers.ReaderReturnBookGet)
	}

	return r
}
