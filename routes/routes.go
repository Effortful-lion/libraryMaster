package routes

import (
	"librarysystem/controllers"
	"librarysystem/middleware"
	"librarysystem/utils"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 创建路由
	r := gin.Default()

	// 设置模板目录和静态文件目录
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "./static")

	 // 加载自定义函数到模板
	r.SetFuncMap(utils.TemplateFunctions())

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