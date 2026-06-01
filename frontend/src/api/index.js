import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// 请求拦截器：自动带上 Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器：统一处理错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // 登录和注册接口的 401 不做跳转，让页面自己显示错误
    const url = error.config?.url || ''
    const isAuthPage = url.includes('/login') || url.includes('/register')

    if (error.response?.status === 401 && !isAuthPage) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
      return Promise.reject(new Error('登录已过期，请重新登录'))
    }
    const msg = error.response?.data?.error || '网络错误'
    return Promise.reject(new Error(msg))
  }
)

export default api
