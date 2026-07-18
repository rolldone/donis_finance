# Donis Finance

> **Family financial management application — Go + React + PostgreSQL**

---

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Node.js 22+ (for frontend dev)
- Go 1.26+ (optional, for local dev outside Docker)

### Start Everything

```bash
# Clone and start
cd donis_finance
docker compose up -d

# Start the Go server inside the container
docker exec -d donis-finance-app-1 sh -c "cd /app && GIN_MODE=release ./server"
```

### Verify

```bash
curl -s http://localhost:8200/api/admin/health
```

### Access

| URL | Description |
|-----|-------------|
| `http://localhost:8200/admin` | Admin Panel |
| `http://localhost:8200/member` | Member Panel |
| `http://localhost:8201` | pgAdmin (Database) |

### Default Credentials

| Service | Username | Password |
|---------|----------|----------|
| Admin | `admin` | `admin123` |
| pgAdmin | `donny@donis.finance` | `donis_admin` |

> 📖 **End-user documentation**: [docs/USER_GUIDE.md](docs/USER_GUIDE.md) (GUI + CLI guide in Indonesian)
> 📡 **API documentation**: [docs/API.md](docs/API.md)

---

## Tech Stack

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
| Auth | JWT (HS256) | — |
| CLI | Cobra | — |
| Container | Docker + Docker Compose | — |

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Browser                                                 │
│  ├── Admin Panel  (/admin)                               │
│  └── Member Panel (/member)                              │
└─────────────────────┬───────────────────────────────────┘
                      │ HTTP
