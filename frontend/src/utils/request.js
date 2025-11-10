import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useUserStore } from '@/stores/user'

// 根据环境变量设置 API 基础地址
// 开发环境使用代理，生产环境使用实际的后端地址
const getBaseURL = () => {
  // 如果设置了 VITE_API_BASE_URL，使用它（生产环境）
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL
  }
  // 开发环境使用代理
  return '/api'
}

const service = axios.create({
  baseURL: getBaseURL(),
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
    
    // 如果是 FormData（文件上传），不设置 Content-Type，让浏览器自动设置
    if (config.data instanceof FormData) {
      // FormData 会自动设置 Content-Type 为 multipart/form-data，并包含 boundary
      // 不要手动设置，否则会丢失 boundary
      delete config.headers['Content-Type']
      // FormData 中已经包含了 token，不需要额外处理
    } else if (config.data && typeof config.data === 'object' && token) {
      // 普通对象请求，添加token到请求体
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
