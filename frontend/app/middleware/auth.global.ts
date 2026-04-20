export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()

  if (to.path === '/login') {
    const isAuthenticated = await auth.ensureSession()

    if (isAuthenticated) {
      return navigateTo('/wireguard')
    }

    return
  }

  const isAuthenticated = await auth.ensureSession()

  if (!isAuthenticated) {
    return navigateTo('/login')
  }
})
