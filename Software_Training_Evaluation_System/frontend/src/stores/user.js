import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => userInfo.value?.role || '')
  const username = computed(() => userInfo.value?.username || '')
  const realName = computed(() => userInfo.value?.real_name || '')

  async function doLogin(loginForm) {
    const res = await login(loginForm)
    token.value = res.data.token
    userInfo.value = res.data.user
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('user', JSON.stringify(res.data.user))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchUserInfo() {
    const res = await getUserInfo()
    userInfo.value = res.data
    localStorage.setItem('user', JSON.stringify(res.data))
  }

  return { token, userInfo, isLoggedIn, role, username, realName, doLogin, logout, fetchUserInfo }
})
