.PHONY: build test run-api run-bot docker-build docker-run lint clean

build-api:
	go build -o bin/api cmd/api/main.go

build-bot:
	go build -o bin/bot cmd/bot/main.go

build: build-api build-bot

test:
	go test -v ./...

run-api:
	go run cmd/api/main.go

run-bot:
	go run cmd/bot/main.go

# Docker
docker-build:
	docker-compose -f deployments/docker/docker-compose.yml build

docker-run:
	docker-compose -f deployments/docker/docker-compose.yml up

# Clean
clean:
	rm -rf bin/
