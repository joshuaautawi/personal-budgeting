# Personal Budgeting Backend (Go + Fiber)

This backend is designed to pair with the frontend in `../fe` and follows a layered architecture:

- `router` → `handlers` → `services` → `repositories`

## Run

From the repo root:

```bash
cd be
go mod tidy
go run ./cmd/api
```

## Run with Docker Compose (recommended)

From the repo root:

```bash
docker compose up --build
```

- To change Postgres credentials, create a `.env` file next to `docker-compose.yml` (see `docker-compose.env.example`).

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:7080`
- Postgres: `localhost:5410` (db: `budgeting`)

Environment variables:

- `PORT` (default `8080`)
- `DATABASE_URL` (optional) - when set, the backend uses Postgres via GORM.

### Connect to your local Postgres

You said: Postgres at `localhost:5432`, database `budgeting`.

Set one of these:

- `DATABASE_URL` (recommended):

```bash
export DATABASE_URL="postgres://YOURUSER:YOURPASS@localhost:5432/budgeting?sslmode=disable"
```

- Or PG* vars:

```bash
export PGHOST=localhost
export PGPORT=5432
export PGDATABASE=budgeting
export PGUSER=YOURUSER
export PGPASSWORD=YOURPASS
```

Schema:
- Run `be/migrations/001_init.sql` in your DB, or let the app auto-create tables on startup (GORM AutoMigrate).

## API (v1)

- `GET /api/v1/health`
- `GET /api/v1/state`
- `PUT /api/v1/state`
- `GET /api/v1/categories`
- `POST /api/v1/categories`
- `PATCH /api/v1/categories/:id`
- `DELETE /api/v1/categories/:id`
- `GET /api/v1/budgets`
- `PUT /api/v1/budgets` (upsert)
- `DELETE /api/v1/budgets/:id`
- `GET /api/v1/transactions`
- `POST /api/v1/transactions`
- `PATCH /api/v1/transactions/:id`
- `DELETE /api/v1/transactions/:id`

## Test

```bash
cd be
go test ./...
```


