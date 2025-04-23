package models

import (
	"errors"
	"strings"
	"sync"
)

// Book 图书模型
type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	PublishedYear int    `json:"published_year"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	CoverURL      string `json:"cover_url"`
	Quantity      int    `json:"quantity"`
}

// Books 全局图书列表
var (
	Books       []*Book
	NextBookID  = 1
	bookMutex   sync.Mutex
	categories  map[string]bool
)

// CreateBook 创建新图书
func CreateBook(title, author, isbn string, publishedYear int, category, description, coverURL string, quantity int) (*Book, error) {
	bookMutex.Lock()
	defer bookMutex.Unlock()
	
	// 验证参数
	if title == "" || author == "" || isbn == "" || category == "" || description == "" || coverURL == "" {
		return nil, errors.New("所有字段都不能为空")
	}
	
	if publishedYear < 1000 || publishedYear > 2100 {
		return nil, errors.New("出版年份必须在1000-2100之间")
	}
	
	if quantity < 1 {
		return nil, errors.New("数量必须大于0")
	}
	
	// 检查ISBN是否已存在
	for _, book := range Books {
		if book.ISBN == isbn {
			return nil, errors.New("ISBN已存在")
		}
	}
	
	// 创建图书
	book := &Book{
		ID:            NextBookID,
		Title:         title,
		Author:        author,
		ISBN:          isbn,
		PublishedYear: publishedYear,
		Category:      category,
		Description:   description,
		CoverURL:      coverURL,
		Quantity:      quantity,
	}
	
	// 添加到分类映射
	if categories == nil {
		categories = make(map[string]bool)
	}
	categories[category] = true
	
	// 添加到列表并递增ID
	Books = append(Books, book)
	NextBookID++
	
	return book, nil
}

// GetBookByID 根据ID获取图书
func GetBookByID(id int) (*Book, error) {
	for _, book := range Books {
		if book.ID == id {
			return book, nil
		}
	}
	return nil, errors.New("图书不存在")
}

// GetAllBooks 获取所有图书
func GetAllBooks() []*Book {
	return Books
}

// GetAllCategories 获取所有分类
func GetAllCategories() []string {
	if categories == nil {
		return []string{}
	}
	
	var result []string
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// GetBooksByCategory 根据分类获取图书
func GetBooksByCategory(category string) []*Book {
	var result []*Book
	for _, book := range Books {
		if book.Category == category {
			result = append(result, book)
		}
	}
	return result
}

// SearchBooks 搜索图书
func SearchBooks(query string) []*Book {
	query = strings.ToLower(query)
	var result []*Book
	
	for _, book := range Books {
		// 匹配标题、作者、ISBN或描述
		if strings.Contains(strings.ToLower(book.Title), query) ||
		   strings.Contains(strings.ToLower(book.Author), query) ||
		   strings.Contains(strings.ToLower(book.ISBN), query) ||
		   strings.Contains(strings.ToLower(book.Description), query) {
			result = append(result, book)
		}
	}
	
	return result
}

// UpdateBook 更新图书
func UpdateBook(id int, title, author, isbn string, publishedYear int, category, description, coverURL string, quantity int) (*Book, error) {
	bookMutex.Lock()
	defer bookMutex.Unlock()
	
	// 验证参数
	if title == "" || author == "" || isbn == "" || category == "" || description == "" || coverURL == "" {
		return nil, errors.New("所有字段都不能为空")
	}
	
	if publishedYear < 1000 || publishedYear > 2100 {
		return nil, errors.New("出版年份必须在1000-2100之间")
	}
	
	if quantity < 1 {
		return nil, errors.New("数量必须大于0")
	}
	
	// 查找图书
	var book *Book
	for _, b := range Books {
		if b.ID == id {
			book = b
			break
		}
	}
	
	if book == nil {
		return nil, errors.New("图书不存在")
	}
	
	// 检查ISBN是否已被其他图书使用
	for _, b := range Books {
		if b.ID != id && b.ISBN == isbn {
			return nil, errors.New("ISBN已被其他图书使用")
		}
	}
	
	// 更新图书信息
	book.Title = title
	book.Author = author
	book.ISBN = isbn
	book.PublishedYear = publishedYear
	book.Category = category
	book.Description = description
	book.CoverURL = coverURL
	book.Quantity = quantity
	
	// 添加到分类映射
	if categories == nil {
		categories = make(map[string]bool)
	}
	categories[category] = true
	
	return book, nil
}

// DeleteBook 删除图书
func DeleteBook(id int) error {
	bookMutex.Lock()
	defer bookMutex.Unlock()
	
	// 检查有无借阅记录
	for _, record := range BorrowRecords {
		if record.BookID == id && record.ReturnDate.IsZero() {
			return errors.New("该图书有未归还的借阅记录，无法删除")
		}
	}
	
	// 查找图书索引
	index := -1
	for i, book := range Books {
		if book.ID == id {
			index = i
			break
		}
	}
	
	if index == -1 {
		return errors.New("图书不存在")
	}
	
	// 删除图书
	Books = append(Books[:index], Books[index+1:]...)
	
	// 重建分类映射
	categories = make(map[string]bool)
	for _, book := range Books {
		categories[book.Category] = true
	}
	
	return nil
}

// IsAvailable 检查图书是否有可用库存
func (b *Book) IsAvailable() bool {
	return b.GetAvailableQuantity() > 0
}

// GetAvailableQuantity 获取可用库存数量
func (b *Book) GetAvailableQuantity() int {
	// 计算当前借出的数量
	borrowedCount := 0
	for _, record := range BorrowRecords {
		if record.BookID == b.ID && record.ReturnDate.IsZero() {
			borrowedCount++
		}
	}
	
	// 返回可用数量
	return b.Quantity - borrowedCount
}

// InitSampleBooks 初始化示例图书数据
func InitSampleBooks() {
	// 清空图书列表
	bookMutex.Lock()
	defer bookMutex.Unlock()
	
	Books = nil
	NextBookID = 1
	categories = make(map[string]bool)
	
	// 创建示例图书
	books := []struct {
		Title         string
		Author        string
		ISBN          string
		PublishedYear int
		Category      string
		Description   string
		CoverURL      string
		Quantity      int
	}{
		{
			Title:         "Python编程：从入门到实践",
			Author:        "埃里克·马瑟斯",
			ISBN:          "9787115428028",
			PublishedYear: 2016,
			Category:      "编程",
			Description:   "本书是一本针对所有层次的Python读者而作的Python入门书。全书分两部分：第一部分介绍用Python编程所必须了解的基本概念，第二部分将理论付诸实践，讲解如何开发三个项目。",
			CoverURL:      "https://img3.doubanio.com/view/subject/l/public/s29424065.jpg",
			Quantity:      5,
		},
		{
			Title:         "Go语言实战",
			Author:        "威廉·肯尼迪",
			ISBN:          "9787115445353",
			PublishedYear: 2017,
			Category:      "编程",
			Description:   "本书首先介绍Go语言的独特之处，然后讲解如何编写地道的Go代码并使用其特有的特性和工具包编写代码。后续章节会介绍测试、Web编程以及与其他主流语言的集成。",
			CoverURL:      "https://img9.doubanio.com/view/subject/l/public/s29446435.jpg",
			Quantity:      3,
		},
		{
			Title:         "明朝那些事儿",
			Author:        "当年明月",
			ISBN:          "9787807023630",
			PublishedYear: 2009,
			Category:      "历史",
			Description:   "《明朝那些事儿》讲述从1344年到1644年，明朝三百年间的历史。以史料为基础，以年代和具体人物为主线，运用小说的手法，对明朝十七帝和其他王公权贵和小人物的命运进行全景展示。",
			CoverURL:      "https://img1.doubanio.com/view/subject/l/public/s27131114.jpg",
			Quantity:      4,
		},
		{
			Title:         "三体",
			Author:        "刘慈欣",
			ISBN:          "9787536692930",
			PublishedYear: 2008,
			Category:      "科幻",
			Description:   "文化大革命如火如荼进行的同时，军方探寻外星文明的绝秘计划'红岸工程'取得了突破性进展。但在按下发射键的那一刻，历经劫难的叶文洁没有意识到，她彻底改变了人类的命运。",
			CoverURL:      "https://img2.doubanio.com/view/subject/l/public/s2768378.jpg",
			Quantity:      2,
		},
		{
			Title:         "围城",
			Author:        "钱钟书",
			ISBN:          "9787020090006",
			PublishedYear: 1991,
			Category:      "文学",
			Description:   "《围城》是钱钟书所著的长篇小说，自问世以来，就以它的犀利的语言、巧妙的结构和象征性的意义在中国文学史上占据重要地位。",
			CoverURL:      "https://img2.doubanio.com/view/subject/l/public/s1070222.jpg",
			Quantity:      3,
		},
	}
	
	for _, bookData := range books {
		book := &Book{
			ID:            NextBookID,
			Title:         bookData.Title,
			Author:        bookData.Author,
			ISBN:          bookData.ISBN,
			PublishedYear: bookData.PublishedYear,
			Category:      bookData.Category,
			Description:   bookData.Description,
			CoverURL:      bookData.CoverURL,
			Quantity:      bookData.Quantity,
		}
		
		// 添加到列表并递增ID
		Books = append(Books, book)
		NextBookID++
		
		// 添加到分类映射
		categories[bookData.Category] = true
	}
}