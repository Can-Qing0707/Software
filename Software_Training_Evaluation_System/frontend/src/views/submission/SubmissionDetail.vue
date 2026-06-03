<template>
  <div class="app-container">
    <el-button style="margin-bottom:16px" @click="$router.back()">
      <el-icon><ArrowLeft /></el-icon> 返回
    </el-button>

    <el-row :gutter="16">
      <el-col :span="isTeacher ? 12 : 16">
        <el-card v-loading="loading">
          <template #header>
            <span>{{ detail.student_name }} - {{ detail.task_title }}</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="提交时间">{{ detail.submit_time }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusType" size="small">{{ statusLabel }}</el-tag>
            </el-descriptions-item>
          </el-descriptions>

          <h4 style="margin:16px 0 8px">提交文件</h4>
          <div v-for="(f, idx) in detail.files" :key="idx" class="file-row">
            <el-icon><Document /></el-icon>
            <span style="flex:1">{{ f.name }}</span>
            <el-button type="primary" link size="small" @click="handleDownload(idx, f.name)">下载</el-button>
          </div>
          <div v-if="!detail.files?.length" style="color:#909399;font-size:13px">无文件</div>

          <h4 style="margin:16px 0 8px">解析内容</h4>
          <div v-if="detail.content_text" style="background:#f5f7fa;padding:12px;border-radius:4px;white-space:pre-wrap;font-size:13px;max-height:300px;overflow-y:auto">
            {{ detail.content_text }}
          </div>
          <div v-else style="color:#909399;font-size:13px">暂未解析</div>
        </el-card>

        <el-card v-if="verification" style="margin-top:16px">
          <template #header>
            <span>核查结果</span>
            <el-button v-if="isTeacher" size="small" type="primary" style="float:right" :loading="verifying" @click="handleVerify">
              重新核查
            </el-button>
          </template>
          <div style="text-align:center;margin-bottom:12px">
            <el-tag :type="verification.overall_pass ? 'success' : 'danger'" size="large">
              {{ verification.overall_pass ? '通过' : '不通过' }}
            </el-tag>
          </div>
          <div v-if="verification.completeness" style="margin-bottom:8px">
            <div style="font-weight:bold;margin-bottom:4px">步骤完整性</div>
            <el-progress :percentage="Math.round((verification.completeness.completeness_ratio || 0) * 100)" />
          </div>
          <div v-if="verification.requirement_match">
            <div style="font-weight:bold;margin-bottom:4px">要求匹配度</div>
            <el-progress :percentage="Math.round((verification.requirement_match.match_ratio || 0) * 100)" :status="verification.requirement_match.match_ratio >= 0.7 ? 'success' : 'warning'" />
          </div>
        </el-card>

        <el-card v-else-if="isTeacher" style="margin-top:16px">
          <template #header>核查操作</template>
          <el-button type="primary" :loading="verifying" @click="handleVerify">
            执行智能核查
          </el-button>
        </el-card>
      </el-col>

      <el-col :span="isTeacher ? 12 : 8">
        <template v-if="isTeacher">
          <el-card style="margin-bottom:16px">
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

          <el-card>
            <template #header>教师评分</template>
            <el-form :model="teacherForm" label-width="100px">
              <el-form-item v-for="item in scoreIndicators" :key="item.indicator_id" :label="item.name">
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
        </template>

        <template v-else>
          <el-card style="margin-top:16px" v-if="scores.length">
            <template #header>评分结果</template>
            <div v-for="s in scores" :key="s.id" class="score-item">
              <div style="display:flex;justify-content:space-between">
                <span>{{ s.indicator_name }}</span>
                <span style="font-weight:bold;color:#409eff">{{ s.final_score ?? s.llm_score ?? '-' }}</span>
              </div>
            </div>
            <el-divider />
            <div style="display:flex;justify-content:space-between;font-weight:bold;font-size:16px">
              <span>总分</span>
              <span style="color:#409eff">{{ totalScore }}</span>
            </div>
          </el-card>

          <el-card style="margin-top:16px" v-if="!scores.length">
            <template #header>评分状态</template>
            <div style="text-align:center;padding:16px 0;color:#909399">
              <el-icon style="font-size:36px;margin-bottom:8px"><Clock /></el-icon>
              <p>暂未评分</p>
              <p style="font-size:12px">教师尚未对此提交进行评分</p>
            </div>
          </el-card>
        </template>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getSubmissionDetail, verifySubmission, getVerification, downloadFile } from '@/api/submission'
