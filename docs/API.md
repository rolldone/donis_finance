# API Documentation — Donis Finance

> **Base URL**: `{{BASE_URL}}/api`
> **Content-Type**: `application/json` (unless noted)

---

## Variable Reference

Dokumen ini menggunakan variable placeholder untuk menggantikan data sensitif. Ganti dengan nilai sesuai lingkungan kamu sebelum menggunakan contoh.

| Variable | Deskripsi | Contoh Nilai |
|----------|-----------|-------------|
| `{{BASE_URL}}` | Base URL server | `http://192.168.1.100:8200` |
| `{{SERVER_IP}}` | IP Address server | `192.168.1.100` |
| `{{ADMIN_USERNAME}}` | Username admin | `admin` |
| `{{ADMIN_PASSWORD}}` | Password admin | `pilih_password_kuat` |
| `{{ADMIN_EMAIL}}` | Email admin | `admin@example.com` |
| `{{MEMBER_USERNAME}}` | Username anggota | `nama_anggota` |
| `{{MEMBER_PASSWORD}}` | Password anggota | `pilih_password_kuat` |
| `{{MEMBER_EMAIL}}` | Email anggota | `anggota@example.com` |
| `{{MEMBER_NAME}}` | Nama lengkap anggota | `Nama Anggota` |
| `{{FROM_EMAIL}}` | Email pengirim (SMTP) | `noreply@example.com` |
| `{{SMTP_USER}}` | Username SMTP | `smtp_user` |
| `{{SMTP_PASS}}` | Password SMTP | `smtp_password` |
| `{{SMTP_HOST}}` | Host SMTP | `smtp.example.com` |

---

## Authentication

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Token is obtained from login endpoints. Default access token TTL is **15 minutes**.

---

## Table of Contents

