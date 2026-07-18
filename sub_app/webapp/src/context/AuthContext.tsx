import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { loginAdmin, loginMember } from '../api'

interface User {
  id: string
  username: string
  role: 'admin' | 'member'
}

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  isAuthenticated: boolean
  isAdmin: boolean
  isMember: boolean
  loginAdmin: (username: string, password: string) => Promise<unknown>
  loginMember: (username: string, password: string, rememberMe?: boolean) => Promise<unknown>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

interface TokenPayload {
  exp: number
  sub?: string
  user_id?: string
}

// Parse JWT to get expiry
function parseToken(token: string): TokenPayload | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1])) as TokenPayload
    return { exp: payload.exp * 1000, user_id: payload.sub || payload.user_id }
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    if (savedToken && savedUser) {
      const parsed = parseToken(savedToken)
      if (parsed && parsed.exp > Date.now()) {
        setToken(savedToken)
        setUser(JSON.parse(savedUser) as User)
      } else {
        // Token expired — try refresh
        // (refresh will be handled by API interceptor on next request)
        setToken(savedToken)
        setUser(JSON.parse(savedUser) as User)
      }
    }
    setLoading(false)
  }, [])

  function saveAuth(token: string, user: User, refreshToken?: string) {
    setToken(token)
    setUser(user)
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    if (refreshToken) {
      localStorage.setItem('refresh_token', refreshToken)
    } else {
      localStorage.removeItem('refresh_token')
    }
  }

  function clearAuth() {
    setToken(null)
    setUser(null)
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('expires_at')
  }

  const handleLoginAdmin = useCallback(async (username: string, password: string) => {
    const res = await loginAdmin(username, password)
    saveAuth(res.token, {
      id: res.user_id,
      username: res.username,
      role: res.role,
    })
    return res
  }, [])

  const handleLoginMember = useCallback(async (username: string, password: string, rememberMe?: boolean) => {
    const res = await loginMember(username, password, rememberMe)
    saveAuth(res.token, {
      id: res.user_id,
      username: res.username,
      role: res.role,
    }, res.refresh_token)
    return res
  }, [])

  const handleLogout = useCallback(() => {
    clearAuth()
  }, [])

  const isAuthenticated = !!token && !!user
  const isAdmin = isAuthenticated && user?.role === 'admin'
  const isMember = isAuthenticated && user?.role === 'member'

  const value: AuthContextType = {
    user, token, loading,
    isAuthenticated, isAdmin, isMember,
    loginAdmin: handleLoginAdmin,
    loginMember: handleLoginMember,
    logout: handleLogout,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
