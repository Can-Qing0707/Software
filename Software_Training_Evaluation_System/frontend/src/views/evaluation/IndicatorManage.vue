<template>
  <div class="app-container">
    <div style="display:flex;justify-content:space-between;margin-bottom:16px">
      <h3>评价指标管理</h3>
      <el-button type="primary" @click="showDialog(null)">新增指标</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="name" label="指标名称" width="160" />
      <el-table-column prop="description" label="说明" min-width="240" />
      <el-table-column prop="default_weight" label="默认权重(%)" width="120">
        <template #default="{ row }">
          <el-progress :percentage="row.default_weight" :stroke-width="16" />
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="70" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑指标' : '新增指标'" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="指标名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="默认权重(%)" prop="default_weight">
          <el-input-number v-model="form.default_weight" :min="0" :max="100" :precision="2" style="width:200px" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" style="width:200px" />
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
import { getIndicatorList, createIndicator, updateIndicator, deleteIndicator } from '@/api/evaluation'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const editId = ref(null)
const form = ref({ name: '', description: '', default_weight: 0, sort_order: 0 })
const rules = {
  name: [{ required: true, message: '请输入指标名称', trigger: 'blur' }]
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getIndicatorList()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

function showDialog(row) {
  if (row) {
    isEdit.value = true
    editId.value = row.id
    form.value = { name: row.name, description: row.description || '', default_weight: Number(row.default_weight), sort_order: row.sort_order }
  } else {
    isEdit.value = false
    editId.value = null
    form.value = { name: '', description: '', default_weight: 0, sort_order: 0 }
  }
  dialogVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (isEdit.value) {
      await updateIndicator(editId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await createIndicator(form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await deleteIndicator(id)
  ElMessage.success('删除成功')
  await fetchData()
}

onMounted(fetchData)
</script>
