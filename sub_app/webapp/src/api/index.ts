const BASE_URL = window.location.origin

let refreshPromise: Promise<boolean> | null = null

async function tryRefreshToken(): Promise<boolean> {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) return false

  try {
    const res = await fetch(BASE_URL + '/api/member/token/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + refreshToken,
      },
    })
    if (!res.ok) return false
    const data = await res.json()
    localStorage.setItem('token', data.token)
    if (data.expires_at) {
      localStorage.setItem('expires_at', data.expires_at)
    }
    return true
  } catch {
    return false
  }
}

async function request<T = unknown>(method: string, path: string, data: unknown = null): Promise<T> {
  const token = localStorage.getItem('token')
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = 'Bearer ' + token

  const res = await fetch(BASE_URL + path, {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
  })

  if (res.status === 401 && !path.includes('/token/refresh')) {
    if (!refreshPromise) {
      refreshPromise = tryRefreshToken().finally(() => { refreshPromise = null })
    }
    const refreshed = await refreshPromise
    if (refreshed) {
      const newToken = localStorage.getItem('token')
      headers['Authorization'] = 'Bearer ' + newToken
      const retryRes = await fetch(BASE_URL + path, {
        method,
        headers,
        body: data ? JSON.stringify(data) : undefined,
      })
      const retryJson = await retryRes.json()
      if (!retryRes.ok) throw new Error(retryJson.error || 'Request failed')
      return retryJson as T
    }

    // Refresh failed — clear auth and redirect
    const savedUser = localStorage.getItem('user')
    const role = savedUser ? JSON.parse(savedUser).role : null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('expires_at')
    window.location.href = role === 'admin' ? '/admin/auth/login' : '/member/auth/login'
    throw new Error('Session expired')
  }

  const json = await res.json()
  if (!res.ok) throw new Error(json.error || 'Request failed')
  return json as T
}

// ─── Types ────────────────────────────────────────────────────────────────────

export interface LoginResponse {
  token: string
  refresh_token?: string
  expires_at: string
  refresh_expires_at?: string
  user_id: string
  username: string
  role: 'admin' | 'member'
}

export interface MemberData {
  id: string
  name?: string
  username?: string
  email?: string
  status?: string
  role: string
  created_at?: string
}

export interface TransactionData {
  id: string
  member_id: string
  account_id: string
  to_account_id?: string
  type: string
  amount: string | number
  description: string
  date: string
  category_name?: string
  account_name?: string
  to_account_name?: string
  category_id?: string
  created_at?: string
}

export interface BudgetStatus {
  id: string
  category_id: string
  category_name: string
  limit: number
  spent: number
  remaining: number
  percentage: number
}

export interface AccountData {
  id: string
  name: string
  type: string
  balance: string | number
}

export interface CategoryData {
  id: string
  name: string
  type: string
  icon?: string
}

export interface SummaryData {
  total_income: string | number
  total_expense: string | number
  net?: string | number
}

export interface CategoryBreakdown {
  category_name: string
  category_color: string
  type: string
  total: number
  count: number
}

export interface SummaryWithBreakdown extends SummaryData {
  category_breakdown: CategoryBreakdown[]
}

