# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run development server
make dev

# Build binaries to bin/
make build

# Database migrations
make migrate-up
make migrate-down
make migrate-fresh   # drop all tables + rerun migrations

# Seed database
make seed

# Generate Swagger/Scalar docs
make docs
```

Standard Go commands also apply:
```bash
go test ./...
go build ./...
```

API docs are served at `http://localhost:{APP_PORT}/docs` when the server is running.

## Architecture

This is a **Go REST API** for syncing and querying Meta Ads data, with a user/role/permission system.

**Stack:** Gin (HTTP), GORM + MySQL (persistence), Redis + Asynq (job queue), Centrifugo (WebSocket), JWT (auth).

### Entry Points

- `cmd/api/main.go` — HTTP API server
- `cmd/worker/main.go` — Asynq background job worker
- `cmd/seed/main.go` — Database seeder

### Layer Pattern

All domain modules follow **Repository → Service → Handler** with manual dependency injection wired in `routes/api.go`.

### Domain Modules

**`internal/core/`** — Multi-tenant user management:
- `auth/` — JWT login/logout, permission validation
- `user/`, `role/`, `permission/` — RBAC with fine-grained permissions using `domain.module.action` naming

**`internal/meta/`** — Meta Ads integration (each sub-module has its own repo/service/handler):
- `campaign/`, `adset/`, `ads/`, `insight/` — Sync from Meta Graph API and serve via REST
- `ad_account/` — Ad account info
- `sync_logs/` — Tracks sync batches and steps for auditing

**`internal/jobs/meta_ads_sync_job.go`** — Background ticker (every 15 min) that calls Meta Graph API and upserts campaigns, adsets, ads, and insights into the DB. Uses `date_preset=last_30d`.

### Key Packages

- `pkg/meta_client/` — HTTP client wrapping Meta Graph API v25.0 with pagination support
- `pkg/response/` — Standard response and pagination envelope structs
- `middleware/auth.go` — JWT auth middleware
- `config/` — Config loaders for env, DB, Redis, Meta API, Centrifugo

### Database

Migrations live in `database/migrations/` as plain SQL (`.up.sql` / `.down.sql`):
- `000001` — Core tables: `users`, `roles`, `permissions`, `user_roles`, `role_permissions`
- `000002` — Meta tables: `meta_campaigns`, `meta_ad_sets`, `meta_ads`, `meta_insights`
- `000003/000004` — Sync logging: `meta_sync_batches`, `meta_sync_steps`

`meta_insights` has a unique constraint on `(campaign_id, adset_id, ad_id, level, date_start, date_stop)` to support upserts. Soft deletes (`deleted_at`) are used on core tables.

### Environment

Copy `.env.example` to `.env`. Required variable groups:
- `APP_*` — Server config
- `DB_*` — MySQL connection
- `REDIS_*` — Redis connection
- `JWT_SECRET`, `JWT_EXPIRATION`
- `CENTRIFUGO_*` — WebSocket server
- `META_*` — Meta Graph API credentials (`META_ACCESS_TOKEN`, `META_AD_ACCOUNT_ID`, etc.)

Local services (MySQL, Redis, Centrifugo) can be started with `docker-compose up`.
