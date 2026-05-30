APP_NAME := petverse-api

.PHONY: run test vet fmt build tidy docker-up docker-down

run:
	go run ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

build:
	go build -trimpath -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

docker-up:
	docker compose -f deployments/docker/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker/docker-compose.yml down