import { getScores, getTaskIndicators, submitTeacherScore as submitTeacherScoreApi, triggerLlmScore } from '@/api/evaluation'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()
const detail = ref({})
const verification = ref(null)
const scores = ref([])
const loading = ref(false)
const verifying = ref(false)

const isTeacher = computed(() => ['admin', 'teacher'].includes(userStore.role))

const statusMap = {
  uploaded: { label: '已上传', type: 'info' },
  parsing: { label: '解析中', type: 'warning' },
  parsed: { label: '已解析', type: '' },
  verified: { label: '已核查', type: 'success' },
  evaluated: { label: '已评分', type: 'success' }
}
const statusLabel = computed(() => statusMap[detail.value.status]?.label || detail.value.status)
const statusType = computed(() => statusMap[detail.value.status]?.type || 'info')

const totalScore = computed(() => {
  const total = scores.value.reduce((sum, s) => sum + (s.final_score ?? s.llm_score ?? 0), 0)
  return scores.value.length ? (total / scores.value.length).toFixed(1) : '-'
})

const llmScores = ref([])
const llmLoading = ref(false)
const scoreIndicators = ref([])
const submitting = ref(false)
const teacherForm = reactive({})

async function runLlmEval() {
  llmLoading.value = true
  try {
    await triggerLlmScore(route.params.id)
    ElMessage.success('LLM评分完成')
    await refreshScoreData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.message || 'LLM评分失败')
  } finally {
    llmLoading.value = false
  }
}

async function submitTeacherScore() {
  submitting.value = true
  try {
    const items = Object.entries(teacherForm)
      .filter(([, v]) => v !== null && v !== undefined)
      .map(([indicatorId, score]) => ({
        indicator_id: Number(indicatorId),
        teacher_score: Number(score)
      }))
    if (!items.length) {
      ElMessage.warning('请至少输入一项评分')
      return
    }
    await submitTeacherScoreApi({ submission_id: Number(route.params.id), scores: items })
    ElMessage.success('评分已保存')
  } finally {
    submitting.value = false
  }
}

async function refreshScoreData() {
  const res = await getScores(route.params.id)
  scores.value = res.data || []
  llmScores.value = scores.value.filter(s => s.llm_score != null)

  const existingMap = {}
  scores.value.forEach(s => {
    if (s.teacher_score != null) existingMap[s.indicator_id] = s.teacher_score
  })

  const subRes = await getSubmissionDetail(route.params.id)
  const taskId = subRes.data?.task_id
  if (taskId) {
    const indRes = await getTaskIndicators(taskId)
    scoreIndicators.value = indRes.data || []
    scoreIndicators.value.forEach(item => {
      if (teacherForm[item.indicator_id] === undefined) {
        teacherForm[item.indicator_id] = existingMap[item.indicator_id] ?? null
      }
    })
  }
}

async function handleVerify() {
  verifying.value = true
  try {
    const res = await verifySubmission(route.params.id)
    verification.value = res.data || null
    detail.value.status = 'verified'
    ElMessage.success('核查完成')
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '核查失败')
  } finally {
    verifying.value = false
  }
}

async function handleDownload(idx, fileName) {
  try {
    await downloadFile(route.params.id, idx, fileName)
  } catch (e) {
    ElMessage.error('文件下载失败')
  }
}

async function loadVerification() {
  try {
    const res = await getVerification(route.params.id)
    verification.value = res.data || null
  } catch {
    // 核查记录不存在时忽略
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getSubmissionDetail(route.params.id)
    detail.value = res.data || {}
    verification.value = res.data?.verification || null
    if (!verification.value) {
      await loadVerification()
    }
    await refreshScoreData()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.file-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 0;
  font-size: 13px;
}
.score-item {
  padding: 6px 0;
  font-size: 14px;
}
.score-row {
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
.score-row:last-child {
  border-bottom: none;
}
</style>
