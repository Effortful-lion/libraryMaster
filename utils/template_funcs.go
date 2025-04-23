package utils

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// TemplateFunctions 返回可在模板中使用的自定义函数
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		// 格式化日期时间
		"formatDate": formatDate,
		"formatDateTime": formatDateTime,
		
		// 日期计算
		"daysBetween": daysBetween,
		"daysAgo": daysAgo,
		"daysFromNow": daysFromNow,
		
		// 字符串操作
		"truncate": truncate,
		"contains": strings.Contains,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		
		// 数值操作
		"add": add,
		"subtract": subtract,
		"multiply": multiply,
		"divide": divide,
		
		// 条件判断
		"eq": eq,
		"ne": ne,
		"lt": lt,
		"lte": lte,
		"gt": gt,
		"gte": gte,
		
		// 切片和映射操作
		"join": strings.Join,
		"split": strings.Split,
		
		// 时间判断
		"isOverdue": isOverdue,
	}
}

// 格式化日期: 2006-01-02
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 格式化日期时间: 2006-01-02 15:04:05
func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 计算两个日期之间的天数
func daysBetween(t1, t2 time.Time) int {
	hours := t2.Sub(t1).Hours()
	return int(hours / 24)
}

// 计算距离今天的天数（过去）
func daysAgo(t time.Time) int {
	return daysBetween(t, time.Now())
}

// 计算距离今天的天数（未来）
func daysFromNow(t time.Time) int {
	return daysBetween(time.Now(), t)
}

// 截断字符串
func truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// 数值操作
func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// 条件判断
func eq(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func ne(a, b interface{}) bool {
	return !eq(a, b)
}

func lt(a, b int) bool {
	return a < b
}

func lte(a, b int) bool {
	return a <= b
}

func gt(a, b int) bool {
	return a > b
}

func gte(a, b int) bool {
	return a >= b
}

// 判断日期是否已过期（相对于当前时间）
func isOverdue(t time.Time) bool {
	return time.Now().After(t)
}