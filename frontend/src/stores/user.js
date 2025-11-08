import { defineStore } from 'pinia'
import { login, getProfile } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null'),
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.userInfo?.role === 'admin',
    email: (state) => state.userInfo?.email || '',
  },

  actions: {
    async login(loginData) {
      try {
        const res = await login(loginData)
        this.token = res.data.token
        this.userInfo = res.data.user_info
        
        localStorage.setItem('token', this.token)
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        
        return res
      } catch (error) {
        throw error
      }
    },

    async fetchProfile() {
      try {
        const res = await getProfile({ token: this.token })
        this.userInfo = res.data.user_info
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        return res
      } catch (error) {
        throw error
      }
    },

    logout() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    },
  },
})
