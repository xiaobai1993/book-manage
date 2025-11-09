import request from '@/utils/request'

// 获取验证码记录列表
export function getEmailCodeList(data) {
  return request({
    url: '/admin/emailCodeList',
    method: 'post',
    data
  })
}

// 获取验证码统计信息
export function getEmailCodeStats() {
  return request({
    url: '/admin/emailCodeStats',
    method: 'post',
    data: {}
  })
}

