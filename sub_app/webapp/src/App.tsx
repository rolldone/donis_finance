import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { ThemeProvider } from './context/ThemeContext'
import { MemberRoute, AdminRoute, GuestRoute } from './components/ProtectedRoute'
import DashboardLayout from './components/DashboardLayout'
import AdminLayout from './components/AdminLayout'

// Auth pages
import Login from './pages/Login'
import Register from './pages/Register'
import ForgotPassword from './pages/ForgotPassword'
import ResetPassword from './pages/ResetPassword'

// Admin auth
import AdminLogin from './pages/admin/Login'

// Member pages
import MemberDashboard from './pages/member/Dashboard'
import MemberTransactions from './pages/member/Transactions'
import MemberBudget from './pages/member/Budget'
import MemberProfile from './pages/member/Profile'

// Admin pages
import AdminDashboard from './pages/admin/Dashboard'
import AdminMembers from './pages/admin/Members'
import AdminCategories from './pages/admin/Categories'
import AdminTransactions from './pages/admin/Transactions'
import AdminBudget from './pages/admin/Budget'
import AdminSettings from './pages/admin/Settings'

function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
      <AuthProvider>
        <Routes>
          {/* Default redirect */}
          <Route path="/" element={<Navigate to="/member/auth/login" replace />} />

          {/* Auth routes (guest only — redirect if already logged in) */}
          <Route path="/member/auth/login" element={<GuestRoute><Login /></GuestRoute>} />
          <Route path="/member/auth/register" element={<GuestRoute><Register /></GuestRoute>} />
          <Route path="/member/auth/forgot-password" element={<GuestRoute><ForgotPassword /></GuestRoute>} />
          <Route path="/member/auth/reset-password" element={<GuestRoute><ResetPassword /></GuestRoute>} />

          {/* Admin auth */}
          <Route path="/admin/auth/login" element={<GuestRoute><AdminLogin /></GuestRoute>} />

          {/* Member dashboard routes (protected) */}
          <Route path="/member" element={<MemberRoute><DashboardLayout><MemberDashboard /></DashboardLayout></MemberRoute>} />
          <Route path="/member/transactions" element={<MemberRoute><DashboardLayout><MemberTransactions /></DashboardLayout></MemberRoute>} />
          <Route path="/member/budget" element={<MemberRoute><DashboardLayout><MemberBudget /></DashboardLayout></MemberRoute>} />
          <Route path="/member/profile" element={<MemberRoute><DashboardLayout><MemberProfile /></DashboardLayout></MemberRoute>} />

          {/* Admin dashboard routes (protected) */}
          <Route path="/admin" element={<AdminRoute><AdminLayout><AdminDashboard /></AdminLayout></AdminRoute>} />
          <Route path="/admin/members" element={<AdminRoute><AdminLayout><AdminMembers /></AdminLayout></AdminRoute>} />
          <Route path="/admin/categories" element={<AdminRoute><AdminLayout><AdminCategories /></AdminLayout></AdminRoute>} />
          <Route path="/admin/transactions" element={<AdminRoute><AdminLayout><AdminTransactions /></AdminLayout></AdminRoute>} />
          <Route path="/admin/budget" element={<AdminRoute><AdminLayout><AdminBudget /></AdminLayout></AdminRoute>} />
          <Route path="/admin/settings" element={<AdminRoute><AdminLayout><AdminSettings /></AdminLayout></AdminRoute>} />

          {/* Catch all */}
          <Route path="*" element={<Navigate to="/member/auth/login" replace />} />
        </Routes>
      </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  )
}

export default App
