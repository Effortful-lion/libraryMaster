package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"librarysystem/models"
)

// APIBooksGet 处理GET /api/books
func APIBooksGet(c *gin.Context) {
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

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"books": books,
		},
	})
}

// APIBookGet 处理GET /api/books/:id
func APIBookGet(c *gin.Context) {
	// 获取图书ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "无效的图书ID",
		})
		return
	}

	// 获取图书
	book, err := models.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"message": "图书不存在",
		})
		return
	}

	// 获取可用数量
	availableCount := book.GetAvailableQuantity()

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"book": book,
			"available": availableCount > 0,
			"available_count": availableCount,
		},
	})
}

// APICategoriesGet 处理GET /api/categories
func APICategoriesGet(c *gin.Context) {
	// 获取所有分类
	categories := models.GetAllCategories()

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"categories": categories,
		},
	})
}

// APIBorrowRecordsGet 处理GET /api/borrow-records
func APIBorrowRecordsGet(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Query("user_id")
	bookIDStr := c.Query("book_id")
	activeOnly := c.Query("active_only") == "true"

	var records []*models.BorrowRecord

	// 根据参数获取不同的记录
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"message": "无效的用户ID",
			})
			return
		}

		if activeOnly {
			records = models.GetActiveBorrowRecordsByUserID(userID)
		} else {
			records = models.GetBorrowRecordsByUserID(userID)
		}
	} else if bookIDStr != "" {
		bookID, err := strconv.Atoi(bookIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"message": "无效的图书ID",
			})
			return
		}

		if activeOnly {
			records = models.GetActiveBorrowRecordsByBookID(bookID)
		} else {
			records = models.GetBorrowRecordsByBookID(bookID)
		}
	} else {
		if activeOnly {
			records = models.GetAllActiveBorrowRecords()
		} else {
			records = models.GetAllBorrowRecords()
		}
	}

	// 获取用户和图书的映射
	userMap := make(map[int]string)
	bookMap := make(map[int]string)

	for _, record := range records {
		// 添加用户名
		if _, exists := userMap[record.UserID]; !exists {
			user, err := models.GetUserByID(record.UserID)
			if err == nil {
				userMap[record.UserID] = user.Username
			} else {
				userMap[record.UserID] = "未知用户"
			}
		}

		// 添加图书标题
		if _, exists := bookMap[record.BookID]; !exists {
			book, err := models.GetBookByID(record.BookID)
			if err == nil {
				bookMap[record.BookID] = book.Title
			} else {
				bookMap[record.BookID] = "未知图书"
			}
		}
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"records": records,
			"users":   userMap,
			"books":   bookMap,
		},
	})
}