- [Public Endpoints](#public-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Member Endpoints](#member-endpoints)
- [Common Response Formats](#common-response-formats)

---

## Public Endpoints

### POST `/api/admin/login`

Admin login. Returns JWT token.

**Request:**
```json
{
  "username": "admin",
  "password": "{{ADMIN_PASSWORD}}"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "username": "admin",
    "email": "{{ADMIN_EMAIL}}"
  }
}
```

---

### POST `/api/member/login`

Member login. Member must have `active` status.

**Request:**
```json
{
  "username": "donny",
  "password": "{{MEMBER_PASSWORD}}"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "username": "donny",
    "name": "Donny Rolanda"
  }
}
```

---

### POST `/api/member/auth/register`

Public member registration. Creates member with `pending` status.

**Request:**
```json
{
  "name": "Donny Rolanda",
  "username": "donny",
  "email": "{{MEMBER_EMAIL}}",
  "password": "{{MEMBER_PASSWORD}}"
}
```

**Response (201):**
```json
{
  "message": "Registration successful. Waiting for admin approval."
}
```

---

### POST `/api/member/auth/forgot-password`

Request password reset email. Requires valid member email.

**Request:**
```json
{
  "email": "{{MEMBER_EMAIL}}"
}
```

**Response (200):**
```json
{
  "message": "Password reset email sent"
}
```

---

### POST `/api/member/auth/reset-password`

Reset password using token from email.

**Request:**
```json
{
  "token": "reset-token-from-email",
  "password": "{{MEMBER_PASSWORD}}"
}
```

**Response (200):**
```json
{
  "message": "Password reset successful"
}
```

---

### GET `/api/admin/health`

Health check endpoint. No auth required.

**Response (200):**
```json
{
  "status": "ok",
  "plugin": "donisfinance"
}
```

---

## Admin Endpoints

All admin endpoints require `Authorization: Bearer <admin_token>` header.

---

### Profile

#### GET `/api/admin/profile`

Get current admin profile.

**Response (200):**
```json
{
  "id": "uuid",
  "username": "admin",
  "email": "{{ADMIN_EMAIL}}",
  "created_at": "2026-07-12T00:00:00Z"
}
```

#### PUT `/api/admin/profile`

Update admin profile.

**Request:**
```json
{
  "username": "admin",
  "email": "newemail@donis.finance"
}
```

#### PUT `/api/admin/password`

Change admin password.

**Request:**
```json
{
  "old_password": "{{ADMIN_PASSWORD}}",
  "new_password": "{{MEMBER_PASSWORD}}"
}
```

---

### SMTP Settings

#### GET `/api/admin/settings/smtp`

Get SMTP configuration (DB override + env fallback).

**Response (200):**
```json
{
  "smtp": {
    "host": "sandbox.smtp.mailtrap.io",
    "port": "2525",
    "from": "{{FROM_EMAIL}}",
    "sender_name": "donis_finance",
    "username": "",
    "password": ""
  },
  "env_smtp": {
    "host": "sandbox.smtp.mailtrap.io",
    "port": "2525",
    "from": "{{FROM_EMAIL}}",
    "sender_name": "donis_finance",
    "username": "your_user",
    "password": "your_pass"
  },
  "override": true
}
```

#### PUT `/api/admin/settings/smtp`

Save SMTP configuration. If all fields empty, deletes DB settings (reverts to .env).

**Request:**
```json
{
  "host": "smtp.gmail.com",
  "port": "587",
  "from": "{{FROM_EMAIL}}",
  "sender_name": "Donis Finance",
  "username": "{{SMTP_USER}}",
  "password": "app-password"
}
```

---

### Members

#### GET `/api/admin/members`

List all members for current admin.

**Response (200):**
```json
{
  "members": [
    {
      "id": "uuid",
      "name": "Donny Rolanda",
      "username": "donny",
      "email": "{{MEMBER_EMAIL}}",
      "status": "active",
      "created_at": "2026-07-12T00:00:00Z"
    }
  ]
}
```

#### POST `/api/admin/members`

Create a new member (auto-approved).

**Request:**
```json
{
  "name": "Donny Rolanda",
  "username": "donny",
  "email": "{{MEMBER_EMAIL}}",
  "password": "{{MEMBER_PASSWORD}}"
}
```

#### PUT `/api/admin/members/:id`

Update member name/username.

**Request:**
```json
{
  "name": "Donny Rolanda",
  "username": "donny_r"
}
```

#### DELETE `/api/admin/members/:id`

Delete a member and all associated data.

#### PATCH `/api/admin/members/:id/approve`

Approve a pending member registration.

#### PATCH `/api/admin/members/:id/reject`

Reject a pending member registration.

---

### Categories

#### GET `/api/admin/categories`

List all categories. Optional query: `?type=income` or `?type=expense`.

**Response (200):**
```json
{
  "categories": [
    {
      "id": "uuid",
      "name": "Gaji",
      "type": "income",
      "icon": "💰",
      "color": "#22c55e"
    }
  ]
}
```

#### POST `/api/admin/categories`

Create a category.

**Request:**
```json
{
  "name": "Gaji",
  "type": "income",
  "icon": "💰",
  "color": "#22c55e"
}
```

> **Duplicate Check:** Jika nama + type sudah ada, return `409 Conflict` dengan `existing_id`. Ini adalah soft check — jika kategori sudah dihapus, nama yang sama bisa dibuat ulang.

#### PUT `/api/admin/categories/:id`

Update a category.

#### DELETE `/api/admin/categories/:id`

Delete a category.

---

### Transactions

#### GET `/api/admin/transactions`

List transactions with filters.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `month` | int | 1-12 |
| `year` | int | e.g. 2026 |
| `type` | string | `income`, `expense`, `transfer` |
| `q` | string | Search in description |
| `member_id` | UUID | Filter by member |
| `sort_by` | string | `date`, `amount`, `created_at` |
| `sort_order` | string | `asc`, `desc` |
| `limit` | int | Max 100, default 20 |
| `offset` | int | Pagination offset |

**Response (200):**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "member_id": "uuid",
      "member_name": "Donny Rolanda",
      "account_id": "uuid",
      "account_name": "Cash",
      "category_id": "uuid",
      "category_name": "Gaji",
      "category_icon": "💰",
      "amount": 5000000,
      "type": "income",
      "description": "Gaji bulanan",
      "notes": "",
      "attachment_path": "",
      "date": "2026-07-15"
    }
  ],
  "total": 25
}
```

#### GET `/api/admin/transactions/summary`

Monthly income/expense summary by category.

**Query:** `?month=7&year=2026`

**Response (200):**
```json
{
  "summary": {
    "income": [
      { "category_id": "uuid", "category_name": "Gaji", "total": 5000000 }
    ],
    "expense": [
      { "category_id": "uuid", "category_name": "Makanan", "total": 1500000 }
    ]
  }
}
```

#### GET `/api/admin/transactions/monthly`

Monthly income/expense series for charts.

**Query:** `?months=6`

**Response (200):**
```json
{
  "monthly": [
    { "month": "2026-02", "income": 5000000, "expense": 3000000 },
    { "month": "2026-03", "income": 5000000, "expense": 2500000 }
  ]
}
```

#### GET `/api/admin/transactions/export`

Export transactions as CSV file.

**Query:** Same filters as list + `?member_id=uuid`

**Response:** CSV file download

#### PUT `/api/admin/transactions/:id`

Update a transaction.

**Request:**
```json
{
  "amount": 75000,
  "description": "Makan siang (updated)",
  "date": "2026-07-15"
}
```

#### DELETE `/api/admin/transactions/:id`

Delete a transaction (and its attachment file).

#### POST `/api/admin/transactions/:id/attachment`

Upload attachment. Max 10MB. Accepts: images, PDF, DOC, DOCX, XLS, XLSX.

**Request:** `multipart/form-data` with `file` field

**Response (200):**
```json
{
  "attachment_path": "/uploads/transactions/uuid/file.jpg"
}
```

#### GET `/api/admin/transactions/:id/attachment`

Download attachment file.

**Response:** Binary file download

---

### Accounts

#### GET `/api/admin/accounts`

List accounts. Optional: `?member_id=uuid`

**Response (200):**
```json
{
  "accounts": [
    {
      "id": "uuid",
      "member_id": "uuid",
      "name": "Cash",
      "type": "cash",
      "balance": 500000
    }
  ]
}
```

---

### Budgets

#### POST `/api/admin/budgets`

Set budget (upsert: creates or updates by member+category+month+year).

**Request:**
```json
{
  "member_id": "uuid",
  "category_id": "uuid",
  "month": 7,
  "year": 2026,
  "amount": 2000000
}
```

#### GET `/api/admin/budgets/status`

Get budget status for a month.

**Query:** `?month=7&year=2026&member_id=uuid`

**Response (200):**
```json
{
  "budgets": [
    {
      "id": "uuid",
      "category_id": "uuid",
      "category_name": "Makanan",
      "category_icon": "🍽️",
      "budget_amount": 2000000,
      "spent_amount": 1500000,
      "remaining": 500000,
      "percentage": 75.0,
      "is_over": false
    }
  ]
}
```

#### DELETE `/api/admin/budgets/:id`

Delete a budget entry.

---

## Member Endpoints

All member endpoints require `Authorization: Bearer <member_token>` header.

Member endpoints are similar to admin but scoped to the authenticated member's data only.

### Key Differences from Admin

| Admin | Member | Description |
|-------|--------|-------------|
| `POST /admin/members` | ❌ | Only admin can create members |
| `GET /admin/transactions?member_id=X` | `GET /member/transactions` | Auto-scoped to member |
| `POST /member/transactions` | ✅ | Members can create their own transactions |
| `POST /member/accounts` | ✅ | Members can create their own accounts |
| `POST`/`GET`/`DELETE /member/budgets` | ✅ | Members can manage their own budgets |

### Member-Only Endpoints

#### POST `/api/member/accounts`

Create an account.

**Request:**
```json
{
  "name": "My Cash",
  "type": "cash",
  "initial_balance": 1000000
}
```

#### PUT `/api/member/accounts/:id`

Update an account (partial update). Fully backward compatible — existing clients without `balance`/`balance_reason` work unchanged.

**Request:**
```json
{
  "name": "My Cash Updated",
  "type": "savings"
}
```

With balance adjustment (requires `balance_reason`):
```json
{
  "balance": 5000000,
  "balance_reason": "Reconcile bank statement 18 Juli"
}
```

| Field | Wajib | Deskripsi |
|-------|-------|-----------|
| `name` | ❌ | Nama baru akun |
| `type` | ❌ | Tipe baru (`cash`/`bank`/`e_wallet`/`savings`/`investment`) |
| `balance` | ❌ | Saldo baru langsung (reconcile). Dicatat di `balance_adjustments` audit trail |
| `balance_reason` | ❌ | **Wajib** jika `balance` diisi. Alasan perubahan saldo |

#### POST `/api/member/transactions`

Create a transaction.

**Request:**
```json
{
  "account_id": "uuid",
  "category_id": "uuid",
  "amount": 75000,
  "type": "expense",
  "description": "Makan siang",
  "notes": "Di warteg",
  "date": "2026-07-15"
}
```

For transfers:
```json
{
  "account_id": "uuid-source",
  "to_account_id": "uuid-destination",
  "amount": 500000,
  "type": "transfer",
  "description": "Transfer ke tabungan",
  "date": "2026-07-15"
}
```

#### POST `/api/member/budgets`

Set budget (upsert: creates or updates by category+month+year). Auto-scoped to authenticated member.

**Request:**
```json
{
  "category_id": "uuid",
  "month": 7,
  "year": 2026,
  "amount": 2000000
}
```

#### GET `/api/member/budgets/status`

Get budget status for a month. Auto-scoped to authenticated member.

**Query:** `?month=7&year=2026`

**Response (200):**
```json
{
  "budgets": [
    {
      "id": "uuid",
      "category_id": "uuid",
      "category_name": "Makanan",
      "category_icon": "🍽️",
      "budget_amount": 2000000,
      "spent_amount": 1500000,
      "remaining": 500000,
      "percentage": 75.0,
      "is_over": false
    }
  ]
}
```

#### DELETE `/api/member/budgets/:id`

Delete a budget entry. Only works if the budget belongs to the authenticated member.

---

## Common Response Formats

### Success
```json
{
  "message": "Success message"
}
```

### Error
```json
{
  "error": "Error description"
}
```

### Validation Error
```json
{
  "error": "Validation failed",
  "details": {
    "field_name": "error message"
  }
}
```

### Pagination (list endpoints)
```json
{
  "items": [...],
  "total": 100
}
```

---

## Account Types

| Value | Description |
|-------|-------------|
| `cash` | Uang tunai |
| `bank` | Rekening bank |
| `e_wallet` | Dompet digital (GoPay, OVO, etc.) |
| `savings` | Tabungan |
| `investment` | Investasi |

---

## Transaction Types

| Value | Description |
|-------|-------------|
| `income` | Pemasukan |
| `expense` | Pengeluaran |
| `transfer` | Transfer antar akun |

---

## Example: Full Flow

### 1. Login as Admin
```bash
TOKEN=$(curl -s -X POST {{BASE_URL}}/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"{{ADMIN_PASSWORD}}"}' | jq -r '.token')
```

### 2. Create a Member
```bash
curl -X POST {{BASE_URL}}/api/admin/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Donny Rolanda","username":"donny","email":"{{MEMBER_EMAIL}}","password":"{{MEMBER_PASSWORD}}"}'
```

### 3. Create Categories
```bash
curl -X POST {{BASE_URL}}/api/admin/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Gaji","type":"income","icon":"💰","color":"#22c55e"}'
```

### 4. Login as Member
```bash
MEMBER_TOKEN=$(curl -s -X POST {{BASE_URL}}/api/member/login \
  -H "Content-Type: application/json" \
  -d '{"username":"donny","password":"{{MEMBER_PASSWORD}}"}' | jq -r '.token')
```

### 5. Create Account & Transaction
```bash
# Create account
curl -X POST {{BASE_URL}}/api/member/accounts \
  -H "Authorization: Bearer $MEMBER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Cash","type":"cash","initial_balance":1000000}'

# Create transaction
curl -X POST {{BASE_URL}}/api/member/transactions \
  -H "Authorization: Bearer $MEMBER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"account_id":"uuid","category_id":"uuid","amount":75000,"type":"expense","description":"Makan siang","date":"2026-07-15"}'
```
