import request from '@/utils/request'

export function getReports(params) {
  return request.get('/reports', { params })
}

export function generateReport(data) {
  return request.post('/reports/generate', data)
}

export function exportReport(id, format) {
  return request.get(`/reports/export/${id}`, {
    params: { format },
    responseType: 'blob'
  })
}
