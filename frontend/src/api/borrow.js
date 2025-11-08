import request from '@/utils/request'

// 借书
export function borrowBook(data) {
  return request({
    url: '/borrow/borrow',
    method: 'post',
    data
  })
}

// 还书
export function returnBook(data) {
  return request({
    url: '/borrow/return',
    method: 'post',
    data
  })
}

// 获取个人借阅记录
export function getBorrowRecords(data) {
  return request({
    url: '/borrow/records',
    method: 'post',
    data
  })
}

// 获取全量借阅记录（管理员）
export function getAllRecords(data) {
  return request({
    url: '/borrow/allRecords',
    method: 'post',
    data
  })
}
