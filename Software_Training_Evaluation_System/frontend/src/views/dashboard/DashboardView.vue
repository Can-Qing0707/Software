<template>
  <div class="app-container">
    <div class="stat-cards">
      <el-card shadow="hover">
        <div class="stat-item">
          <div class="stat-value">{{ stats.courses }}</div>
          <div class="stat-label">课程总数</div>
        </div>
      </el-card>
      <el-card shadow="hover" v-if="isTeacher">
        <div class="stat-item">
          <div class="stat-value">{{ stats.submissions }}</div>
          <div class="stat-label">提交总数</div>
        </div>
      </el-card>
      <el-card shadow="hover" v-else>
        <div class="stat-item">
          <div class="stat-value">{{ stats.tasks }}</div>
          <div class="stat-label">任务总数</div>
        </div>
      </el-card>
      <el-card shadow="hover">
        <div class="stat-item">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">{{ isTeacher ? '待评分' : '待完成' }}</div>
        </div>
      </el-card>
      <el-card shadow="hover">
        <div class="stat-item">
          <div class="stat-value">{{ stats.evaluated }}</div>
          <div class="stat-label">{{ isTeacher ? '已评分' : '已完成' }}</div>
        </div>
      </el-card>
    </div>
    <el-row :gutter="16">
      <el-col :span="14">
        <el-card>
          <template #header>{{ isTeacher ? '近期任务' : '待完成任务' }}</template>
          <el-table :data="tableTasks" stripe v-loading="taskLoading">
            <el-table-column prop="title" label="任务名称" min-width="160" />
            <el-table-column prop="course_name" label="所属课程" width="140" />
            <el-table-column prop="deadline" label="截止时间" width="160" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="isExpired(row.deadline) ? 'danger' : 'success'" size="small">
                  {{ isExpired(row.deadline) ? '已截止' : '进行中' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!tableTasks.length && !taskLoading" description="暂无数据" />
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card>
          <template #header>{{ isTeacher ? '评分概览' : '完成概览' }}</template>
          <div style="height:260px">
            <v-chart :option="chartOption" autoresize />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getCourseList, getMyCourses } from '@/api/course'
import { getTaskList } from '@/api/task'
import { getSubmissionList } from '@/api/submission'
import { useUserStore } from '@/stores/user'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([PieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const userStore = useUserStore()
const isTeacher = computed(() => userStore.role === 'teacher')
const isStudent = computed(() => userStore.role === 'student')

const stats = ref({ courses: 0, tasks: 0, submissions: 0, evaluated: 0, pending: 0 })
const allTasks = ref([])
const allSubmissions = ref([])
const taskLoading = ref(false)

function isExpired(deadline) {
  return deadline && new Date(deadline) < new Date()
}

const tableTasks = computed(() => {
  if (isTeacher.value) {
    return allTasks.value.slice(0, 5)
  }
  const submittedIds = new Set(allSubmissions.value.map(s => s.task_id))
  return allTasks.value.filter(t => !submittedIds.has(t.id)).slice(0, 5)
})

const chartOption = computed(() => {
  if (isTeacher.value) {
    return {
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0, orient: 'horizontal' },
      series: [{
        type: 'pie',
        radius: ['40%', '65%'],
        avoidLabelOverlap: true,
        label: { show: false },
        data: [
          { value: stats.value.evaluated, name: '已评分', itemStyle: { color: '#67c23a' } },
          { value: stats.value.pending, name: '待评分', itemStyle: { color: '#e6a23c' } }
        ]
      }]
    }
  }
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, orient: 'horizontal' },
    series: [{
      type: 'pie',
      radius: ['40%', '65%'],
      avoidLabelOverlap: true,
      label: { show: false },
      data: [
        { value: stats.value.evaluated, name: '已完成', itemStyle: { color: '#67c23a' } },
        { value: stats.value.pending, name: '待完成', itemStyle: { color: '#e6a23c' } }
      ]
    }]
  }
})

onMounted(async () => {
  taskLoading.value = true
  try {
    const [courses, tasks, submissions] = await Promise.all([
      isStudent.value ? getMyCourses() : getCourseList(),
      getTaskList(),
      getSubmissionList()
    ])
    const courseList = courses.data || []
    const taskList = tasks.data || []
    const submissionList = submissions.data || []
    allTasks.value = taskList
    allSubmissions.value = submissionList

    stats.value.courses = courseList.length
    stats.value.tasks = taskList.length
    stats.value.submissions = submissionList.length
    stats.value.evaluated = submissionList.filter(s => s.status === 'evaluated').length

    if (isTeacher.value) {
      stats.value.pending = submissionList.length - stats.value.evaluated
    } else {
      const submittedTaskIds = new Set(submissionList.map(s => s.task_id))
      stats.value.evaluated = submittedTaskIds.size
      stats.value.pending = Math.max(0, taskList.length - submittedTaskIds.size)
    }
  } catch {
    // ignore
  } finally {
    taskLoading.value = false
  }
})
</script>

<style scoped>
.stat-item {
  text-align: center;
  padding: 8px 0;
}
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
}
.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}
</style>
