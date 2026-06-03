import request from '@/utils/request'

export function getLlmConfig() {
  return request.get('/config/llm')
}

export function updateLlmConfig(data) {
  return request.put('/config/llm', data)
}
