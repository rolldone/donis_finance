# donis_finance — Rencana Kerja

> Dibuat: 17 Juli 2026
> Prioritas: berdasarkan dampak vs effort

---

## ✅ Phase 1: UI Polish (Selesai — 17 Juli 2026)

### 1.1 Loading Skeleton ✅
- **Deskripsi:** Animasi placeholder saat fetching data (table, cards, dashboard widgets)
- **File affected:**
  - `sub_app/webapp/src/components/Skeleton.tsx` — komponen reusable (Text, Title, Avatar, Card, Table, BudgetCard, ChartCard, TxItem)
  - 8 pages di-update: Dashboard (member+admin), Transactions (member+admin), Budget (member+admin), Members, Categories, Settings
- **Acceptance criteria:**
  - ✅ Muncul saat `loading = true`
  - ✅ Hilang smooth setelah data tampil
  - ✅ Bentuk sesuai konteks: baris table, card, circle avatar

### 1.2 Responsive Mobile ✅
- **Deskripsi:** Padding responsif, table horizontal scroll
- **File affected:**
  - `DashboardLayout.tsx` — `p-6` → `p-4 sm:p-6`, `px-6` → `px-4 sm:px-6`
  - `AdminLayout.tsx` — same padding changes
  - Member/Admin Transactions — `overflow-x-auto` + `min-w-[600px]`
  - Member Dashboard — `overflow-x-auto` untuk recent transactions
- **Acceptance criteria:**
  - ✅ Sidebar sudah collapse via hamburger (dari awal)
  - ✅ Table horizontal scroll di mobile
  - ✅ Padding mengecil di mobile (`p-4` → `p-6` di `sm:`)

### 1.3 Dark/Light Toggle ✅
- **Deskripsi:** Switch theme di seluruh app
- **File affected:**
  - `sub_app/webapp/src/context/ThemeContext.tsx` — provider + toggle + localStorage + system preference
  - `sub_app/webapp/src/index.css` — CSS variables (`:root` + `.dark`) + `@custom-variant dark` + utility classes + global overrides
  - `App.tsx` — wrapper `<ThemeProvider>`
  - `DashboardLayout.tsx` — toggle button (sun/moon), CSS variables
  - `AdminLayout.tsx` — same toggle + variables
  - `AuthLayout.tsx`, `Login.tsx`, `Register.tsx`, `ForgotPassword.tsx` — CSS variables
  - `ProtectedRoute.tsx` — CSS variables
- **Acceptance criteria:**
  - ✅ Toggle di topbar (member + admin)
  - ✅ State persist ke localStorage
  - ✅ Default ikut system preference `prefers-color-scheme`

---

## ✅ Phase 2: Charts & Visualisasi (Selesai — 17 Juli 2026)

### 2.1 Charts — Pengeluaran per Kategori ✅
- **Deskripsi:** Grafik pie pengeluaran per kategori + bar chart income vs expense (Recharts)
- **File affected:**

  **Backend:**
  - `plugins/donisfinance/services/transactions.go` — tambah `GetMonthlySeries()` (income/expense per bulan), struct `MonthlySeriesPoint`
  - `plugins/donisfinance/handlers/transactions.go` — tambah handler `GetMonthlySeries`
  - `plugins/donisfinance/plugin.go` — register route `GET /transactions/monthly` (member + admin)

  **Frontend:**
  - `sub_app/webapp/` — install `recharts` (39 packages)
  - `sub_app/webapp/src/api/index.ts` — tambah types `CategoryBreakdown`, `SummaryWithBreakdown`, `MonthlySeriesPoint` + API functions `getTransactionSummaryWithBreakdown`, `getMonthlySeries`, dll
  - `sub_app/webapp/src/components/ExpenseChart.tsx` — pie chart (donut) pengeluaran per kategori dengan label persentase + legend + tooltip
  - `sub_app/webapp/src/components/IncomeVsExpenseChart.tsx` — bar chart pemasukan vs pengeluaran 6 bulan terakhir
  - `sub_app/webapp/src/pages/member/Dashboard.tsx` — integrasi kedua chart, fetch `getTransactionSummaryWithBreakdown` + `getMonthlySeries`
  - `sub_app/webapp/src/pages/admin/Dashboard.tsx` — integrasi kedua chart + fetch admin variants

