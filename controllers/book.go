package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"librarysystem/models"
	"librarysystem/utils"
)

// BookForm 图书表单结构
type BookForm struct {
	Title         string `form:"title" binding:"required"`
	Author        string `form:"author" binding:"required"`
	ISBN          string `form:"isbn" binding:"required"`
	PublishedYear int    `form:"published_year" binding:"required,min=1000,max=2100"`
	Category      string `form:"category" binding:"required"`
	Description   string `form:"description" binding:"required"`
	CoverURL      string `form:"cover_url" binding:"required,url"`
	Quantity      int    `form:"quantity" binding:"required,min=1"`
	CSRFToken     string `form:"csrf_token"`
}

// IndexGet 处理GET /
func IndexGet(c *gin.Context) {
	// 获取特色图书（前4本）
	allBooks := models.GetAllBooks()
	var featuredBooks []*models.Book
	if len(allBooks) > 4 {
		featuredBooks = allBooks[:4]
	} else {
		featuredBooks = allBooks
	}

	// 获取所有分类
	categories := models.GetAllCategories()

	// 渲染首页
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "首页",
		"featured_books": featuredBooks,
		"categories":     categories,
	})
}

// BooksGet 处理GET /books
func BooksGet(c *gin.Context) {
	// 获取查询参数
	query := c.Query("q")
	category := c.Query("category")

	// 根据查询参数获取图书
	var books []*models.Book
	if query != "" {
		books = models.SearchBooks(query)
	} else if category != "" {
		books = models.GetBooksByCategory(category)
	} else {
		books = models.GetAllBooks()
	}

	// 获取所有分类
	categories := models.GetAllCategories()

	// 渲染图书列表页面
	c.HTML(http.StatusOK, "book_list.html", gin.H{
		"title":      "图书列表",
		"books":      books,
		"categories": categories,
		"query":      query,
		"category":   category,
	})
}

// BookDetailGet 处理GET /books/:id
func BookDetailGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取图书ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的图书ID",
		})
		return
	}

	// 获取图书
	book, err := models.GetBookByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "图书不存在",
		})
		return
	}

	// 检查图书是否可借
	available := book.IsAvailable()
	availableCount := book.GetAvailableQuantity()

	// 获取用户信息
	userID := mg.GetUserIDFromSession(c)
	userRole := mg.GetUserRoleFromSession(c)

	// 渲染图书详情页面
	c.HTML(http.StatusOK, "book_detail.html", gin.H{
		"title":           book.Title,
		"book":            book,
		"available":       available,
		"availableCount":  availableCount,
		"user_id":         userID,
		"user_role":       userRole,
	})
}

// AdminBooksGet 处理GET /admin/books
func AdminBooksGet(c *gin.Context) {
	// 获取所有图书
	books := models.GetAllBooks()

	// 渲染管理员图书管理页面
	c.HTML(http.StatusOK, "admin/books.html", gin.H{
		"title": "图书管理",
		"books": books,
	})
}

// AdminAddBookGet 处理GET /admin/add-book
func AdminAddBookGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 生成CSRF令牌
	token := mg.GenerateCSRFToken(c)

	// 获取错误消息，如果有
	errorMsg := mg.GetFlashMessage(c, "error")

	// 渲染添加图书页面
	c.HTML(http.StatusOK, "admin/edit_book.html", gin.H{
		"title":      "添加图书",
		"csrf_token": token,
		"error":      errorMsg,
		"is_add":     true,
	})
}

