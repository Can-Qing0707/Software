import request from '@/utils/request'

export function getSubmissionList(params) {
  return request.get('/submissions', { params })
}

export function getSubmissionDetail(id) {
  return request.get(`/submissions/${id}`)
}

export function uploadFile(formData) {
  return request.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function createSubmission(data) {
  return request.post('/submissions', data)
}

export function resubmitSubmission(data) {
  return request.post('/submissions/resubmit', data)
}

export async function downloadFile(submissionId, fileIdx, fileName) {
  const blob = await request.get(`/submissions/${submissionId}/download/${fileIdx}`, {
    responseType: 'blob'
  })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

export function verifySubmission(submissionId) {
  return request.post(`/verify/${submissionId}`)
}

export function getVerification(submissionId) {
  return request.get(`/verify/${submissionId}`)
}
