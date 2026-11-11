# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git

Never commit or push on your own; the user handles all git commits and pushes. Only do so when explicitly told to.

## Overview

Podlet is a self-hosted deploy platform written in Go (module `github.com/SudiMango/podlet`). Three binaries share one module:

- `cmd/api` — HTTP API (chi router, uber/fx DI, Postgres via pgx + sqlc, auth via `limen` with GitHub OAuth/device flow).
- `cmd/cli` — user-facing `podlet` CLI (cobra, in `internal/cli`).
- `cmd/agent` — daemon installed on user servers (`internal/agent`); heartbeats to the API, claims and executes jobs.

## Commands

- `make setup` — `go mod download`
- `make api` — run the API (loads `.env`; see `.env.example`)
- `make podlet args="..."` — run the CLI
- `make agent args="..."` — build and run the agent with sudo (needs a real server environment)
- `make db` / `make db-down` — local Postgres via `compose.db.yaml` (port 5433)
- `make migration-new name=x`, `make migration-up`, `make migration-down` — goose migrations in `internal/sql/migrations`
- `make sqlc` — regenerate `internal/sql/database` from `internal/sql/queries` (generated code is committed; never hand-edit)
- `make seed file=x.sql` — run a seed from `internal/sql/seeds`

CI (`.github/workflows/ci.yml`) runs: `gofmt -l .` must be empty, `go tool sqlc compile`, `go test ./...`, `go vet ./...`, `bash -n` on `internal/scripts/*.sh`, and a Docker build of the API. Single test: `go test ./internal/path/... -run TestName`. CD builds/pushes the API image to GHCR on main.

## Architecture

**API layering** (`internal/api`): `handler` → `service` → `repository` → `internal/sql/database` (sqlc). Each package has a `module.go` with an fx module; `server/server.go` wires the router. `dto` holds API-side types; `cloudflare` and `github` are external clients. `internal/contract` holds request/response types and enums **shared between API, CLI and agent** — change wire formats there.

**Agent ↔ API job system** (the core flow):
1. The API creates jobs (`internal/api/service/jobs.go`, `repository/jobs.go`); the agent polls `GET /agent/jobs/claim` (`internal/agent/jobs.go`) and runs claimed jobs one at a time (`busy` flag).
2. Single jobs are `JobExecutor`s registered in `JobRegistry` keyed by `dto.JobType` (`internal/agent/dto/enums.go`). Multi-step operations are `JobChain`s in `ChainRegistry` keyed by `contract.JobType` (e.g. `DeployNewRepositoryChain` in `chain_deploy_new_repo.go`: clone → build → create DNS record → remove previous deploy), which call executors from `a.executors` sequentially.
3. Executors live under `internal/agent/jobs/`: `system/initialize` and `system/uninstall` (Docker, podlet network, Cloudflare tunnel/container, Railpack, BuildKit) and `runtime` (clone/build/DNS/remove-prev-deploy for deployments).

Adding a job type usually touches: `contract` enum (if API-visible), `internal/agent/dto/enums.go`, the executor, `JobRegistry`/`ChainRegistry`, and the API service. Job types are stored as text in the DB (migration 00009), so no DB enum change is needed.

**Deploy mechanics**: apps are built with Railpack on BuildKit and run as Docker containers on a `podlet` network, exposed through a Cloudflare tunnel with DNS records created per app (`CLOUDFLARE_*` env vars).

**Install/uninstall**: `internal/scripts/*.sh` are the CLI/agent installers (embedded/served for download). `stuff/` has manual release-testing notes and check scripts (`manual-release-testing.md` references `docs/build-test-artifacts.sh` but the script lives at `stuff/build-test-artifacts.sh`).