- **Acceptance criteria:**
  - ✅ Pie chart: pengeluaran per kategori (member + admin)
  - ✅ Bar chart: income vs expense per bulan
  - ✅ Warna konsisten dengan kategori (pakai `category_color` dari database, fallback ke palet 20 warna)
  - ✅ Tooltip + legend
  - ✅ Label persentase di pie (sembunyi jika < 3%)
  - ✅ Zero-fill bulan tanpa transaksi  
  - ✅ Dark mode compatible (CSS variables)
  - ✅ Backend: 0 error, Frontend: TypeScript 0 error + build sukses (612 modules)

### 2.1 Charts - Pengeluaran per Kategori
- **Deskripsi:** Grafik pie/bar pengeluaran per kategori pakai Recharts
- **File affected:**
  - `sub_app/webapp/src/components/` → buat `ExpenseChart.tsx`, `IncomeVsExpenseChart.tsx`
  - `sub_app/webapp/src/pages/member/Overview.tsx` — tambah chart section
  - `sub_app/webapp/src/pages/admin/Dashboard.tsx` — tambah chart section
  - Backend: mungkin perlu endpoint aggregasi `GET /api/transactions/summary/by-category`
- **Acceptance criteria:**
  - Pie chart: pengeluaran per kategori (member + admin)
  - Bar chart: income vs expense per bulan
  - Warna konsisten dengan kategori

---

## Phase 3: Fitur Non-Kritis (Medium-High Effort)

### 3.1 Forgot Password (Real)
- **Deskripsi:** Kirim email reset password beneran (pakai Mailer yang sudah ada)
- **File affected:**
  - Backend: handler & service forgot/reset password
  - Frontend: halaman forgot password, reset password form
  - Email template: `templates/email/reset_password.html`
- **Acceptance criteria:**
  - User input email → dapat link reset via email
  - Link expire (misal 1 jam)
  - Bisa set password baru

### 3.2 Register Member Public
- **Deskripsi:** Endpoint register untuk public + approval flow oleh admin
- **File affected:**
  - Backend: handler register, service, model status `pending`
  - Frontend: halaman register real (bukan dummy)
  - Admin: halaman approval list
- **Acceptance criteria:**
  - Siapa aja bisa daftar
  - Status `pending` sampai di-approve admin
  - Admin bisa approve/reject dari dashboard

### 3.3 Scheduler Email Report
- **Deskripsi:** Cron job kirim laporan otomatis tiap tanggal 1
- **File affected:**
  - Backend: scheduler package / goroutine + cron expression
  - Service: generate report + send email ke semua member
  - Console command: test trigger manual
- **Acceptance criteria:**
  - Jalan otomatis tiap 1st day of month
  - Bisa di-trigger manual lewat CLI
  - Log sukses/gagal

---

## Phase 4: Enhancement (Low Effort)

### 4.1 Sorting Transaksi
- **Deskripsi:** Klik header kolom table untuk sort by amount, date, description
- **File affected:**
  - Backend: query param `sort_by`, `sort_order`
  - Frontend: table header click handler, indicator arrow
- **Acceptance criteria:**
  - Klik kolom → sort ascending, klik lagi → descending
  - Ada icon panah di header
  - Default: sort by date descending

### 4.2 Multi-language (i18n)
- **Deskripsi:** Dukungan multi bahasa (ID + EN)
- **File affected:**
  - `package.json` — tambah `react-i18next` / `i18next`
  - `src/i18n/` — folder locales
  - Semua komponen — ubah static text pake `t()`
- **Acceptance criteria:**
  - Bisa ganti bahasa dari dropdown
  - Default ID
  - Tidak ada hardcoded text (kecuali constants)

---

## Cara Kerja

1. **Satu task dikerjakan per sesi** — mulai dari Phase 1
2. Setiap task buat branch baru: `feat/{nama-fitur}`
3. Frontend + Backend dikerjakan bareng dalam satu task
4. Update `#file:.hermes.md` setelah selesai
