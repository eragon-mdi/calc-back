include .env
export

CONTAINER_DB_NAME=my-postgres

# PROD build
start-quiet:
	docker compose build && \
	docker compose up -d && \
	clear && \
	docker ps -a

restart-quiet: down start-quiet

down:
	docker compose down


# Need new migrate files
migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Error: укажи имя миграции через 'name=...'" && exit 1; \
	fi
	migrate create -ext sql -dir ./migrations $(name)
# p.s. 
#curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
#sudo mv migrate /usr/local/bin


# DEV fast build
dev-up: run-postrges wait migrate-up run-app

wait:
	@sleep 2

dev-down: rm-postgres

run-postrges:
	docker run \
		--name $(CONTAINER_DB_NAME) \
		-e POSTGRES_USER=$(STORAGE_USER) \
		-e POSTGRES_PASSWORD=$(STORAGE_PASS) \
		-e POSTGRES_DB=$(STORAGE_NAME) \
		-p $(STORAGE_PORT):5432 \
		-d \
		postgres:latest

rm-postgres:
	docker rm -f $(CONTAINER_DB_NAME)

migrate-up:
	migrate -path ./migrations -database "postgres://$(STORAGE_USER):$(STORAGE_PASS)@$(STORAGE_HOST):$(STORAGE_PORT)/$(STORAGE_NAME)?sslmode=$(STORAGE_SSLM)" up

run-app:
	go run -tags=dev cmd/calc-app/main.go

clear-port: 					# if don't correct close app
	kill -9 $(lsof -ti :$(SERVER_PORT))


# DEV tools
lint:
	golangci-lint run ./cmd/... ./internal/... ./pkg/...

swag:
	swag init -g ./cmd/calc-app/main.go

gen-mocks:
	go generate ./internal/...

gen-base-tests-transport:
	gotests -w -all ./internal/transport/http/rest/calculation.go

gen-base-tests-service:
	gotests -w -all ./internal/service/calculation.go

go-tests-coverage:
	mkdir -p ./tmp
	go test \
		-short \
		-count=1 \
		-race \
		-coverprofile=./tmp/coverage.out \
		./internal/service/... \
		./internal/transport/http/rest/...
	go tool cover \
		-html=./tmp/coverage.out \
		-o ./tmp/coverage.html
	firefox tmp/coverage.html