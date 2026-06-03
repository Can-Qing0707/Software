const roleHierarchy = {
  admin: 3,
  teacher: 2,
  student: 1
}

export function hasPermission(userRole, requiredRole) {
  return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

export default {
  install(app) {
    app.directive('role', (el, binding) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      const userRole = user.role || ''
      const requiredRole = binding.value
      if (requiredRole && !hasPermission(userRole, requiredRole)) {
        el.parentNode?.removeChild(el)
      }
    })
  }
}
