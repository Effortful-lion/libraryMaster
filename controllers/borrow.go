package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"librarysystem/models"
	"librarysystem/utils"
)

// BorrowForm 借阅表单结构
type BorrowForm struct {
	UserID    int    `form:"user_id" binding:"required"`
	BookID    int    `form:"book_id" binding:"required"`
	CSRFToken string `form:"csrf_token"`
}

// AdminUsersGet 处理GET /admin/users
func AdminUsersGet(c *gin.Context) {
	// 获取所有用户
	users := models.GetAllUsers()

	// 渲染用户管理页面
	c.HTML(http.StatusOK, "admin/users.html", gin.H{
		"title": "用户管理",
		"users": users,
	})
}

// AdminChangeUserRoleGet 处理GET /admin/change-user-role/:id/:role
func AdminChangeUserRoleGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	// 获取角色
	role := c.Param("role")
	if role != "admin" && role != "librarian" && role != "reader" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "无效的角色",
		})
		return
	}

	// 获取用户
	user, err := models.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 更新角色
	err = user.UpdateRole(models.UserRole(role))
	if err != nil {
		mg.SetFlashMessage(c, "error", err.Error())
	} else {
		mg.SetFlashMessage(c, "success", "用户角色已更新为"+role)
	}

	// 重定向回用户列表
	c.Redirect(http.StatusFound, "/admin/users")
}

// LibrarianBorrowGet 处理GET /librarian/borrow
func LibrarianBorrowGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取所有借阅记录
	allRecords := models.GetAllBorrowRecords()

	// 获取所有用户
	users := models.GetAllUsers()
	userMap := make(map[int]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 获取所有图书
	books := models.GetAllBooks()
	bookMap := make(map[int]*models.Book)
	for _, book := range books {
		bookMap[book.ID] = book
	}

	// 生成CSRF令牌
	token := mg.GenerateCSRFToken(c)

	// 获取有库存的图书（用于新借阅）
	var availableBooks []*models.Book
	for _, book := range books {
		if book.GetAvailableQuantity() > 0 {
			availableBooks = append(availableBooks, book)
		}
	}

	// 渲染借阅管理页面
	c.HTML(http.StatusOK, "librarian/borrow.html", gin.H{
		"title":           "借阅管理",
		"all_records":     allRecords,
		"users":           users,
		"books":           books,
		"user_map":        userMap,
		"book_map":        bookMap,
		"csrf_token":      token,
		"available_books": availableBooks,
	})
}

// LibrarianCreateBorrowPost 处理POST /librarian/create-borrow
func LibrarianCreateBorrowPost(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	var form BorrowForm
	if err := c.ShouldBind(&form); err != nil {
		// 表单验证失败
		mg.SetFlashMessage(c, "error", "请选择用户和图书")
		c.Redirect(http.StatusFound, "/librarian/borrow")
		return
	}

	// 验证CSRF令牌
	if !mg.VerifyCSRFToken(c, form.CSRFToken) {
		mg.SetFlashMessage(c, "error", "安全验证失败，请重试")
		c.Redirect(http.StatusFound, "/librarian/borrow")
		return
	}

	// 设置借阅日期和到期日期
	now := time.Now()
	dueDate := now.AddDate(0, 0, 14) // 14天后到期

	// 创建借阅记录
	_, err := models.CreateBorrowRecord(form.UserID, form.BookID, now, dueDate)
	if err != nil {
		// 创建失败
		mg.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/librarian/borrow")
		return
	}

	// 设置成功消息
	mg.SetFlashMessage(c, "success", "借阅记录已创建")
	
	// 重定向回借阅管理页面
	c.Redirect(http.StatusFound, "/librarian/borrow")
}

// LibrarianReturnBookGet 处理GET /librarian/return-book/:id
func LibrarianReturnBookGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取借阅记录ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的借阅记录ID",
		})
		return
	}

	// 归还图书
	_, err = models.ReturnBook(id)
	if err != nil {
		mg.SetFlashMessage(c, "error", err.Error())
	} else {
		mg.SetFlashMessage(c, "success", "图书已成功归还")
	}

	// 重定向回借阅管理页面
	c.Redirect(http.StatusFound, "/librarian/borrow")
}

