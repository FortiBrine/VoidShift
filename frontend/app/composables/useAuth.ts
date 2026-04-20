type AuthStatus = 'unknown' | 'authenticated' | 'guest'

export const useAuth = () => {
  const api = useApiClient()
  const status = useState<AuthStatus>('auth-status', () => 'unknown')
  const checkInFlight = useState<Promise<boolean> | null>('auth-check-in-flight', () => null)

  const login = async (username: string, password: string) => {
    await api.request<void, { username: string; password: string }>('/api/auth/login', {
      method: 'POST',
      body: { username, password },
    })

    status.value = 'authenticated'
  }

  const logout = async () => {
    try {
      await api.request<void>('/api/auth/logout', {
        method: 'POST',
      })
    } finally {
      status.value = 'guest'
    }
  }

  const probeSession = async () => {
    try {
      await api.request('/api/vpn/wireguard/networks')
      status.value = 'authenticated'
      return true
    } catch {
      status.value = 'guest'
      return false
    }
  }

  const ensureSession = async (force = false) => {
    if (!force && status.value !== 'unknown') {
      return status.value === 'authenticated'
    }

    if (checkInFlight.value) {
      return await checkInFlight.value
    }

    checkInFlight.value = probeSession()

    try {
      return await checkInFlight.value
    } finally {
      checkInFlight.value = null
    }
  }

  return {
    status,
    login,
    logout,
    ensureSession,
  }
}
