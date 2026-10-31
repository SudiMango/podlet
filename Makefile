ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: api podlet sqlc migration-up migration-down migration-new db db-down seed DB_CHECK setup agent
.DEFAULT_GOAL := api

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5433
POSTGRES_ENV = PGHOST="$(POSTGRES_HOST)" PGPORT="$(POSTGRES_PORT)" PGUSER="$(POSTGRES_USER)" PGPASSWORD="$(POSTGRES_PASSWORD)" PGDATABASE="$(POSTGRES_DB)" PGSSLMODE=disable


# setup stuff

setup:
	go mod download

# run apps

api:
	go run cmd/api/main.go


podlet:
	go run cmd/cli/main.go $(args)


agent:
	go run cmd/agent/main.go $(args)


# database related commands

DB_CHECK:
ifndef POSTGRES_USER
	$(error Error: POSTGRES_USER is not set.)
endif
ifndef POSTGRES_PASSWORD
	$(error Error: POSTGRES_PASSWORD is not set.)
endif
ifndef POSTGRES_DB
	$(error Error: POSTGRES_DB is not set.)
endif


db:
	docker compose -f compose.db.yaml up -d


db-down:
	docker compose -f compose.db.yaml down


sqlc:
	go tool sqlc generate


migration-new:
ifndef name
	$(error Error: Please provide a migration name. Example: make migration-new name=add_billing_table)
endif
	go tool goose -dir internal/sql/migrations -s create $(name) sql


migration-up: DB_CHECK
	$(POSTGRES_ENV) go tool goose -dir internal/sql/migrations postgres "" up


migration-down: DB_CHECK
	$(POSTGRES_ENV) go tool goose -dir internal/sql/migrations postgres "" down


seed: DB_CHECK
ifndef file
	$(error Error: Please provide a seed filename. Example: make seed file=mock_seed.sql)
endif
	@test "$(notdir $(file))" = "$(file)" || \
		(echo "Error: file must be a filename from internal/sql/seeds"; exit 1)
	@test -f "internal/sql/seeds/$(file)" || \
		(echo "Error: seed file not found: internal/sql/seeds/$(file)"; exit 1)
	@$(POSTGRES_ENV) psql -v ON_ERROR_STOP=1 -f "internal/sql/seeds/$(file)"
