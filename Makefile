APP_NAME := app
GO := go

.PHONY: build test lint docker-build run gen migrate-up migrate-down

build:
	$(GO) build -o $(APP_NAME) .

test:
	$(GO) test ./... -race -count=1

lint:
	golangci-lint run ./...

docker-build:
	docker build -t grinex-rates-service .

run:
	docker-compose --profile app up --build

gen: gen-proto gen-sqlc

gen-proto:
	buf generate

gen-sqlc:
	sqlc generate

migrate-up:
	$(GO) run . migrate up

migrate-down:
	$(GO) run . migrate down