┌─────────────────────▼───────────────────────────────────┐
│  Go Server (:8200)                                       │
│  ├── SPA Static Files (sub_app/webapp/dist/)             │
│  ├── REST API (/api/*)                                   │
│  │   ├── Auth Middleware (JWT Bearer)                    │
│  │   ├── Handlers → Services → GORM → PostgreSQL        │
│  │   └── File Storage (./storage/)                       │
│  └── Cron Jobs (monthly email reports)                   │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│  PostgreSQL 16              KeyDB (Redis)                │
│  DB: donis_finance         Flash messages, sessions      │
└─────────────────────────────────────────────────────────┘
```

**Key patterns:**
- Go serves BOTH the SPA static files AND the REST API on the same port
- SPA catch-all: any non-API, non-file route returns `index.tmpl` (React Router handles client-side routing)
- Plugin architecture: all business code in `plugins/donisfinance/`, core framework in `internal/`
- Multi-tenant: each admin owns their members; members scoped by `admin_id`

---

## Project Structure

```
donis_finance/
├── cmd/
│   ├── server/main.go              # HTTP server entrypoint
│   └── console/main.go             # CLI tools (Cobra)
│
├── internal/                       # Core framework
│   ├── app/bootstrap.go            # App initialization, plugin loading
│   ├── auth/                       # JWT sign/parse (HS256)
│   ├── db/                         # GORM setup, migrations, transactions
│   ├── storage/                    # Local/S3 file storage
│   ├── keydb/                      # Redis-compatible client
│   ├── mail/                       # SMTP email (async queue)
│   ├── events/                     # Pub/sub event system
│   ├── console/                    # CLI root commands (Cobra)
│   └── secrets/                    # Env variable loading
│
├── plugins/
│   └── donisfinance/               # ★ ALL business logic lives here
│       ├── plugin.go               # Route registration + plugin init
│       ├── handlers/               # HTTP handlers (request/response)
│       ├── services/               # Business logic layer
│       ├── models/                 # GORM models (DB schema)
│       ├── middleware/             # JWT auth middleware
│       ├── console/                # CLI commands (Cobra)
│       ├── migrations/postgres/    # SQL migration files
│       └── templates/email/        # Email HTML templates
│
├── sub_app/webapp/                 # ★ React frontend
│   ├── src/
│   │   ├── App.tsx                 # Route definitions
│   │   ├── api/index.ts            # API client (fetch wrapper)
│   │   ├── components/             # Layout components
│   │   ├── context/                # Auth, Theme, Language contexts
│   │   └── pages/                  # Page components
│   └── dist/                       # Built assets (served by Go)
│
├── templates/index.tmpl            # SPA HTML template (with asset hashes)
├── docker-compose.yml              # All services
├── .env                            # Environment variables
└── AGENTS.md                       # AI agent context document
```

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

## Development Guide

### Build Frontend

```bash
cd sub_app/webapp && npm run build
```

### ⚠️ After Frontend Build — MUST Update Asset Hashes

Vite generates new filenames on each build. You **must**:

1. Note the new filenames from build output (e.g., `index-XXXX.js`, `index-XXXX.css`)
2. Update `templates/index.tmpl` with the new hashes
3. Rebuild the Go binary (it embeds the template)

```bash
# Example: update templates/index.tmpl hashes, then:
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server"
```

### Build Backend

```bash
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server"
```

### Run Migrations

```bash
docker exec donis-finance-app-1 sh -c "cd /app && ./console migrate up"
```

### Full Restart Cycle

```bash
# Kill stale processes (zombie servers accumulate!)
docker exec donis-finance-app-1 sh -c "kill -9 \$(pgrep -f './server') 2>/dev/null"

# Rebuild
docker exec donis-finance-app-1 sh -c "cd /app && go build -o /app/server ./cmd/server"

# Start fresh
docker exec -d donis-finance-app-1 sh -c "cd /app && GIN_MODE=release ./server"
```

---

## CLI Commands

The project includes a CLI tool for admin operations:

```bash
docker exec -it donis-finance-app-1 sh
cd /app
./console <command> [flags]
```

### Available Commands

| Command | Description |
|---------|-------------|
| `donisfinance:create-admin` | Create a new admin user |
| `donisfinance:create-member` | Create a member under an admin |
| `donisfinance:list-admins` | List all admin users |
| `donisfinance:list-members` | List all members |
| `donisfinance:tx-add` | Add a transaction |
| `donisfinance:tx-list` | List transactions |
| `donisfinance:tx-edit` | Edit a transaction (partial update) |
| `donisfinance:tx-delete` | Delete a transaction |
| `donisfinance:tx-transfer` | Transfer between accounts |
| `donisfinance:tx-summary` | Monthly income/expense summary |
| `donisfinance:tx-export` | Export transactions to CSV |
| `donisfinance:tx-import` | Import from CSV (BCA/BLU format) |
| `donisfinance:account-create` | Create an account |
| `donisfinance:account-list` | List accounts |
| `donisfinance:account-update` | Update account name/type/balance |
| `donisfinance:account-adjust` | Adjust account balance (with audit) |
| `donisfinance:budget-set` | Set monthly budget |
| `donisfinance:budget-status` | Show budget vs actual |
| `donisfinance:budget-check` | Preview if transaction fits budget |
| `donisfinance:dashboard` | Show financial dashboard |
| `donisfinance:send-report` | Send report to one member |
| `donisfinance:send-bulk-reports` | Send report to all members |
| `migrate up/down/list` | Database migrations |
| `seed` | Seed database |

> 📖 Full CLI reference with examples: [docs/USER_GUIDE.md](docs/USER_GUIDE.md#perintah-cli-command-line)

---

## Deployment

### Docker Services

| Container | Port | Purpose |
|-----------|------|---------|
| `donis-finance-app-1` | **8200** | Go server (SPA + API) |
| `donis-finance-postgres-1` | 5432 (internal) | PostgreSQL 16 |
| `donis-finance-keydb-1` | 6379 (internal) | KeyDB cache |
| `donis-finance-pgadmin-1` | **8201** | pgAdmin web UI |

### Environment Variables

Key variables in `.env`:

| Variable | Example | Description |
|----------|---------|-------------|
| `APP_PORT` | `8200` | Server port |
| `DB_HOST` | `postgres` | Database host (Docker service name) |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `donis_finance` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `secret` | Database password |
| `AUTH_JWT_SECRET` | `random-string` | JWT signing secret |
| `JWT_ACCESS_EXP_SECONDS` | `900` | Access token TTL (15 min) |
| `JWT_REFRESH_EXP_SECONDS` | `1209600` | Refresh token TTL (14 days) |
| `SMTP_HOST` | `smtp.example.com` | Email server |
| `SMTP_PORT` | `587` | Email port |
| `STORAGE_DRIVER` | `local` | `local` or `s3` |
| `STORAGE_ROOT` | `./storage` | Local storage path |

---

## API Overview

Base URL: `http://localhost:8200/api`

### Authentication

```bash
# Login (admin)
curl -X POST http://localhost:8200/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Use token
curl http://localhost:8200/api/admin/profile \
  -H "Authorization: Bearer <token>"
```

### Endpoints Summary

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/admin/login` | — | Admin login → JWT |
| `POST` | `/member/login` | — | Member login → JWT |
| `POST` | `/member/auth/register` | — | Public registration |
| `GET/PUT` | `/admin/profile` | Admin | Admin profile |
| `GET/POST` | `/admin/members` | Admin | Manage members |
| `GET/POST` | `/admin/categories` | Admin | Manage categories |
| `GET` | `/admin/transactions` | Admin | List all transactions |
| `GET` | `/admin/transactions/summary` | Admin | Monthly summary |
| `GET/POST` | `/admin/budgets` | Admin | Manage budgets |
| `GET/PUT` | `/member/profile` | Member | Member profile |
| `GET/POST` | `/member/transactions` | Member | Own transactions |
| `GET` | `/member/transactions/summary` | Member | Monthly summary |
| `POST/GET` | `/member/budgets` | Member | Budget status |

> 📖 Full API docs with request/response examples: [docs/API.md](docs/API.md)

---

## Common Pitfalls

1. **Asset hashes**: After `npm run build`, JS/CSS filenames change. You MUST update `templates/index.tmpl` with new hashes, then rebuild the Go binary.

2. **Zombie server processes**: Multiple `./server` processes can accumulate inside Docker. Use `kill -9 $(pgrep -f './server')` to clean up.

3. **Port binding after restart**: After `docker restart`, the container's internal `./server` process is killed. You must manually start it with `docker exec -d`.

4. **SPA routing**: Go serves `index.tmpl` for all non-API routes. React Router handles client-side routing. Do NOT add server-side routes for SPA pages.

5. **Amounts**: All monetary amounts are `BIGINT` (Rupiah, no decimal). Display as `Rp X.XXX`.

---

## Documentation

| File | Target | Language |
|------|--------|----------|
| [README.md](README.md) | Developers | English |
| [AGENTS.md](AGENTS.md) | AI agents (Copilot) | English |
| [docs/API.md](docs/API.md) | API consumers | English |
| [docs/USER_GUIDE.md](docs/USER_GUIDE.md) | End users + CLI | Indonesian |

---

## License

Private — Donis Finance
