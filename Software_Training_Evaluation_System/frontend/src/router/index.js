import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { noAuth: true }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' }
      },
      {
        path: 'course',
        name: 'CourseList',
        component: () => import('@/views/course/CourseList.vue'),
        meta: { title: '课程管理', icon: 'Reading' }
      },
      {
        path: 'course/:id',
        name: 'CourseDetail',
        component: () => import('@/views/course/CourseDetail.vue'),
        meta: { title: '课程详情', hidden: true }
      },
      {
        path: 'task',
        name: 'TaskList',
        component: () => import('@/views/task/TaskList.vue'),
        meta: { title: '实训任务', icon: 'List' }
      },
      {
        path: 'task/:id',
        name: 'TaskDetail',
        component: () => import('@/views/task/TaskDetail.vue'),
        meta: { title: '任务详情', hidden: true }
      },
      {
        path: 'submission',
        name: 'SubmissionList',
        component: () => import('@/views/submission/SubmissionList.vue'),
        meta: { title: '成果提交', icon: 'Upload' }
      },
      {
        path: 'submission/upload/:taskId?',
        name: 'SubmissionUpload',
        component: () => import('@/views/submission/SubmissionUpload.vue'),
        meta: { title: '提交成果', hidden: true }
      },
      {
        path: 'submission/:id',
        name: 'SubmissionDetail',
        component: () => import('@/views/submission/SubmissionDetail.vue'),
        meta: { title: '提交详情', hidden: true }
      },
      {
        path: 'evaluation/indicator',
        name: 'IndicatorManage',
        component: () => import('@/views/evaluation/IndicatorManage.vue'),
        meta: { title: '评价指标', icon: 'SetUp' }
      },
      {
        path: 'evaluation/score/:submissionId',
        name: 'ScoreDetail',
        component: () => import('@/views/evaluation/ScoreDetail.vue'),
        meta: { title: '评分详情', hidden: true }
      },
      {
        path: 'report',
        name: 'ReportView',
        component: () => import('@/views/report/ReportView.vue'),
        meta: { title: '报表管理', icon: 'DataAnalysis' }
      },
      {
        path: 'user',
        name: 'UserManage',
        component: () => import('@/views/user/UserManage.vue'),
        meta: { title: '用户管理', icon: 'User', role: 'admin' }
      },
      {
        path: 'system/llm',
        name: 'LlmConfig',
        component: () => import('@/views/system/LlmConfig.vue'),
        meta: { title: 'LLM配置', icon: 'Setting', role: 'admin' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.noAuth) {
    next()
  } else if (!token) {
    next('/login')
  } else {
    next()
  }
})

export default router