// ReaderBooksGet 处理GET /reader/books
func ReaderBooksGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
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

	// 获取用户ID
	userID := mg.GetUserIDFromSession(c)

	// 获取用户当前借阅的图书ID
	borrowedBookIDs := make(map[int]bool)
	activeRecords := models.GetActiveBorrowRecordsByUserID(userID)
	for _, record := range activeRecords {
		borrowedBookIDs[record.BookID] = true
	}

	// 渲染读者图书列表页面
	c.HTML(http.StatusOK, "reader/books.html", gin.H{
		"title":           "查找图书",
		"books":           books,
		"categories":      categories,
		"query":           query,
		"category":        category,
		"borrowed_books":  borrowedBookIDs,
	})
}

// ReaderBorrowGet 处理GET /reader/borrow/:id
func ReaderBorrowGet(c *gin.Context) {
	// 获取图书ID
	mg := utils.NewSessionManager(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的图书ID",
		})
		return
	}

	// 获取用户ID
	userID := mg.GetUserIDFromSession(c)

	// 设置借阅日期和到期日期
	now := time.Now()
	dueDate := now.AddDate(0, 0, 14) // 14天后到期

	// 创建借阅记录
	_, err = models.CreateBorrowRecord(userID, id, now, dueDate)
	if err != nil {
		mg.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/books/"+idStr)
		return
	}

	// 设置成功消息
	mg.SetFlashMessage(c, "success", "图书借阅成功，请在两周内归还")
	
	// 重定向到借阅记录页面
	c.Redirect(http.StatusFound, "/reader/borrowed")
}

// ReaderBorrowedGet 处理GET /reader/borrowed
func ReaderBorrowedGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取用户ID
	userID := mg.GetUserIDFromSession(c)

	// 获取用户的所有借阅记录
	records := models.GetBorrowRecordsByUserID(userID)

	// 获取所有图书
	books := models.GetAllBooks()
	bookMap := make(map[int]*models.Book)
	for _, book := range books {
		bookMap[book.ID] = book
	}

	// 处理记录，计算逾期等信息
	type EnhancedRecord struct {
		Record         *models.BorrowRecord
		Book           *models.Book
		IsActive       bool
		IsOverdue      bool
		DaysUntilDue   int
		OverdueDays    int
		StatusText     string
	}

	var enhancedRecords []EnhancedRecord
	var activeCount, overdueCount int

	for _, record := range records {
		book := bookMap[record.BookID]
		isActive := record.ReturnDate.IsZero()
		isOverdue := record.IsOverdue()
		
		// 计算状态文本
		var statusText string
		if !isActive {
			statusText = "已归还"
		} else if isOverdue {
			statusText = "已逾期"
		} else {
			statusText = "借阅中"
		}

		// 创建增强记录
		enhancedRecord := EnhancedRecord{
			Record:         record,
			Book:           book,
			IsActive:       isActive,
			IsOverdue:      isOverdue,
			DaysUntilDue:   record.DaysUntilDue(),
			OverdueDays:    record.OverdueDays(),
			StatusText:     statusText,
		}

		enhancedRecords = append(enhancedRecords, enhancedRecord)

		// 统计
		if isActive {
			activeCount++
			if isOverdue {
				overdueCount++
			}
		}
	}

	// 渲染借阅记录页面
	c.HTML(http.StatusOK, "reader/borrowed.html", gin.H{
		"title":            "我的借阅",
		"records":          enhancedRecords,
		"active_count":     activeCount,
		"overdue_count":    overdueCount,
		"book_map":         bookMap,
	})
}

// ReaderReturnBookGet 处理GET /reader/return-book/:id
func ReaderReturnBookGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取借阅记录ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "无效的借阅记录ID",
		})
		return
	}

	// 获取借阅记录
	record, err := models.GetBorrowRecordByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "借阅记录不存在",
		})
		return
	}

	// 检查记录是否属于当前用户
	userID := mg.GetUserIDFromSession(c)
	if record.UserID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "无权操作此借阅记录",
		})
		return
	}

	// 设置成功消息（读者暂时不能自己归还图书，只提供通知）
	mg.SetFlashMessage(c, "info", "请将图书交给图书管理员归还")
	
	// 重定向回借阅记录页面
	c.Redirect(http.StatusFound, "/reader/borrowed")
}