export interface MonthlySeriesPoint {
  year: number
  month: number
  income: number
  expense: number
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

export function loginAdmin(username: string, password: string) {
  return request<LoginResponse>('POST', '/api/admin/login', { username, password })
}

export function loginMember(username: string, password: string, rememberMe = false) {
  return request<LoginResponse>('POST', '/api/member/login', { username, password, remember_me: rememberMe })
}

export function forgotPassword(email: string) {
  return request<{ message: string }>('POST', '/api/member/auth/forgot-password', { email })
}

export function resetPassword(token: string, password: string) {
  return request<{ message: string }>('POST', '/api/member/auth/reset-password', { token, password })
}

export function registerMember(name: string, email: string, password: string) {
  return request<{ message: string; member: MemberData }>('POST', '/api/member/auth/register', { name, email, password })
}

// ─── Categories ───────────────────────────────────────────────────────────────

export function getCategories(type = '') {
  const query = type ? `?type=${type}` : ''
  return request<{ categories: CategoryData[] }>('GET', `/api/member/categories${query}`)
}

export function getAdminCategories(type = '') {
  const query = type ? `?type=${type}` : ''
  return request<{ categories: CategoryData[] }>('GET', `/api/admin/categories${query}`)
}

export function createAdminCategory(data: Record<string, unknown>) {
  return request<{ category: CategoryData }>('POST', '/api/admin/categories', data)
}

export function updateAdminCategory(id: string, data: Record<string, unknown>) {
  return request<{ category: CategoryData }>('PUT', `/api/admin/categories/${id}`, data)
}

export function deleteAdminCategory(id: string) {
  return request('DELETE', `/api/admin/categories/${id}`)
}

// ─── Accounts ─────────────────────────────────────────────────────────────────

export function getAccounts() {
  return request<{ accounts: AccountData[] }>('GET', '/api/member/accounts')
}

export function getAdminAccounts(memberId = '') {
  const query = memberId ? `?member_id=${memberId}` : ''
  return request<{ accounts: AccountData[] }>('GET', `/api/admin/accounts${query}`)
}

export function createAccount(name: string, type = 'bank', initialBalance = 0) {
  return request<AccountData>('POST', '/api/member/accounts', { name, type, initial_balance: initialBalance })
}

export function updateAccount(id: string, name: string, type: string) {
  return request<{ account: AccountData }>('PUT', `/api/member/accounts/${id}`, { name, type })
}

export function deleteAccount(id: string) {
  return request('DELETE', `/api/member/accounts/${id}`)
}

// ─── Transactions ─────────────────────────────────────────────────────────────

interface QueryParams {
  month?: string | number
  year?: string | number
  type?: string
  q?: string
  member_id?: string
  sort_by?: string
  sort_order?: string
  limit?: string | number
  offset?: string | number
}

function toQString(params: QueryParams): string {
  const q = new URLSearchParams()
  if (params.month) q.set('month', String(params.month))
  if (params.year) q.set('year', String(params.year))
  if (params.type) q.set('type', params.type)
  if (params.q) q.set('q', params.q)
  if (params.member_id) q.set('member_id', params.member_id)
  if (params.sort_by) q.set('sort_by', params.sort_by)
  if (params.sort_order) q.set('sort_order', params.sort_order)
  if (params.limit) q.set('limit', String(params.limit))
  if (params.offset) q.set('offset', String(params.offset))
  return q.toString() ? `?${q.toString()}` : ''
}

export function getTransactions(params: QueryParams = {}) {
  return request<{ transactions: TransactionData[]; total: number }>('GET', `/api/member/transactions${toQString(params)}`)
}

export function getAdminTransactions(params: QueryParams = {}) {
  return request<{ transactions: TransactionData[]; total: number }>('GET', `/api/admin/transactions${toQString(params)}`)
}

export function createTransaction(data: Record<string, unknown>) {
  return request<TransactionData>('POST', '/api/member/transactions', data)
}

export function updateTransaction(id: string, data: Record<string, unknown>) {
  return request<TransactionData>('PUT', `/api/member/transactions/${id}`, data)
}

export function updateAdminTransaction(id: string, data: Record<string, unknown>) {
  return request<TransactionData>('PUT', `/api/admin/transactions/${id}`, data)
}

export function deleteTransaction(id: string) {
  return request('DELETE', `/api/admin/transactions/${id}`)
}

export function deleteMemberTransaction(id: string) {
  return request('DELETE', `/api/member/transactions/${id}`)
}

export function getTransactionSummary(params: QueryParams = {}) {
  return request<SummaryData>('GET', `/api/member/transactions/summary${toQString(params)}`)
}

export function getAdminTransactionSummary(params: QueryParams = {}) {
  return request<SummaryData>('GET', `/api/admin/transactions/summary${toQString(params)}`)
}

export function getTransactionSummaryWithBreakdown(params: QueryParams = {}) {
  return request<SummaryWithBreakdown>('GET', `/api/member/transactions/summary${toQString(params)}`)
}

export function getAdminTransactionSummaryWithBreakdown(params: QueryParams = {}) {
  return request<SummaryWithBreakdown>('GET', `/api/admin/transactions/summary${toQString(params)}`)
}

export function getMonthlySeries(params: QueryParams & { months?: number } = {}) {
  const q = new URLSearchParams()
  if (params.month) q.set('month', String(params.month))
  if (params.year) q.set('year', String(params.year))
  if (params.months) q.set('months', String(params.months))
  const query = q.toString() ? `?${q.toString()}` : ''
  return request<{ series: MonthlySeriesPoint[] }>('GET', `/api/member/transactions/monthly${query}`)
}

export function getAdminMonthlySeries(params: QueryParams & { months?: number; member_id?: string } = {}) {
  const q = new URLSearchParams()
  if (params.months) q.set('months', String(params.months))
  if (params.member_id) q.set('member_id', params.member_id)
  const query = q.toString() ? `?${q.toString()}` : ''
  return request<{ series: MonthlySeriesPoint[] }>('GET', `/api/admin/transactions/monthly${query}`)
}

// ─── Budgets ──────────────────────────────────────────────────────────────────

export function getBudgetStatus(params: QueryParams = {}) {
  return request<{ budgets: BudgetStatus[] }>('GET', `/api/member/budgets/status${toQString(params)}`)
}

export function getAdminBudgetStatus(params: QueryParams = {}) {
  return request<{ budgets: BudgetStatus[] }>('GET', `/api/admin/budgets/status${toQString(params)}`)
}

export function setBudget(data: Record<string, unknown>) {
  return request<BudgetStatus>('POST', '/api/member/budgets', data)
}

export function setAdminBudget(data: Record<string, unknown>) {
  return request<BudgetStatus>('POST', '/api/admin/budgets', data)
}

export function deleteAdminBudget(id: string) {
  return request('DELETE', `/api/admin/budgets/${id}`)
}

export function deleteMemberBudget(id: string) {
  return request('DELETE', `/api/member/budgets/${id}`)
}

// ─── Members (Admin) ──────────────────────────────────────────────────────────

export function getMembers() {
  return request<{ members: MemberData[] }>('GET', '/api/admin/members')
}

export function createMember(data: Record<string, unknown>) {
  return request<MemberData>('POST', '/api/admin/members', data)
}

export function updateMember(id: string, data: Record<string, unknown>) {
  return request<MemberData>('PUT', `/api/admin/members/${id}`, data)
}

export function deleteMember(id: string) {
  return request('DELETE', `/api/admin/members/${id}`)
}

export function approveMember(id: string) {
  return request<{ message: string }>('PATCH', `/api/admin/members/${id}/approve`)
}

export function rejectMember(id: string) {
  return request<{ message: string }>('PATCH', `/api/admin/members/${id}/reject`)
}

// ─── Export ────────────────────────────────────────────────────────────────────

export async function exportTransactionsCSV(params: { month?: number; year?: number; member_id?: string; member_name?: string }) {
  const q = new URLSearchParams()
  if (params.month) q.set('month', String(params.month))
  if (params.year) q.set('year', String(params.year))
  if (params.member_id) q.set('member_id', params.member_id)
  if (params.member_name) q.set('member_name', params.member_name)
  const query = q.toString() ? `?${q.toString()}` : ''

  const token = localStorage.getItem('token')
  const res = await fetch(`${BASE_URL}/api/admin/transactions/export${query}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Export failed' }))
    throw new Error(err.error || 'Export failed')
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `transactions_${params.year || new Date().getFullYear()}_${String(params.month || new Date().getMonth() + 1).padStart(2, '0')}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

// ─── Health ───────────────────────────────────────────────────────────────────

export function healthCheck() {
  return request<{ status: string }>('GET', '/api/admin/health')
}

// ─── Profile & Password ───────────────────────────────────────────────────────

export interface MemberProfile {
  id: string
  name: string
  username: string
  created_at: string
}

export interface AdminProfile {
  id: string
  username: string
  email?: string
  created_at: string
}

export function getMemberProfile() {
  return request<MemberProfile>('GET', '/api/member/profile')
}

export function updateMemberProfile(data: { name?: string; username?: string }) {
  return request<MemberProfile>('PUT', '/api/member/profile', data)
}

export function changeMemberPassword(oldPassword: string, newPassword: string) {
  return request<{ message: string }>('PUT', '/api/member/password', { old_password: oldPassword, new_password: newPassword })
}

export function getAdminProfile() {
  return request<AdminProfile>('GET', '/api/admin/profile')
}

export function updateAdminProfile(data: { username?: string; email?: string }) {
  return request<AdminProfile>('PUT', '/api/admin/profile', data)
}

export function changeAdminPassword(oldPassword: string, newPassword: string) {
  return request<{ message: string }>('PUT', '/api/admin/password', { old_password: oldPassword, new_password: newPassword })
}

// ─── SMTP Settings (Admin) ──────────────────────────────────────────────────

export interface SMTPConfig {
  host: string
  port: string
  user: string
  pass: string
  from_email: string
  from_name: string
  use_tls: boolean
  use_starttls: boolean
  skip_verify: boolean
}

export interface SettingsResponse {
  smtp: SMTPConfig
  env_smtp: SMTPConfig
  override: boolean
  notif_email: string
}

export function getAdminSMTPConfig() {
  return request<SettingsResponse>('GET', '/api/admin/settings/smtp')
}

export function saveAdminSMTPConfig(cfg: Partial<SMTPConfig> & { notif_email?: string }) {
  return request<SettingsResponse & { message: string }>('PUT', '/api/admin/settings/smtp', cfg)
}
