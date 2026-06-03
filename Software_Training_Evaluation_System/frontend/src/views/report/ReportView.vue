<template>
  <div class="app-container">
    <div class="filter-container">
      <el-select v-model="query.course_id" placeholder="选择课程" clearable @change="fetchData" style="width:180px">
        <el-option v-for="c in courses" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
      <el-select v-model="query.task_id" placeholder="选择任务" clearable @change="fetchData" style="width:200px">
        <el-option v-for="t in tasks" :key="t.id" :label="t.title" :value="t.id" />
      </el-select>
      <el-button type="primary" @click="fetchData">搜索</el-button>
      <el-button type="success" @click="generateReport">生成报表</el-button>
    </div>

    <el-card style="margin-bottom:16px">
      <template #header>数据概览</template>
      <div class="stat-cards">
        <el-statistic title="总提交" :value="stats.total" />
        <el-statistic title="已评分" :value="stats.evaluated" />
        <el-statistic title="平均分" :value="stats.avgScore" :precision="1" />
        <el-statistic title="通过率" :value="stats.passRate" suffix="%" :precision="1" />
      </div>
    </el-card>

    <el-row :gutter="16">
      <el-col :span="12">
        <el-card>
          <template #header>评分分布</template>
          <v-chart :option="chartOption" style="height:300px" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>导出记录</template>
          <el-table :data="reports" stripe>
            <el-table-column prop="title" label="报告名称" min-width="160" />
            <el-table-column prop="type" label="类型" width="80">
              <template #default="{ row }">{{ typeMap[row.type] }}</template>
            </el-table-column>
            <el-table-column prop="format" label="格式" width="70">
              <template #default="{ row }">
                <el-tag size="small">{{ row.format.toUpperCase() }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="generated_at" label="生成时间" width="150" />
            <el-table-column label="操作" width="80">
              <template #default="{ row }">
                <el-button type="primary" link @click="downloadReport(row)">下载</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getCourseList } from '@/api/course'
import { getTaskList } from '@/api/task'
import { getSubmissionList } from '@/api/submission'
import { getReports, generateReport as generateReportApi, exportReport } from '@/api/report'
import { ElMessage } from 'element-plus'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([BarChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent, CanvasRenderer])

const courses = ref([])
const tasks = ref([])
const reports = ref([])
const loading = ref(false)
const query = ref({ course_id: '', task_id: '' })
const stats = ref({ total: 0, evaluated: 0, avgScore: 0, passRate: 0 })
const typeMap = { individual: '个人', class: '班级', course: '课程' }

const chartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 50, right: 20, bottom: 30 },
  xAxis: { type: 'category', data: ['优秀(≥90)', '良好(80-89)', '中等(70-79)', '及格(60-69)', '不及格(<60)'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [stats.value.total * 0.2, stats.value.total * 0.4, stats.value.total * 0.25, stats.value.total * 0.1, stats.value.total * 0.05], itemStyle: { color: '#409eff' } }]
}))

async function fetchData() {
  loading.value = true
  try {
    const [subRes, reportRes] = await Promise.all([
      getSubmissionList(query.value),
      getReports(query.value)
    ])
    const data = subRes.data || []
    reports.value = reportRes.data || []
    stats.value.total = data.length
    stats.value.evaluated = data.filter(s => s.status === 'evaluated').length
    const scored = data.filter(s => s.final_score != null)
    stats.value.avgScore = scored.length ? scored.reduce((a, b) => a + (b.final_score || 0), 0) / scored.length : 0
    stats.value.passRate = scored.length ? (scored.filter(s => (s.final_score || 0) >= 60).length / scored.length * 100) : 0
  } finally {
    loading.value = false
  }
}

async function generateReport() {
  try {
      await generateReportApi(query.value)
    ElMessage.success('报表生成中，请稍后刷新')
    await fetchData()
  } catch {
    // handled by interceptor
  }
}

async function downloadReport(row) {
  try {
    const blob = await exportReport(row.id, row.format)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${row.title}.${row.format}`
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    // handled by interceptor
  }
}

onMounted(async () => {
  const [courseRes, taskRes] = await Promise.all([getCourseList(), getTaskList()])
  courses.value = courseRes.data || []
  tasks.value = taskRes.data || []
  await fetchData()
})
</script>
