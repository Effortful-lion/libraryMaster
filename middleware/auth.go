package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"librarysystem/utils"
)

// RequireAuth 验证用户是否登录
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已登录
		mg := utils.NewSessionManager(c)
		if !mg.IsLoggedIn(c) {
			// 存储当前URL以便登录后重定向
			if c.Request.Method == "GET" {
				mg.SetFlashMessage(c, "redirect", c.Request.URL.Path)
			}
			
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "请先登录")
			
			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// 继续处理请求
		c.Next()
	}
}

// RequireAdmin 验证用户是否为管理员
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已登录
		mg := utils.NewSessionManager(c)
		if !mg.IsLoggedIn(c) {
			// 存储当前URL以便登录后重定向
			if c.Request.Method == "GET" {
				mg.SetFlashMessage(c, "redirect", c.Request.URL.Path)
			}
			
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "请先登录")
			
			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// 检查用户角色
		userRole := mg.GetUserRoleFromSession(c)
		if userRole != "admin" {
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "需要管理员权限")
			
			// 重定向到仪表板
			c.Redirect(http.StatusFound, "/dashboard")
			c.Abort()
			return
		}
		
		// 继续处理请求
		c.Next()
	}
}

// RequireLibrarian 验证用户是否为图书管理员
func RequireLibrarian() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已登录
		mg := utils.NewSessionManager(c)
		if !mg.IsLoggedIn(c) {
			// 存储当前URL以便登录后重定向
			if c.Request.Method == "GET" {
				mg.SetFlashMessage(c, "redirect", c.Request.URL.Path)
			}
			
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "请先登录")
			
			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// 检查用户角色
		userRole := mg.GetUserRoleFromSession(c)
		if userRole != "librarian" && userRole != "admin" {
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "需要图书管理员权限")
			
			// 重定向到仪表板
			c.Redirect(http.StatusFound, "/dashboard")
			c.Abort()
			return
		}
		
		// 继续处理请求
		c.Next()
	}
}

// RequireReader 验证用户是否为读者
func RequireReader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已登录
		mg := utils.NewSessionManager(c)
		if !mg.IsLoggedIn(c) {
			// 存储当前URL以便登录后重定向
			if c.Request.Method == "GET" {
				mg.SetFlashMessage(c, "redirect", c.Request.URL.Path)
			}
			
			// 设置错误信息
			mg.SetFlashMessage(c, "error", "请先登录")
			
			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// 所有登录用户都可以访问读者页面
		// 继续处理请求
		c.Next()
	}
}