# AGENTS.md — AI Agent Context for Donis Finance

> **Fixed deployment**: `http://100.104.55.66:8200`
> **Admin panel**: `http://100.104.55.66:8200/admin`
> **Member panel**: `http://100.104.55.66:8200/member`
> **DB admin (pgAdmin)**: `http://100.104.55.66:8201`

---

## Quick Reference

### Tech Stack
| Layer | Technology | Version |
|-------|-----------|---------|
| Language | Go | 1.26 |
| HTTP Framework | Gin | — |
| ORM | GORM | (PostgreSQL driver) |
| Database | PostgreSQL | 16 |
| Cache | KeyDB (Redis-compatible) | alpine |
| Frontend | React + TypeScript + Vite | 19 / 7 / 8.1.4 |
| CSS | Tailwind CSS | 4 |
| Charts | Recharts | — |
| CLI | Cobra | — |

### Default Credentials
| Role | Username | Password | Email |
|------|----------|----------|-------|
| Admin | `admin` | `admin123` | `admin@donis.finance` |
| pgAdmin | `donny@donis.finance` | `donis_admin` | — |

### Docker Services
| Container | Port | Purpose |
|-----------|------|---------|
| `donis-finance-app-1` | **8200** | Go server (SPA + API) |
| `donis-finance-postgres-1` | 5432 (internal) | PostgreSQL 16 |
| `donis-finance-keydb-1` | 6379 (internal) | KeyDB cache |
| `donis-finance-pgadmin-1` | **8201** | pgAdmin web UI |

---

## Project Structure

```
donis_finance/
├── cmd/
│   ├── server/main.go          # HTTP server entrypoint
│   └── console/main.go         # CLI entrypoint (Cobra)
├── internal/
│   ├── app/bootstrap.go        # App initialization, plugin loading
│   ├── auth/                   # JWT sign/parse (HS256)
│   ├── db/                     # GORM setup, migrations, transactions
│   ├── storage/                # Local/S3 file storage
│   ├── keydb/                  # Redis-compatible client
│   ├── mail/                   # SMTP email (async queue)
│   ├── events/                 # Pub/sub event system
│   └── secrets/                # Env variable loading
├── plugins/
│   └── donisfinance/           # ★ ALL business logic lives here
│       ├── plugin.go           # Route registration + plugin init
│       ├── handlers/           # HTTP handlers (request/response)
│       ├── services/           # Business logic layer
│       ├── models/             # GORM models (DB schema)
│       ├── middleware/         # JWT auth middleware
│       ├── console/            # CLI commands
│       ├── migrations/         # SQL migration files
│       └── templates/          # Email HTML templates
├── sub_app/webapp/             # ★ React frontend
│   ├── src/
│   │   ├── App.tsx             # Route definitions
│   │   ├── api/index.ts        # API client (fetch wrapper)
│   │   ├── components/         # Layout components
│   │   ├── context/            # Auth, Theme, Language contexts
│   │   └── pages/              # Page components
│   └── dist/                   # Built assets (served by Go)
├── templates/index.tmpl        # HTML template (SPA catch-all)
├── docker-compose.yml          # All services
├── .env                        # Environment variables
└── build.sh                    # Build script
```

---

## Architecture

```
Browser ──► Go Server (:8200)
              ├── SPA (serves sub_app/webapp/dist/)
              ├── API (/api/*)
              │    ├── Auth middleware (JWT HS256)
              │    ├── Handlers → Services → GORM → PostgreSQL
              │    └── File storage (./storage/)
              └── Cron (monthly email reports)
```

**Key patterns:**
- Go serves BOTH the SPA static files AND the REST API
- SPA catch-all: any non-API, non-file route returns `index.tmpl` (React Router handles routing)
- Plugin architecture: all business code in `plugins/donisfinance/`, core framework in `internal/`
- Multi-tenant: each admin owns their members; members scoped by `admin_id`

---

## API Overview

### Authentication
- JWT Bearer token in `Authorization` header
- Access token TTL: 15 min (`JWT_ACCESS_EXP_SECONDS=900`)
- Refresh token TTL: 14 days
- Token claims: `sub` (user ID), `iat`, `exp`, `typ` (access/refresh/reset)

### Endpoints (all prefixed with `/api`)

