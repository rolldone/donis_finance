import { useState } from 'react'
import { Link } from 'react-router-dom'
import AuthLayout from '../components/AuthLayout'

const DUMMY_USER = {
  name: 'Donny',
  email: 'donny@donis.finance',
}

export default function Profile() {
  const [form, setForm] = useState({ name: DUMMY_USER.name, email: DUMMY_USER.email })
  const [isEditing, setIsEditing] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Update profile:', form)
    setIsEditing(false)
    alert('Profile updated (dummy)')
  }

  return (
    <AuthLayout>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Profil Saya</h2>
          <p className="text-gray-400 text-sm">Kelola informasi akun Anda</p>
        </div>
        <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center text-gray-600 font-semibold text-sm">
          {DUMMY_USER.name.charAt(0)}
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Name */}
        <div>
          <label className="block text-sm font-medium text-gray-600 mb-1.5">Nama</label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            disabled={!isEditing}
            className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>

        {/* Email */}
        <div>
          <label className="block text-sm font-medium text-gray-600 mb-1.5">Email</label>
          <input
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            disabled={!isEditing}
            className="w-full px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-300 transition text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>

        {/* Buttons */}
        {isEditing ? (
          <div className="flex gap-3">
            <button
              type="button"
              onClick={() => { setIsEditing(false); setForm({ name: DUMMY_USER.name, email: DUMMY_USER.email }) }}
              className="flex-1 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition"
            >
              Batal
            </button>
            <button
              type="submit"
              className="flex-1 py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-lg transition active:scale-[0.98]"
            >
              Simpan
            </button>
          </div>
        ) : (
          <button
            type="button"
            onClick={() => setIsEditing(true)}
            className="w-full py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition"
          >
            ✏️ Edit Profil
          </button>
        )}
      </form>

      {/* Divider */}
      <div className="border-t border-gray-100 my-6" />

      {/* Actions */}
      <div className="space-y-2">
        <Link
          to="/member/auth/forgot-password"
          className="flex items-center gap-3 px-3 py-2.5 bg-gray-50 hover:bg-gray-100 rounded-lg transition text-gray-600 text-sm"
        >
          🔒 <span>Ganti Password</span>
        </Link>
        <Link
          to="/member/auth/login"
          className="flex items-center gap-3 px-3 py-2.5 bg-red-50 hover:bg-red-100 rounded-lg transition text-red-600 text-sm"
        >
          🚪 <span>Keluar</span>
        </Link>
      </div>
    </AuthLayout>
  )
}
