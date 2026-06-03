<template>
  <div class="app-container">
    <div class="filter-container">
      <el-select v-model="query.role" placeholder="角色" clearable @change="fetchData" style="width:140px">
        <el-option label="管理员" value="admin" />
        <el-option label="教师" value="teacher" />
        <el-option label="学生" value="student" />
      </el-select>
      <el-input v-model="query.keyword" placeholder="搜索用户名/姓名" clearable style="width:200px" @keyup.enter="fetchData" />
      <el-button type="primary" @click="fetchData">搜索</el-button>
      <el-button type="success" @click="showDialog(null)">新增用户</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="real_name" label="姓名" width="120" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="roleMap[row.role]?.type || ''" size="small">{{ roleMap[row.role]?.label || row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱" min-width="160" />
      <el-table-column prop="phone" label="手机" width="130" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button type="warning" link @click="showDialog(row)">编辑</el-button>
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)">
            <template #reference>
              <el-button type="danger" link>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新增用户'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="form.real_name" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password :placeholder="isEdit ? '留空则不修改' : '请输入密码'" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" style="width:100%">
            <el-option label="管理员" value="admin" />
            <el-option label="教师" value="teacher" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input v-model="form.phone" />
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
import { ref, onMounted } from 'vue'
import { getUserList, createUser, updateUser, deleteUser } from '@/api/user'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const editId = ref(null)
const query = ref({ role: '', keyword: '' })
const form = ref({ username: '', real_name: '', password: '', role: 'student', email: '', phone: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}
const roleMap = { admin: { label: '管理员', type: 'danger' }, teacher: { label: '教师', type: 'warning' }, student: { label: '学生', type: '' } }

async function fetchData() {
  loading.value = true
  try {
    const res = await getUserList(query.value)
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

function showDialog(row) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.value = { username: row.username, real_name: row.real_name, password: '', role: row.role, email: row.email || '', phone: row.phone || '' }
  } else {
    isEdit.value = false
    editId.value = null
    form.value = { username: '', real_name: '', password: '', role: 'student', email: '', phone: '' }
  }
  dialogVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      if (!payload.password) delete payload.password
      await updateUser(editId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createUser(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await deleteUser(id)
  ElMessage.success('删除成功')
  await fetchData()
}

onMounted(fetchData)
</script>
