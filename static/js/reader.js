/**
 * 读者功能JavaScript文件
 */

document.addEventListener('DOMContentLoaded', function() {
    // 设置借阅搜索和过滤功能
    const borrowedBooksFilter = document.getElementById('borrowed-books-filter');
    if (borrowedBooksFilter) {
        borrowedBooksFilter.addEventListener('input', function() {
            filterBorrowedBooks(this.value.toLowerCase());
        });
    }
    
    // 高亮显示读者的逾期图书
    highlightReaderOverdueBooks();
    
    // 初始化书籍搜索和过滤功能
    initBookSearch();
});

/**
 * 过滤已借阅图书
 * @param {string} searchTerm - 搜索关键词
 */
function filterBorrowedBooks(searchTerm) {
    const borrowedItems = document.querySelectorAll('.borrowed-book-item');
    
    borrowedItems.forEach(function(item) {
        const title = item.querySelector('.book-title').textContent.toLowerCase();
        const author = item.querySelector('.book-author').textContent.toLowerCase();
        
        if (title.includes(searchTerm) || author.includes(searchTerm)) {
            item.style.display = '';
        } else {
            item.style.display = 'none';
        }
    });
}

/**
 * 高亮显示读者的逾期图书
 */
function highlightReaderOverdueBooks() {
    const today = new Date();
    
    document.querySelectorAll('.borrowed-book-item').forEach(function(item) {
        const dueDateElement = item.querySelector('.due-date');
        const returnDateElement = item.querySelector('.return-date');
        
        if (!dueDateElement) return;
        
        const dueDate = new Date(dueDateElement.textContent);
        const isReturned = returnDateElement && returnDateElement.textContent.trim() !== '';
        
        if (!isReturned && today > dueDate) {
            // 已逾期但未归还
            item.classList.add('border-danger');
            
            const statusBadge = item.querySelector('.status-badge');
            if (statusBadge) {
                // 计算逾期天数
                const overdueDays = Math.floor((today - dueDate) / (1000 * 60 * 60 * 24));
                statusBadge.textContent = `逾期 ${overdueDays} 天`;
                statusBadge.className = 'badge bg-danger status-badge';
            }
        }
    });
}

/**
 * 初始化书籍搜索和过滤功能
 */
function initBookSearch() {
    // 图书搜索
    const bookSearchForm = document.getElementById('book-search-form');
    if (bookSearchForm) {
        bookSearchForm.addEventListener('submit', function(event) {
            const searchInput = document.getElementById('book-search-input');
            if (searchInput && !searchInput.value.trim()) {
                event.preventDefault();
                alert('请输入搜索内容');
            }
        });
    }
    
    // 图书过滤
    const categoryFilter = document.getElementById('category-filter');
    if (categoryFilter) {
        categoryFilter.addEventListener('change', function() {
            filterBooksByCategory(this.value);
        });
    }
}

/**
 * 按分类过滤图书
 * @param {string} category - 分类名称
 */
function filterBooksByCategory(category) {
    const bookItems = document.querySelectorAll('.book-item');
    
    bookItems.forEach(function(item) {
        if (category === 'all') {
            item.style.display = '';
            return;
        }
        
        const bookCategory = item.dataset.category;
        if (bookCategory === category) {
            item.style.display = '';
        } else {
            item.style.display = 'none';
        }
    });
}