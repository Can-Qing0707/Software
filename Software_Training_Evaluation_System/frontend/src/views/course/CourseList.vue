<template>
  <div class="app-container">
    <div class="filter-bar">
      <div class="filter-left">
        <el-input v-model="query.keyword" placeholder="搜索课程" clearable style="width:200px" @clear="fetchData" @keyup.enter="fetchData" />
        <el-button type="primary" @click="fetchData">搜索</el-button>
      </div>
      <div class="filter-right">
        <el-button type="success" @click="showDialog(null)" v-if="isTeacherOrAdmin">新增课程</el-button>
        <el-button type="primary" @click="joinDialogVisible = true" v-if="isStudent">加入课程</el-button>
      </div>
    </div>
    <el-table :data="pagedList" v-loading="loading" stripe size="small">
      <el-table-column prop="name" label="课程名称" min-width="160" show-overflow-tooltip />
      <el-table-column prop="code" label="课程代码" width="110" />
      <el-table-column prop="teacher_name" label="授课教师" width="120" v-if="!isTeacher" />
      <el-table-column prop="student_count" label="参加人数" width="90" />
      <el-table-column prop="semester" label="学期" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '进行中' : '已结束' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push('/course/' + row.id)">详情</el-button>
          <template v-if="isStudent">
            <el-popconfirm title="确定退出该课程吗？" @confirm="handleLeave(row.id)">
              <template #reference>
                <el-button type="danger" link>退出</el-button>
              </template>
            </el-popconfirm>
          </template>
          <template v-if="isTeacherOrAdmin">
            <el-button type="warning" link @click="showDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑课程' : '新增课程'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-input v-model="form.semester" placeholder="如 2025-2026-2" />
        </el-form-item>
        <el-form-item label="课程描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="joinDialogVisible" title="加入课程" width="400px">
      <el-form @submit.prevent="handleJoin">
        <el-form-item label="课程代码">
          <el-input v-model="joinCode" placeholder="请输入6位课程代码" maxlength="6" style="text-transform:uppercase" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoin" :loading="joining">加入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getCourseList, createCourse, updateCourse, deleteCourse, joinCourse, getMyCourses, leaveCourse } from '@/api/course'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const isStudent = computed(() => userStore.role === 'student')
const isTeacher = computed(() => userStore.role === 'teacher')
const isTeacherOrAdmin = computed(() => ['admin', 'teacher'].includes(userStore.role))

const list = ref([])
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
const query = ref({ keyword: '' })
const form = ref({ name: '', semester: '', description: '' })
const formRef = ref(null)
const rules = {
  name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
  semester: [{ required: true, message: '请输入学期', trigger: 'blur' }]
}

const joinDialogVisible = ref(false)
const joinCode = ref('')
const joining = ref(false)

async function fetchData() {
  loading.value = true
  try {
    let res
    if (isStudent.value) {
      res = await getMyCourses()
    } else {
      res = await getCourseList({ keyword: query.value.keyword })
    }
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
    form.value = { name: row.name, semester: row.semester, description: row.description || '' }
  } else {
    isEdit.value = false
    editId.value = null
    form.value = { name: '', semester: '', description: '' }
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
    if (isEdit.value) {
      await updateCourse(editId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await createCourse(form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await deleteCourse(id)
  ElMessage.success('删除成功')
  await fetchData()
}

async function handleJoin() {
  if (!joinCode.value.trim()) {
    ElMessage.warning('请输入课程代码')
    return
  }
  joining.value = true
  try {
    await joinCourse(joinCode.value.trim().toUpperCase())
    ElMessage.success('加入课程成功')
    joinDialogVisible.value = false
    joinCode.value = ''
    await fetchData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '加入失败')
  } finally {
    joining.value = false
  }
}

async function handleLeave(id) {
  await leaveCourse(id)
  ElMessage.success('已退出课程')
  await fetchData()
}

onMounted(fetchData)
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
