/**
 * 管理员功能JavaScript文件
 */

document.addEventListener('DOMContentLoaded', function() {
    // 图书表单验证
    setupBookFormValidation();
    
    // 用户管理功能
    setupUserManagement();
    
    // 设置删除确认
    setupDeleteConfirmation();
});

/**
 * 设置图书表单验证
 */
function setupBookFormValidation() {
    // 查找图书表单
    const bookForm = document.querySelector('form.book-form');
    if (bookForm) {
        bookForm.addEventListener('submit', function(event) {
            if (!validateBookForm(event)) {
                event.preventDefault();
            }
        });
    }
}

/**
 * 验证图书表单
 * @param {Event} event - 表单提交事件
 * @returns {boolean} 表单是否有效
 */
function validateBookForm(event) {
    const form = event.target;
    let isValid = true;
    
    // 验证标题
    const titleInput = form.querySelector('input[name="title"]');
    if (titleInput && !titleInput.value.trim()) {
        showError(titleInput, '标题不能为空');
        isValid = false;
    }
    
    // 验证作者
    const authorInput = form.querySelector('input[name="author"]');
    if (authorInput && !authorInput.value.trim()) {
        showError(authorInput, '作者不能为空');
        isValid = false;
    }
    
    // 验证ISBN
    const isbnInput = form.querySelector('input[name="isbn"]');
    if (isbnInput) {
        const isbnValue = isbnInput.value.trim();
        if (!isbnValue) {
            showError(isbnInput, 'ISBN不能为空');
            isValid = false;
        } else if (!/^[0-9-]{10,17}$/.test(isbnValue)) {
            showError(isbnInput, 'ISBN格式无效');
            isValid = false;
        }
    }
    
    // 验证出版年份
    const yearInput = form.querySelector('input[name="published_year"]');
    if (yearInput) {
        const year = parseInt(yearInput.value);
        const currentYear = new Date().getFullYear();
        if (isNaN(year) || year < 1000 || year > currentYear) {
            showError(yearInput, `年份必须在1000到${currentYear}之间`);
            isValid = false;
        }
    }
    
    // 验证数量
    const quantityInput = form.querySelector('input[name="quantity"]');
    if (quantityInput) {
        const quantity = parseInt(quantityInput.value);
        if (isNaN(quantity) || quantity < 1) {
            showError(quantityInput, '数量必须大于0');
            isValid = false;
        }
    }
    
    return isValid;
}

/**
 * 显示表单错误
 * @param {HTMLElement} input - 输入元素
 * @param {string} message - 错误消息
 */
function showError(input, message) {
    input.classList.add('is-invalid');
    
    // 查找或创建错误消息元素
    let feedbackElement = input.nextElementSibling;
    if (!feedbackElement || !feedbackElement.classList.contains('invalid-feedback')) {
        feedbackElement = document.createElement('div');
        feedbackElement.className = 'invalid-feedback';
        input.parentNode.insertBefore(feedbackElement, input.nextSibling);
    }
    
    feedbackElement.textContent = message;
    
    // 添加输入事件监听器以清除错误
    input.addEventListener('input', function() {
        input.classList.remove('is-invalid');
    }, { once: true });
}

/**
 * 设置用户管理功能
 */
function setupUserManagement() {
    // 用户角色变更
    const roleChangeLinks = document.querySelectorAll('.role-change-link');
    roleChangeLinks.forEach(function(link) {
        link.addEventListener('click', function(event) {
            if (!confirm('确定要更改该用户的角色吗？')) {
                event.preventDefault();
            }
        });
    });
    
    // 用户筛选
    const userFilterInput = document.getElementById('user-filter');
    if (userFilterInput) {
        userFilterInput.addEventListener('input', function() {
            const filterText = this.value.toLowerCase();
            const userRows = document.querySelectorAll('table.user-table tbody tr');
            
            userRows.forEach(function(row) {
                const username = row.querySelector('td:nth-child(2)').textContent.toLowerCase();
                const email = row.querySelector('td:nth-child(3)').textContent.toLowerCase();
                
                if (username.includes(filterText) || email.includes(filterText)) {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            });
        });
    }
}

/**
 * 设置删除确认
 */
function setupDeleteConfirmation() {
    const deleteLinks = document.querySelectorAll('.delete-link');
    deleteLinks.forEach(function(link) {
        link.addEventListener('click', function(event) {
            if (!confirm('确定要删除吗？此操作不可恢复！')) {
                event.preventDefault();
            }
        });
    });
}