#### Public
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/admin/login` | Admin login → JWT |
| `POST` | `/member/login` | Member login → JWT (must be `active`) |
| `POST` | `/member/auth/register` | Public registration (→ `pending`) |
| `POST` | `/member/auth/forgot-password` | Request reset email |
| `POST` | `/member/auth/reset-password` | Reset with token |
| `GET` | `/admin/health` | Health check |

#### Admin (JWT required, `user_type=admin`)
| Method | Path | Description |
|--------|------|-------------|
| `GET/PUT` | `/admin/profile` | Get/update profile |
| `PUT` | `/admin/password` | Change password |
| `GET/PUT` | `/admin/settings/smtp` | SMTP config (DB + env fallback) |
| `GET/POST` | `/admin/members` | List/create members |
| `PUT/DELETE` | `/admin/members/:id` | Update/delete member |
| `PATCH` | `/admin/members/:id/approve` | Approve pending member |
| `PATCH` | `/admin/members/:id/reject` | Reject pending member |
| `GET/POST` | `/admin/categories` | List/create categories |
| `PUT/DELETE` | `/admin/categories/:id` | Update/delete category |
| `GET` | `/admin/accounts` | List accounts |
| `GET` | `/admin/transactions` | List (filter: month, year, type, q, member_id, sort) |
| `GET` | `/admin/transactions/summary` | Monthly summary by category |
| `GET` | `/admin/transactions/monthly` | Monthly income/expense series |
| `GET` | `/admin/transactions/export` | CSV export |
| `PUT/DELETE` | `/admin/transactions/:id` | Update/delete transaction |
| `POST/GET` | `/admin/transactions/:id/attachment` | Upload/download attachment |
| `POST/GET` | `/admin/budgets` | Set/get budget status |
| `DELETE` | `/admin/budgets/:id` | Delete budget |

#### Member (JWT required, `user_type=member`)
| Method | Path | Description |
|--------|------|-------------|
| `GET/PUT` | `/member/profile` | Get/update profile |
| `PUT` | `/member/password` | Change password |
| `GET` | `/member/categories` | List categories |
| `GET/POST` | `/member/accounts` | List/create accounts |
| `GET/POST` | `/member/transactions` | List/create transactions |
| `GET` | `/member/transactions/summary` | Monthly summary |
| `GET` | `/member/transactions/monthly` | Monthly series |
| `PUT/DELETE` | `/member/transactions/:id` | Update/delete transaction |
| `POST/GET` | `/member/transactions/:id/attachment` | Upload/download |
| `POST/GET` | `/member/budgets` | Set/get budget status |

---

## Database Schema

7 tables in `donis_finance` database:

| Table | Purpose | Key Relations |
|-------|---------|---------------|
| `admins` | Admin users | — |
| `members` | Family members | FK → `admins.id` (CASCADE) |
| `categories` | Income/expense categories | — |
| `accounts` | Member wallets/banks | FK → `members.id` (CASCADE) |
| `transactions` | Financial transactions | FK → members, accounts, categories |
| `budgets` | Monthly budget limits | FK → members, categories. UNIQUE(member,category,month,year) |
| `settings` | Key-value config (SMTP) | PK = `key` |

All PKs are UUIDs. Amounts stored as `BIGINT` (no decimals — Rupiah integer).

---

## How to Build & Run

### Start everything
```bash
cd /home/donny/workspaces/donis_finance
docker compose up -d
docker exec -d donis-finance-app-1 sh -c "cd /app && GIN_MODE=release ./server"
```

### Rebuild frontend
```bash
cd sub_app/webapp && npm run build
```
**IMPORTANT**: After building, update asset hashes in `templates/index.tmpl`:
```html
<link rel="stylesheet" crossorigin href="/app/assets/index-NEWHASH.css">
<script type="module" src="/app/assets/index-NEWHASH.js"></script>
```

### Rebuild backend
```bash
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server"
```

### Full restart cycle
```bash
# Kill old processes
docker exec donis-finance-app-1 sh -c "kill -9 \$(pgrep -f './server') 2>/dev/null"

# Rebuild if needed
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server"

# Start fresh
docker exec -d donis-finance-app-1 sh -c "cd /app && GIN_MODE=release ./server"
```

### Access points
| URL | Purpose |
|-----|---------|
| `http://100.104.55.66:8200/admin` | Admin panel |
| `http://100.104.55.66:8200/admin/auth/login` | Admin login |
| `http://100.104.55.66:8200/member` | Member panel |
| `http://100.104.55.66:8200/member/auth/login` | Member login |
| `http://100.104.55.66:8200/api/admin/health` | Health check |
| `http://100.104.55.66:8201` | pgAdmin (DB admin) |

---

## Business Rules

1. **Multi-tenant by admin**: Each admin owns their members. Data is scoped.
2. **Member registration flow**: Register → `pending` → Admin approves → `active` → Can login
3. **Transaction types**: `income`, `expense`, `transfer` (between accounts)
4. **Balance tracking**: Transactions auto-update account balances on create/update/delete
5. **Budget warnings**: Creating an expense checks budget limits, returns warning if exceeded
6. **Seed data**: Default admin + 15 categories (5 income, 10 expense) auto-created
7. **File attachments**: Max 10MB, stored at `./storage/transactions/{tx_id}/{uuid}.ext`
8. **Monthly reports**: Cron on 1st of month at 08:00, emails all active members
9. **SMTP override**: DB settings override `.env` vars; clearing DB reverts to env

---

## Frontend Routes

| Route | Component | Auth |
|-------|-----------|------|
| `/` | → redirect to login | — |
| `/admin/auth/login` | AdminLogin | Guest |
| `/admin` | AdminDashboard | Admin |
| `/admin/members` | AdminMembers | Admin |
| `/admin/categories` | AdminCategories | Admin |
| `/admin/transactions` | AdminTransactions | Admin |
| `/admin/budget` | AdminBudget | Admin |
| `/admin/settings` | AdminSettings | Admin |
| `/member/auth/login` | Login | Guest |
| `/member/auth/register` | Register | Guest |
| `/member` | MemberDashboard | Member |
| `/member/transactions` | MemberTransactions | Member |
| `/member/budget` | MemberBudget | Member |
| `/member/profile` | MemberProfile | Member |

---

## Common Pitfalls

1. **Asset hashes**: After `npm run build`, the JS/CSS filenames change. You MUST update `templates/index.tmpl` with new hashes, then rebuild the Go binary.
2. **Zombie server processes**: Multiple `./server` processes can accumulate inside Docker. Use `kill -9 $(pgrep -f './server')` to clean up.
3. **Port binding**: After `docker restart`, the container's internal `./server` process is killed. You must manually start it with `docker exec -d`.
4. **SPA routing**: Go serves `index.tmpl` for all non-API routes. React Router handles client-side routing. Do NOT add server-side routes for SPA pages.
5. **Amounts**: All monetary amounts are `BIGINT` (Rupiah, no decimal). Display as `Rp X.XXX`.
