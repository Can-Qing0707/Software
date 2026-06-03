import request from '@/utils/request'

export function getCourseList(params) {
  return request.get('/courses', { params })
}

export function getMyCourses() {
  return request.get('/courses/my')
}

export function getCourseDetail(id) {
  return request.get(`/courses/${id}`)
}

export function createCourse(data) {
  return request.post('/courses', data)
}

export function updateCourse(id, data) {
  return request.put(`/courses/${id}`, data)
}

export function deleteCourse(id) {
  return request.delete(`/courses/${id}`)
}

export function joinCourse(code) {
  return request.post('/courses/join', { code })
}

export function leaveCourse(courseId) {
  return request.delete(`/courses/${courseId}/leave`)
}
