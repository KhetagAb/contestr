.PHONY: build test run-api run-bot docker-build docker-run lint clean

code-gen:
	./scripts/codegen.sh

build: code-gen
	go build -o bin/app cmd/app/main.go

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
