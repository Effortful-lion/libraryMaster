package controllers

import (
	"encoding/json"
	"librarysystem/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// APIBooksGet 获取所有图书
func APIBooksGet(c *gin.Context) {
	books := models.Books
	c.JSON(http.StatusOK, books)
}

// APIBookGet 获取单本图书
func APIBookGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	book, err := models.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// APICategoriesGet 获取所有图书分类
func APICategoriesGet(c *gin.Context) {
	categories := models.GetCategories()
	c.JSON(http.StatusOK, categories)
}

// APIBorrowRecordsGet 获取借阅记录
func APIBorrowRecordsGet(c *gin.Context) {
	records := models.BorrowRecords

	// 创建用户名和图书标题映射
	userMap := make(map[int]string)
	bookMap := make(map[int]string)

	// 遍历记录填充映射
	for _, record := range records {
		if _, exists := userMap[record.UserID]; !exists {
			user, err := models.GetUserByID(record.UserID)
			if err == nil {
				userMap[record.UserID] = user.Username
			} else {
				userMap[record.UserID] = "未知用户"
			}
		}

		if _, exists := bookMap[record.BookID]; !exists {
			book, err := models.GetBookByID(record.BookID)
			if err == nil {
				bookMap[record.BookID] = book.Title
			} else {
				bookMap[record.BookID] = "未知图书"
			}
		}
	}

	// 构建响应
	type RecordResponse struct {
		*models.BorrowRecord
		UserName  string `json:"username"`
		BookTitle string `json:"book_title"`
	}

	var response []RecordResponse
	for _, record := range records {
		response = append(response, RecordResponse{
			BorrowRecord: record,
			UserName:     userMap[record.UserID],
			BookTitle:    bookMap[record.BookID],
		})
	}

	c.JSON(http.StatusOK, response)
}