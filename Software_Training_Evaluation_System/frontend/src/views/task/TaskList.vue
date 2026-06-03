<template>
  <div class="app-container">
    <div class="filter-bar">
      <div class="filter-left">
        <el-select v-model="query.course_id" placeholder="选择课程" clearable @change="fetchData" style="width:160px">
          <el-option v-for="c in courses" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
        <el-input v-model="query.keyword" placeholder="搜索任务" clearable style="width:180px" @keyup.enter="fetchData" />
        <el-button type="primary" @click="fetchData">搜索</el-button>
      </div>
      <div class="filter-right">
        <el-button type="success" @click="showDialog(null)" v-role="'teacher'">新增任务</el-button>
      </div>
    </div>
    <el-table :data="pagedList" v-loading="loading" stripe size="small">
      <el-table-column prop="title" label="任务名称" min-width="160" show-overflow-tooltip />
      <el-table-column prop="course_name" label="所属课程" min-width="120" show-overflow-tooltip />
      <el-table-column label="截止时间" width="100">
        <template #default="{ row }">
          <span v-if="row.deadline">{{ row.deadline.slice(0, 10) }}</span>
          <span v-else style="color:#c0c4cc">未设置</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push('/task/' + row.id)">详情</el-button>
          <el-button type="warning" link @click="showDialog(row)" v-role="'teacher'">编辑</el-button>
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button type="danger" link v-role="'teacher'">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap" v-if="list.length > pageSize">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="list.length"
        layout="prev, pager, next"
        background
        small
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑任务' : '新增任务'" width="700px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="所属课程" prop="course_id">
          <el-select v-model="form.course_id" placeholder="选择课程" style="width:100%">
            <el-option v-for="c in courses" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="任务标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="任务描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="6" />
        </el-form-item>
        <el-form-item label="截止时间" prop="deadline">
          <el-date-picker v-model="form.deadline" type="datetime" placeholder="选择日期" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTaskList, createTask, updateTask, deleteTask } from '@/api/task'
import { getCourseList } from '@/api/course'
import { ElMessage } from 'element-plus'

const list = ref([])
const courses = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 8

const pagedList = computed(() => {
  const start = (page.value - 1) * pageSize
  return list.value.slice(start, start + pageSize)
})
const dialogVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const editId = ref(null)
const query = ref({ course_id: '', keyword: '' })
const form = ref({ course_id: '', title: '', description: '', deadline: null })
const formRef = ref(null)
const rules = {
  course_id: [{ required: true, message: '请选择课程', trigger: 'change' }],
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }]
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getTaskList(query.value)
    list.value = res.data || []
    page.value = 1
  } finally {
    loading.value = false
  }
}

function showDialog(row) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.value = { course_id: row.course_id, title: row.title, description: row.description || '', deadline: row.deadline || null }
  } else {
    isEdit.value = false
    editId.value = null
    form.value = { course_id: query.value.course_id || '', title: '', description: '', deadline: null }
  }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (!payload.deadline) delete payload.deadline
    if (isEdit.value) {
      await updateTask(editId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createTask(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await deleteTask(id)
  ElMessage.success('删除成功')
  await fetchData()
}

onMounted(async () => {
  const res = await getCourseList()
  courses.value = res.data || []
  await fetchData()
})
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.filter-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}
:deep(.el-table__row) {
  height: 44px;
}
:deep(.el-table__cell) .cell {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
