<template>
  <div class="app-container">
    <el-button style="margin-bottom:16px" @click="$router.back()">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>
    <el-card v-loading="loading">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>{{ detail.title }}</span>
          <div>
            <el-tag :type="detail.status === 1 ? 'success' : 'info'" size="small" style="margin-right:8px">
              {{ detail.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
            <el-button type="primary" size="small" @click="goUpload" v-role="'student'">提交成果</el-button>
          </div>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="所属课程">{{ detail.course_name }}</el-descriptions-item>
        <el-descriptions-item label="截止时间">{{ detail.deadline || '无' }}</el-descriptions-item>
        <el-descriptions-item label="任务描述" :span="2">
          <pre style="white-space:pre-wrap;margin:0;font-family:inherit">{{ detail.description }}</pre>
        </el-descriptions-item>
      </el-descriptions>

      <h3 style="margin:24px 0 16px">评价指标权重
        <el-button v-if="isTeacherOrAdmin" size="small" type="primary" style="margin-left:12px" @click="showIndicatorDialog">
          配置指标
        </el-button>
      </h3>
      <el-table :data="indicators" stripe v-if="indicators.length">
        <el-table-column prop="name" label="指标" />
        <el-table-column prop="weight" label="权重(%)" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.weight" :stroke-width="16" />
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="暂未配置评价指标" />

      <el-dialog v-model="indicatorDialogVisible" title="配置任务评价指标" width="600px">
        <div style="margin-bottom:12px;color:#909399;font-size:13px">
          从全局指标库中选择本任务适用的评价指标并设置权重，权重总和建议为 100%。
        </div>
        <el-table :data="allIndicators" stripe>
          <el-table-column label="启用" width="60">
            <template #default="{ row }">
              <el-checkbox v-model="row.checked" />
            </template>
          </el-table-column>
          <el-table-column prop="name" label="指标名称" width="140" />
          <el-table-column prop="description" label="说明" min-width="180" />
          <el-table-column label="权重(%)" width="150">
            <template #default="{ row }">
              <el-input-number v-model="row.weight" :min="0" :max="100" :precision="2" :disabled="!row.checked" size="small" style="width:120px" />
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top:8px;font-size:13px" :style="{ color: weightSum === 100 ? '#67c23a' : '#e6a23c' }">
          当前权重总和：{{ weightSum }}%
        </div>
        <template #footer>
          <el-button @click="indicatorDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveIndicators" :loading="savingIndicators">保存</el-button>
        </template>
      </el-dialog>

      <h3 style="margin:24px 0 16px" v-if="isTeacherOrAdmin">学生提交</h3>
      <el-tabs v-if="isTeacherOrAdmin" v-model="submissionTab" type="border-card">
        <el-tab-pane label="未评分提交" name="unscored">
          <el-table :data="pagedUnscored" stripe size="small">
            <el-table-column prop="student_name" label="学生" width="100" />
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="statusMap[row.status]?.type || 'info'" size="small">
                  {{ statusMap[row.status]?.label || row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="LLM评分" width="90">
              <template #default="{ row }">
                {{ row.llm_total != null ? row.llm_total.toFixed(1) : '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="submit_time" label="提交时间" width="160" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="primary" link @click="$router.push('/submission/' + row.id)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!unscoredSubmissions.length" description="暂无未评分提交" />
          <div class="pagination-wrap" v-if="unscoredSubmissions.length > subPageSize">
            <el-pagination v-model:current-page="unscoredPage" :page-size="subPageSize" :total="unscoredSubmissions.length" layout="prev, pager, next" background small />
          </div>
        </el-tab-pane>
        <el-tab-pane label="已评分提交" name="scored">
          <el-table :data="pagedScored" stripe size="small">
            <el-table-column prop="student_name" label="学生" width="100" />
            <el-table-column label="LLM评分" width="90">
              <template #default="{ row }">
                {{ row.llm_total != null ? row.llm_total.toFixed(1) : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="教师评分" width="90">
              <template #default="{ row }">
                <span style="color:#67c23a;font-weight:bold">{{ row.teacher_total != null ? row.teacher_total.toFixed(1) : '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="submit_time" label="提交时间" width="160" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="primary" link @click="$router.push('/submission/' + row.id)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!scoredSubmissions.length" description="暂无已评分提交" />
          <div class="pagination-wrap" v-if="scoredSubmissions.length > subPageSize">
            <el-pagination v-model:current-page="scoredPage" :page-size="subPageSize" :total="scoredSubmissions.length" layout="prev, pager, next" background small />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getTaskDetail } from '@/api/task'
import { getSubmissionList } from '@/api/submission'
import { getTaskIndicators, saveTaskIndicators, getIndicatorList, getCourseIndicators } from '@/api/evaluation'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const detail = ref({})
const indicators = ref([])
const loading = ref(false)

const isTeacherOrAdmin = computed(() => ['admin', 'teacher'].includes(userStore.role))

const indicatorDialogVisible = ref(false)
const allIndicators = ref([])
const savingIndicators = ref(false)

const submissionTab = ref('unscored')
const allSubmissions = ref([])
const unscoredPage = ref(1)
const scoredPage = ref(1)
const subPageSize = 8

const statusMap = {
  uploaded: { label: '已上传', type: 'info' },
  parsing: { label: '解析中', type: 'warning' },
  parsed: { label: '已解析', type: '' },
  verified: { label: '已核查', type: 'success' },
  evaluated: { label: '已评分', type: 'success' }
}

const unscoredSubmissions = computed(() =>
  allSubmissions.value.filter(s => s.status !== 'evaluated')
)
const scoredSubmissions = computed(() =>
  allSubmissions.value.filter(s => s.status === 'evaluated')
)
const pagedUnscored = computed(() => {
  const start = (unscoredPage.value - 1) * subPageSize
  return unscoredSubmissions.value.slice(start, start + subPageSize)
})
const pagedScored = computed(() => {
  const start = (scoredPage.value - 1) * subPageSize
  return scoredSubmissions.value.slice(start, start + subPageSize)
})

const weightSum = computed(() => {
  return allIndicators.value
    .filter(i => i.checked)
    .reduce((sum, i) => sum + (Number(i.weight) || 0), 0)
})

function goUpload() {
  router.push(`/submission/upload/${route.params.id}`)
}

async function fetchSubmissions() {
  try {
    const res = await getSubmissionList({ task_id: route.params.id })
    allSubmissions.value = res.data || []
  } catch {
    allSubmissions.value = []
  }
}

async function showIndicatorDialog() {
  try {
    const courseId = detail.value.course_id
    const fetches = [
      getIndicatorList(),
      getTaskIndicators(route.params.id)
    ]
    if (courseId) {
      fetches.push(getCourseIndicators(courseId))
    }
    const results = await Promise.all(fetches)
    const globalList = results[0].data || []
    const taskList = results[1].data || []
    const courseList = courseId ? (results[2]?.data || []) : []
    const taskMap = {}
    taskList.forEach(item => {
      taskMap[item.indicator_id] = item.weight
    })
    const courseMap = {}
    courseList.forEach(item => {
      courseMap[item.indicator_id] = item.weight
    })
    allIndicators.value = globalList.map(item => ({
      id: item.id,
      name: item.name,
      description: item.description,
      checked: taskMap.hasOwnProperty(item.id) || courseMap.hasOwnProperty(item.id),
      weight: taskMap[item.id] ?? courseMap[item.id] ?? item.default_weight
    }))
    indicatorDialogVisible.value = true
  } catch {
    ElMessage.warning('加载指标数据失败')
  }
}

async function saveIndicators() {
  const checkedItems = allIndicators.value.filter(i => i.checked)
  if (!checkedItems.length) {
    ElMessage.warning('请至少启用一个评价指标')
    return
  }
  savingIndicators.value = true
  try {
    const payload = {
      indicators: checkedItems.map(item => ({
        indicator_id: item.id,
        weight: Number(item.weight)
      }))
    }
    await saveTaskIndicators(route.params.id, payload)
    ElMessage.success('评价指标保存成功')
    indicatorDialogVisible.value = false
    const indRes = await getTaskIndicators(route.params.id)
    indicators.value = indRes.data || []
  } finally {
    savingIndicators.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getTaskDetail(route.params.id)
    detail.value = res.data || {}
    const [indRes] = await Promise.all([
      getTaskIndicators(route.params.id),
      fetchSubmissions()
    ])
    indicators.value = indRes.data || []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}
:deep(.el-tabs__content) {
  padding: 0;
}
</style>
