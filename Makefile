APP=api-quest
BIN=./bin/$(APP)
MAIN=./cmd/api/main.go
SEED=./cmd/seed/main.go

.PHONY: build run seed docker-up docker-down docker-build logs clean tidy

## Build the binary
build:
	go build -o $(BIN) $(MAIN)

## Run the API locally (requires MongoDB running)
run:
	go run $(MAIN)

## Run the seeder locally (requires MongoDB running)
seed:
	go run $(SEED)

## Start all services (API + MongoDB) via Docker Compose
docker-up:
	docker compose up --build -d

## Start only MongoDB in the background
docker-mongo:
	docker compose up mongo -d

## Stop all Docker Compose services
docker-down:
	docker compose down

## Build the Docker image only
docker-build:
	docker build -t $(APP) .

## Tail logs from the API container
logs:
	docker compose logs -f api

## Download and tidy Go modules
tidy:
	go mod tidy

## Remove compiled binary
clean:
	rm -rf ./bin
