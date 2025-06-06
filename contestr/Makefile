code-gen:
	./scripts/codegen.sh

build: code-gen
	go build -o bin/app cmd/app/main.go

test:
	go test -v ./...

# Clean
clean:
	rm -rf bin/
