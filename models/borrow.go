package models

import (
	"errors"
	"sync"
	"time"
)

// BorrowRecord 借阅记录模型
type BorrowRecord struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	BookID     int       `json:"book_id"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
	ReturnDate time.Time `json:"return_date"`
}

// BorrowRecords 全局借阅记录列表
var (
	BorrowRecords  []*BorrowRecord
	NextBorrowID   = 1
	borrowMutex    sync.Mutex
)

// CreateBorrowRecord 创建新借阅记录
func CreateBorrowRecord(userID, bookID int, borrowDate, dueDate time.Time) (*BorrowRecord, error) {
	borrowMutex.Lock()
	defer borrowMutex.Unlock()
	
	// 验证用户是否存在
	user, err := GetUserByID(userID)
	if user == nil && err != nil {
		return nil, errors.New("用户不存在")
	}
	
	// 验证图书是否存在
	book, err := GetBookByID(bookID)
	if err != nil {
		return nil, errors.New("图书不存在")
	}
	
	// 检查图书是否可借
	if !book.IsAvailable() {
		return nil, errors.New("该图书无可用库存")
	}
	
	// 检查用户是否已借阅该图书
	for _, record := range BorrowRecords {
		if record.UserID == userID && record.BookID == bookID && record.ReturnDate.IsZero() {
			return nil, errors.New("您已借阅该图书，尚未归还")
		}
	}
	
	// 创建借阅记录
	record := &BorrowRecord{
		ID:         NextBorrowID,
		UserID:     userID,
		BookID:     bookID,
		BorrowDate: borrowDate,
		DueDate:    dueDate,
	}
	
	// 添加到列表并递增ID
	BorrowRecords = append(BorrowRecords, record)
	NextBorrowID++
	
	return record, nil
}

// ReturnBook 归还图书
func ReturnBook(recordID int) (*BorrowRecord, error) {
	borrowMutex.Lock()
	defer borrowMutex.Unlock()
	
	// 查找借阅记录
	var record *BorrowRecord
	for _, r := range BorrowRecords {
		if r.ID == recordID {
			record = r
			break
		}
	}
	
	if record == nil {
		return nil, errors.New("借阅记录不存在")
	}
	
	// 检查是否已归还
	if !record.ReturnDate.IsZero() {
		return nil, errors.New("该图书已归还")
	}
	
	// 更新归还日期
	record.ReturnDate = time.Now()
	
	return record, nil
}

// GetBorrowRecordByID 根据ID获取借阅记录
func GetBorrowRecordByID(id int) (*BorrowRecord, error) {
	for _, record := range BorrowRecords {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("借阅记录不存在")
}

// GetAllBorrowRecords 获取所有借阅记录
func GetAllBorrowRecords() []*BorrowRecord {
	return BorrowRecords
}

// GetBorrowRecordsByUserID 获取用户的所有借阅记录
func GetBorrowRecordsByUserID(userID int) []*BorrowRecord {
	var records []*BorrowRecord
	for _, record := range BorrowRecords {
		if record.UserID == userID {
			records = append(records, record)
		}
	}
	return records
}

// GetBorrowRecordsByBookID 获取图书的所有借阅记录
func GetBorrowRecordsByBookID(bookID int) []*BorrowRecord {
	var records []*BorrowRecord
	for _, record := range BorrowRecords {
		if record.BookID == bookID {
			records = append(records, record)
		}
	}
	return records
}

// GetActiveBorrowRecordsByUserID 获取用户的未归还借阅记录
func GetActiveBorrowRecordsByUserID(userID int) []*BorrowRecord {
	var records []*BorrowRecord
	for _, record := range BorrowRecords {
		if record.UserID == userID && record.ReturnDate.IsZero() {
			records = append(records, record)
		}
	}
	return records
}

// GetActiveBorrowRecordsByBookID 获取图书的未归还借阅记录
func GetActiveBorrowRecordsByBookID(bookID int) []*BorrowRecord {
	var records []*BorrowRecord
	for _, record := range BorrowRecords {
		if record.BookID == bookID && record.ReturnDate.IsZero() {
			records = append(records, record)
		}
	}
	return records
}

// GetAllActiveBorrowRecords 获取所有未归还的借阅记录
func GetAllActiveBorrowRecords() []*BorrowRecord {
	var records []*BorrowRecord
	for _, record := range BorrowRecords {
		if record.ReturnDate.IsZero() {
			records = append(records, record)
		}
	}
	return records
}

// GetAllOverdueBorrowRecords 获取所有逾期的借阅记录
func GetAllOverdueBorrowRecords() []*BorrowRecord {
	var records []*BorrowRecord
	now := time.Now()
	for _, record := range BorrowRecords {
		if record.ReturnDate.IsZero() && now.After(record.DueDate) {
			records = append(records, record)
		}
	}
	return records
}

// IsOverdue 检查借阅记录是否逾期
func (r *BorrowRecord) IsOverdue() bool {
	return r.ReturnDate.IsZero() && time.Now().After(r.DueDate)
}

// DaysUntilDue 计算距离到期日还有多少天
func (r *BorrowRecord) DaysUntilDue() int {
	if r.ReturnDate.IsZero() {
		now := time.Now()
		if now.After(r.DueDate) {
			return 0 // 已逾期
		}
		hours := r.DueDate.Sub(now).Hours()
		return int(hours / 24)
	}
	return 0 // 已归还
}

// OverdueDays 计算逾期天数
func (r *BorrowRecord) OverdueDays() int {
	if r.ReturnDate.IsZero() {
		now := time.Now()
		if now.After(r.DueDate) {
			hours := now.Sub(r.DueDate).Hours()
			return int(hours / 24)
		}
	} else if r.ReturnDate.After(r.DueDate) {
		hours := r.ReturnDate.Sub(r.DueDate).Hours()
		return int(hours / 24)
	}
	return 0 // 未逾期或已归还且未逾期
}

// InitSampleBorrowRecords 初始化示例借阅记录
func InitSampleBorrowRecords() {
	// 清空借阅记录列表
	borrowMutex.Lock()
	defer borrowMutex.Unlock()
	
	BorrowRecords = nil
	NextBorrowID = 1
	
	// 确保有足够的用户和图书
	if len(Users) < 3 || len(Books) < 3 {
		return
	}
	
	// 设置时间
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)
	twoWeeksAgo := now.AddDate(0, 0, -14)
	inOneWeek := now.AddDate(0, 0, 7)
	inTwoWeeks := now.AddDate(0, 0, 14)
	overdueDate := now.AddDate(0, 0, -1) // 昨天到期
	
	// 管理员的借阅记录
	record1 := &BorrowRecord{
		ID:         NextBorrowID,
		UserID:     1, // 管理员
		BookID:     1, // 第一本书
		BorrowDate: twoWeeksAgo,
		DueDate:    twoWeeksAgo.AddDate(0, 0, 14),
		ReturnDate: oneWeekAgo,
	}
	BorrowRecords = append(BorrowRecords, record1)
	NextBorrowID++
	
	// 图书管理员的借阅记录
	record2 := &BorrowRecord{
		ID:         NextBorrowID,
		UserID:     2, // 图书管理员
		BookID:     2, // 第二本书
		BorrowDate: oneWeekAgo,
		DueDate:    inOneWeek,
	}
	BorrowRecords = append(BorrowRecords, record2)
	NextBorrowID++
	
	// 读者的借阅记录 (未到期)
	record3 := &BorrowRecord{
		ID:         NextBorrowID,
		UserID:     3, // 读者
		BookID:     3, // 第三本书
		BorrowDate: oneWeekAgo,
		DueDate:    inTwoWeeks,
	}
	BorrowRecords = append(BorrowRecords, record3)
	NextBorrowID++
	
	// 读者的借阅记录 (已逾期)
	record4 := &BorrowRecord{
		ID:         NextBorrowID,
		UserID:     3, // 读者
		BookID:     4, // 第四本书
		BorrowDate: twoWeeksAgo,
		DueDate:    overdueDate,
	}
	BorrowRecords = append(BorrowRecords, record4)
	NextBorrowID++
}