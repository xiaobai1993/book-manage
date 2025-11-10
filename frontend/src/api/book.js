import request from '@/utils/request'

// 添加图书（管理员）
export function addBook(data) {
  return request({
    url: '/book/add',
    method: 'post',
    data
  })
}

// 编辑图书（管理员）
export function editBook(data) {
  return request({
    url: '/book/edit',
    method: 'post',
    data
  })
}

// 删除图书（管理员）
export function deleteBook(data) {
  return request({
    url: '/book/delete',
    method: 'post',
    data
  })
}

// 获取图书详情
export function getBookDetail(data) {
  return request({
    url: '/book/detail',
    method: 'post',
    data
  })
}

// 图书搜索
export function searchBooks(data) {
  return request({
    url: '/book/search',
    method: 'post',
    data
  })
}

// 上传图书封面（管理员）
export function uploadCover(formData) {
  return request({
    url: '/book/uploadCover',
    method: 'post',
    data: formData
    // 注意：不要设置 Content-Type，让浏览器自动设置 multipart/form-data
  })
}

// 删除图书封面（管理员）
export function deleteCover(data) {
  return request({
    url: '/book/deleteCover',
    method: 'post',
    data
  })
}
