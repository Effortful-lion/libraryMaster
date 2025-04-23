package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// 用户角色
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleLibrarian UserRole = "librarian"
	RoleReader    UserRole = "reader"
)

// User 用户模型
type User struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"-"` // 不在JSON中显示密码哈希
	Role         UserRole `json:"role"`
}

// Users 全局用户列表
var (
	Users      []*User
	NextUserID = 1
	userMutex  sync.Mutex
)

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	hashedPassword := hashPassword(password)
	return u.PasswordHash == hashedPassword
}

// UpdateRole 更新用户角色
func (u *User) UpdateRole(role UserRole) error {
	// 验证角色是否有效
	if role != RoleAdmin && role != RoleLibrarian && role != RoleReader {
		return errors.New("无效的用户角色")
	}
	
	// 更新角色
	u.Role = role
	return nil
}

// CreateUser 创建新用户
func CreateUser(username, email, password string, role UserRole) (*User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()
	
	// 验证参数
	if username == "" || email == "" || password == "" {
		return nil, errors.New("用户名、邮箱和密码不能为空")
	}
	
	// 检查用户名是否已存在
	for _, user := range Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return nil, errors.New("用户名已存在")
		}
	}
	
	// 检查邮箱是否已存在
	for _, user := range Users {
		if strings.ToLower(user.Email) == strings.ToLower(email) {
			return nil, errors.New("邮箱已存在")
		}
	}
	
	// 创建用户
	user := &User{
		ID:           NextUserID,
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password),
		Role:         role,
	}
	
	// 添加到列表并递增ID
	Users = append(Users, user)
	NextUserID++
	
	return user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	for _, user := range Users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	for _, user := range Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	for _, user := range Users {
		if strings.ToLower(user.Email) == strings.ToLower(email) {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

// GetAllUsers 获取所有用户
func GetAllUsers() []*User {
	return Users
}

// 私有函数，用于密码哈希
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// InitSampleUsers 初始化示例用户数据
func InitSampleUsers() {
	// 清空用户列表
	userMutex.Lock()
	defer userMutex.Unlock()
	
	Users = nil
	NextUserID = 1
	
	// 创建管理员
	admin := &User{
		ID:           NextUserID,
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: hashPassword("admin123"),
		Role:         RoleAdmin,
	}
	Users = append(Users, admin)
	NextUserID++
	
	// 创建图书管理员
	librarian := &User{
		ID:           NextUserID,
		Username:     "librarian",
		Email:        "librarian@example.com",
		PasswordHash: hashPassword("librarian123"),
		Role:         RoleLibrarian,
	}
	Users = append(Users, librarian)
	NextUserID++
	
	// 创建读者
	reader := &User{
		ID:           NextUserID,
		Username:     "reader",
		Email:        "reader@example.com",
		PasswordHash: hashPassword("reader123"),
		Role:         RoleReader,
	}
	Users = append(Users, reader)
	NextUserID++
}