// AdminAddBookPost 处理POST /admin/add-book
func AdminAddBookPost(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	var form BookForm
	if err := c.ShouldBind(&form); err != nil {
		// 表单验证失败
		mg.SetFlashMessage(c, "error", "请正确填写所有必填字段")
		c.Redirect(http.StatusFound, "/admin/add-book")
		return
	}

	// 验证CSRF令牌
	if !mg.VerifyCSRFToken(c, form.CSRFToken) {
		mg.SetFlashMessage(c, "error", "安全验证失败，请重试")
		c.Redirect(http.StatusFound, "/admin/add-book")
		return
	}

	// 创建新图书
	book, err := models.CreateBook(
		form.Title,
		form.Author,
		form.ISBN,
		form.PublishedYear,
		form.Category,
		form.Description,
		form.CoverURL,
		form.Quantity,
	)
	if err != nil {
		// 图书创建失败
		mg.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/admin/add-book")
		return
	}

	// 设置成功消息
	mg.SetFlashMessage(c, "success", "图书添加成功："+book.Title)
	
	// 重定向到图书列表
	c.Redirect(http.StatusFound, "/admin/books")
}

// AdminEditBookGet 处理GET /admin/edit-book/:id
func AdminEditBookGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取图书ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的图书ID",
		})
		return
	}

	// 获取图书
	book, err := models.GetBookByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "图书不存在",
		})
		return
	}

	// 生成CSRF令牌
	token := mg.GenerateCSRFToken(c)

	// 获取错误消息，如果有
	errorMsg := mg.GetFlashMessage(c, "error")

	// 渲染编辑图书页面
	c.HTML(http.StatusOK, "admin/edit_book.html", gin.H{
		"title":      "编辑图书",
		"book":       book,
		"csrf_token": token,
		"error":      errorMsg,
		"is_add":     false,
	})
}

// AdminEditBookPost 处理POST /admin/edit-book/:id
func AdminEditBookPost(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取图书ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的图书ID",
		})
		return
	}

	var form BookForm
	if err := c.ShouldBind(&form); err != nil {
		// 表单验证失败
		mg.SetFlashMessage(c, "error", "请正确填写所有必填字段")
		c.Redirect(http.StatusFound, "/admin/edit-book/"+idStr)
		return
	}

	// 验证CSRF令牌
	if !mg.VerifyCSRFToken(c, form.CSRFToken) {
		mg.SetFlashMessage(c, "error", "安全验证失败，请重试")
		c.Redirect(http.StatusFound, "/admin/edit-book/"+idStr)
		return
	}

	// 更新图书
	book, err := models.UpdateBook(
		id,
		form.Title,
		form.Author,
		form.ISBN,
		form.PublishedYear,
		form.Category,
		form.Description,
		form.CoverURL,
		form.Quantity,
	)
	if err != nil {
		// 图书更新失败
		mg.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/admin/edit-book/"+idStr)
		return
	}

	// 设置成功消息
	mg.SetFlashMessage(c, "success", "图书更新成功："+book.Title)
	
	// 重定向到图书列表
	c.Redirect(http.StatusFound, "/admin/books")
}

// AdminDeleteBookGet 处理GET /admin/delete-book/:id
func AdminDeleteBookGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取图书ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的图书ID",
		})
		return
	}

	// 删除图书
	err = models.DeleteBook(id)
	if err != nil {
		// 删除失败
		mg.SetFlashMessage(c, "error", err.Error())
	} else {
		// 删除成功
		mg.SetFlashMessage(c, "success", "图书删除成功")
	}
	
	// 重定向到图书列表
	c.Redirect(http.StatusFound, "/admin/books")
}

// LibrarianBooksGet 处理GET /librarian/books
func LibrarianBooksGet(c *gin.Context) {
	// 获取所有图书
	books := models.GetAllBooks()

	// 获取每本书的借阅情况
	bookStatus := make(map[int]map[string]interface{})
	for _, book := range books {
		activeRecords := models.GetActiveBorrowRecordsByBookID(book.ID)
		bookStatus[book.ID] = map[string]interface{}{
			"total":           book.Quantity,
			"available":       book.GetAvailableQuantity(),
			"borrowed":        len(activeRecords),
			"active_records":  activeRecords,
		}
	}

	// 渲染图书管理员图书管理页面
	c.HTML(http.StatusOK, "librarian/books.html", gin.H{
		"title":       "图书库存管理",
		"books":       books,
		"book_status": bookStatus,
	})
}