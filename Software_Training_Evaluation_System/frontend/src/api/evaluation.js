import request from '@/utils/request'

export function getIndicatorList() {
  return request.get('/eval/indicators')
}

export function createIndicator(data) {
  return request.post('/eval/indicators', data)
}

export function updateIndicator(id, data) {
  return request.put(`/eval/indicators/${id}`, data)
}

export function deleteIndicator(id) {
  return request.delete(`/eval/indicators/${id}`)
}

export function getScores(submissionId) {
  return request.get(`/eval/score/${submissionId}`)
}

export function submitTeacherScore(data) {
  return request.post('/eval/score/teacher', data)
}

export function getTaskIndicators(taskId) {
  return request.get(`/eval/task-indicators/${taskId}`)
}

export function saveTaskIndicators(taskId, data) {
  return request.put(`/eval/task-indicators/${taskId}`, data)
}

export function triggerLlmScore(submissionId) {
  return request.post(`/eval/score/llm/${submissionId}`)
}

export function getCourseIndicators(courseId) {
  return request.get(`/eval/course-indicators/${courseId}`)
}

export function saveCourseIndicators(courseId, data) {
  return request.put(`/eval/course-indicators/${courseId}`, data)
}
