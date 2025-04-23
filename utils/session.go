package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 常量
const (
	// Session相关键名
	KeyUserID       = "user_id"
	KeyUsername     = "username"
	KeyUserRole     = "user_role"
	KeyLoggedIn     = "logged_in"
	KeyCSRFToken    = "csrf_token"
	KeyFlashMessage = "flash_message"
)

// SaveUserToSession 保存用户信息到会话
func (sm *SessionManager) SaveUserToSession(c *gin.Context, userID int, username, role string) {
	// 设置会话值
	session := sm.GetSession(c)
	session.Set(KeyUserID, userID)
	session.Set(KeyUsername, username)
	session.Set(KeyUserRole, role)
	session.Set(KeyLoggedIn, true)
	session.Save()
}

// ClearSession 清除会话，用于退出登录
func (sm *SessionManager) ClearSession(c *gin.Context) {
	session := sm.GetSession(c)
	session.Clear()
	session.Save()
}

// IsLoggedIn 检查用户是否已登录
func (sm *SessionManager) IsLoggedIn(c *gin.Context) bool {
	session := sm.GetSession(c)
	loggedIn, exists := session.Get(KeyLoggedIn)
	if !exists {
		return false
	}
	return loggedIn.(bool)
}

// GetUserIDFromSession 从会话中获取用户ID
func (sm *SessionManager) GetUserIDFromSession(c *gin.Context) int {
	session := sm.GetSession(c)
	userID, exists := session.Get(KeyUserID)
	if !exists {
		return 0
	}
	return userID.(int)
}

// GetUsernameFromSession 从会话中获取用户名
func (sm *SessionManager) GetUsernameFromSession(c *gin.Context) string {
	session := sm.GetSession(c)
	username, exists := session.Get(KeyUsername)
	if !exists {
		return ""
	}
	return username.(string)
}

// GetUserRoleFromSession 从会话中获取用户角色
func (sm *SessionManager) GetUserRoleFromSession(c *gin.Context) string {
	session := sm.GetSession(c)
	role, exists := session.Get(KeyUserRole)
	if !exists {
		return ""
	}
	return role.(string)
}

// GenerateCSRFToken 生成CSRF令牌并保存到会话
func (sm *SessionManager) GenerateCSRFToken(c *gin.Context) string {
	// 生成随机令牌
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.StdEncoding.EncodeToString(b)

	// 保存到会话
	session := sm.GetSession(c)
	session.Set(KeyCSRFToken, token)
	session.Save()

	return token
}

// VerifyCSRFToken 验证CSRF令牌
func (sm *SessionManager) VerifyCSRFToken(c *gin.Context, token string) bool {
	session := sm.GetSession(c)
	storedToken, exists := session.Get(KeyCSRFToken)
	if !exists {
		return false
	}

	// 验证令牌
	valid := storedToken.(string) == token

	// 验证后清除令牌，防止重放攻击
	if valid {
		session.Delete(KeyCSRFToken)
		session.Save()
	}

	return valid
}

// SetFlashMessage 设置闪存消息
func (sm *SessionManager) SetFlashMessage(c *gin.Context, key, message string) {
	session := sm.GetSession(c)
	session.Set(fmt.Sprintf("%s_%s", KeyFlashMessage, key), message)
	session.Save()
}

// GetFlashMessage 获取并清除闪存消息
func (sm *SessionManager) GetFlashMessage(c *gin.Context, key string) string {
	session := sm.GetSession(c)
	fullKey := fmt.Sprintf("%s_%s", KeyFlashMessage, key)
	message, exists := session.Get(fullKey)

	if !exists {
		return ""
	}

	// 获取后清除消息
	session.Delete(fullKey)
	session.Save()

	return message.(string)
}

// SessionStore 会话存储接口
type SessionStore interface {
	Get(sessionID string) (map[string]interface{}, error)
	Save(sessionID string, data map[string]interface{}, expiry time.Time) error
	Delete(sessionID string) error
	ClearExpired() error
}

// MemorySessionStore 内存会话存储实现
type MemorySessionStore struct {
	sessions map[string]*sessionItem
}

type sessionItem struct {
	data   map[string]interface{}
	expiry time.Time
}

// Session 会话操作接口
type Session interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Delete(key string)
	Clear()
	Save()
}

// SessionManager 会话管理结构
type SessionManager struct {
	store     SessionStore
	ctx       *gin.Context
	encryptor DataEncryptor
}

// NewSessionManager 创建会话管理器
func NewSessionManager(c *gin.Context, store SessionStore, encryptor DataEncryptor) *SessionManager {
	return &SessionManager{
		store:     store,
		ctx:       c,
		encryptor: encryptor,
	}
}

// 全局会话存储
var sessions = make(map[string]*MemorySession)

// GetSessionData 获取会话数据
func (sm *SessionManager) GetSessionData() map[string]interface{} {
	session := sm.GetSession(sm.ctx)
	return session.GetAll()
}

// GetAll 获取所有会话数据
func (s *MemorySession) GetAll() map[string]interface{} {
	return s.data
}

// GetSession 获取当前会话
func (sm *SessionManager) GetSession(c *gin.Context) Session {
	// 尝试从Cookie获取会话ID
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		// 创建新会话
		sessionID = generateSessionID()
		c.SetCookie("session_id", sessionID, 86400, "/", "", false, true)
	}

	// 检查会话是否存在或已过期
	session, exists := sessions[sessionID]
	if !exists || time.Now().After(session.expiry) {
		// 创建新会话
		session = &MemorySession{
			data:   make(map[string]interface{}),
			id:     sessionID,
			expiry: time.Now().Add(24 * time.Hour),
		}
		sessionStore.Save(sessionID, session.data, session.expiry)
	} else {
		// 延长会话有效期
		session.expiry = time.Now().Add(24 * time.Hour)
	}

	return session
}

// 会话方法实现
func (s *MemorySession) Get(key string) (interface{}, bool) {
	val, exists := s.data[key]
	return val, exists
}

func (s *MemorySession) Set(key string, val interface{}) {
	s.data[key] = val
}

func (s *MemorySession) Delete(key string) {
	delete(s.data, key)
}

func (s *MemorySession) Clear() {
	s.data = make(map[string]interface{})
}

func (s *MemorySession) Save() {
	// 内存会话无需保存
	sessionStore.Save(s.id, s.data, s.expiry)
}

// generateSessionID 生成随机会话ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// CleanupExpiredSessions 清理过期会话
func CleanupExpiredSessions() {
	now := time.Now()
	for id, session := range sessions {
		if now.After(session.expiry) {
			delete(sessions, id)
		}
	}
}
