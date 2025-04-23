/**
 * 图书管理系统主JavaScript文件
 */

document.addEventListener('DOMContentLoaded', function() {
    // 初始化提示工具
    initTooltips();
    
    // 初始化日期格式化
    formatDates();
    
    // 设置闪存消息自动关闭
    setupFlashMessages();
    
    // 表单验证
    setupFormValidation();
});

/**
 * 初始化Bootstrap提示工具
 */
function initTooltips() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function(tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
}

/**
 * 初始化日期格式化
 */
function formatDates() {
    document.querySelectorAll('.format-date').forEach(function(element) {
        const date = new Date(element.textContent);
        if (!isNaN(date)) {
            element.textContent = formatDate(date);
        }
    });
    
    document.querySelectorAll('.format-datetime').forEach(function(element) {
        const date = new Date(element.textContent);
        if (!isNaN(date)) {
            element.textContent = formatDateTime(date);
        }
    });
}

/**
 * 设置闪存消息自动关闭
 */
function setupFlashMessages() {
    document.querySelectorAll('.alert-dismissible').forEach(function(alert) {
        setTimeout(function() {
            const closeButton = alert.querySelector('.btn-close');
            if (closeButton) {
                closeButton.click();
            }
        }, 5000); // 5秒后自动关闭
    });
}

/**
 * 设置表单验证
 */
function setupFormValidation() {
    document.querySelectorAll('form').forEach(function(form) {
        form.addEventListener('submit', function(event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
                
                // 标记所有必填字段
                form.querySelectorAll('input, select, textarea').forEach(function(input) {
                    if (input.hasAttribute('required') && !input.value.trim()) {
                        input.classList.add('is-invalid');
                    } else {
                        input.classList.remove('is-invalid');
                    }
                });
            }
            
            form.classList.add('was-validated');
        });
        
        // 实时验证
        form.querySelectorAll('input, select, textarea').forEach(function(input) {
            input.addEventListener('input', function() {
                if (input.checkValidity()) {
                    input.classList.remove('is-invalid');
                    input.classList.add('is-valid');
                } else {
                    input.classList.remove('is-valid');
                    input.classList.add('is-invalid');
                }
            });
        });
    });
}

/**
 * 格式化日期为YYYY-MM-DD格式
 * @param {Date} date - 日期对象
 * @returns {string} 格式化后的日期字符串
 */
function formatDate(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}

/**
 * 格式化日期时间为YYYY-MM-DD HH:MM:SS格式
 * @param {Date} date - 日期对象
 * @returns {string} 格式化后的日期时间字符串
 */
function formatDateTime(date) {
    const formattedDate = formatDate(date);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    return `${formattedDate} ${hours}:${minutes}:${seconds}`;
}

/**
 * 计算两个日期之间的天数
 * @param {Date} date1 - 第一个日期
 * @param {Date} date2 - 第二个日期
 * @returns {number} 两个日期之间的天数
 */
function daysBetween(date1, date2) {
    const oneDay = 24 * 60 * 60 * 1000; // 一天的毫秒数
    // 取两个日期的时间戳，并计算差值
    const diffTime = Math.abs(date2.getTime() - date1.getTime());
    return Math.round(diffTime / oneDay);
}

/**
 * 计算逾期天数
 * @param {string} dueDate - 到期日期字符串
 * @returns {number} 逾期天数，如未逾期返回0
 */
function calculateOverdueDays(dueDate) {
    const today = new Date();
    const due = new Date(dueDate);
    
    // 如果今天在到期日之后，则计算天数
    if (today > due) {
        return daysBetween(due, today);
    }
    
    return 0;
}

/**
 * 检查是否逾期
 * @param {string} dueDate - 到期日期字符串
 * @returns {boolean} 是否逾期
 */
function isOverdue(dueDate) {
    return calculateOverdueDays(dueDate) > 0;
}

/**
 * 显示或隐藏加载指示器
 * @param {boolean} show - 是否显示加载指示器
 */
function toggleLoading(show) {
    const loadingElement = document.getElementById('loading-indicator');
    
    if (!loadingElement) {
        // 创建加载指示器
        const loading = document.createElement('div');
        loading.id = 'loading-indicator';
        loading.className = 'position-fixed w-100 h-100 d-flex justify-content-center align-items-center bg-dark bg-opacity-50';
        loading.style.top = '0';
        loading.style.left = '0';
        loading.style.zIndex = '9999';
        
        const spinner = document.createElement('div');
        spinner.className = 'spinner-border text-light';
        spinner.setAttribute('role', 'status');
        
        const srOnly = document.createElement('span');
        srOnly.className = 'visually-hidden';
        srOnly.textContent = '加载中...';
        
        spinner.appendChild(srOnly);
        loading.appendChild(spinner);
        
        document.body.appendChild(loading);
    } else {
        loadingElement.style.display = show ? 'flex' : 'none';
    }
}

/**
 * 显示通知消息
 * @param {string} message - 消息内容
 * @param {string} type - 消息类型（success, error, warning, info）
 * @param {number} duration - 持续时间（毫秒）
 */
function showNotification(message, type = 'info', duration = 3000) {
    // 创建通知容器（如果不存在）
    let container = document.getElementById('notification-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'notification-container';
        container.className = 'position-fixed top-0 end-0 p-3';
        container.style.zIndex = '9999';
        document.body.appendChild(container);
    }
    
    // 设置样式映射
    const typeClass = {
        success: 'bg-success',
        error: 'bg-danger',
        warning: 'bg-warning',
        info: 'bg-info'
    };
    
    // 创建通知元素
    const notification = document.createElement('div');
    notification.className = `toast show ${typeClass[type] || typeClass.info} text-white`;
    notification.setAttribute('role', 'alert');
    notification.setAttribute('aria-live', 'assertive');
    notification.setAttribute('aria-atomic', 'true');
    
    const notificationBody = document.createElement('div');
    notificationBody.className = 'toast-body d-flex align-items-center';
    
    // 添加图标
    const icon = document.createElement('i');
    icon.className = 'me-2 fas ' + 
        (type === 'success' ? 'fa-check-circle' : 
        (type === 'error' ? 'fa-exclamation-circle' : 
        (type === 'warning' ? 'fa-exclamation-triangle' : 'fa-info-circle')));
    
    notificationBody.appendChild(icon);
    notificationBody.appendChild(document.createTextNode(message));
    
    // 添加关闭按钮
    const closeButton = document.createElement('button');
    closeButton.type = 'button';
    closeButton.className = 'btn-close btn-close-white ms-auto';
    closeButton.setAttribute('data-bs-dismiss', 'toast');
    closeButton.setAttribute('aria-label', '关闭');
    
    notificationBody.appendChild(closeButton);
    notification.appendChild(notificationBody);
    
    // 添加到容器
    container.appendChild(notification);
    
    // 设置自动关闭
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            container.removeChild(notification);
        }, 300);
    }, duration);
    
    // 点击关闭按钮时移除通知
    closeButton.addEventListener('click', () => {
        notification.classList.remove('show');
        setTimeout(() => {
            container.removeChild(notification);
        }, 300);
    });
}