<template>
  <div class="app-container">
    <el-button style="margin-bottom:16px" @click="$router.back()">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card v-loading="loading">
          <template #header>LLM 自动评分</template>
          <div v-if="!llmScores.length" style="color:#909399;text-align:center;padding:40px 0">
            <el-icon style="font-size:48px;margin-bottom:12px"><InfoFilled /></el-icon>
            <p>LLM评分尚未生成</p>
            <el-button type="primary" style="margin-top:12px" @click="runLlmEval" :loading="llmLoading">执行LLM评分</el-button>
          </div>
          <div v-else>
            <div v-for="s in llmScores" :key="s.id" class="score-row">
              <div style="display:flex;justify-content:space-between;align-items:center">
                <span style="font-weight:bold">{{ s.indicator_name }}</span>
                <span style="font-size:18px;color:#409eff;font-weight:bold">{{ s.llm_score }}</span>
              </div>
              <div v-if="s.llm_comment" style="font-size:12px;color:#909399;margin-top:4px">{{ s.llm_comment }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>教师评分</template>
          <el-form :model="teacherForm" label-width="100px">
            <el-form-item v-for="item in allIndicators" :key="item.indicator_id" :label="item.name">
              <div style="display:flex;gap:8px;align-items:center">
                <el-input-number v-model="teacherForm[item.indicator_id]" :min="0" :max="100" :precision="2" size="small" style="width:160px" />
                <span style="color:#909399;font-size:12px">权重 {{ item.weight }}%</span>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitTeacherScore" :loading="submitting">保存评分</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getScores, getTaskIndicators, submitTeacherScore as submitTeacherScoreApi, triggerLlmScore } from '@/api/evaluation'
import { getSubmissionDetail } from '@/api/submission'
import { ElMessage } from 'element-plus'

const route = useRoute()
const submissionId = route.params.submissionId
const llmScores = ref([])
const allIndicators = ref([])
const loading = ref(false)
const llmLoading = ref(false)
const submitting = ref(false)
const teacherForm = reactive({})

async function runLlmEval() {
  llmLoading.value = true
  try {
    await triggerLlmScore(submissionId)
    ElMessage.success('LLM评分完成')
    await fetchData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.message || 'LLM评分失败')
  } finally {
    llmLoading.value = false
  }
}

async function submitTeacherScore() {
  submitting.value = true
  try {
    const scores = Object.entries(teacherForm)
      .filter(([k, v]) => v !== null && v !== undefined)
      .map(([indicatorId, score]) => ({
        indicator_id: Number(indicatorId),
        teacher_score: Number(score)
      }))
      await submitTeacherScoreApi({ submission_id: Number(submissionId), scores })
    ElMessage.success('评分已保存')
  } finally {
    submitting.value = false
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getScores(submissionId)
    llmScores.value = res.data || []
    const subRes = await getSubmissionDetail(submissionId)
    const taskId = subRes.data?.task_id
    if (taskId) {
      const indRes = await getTaskIndicators(taskId)
      allIndicators.value = indRes.data || []
      const existingMap = {}
      ;(res.data || []).forEach(s => {
        if (s.teacher_score != null) existingMap[s.indicator_id] = s.teacher_score
      })
      allIndicators.value.forEach(item => {
        if (teacherForm[item.indicator_id] === undefined) {
          teacherForm[item.indicator_id] = existingMap[item.indicator_id] ?? null
        }
      })
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped>
.score-row {
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
.score-row:last-child {
  border-bottom: none;
}
</style>
