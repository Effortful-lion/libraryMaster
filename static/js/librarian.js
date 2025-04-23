/**
 * 图书管理员功能JavaScript文件
 */

document.addEventListener('DOMContentLoaded', function() {
    // 初始化库存管理功能
    setupInventoryManagement();
    
    // 初始化借阅管理功能
    setupBorrowManagement();
    
    // 高亮显示逾期图书
    highlightOverdueBooks();
});

/**
 * 初始化库存管理功能
 */
function setupInventoryManagement() {
    // 库存筛选功能
    const inventoryFilterInput = document.getElementById('inventory-filter');
    if (inventoryFilterInput) {
        inventoryFilterInput.addEventListener('input', function() {
            const filterText = this.value.toLowerCase();
            filterInventory(filterText);
        });
    }
    
    // 库存排序功能
    const inventorySortSelect = document.getElementById('inventory-sort');
    if (inventorySortSelect) {
        inventorySortSelect.addEventListener('change', function() {
            sortInventory(this.value);
        });
    }
}

/**
 * 筛选库存
 * @param {string} filterText - 筛选文本
 */
function filterInventory(filterText) {
    const bookItems = document.querySelectorAll('.book-item');
    
    bookItems.forEach(function(item) {
        const title = item.querySelector('.book-title').textContent.toLowerCase();
        const author = item.querySelector('.book-author').textContent.toLowerCase();
        const isbn = item.querySelector('.book-isbn').textContent.toLowerCase();
        
        if (title.includes(filterText) || author.includes(filterText) || isbn.includes(filterText)) {
            item.style.display = '';
        } else {
            item.style.display = 'none';
        }
    });
}

/**
 * 排序库存
 * @param {string} sortOption - 排序选项
 */
function sortInventory(sortOption) {
    const bookList = document.querySelector('.book-list');
    if (!bookList) return;
    
    const bookItems = Array.from(bookList.querySelectorAll('.book-item'));
    
    // 根据选项排序
    bookItems.sort(function(a, b) {
        switch (sortOption) {
            case 'title-asc':
                return a.querySelector('.book-title').textContent.localeCompare(b.querySelector('.book-title').textContent);
            case 'title-desc':
                return b.querySelector('.book-title').textContent.localeCompare(a.querySelector('.book-title').textContent);
            case 'available-asc':
                return parseInt(a.dataset.available) - parseInt(b.dataset.available);
            case 'available-desc':
                return parseInt(b.dataset.available) - parseInt(a.dataset.available);
            default:
                return 0;
        }
    });
    
    // 重新添加到DOM
    bookItems.forEach(function(item) {
        bookList.appendChild(item);
    });
}

/**
 * 初始化借阅管理功能
 */
function setupBorrowManagement() {
    // 借阅记录筛选功能
    const borrowFilterInput = document.getElementById('borrow-filter');
    if (borrowFilterInput) {
        borrowFilterInput.addEventListener('input', function() {
            const filterText = this.value.toLowerCase();
            filterBorrowRecords(filterText);
        });
    }
    
    // 借阅状态筛选
    const statusFilterSelect = document.getElementById('status-filter');
    if (statusFilterSelect) {
        statusFilterSelect.addEventListener('change', function() {
            filterBorrowRecordsByStatus(this.value);
        });
    }
    
    // 借阅表单用户选择
    const userSelect = document.getElementById('user_id');
    if (userSelect) {
        userSelect.addEventListener('change', function() {
            updateBorrowInfo();
        });
    }
    
    // 借阅表单图书选择
    const bookSelect = document.getElementById('book_id');
    if (bookSelect) {
        bookSelect.addEventListener('change', function() {
            updateBorrowInfo();
        });
    }
}

/**
 * 过滤借阅记录
 * @param {string} filterText - 筛选文本
 */
function filterBorrowRecords(filterText) {
    const borrowRows = document.querySelectorAll('table.borrow-table tbody tr');
    
    borrowRows.forEach(function(row) {
        const username = row.querySelector('.borrow-username').textContent.toLowerCase();
        const bookTitle = row.querySelector('.borrow-book-title').textContent.toLowerCase();
        
        if (username.includes(filterText) || bookTitle.includes(filterText)) {
            row.style.display = '';
        } else {
            row.style.display = 'none';
        }
    });
}

/**
 * 按状态过滤借阅记录
 * @param {string} status - 状态值 (all, active, returned, overdue)
 */
function filterBorrowRecordsByStatus(status) {
    const borrowRows = document.querySelectorAll('table.borrow-table tbody tr');
    
    borrowRows.forEach(function(row) {
        if (status === 'all') {
            row.style.display = '';
            return;
        }
        
        const isActive = !row.querySelector('.borrow-return-date').textContent.trim();
        const isOverdue = row.classList.contains('overdue');
        
        switch (status) {
            case 'active':
                row.style.display = isActive && !isOverdue ? '' : 'none';
                break;
            case 'returned':
                row.style.display = !isActive ? '' : 'none';
                break;
            case 'overdue':
                row.style.display = isOverdue ? '' : 'none';
                break;
        }
    });
}

/**
 * 更新借阅信息
 */
function updateBorrowInfo() {
    const userSelect = document.getElementById('user_id');
    const bookSelect = document.getElementById('book_id');
    const borrowInfoDiv = document.getElementById('borrow-info');
    
    if (!userSelect || !bookSelect || !borrowInfoDiv) return;
    
    const selectedUser = userSelect.options[userSelect.selectedIndex];
    const selectedBook = bookSelect.options[bookSelect.selectedIndex];
    
    if (selectedUser && selectedBook) {
        const userName = selectedUser.textContent;
        const bookTitle = selectedBook.textContent;
        const availableCount = parseInt(selectedBook.dataset.available || '0');
        
        let infoHTML = `
            <div class="alert ${availableCount > 0 ? 'alert-info' : 'alert-warning'}">
                <p><strong>借阅信息</strong></p>
                <p>用户: ${userName}</p>
                <p>图书: ${bookTitle}</p>
                <p>可借数量: ${availableCount}</p>
            </div>
        `;
        
        borrowInfoDiv.innerHTML = infoHTML;
        
        // 禁用提交按钮，如果没有可用数量
        const submitButton = document.querySelector('button[type="submit"]');
        if (submitButton) {
            submitButton.disabled = availableCount <= 0;
        }
    }
}

/**
 * 高亮显示逾期图书
 */
function highlightOverdueBooks() {
    const today = new Date();
    
    // 检查所有借阅记录行
    document.querySelectorAll('table.borrow-table tbody tr').forEach(function(row) {
        const dueDateCell = row.querySelector('.borrow-due-date');
        const returnDateCell = row.querySelector('.borrow-return-date');
        
        if (!dueDateCell) return;
        
        const dueDate = new Date(dueDateCell.textContent);
        const returnDate = returnDateCell && returnDateCell.textContent.trim() 
            ? new Date(returnDateCell.textContent) 
            : null;
        
        // 如果尚未归还且已过期
        if (!returnDate && today > dueDate) {
            row.classList.add('overdue');
            row.classList.add('table-danger');
            
            // 计算逾期天数
            const overdueDays = Math.floor((today - dueDate) / (1000 * 60 * 60 * 24));
            const statusCell = row.querySelector('.borrow-status');
            if (statusCell) {
                statusCell.innerHTML = `<span class="badge bg-danger">逾期 ${overdueDays} 天</span>`;
            }
        }
    });
}