import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useUserStore } from '@/stores/user'

const service = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    const token = userStore.token
    
    // 优先使用Header方式传递token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 如果请求数据是对象，添加token到请求体（后端也支持从请求体读取token）
    if (config.data && typeof config.data === 'object' && token) {
      config.data.token = token
    }
    
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 业务状态码判断
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      
      // token过期或无效，跳转到登录页
      if (res.code === 10001 && res.message.includes('token')) {
        const userStore = useUserStore()
        userStore.logout()
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res
  },
  error => {
    console.error('响应错误:', error)
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default service
