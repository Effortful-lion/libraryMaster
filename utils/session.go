
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

// Session 会话操作接口
type Session interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Delete(key string)
	Clear()
	Save()
}

// MemorySession 内存会话实现
type MemorySession struct {
	id     string
	data   map[string]interface{}
	expiry time.Time
}

// SessionStore 会话存储接口
type SessionStore interface {
	Get(sessionID string) (map[string]interface{}, error)
	Save(sessionID string, data map[string]interface{}, expiry time.Time) error
	Delete(sessionID string) error
	ClearExpired() error
}

// SessionManager 会话管理结构
type SessionManager struct {
	store SessionStore
	ctx   *gin.Context
}

// NewSessionManager 创建会话管理器
func NewSessionManager(c *gin.Context) *SessionManager {
	return &SessionManager{
		store: &MemorySessionStore{
			sessions: make(map[string]*sessionItem),
		},
		ctx: c,
	}
}

type sessionItem struct {
	data   map[string]interface{}
	expiry time.Time
}

// MemorySessionStore 内存会话存储实现
type MemorySessionStore struct {
	sessions map[string]*sessionItem
}

// 会话存储相关方法实现
func (s *MemorySessionStore) Get(sessionID string) (map[string]interface{}, error) {
	if item, exists := s.sessions[sessionID]; exists && time.Now().Before(item.expiry) {
		return item.data, nil
	}
	return nil, fmt.Errorf("session not found or expired")
}

func (s *MemorySessionStore) Save(sessionID string, data map[string]interface{}, expiry time.Time) error {
	s.sessions[sessionID] = &sessionItem{
		data:   data,
		expiry: expiry,
	}
	return nil
}

func (s *MemorySessionStore) Delete(sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

func (s *MemorySessionStore) ClearExpired() error {
	now := time.Now()
	for id, session := range s.sessions {
		if now.After(session.expiry) {
			delete(s.sessions, id)
		}
	}
	return nil
}

// MemorySession方法实现
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
	// 实现保存逻辑
}

// GetSession 获取当前会话
func (sm *SessionManager) GetSession(c *gin.Context) Session {
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		sessionID = generateSessionID()
		c.SetCookie("session_id", sessionID, 86400, "/", "", false, true)
	}

	data, err := sm.store.Get(sessionID)
	if err != nil {
		data = make(map[string]interface{})
	}

	session := &MemorySession{
		id:     sessionID,
		data:   data,
		expiry: time.Now().Add(24 * time.Hour),
	}

	return session
}

// SaveUserToSession 保存用户信息到会话
func (sm *SessionManager) SaveUserToSession(c *gin.Context, userID int, username, role string) {
	session := sm.GetSession(c)
	session.Set(KeyUserID, userID)
	session.Set(KeyUsername, username)
	session.Set(KeyUserRole, role)
	session.Set(KeyLoggedIn, true)
	session.Save()
}

// ClearSession 清除会话
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

// GenerateCSRFToken 生成CSRF令牌
func (sm *SessionManager) GenerateCSRFToken(c *gin.Context) string {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.StdEncoding.EncodeToString(b)

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

	valid := storedToken.(string) == token
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

	session.Delete(fullKey)
	session.Save()

	return message.(string)
}

// generateSessionID 生成随机会话ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
// CleanupExpiredSessions 清理过期的会话
func CleanupExpiredSessions() {
	// 获取默认存储实例
	store := &MemorySessionStore{
		sessions: make(map[string]*sessionItem),
	}
	
	// 调用清理方法
	if err := store.ClearExpired(); err != nil {
		log.Printf("清理过期会话时出错: %v\n", err)
	}
}
