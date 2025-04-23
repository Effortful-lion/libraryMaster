package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"librarysystem/models"
	"librarysystem/utils"
)

// AdminUsers 管理员用户管理页面
func AdminUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mg := utils.NewSessionManager(c)
		session := mg.GetSession(c)

		// 查询所有用户
		var users []models.User
		db.Order("id ASC").Find(&users)
		
		// 设置模板数据
		data := gin.H{
			"title":       "用户管理",
			"users":       users,
			"currentYear": c.MustGet("currentYear"),
			"username":    mg.GetUsernameFromSession(c),
			"user_role":   mg.GetUserRoleFromSession(c),
		}
		
		c.HTML(http.StatusOK, "admin/users.html", data)
	}
}

// ChangeUserRole 更改用户角色
func ChangeUserRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"error":       "无效的用户ID",
				"title":       "错误",
				"currentYear": c.MustGet("currentYear"),
			})
			return
		}
		
		// 查询用户
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"error":       "找不到用户",
				"title":       "错误",
				"currentYear": c.MustGet("currentYear"),
			})
			return
		}
		
		// 获取当前会话用户ID
		sessionManager := utils.NewSessionManager(c)
		currentUserID := sessionManager.GetUserIDFromSession(c)
		
		// 不允许管理员更改自己的角色
		if uint(userID) == uint(currentUserID) {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":       "不能更改自己的角色",
				"title":       "错误",
				"currentYear": c.MustGet("currentYear"),
			})
			return
		}
		
		// 根据当前角色循环切换角色
		switch user.Role {
		case "reader":
			user.Role = "librarian"
		case "librarian":
			user.Role = "admin"
		case "admin":
			user.Role = "reader"
		}
		
		// 保存更改
		if err := db.Save(&user).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error":       "更新用户角色时出错",
				"title":       "错误",
				"currentYear": c.MustGet("currentYear"),
			})
			return
		}
		
		// 重定向到用户管理页面
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}
}