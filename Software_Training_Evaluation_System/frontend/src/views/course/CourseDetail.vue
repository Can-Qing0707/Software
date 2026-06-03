<template>
  <div class="app-container">
    <el-button style="margin-bottom:16px" @click="$router.back()">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>
    <el-card v-loading="loading">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>{{ detail.name }}</span>
          <el-tag :type="detail.status === 1 ? 'success' : 'info'">
            {{ detail.status === 1 ? '进行中' : '已结束' }}
          </el-tag>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="课程代码">
          <el-tag>{{ detail.code }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="授课教师">{{ detail.teacher_name }}</el-descriptions-item>
        <el-descriptions-item label="参加人数">{{ detail.student_count ?? 0 }} 人</el-descriptions-item>
        <el-descriptions-item label="学期">{{ detail.semester }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ detail.created_at }}</el-descriptions-item>
        <el-descriptions-item label="课程描述" :span="2">{{ detail.description }}</el-descriptions-item>
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

      <el-dialog v-model="indicatorDialogVisible" title="配置课程评价指标" width="600px">
        <div style="margin-bottom:12px;color:#909399;font-size:13px">
          设置本课程默认的评价指标及权重。任务的指标权重默认继承课程配置，可在任务详情中覆盖。
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

      <h3 style="margin:24px 0 16px">实训任务
        <el-button v-if="isTeacherOrAdmin" size="small" type="success" style="margin-left:12px" @click="showTaskDialog">
          添加任务
        </el-button>
      </h3>
      <el-table :data="tasks" stripe>
        <el-table-column prop="title" label="任务名称" min-width="180" />
        <el-table-column prop="deadline" label="截止时间" width="160" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push('/task/' + row.id)">查看</el-button>
            <el-button v-if="isTeacherOrAdmin" type="danger" link @click="handleDeleteTask(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="taskDialogVisible" title="添加任务" width="600px">
      <el-form ref="taskFormRef" :model="taskForm" :rules="taskRules" label-width="100px">
        <el-form-item label="任务标题" prop="title">
          <el-input v-model="taskForm.title" />
        </el-form-item>
        <el-form-item label="任务描述" prop="description">
          <el-input v-model="taskForm.description" type="textarea" :rows="5" />
        </el-form-item>
        <el-form-item label="截止时间" prop="deadline">
          <el-date-picker v-model="taskForm.deadline" type="datetime" placeholder="选择日期" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taskDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTask" :loading="creatingTask">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getCourseDetail } from '@/api/course'
import { getTaskList, createTask, deleteTask } from '@/api/task'
import { getCourseIndicators, saveCourseIndicators, getIndicatorList } from '@/api/evaluation'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const userStore = useUserStore()
const detail = ref({})
const tasks = ref([])
const indicators = ref([])
const loading = ref(false)

const isTeacherOrAdmin = computed(() => ['admin', 'teacher'].includes(userStore.role))

const indicatorDialogVisible = ref(false)
const allIndicators = ref([])
const savingIndicators = ref(false)

const taskDialogVisible = ref(false)
const creatingTask = ref(false)
const taskFormRef = ref(null)
const taskForm = ref({ title: '', description: '', deadline: null })
const taskRules = {
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }]
}

const weightSum = computed(() => {
  return allIndicators.value
    .filter(i => i.checked)
    .reduce((sum, i) => sum + (Number(i.weight) || 0), 0)
})

async function fetchTasks() {
  const res = await getTaskList({ course_id: route.params.id })
  tasks.value = res.data || []
}

function showTaskDialog() {
  taskForm.value = { title: '', description: '', deadline: null }
  taskDialogVisible.value = true
}

async function handleCreateTask() {
  try {
    await taskFormRef.value.validate()
  } catch {
    return
  }
  creatingTask.value = true
  try {
    const payload = { ...taskForm.value, course_id: Number(route.params.id) }
    if (!payload.deadline) delete payload.deadline
    await createTask(payload)
    ElMessage.success('任务创建成功')
    taskDialogVisible.value = false
    await fetchTasks()
  } finally {
    creatingTask.value = false
  }
}

async function handleDeleteTask(id) {
  try {
    await deleteTask(id)
    ElMessage.success('删除成功')
    await fetchTasks()
  } catch {
    ElMessage.error('删除失败')
  }
}

async function showIndicatorDialog() {
  try {
    const [globalRes, courseRes] = await Promise.all([
      getIndicatorList(),
      getCourseIndicators(route.params.id)
    ])
    const globalList = globalRes.data || []
    const courseList = courseRes.data || []
    const courseMap = {}
    courseList.forEach(item => {
      courseMap[item.indicator_id] = item.weight
    })
    allIndicators.value = globalList.map(item => ({
      id: item.id,
      name: item.name,
      description: item.description,
      checked: courseMap.hasOwnProperty(item.id),
      weight: courseMap[item.id] ?? item.default_weight
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
    await saveCourseIndicators(route.params.id, payload)
    ElMessage.success('课程评价指标保存成功')
    indicatorDialogVisible.value = false
    const indRes = await getCourseIndicators(route.params.id)
    indicators.value = indRes.data || []
  } finally {
    savingIndicators.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getCourseDetail(route.params.id)
    detail.value = res.data || {}
    const [taskRes, indRes] = await Promise.all([
      getTaskList({ course_id: route.params.id }),
      getCourseIndicators(route.params.id)
    ])
    tasks.value = taskRes.data || []
    indicators.value = indRes.data || []
  } finally {
    loading.value = false
  }
})
</script>
