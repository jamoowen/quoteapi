.PHONY: test build start-server clean


TEST_PKGS = ./internal/utils ./internal/cache/

test:
	go test -v $(TEST_PKGS)
run:
	go run cmd/server/main.go

start-server:
	./bin/quoteapi

build:
	go build -o bin/quoteapi ./cmd/server/main.go

clean:
	rm -rf bin 
