<template>
  <div class="app-container">
    <el-card>
      <template #header>LLM 服务配置</template>
      <el-form :model="form" label-width="120px" v-loading="loading">
        <el-form-item label="提供商">
          <el-select v-model="form.provider" style="width:300px">
            <el-option label="OpenAI 兼容" value="openai" />
            <el-option label="通义千问" value="qwen" />
            <el-option label="DeepSeek" value="deepseek" />
          </el-select>
        </el-form-item>
        <el-form-item label="API 地址">
          <el-input v-model="form.api_url" placeholder="https://api.openai.com/v1" style="width:500px" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" type="password" show-password style="width:500px" placeholder="sk-..." />
        </el-form-item>
        <el-form-item label="模型名称">
          <el-input v-model="form.model" placeholder="gpt-4o" style="width:300px" />
        </el-form-item>
        <el-form-item label="最大 Token">
          <el-input-number v-model="form.max_tokens" :min="512" :max="32768" :step="512" />
        </el-form-item>
        <el-form-item label="温度">
          <el-slider v-model="form.temperature" :min="0" :max="2" :step="0.1" style="width:300px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving">保存配置</el-button>
          <el-button @click="testConnection" :loading="testing">测试连接</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLlmConfig, updateLlmConfig } from '@/api/config'
import { ElMessage } from 'element-plus'

const form = ref({ provider: 'openai', api_url: '', api_key: '', model: 'gpt-4o', max_tokens: 4096, temperature: 0.3 })
const loading = ref(false)
const saving = ref(false)
const testing = ref(false)

async function handleSave() {
  saving.value = true
  try {
    await updateLlmConfig(form.value)
    ElMessage.success('配置已保存')
  } finally {
    saving.value = false
  }
}

async function testConnection() {
  testing.value = true
  ElMessage.info('测试连接功能需后端支持')
  setTimeout(() => { testing.value = false }, 1000)
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getLlmConfig()
    if (res.data) {
      form.value = { ...form.value, ...res.data }
    }
  } finally {
    loading.value = false
  }
})
</script>
