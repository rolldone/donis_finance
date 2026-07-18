# Panduan Pengguna — Donis Finance

> **Aplikasi manajemen keuangan keluarga**
> 🔗 {{BASE_URL}}

---

## 📖 Daftar Isi

1. [Mengenal Donis Finance](#mengenal-donis-finance)
2. [Cara Login](#cara-login)
3. [Panel Admin](#panel-admin)
4. [Panel Member (Anggota Keluarga)](#panel-member)
5. [Perintah CLI (Command Line)](#perintah-cli-command-line)
6. [Tips & Trik](#tips--trik)

---

## Mengenal Donis Finance

Donis Finance adalah aplikasi pencatatan keuangan untuk keluarga. Ada 2 jenis pengguna:

| Pengguna | Fungsi | Akses |
|----------|--------|-------|
| **Admin** | Mengelola anggota keluarga, kategori, melihat laporan | `/admin` |
| **Member (Anggota)** | Mencatat pemasukan, pengeluaran, transfer uang | `/member` |

### Alur Penggunaan

```
1. Admin login → membuat akun anggota (members)
2. Admin menyetujui pendaftaran anggota
3. Anggota login → mulai mencatat keuangan
4. Admin bisa melihat semua data anggota
5. Anggota hanya melihat data sendiri
```

---

## Cara Login

### Login Admin

1. Buka browser, kunjungi: **{{BASE_URL}}/admin**
2. Masukkan:
   - **Username**: `admin`
   - **Password**: `{{ADMIN_PASSWORD}}`
3. Klik **Login**

### Login Anggota (Member)

1. Buka browser, kunjungi: **{{BASE_URL}}/member**
2. Masukkan **Username** dan **Password** yang sudah dibuat admin
3. Klik **Login**

### Daftar Akun Baru (Anggota)

Jika kamu belum punya akun anggota:

1. Buka **{{BASE_URL}}/member/auth/register**
2. Isi:
   - **Nama lengkap**: Donny Rolanda
   - **Username**: donny
   - **Email**: {{MEMBER_EMAIL}}
   - **Password**: pilih password minimal 6 karakter
3. Klik **Register**
4. **Menunggu persetujuan admin** — status akan "pending" sampai admin menyetujui

---

## Panel Admin

Setelah login, kamu akan melihat **Dashboard** dengan ringkasan keuangan semua anggota.

### 📊 Dashboard

Halaman utama menampilkan:
- **Total Members** — jumlah anggota keluarga
- **Total Transactions** — total transaksi bulan ini
- **Income** — total pemasukan bulan ini
- **Expense** — total pengeluaran bulan ini
- **Grafik** — pie chart pengeluaran per kategori, grafik pemasukan vs pengeluaran
- **Recent Transactions** — transaksi terakhir

---

### 👥 Members (Kelola Anggota)

Halaman untuk mengelola anggota keluarga.

#### Menambah Anggota Baru

1. Klik tombol **+ Add Member**
2. Isi **Name**, **Username**, **Email**, **Password**
3. Klik **Save**
4. Anggota baru langsung bisa login (status: active)

#### Menyetujui Pendaftaran

Jika ada anggota yang mendaftar sendiri (status: pending):
1. Klik tombol **✓ (Approve)** di samping nama anggota
2. Anggota sekarang bisa login

#### Menolak Pendaftaran

1. Klik tombol **✕ (Reject)**
2. Anggota tidak bisa login

#### Menghapus Anggota

1. Klik tombol **🗑️ (Delete)**
2. **Peringatan**: Semua data transaksi anggota juga akan terhapus!

---

### 🏷️ Categories (Kategori)

Kategori membantu mengelompokkan pemasukan dan pengeluaran.

#### Kategori Default

**Pemasukan (Income):**
| Kategori | Icon |
|----------|------|
| Gaji | 💰 |
| Freelance | 💻 |
| Investasi | 📈 |
| Hadiah | 🎁 |
| Lainnya | 📥 |

**Pengeluaran (Expense):**
| Kategori | Icon |
|----------|------|
| Makanan | 🍽️ |
| Transportasi | 🚗 |
| Belanja | 🛒 |
| Tagihan | 📄 |
| Kesehatan | 🏥 |
| Pendidikan | 📚 |
| Hiburan | 🎬 |
| Rumah Tangga | 🏠 |
| Pakaian | 👔 |
| Lainnya | 📤 |

#### Menambah Kategori Baru

1. Klik **+ Add Category**
2. Isi:
   - **Name**: "Gaji Bonus"
   - **Type**: pilih `income` atau `expense`
   - **Icon**: pilih emoji (opsional)
   - **Color**: pilih warna (opsional)
3. Klik **Save**

#### Edit/Hapus Kategori

- Klik **✏️** untuk edit
- Klik **🗑️** untuk hapus

> ⚠️ **Catatan**: Kategori yang sudah digunakan di transaksi tidak bisa dihapus. Hapus transaksi yang terkait terlebih dahulu.

---

### 💳 Transactions (Transaksi)

Halaman untuk melihat dan mengelola SEMUA transaksi dari semua anggota.

#### Melihat Transaksi

- Gunakan **filter** di bagian atas:
  - **Bulan & Tahun**: pilih periode
  - **Type**: Income / Expense / Transfer
  - **Search**: cari di deskripsi
  - **Member**: filter per anggota
- Klik ikon **↑↓** untuk mengurutkan (tanggal, jumlah)

#### Mengedit Transaksi

1. Klik tombol **✏️** di samping transaksi
2. Ubah data yang diperlukan
3. Klik **Save**

#### Menghapus Transaksi

1. Klik tombol **🗑️**
2. Konfirmasi hapus

#### Export CSV

1. Klik tombol **Export CSV** di pojok kanan atas
2. File `.csv` akan terdownload
3. Bisa dibuka di Excel atau Google Sheets

---

### 🎯 Budget (Anggaran)

Atur batas pengeluaran per kategori per bulan.

#### Menambah Budget

1. Klik **+ Add Budget**
2. Pilih:
   - **Member**: anggota yang akan diatur
   - **Category**: kategori pengeluaran
   - **Month & Year**: periode
   - **Amount**: batas pengeluaran (contoh: 2000000 untuk Rp 2.000.000)
3. Klik **Save**

#### Melihat Status Budget

Di halaman Budget, kamu bisa melihat:
- **Budget Amount** — batas yang ditetapkan
- **Spent** — sudah terpakai berapa
- **Remaining** — sisa budget
- **Percentage** — persentase penggunaan
- **Indikator warna** — hijau (aman), kuning (hampir habis), merah (melebihi)

---

### ⚙️ Settings (Pengaturan)

#### Profil Admin

- Ubah **Username** dan **Email** admin
- Klik **Save Profile**

#### Ganti Password

1. Klik **Change** di section Security
2. Masukkan **password lama** dan **password baru**
3. Klik **Save**

#### SMTP Settings (Pengaturan Email)

Bagian ini mengatur pengiriman email laporan bulanan.

- Klik **View .env defaults** untuk melihat setting default
- Isi field SMTP (host, port, from, username, password)
- Klik **Save SMTP Settings**

> 💡 **Tip**: Jika semua field dikosongkan dan disimpan, setting akan kembali ke default dari server (.env).

---

## Panel Member (Anggota Keluarga)

### 📊 Dashboard

Halaman utama menampilkan ringkasan keuangan pribadi:
- Total saldo semua akun
- Pemasukan & pengeluaran bulan ini
- Grafik pie pengeluaran per kategori
- Grafik tren bulanan
- Transaksi terakhir

---

### 💳 Transactions (Transaksi Saya)

#### Mencatat Pemasukan

1. Klik **+ Add Transaction**
2. Isi:
   - **Account**: pilih akun (Cash, Bank, dll)
   - **Category**: pilih kategori income (Gaji, Freelance, dll)
   - **Amount**: jumlah (contoh: 5000000 untuk Rp 5.000.000)
   - **Type**: `Income`
   - **Description**: "Gaji bulanan Juli"
   - **Date**: pilih tanggal
3. Klik **Save**

#### Mencatat Pengeluaran

1. Klik **+ Add Transaction**
2. Isi:
   - **Account**: pilih akun yang dikeluarkan
   - **Category**: pilih kategori expense (Makanan, Transportasi, dll)
   - **Amount**: jumlah pengeluaran
   - **Type**: `Expense`
   - **Description**: "Makan siang di warteg"
   - **Date**: pilih tanggal
3. Klik **Save**

#### Transfer Antar Akun

Misalnya: transfer dari Cash ke Tabungan.

1. Klik **+ Add Transaction**
2. Isi:
   - **Account**: akun sumber (Cash)
   - **To Account**: akun tujuan (Tabungan)
   - **Amount**: jumlah transfer
   - **Type**: `Transfer`
   - **Description**: "Transfer ke tabungan"
   - **Date**: pilih tanggal
3. Klik **Save**

#### Menambah Lampiran (Attachment)

1. Buka detail transaksi (atau edit transaksi)
2. Klik **Upload Attachment**
3. Pilih file (gambar, PDF, DOC, XLS — maks 10MB)
4. File akan terupload dan bisa didownload kapan saja

---

### 🎯 Budget (Anggaran Saya)

Melihat status budget yang sudah diatur admin.

- Lihat berapa budget yang sudah terpakai
- Lihat sisa budget per kategori
- Peringatan jika pengeluaran mendekati/melebihi budget

---

### 👤 Profile (Profil Saya)

#### Update Profil

1. Buka halaman Profile
2. Ubah **Nama** atau **Username**
3. Klik **Save**

#### Ganti Password

1. Klik **Change Password**
2. Masukkan **password lama** dan **password baru**
3. Klik **Save**

---
# Perintah CLI (Command Line)

Selain melalui web, Donis Finance juga bisa dikendalikan lewat terminal/CLI. Cocok untuk admin server yang ingin mengelola data tanpa buka browser.

> **Cara pakai:** Jalankan dari dalam container Docker:
> ```bash
> docker exec -it donis-finance-app-1 sh
> cd /app
> ./console <perintah> [flags]
> ```

---

## 👤 Manajemen Pengguna

### Buat Admin Baru

```bash
./console donisfinance:create-admin \
  --username budi \
  --password {{MEMBER_PASSWORD}} \
  --email {{ADMIN_EMAIL}}
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-u, --username` | ✅ | Username admin |
| `-p, --password` | ✅ | Password admin |
| `-e, --email` | ❌ | Email admin |

### Buat Member Baru

```bash
./console donisfinance:create-member \
  --admin budi \
  --name "Istri Budi" \
  --username istri \
  --password {{MEMBER_PASSWORD}}
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-a, --admin` | ✅ | Username admin pemilik |
| `-n, --name` | ✅ | Nama tampilan member |
| `-u, --username` | ✅ | Username member |
| `-p, --password` | ✅ | Password member |

### Lihat Daftar Admin

```bash
./console donisfinance:list-admins
```

Contoh output:
```
ID                                     Username             Email
------------------------------------------------------------------
a1b2c3d4-...                           budi                 {{ADMIN_EMAIL}}
```

### Lihat Daftar Member

```bash
./console donisfinance:list-members
```

---

## 💰 Transaksi

### Tambah Transaksi

```bash
# Pemasukan
./console donisfinance:tx-add \
  --member istri \
  --type income \
  --amount 5000000 \
  --category "Gaji" \
  --account "BCA" \
  --desc "Gaji bulanan Juli" \
  --date 2026-07-01

# Pengeluaran
./console donisfinance:tx-add \
  --member istri \
  --type expense \
  --amount 250000 \
  --category "Makan" \
  --account "Cash" \
  --desc "Makan siang" \
  --date 2026-07-15
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-m, --member` | ✅ | Username member |
| `-a, --amount` | ✅ | Jumlah (Rupiah) |
| `-t, --type` | ✅ | `income`, `expense`, atau `transfer` |
| `-c, --category` | ❌ | Nama kategori |
| `-k, --account` | ❌ | Nama akun |
| `-d, --desc` | ❌ | Deskripsi |
| `-n, --notes` | ❌ | Catatan panjang |
| `-D, --date` | ❌ | Tanggal (YYYY-MM-DD, default hari ini) |

### Lihat Daftar Transaksi

```bash
# Semua transaksi bulan ini
./console donisfinance:tx-list --member istri

# Filter bulan tertentu
./console donisfinance:tx-list --member istri --month 7 --year 2026

# Hanya pengeluaran
./console donisfinance:tx-list --member istri --type expense --month 7 --year 2026

# Batasi jumlah hasil
./console donisfinance:tx-list --member istri --limit 10
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-m, --member` | ✅ | Username member |
| `-M, --month` | ❌ | Filter bulan (1-12) |
| `-Y, --year` | ❌ | Filter tahun |
| `-t, --type` | ❌ | Filter tipe (`income`/`expense`) |
| `-l, --limit` | ❌ | Jumlah max hasil (default: 20) |

### Hapus Transaksi

```bash
./console donisfinance:tx-delete --id <transaction-id>
```

### Transfer Antar Akun

```bash
./console donisfinance:tx-transfer \
  --member istri \
  --from "BCA" \
  --to "Cash" \
  --amount 500000 \
  --desc "Tarik ATM" \
  --date 2026-07-15
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-m, --member` | ✅ | Username member |
| `-f, --from` | ✅ | Akun sumber |
| `-t, --to` | ✅ | Akun tujuan |
| `-a, --amount` | ✅ | Jumlah |
| `-d, --desc` | ❌ | Deskripsi |
| `-D, --date` | ❌ | Tanggal |

### Ringkasan Bulanan

```bash
./console donisfinance:tx-summary --member istri --month 7 --year 2026
```

Contoh output:
```
📊 Summary Jul 2026 — Istri Budi
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💚 Income:  Rp5000000
❤️ Expense: Rp1250000
💰 Balance: Rp3750000
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Category                  Amount Count
─────────────────────────────────────
💚 Gaji                 Rp5000000   1x
❤️ Makan                Rp750000   3x
❤️ Transport             Rp500000   2x
```

### Export ke CSV

```bash
./console donisfinance:tx-export \
  --member istri \
  --month 7 \
  --year 2026 \
  --file laporan_juli.csv
```

File CSV akan berisi kolom: Date, Type, Category, Amount, Description, Notes.

---

## 🏦 Akun

### Buat Akun Baru

```bash
./console donisfinance:account-create \
  --member istri \
  --name "BCA" \
  --type bank \
  --balance 10000000
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-m, --member` | ✅ | Username member |
| `-n, --name` | ✅ | Nama akun |
| `-t, --type` | ❌ | Tipe: `cash`, `bank`, `e_wallet`, `savings`, `investment` (default: `cash`) |
| `-b, --balance` | ❌ | Saldo awal (default: 0) |

### Lihat Daftar Akun

```bash
./console donisfinance:account-list --member istri
```

Contoh output:
```
ID                                     Name                 Type           Balance
──────────────────────────────────────────────────────────────────────────
a1b2c3d4-...                           BCA                  bank        Rp10000000
e5f6g7h8-...                           Cash                 cash         Rp2500000
```

---

## 📊 Budget (Anggaran)

### Set Budget Bulanan

```bash
./console donisfinance:budget-set \
  --member istri \
  --category "Makan" \
  --month 7 \
  --year 2026 \
  --amount 1000000
```

| Flag | Wajib | Keterangan |
|------|-------|-----------|
| `-m, --member` | ✅ | Username member |
| `-M, --month` | ✅ | Bulan (1-12) |
| `-Y, --year` | ✅ | Tahun |
| `-a, --amount` | ✅ | Limit budget (Rupiah) |
| `-c, --category` | ❌ | Nama kategori (kosongkan = semua kategori) |

### Lihat Status Budget

```bash
./console donisfinance:budget-status --member istri --month 7 --year 2026
```

Contoh output:
```
📊 Budget Jul 2026 — Istri Budi
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Category              Budget       Spent    Remaining    Used
──────────────────────────────────────────────────────────────
Makan              Rp1000000    Rp750000    Rp250000   75% ███████░░░
Transport           Rp500000    Rp500000         Rp0  100% ██████████ 🔴 OVER!
Belanja             Rp800000    Rp200000    Rp600000   25% ██░░░░░░░░
```

### Cek Budget Sebelum Transaksi

```bash
./console donisfinance:budget-check \
  --member istri \
  --category "Makan" \
  --month 7 \
  --year 2026 \
  --amount 300000
```

Output:
- ✅ `Sisa budget: Rp250000` — Masih muat
- 🔴 `OVER BUDGET by Rp50000!` — Melebihi budget

---

## 📧 Laporan Email

### Kirim Laporan ke 1 Member

```bash
./console donisfinance:send-report \
  --member istri \
  --to {{MEMBER_EMAIL}} \
  --month 7 \
  --year 2026
```

### Kirim Laporan ke Semua Member

```bash
./console donisfinance:send-bulk-reports --month 7 --year 2026
```

| Flag | Keterangan |
|------|-----------|
| `-M, --month` | Bulan (default: bulan lalu) |
| `-Y, --year` | Tahun |
| `--dry-run` | Cek siapa saja yang akan dikirim, tanpa mengirim |

---

## 📋 Dashboard (Ringkasan Keuangan)

```bash
./console donisfinance:dashboard --member istri
```

Flag tambahan:
- `-M, --month` — Bulan (default: bulan ini)
- `-Y, --year` — Tahun

---

## 🗃️ Database Migration (Developer)

Perintah migrasi untuk pengembang:

```bash
# Lihat status migrasi
./console migrate list

# Jalankan migrasi baru
./console migrate up

# Rollback 1 migrasi
./console migrate down

# Rollback semua
./console migrate down-all

# Buat file migrasi baru
./console migrate make tambah_kolom_baru --plugin donisfinance

# Seed database
./console seed
```

| Flag | Keterangan |
|------|-----------|
| `--plugin` | Target plugin (`core`, `donisfinance`, atau `all`) |
| `--db` | Database target |

---

## 🎯 Contoh Skenario CLI Lengkap

Berikut contoh mengelola keuangan sepenuhnya lewat CLI:

```bash
# 1. Masuk ke container
docker exec -it donis-finance-app-1 sh
cd /app

# 2. Buat admin
./console donisfinance:create-admin \
  --username pak_budi --password rahasia --email {{ADMIN_EMAIL}}

# 3. Buat akun member
./console donisfinance:create-member \
  --admin pak_budi --name "Ibu Sari" \
  --username sari --password rahasia

# 4. Buat akun untuk member
./console donisfinance:account-create \
  --member sari --name "BCA" --type bank --balance 10000000
./console donisfinance:account-create \
  --member sari --name "Cash" --type cash --balance 2000000

# 5. Set budget
./console donisfinance:budget-set \
  --member sari --category "Makan" --month 7 --year 2026 --amount 1500000

# 6. Catat transaksi
./console donisfinance:tx-add \
  --member sari --type income --amount 5000000 \
  --category "Gaji" --account "BCA" --desc "Gaji Juli"
./console donisfinance:tx-add \
  --member sari --type expense --amount 350000 \
  --category "Makan" --account "Cash" --desc "Groceries"

# 7. Lihat ringkasan
./console donisfinance:tx-summary --member sari --month 7 --year 2026

# 8. Cek budget
./console donisfinance:budget-status --member sari --month 7 --year 2026

# 9. Export laporan
./console donisfinance:tx-export \
  --member sari --month 7 --year 2026 --file laporan_sari.csv

# 10. Kirim email laporan
./console donisfinance:send-report \
  --member sari --to {{MEMBER_EMAIL}} --month 7 --year 2026
```

---
## Tips & Trik

### 💡 Tips Menggunakan Aplikasi

1. **Selalu update saldo akun** — Saat mencatat pengeluaran, pilih akun yang benar agar saldo akurat

2. **Gunakan deskripsi yang jelas** — "Makan siang di warteg" lebih berguna dari "Makan"

3. **Atur budget di awal bulan** — Sebelum bulan dimulai, tetapkan budget untuk setiap kategori

4. **Lampirkan bukti transaksi** — Upload foto struk atau bukti transfer untuk referensi

5. **Rutin cek dashboard** — Dashboard memberikan gambaran cepat kondisi keuangan

### 🔐 Keamanan

1. **Ganti password default** — Setelah pertama login, ganti password `{{ADMIN_PASSWORD}}` ke yang lebih kuat
2. **Jangan share token** — Token login hanya berlaku 15 menit
3. **Review anggota** — Admin harus menyetujui setiap pendaftaran baru

### ⚠️ Hal yang Perlu Diperhatikan

1. **Anggota yang didaftarkan admin langsung aktif** — Tidak perlu menunggu persetujuan

2. **Pendaftaran mandiri perlu disetujui** — Jika anggota mendaftar sendiri, admin harus approve dulu

3. **Hapus transaksi hati-hati** — Menghapus transaksi akan mengubah saldo akun

4. **Kategori tidak bisa dihapus jika ada transaksi** — Hapus transaksi terkait dulu

5. **Budget hanya untuk pengeluaran** — Budget tidak berlaku untuk pemasukan

### 📱 Akses dari HP

Aplikasi ini **responsive** — bisa diakses dari HP. Cukup buka browser di HP dan kunjungi:

```
{{BASE_URL}}
```

- Admin: `/admin`
- Member: `/member`

### 🔄 Lupa Password

1. Buka halaman login member
2. Klik **"Forgot Password?"**
3. Masukkan email yang terdaftar
4. Cek email untuk link reset password
5. Klik link dan masukkan password baru

---

## Contoh Penggunaan

### Skenario: Keluarga Budi

1. **Budi** (admin) login ke `/admin` dengan `admin`/`{{ADMIN_PASSWORD}}`

2. Budi membuat akun untuk istrinya:
   - Nama: Ani
   - Username: ani
   - Password: ani123

3. Budi membuat akun untuk anaknya:
   - Nama: Rina
   - Username: rina
   - Password: rina123

4. **Ani** login ke `/member` → mencatat:
   - Pemasukan: Gaji Budi Rp 8.000.000
   - Pengeluaran: Belanja bulanan Rp 2.500.000
   - Pengeluaran: SPP anak Rp 500.000
   - Transfer: Rp 3.000.000 dari Cash ke Tabungan

5. **Rina** login ke `/member` → mencatat:
   - Pengeluaran: Jajan Rp 50.000
   - Pengeluaran: Fotocopy tugas Rp 25.000

6. **Budi** di `/admin` → melihat:
   - Dashboard: total income Rp 8.000.000, total expense Rp 3.075.000
   - Semua transaksi Ani dan Rina di halaman Transactions
   - Budget: apakah pengeluaran makanan sudah sesuai budget

---

## Bantuan

Jika mengalami masalah:

1. **Tidak bisa login** — Hubungi admin untuk memeriksa status akun
2. **Halaman error** — Coba refresh (Ctrl+R atau Cmd+R)
3. **Transaksi salah** — Edit atau hapus transaksi di halaman Transactions
4. **Lupa password** — Gunakan fitur "Forgot Password" di halaman login

---

*Donis Finance — Kelola Keuangan Keluarga dengan Mudah* 💰
