package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"librarysystem/models"
	"librarysystem/utils"
)

// LoginForm 登录表单结构
type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	CSRFToken string `form:"csrf_token"`
}

// RegisterForm 注册表单结构
type RegisterForm struct {
	Username        string `form:"username" binding:"required,min=3,max=20"`
	Email           string `form:"email" binding:"required,email"`
	Password        string `form:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
	CSRFToken       string `form:"csrf_token"`
}

// LoginGet 处理GET /login
func LoginGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 如果用户已登录，重定向到仪表板
	if mg.IsLoggedIn(c) {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	// 生成CSRF令牌
	token := mg.GenerateCSRFToken(c)

	// 获取错误消息，如果有
	errorMsg := mg.GetFlashMessage(c, "error")
	
	// 渲染登录页面
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":      "用户登录",
		"csrf_token": token,
		"error":      errorMsg,
	})
}

// LoginPost 处理POST /login
func LoginPost(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		// 表单验证失败
		mg.SetFlashMessage(c, "error", "用户名和密码不能为空")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 验证CSRF令牌
	if !mg.VerifyCSRFToken(c, form.CSRFToken) {
		mg.SetFlashMessage(c, "error", "安全验证失败，请重试")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 获取用户
	user, err := models.GetUserByUsername(form.Username)
	if err != nil || !user.CheckPassword(form.Password) {
		// 用户名或密码错误
		mg.SetFlashMessage(c, "error", "用户名或密码错误")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 登录成功，保存会话
	mg.SaveUserToSession(c, user.ID, user.Username, string(user.Role))
	
	// 设置成功消息
	mg.SetFlashMessage(c, "success", "登录成功！欢迎回来，"+user.Username)
	
	// 重定向到仪表板
	c.Redirect(http.StatusFound, "/dashboard")
}

// RegisterGet 处理GET /register
func RegisterGet(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 如果用户已登录，重定向到仪表板
	if mg.IsLoggedIn(c) {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	// 生成CSRF令牌
	token := mg.GenerateCSRFToken(c)

	// 获取错误消息，如果有
	errorMsg := mg.GetFlashMessage(c, "error")
	
	// 渲染注册页面
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":      "用户注册",
		"csrf_token": token,
		"error":      errorMsg,
	})
}

// RegisterPost 处理POST /register
func RegisterPost(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		// 表单验证失败
		mg.SetFlashMessage(c, "error", "请填写所有必填字段并确保密码一致")
		c.Redirect(http.StatusFound, "/register")
		return
	}

	// 验证CSRF令牌
	if !mg.VerifyCSRFToken(c, form.CSRFToken) {
		mg.SetFlashMessage(c, "error", "安全验证失败，请重试")
		c.Redirect(http.StatusFound, "/register")
		return
	}

	// 创建新用户
	user, err := models.CreateUser(form.Username, form.Email, form.Password, models.RoleReader)
	if err != nil {
		// 用户创建失败
		mg.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/register")
		return
	}

	// 注册成功，保存会话
	mg.SaveUserToSession(c, user.ID, user.Username, string(user.Role))
	
	// 设置成功消息
	mg.SetFlashMessage(c, "success", "注册成功！欢迎，"+user.Username)
	
	// 重定向到仪表板
	c.Redirect(http.StatusFound, "/dashboard")
}

// Logout 处理GET /logout
func Logout(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 清除会话
	mg.ClearSession(c)
	
	// 设置成功消息
	mg.SetFlashMessage(c, "success", "您已成功退出登录")
	
	// 重定向到首页
	c.Redirect(http.StatusFound, "/")
}

// Dashboard 处理GET /dashboard
func Dashboard(c *gin.Context) {
	mg := utils.NewSessionManager(c)
	// 获取用户ID
	userID := mg.GetUserIDFromSession(c)
	if userID <= 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 获取用户角色
	userRole := mg.GetUserRoleFromSession(c)
	username := mg.GetUsernameFromSession(c)

	// 准备仪表板数据
	data := gin.H{
		"title":     "仪表板",
		"username":  username,
		"user_role": userRole,
		"now":       time.Now(),
	}

	// 根据用户角色添加不同的数据
	switch userRole {
	case "admin":
		// 管理员仪表板数据
		data["user_count"] = len(models.GetAllUsers())
		data["book_count"] = len(models.GetAllBooks())
		data["borrow_count"] = len(models.GetAllBorrowRecords())
	case "librarian":
		// 图书管理员仪表板数据
		data["book_count"] = len(models.GetAllBooks())
		data["active_borrow_count"] = len(models.GetAllActiveBorrowRecords())
		data["overdue_count"] = len(models.GetAllOverdueBorrowRecords())
	default:
		// 读者仪表板数据
		activeRecords := models.GetActiveBorrowRecordsByUserID(userID)
		overdueCount := 0
		for _, record := range activeRecords {
			if record.IsOverdue() {
				overdueCount++
			}
		}
		data["borrowed_books_count"] = len(activeRecords)
		data["overdue_books_count"] = overdueCount
	}

	// 渲染仪表板页面
	c.HTML(http.StatusOK, "dashboard.html", data)
}