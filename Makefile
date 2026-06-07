ENV ?= dev

ENV_FILE := .env.$(ENV)
MIGRATIONS_DIR := ./internal/migrations
ifneq ($(wildcard $(ENV_FILE)),)
    include $(ENV_FILE)
    export $(shell sed 's/=.*//' $(ENV_FILE))
endif

.PHONY: status up down create

status:
		@echo "Checking goose migration status for [$(ENV)]"
		goose -dir $(MIGRATIONS_DIR) status

up:
		@echo "Running migrations for [$(ENV)]"
		goose -dir $(MIGRATIONS_DIR) up

down:
		@echo "Rolling back to migration for [$(ENV)]"
		goose -dir $(MIGRATIONS_DIR) down