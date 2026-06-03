<template>
  <el-menu
    :default-active="activeMenu"
    :collapse="appStore.sidebarCollapsed"
    background-color="#304156"
    text-color="#bfcbd9"
    active-text-color="#409eff"
    router
    unique-opened
  >
    <div class="logo-area">
      <span class="logo-text" v-show="!appStore.sidebarCollapsed">实训评价系统</span>
      <span class="logo-short" v-show="appStore.sidebarCollapsed">实训</span>
    </div>
    <template v-for="item in filteredMenus" :key="item.path">
      <el-menu-item :index="item.path">
        <el-icon><component :is="item.icon" /></el-icon>
        <template #title>{{ item.title }}</template>
      </el-menu-item>
    </template>
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()

const menuItems = [
  { path: '/dashboard', title: '仪表盘', icon: 'Odometer', roles: ['admin', 'teacher', 'student'] },
  { path: '/course', title: '课程管理', icon: 'Reading', roles: ['admin', 'teacher', 'student'] },
  { path: '/task', title: '实训任务', icon: 'List', roles: ['admin', 'teacher'] },
  { path: '/submission', title: '实训任务', icon: 'List', roles: ['student'] },
  { path: '/submission', title: '成果提交', icon: 'Upload', roles: ['teacher'] },
  { path: '/evaluation/indicator', title: '评价指标', icon: 'SetUp', roles: ['admin'] },
  { path: '/report', title: '报表管理', icon: 'DataAnalysis', roles: ['admin', 'teacher', 'student'] },
  { path: '/user', title: '用户管理', icon: 'User', roles: ['admin'] },
  { path: '/system/llm', title: 'LLM配置', icon: 'Setting', roles: ['admin'] }
]

const filteredMenus = computed(() => {
  const role = userStore.role
  return menuItems.filter(item => item.roles.includes(role))
})

const activeMenu = computed(() => route.path)
</script>

<style scoped>
.logo-area {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: bold;
  font-size: 16px;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}
.logo-short {
  font-size: 14px;
}
.el-menu {
  border-right: none;
}
</style>
