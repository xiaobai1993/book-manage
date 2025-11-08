import request from '@/utils/request'

// 用户注册
export function register(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

// 用户登录
export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

// 发送邮箱验证码
export function sendEmailCode(data) {
  return request({
    url: '/user/sendEmailCode',
    method: 'post',
    data
  })
}

// 密码找回
export function forgetPassword(data) {
  return request({
    url: '/user/forgetPassword',
    method: 'post',
    data
  })
}

// 获取个人信息
export function getProfile(data) {
  return request({
    url: '/user/profile',
    method: 'post',
    data
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/user/changePassword',
    method: 'post',
    data
  })
}

// 获取个人借阅记录
export function getBorrowRecords(data) {
  return request({
    url: '/user/borrowRecords',
    method: 'post',
    data
  })
}
