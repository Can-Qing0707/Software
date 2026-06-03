<template>
  <div class="app-container">
    <el-button style="margin-bottom:16px" @click="$router.back()">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>
    <el-card>
      <template #header>{{ isResubmit ? '重新提交成果' : '提交实训成果' }}</template>
      <el-form label-width="100px" v-loading="loading">
        <el-form-item label="实训任务" required>
          <el-select v-model="selectedTaskId" placeholder="请选择要提交的实训任务" style="width:100%" @change="onTaskChange">
            <el-option v-for="t in availableTasks" :key="t.id" :label="`${t.title} (${t.course_name})`" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="selectedTask && selectedTask.deadline" label="截止时间">
          <el-tag :type="new Date(selectedTask.deadline) < new Date() ? 'danger' : ''">
            {{ selectedTask.deadline }}
          </el-tag>
        </el-form-item>
        <el-form-item label="上传文件">
          <el-upload
            ref="uploadRef"
            drag
            multiple
            :auto-upload="false"
            :on-change="handleFileChange"
            accept=".doc,.docx,.pdf,.png,.jpg,.jpeg,.zip"
          >
            <el-icon class="el-icon--upload" size="48"><UploadFilled /></el-icon>
            <div class="el-upload__text">拖拽文件到此处或 <em>点击上传</em></div>
            <template #tip>
              <div class="el-upload__tip">支持格式：Word(.doc/.docx)、PDF、图片、ZIP压缩包</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item v-if="files.length">
          <div style="width:100%">
            <div v-for="(f, idx) in files" :key="idx" class="file-item">
              <el-icon><Document /></el-icon>
              <span>{{ f.name }}</span>
              <span class="file-size">({{ (f.size / 1024).toFixed(1) }} KB)</span>
              <el-button type="danger" link @click="removeFile(idx)">移除</el-button>
            </div>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting" :disabled="!selectedTaskId || !files.length">
            {{ isResubmit ? '重新提交' : '提交' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getTaskList } from '@/api/task'
import { uploadFile, createSubmission, getSubmissionList, resubmitSubmission } from '@/api/submission'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const availableTasks = ref([])
const selectedTaskId = ref(null)
const submittedTaskIds = ref(new Set())
const files = ref([])
const submitting = ref(false)
const loading = ref(false)

const preselectedTaskId = computed(() => route.params.taskId ? Number(route.params.taskId) : null)

const selectedTask = computed(() => availableTasks.value.find(t => t.id === selectedTaskId.value) || null)

const isResubmit = computed(() => selectedTaskId.value && submittedTaskIds.value.has(selectedTaskId.value))

function onTaskChange() {
  files.value = []
}

function handleFileChange(uploadFile) {
  files.value.push(uploadFile.raw)
}

function removeFile(idx) {
  files.value.splice(idx, 1)
}

async function handleSubmit() {
  if (!selectedTaskId.value) {
    ElMessage.warning('请选择实训任务')
    return
  }
  submitting.value = true
  try {
    const uploadedFiles = []
    for (const file of files.value) {
      const formData = new FormData()
      formData.append('file', file)
      const res = await uploadFile(formData)
      uploadedFiles.push(res.data)
    }
    const payload = { task_id: selectedTaskId.value, files: uploadedFiles }
    if (isResubmit.value) {
      await resubmitSubmission(payload)
      ElMessage.success('重新提交成功')
    } else {
      await createSubmission(payload)
      ElMessage.success('提交成功')
    }
    router.push('/submission')
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const [taskRes, subRes] = await Promise.all([
      getTaskList(),
      getSubmissionList()
    ])
    const allTasks = taskRes.data || []
    const submissions = subRes.data || []
    const ids = new Set(submissions.map(s => s.task_id))
    submittedTaskIds.value = ids
    availableTasks.value = allTasks.filter(t => !ids.has(t.id))

    if (preselectedTaskId.value) {
      const preTask = allTasks.find(t => t.id === preselectedTaskId.value)
      if (preTask && ids.has(preselectedTaskId.value)) {
        availableTasks.value.unshift(preTask)
      }
      if (preTask) {
        selectedTaskId.value = preselectedTaskId.value
      }
    }
  } catch {
    ElMessage.warning('加载任务列表失败')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  margin-bottom: 6px;
}
.file-size {
  color: #909399;
  font-size: 12px;
}
</style>
