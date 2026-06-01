import { reactive } from 'vue'

// 简单的全局认证状态（不依赖 Pinia，足够这个项目使用）
const state = reactive({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  token: localStorage.getItem('token') || '',
})

export function useAuth() {
  return state
}

export function isLoggedIn() {
  return !!state.token
}

export function login(token, user) {
  state.token = token
  state.user = user
  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
}

export function logout() {
  state.token = ''
  state.user = null
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}
